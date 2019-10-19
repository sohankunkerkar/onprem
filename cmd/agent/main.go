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

package main

import (
	"flag"
	"os"

	clustermanagerv1alpha1 "github.com/font/onprem/api/v1alpha1"
	"github.com/font/onprem/pkg/controllers/agent"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = clustermanagerv1alpha1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()

	ctrl.SetLogger(zap.Logger(true))

	spokeConfig := ctrl.GetConfigOrDie()
	spokeClient, err := client.New(spokeConfig, client.Options{})
	if err != nil {
		setupLog.Error(err, "Unable to create spoke client")
		os.Exit(1)
	}

	joinedClusterCoordinates, err := agent.GetJoinedClusterCoordinates(spokeClient)
	if err != nil {
		setupLog.Error(err, "Unable to get join cluster coordinates")
		os.Exit(1)
	}

	hubClusterConfig, err := agent.BuildHubClusterConfig(spokeClient)
	if err != nil {
		setupLog.Error(err, "Unable to build hub cluster config")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(hubClusterConfig, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		SyncPeriod:         &agent.HeartBeatDelay,
		Namespace:          joinedClusterCoordinates.Namespace,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&agent.JoinedClusterReconciler{
		HubClient:   mgr.GetClient(),
		SpokeClient: spokeClient,
		Log:         ctrl.Log.WithName("agent").WithName("JoinedCluster"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "JoinedCluster")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
