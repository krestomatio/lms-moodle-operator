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
	// ServerImage defines image for server container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	ServerImage string `json:"serverImage,omitempty"`

	// ServerPvcDataSize defines server storage size
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=20
	// +optional
	ServerPvcDataSize string `json:"serverPvcDataSize,omitempty"`

	// ServerPvcDataStorageAccessMode defines server storage access modes
	// +optional
	ServerPvcDataStorageAccessMode StorageAccessMode `json:"serverPvcDataStorageAccessMode,omitempty"`

	// ServerPvcDataStorageClassName defines server storage class
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=63
	// +optional
	ServerPvcDataStorageClassName string `json:"serverPvcDataStorageClassName,omitempty"`

	// ServerResourceRequests whether resource requests are set
	// +optional
	ServerResourceRequests bool `json:"serverResourceRequests,omitempty"`

	// ServerResourceRequestsCpu set cpu for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	ServerResourceRequestsCpu string `json:"serverResourceRequestsCpu,omitempty"`

	// ServerResourceRequestsMemory set memory for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	ServerResourceRequestsMemory string `json:"serverResourceRequestsMemory,omitempty"`

	// ServerResourceLimits whether resource limits are set
	// +optional
	ServerResourceLimits bool `json:"serverResourceLimits,omitempty"`

	// ServerResourceLimitsCpu set cpu for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	ServerResourceLimitsCpu string `json:"serverResourceLimitsCpu,omitempty"`

	// ServerResourceLimitsMemory set memory for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	ServerResourceLimitsMemory string `json:"serverResourceLimitsMemory,omitempty"`

	// ServerTolerations defines any tolerations for Server pods.
	// +optional
	ServerTolerations []corev1.Toleration `json:"serverTolerations,omitempty"`

	// ServerExportUserid defines export folder userid
	// +optional defines export folder userid
	ServerExportUserid int32 `json:"serverExportUserid,omitempty"`

	// ServerExportGroupid defines export folder groupid
	// +optional
	ServerExportGroupid int32 `json:"serverExportGroupid,omitempty"`

	// ServerExportMode defines folder permissions mode
	// +kubebuilder:validation:Pattern="[0-7]{4}"
	// +optional
	ServerExportMode string `json:"serverExportMode,omitempty"`

	// ServerPvcAutoexpansion enables autoexpansion
	// +optional
	ServerPvcAutoexpansion bool `json:"serverPvcAutoexpansion,omitempty"`

	// ServerPvcAutoexpansionIncrementGib defines Gib to increment
	// +optional
	ServerPvcAutoexpansionIncrementGib int32 `json:"serverPvcAutoexpansionIncrementGib,omitempty"`

	// ServerPvcAutoexpansionCapGib defines limit for autoexpansion increments
	// +optional
	ServerPvcAutoexpansionCapGib int32 `json:"serverPvcAutoexpansionCapGib,omitempty"`

	// ServerPvcAutoexpansionInitSizeGib defines initial pvc size
	// +optional
	ServerPvcAutoexpansionInitSizeGib int32 `json:"serverPvcAutoexpansionInitSizeGib,omitempty"`

	// ServerGaneshaExtraBlockConfig contains extra block in ganesha server config
	// +optional
	ServerGaneshaExtraBlockConfig string `json:"serverGaneshaExtraBlockConfig,omitempty"`

	// ServerConfLogLevel defines nfs log level. Default: EVENT
	// +kubebuilder:validation:Enum=NULL;FATAL;MAJ;CRIT;WARN;EVENT;INFO;DEBUG;MID_DEBUG;M_DBG;FULL_DEBUG;F_DBG
	// +optional
	ServerConfLogLevel string `json:"serverConfLogLevel,omitempty"`
}
