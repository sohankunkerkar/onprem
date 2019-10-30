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

// SingleClusterWorkloadSpec defines the desired state of SingleClusterWorkload
type SingleClusterWorkloadSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ServiceAccount string `json:"serviceAccount"`
}

// SingleClusterWorkloadConditionType describes the possible type of conditions that can occur for this resource
type SingleClusterWorkloadConditionType string

const (
	// ConditionTypeRBACPermissions means the RBAC permissions have been created for the cluster agent to perform workload operations.
	ConditionTypeRBACPermissions SingleClusterWorkloadConditionType = "RBACPermissionsReady "
	// ConditionTypeDeploymentSuccess means the deployment of kube resources is successful.
	ConditionTypeDeploymentSuccess SingleClusterWorkloadConditionType = "ResourceDeploymentSuccessful"
	// ConditionTypeDeploymentFailed means the deployement of kube resources failed
	ConditionTypeDeploymentFailed SingleClusterWorkloadConditionType = "ResourceDeploymentFailed"
)

// SingleClusterWorkloadConditions provides the detail about each condition type at any given time.
type SingleClusterWorkloadConditions struct {
	// Type defines the type of JoinedClusterCondition being populated by the controller
	Type SingleClusterWorkloadConditionType `json:"type"`

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

// SingleClusterWorkloadStatus defines the observed state of SingleClusterWorkload
type SingleClusterWorkloadStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []SingleClusterWorkloadConditions `json:"conditions"`
}

// +kubebuilder:object:root=true

// SingleClusterWorkload is the Schema for the singleclusterworkloads API
type SingleClusterWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SingleClusterWorkloadSpec   `json:"spec,omitempty"`
	Status SingleClusterWorkloadStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SingleClusterWorkloadList contains a list of SingleClusterWorkload
type SingleClusterWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SingleClusterWorkload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SingleClusterWorkload{}, &SingleClusterWorkloadList{})
}
