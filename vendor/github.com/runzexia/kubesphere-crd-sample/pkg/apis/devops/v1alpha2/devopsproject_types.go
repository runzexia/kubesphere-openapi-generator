/*
Copyright 2019 The KubeSphere Authors.

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

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindDevOpsProject     = "DevOpsProject"
	ResourceSingularDevOpsProject = "devopsproject"
	ResourcePluralDevOpsProject   = "devopsprojects"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DevOpsProjectSpec defines the desired state of DevOpsProject
type DevOpsProjectSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// DisplayName is DisplayName of DevOpsProject
	DisplayName string `json:"displayName,omitempty"`

	// Description is Description of DevOpsProject
	Description string `json:"description,omitempty"`
}

// DevOpsProjectStatus defines the observed state of DevOpsProject
type DevOpsProjectStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase is Status of DevOpsProject
	Phase string `json:"phase,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// DevOpsProject is the Schema for the devopsprojects API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type DevOpsProject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevOpsProjectSpec   `json:"spec,omitempty"`
	Status DevOpsProjectStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient:nonNamespaced

// DevOpsProjectList contains a list of DevOpsProject
type DevOpsProjectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevOpsProject `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevOpsProject{}, &DevOpsProjectList{})
}
