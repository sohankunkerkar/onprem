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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JoinedClusterSpec defines the desired state of JoinedCluster
type JoinedClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9]))*$

	// Optional service account name to allow spoke cluster to communicate with the hub when joining
	// If the service account by this name doesn't exist, it will be created in the hub cluster
	// If not specified, a service account will be generated for the spoke cluster to use.
	// +optional
	ServiceAccount *string `json:"serviceAccount,omitempty"`

	// Optional stale duration used to configure the time to wait before
	// determining that the spoke cluster connection has gone stale by not
	// heartbeating back to the hub.
	// +optional
	StaleDuration *metav1.Duration `json:"staleDuration,omitempty"`

	// Optional disconnect duration used to configure the time to wait before
	// determining that the spoke cluster has disconnected by not heartbeating
	// back to the hub after the connection became stale.
	// +optional
	DisconnectDuration *metav1.Duration `json:"disconnectDuration,omitempty"`
}

// ConditionStatus describes the status of the condition as described by the constants below
// +kubebuilder:validation:Enum=True;False;Unknown
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

// JoinedClusterConditionType describes the possible type of conditions that can occur for this resource
// +kubebuilder:validation:Enum=ReadyToJoin;AgentConnected;AgentStale;AgentDisconnected
type JoinedClusterConditionType string

const (
	ConditionTypeReadyToJoin       JoinedClusterConditionType = "ReadyToJoin"
	ConditionTypeAgentConnected    JoinedClusterConditionType = "AgentConnected"
	ConditionTypeAgentStale        JoinedClusterConditionType = "AgentStale"
	ConditionTypeAgentDisconnected JoinedClusterConditionType = "AgentDisconnected"
)

type JoinedClusterConditions struct {
	// Type defines the type of JoinedClusterCondition being populated by the controller
	Type JoinedClusterConditionType `json:"type"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	Status ConditionStatus `json:"status"`
	// Last transition time when this condition got set
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason *string `json:"reason,omitempty"`
	// Human readable message indicating details about last transition
	// +optional
	Message *string `json:"message,omitempty"`
}

// ClusterAgentInfo describes the metadata reported by the cluster agent in the
// spoke cluster.
type ClusterAgentInfo struct {
	// Version of the cluster agent running in the spoke cluster.
	Version string `json:"version"`
	// Image of the cluster agent running int he spoke cluster.
	Image string `json:"image"`
	// Last update time written by cluster agent.
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

// JoinedClusterStatus defines the observed state of JoinedCluster
type JoinedClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//Conditions
	Conditions []JoinedClusterConditions `json:"conditions"`

	// JoinCommand
	// +optional
	JoinCommand *string `json:"joinCommand,omitempty"`

	// +kubebuilder:validation:253
	// +kubebuilder:validation:^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9]))*$

	// ServiceAccount name chosen by the hub for the spoke to use
	// +optional
	ServiceAccountName *string `json:"serviceAccountName,omitempty"`

	// When the cluster agent starts running and heartbeating, it will report
	// metadata information in this field.
	// +optional
	ClusterAgentInfo *ClusterAgentInfo `json:"clusterAgentInfo,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// JoinedCluster is the Schema for the joinedclusters API
type JoinedCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JoinedClusterSpec   `json:"spec,omitempty"`
	Status JoinedClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// JoinedClusterList contains a list of JoinedCluster
type JoinedClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JoinedCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JoinedCluster{}, &JoinedClusterList{})
}
