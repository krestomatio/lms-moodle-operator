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

	// MoodleSpec defines Moodle spec
	MoodleSpec MoodleSpec `json:"moodleSpec"`

	// PostgresSpec defines Postgres spec to deploy optionally
	// +optional
	PostgresSpec PostgresSpec `json:"postgresSpec"`

	// NfsSpec defines (NFS) Ganesha server spec to deploy optionally
	// +optional
	NfsSpec NfsSpec `json:"nfsSpec"`

	// KeydbSpec defines Keydb spec to deploy optionally
	// +optional
	KeydbSpec KeydbSpec `json:"keydbSpec"`
}

// FlavorStatus defines the observed state of Flavor
type FlavorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// state describes the site state
	// +kubebuilder:default:="Unknown"
	// +optional
	State StatusState `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="Age of the resource",priority=0
//+kubebuilder:printcolumn:name="STATUS",type="string",description="Flavor status such as Unknown/Used/NotUsed/Terminating etc",JSONPath=".status.state",priority=0

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
