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
	corev1 "k8s.io/api/core/v1"
)

// MoodleSpec defines the desired state of Moodle
type MoodleSpec struct {
	// MoodleImage defines image for moodle container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	MoodleImage string `json:"moodleImage,omitempty"`

	// MoodleNewInstance whether new instance job runs
	// +optional
	MoodleNewInstance bool `json:"moodleNewInstance,omitempty"`

	// MoodleNewInstanceAgreeLicense whether agree to Moodle license. Required
	MoodleNewInstanceAgreeLicense bool `json:"moodleNewInstanceAgreeLicense"`

	// MoodleNewInstanceLang set moodle language code
	// +kubebuilder:validation:Pattern="^[a-z_]+$"
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=15
	// +optional
	MoodleNewInstanceLang string `json:"moodleNewInstanceLang"`

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

	// MoodleCronjobNodeSelector defines any node labels selectors for Moodle cronjob pods.
	// +optional
	MoodleCronjobNodeSelector string `json:"moodleCronjobNodeSelector,omitempty"`

	// MoodleCronjobAffinity defines any affinity rules for Moodle cronjob pods.
	// +optional
	MoodleCronjobAffinity string `json:"moodleCronjobAffinity,omitempty"`

	// MoodleCronjobResourceRequests whether moodle cronjob resource requests are added. Default: true
	// +optional
	MoodleCronjobResourceRequests bool `json:"moodleCronjobResourceRequests,omitempty"`

	// MoodleCronjobResourceRequestsCpu set moodle cronjob resource requests cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleCronjobResourceRequestsCpu string `json:"moodleCronjobResourceRequestsCpu,omitempty"`

	// MoodleCronjobResourceRequestsMemory set moodle cronjob resource requests memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleCronjobResourceRequestsMemory string `json:"moodleCronjobResourceRequestsMemory,omitempty"`

	// MoodleCronjobResourceLimits whether moodle cronjob resource limits are added. Default: false
	// +optional
	MoodleCronjobResourceLimits bool `json:"moodleCronjobResourceLimits,omitempty"`

	// MoodleCronjobResourceLimitsCpu set moodle cronjob resource limits cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleCronjobResourceLimitsCpu string `json:"moodleCronjobResourceLimitsCpu,omitempty"`

	// MoodleCronjobResourceLimitsMemory set moodle cronjob resource limits memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleCronjobResourceLimitsMemory string `json:"moodleCronjobResourceLimitsMemory,omitempty"`

	// MoodleCronjobVpaSpec set moodle cronjob vertical pod autoscaler spec
	// +optional
	MoodleCronjobVpaSpec string `json:"moodleCronjobVpaSpec,omitempty"`

	// MoodleUpdateJobTolerations defines any tolerations for Moodle cronjob pods.
	// +optional
	MoodleUpdateJobTolerations []corev1.Toleration `json:"moodleUpdateJobTolerations,omitempty"`

	// MoodleUpdateJobNodeSelector defines any node labels selectors for Moodle cronjob pods.
	// +optional
	MoodleUpdateJobNodeSelector string `json:"moodleUpdateJobNodeSelector,omitempty"`

	// MoodleUpdateJobAffinity defines any affinity rules for Moodle cronjob pods.
	// +optional
	MoodleUpdateJobAffinity string `json:"moodleUpdateJobAffinity,omitempty"`

	// MoodleUpdateJobResourceRequests whether moodle update job resource requests are added. Default: true
	// +optional
	MoodleUpdateJobResourceRequests bool `json:"moodleUpdateJobResourceRequests,omitempty"`

	// MoodleUpdateJobResourceRequestsCpu set moodle update job resource requests cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleUpdateJobResourceRequestsCpu string `json:"moodleUpdateJobResourceRequestsCpu,omitempty"`

	// MoodleUpdateJobResourceRequestsMemory set moodle update job resource requests memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleUpdateJobResourceRequestsMemory string `json:"moodleUpdateJobResourceRequestsMemory,omitempty"`

	// MoodleUpdateJobResourceLimits whether moodle update job resource limits are added. Default: false
	// +optional
	MoodleUpdateJobResourceLimits bool `json:"moodleUpdateJobResourceLimits,omitempty"`

	// MoodleUpdateJobResourceLimitsCpu set moodle update job resource limits cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleUpdateJobResourceLimitsCpu string `json:"moodleUpdateJobResourceLimitsCpu,omitempty"`

	// MoodleUpdateJobResourceLimitsMemory set moodle cronjob resource limits memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleUpdateJobResourceLimitsMemory string `json:"moodleUpdateJobResourceLimitsMemory,omitempty"`

	// MoodleNewInstanceJobTolerations defines any tolerations for Moodle cronjob pods.
	// +optional
	MoodleNewInstanceJobTolerations []corev1.Toleration `json:"moodleNewInstanceJobTolerations,omitempty"`

	// MoodleNewInstanceJobNodeSelector defines any node labels selectors for Moodle cronjob pods.
	// +optional
	MoodleNewInstanceJobNodeSelector string `json:"moodleNewInstanceJobNodeSelector,omitempty"`

	// MoodleNewInstanceJobAffinity defines any affinity rules for Moodle cronjob pods.
	// +optional
	MoodleNewInstanceJobAffinity string `json:"moodleNewInstanceJobAffinity,omitempty"`

	// MoodleNewInstanceJobResourceRequests whether moodle new instance job resource requests are added. Default: true
	// +optional
	MoodleNewInstanceJobResourceRequests bool `json:"moodleNewInstanceJobResourceRequests,omitempty"`

	// MoodleNewInstanceJobResourceRequestsCpu set moodle new instance job resource requests cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleNewInstanceJobResourceRequestsCpu string `json:"moodleNewInstanceJobResourceRequestsCpu,omitempty"`

	// MoodleNewInstanceJobResourceRequestsMemory set moodle new instance job resource requests memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleNewInstanceJobResourceRequestsMemory string `json:"moodleNewInstanceJobResourceRequestsMemory,omitempty"`

	// MoodleNewInstanceJobResourceLimits whether moodle new instance job resource limits are added. Default: false
	// +optional
	MoodleNewInstanceJobResourceLimits bool `json:"moodleNewInstanceJobResourceLimits,omitempty"`
	// MoodleNewInstanceJobResourceLimitsCpu set moodle new instance job resource limits cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleNewInstanceJobResourceLimitsCpu string `json:"moodleNewInstanceJobResourceLimitsCpu,omitempty"`
	// MoodleNewInstanceJobResourceLimitsMemory set moodle new instance job resource limits memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	MoodleNewInstanceJobResourceLimitsMemory string `json:"moodleNewInstanceJobResourceLimitsMemory,omitempty"`

	// MoodleConfigAdditionalCfg defines moodle extra config properties in config.php
	// +kubebuilder:pruning:PreserveUnknownFields
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

	// MoodleNetpolOmit whether to omit default moodle network policy. Default: true
	// +optional
	MoodleNetpolOmit bool `json:"moodleNetpolOmit,omitempty"`

	// MoodleNetpolIngressIpblock defines ingress ip block for moodle default network policy
	// +optional
	MoodleNetpolIngressIpblock string `json:"moodleNetpolIngressIpblock,omitempty"`

	// MoodleNetpolEgressIpblock defines egress ip block for moodle default network policy
	// +optional
	MoodleNetpolEgressIpblock string `json:"moodleNetpolEgressIpblock,omitempty"`

	// MoodleNetpolIngressExtraPorts defines extra ingress ports for moodle default network policy
	// +optional
	MoodleNetpolIngressExtraPorts []NetworkPolicyExtraPort `json:"moodleNetpolIngressExtraPorts,omitempty"`

	// MoodleNetpolEgressExtraPorts defines extra egress ports for moodle default network policy
	// +optional
	MoodleNetpolEgressExtraPorts []NetworkPolicyExtraPort `json:"moodleNetpolEgressExtraPorts,omitempty"`

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

	// NginxNodeSelector defines any node labels selectors for Nginx pods.
	// +optional
	NginxNodeSelector string `json:"nginxNodeSelector,omitempty"`

	// NginxAffinity defines any affinity rules for Nginx pods.
	// +optional
	NginxAffinity string `json:"nginxAffinity,omitempty"`

	// NginxExtraConfig contains extra Nginx config
	// +optional
	NginxExtraConfig string `json:"nginxExtraConfig,omitempty"`

	// NginxResourceRequests whether nginx resource requests are added. Default: true
	// +optional
	NginxResourceRequests bool `json:"nginxResourceRequests,omitempty"`

	// NginxResourceRequestsCpu set nginx resource requests cpu
	// +optional
	NginxResourceRequestsCpu string `json:"nginxResourceRequestsCpu,omitempty"`

	// NginxResourceRequestsMemory set nginx resource requests memory
	// +optional
	NginxResourceRequestsMemory string `json:"nginxResourceRequestsMemory,omitempty"`

	// NginxResourceLimits whether nginx resource limits are added. Default: false
	// +optional
	NginxResourceLimits bool `json:"nginxResourceLimits,omitempty"`

	// NginxResourceLimitsCpu set nginx resource limits cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	NginxResourceLimitsCpu string `json:"nginxResourceLimitsCpu,omitempty"`

	// NginxResourceLimitsMemory set nginx resource limits memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	NginxResourceLimitsMemory string `json:"nginxResourceLimitsMemory,omitempty"`

	// NginxHpaSpec set nginx horizontal pod autoscaler spec
	// +optional
	NginxHpaSpec string `json:"nginxHpaSpec,omitempty"`

	// NginxVpaSpec set nginx vertical pod autoscaler spec
	// +optional
	NginxVpaSpec string `json:"nginxVpaSpec,omitempty"`

	// NginxNetpolOmit whether to omit default network policy for nginx. Default: true
	// +optional
	NginxNetpolOmit bool `json:"nginxNetpolOmit,omitempty"`

	// NginxNetpolIngressIpblock defines ingress ip block for nginx default network policy
	// +optional
	NginxNetpolIngressIpblock string `json:"nginxNetpolIngressIpblock,omitempty"`

	// NginxNetpolEgressIpblock defines egress ip block for nginx default network policy
	// +optional
	NginxNetpolEgressIpblock string `json:"nginxNetpolEgressIpblock,omitempty"`

	// NginxNetpolIngressExtraPorts defines extra ingress ports for nginx default network policy
	// +optional
	NginxNetpolIngressExtraPorts []NetworkPolicyExtraPort `json:"nginxNetpolIngressExtraPorts,omitempty"`

	// NginxNetpolEgressExtraPorts defines extra egress ports for nginx default network policy
	// +optional
	NginxNetpolEgressExtraPorts []NetworkPolicyExtraPort `json:"nginxNetpolEgressExtraPorts,omitempty"`

	// PhpFpmSize defines php-fpm number of replicas between 0 and 255
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=255
	// +optional
	PhpFpmSize int32 `json:"phpFpmSize,omitempty"`

	// PhpFpmImage defines image for php-fpm container
	// +kubebuilder:validation:MaxLength=255
	// +optional
	PhpFpmImage string `json:"phpFpmImage,omitempty"`

	// PhpFpmIngressAnnotations defines php-fpm annotations
	// +optional
	PhpFpmIngressAnnotations string `json:"phpFpmIngressAnnotations,omitempty"`

	// PhpFpmTolerations defines any tolerations for php-fpm pods.
	// +optional
	PhpFpmTolerations []corev1.Toleration `json:"phpFpmTolerations,omitempty"`

	// PhpFpmNodeSelector defines any node labels selectors for PhpFpm pods.
	// +optional
	PhpFpmNodeSelector string `json:"phpFpmNodeSelector,omitempty"`

	// PhpFpmAffinity defines any affinity rules for PhpFpm pods.
	// +optional
	PhpFpmAffinity string `json:"phpFpmAffinity,omitempty"`

	// PhpFpmPhpExtraIni contains extra php ini config
	// +optional
	PhpFpmPhpExtraIni string `json:"phpFpmPhpExtraIni,omitempty"`

	// PhpFpmExtraConfig contains extra php-fpm config
	// +optional
	PhpFpmExtraConfig string `json:"phpFpmExtraConfig,omitempty"`

	// PhpFpmResourceRequests whether php-fpm resource requests are added. Default: true
	// +optional
	PhpFpmResourceRequests bool `json:"phpFpmResourceRequests,omitempty"`

	// PhpFpmResourceRequestsCpu set php-fpm resource requests cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PhpFpmResourceRequestsCpu string `json:"phpFpmResourceRequestsCpu,omitempty"`

	// PhpFpmResourceRequestsMemory set php-fpm resource requests memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PhpFpmResourceRequestsMemory string `json:"phpFpmResourceRequestsMemory,omitempty"`

	// PhpFpmResourceLimits whether php-fpm resource limits are added. Default: false
	// +optional
	PhpFpmResourceLimits bool `json:"phpFpmResourceLimits,omitempty"`

	// PhpFpmResourceLimitsCpu set php-fpm resource limits cpu
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PhpFpmResourceLimitsCpu string `json:"phpFpmResourceLimitsCpu,omitempty"`

	// PhpFpmResourceLimitsMemory set php-fpm resource limits memory
	// +kubebuilder:validation:MaxLength=20
	// +optional
	PhpFpmResourceLimitsMemory string `json:"phpFpmResourceLimitsMemory,omitempty"`

	// PhpFpmHpaSpec set php-fpm horizontal pod autoscaler spec
	// +optional
	PhpFpmHpaSpec string `json:"phpFpmHpaSpec,omitempty"`

	// PhpFpmVpaSpec set php-fpm vertical pod autoscaler spec
	// +optional
	PhpFpmVpaSpec string `json:"phpFpmVpaSpec,omitempty"`

	// PhpFpmNetpolOmit whether to omit default network policy for php-fpm. Default: true
	// +optional
	PhpFpmNetpolOmit bool `json:"phpFpmNetpolOmit,omitempty"`

	// PhpFpmNetpolIngressIpblock defines ingress ip block for php-fpm default network policy
	// +optional
	PhpFpmNetpolIngressIpblock string `json:"phpFpmNetpolIngressIpblock,omitempty"`

	// PhpFpmNetpolEgressIpblock defines egress ip block for php-fpm default network policy
	// +optional
	PhpFpmNetpolEgressIpblock string `json:"phpFpmNetpolEgressIpblock,omitempty"`

	// PhpFpmNetpolIngressExtraPorts defines extra ingress ports for php-fpm default network policy
	// +optional
	PhpFpmNetpolIngressExtraPorts []NetworkPolicyExtraPort `json:"phpFpmNetpolIngressExtraPorts,omitempty"`

	// PhpFpmNetpolEgressExtraPorts defines extra egress ports for php-fpm default network policy
	// +optional
	PhpFpmNetpolEgressExtraPorts []NetworkPolicyExtraPort `json:"phpFpmNetpolEgressExtraPorts,omitempty"`

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

	// RoutineStatusCrNotifyTermination specification using ansible URI module
	// +optional
	RoutineStatusCrNotifyTermination RoutineStatusCrNotify `json:"routineStatusCrNotifyTermination,omitempty"`
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

// RoutineStatusCrNotifyHeaders used when notifying status to an endpoint
// +optional
type RoutineStatusCrNotifyHeaders struct{}

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
	// Headers used when notifying status to an endpoint
	// +optional
	Headers RoutineStatusCrNotifyHeaders `json:"headers,omitempty"`
	// JwtSecretEnvName environment variable name that holds secret to generate jwt tokens
	// +optional
	JwtSecretEnvName string `json:"jwtSecretEnvName,omitempty"`
}

type NetworkPolicyExtraPort struct {
	// Port number
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`
	// Protocol TCP or UDP
	// +kubebuilder:validation:Enum=TCP;UDP
	// +optional
	Protocol string `json:"protocol,omitempty"`
}
