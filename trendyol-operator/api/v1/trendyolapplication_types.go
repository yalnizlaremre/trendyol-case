/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TrendyolApplicationSpec defines the desired state of TrendyolApplication.
type TrendyolApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Namespace   string            `json:"namespace"`
	Image       string            `json:"image"`
	PullSecret  string            `json:"pullSecret"`
	Command     []string          `json:"command"`
	Arguments   []string          `json:"arguments"`
	Replicas    *int32            `json:"replicas"`
	Environment map[string]string `json:"environment"`

	// Foo is an example field of TrendyolApplication. Edit trendyolapplication_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// TrendyolApplicationStatus defines the observed state of TrendyolApplication.
type TrendyolApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Phase      string `json:"phase,omitempty"`
	DeployedAs string `json:"deployedAs,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// TrendyolApplication is the Schema for the trendyolapplications API.
type TrendyolApplication struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TrendyolApplicationSpec   `json:"spec,omitempty"`
	Status TrendyolApplicationStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TrendyolApplicationList contains a list of TrendyolApplication.
type TrendyolApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TrendyolApplication `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TrendyolApplication{}, &TrendyolApplicationList{})
}
