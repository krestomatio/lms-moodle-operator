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

	// GaneshaResourceRequests whether ganesha resource requests are added. Default: true
	// +optional
	GaneshaResourceRequests bool `json:"ganeshaResourceRequests,omitempty"`

	// GaneshaResourceRequestsCpu set ganesha resource requests cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaResourceRequestsCpu string `json:"ganeshaResourceRequestsCpu,omitempty"`

	// GaneshaResourceRequestsMemory set ganesha resource requests memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaResourceRequestsMemory string `json:"ganeshaResourceRequestsMemory,omitempty"`

	// GaneshaResourceLimits whether ganesha resource limits are added. Default: false
	// +optional
	GaneshaResourceLimits bool `json:"ganeshaResourceLimits,omitempty"`

	// GaneshaResourceLimitsCpu set ganesha resource limits cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaResourceLimitsCpu string `json:"ganeshaResourceLimitsCpu,omitempty"`

	// GaneshaResourceLimitsMemory set ganesha resource limits memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	GaneshaResourceLimitsMemory string `json:"ganeshaResourceLimitsMemory,omitempty"`

	// GaneshaTolerations defines any tolerations for Ganesha server pods.
	// +optional
	GaneshaTolerations []corev1.Toleration `json:"ganeshaTolerations,omitempty"`

	// GaneshaNodeSelector defines any node labels selectors for Ganesha pods.
	// +optional
	GaneshaNodeSelector string `json:"ganeshaNodeSelector,omitempty"`

	// GaneshaAffinity defines any affinity rules for Ganesha pods.
	// +optional
	GaneshaAffinity string `json:"ganeshaAffinity,omitempty"`

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

	// GaneshaPvcDataAutoexpansion enables autoexpansion
	// +optional
	GaneshaPvcDataAutoexpansion bool `json:"ganeshaPvcDataAutoexpansion,omitempty"`

	// GaneshaPvcDataAutoexpansionIncrementGib defines Gib to increment
	// +optional
	GaneshaPvcDataAutoexpansionIncrementGib int32 `json:"ganeshaPvcDataAutoexpansionIncrementGib,omitempty"`

	// GaneshaPvcDataAutoexpansionCapGib defines limit for autoexpansion increments
	// +optional
	GaneshaPvcDataAutoexpansionCapGib int32 `json:"ganeshaPvcDataAutoexpansionCapGib,omitempty"`

	// GaneshaExtraBlockConfig contains extra block in ganesha server ganesha config
	// +optional
	GaneshaExtraBlockConfig string `json:"ganeshaExtraBlockConfig,omitempty"`

	// GaneshaConfLogLevel defines nfs log level. Default: EVENT
	// +kubebuilder:validation:Enum=NULL;FATAL;MAJ;CRIT;WARN;EVENT;INFO;DEBUG;MID_DEBUG;M_DBG;FULL_DEBUG;F_DBG
	// +optional
	GaneshaConfLogLevel string `json:"ganeshaConfLogLevel,omitempty"`

	// GaneshaVpaSpec set ganesha horizontal pod autoscaler spec
	// +optional
	GaneshaVpaSpec string `json:"ganeshaVpaSpec,omitempty"`
}
