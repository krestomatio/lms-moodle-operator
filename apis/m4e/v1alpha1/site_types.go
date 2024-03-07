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

// SiteSpec defines the desired state of Site
type SiteSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Flavor defines what Moodle flavor to use
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Flavor string `json:"flavor"`

	// FlavorSpec to set same fields as Flavor
	FlavorSpec `json:",inline"`
}

// SiteStatus defines the observed state of Site
type SiteStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions represent the latest available observations of the resource state
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// state describes the site state
	// +kubebuilder:default:="Unknown"
	// +optional
	State string `json:"state,omitempty"`

	// Site url
	// +optional
	Url string `json:"url,omitempty"`

	// Site number of current GB for storage capacity
	// +kubebuilder:default:="0"
	// +optional
	StorageGb string `json:"storageGb,omitempty"`

	// Site number of current registered users for user capacity
	// +kubebuilder:default:=0
	// +optional
	RegisteredUsers int64 `json:"registeredUsers,omitempty"`

	// Site moodle version
	// +optional
	Release string `json:"release,omitempty"`
}

const (
	// Resource is in an unknown
	UnknownState string = "Unknown"

	// Resource has failed
	FailedState string = "Failed"

	// Resource is ready
	ReadyState string = "Ready"

	// Resource is being deleted
	TerminatingState string = "Terminating"

	// Resource has being deleted
	TerminatedState string = "Terminated"

	// Resource is successful
	SuccessfulState string = "Successful"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="Age of the resource",priority=0
//+kubebuilder:printcolumn:name="STATUS",type="string",description="Site status such as Unknown/SettingUp/Ready/Failed/Terminating etc",JSONPath=".status.state",priority=0
//+kubebuilder:printcolumn:name="SINCE",type="date",JSONPath=".status.conditions[?(@.type=='Ready')].lastTransitionTime",description="Time of latest transition",priority=0
//+kubebuilder:printcolumn:name="FLAVOR",type="string",description="Flavor name",JSONPath=".spec.flavor",priority=0
//+kubebuilder:printcolumn:name="HOST",type="string",JSONPath=".spec.moodleSpec.moodleHost",description="Site URL",priority=0

// Site is the Schema for the sites API
type Site struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SiteSpec   `json:"spec,omitempty"`
	Status SiteStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SiteList contains a list of Site
type SiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Site `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Site{}, &SiteList{})
}
