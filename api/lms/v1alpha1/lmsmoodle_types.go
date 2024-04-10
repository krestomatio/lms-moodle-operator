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

// LMSMoodleSpec defines the desired state of LMSMoodle
type LMSMoodleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// LMSMoodleTemplateName defines what LMS Moodle template to use
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	LMSMoodleTemplateName string `json:"lmsMoodleTemplateName"`

	// DesiredState defines the desired state to put a LMSMoodle
	// +kubebuilder:validation:Enum=Ready;Suspended
	// +kubebuilder:default:="Ready"
	// +optional
	DesiredState string `json:"desiredState,omitempty"`

	// LMSMoodleTemplateSpec to set same fields as LMSMoodleTemplate
	LMSMoodleTemplateSpec `json:",inline"`
}

// LMSMoodleStatus defines the observed state of LMSMoodle
type LMSMoodleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions represent the latest available observations of the resource state
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// state describes the LMSMoodle state
	// +kubebuilder:default:="Unknown"
	// +optional
	State string `json:"state,omitempty"`

	// Url defines LMSMoodle url
	// +optional
	Url string `json:"url,omitempty"`

	// StorageGb defines LMSMoodle number of current GB for storage capacity
	// +kubebuilder:default:="0"
	// +optional
	StorageGb string `json:"storageGb,omitempty"`

	// RegisteredUsers defines LMSMoodle number of current registered users for user capacity
	// +kubebuilder:default:=0
	// +optional
	RegisteredUsers int64 `json:"registeredUsers,omitempty"`

	// Release defines LMSMoodle moodle version
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

	// Resource is successful
	SuspendedState string = "Suspended"
)

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster,categories={lms},shortName=lm
//+kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="Age of the resource",priority=0
//+kubebuilder:printcolumn:name="STATUS",type="string",description="LMSMoodle status such as Unknown/SettingUp/Ready/Failed/Terminating etc",JSONPath=".status.state",priority=0
//+kubebuilder:printcolumn:name="SINCE",type="date",JSONPath=".status.conditions[?(@.type=='Ready')].lastTransitionTime",description="Time of latest transition",priority=0
//+kubebuilder:printcolumn:name="TEMPLATE",type="string",description="LMSMoodleTemplate name",JSONPath=".spec.lmsMoodleTemplate",priority=0
//+kubebuilder:printcolumn:name="URL",type="string",JSONPath=".status.url",description="LMSMoodle URL",priority=0
//+kubebuilder:printcolumn:name="USERS",type="integer",JSONPath=".status.registeredUsers",description="LMSMoodle registered users",priority=0
//+kubebuilder:printcolumn:name="GB",type="string",JSONPath=".status.storageGb",description="LMSMoodle storage usage in GB",priority=0

// LMSMoodle is the Schema for the lmsmoodles API
type LMSMoodle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LMSMoodleSpec   `json:"spec,omitempty"`
	Status LMSMoodleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LMSMoodleList contains a list of LMSMoodle
type LMSMoodleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LMSMoodle `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LMSMoodle{}, &LMSMoodleList{})
}
