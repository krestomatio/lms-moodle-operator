/*
Copyright 2021.

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

// FlavorSpec defines the desired state of Flavor
type FlavorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// MoodleNewInstance whether new instance job runs
	// +optional
	MoodleNewInstance bool `json:"moodleNewInstance,omitempty"`

	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleNewInstanceFullname string `json:"moodleNewInstanceFullname,omitempty"`

	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleNewInstanceShortname string `json:"moodleNewInstanceShortname,omitempty"`

	// +kubebuilder:validation:MaxLength=300
	// +optional
	MoodleNewInstanceSummary string `json:"moodleNewInstanceSummary,omitempty"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleNewInstanceAdminuser string `json:"moodleNewInstanceAdminuser,omitempty"`

	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=100
	// MoodleNewInstanceAdminPass is the admin password to set in new instance. Required
	MoodleNewInstanceAdminpass string `json:"moodleNewInstanceAdminpass"`

	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=100
	// MoodleNewInstanceAdminMail is the admin email to set in new instance. Required
	MoodleNewInstanceAdminmail string `json:"moodleNewInstanceAdminmail"`

	// MoodleNewInstanceAgreeLicense whether agree to Moodle license. Required
	MoodleNewInstanceAgreeLicense bool `json:"moodleNewInstanceAgreeLicense"`
}

// FlavorStatus defines the observed state of Flavor
type FlavorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Sites store the name of the sites which are using this flavor
	Sites []string `json:"nodes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Flavor is the Schema for the flavors API
type Flavor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FlavorSpec   `json:"spec,omitempty"`
	Status FlavorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// FlavorList contains a list of Flavor
type FlavorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Flavor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Flavor{}, &FlavorList{})
}
