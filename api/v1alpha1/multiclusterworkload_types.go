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

// MultiClusterWorkloadSpec defines the desired state of MultiClusterWorkload
type MultiClusterWorkloadSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// This service account is used by the impersonation controller
	ServiceAccount string `json:"serviceAccount"`
}

// MultiClusterWorkloadStatus defines the observed state of MultiClusterWorkload
type MultiClusterWorkloadStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	//+optional
	// It shows the list of all spoke clusters connected to the hub
	ListOfClusters []string `json:"listOfClusters,omitempty"`
}

// +kubebuilder:object:root=true

// MultiClusterWorkload is the Schema for the multiclusterworkloads API
type MultiClusterWorkload struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MultiClusterWorkloadSpec   `json:"spec,omitempty"`
	Status MultiClusterWorkloadStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MultiClusterWorkloadList contains a list of MultiClusterWorkload
type MultiClusterWorkloadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MultiClusterWorkload `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MultiClusterWorkload{}, &MultiClusterWorkloadList{})
}
