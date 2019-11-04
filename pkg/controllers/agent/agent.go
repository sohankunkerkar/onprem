/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package agent

import (
	"context"
	"io/ioutil"
	"time"

	v1alpha1 "github.com/font/onprem/api/v1alpha1"
	"github.com/go-logr/logr"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	HeartBeatDelay = time.Second * 5
	kubeAPIQPS     = 20.0
	kubeAPIBurst   = 30

	defaultConfigMapName      = "hub-cluster"
	defaultSecretName         = "hub-cluster"
	apiEndpointKey            = "server"
	joinedClusterNameKey      = "joinClusterName"
	joinedClusterNamespaceKey = "joinClusterNamespace"
	caBundleKey               = "caBundle"
	tokenKey                  = "token"
)

type JoinedClusterCoordinates struct {
	APIEndpoint string
	Name        string
	Namespace   string
	Token       []byte
	CABundle    []byte
}

func GetJoinedClusterCoordinates(spokeClient client.Client) (*JoinedClusterCoordinates, error) {
	namespaceBytes, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get my namespace")
	}
	namespace := string(namespaceBytes)

	// Extract contents from ConfigMap.
	namespacedName := client.ObjectKey{Namespace: namespace, Name: defaultConfigMapName}
	configMap := &corev1.ConfigMap{}
	err = spokeClient.Get(context.TODO(), namespacedName, configMap)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get configmap %+v", namespacedName)
	}

	apiEndpoint, found := configMap.Data[apiEndpointKey]
	if !found || len(apiEndpoint) == 0 {
		return nil, errors.Errorf("The configmap for the hub cluster is missing a non-empty value for %q", apiEndpointKey)
	}

	joinedClusterName, found := configMap.Data[joinedClusterNameKey]
	if !found || len(joinedClusterName) == 0 {
		return nil, errors.Errorf("The configmap for the hub cluster is missing a non-empty value for %q", joinedClusterNamespaceKey)
	}

	joinedClusterNamespace, found := configMap.Data[joinedClusterNamespaceKey]
	if !found || len(joinedClusterNamespace) == 0 {
		return nil, errors.Errorf("The configmap for the hub cluster is missing a non-empty value for %q", joinedClusterNamespaceKey)
	}

	// Extract contents from Secret.
	namespacedName.Name = defaultSecretName
	secret := &corev1.Secret{}
	err = spokeClient.Get(context.TODO(), namespacedName, secret)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get secret %+v", namespacedName)
	}

	token, found := secret.Data[tokenKey]
	if !found || len(token) == 0 {
		return nil, errors.Errorf("The secret for the hub cluster is missing a non-empty value for %q", tokenKey)
	}

	caBundle, found := secret.Data[caBundleKey]
	if !found || len(caBundle) == 0 {
		return nil, errors.Errorf("The secret for the hub cluster is missing a non-empty value for %q", caBundleKey)
	}

	return &JoinedClusterCoordinates{
		APIEndpoint: apiEndpoint,
		Name:        joinedClusterName,
		Namespace:   joinedClusterNamespace,
		Token:       token,
		CABundle:    caBundle,
	}, nil
}

func BuildHubClusterConfig(spokeClient client.Client, jcc *JoinedClusterCoordinates) (*rest.Config, error) {
	// Build config using contents extracted from ConfigMap and Secret.
	clusterConfig, err := clientcmd.BuildConfigFromFlags(jcc.APIEndpoint, "")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to build config using API endpoint")
	}

	clusterConfig.CAData = jcc.CABundle
	clusterConfig.BearerToken = string(jcc.Token)
	clusterConfig.QPS = float32(kubeAPIQPS)
	clusterConfig.Burst = kubeAPIBurst

	return clusterConfig, nil
}

// JoinedClusterReconciler reconciles a JoinedCluster object
type JoinedClusterReconciler struct {
	HubClient       client.Client
	SpokeClient     client.Client
	DiscoveryClient *discovery.DiscoveryClient
	Log             logr.Logger
}

func (r *JoinedClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("joinedcluster", req.NamespacedName)

	jcc, err := GetJoinedClusterCoordinates(r.SpokeClient)
	if err != nil {
		log.Error(err, "unable to get joinedcluster coordinates")
		return ctrl.Result{}, err
	}

	// Ignore JoinedClusters that do not belong to us. Only check name as we're
	// only watching one namespace.
	if jcc.Name != req.NamespacedName.Name {
		return ctrl.Result{}, nil
	}

	log.V(2).Info("reconciling")

	var joinedCluster v1alpha1.JoinedCluster
	if err := r.HubClient.Get(ctx, req.NamespacedName, &joinedCluster); err != nil {
		log.Error(err, "unable to get joinedcluster")
		return ctrl.Result{}, err
	}

	reason := "AgentHeartBeat"
	message := "Spoke agent successfully connected to hub"
	agentConnected := v1alpha1.JoinedClusterConditions{
		Type:    v1alpha1.ConditionTypeAgentConnected,
		Status:  v1alpha1.ConditionTrue,
		Reason:  &reason,
		Message: &message,
	}

	currentTime := metav1.Now()
	statusTransitioned := joinedCluster.Status.Conditions != nil && joinedCluster.Status.Conditions[0].Type != v1alpha1.ConditionTypeAgentConnected
	if statusTransitioned {
		agentConnected.LastTransitionTime = &currentTime
	}
	joinedCluster.Status.Conditions = []v1alpha1.JoinedClusterConditions{agentConnected}

	clusterName, err := getClusterName(r, log)
	if err != nil {
		log.Error(err, "unable to get the infrastructure for the spoke cluster")
		return ctrl.Result{}, err
	}
	clusterVersion, err := r.DiscoveryClient.ServerVersion()
	if err != nil {
		log.Error(err, "unable to get clusterVerison object for getting the version of openshift distribution running in the spoke cluster")
		return ctrl.Result{}, err
	}
	nodeList := &corev1.NodeList{}
	if err = r.SpokeClient.List(context.TODO(), nodeList, &client.ListOptions{}); err != nil {
		log.Error(err, "unable to get the nodelist from the spoke cluster")
		return ctrl.Result{}, err
	}
	agentInfo := &v1alpha1.ClusterAgentInfo{
		Version:        "v0.0.1",
		Image:          "quay.io/ifont/onprem-agent:latest",
		LastUpdateTime: currentTime,
		ClusterName:    clusterName,
		ClusterVersion: clusterVersion.GitVersion,
		NodeCount:      len(nodeList.Items),
	}
	joinedCluster.Status.ClusterAgentInfo = agentInfo

	if err := r.HubClient.Status().Update(ctx, &joinedCluster); err != nil {
		log.Error(err, "unable to update joinedcluster status")
		return ctrl.Result{}, err
	}

	// Force a reconcile even if no changes.
	return ctrl.Result{RequeueAfter: HeartBeatDelay}, nil
}

func (r *JoinedClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.JoinedCluster{}).
		Complete(r)
}

func getClusterName(r *JoinedClusterReconciler, log logr.Logger) (string, error) {
	infrastructure := &configv1.Infrastructure{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cluster",
			Namespace: "",
		},
	}
	infraObjectKey, err := client.ObjectKeyFromObject(infrastructure)
	if err != nil {
		log.Error(err, "Error getting the object key for infrastructure object", "name", infrastructure.Name)
		return "", err
	}
	err = r.SpokeClient.Get(context.Background(), infraObjectKey, infrastructure)
	if err != nil {
		log.Error(err, "Error getting infrastructure from API server", "name", infrastructure.Name)
		return "", err
	}
	return infrastructure.Status.InfrastructureName, nil
}
