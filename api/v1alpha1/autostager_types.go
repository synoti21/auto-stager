/*
Copyright 2024.

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

// AutostagerSpec defines the desired state of Autostager
type AutostagerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Autostager. Edit autostager_types.go to remove/update
	HelmMode      bool   `json:"helmMode"` // true: use helm, false: use manifest from ci/cd
	Image         string `json:"image"`
	Namespace     string `json:"namespace"`
	ContainerPort int32  `json:"containerPort"`
	Replicas      *int32 `json:"replicas,omitempty"`
	IngressHost   string `json:"ingressHost"`
}

type ServiceAccountSpec struct {
	Create      *bool             `json:"create,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

// AutostagerStatus defines the observed state of Autostager
type AutostagerStatus struct {
	Conditions   []metav1.Condition `json:"conditions,omitempty"`
	Replicas     string             `json:"replicas,omitempty"`
	LastSyncTime metav1.Time        `json:"lastSyncTime"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Autostager is the Schema for the autostagers API
type Autostager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutostagerSpec   `json:"spec,omitempty"`
	Status AutostagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AutostagerList contains a list of Autostager
type AutostagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Autostager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Autostager{}, &AutostagerList{})
}
