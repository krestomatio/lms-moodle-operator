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

// KeydbSpec defines the desired state of Keydb
// +optional
type KeydbSpec struct {
	// KeydbMode describes mode keydb runs
	// +optional
	KeydbMode KeydbMode `json:"keydbMode,omitempty"`

	// KeydbExtraConfig contains extra keydb config
	// +optional
	KeydbExtraConfig string `json:"keydbExtraConfig,omitempty"`

	// KeydbSize defines keydb number of replicas
	// +optional
	KeydbSize int32 `json:"keydbSize,omitempty"`

	// KeydbImage defines image for keydb container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	KeydbImage string `json:"keydbImage,omitempty"`

	// KeydbPvcDataSize defines keydb storage size
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=20
	// +optional
	KeydbPvcDataSize string `json:"keydbPvcDataSize,omitempty"`

	// KeydbPvcDataStorageAccessMode defines keydb storage access modes
	// +optional
	KeydbPvcDataStorageAccessMode StorageAccessMode `json:"keydbPvcDataStorageAccessMode,omitempty"`

	// KeydbPvcDataStorageClassName defines keydb storage class
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=63
	// +optional
	KeydbPvcDataStorageClassName string `json:"keydbPvcDataStorageClassName,omitempty"`

	// KeydbResourceRequests whether keydb resource requests are added. Default: true
	// +optional
	KeydbResourceRequests bool `json:"keydbResourceRequests,omitempty"`

	// KeydbResourceRequestsCpu set keydb resource requests cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	KeydbResourceRequestsCpu string `json:"keydbResourceRequestsCpu,omitempty"`

	// KeydbResourceRequestsMemory set keydb resource requests memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	KeydbResourceRequestsMemory string `json:"keydbResourceRequestsMemory,omitempty"`

	// KeydbResourceLimits whether keydb resource limits are added. Default: false
	// +optional
	KeydbResourceLimits bool `json:"keydbResourceLimits,omitempty"`

	// KeydbResourceLimitsCpu set keydb resource limits cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	KeydbResourceLimitsCpu string `json:"keydbResourceLimitsCpu,omitempty"`

	// KeydbResourceLimitsMemory set keydb resource limits memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	KeydbResourceLimitsMemory string `json:"keydbResourceLimitsMemory,omitempty"`

	// KeydbTolerations defines any tolerations for Keydb pods.
	// +optional
	KeydbTolerations []corev1.Toleration `json:"keydbTolerations,omitempty"`
}

// KeydbMode describes mode keydb runs
// +kubebuilder:validation:Enum=standalone;multimaster;custom
type KeydbMode string

const (
	// Standalone runs keydb as standlone, single node
	KeydbStandalone KeydbMode = "standalone"

	// Multimaster runs keydb as multimaster, three nodes by default
	KeydbMultimaster KeydbMode = "multimaster"

	// Custom do not set force any mode, the user must configure keydb
	KeydbCustom KeydbMode = "custom"
)
