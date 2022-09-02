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
	corev1 "k8s.io/api/core/v1"
)

// M4eSpec defines the desired state of M4e
type M4eSpec struct {
	// MoodleImage defines image for moodle container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	MoodleImage string `json:"moodleImage,omitempty"`

	// MoodleNewInstance whether new instance job runs
	// +optional
	MoodleNewInstance bool `json:"moodleNewInstance,omitempty"`

	// MoodleNewInstanceAgreeLicense whether agree to Moodle license. Required
	MoodleNewInstanceAgreeLicense bool `json:"moodleNewInstanceAgreeLicense"`

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

	// MoodleNewInstanceAdminMail is the admin email to set in new instance. Required
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=100
	MoodleNewInstanceAdminmail string `json:"moodleNewInstanceAdminmail"`

	// MoodleNewAdminPassHash is the bcrypt compatible admin password to set in new instance. Required
	// +kubebuilder:validation:MinLength=60
	// +kubebuilder:validation:MaxLength=60
	// +kubebuilder:validation:Pattern="^\\$2[ayb]\\$.{56}$"
	MoodleNewAdminpassHash string `json:"moodleNewAdminpassHash"`

	// MoodlePvcDataSize defines moodledata storage size
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodlePvcDataSize string `json:"moodlePvcDataSize,omitempty"`

	// MoodlePvcDataStorageAccessMode defines moodledata storage access modes
	// +optional
	MoodlePvcDataStorageAccessMode StorageAccessMode `json:"moodlePvcDataStorageAccessMode,omitempty"`

	// MoodlePvcDataStorageClassName defines moodledata storage class
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=63
	// +optional
	MoodlePvcDataStorageClassName string `json:"moodlePvcDataStorageClassName,omitempty"`

	// MoodleHost defines Moodle host for url
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=100
	MoodleHost string `json:"moodleHost,omitempty"`

	// MoodlePort defines Moodle port for url
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	MoodlePort int32 `json:"moodlePort,omitempty"`

	// MoodleSubpath defines Moodle subpath for url
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleSubpath string `json:"moodleSubpath,omitempty"`

	// MoodleHealthcheckSubpath defines Moodle subpath for nginx check
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleHealthcheckSubpath string `json:"moodleHealthcheckSubpath,omitempty"`

	// MoodleProtocol whether to use http or https
	// +optional
	MoodleProtocol MoodleProtocol `json:"moodleProtocol,omitempty"`

	// MoodleCronjobTolerations defines any tolerations for Moodle cronjob pods.
	// +optional
	MoodleCronjobTolerations []corev1.Toleration `json:"moodleCronjobTolerations,omitempty"`

	// MoodleConfigAdditionalCfg defines moodle extra config properties in config.php
	// +optional
	MoodleConfigAdditionalCfg MoodleConfigProperty `json:"moodleConfigAdditionalCfg,omitempty"`

	// MoodleConfigAdditionalBlock defines moodle extra block in config.php
	// +optional
	MoodleConfigAdditionalBlock string `json:"moodleConfigAdditionalBlock,omitempty"`

	// MoodleUpdateMinor whether minor updates are automatically applied. Default: true
	// +optional
	MoodleUpdateMinor bool `json:"moodleUpdateMinor,omitempty"`

	// MoodleUpdateMajor whether major updates are automatically applied. Default: false
	// +optional
	MoodleUpdateMajor bool `json:"moodleUpdateMajor,omitempty"`

	// NginxSize defines nginx number of replicas between 0 and 255
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=255
	// +optional
	NginxSize int32 `json:"nginxSize,omitempty"`

	// NginxImage defines image for nginx container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	NginxImage string `json:"nginxImage,omitempty"`

	// NginxIngressAnnotations defines nginx annotations
	// +optional
	NginxIngressAnnotations string `json:"nginxIngressAnnotations,omitempty"`

	// NginxTolerations defines any tolerations for Nginx pods.
	// +optional
	NginxTolerations []corev1.Toleration `json:"nginxTolerations,omitempty"`

	// NginxExtraConfig contains extra Nginx config
	// +optional
	NginxExtraConfig string `json:"nginxExtraConfig,omitempty"`

	// PhpFpmSize defines Phpfpm number of replicas between 0 and 255
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=255
	// +optional
	PhpFpmSize int32 `json:"phpFpmSize,omitempty"`

	// PhpFpmImage defines image for Phpfpm container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	PhpFpmImage string `json:"phpFpmImage,omitempty"`

	// PhpFpmIngressAnnotations defines Phpfpm annotations
	// +optional
	PhpFpmIngressAnnotations string `json:"phpFpmIngressAnnotations,omitempty"`

	// PhpFpmTolerations defines any tolerations for Phpfpm pods.
	// +optional
	PhpFpmTolerations []corev1.Toleration `json:"phpFpmTolerations,omitempty"`

	// PhpFpmPhpExtraIni contains extra php ini config
	// +optional
	PhpFpmPhpExtraIni string `json:"phpFpmPhpExtraIni,omitempty"`

	// PhpFpmExtraConfig contains extra php-fpm config
	// +optional
	PhpFpmExtraConfig string `json:"phpFpmExtraConfig,omitempty"`

	// MoodlePostgresMetaName defines Postgres CR name to use as database.
	// +kubebuilder:validation:MaxLength=63
	// +optional
	MoodlePostgresMetaName string `json:"moodlePostgresMetaName,omitempty"`

	// MoodleNfsMetaName defines (NFS) Ganesha server CR name to use as shared storage for moodledata.
	// +kubebuilder:validation:MaxLength=63
	// +optional
	MoodleNfsMetaName string `json:"moodleNfsMetaName,omitempty"`

	// MoodleKeydbMetaName defines Keydb CR name to use as redis cache.
	// +kubebuilder:validation:MaxLength=63
	// +optional
	MoodleKeydbMetaName string `json:"moodleKeydbMetaName,omitempty"`

	// MoodleRedisSessionStore whether redis is configured as session store. Default: false
	// +optional
	MoodleRedisSessionStore bool `json:"moodleRedisSessionStore,omitempty"`

	// MoodleRedisMucStore whether redis is configured as MUC store. Default: false
	// +optional
	MoodleRedisMucStore bool `json:"moodleRedisMucStore,omitempty"`

	// MoodleRedisHost defines redis host. Default: '127.0.0.1'
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleRedisHost string `json:"moodleRedisHost,omitempty"`

	// MoodleRedisSecret defines redis auth secret name. Default: ''
	// +kubebuilder:validation:MaxLength=255
	// +optional
	MoodleRedisSecret string `json:"moodleRedisSecret,omitempty"`

	// MoodleRedisSecretAuthKey defines key inside auth secret name. Default: 'keydb_password'
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleRedisSecretAuthKey string `json:"moodleRedisSecretAuthKey,omitempty"`

	// MoodleConfigSessionRedisPrefix defines prefix for redis session. Default: ''
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleConfigSessionRedisPrefix string `json:"moodleConfigSessionRedisPrefix,omitempty"`

	// MoodleConfigSessionRedisSerializerUseIgbinary whether igbinary is used for redis session. Default: false
	// +optional
	MoodleConfigSessionRedisSerializerUseIgbinary bool `json:"moodleConfigSessionRedisSerializerUseIgbinary,omitempty"`

	// MoodleConfigSessionRedisCompressor defines redis session compresor
	// +optional
	MoodleConfigSessionRedisCompressor SessionRedisCompressor `json:"moodleConfigSessionRedisCompressor,omitempty"`

	// MoodleRedisMucStorePrefix defines prefix for redis MUC store. Default: ''
	// +kubebuilder:validation:MaxLength=100
	// +optional
	MoodleRedisMucStorePrefix string `json:"moodleRedisMucStorePrefix,omitempty"`

	// MoodleRedisMucStoreSerializer defines serializer for redis MUC store. Default: 1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=2
	// +optional
	MoodleRedisMucStoreSerializer int8 `json:"moodleRedisMucStoreSerializer,omitempty"`

	// MoodleRedisMucStoreCompressor defines compressor for redis MUC store. Default: 0
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=1
	// +optional
	MoodleRedisMucStoreCompressor int8 `json:"moodleRedisMucStoreCompressor,omitempty"`

	// RoutineStatusCrNotify specification using ansible URI module
	// +optional
	RoutineStatusCrNotify RoutineStatusCrNotify `json:"routineStatusCrNotify,omitempty"`
}

// StorageAccessMode describes storage access modes
// +kubebuilder:validation:Enum=ReadWriteOnce;ReadOnlyMany;ReadWriteMany
type StorageAccessMode string

// MoodleProtocol describes Moodle access protocol
// +kubebuilder:validation:Enum=http;https
type MoodleProtocol string

// MoodleConfigAdditionalCfg defines moodle extra config properties in config.php
type MoodleConfigProperty struct{}

// SessionRedisCompressor describes Moodle redis session compresor
// +kubebuilder:validation:Enum=none;gzip;zstd
type SessionRedisCompressor string

const (
	// ReadWriteOnce can be mounted as read-write by a single node
	ReadWriteOnce StorageAccessMode = "ReadWriteOnce"

	// ReadOnlyMany can be mounted read-only by many nodes
	ReadOnlyMany StorageAccessMode = "ReadOnlyMany"

	// ReadWriteMany the volume can be mounted as read-write by many nodes
	ReadWriteMany StorageAccessMode = "ReadWriteMany"
)

// RoutineStatusCrNotifyUUID used when notifying status to an endpoint
// +kubebuilder:validation:MinLength=36
// +kubebuilder:validation:MaxLength=36
// +optional
type RoutineStatusCrNotifyUUID string

// RoutineStatusCrNotify specification using ansible URI module
type RoutineStatusCrNotify struct {
	// HTTP or HTTPS URL in the form (http|https)://host.domain[:port]/path
	Url string `json:"url"`
	// StatusCode A list of valid, numeric, HTTP status codes that signifies success of the request.
	// +optional
	StatusCode []int8 `json:"statusCode,omitempty"`
	// Method The HTTP method of the request or response.
	// +kubebuilder:validation:Enum=GET;POST;PUT;PATCH;DELETE
	// +optional
	Method string `json:"method,omitempty"`
	// UUID used when notifying status to an endpoint
	// +optional
	UUID RoutineStatusCrNotifyUUID `json:"uuid,omitempty"`
}
