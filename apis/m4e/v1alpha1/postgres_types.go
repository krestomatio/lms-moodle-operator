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

// PostgresSpec defines the desired state of Postgres
// +optional
type PostgresSpec struct {
	// PostgresMode describes mode postgres runs
	// +optional
	PostgresMode PostgresMode `json:"postgresMode,omitempty"`

	// PostgresExtraConfig contains extra postgres config
	// +optional
	PostgresExtraConfig string `json:"postgresExtraConfig,omitempty"`

	// PostgresSize defines postgres number of replicas
	// +optional
	PostgresSize int32 `json:"postgresSize,omitempty"`

	// PostgresImage defines image for postgres container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	PostgresImage string `json:"postgresImage,omitempty"`

	// PostgresPvcDataSize defines postgres storage size
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresPvcDataSize string `json:"postgresPvcDataSize,omitempty"`

	// PostgresPvcDataStorageAccessMode defines postgres storage access modes
	// +optional
	PostgresPvcDataStorageAccessMode StorageAccessMode `json:"postgresPvcDataStorageAccessMode,omitempty"`

	// PostgresPvcDataStorageClassName defines postgres storage class
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=63
	// +optional
	PostgresPvcDataStorageClassName string `json:"postgresPvcDataStorageClassName,omitempty"`

	// PostgresResourceRequests whether resource requests are set
	// +optional
	PostgresResourceRequests bool `json:"postgresResourceRequests,omitempty"`

	// PostgresResourceRequestsCpu set cpu for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresResourceRequestsCpu string `json:"postgresResourceRequestsCpu,omitempty"`

	// PostgresResourceRequestsMemory set memory for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresResourceRequestsMemory string `json:"postgresResourceRequestsMemory,omitempty"`

	// PostgresResourceLimits whether resource limits are set
	// +optional
	PostgresResourceLimits bool `json:"postgresResourceLimits,omitempty"`

	// PostgresResourceLimitsCpu set cpu for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresResourceLimitsCpu string `json:"postgresResourceLimitsCpu,omitempty"`

	// PostgresResourceLimitsMemory set memory for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresResourceLimitsMemory string `json:"postgresResourceLimitsMemory,omitempty"`

	// PostgresTolerations defines any tolerations for Postgres pods.
	// +optional
	PostgresTolerations []corev1.Toleration `json:"postgresTolerations,omitempty"`

	// PostgresReadreplicasSize defines postgres readreplicas number of replicas
	// +optional
	PostgresReadreplicasSize int32 `json:"postgresReadreplicasSize,omitempty"`

	// PostgresReadreplicasPvcDataSize defines postgres readreplicas storage size
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresReadreplicasPvcDataSize string `json:"postgresReadreplicasPvcDataSize,omitempty"`

	// PostgresReadreplicasPvcDataStorageAccessMode defines postgres readreplicas storage access modes
	// +optional
	PostgresReadreplicasPvcDataStorageAccessMode StorageAccessMode `json:"postgresReadreplicasPvcDataStorageAccessMode,omitempty"`

	// PostgresReadreplicasPvcDataStorageClassName defines postgres readreplicas storage class
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=63
	// +optional
	PostgresReadreplicasPvcDataStorageClassName string `json:"postgresReadreplicasPvcDataStorageClassName,omitempty"`

	// PostgresReadreplicasResourceRequests whether resource requests are set
	// +optional
	PostgresReadreplicasResourceRequests bool `json:"postgresReadreplicasResourceRequests,omitempty"`

	// PostgresReadreplicasResourceRequestsCpu set cpu for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresReadreplicasResourceRequestsCpu string `json:"postgresReadreplicasResourceRequestsCpu,omitempty"`

	// PostgresReadreplicasResourceRequestsMemory set memory for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresReadreplicasResourceRequestsMemory string `json:"postgresReadreplicasResourceRequestsMemory,omitempty"`

	// PostgresReadreplicasResourceLimits whether resource limits are set
	// +optional
	PostgresReadreplicasResourceLimits bool `json:"postgresReadreplicasResourceLimits,omitempty"`

	// PostgresReadreplicasResourceLimitsCpu set cpu for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresReadreplicasResourceLimitsCpu string `json:"postgresReadreplicasResourceLimitsCpu,omitempty"`

	// PostgresReadreplicasResourceLimitsMemory set memory for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PostgresReadreplicasResourceLimitsMemory string `json:"postgresReadreplicasResourceLimitsMemory,omitempty"`

	// PostgresReadreplicasTolerations defines any tolerations for PostgresReadreplicas pods.
	// +optional
	PostgresReadreplicasTolerations []corev1.Toleration `json:"postgresReadreplicasTolerations,omitempty"`

	// PgbouncerExtraConfig contains extra pgbouncer config
	// +optional
	PgbouncerExtraConfig string `json:"pgbouncerExtraConfig,omitempty"`

	// PgbouncerResourceRequests whether resource requests are set
	// +optional
	PgbouncerResourceRequests bool `json:"pgbouncerResourceRequests,omitempty"`

	// PgbouncerResourceRequestsCpu set cpu for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PgbouncerResourceRequestsCpu string `json:"pgbouncerResourceRequestsCpu,omitempty"`

	// PgbouncerResourceRequestsMemory set memory for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PgbouncerResourceRequestsMemory string `json:"pgbouncerResourceRequestsMemory,omitempty"`

	// PgbouncerResourceLimits whether resource limits are set
	// +optional
	PgbouncerResourceLimits bool `json:"pgbouncerResourceLimits,omitempty"`

	// PgbouncerResourceLimitsCpu set cpu for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PgbouncerResourceLimitsCpu string `json:"pgbouncerResourceLimitsCpu,omitempty"`

	// PgbouncerResourceLimitsMemory set memory for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PgbouncerResourceLimitsMemory string `json:"pgbouncerResourceLimitsMemory,omitempty"`

	// PgbouncerTolerations defines any tolerations for Pgbouncer pods.
	// +optional
	PgbouncerTolerations []corev1.Toleration `json:"pgbouncerTolerations,omitempty"`

	// PgbouncerReadonlyExtraConfig contains extra pgbouncerReadonly config
	// +optional
	PgbouncerReadonlyExtraConfig string `json:"pgbouncerReadonlyExtraConfig,omitempty"`

	// PgbouncerReadonlyResourceRequests whether resource requests are set
	// +optional
	PgbouncerReadonlyResourceRequests bool `json:"pgbouncerReadonlyResourceRequests,omitempty"`

	// PgbouncerReadonlyResourceRequestsCpu set cpu for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PgbouncerReadonlyResourceRequestsCpu string `json:"pgbouncerReadonlyResourceRequestsCpu,omitempty"`

	// PgbouncerReadonlyResourceRequestsMemory set memory for resource requests
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PgbouncerReadonlyResourceRequestsMemory string `json:"pgbouncerReadonlyResourceRequestsMemory,omitempty"`

	// PgbouncerReadonlyResourceLimits whether resource limits are set
	// +optional
	PgbouncerReadonlyResourceLimits bool `json:"pgbouncerReadonlyResourceLimits,omitempty"`

	// PgbouncerReadonlyResourceLimitsCpu set cpu for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PgbouncerReadonlyResourceLimitsCpu string `json:"pgbouncerReadonlyResourceLimitsCpu,omitempty"`

	// PgbouncerReadonlyResourceLimitsMemory set memory for resource limits
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PgbouncerReadonlyResourceLimitsMemory string `json:"pgbouncerReadonlyResourceLimitsMemory,omitempty"`

	// PgbouncerReadonlyTolerations defines any tolerations for PgbouncerReadonly pods.
	// +optional
	PgbouncerReadonlyTolerations []corev1.Toleration `json:"pgbouncerReadonlyTolerations,omitempty"`
}

// PostgresMode describes mode postgres runs
// +kubebuilder:validation:Enum=standalone;readreplicas
type PostgresMode string

const (
	// Standalone runs postgres as standlone, single node
	PostgresStandalone PostgresMode = "standalone"

	// Readreplicas runs postgres with readreplicas,  1 primary and 1 replica by default
	PostgresReadreplicas PostgresMode = "readreplicas"
)
