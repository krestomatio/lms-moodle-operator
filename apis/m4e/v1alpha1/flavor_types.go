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

	// M4eSpec defines M4e spec
	M4eSpec FlavorM4eSpec `json:"m4eSpec"`

	// NfsSpec defines NFS Server spec to deploy optionally
	// +optional
	NfsSpec FlavorNfsSpec `json:"nfsSpec"`
}

// FlavorM4eSpec defines the desired state of M4e
type FlavorM4eSpec struct {
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

	// MoodleNewInstanceAdminPass is the admin password to set in new instance. Required
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=100
	MoodleNewInstanceAdminpass string `json:"moodleNewInstanceAdminpass"`

	// MoodleNewInstanceAdminMail is the admin email to set in new instance. Required
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=100
	MoodleNewInstanceAdminmail string `json:"moodleNewInstanceAdminmail"`

	// MoodleNewInstanceAgreeLicense whether agree to Moodle license. Required
	MoodleNewInstanceAgreeLicense bool `json:"moodleNewInstanceAgreeLicense"`

	// MoodlePvcMoodledataSize defines moodledata storage size
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodlePvcMoodledataSize string `json:"moodlePvcMoodledataSize,omitempty"`

	// MoodlePvcMoodledataStorageAccessMode defines moodledata storage access modes
	// +optional
	MoodlePvcMoodledataStorageAccessMode StorageAccessMode `json:"moodlePvcMoodledataStorageAccessMode,omitempty"`

	// MoodlePvcMoodledataStorageClassName defines moodledata storage class
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=63
	// +optional
	MoodlePvcMoodledataStorageClassName string `json:"moodlePvcMoodledataStorageClassName,omitempty"`
}

// StorageAccessMode describes storage access modes
// +kubebuilder:validation:Enum=ReadWriteOnce;ReadOnlyMany;ReadWriteMany
type StorageAccessMode string

const (
	// ReadWriteOnce can be mounted as read-write by a single node
	ReadWriteOnce StorageAccessMode = "ReadWriteOnce"

	// ReadOnlyMany can be mounted read-only by many nodes
	ReadOnlyMany StorageAccessMode = "ReadOnlyMany"

	// ReadWriteMany the volume can be mounted as read-write by many nodes
	ReadWriteMany StorageAccessMode = "ReadWriteMany"
)

// FlavorNfsSpec defines the desired state of Nfs
// +optional
type FlavorNfsSpec struct {

	// +optional
	ServerExportUserid int32 `json:"serverExportUserid,omitempty"`

	// +optional
	ServerExportGroupid int32 `json:"serverExportGroupid,omitempty"`

	// +kubebuilder:validation:Pattern="[0-7]{4}"
	// +optional
	ServerExportPermissions string `json:"serverExportPermissions,omitempty"`

	// +optional
	ServerPvcAutoexpansion bool `json:"serverPvcAutoexpansion,omitempty"`

	// +optional
	ServerPvcAutoexpansionGib int32 `json:"serverPvcAutoexpansionGib,omitempty"`

	// +optional
	ServerPvcAutoexpansionCapGib int32 `json:"serverPvcAutoexpansionCapGib,omitempty"`

	// +optional
	ServerPvcInitSizeGib int32 `json:"serverPvcInitSizeGib,omitempty"`

	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	// +optional
	ServerPvcStorageClassName string `json:"serverPvcStorageClassName,omitempty"`
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
