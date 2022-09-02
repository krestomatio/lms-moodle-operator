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

import corev1 "k8s.io/api/core/v1"

// NfsSpec defines the desired state of Nfs
// +optional
type NfsSpec struct {
	// GaneshaImage defines image for ganesha server container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	GaneshaImage string `json:"ganeshaImage,omitempty"`

	// GaneshaPvcDataSize defines ganesha server storage size
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaPvcDataSize string `json:"ganeshaPvcDataSize,omitempty"`

	// GaneshaPvcDataStorageAccessMode defines ganesha server storage access modes
	// +optional
	GaneshaPvcDataStorageAccessMode StorageAccessMode `json:"ganeshaPvcDataStorageAccessMode,omitempty"`

	// GaneshaPvcDataStorageClassName defines ganesha server storage class
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=63
	// +optional
	GaneshaPvcDataStorageClassName string `json:"ganeshaPvcDataStorageClassName,omitempty"`

	// GaneshaResourceRequests whether resource requests are set
	// +optional
	GaneshaResourceRequests bool `json:"ganeshaResourceRequests,omitempty"`

	// GaneshaResourceRequestsCpu set cpu for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaResourceRequestsCpu string `json:"ganeshaResourceRequestsCpu,omitempty"`

	// GaneshaResourceRequestsMemory set memory for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaResourceRequestsMemory string `json:"ganeshaResourceRequestsMemory,omitempty"`

	// GaneshaResourceLimits whether resource limits are set
	// +optional
	GaneshaResourceLimits bool `json:"ganeshaResourceLimits,omitempty"`

	// GaneshaResourceLimitsCpu set cpu for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaResourceLimitsCpu string `json:"ganeshaResourceLimitsCpu,omitempty"`

	// GaneshaResourceLimitsMemory set memory for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaResourceLimitsMemory string `json:"ganeshaResourceLimitsMemory,omitempty"`

	// GaneshaTolerations defines any tolerations for Ganesha server pods.
	// +optional
	GaneshaTolerations []corev1.Toleration `json:"ganeshaTolerations,omitempty"`

	// GaneshaExportUserid defines export folder userid
	// +optional defines export folder userid
	GaneshaExportUserid int32 `json:"ganeshaExportUserid,omitempty"`

	// GaneshaExportGroupid defines export folder groupid
	// +optional
	GaneshaExportGroupid int32 `json:"ganeshaExportGroupid,omitempty"`

	// GaneshaExportMode defines folder permissions mode
	// +kubebuilder:validation:Pattern="[0-7]{4}"
	// +optional
	GaneshaExportMode string `json:"ganeshaExportMode,omitempty"`

	// GaneshaPvcAutoexpansion enables autoexpansion
	// +optional
	GaneshaPvcAutoexpansion bool `json:"ganeshaPvcAutoexpansion,omitempty"`

	// GaneshaPvcAutoexpansionIncrementGib defines Gib to increment
	// +optional
	GaneshaPvcAutoexpansionIncrementGib int32 `json:"ganeshaPvcAutoexpansionIncrementGib,omitempty"`

	// GaneshaPvcAutoexpansionCapGib defines limit for autoexpansion increments
	// +optional
	GaneshaPvcAutoexpansionCapGib int32 `json:"ganeshaPvcAutoexpansionCapGib,omitempty"`

	// GaneshaExtraBlockConfig contains extra block in ganesha server ganesha config
	// +optional
	GaneshaExtraBlockConfig string `json:"ganeshaExtraBlockConfig,omitempty"`

	// GaneshaConfLogLevel defines nfs log level. Default: EVENT
	// +kubebuilder:validation:Enum=NULL;FATAL;MAJ;CRIT;WARN;EVENT;INFO;DEBUG;MID_DEBUG;M_DBG;FULL_DEBUG;F_DBG
	// +optional
	GaneshaConfLogLevel string `json:"ganeshaConfLogLevel,omitempty"`
}
