package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RegisteredCluster represents a cluster known to this hub. The name does NOT match the cluster itself.  Instead we
// recommend using a generated name to avoid conflicts. The name will get a `cluster-` prefix to map to its namespace.
type RegisteredCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RegisteredClusterSpec   `json:"spec"`
	Status RegisteredClusterStatus `json:"status"`
}

type RegisteredClusterSpec struct {
	ClusterName string `json:"clusterName"`
}

type RegisteredClusterStatus struct {
	// conditions describes the state of the controller's reconciliation functionality.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlacementPolicy matches a set of MultiClusterWorkloads with a set of RegisteredClusters.  Keep in mind that it is possible
// for a single MultiClusterWorkload,RegisteredClusters tuple to be produced by different PlacementPolicies.
// Access control inside the hub is uniform, so if a user has the power to create a PlacementPolicy, he can put ANY
// MultiClusterWorkload on any RegisteredCluster.  However, he cannot control the CONTENT of that MultiClusterWorkload
// unless he has the power to "use" the "asUser" in that MultiClusterWorkload.
type PlacementPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlacementPolicySpec   `json:"spec"`
	Status PlacementPolicyStatus `json:"status"`
}

type PlacementPolicySpec struct {
	// multiClusterWorkloadSelector describes how to select MultiClusterWorkloads.  Every match for every selector in
	// the list is added to the list of workloads. This makes the slice a logical OR
	MultiClusterWorkloadSelectors []MultiClusterWorkloadSelector `json:"multiClusterWorkloadSelector"`
	// registeredClusterSelector describes how to select RegisteredClusters.  Every match for every selector in
	// the list is added to the list of clusters. This makes the slice a logical OR
	RegisteredClusterSelectors []RegisteredClusterSelector `json:"registeredClusterSelector"`
}

type MultiClusterWorkloadSelector struct {
	// type indicates the type of workload to put into a spoke cluster.  Valid values:
	//  1. "LabelSelector"
	//  1. "Names"
	Type string `json:"type"`

	// labelSelector specifies a label selector to evaluate against all MultiClusterWorkloads.  It is mutually exclusive with names.
	LabelSelector *metav1.LabelSelector `json:"labelSelector"`
	// names specifies specific names to match MultiClusterWorkload.  It is mutually exclusive with labelSelector.
	Names []string `json:"names"`
}

type RegisteredClusterSelector struct {
	// type indicates the type of workload to put into a spoke cluster.  Valid values:
	//  1. "LabelSelector"
	//  1. "Names"
	Type string `json:"type"`

	// labelSelector specifies a label selector to evaluate against all RegisteredClusterSelectors.  It is mutually exclusive with names.
	LabelSelector *metav1.LabelSelector `json:"labelSelector"`
	// names specifies specific names to match RegisteredClusterSelector.  It is mutually exclusive with labelSelector.
	Names []string `json:"names"`
}

type PlacementPolicyStatus struct {
	// conditions describes the state of the controller's reconciliation functionality.  This only reflects the ability to
	// resolve the MultiClusterWorkloads and RegisteredClusters.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`

	// multiClusterWorkloads is the list of matching MultiClusterWorkloads
	MultiClusterWorkloads []string `json:"multiClusterWorkloads"`
	// registeredClusters is the list of matching RegisteredClusters
	RegisteredClusters []string `json:"registeredClusters"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MultiClusterWorkload describes what could be placed into a spoke cluster.  It does not describe which (if any) clusters should
// create the workload.  That is done based on a PlacementPolicy.  If a PlacementPolicy matches a MultiClusterWorkload
// and a RegisteredCluster, then a namespaced Workload resource will be created in the cluster's namespace.
type MultiClusterWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultiClusterWorkloadSpec   `json:"spec"`
	Status MultiClusterWorkloadStatus `json:"status"`
}

type MultiClusterWorkloadSpec struct {
	// workload describes the resources to be placed on the spoke.
	Workload WorkloadSpec `json:"workload"`

	// asUser indicates which users should be used in the spoke cluster to create this workload.  Because the workload's
	// ability to be successfully created in the cluster is tightly coupled to the subject being used to create that workload
	// in the spoke cluster, this is a property of the MultiClusterWorkloadSpec.
	// This name must match a DefinedSubject.metadata.name to pass an admission validation check, but the reference can
	// become stale at a later date.
	// When a client creates or updates a MultiClusterWorkload, a secondary ACL check is performed to see if the client
	// can "use" the referenced subject.  The "use" check is uniform across ALL clusters registered in the hub. It doesn't
	// mean that every cluster honors this subject, but for each cluster does honor it, it is possible for this MultiClusterWorkload
	// to be placed and created there.
	AsUser Subject `json:"asUser"`
}

// Subject describes a single subject by user and groups.
// +union
type Subject struct {
	User   string   `json:"user"`
	Groups []string `json:"groups"`
}

type WorkloadSpec struct {
	// type indicates the type of workload to put into a spoke cluster.  Valid values:
	//  1. "Manifests"
	// in the future, more types could be added by helm chart or some other higher level entity.
	Type string `json:"type"`

	// manifests holds a metav1.List that contains items for the individual manifests that you want created in the
	// spoke cluster.
	Manifests *Manifests `json:"manifests"`
}

// Manifests holds a metav1.List that contains items for the individual manifests that you want created in the
// spoke cluster.
type Manifests struct {
	metav1.List `json:",inline"`
}

type MultiClusterWorkloadStatus struct {
	// conditions describes the state of the controller's reconciliation functionality.  This is a union of the ClusterStatuses
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`

	// clusterStatuses track the deployment values and errors across individual clusters
	// +optional
	ClusterStatuses []ClusterStatus `json:"clusterStatuses"`
}

type ClusterStatus struct {
	// registeredCluster matches the name of a RegisteredCluster resource
	RegisteredCluster string `json:"registeredCluster"`

	// conditions describes the state of the controller's reconciliation functionality.  It only includes exception statuses
	// like failures, not healthy states.  Anna Karenina and all that.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterWorkload describes what could be placed into a spoke cluster.  It does not describe which (if any) clusters should
// create the workload.  That is done based on a PlacementPolicy.  If a PlacementPolicy matches a MultiClusterWorkload
// and a RegisteredCluster, then a namespaced Workload resource will be created in the cluster's namespace.
type ClusterWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterWorkloadSpec   `json:"spec"`
	Status ClusterWorkloadStatus `json:"status"`
}

type ClusterWorkloadSpec struct {
	// workload describes the resources to be placed on the spoke.
	Workload WorkloadSpec `json:"workload"`

	// asUser indicates which users should be used in the spoke cluster to create this workload.  Because the workload's
	// ability to be successfully created in the cluster is tightly coupled to the subject being used to create that workload
	// in the spoke cluster, this is a property of the MultiClusterWorkloadSpec.
	// This name must match a DefinedSubject.metadata.name to pass an admission validation check, but the reference can
	// become stale at a later date.
	// When a client creates or updates a MultiClusterWorkload, a secondary ACL check is performed to see if the client
	// can "use" the referenced subject.  The "use" check is uniform across ALL clusters registered in the hub. It doesn't
	// mean that every cluster honors this subject, but for each cluster does honor it, it is possible for this MultiClusterWorkload
	// to be placed and created there.
	AsUser Subject `json:"asUser"`
}

type ClusterWorkloadStatus struct {
	// conditions describes the state of the controller's reconciliation functionality.  This is a union of the WorkloadStatus and
	// overall conditions like...
	//  1. SubjectHonored
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`

	// workloadStatus track the deployment values and errors for particular kinds of workloads
	// +optional
	WorkloadStatus WorkloadStatus `json:"workloadStatus"`
}

type WorkloadStatus struct {
	// type indicates the type of workload to put into a spoke cluster.  These should match the types of workloads.
	// Valid values:
	//  1. "Manifests"
	Type string `json:"type"`

	// manifestsStatus holds a metav1.List that contains items for the individual manifests that you want created in the
	// spoke cluster.
	ManifestsStatus *ManifestStatus `json:"manifestsStatus"`
}

type ManifestStatus struct {
	// conditions describes the state of the controller's reconciliation functionality.  This is a union of the ItemStatuses
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`

	// ItemStatuses track the deployment values by item
	// +optional
	ItemStatuses []ItemStatus `json:"itemStatuses"`
}

type ItemStatus struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name"`

	// conditions describes the state of the controller's reconciliation functionality.
	// known conditions include:
	//  1. Pending
	//  2. Reconciled
	//  3. Complete - for "known" resources like deployments where we wait for a condition like available pods, indicates that
	//                the logically known-good state has been reached.  For all other resource, this simply means that
	//                the apply call succeeded.
	//  4. Errored - errored indicates that REST returned a non-2xx return code OR that complete has taken too long.
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +optional
	Conditions []Condition `json:"conditions,omitempty"  patchStrategy:"merge" patchMergeKey:"type"`

	// ItemStatuses track the deployment values by item
	// +optional
	ItemStatuses []ItemStatus `json:"itemStatuses"`
}

type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition.
// "ConditionFalse" means a resource is not in the condition. "ConditionUnknown" means kubernetes
// can't decide if a resource is in the condition or not. In the future, we could add other
// intermediate conditions, e.g. ConditionDegraded.
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// ClusterOperatorStatusCondition represents the state of the operator's
// reconciliation functionality.
// +k8s:deepcopy-gen=true
type Condition struct {
	// type specifies the state of the operator's reconciliation functionality.
	Type string `json:"type"`

	// status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status"`

	// lastTransitionTime is the time of the last update to the current status object.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// reason is the reason for the condition's last transition.  Reasons are CamelCase
	Reason string `json:"reason,omitempty"`

	// message provides additional information about the current condition.
	// This is only to be consumed by humans.
	Message string `json:"message,omitempty"`
}
