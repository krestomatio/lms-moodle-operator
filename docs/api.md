# API Reference

## Packages
- [lms.krestomat.io/v1alpha1](#lmskrestomatiov1alpha1)


## lms.krestomat.io/v1alpha1

Package v1alpha1 contains API Schema definitions for the lms v1alpha1 API group

### Resource Types
- [LMSMoodle](#lmsmoodle)
- [LMSMoodleList](#lmsmoodlelist)
- [LMSMoodleTemplate](#lmsmoodletemplate)
- [LMSMoodleTemplateList](#lmsmoodletemplatelist)



#### KeydbMode

_Underlying type:_ _string_

KeydbMode describes mode keydb runs

_Validation:_
- Enum: [standalone multimaster custom]

_Appears in:_
- [KeydbSpec](#keydbspec)



#### KeydbSpec



KeydbSpec defines the desired state of Keydb



_Appears in:_
- [LMSMoodleSpec](#lmsmoodlespec)
- [LMSMoodleTemplateSpec](#lmsmoodletemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `keydbMode` _[KeydbMode](#keydbmode)_ | KeydbMode describes mode keydb runs |  | Enum: [standalone multimaster custom] <br /> |
| `keydbExtraConfig` _string_ | KeydbExtraConfig contains extra keydb config |  |  |
| `keydbSize` _integer_ | KeydbSize defines keydb number of replicas |  |  |
| `keydbImage` _string_ | KeydbImage defines image for keydb container |  | MaxLength: 255 <br /> |
| `keydbPvcDataSize` _string_ | KeydbPvcDataSize defines keydb storage size |  | MaxLength: 20 <br />MinLength: 2 <br /> |
| `keydbPvcDataStorageAccessMode` _[StorageAccessMode](#storageaccessmode)_ | KeydbPvcDataStorageAccessMode defines keydb storage access modes |  | Enum: [ReadWriteOnce ReadOnlyMany ReadWriteMany] <br /> |
| `keydbPvcDataStorageClassName` _string_ | KeydbPvcDataStorageClassName defines keydb storage class |  | MaxLength: 63 <br />MinLength: 2 <br /> |
| `keydbPvcDataAutoexpansion` _boolean_ | KeydbPvcDataAutoexpansion enables autoexpansion |  |  |
| `keydbPvcDataAutoexpansionIncrementGib` _integer_ | KeydbPvcDataAutoexpansionIncrementGib defines Gib to increment |  |  |
| `keydbPvcDataAutoexpansionCapGib` _integer_ | KeydbPvcDataAutoexpansionCapGib defines limit for autoexpansion increments |  |  |
| `keydbResourceRequests` _boolean_ | KeydbResourceRequests whether keydb resource requests are added. Default: true |  |  |
| `keydbResourceRequestsCpu` _string_ | KeydbResourceRequestsCpu set keydb resource requests cpu |  | MaxLength: 20 <br /> |
| `keydbResourceRequestsMemory` _string_ | KeydbResourceRequestsMemory set keydb resource requests memory |  | MaxLength: 20 <br /> |
| `keydbResourceLimits` _boolean_ | KeydbResourceLimits whether keydb resource limits are added. Default: false |  |  |
| `keydbResourceLimitsCpu` _string_ | KeydbResourceLimitsCpu set keydb resource limits cpu |  | MaxLength: 20 <br /> |
| `keydbResourceLimitsMemory` _string_ | KeydbResourceLimitsMemory set keydb resource limits memory |  | MaxLength: 20 <br /> |
| `keydbTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | KeydbTolerations defines any tolerations for Keydb pods. |  |  |
| `keydbNodeSelector` _string_ | KeydbNodeSelector defines any node labels selectors for Keydb pods. |  |  |
| `keydbAffinity` _string_ | KeydbAffinity defines any affinity rules for Keydb pods. |  |  |
| `keydbVpaSpec` _string_ | KeydbVpaSpec set keydb horizontal pod autoscaler spec |  |  |
| `keydbNetpolOmit` _boolean_ | KeydbNetpolOmit whether to omit default keydb network policy. Default: true |  |  |
| `keydbNetpolIngressIpblock` _string_ | GaneshaNetpolIngressIpblock defines ingress ip block for keydb default network policy |  |  |
| `keydbNetpolEgressIpblock` _string_ | KeydbNetpolEgressIpblock defines egress ip block for keydb default network policy |  |  |


#### LMSMoodle



LMSMoodle is the Schema for the lmsmoodles API



_Appears in:_
- [LMSMoodleList](#lmsmoodlelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `lms.krestomat.io/v1alpha1` | | |
| `kind` _string_ | `LMSMoodle` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[LMSMoodleSpec](#lmsmoodlespec)_ |  |  |  |
| `status` _[LMSMoodleStatus](#lmsmoodlestatus)_ |  |  |  |


#### LMSMoodleList



LMSMoodleList contains a list of LMSMoodle





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `lms.krestomat.io/v1alpha1` | | |
| `kind` _string_ | `LMSMoodleList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[LMSMoodle](#lmsmoodle) array_ |  |  |  |


#### LMSMoodleSpec



LMSMoodleSpec defines the desired state of LMSMoodle



_Appears in:_
- [LMSMoodle](#lmsmoodle)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `lmsMoodleTemplateName` _string_ | LMSMoodleTemplateName defines what LMS Moodle template to use |  | MaxLength: 255 <br />MinLength: 1 <br /> |
| `lmsMoodleNetpolOmit` _boolean_ | LMSMoodleNetpolOmit whether to omit default network policy for the namespace. Default: false<br />It will deny all ingress and egress traffic to the namespace<br />Intended to be used with custom network policies already in place or<br />by not omitting default network policies of each dependant resource |  |  |
| `desiredState` _string_ | DesiredState defines the desired state to put a LMSMoodle | Ready | Enum: [Ready Suspended] <br /> |
| `moodleSpec` _[MoodleSpec](#moodlespec)_ | MoodleSpec defines Moodle spec |  |  |
| `postgresSpec` _[PostgresSpec](#postgresspec)_ | PostgresSpec defines Postgres spec to deploy optionally |  |  |
| `nfsSpec` _[NfsSpec](#nfsspec)_ | NfsSpec defines (NFS) Ganesha server spec to deploy optionally |  |  |
| `keydbSpec` _[KeydbSpec](#keydbspec)_ | KeydbSpec defines Keydb spec to deploy optionally |  |  |


#### LMSMoodleStatus



LMSMoodleStatus defines the observed state of LMSMoodle



_Appears in:_
- [LMSMoodle](#lmsmoodle)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `conditions` _[Condition](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#condition-v1-meta) array_ | Conditions represent the latest available observations of the resource state |  |  |
| `state` _string_ | state describes the LMSMoodle state | Unknown |  |
| `url` _string_ | Url defines LMSMoodle url |  |  |
| `storageGb` _string_ | StorageGb defines LMSMoodle number of current GB for storage capacity | 0 |  |
| `registeredUsers` _integer_ | RegisteredUsers defines LMSMoodle number of current registered users for user capacity | 0 |  |
| `release` _string_ | Release defines LMSMoodle moodle version |  |  |


#### LMSMoodleTemplate



LMSMoodleTemplate is the Schema for the lmsmoodletemplates API



_Appears in:_
- [LMSMoodleTemplateList](#lmsmoodletemplatelist)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `lms.krestomat.io/v1alpha1` | | |
| `kind` _string_ | `LMSMoodleTemplate` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `spec` _[LMSMoodleTemplateSpec](#lmsmoodletemplatespec)_ |  |  |  |
| `status` _[LMSMoodleTemplateStatus](#lmsmoodletemplatestatus)_ |  |  |  |


#### LMSMoodleTemplateList



LMSMoodleTemplateList contains a list of Moodle Template





| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `apiVersion` _string_ | `lms.krestomat.io/v1alpha1` | | |
| `kind` _string_ | `LMSMoodleTemplateList` | | |
| `kind` _string_ | Kind is a string value representing the REST resource this object represents.<br />Servers may infer this from the endpoint the client submits requests to.<br />Cannot be updated.<br />In CamelCase.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds |  |  |
| `apiVersion` _string_ | APIVersion defines the versioned schema of this representation of an object.<br />Servers should convert recognized schemas to the latest internal value, and<br />may reject unrecognized values.<br />More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources |  |  |
| `metadata` _[ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#listmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |  |  |
| `items` _[LMSMoodleTemplate](#lmsmoodletemplate) array_ |  |  |  |


#### LMSMoodleTemplateSpec



LMSMoodleTemplateSpec defines the desired state of LMSMoodleTemplate



_Appears in:_
- [LMSMoodleSpec](#lmsmoodlespec)
- [LMSMoodleTemplate](#lmsmoodletemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `moodleSpec` _[MoodleSpec](#moodlespec)_ | MoodleSpec defines Moodle spec |  |  |
| `postgresSpec` _[PostgresSpec](#postgresspec)_ | PostgresSpec defines Postgres spec to deploy optionally |  |  |
| `nfsSpec` _[NfsSpec](#nfsspec)_ | NfsSpec defines (NFS) Ganesha server spec to deploy optionally |  |  |
| `keydbSpec` _[KeydbSpec](#keydbspec)_ | KeydbSpec defines Keydb spec to deploy optionally |  |  |


#### LMSMoodleTemplateStatus



LMSMoodleTemplateStatus defines the observed state of LMSMoodleTemplate



_Appears in:_
- [LMSMoodleTemplate](#lmsmoodletemplate)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `state` _string_ | state describes the LMSMoodleTemplate state | Unknown |  |


#### MoodleConfigProperty



MoodleConfigAdditionalCfg defines moodle extra config properties in config.php



_Appears in:_
- [MoodleSpec](#moodlespec)



#### MoodleProtocol

_Underlying type:_ _string_

MoodleProtocol describes Moodle access protocol

_Validation:_
- Enum: [http https]

_Appears in:_
- [MoodleSpec](#moodlespec)



#### MoodleSpec



MoodleSpec defines the desired state of Moodle



_Appears in:_
- [LMSMoodleSpec](#lmsmoodlespec)
- [LMSMoodleTemplateSpec](#lmsmoodletemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `moodleImage` _string_ | MoodleImage defines image for moodle container |  | MaxLength: 255 <br /> |
| `moodleNewInstance` _boolean_ | MoodleNewInstance whether new instance job runs |  |  |
| `moodleNewInstanceAgreeLicense` _boolean_ | MoodleNewInstanceAgreeLicense whether agree to Moodle license. Required |  |  |
| `moodleNewInstanceLang` _string_ | MoodleNewInstanceLang set moodle language code |  | MaxLength: 15 <br />MinLength: 2 <br />Pattern: `^[a-z_]+$` <br /> |
| `moodleNewInstanceFullname` _string_ |  |  | MaxLength: 100 <br /> |
| `moodleNewInstanceShortname` _string_ |  |  | MaxLength: 100 <br /> |
| `moodleNewInstanceSummary` _string_ |  |  | MaxLength: 300 <br /> |
| `moodleNewInstanceAdminuser` _string_ |  |  | MaxLength: 100 <br />MinLength: 1 <br /> |
| `moodleNewInstanceAdminmail` _string_ | MoodleNewInstanceAdminMail is the admin email to set in new instance. Required |  | MaxLength: 100 <br />MinLength: 3 <br /> |
| `moodleNewAdminpassHash` _string_ | MoodleNewAdminPassHash is the bcrypt compatible admin password to set in new instance. Required |  | MaxLength: 60 <br />MinLength: 60 <br />Pattern: `^\$2[ayb]\$.{56}$` <br /> |
| `moodlePvcDataSize` _string_ | MoodlePvcDataSize defines moodledata storage size |  | MaxLength: 100 <br />MinLength: 2 <br /> |
| `moodlePvcDataStorageAccessMode` _[StorageAccessMode](#storageaccessmode)_ | MoodlePvcDataStorageAccessMode defines moodledata storage access modes |  | Enum: [ReadWriteOnce ReadOnlyMany ReadWriteMany] <br /> |
| `moodlePvcDataStorageClassName` _string_ | MoodlePvcDataStorageClassName defines moodledata storage class |  | MaxLength: 63 <br />MinLength: 2 <br /> |
| `moodleHost` _string_ | MoodleHost defines Moodle host for url |  | MaxLength: 100 <br />MinLength: 2 <br /> |
| `moodlePort` _integer_ | MoodlePort defines Moodle port for url |  | Maximum: 65535 <br />Minimum: 1 <br /> |
| `moodleSubpath` _string_ | MoodleSubpath defines Moodle subpath for url |  | MaxLength: 100 <br />MinLength: 2 <br /> |
| `moodleHealthcheckSubpath` _string_ | MoodleHealthcheckSubpath defines Moodle subpath for nginx check |  | MaxLength: 100 <br />MinLength: 2 <br /> |
| `moodleProtocol` _[MoodleProtocol](#moodleprotocol)_ | MoodleProtocol whether to use http or https |  | Enum: [http https] <br /> |
| `moodleCronjobTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | MoodleCronjobTolerations defines any tolerations for Moodle cronjob pods. |  |  |
| `moodleCronjobNodeSelector` _string_ | MoodleCronjobNodeSelector defines any node labels selectors for Moodle cronjob pods. |  |  |
| `moodleCronjobAffinity` _string_ | MoodleCronjobAffinity defines any affinity rules for Moodle cronjob pods. |  |  |
| `moodleCronjobResourceRequests` _boolean_ | MoodleCronjobResourceRequests whether moodle cronjob resource requests are added. Default: true |  |  |
| `moodleCronjobResourceRequestsCpu` _string_ | MoodleCronjobResourceRequestsCpu set moodle cronjob resource requests cpu |  | MaxLength: 20 <br /> |
| `moodleCronjobResourceRequestsMemory` _string_ | MoodleCronjobResourceRequestsMemory set moodle cronjob resource requests memory |  | MaxLength: 20 <br /> |
| `moodleCronjobResourceLimits` _boolean_ | MoodleCronjobResourceLimits whether moodle cronjob resource limits are added. Default: false |  |  |
| `moodleCronjobResourceLimitsCpu` _string_ | MoodleCronjobResourceLimitsCpu set moodle cronjob resource limits cpu |  | MaxLength: 20 <br /> |
| `moodleCronjobResourceLimitsMemory` _string_ | MoodleCronjobResourceLimitsMemory set moodle cronjob resource limits memory |  | MaxLength: 20 <br /> |
| `moodleCronjobVpaSpec` _string_ | MoodleCronjobVpaSpec set moodle cronjob vertical pod autoscaler spec |  |  |
| `moodleUpdateJobTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | MoodleUpdateJobTolerations defines any tolerations for Moodle cronjob pods. |  |  |
| `moodleUpdateJobNodeSelector` _string_ | MoodleUpdateJobNodeSelector defines any node labels selectors for Moodle cronjob pods. |  |  |
| `moodleUpdateJobAffinity` _string_ | MoodleUpdateJobAffinity defines any affinity rules for Moodle cronjob pods. |  |  |
| `moodleUpdateJobResourceRequests` _boolean_ | MoodleUpdateJobResourceRequests whether moodle update job resource requests are added. Default: true |  |  |
| `moodleUpdateJobResourceRequestsCpu` _string_ | MoodleUpdateJobResourceRequestsCpu set moodle update job resource requests cpu |  | MaxLength: 20 <br /> |
| `moodleUpdateJobResourceRequestsMemory` _string_ | MoodleUpdateJobResourceRequestsMemory set moodle update job resource requests memory |  | MaxLength: 20 <br /> |
| `moodleUpdateJobResourceLimits` _boolean_ | MoodleUpdateJobResourceLimits whether moodle update job resource limits are added. Default: false |  |  |
| `moodleUpdateJobResourceLimitsCpu` _string_ | MoodleUpdateJobResourceLimitsCpu set moodle update job resource limits cpu |  | MaxLength: 20 <br /> |
| `moodleUpdateJobResourceLimitsMemory` _string_ | MoodleUpdateJobResourceLimitsMemory set moodle cronjob resource limits memory |  | MaxLength: 20 <br /> |
| `moodleNewInstanceJobTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | MoodleNewInstanceJobTolerations defines any tolerations for Moodle cronjob pods. |  |  |
| `moodleNewInstanceJobNodeSelector` _string_ | MoodleNewInstanceJobNodeSelector defines any node labels selectors for Moodle cronjob pods. |  |  |
| `moodleNewInstanceJobAffinity` _string_ | MoodleNewInstanceJobAffinity defines any affinity rules for Moodle cronjob pods. |  |  |
| `moodleNewInstanceJobResourceRequests` _boolean_ | MoodleNewInstanceJobResourceRequests whether moodle new instance job resource requests are added. Default: true |  |  |
| `moodleNewInstanceJobResourceRequestsCpu` _string_ | MoodleNewInstanceJobResourceRequestsCpu set moodle new instance job resource requests cpu |  | MaxLength: 20 <br /> |
| `moodleNewInstanceJobResourceRequestsMemory` _string_ | MoodleNewInstanceJobResourceRequestsMemory set moodle new instance job resource requests memory |  | MaxLength: 20 <br /> |
| `moodleNewInstanceJobResourceLimits` _boolean_ | MoodleNewInstanceJobResourceLimits whether moodle new instance job resource limits are added. Default: false |  |  |
| `moodleNewInstanceJobResourceLimitsCpu` _string_ | MoodleNewInstanceJobResourceLimitsCpu set moodle new instance job resource limits cpu |  | MaxLength: 20 <br /> |
| `moodleNewInstanceJobResourceLimitsMemory` _string_ | MoodleNewInstanceJobResourceLimitsMemory set moodle new instance job resource limits memory |  | MaxLength: 20 <br /> |
| `moodleConfigAdditionalCfg` _[MoodleConfigProperty](#moodleconfigproperty)_ | MoodleConfigAdditionalCfg defines moodle extra config properties in config.php |  |  |
| `moodleConfigAdditionalBlock` _string_ | MoodleConfigAdditionalBlock defines moodle extra block in config.php |  |  |
| `moodleUpdateMinor` _boolean_ | MoodleUpdateMinor whether minor updates are automatically applied. Default: true |  |  |
| `moodleUpdateMajor` _boolean_ | MoodleUpdateMajor whether major updates are automatically applied. Default: false |  |  |
| `moodleNetpolOmit` _boolean_ | MoodleNetpolOmit whether to omit default moodle network policy. Default: true |  |  |
| `moodleNetpolIngressIpblock` _string_ | MoodleNetpolIngressIpblock defines ingress ip block for moodle default network policy |  |  |
| `moodleNetpolEgressIpblock` _string_ | MoodleNetpolEgressIpblock defines egress ip block for moodle default network policy |  |  |
| `nginxSize` _integer_ | NginxSize defines nginx number of replicas between 0 and 255 |  | Maximum: 255 <br />Minimum: 0 <br /> |
| `nginxImage` _string_ | NginxImage defines image for nginx container |  | MaxLength: 255 <br /> |
| `nginxIngressAnnotations` _string_ | NginxIngressAnnotations defines nginx annotations |  |  |
| `nginxTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | NginxTolerations defines any tolerations for Nginx pods. |  |  |
| `nginxNodeSelector` _string_ | NginxNodeSelector defines any node labels selectors for Nginx pods. |  |  |
| `nginxAffinity` _string_ | NginxAffinity defines any affinity rules for Nginx pods. |  |  |
| `nginxExtraConfig` _string_ | NginxExtraConfig contains extra Nginx config |  |  |
| `nginxResourceRequests` _boolean_ | NginxResourceRequests whether nginx resource requests are added. Default: true |  |  |
| `nginxResourceRequestsCpu` _string_ | NginxResourceRequestsCpu set nginx resource requests cpu |  |  |
| `nginxResourceRequestsMemory` _string_ | NginxResourceRequestsMemory set nginx resource requests memory |  |  |
| `nginxResourceLimits` _boolean_ | NginxResourceLimits whether nginx resource limits are added. Default: false |  |  |
| `nginxResourceLimitsCpu` _string_ | NginxResourceLimitsCpu set nginx resource limits cpu |  | MaxLength: 20 <br /> |
| `nginxResourceLimitsMemory` _string_ | NginxResourceLimitsMemory set nginx resource limits memory |  | MaxLength: 20 <br /> |
| `nginxHpaSpec` _string_ | NginxHpaSpec set nginx horizontal pod autoscaler spec |  |  |
| `nginxVpaSpec` _string_ | NginxVpaSpec set nginx vertical pod autoscaler spec |  |  |
| `nginxNetpolOmit` _boolean_ | NginxNetpolOmit whether to omit default network policy for nginx. Default: true |  |  |
| `nginxNetpolIngressIpblock` _string_ | NginxNetpolIngressIpblock defines ingress ip block for nginx default network policy |  |  |
| `nginxNetpolEgressIpblock` _string_ | NginxNetpolEgressIpblock defines egress ip block for nginx default network policy |  |  |
| `phpFpmSize` _integer_ | PhpFpmSize defines php-fpm number of replicas between 0 and 255 |  | Maximum: 255 <br />Minimum: 0 <br /> |
| `phpFpmImage` _string_ | PhpFpmImage defines image for php-fpm container |  | MaxLength: 255 <br /> |
| `phpFpmIngressAnnotations` _string_ | PhpFpmIngressAnnotations defines php-fpm annotations |  |  |
| `phpFpmTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | PhpFpmTolerations defines any tolerations for php-fpm pods. |  |  |
| `phpFpmNodeSelector` _string_ | PhpFpmNodeSelector defines any node labels selectors for PhpFpm pods. |  |  |
| `phpFpmAffinity` _string_ | PhpFpmAffinity defines any affinity rules for PhpFpm pods. |  |  |
| `phpFpmPhpExtraIni` _string_ | PhpFpmPhpExtraIni contains extra php ini config |  |  |
| `phpFpmExtraConfig` _string_ | PhpFpmExtraConfig contains extra php-fpm config |  |  |
| `phpFpmResourceRequests` _boolean_ | PhpFpmResourceRequests whether php-fpm resource requests are added. Default: true |  |  |
| `phpFpmResourceRequestsCpu` _string_ | PhpFpmResourceRequestsCpu set php-fpm resource requests cpu |  | MaxLength: 20 <br /> |
| `phpFpmResourceRequestsMemory` _string_ | PhpFpmResourceRequestsMemory set php-fpm resource requests memory |  | MaxLength: 20 <br /> |
| `phpFpmResourceLimits` _boolean_ | PhpFpmResourceLimits whether php-fpm resource limits are added. Default: false |  |  |
| `phpFpmResourceLimitsCpu` _string_ | PhpFpmResourceLimitsCpu set php-fpm resource limits cpu |  | MaxLength: 20 <br /> |
| `phpFpmResourceLimitsMemory` _string_ | PhpFpmResourceLimitsMemory set php-fpm resource limits memory |  | MaxLength: 20 <br /> |
| `phpFpmHpaSpec` _string_ | PhpFpmHpaSpec set php-fpm horizontal pod autoscaler spec |  |  |
| `phpFpmVpaSpec` _string_ | PhpFpmVpaSpec set php-fpm vertical pod autoscaler spec |  |  |
| `phpFpmNetpolOmit` _boolean_ | PhpFpmNetpolOmit whether to omit default network policy for php-fpm. Default: true |  |  |
| `phpFpmNetpolIngressIpblock` _string_ | PhpFpmNetpolIngressIpblock defines ingress ip block for php-fpm default network policy |  |  |
| `phpFpmNetpolEgressIpblock` _string_ | PhpFpmNetpolEgressIpblock defines egress ip block for php-fpm default network policy |  |  |
| `moodlePostgresMetaName` _string_ | MoodlePostgresMetaName defines Postgres CR name to use as database. |  | MaxLength: 63 <br /> |
| `moodleNfsMetaName` _string_ | MoodleNfsMetaName defines (NFS) Ganesha server CR name to use as shared storage for moodledata. |  | MaxLength: 63 <br /> |
| `moodleKeydbMetaName` _string_ | MoodleKeydbMetaName defines Keydb CR name to use as redis cache. |  | MaxLength: 63 <br /> |
| `moodleRedisSessionStore` _boolean_ | MoodleRedisSessionStore whether redis is configured as session store. Default: false |  |  |
| `moodleRedisMucStore` _boolean_ | MoodleRedisMucStore whether redis is configured as MUC store. Default: false |  |  |
| `moodleRedisHost` _string_ | MoodleRedisHost defines redis host. Default: '127.0.0.1' |  | MaxLength: 100 <br /> |
| `moodleRedisSecret` _string_ | MoodleRedisSecret defines redis auth secret name. Default: '' |  | MaxLength: 255 <br /> |
| `moodleRedisSecretAuthKey` _string_ | MoodleRedisSecretAuthKey defines key inside auth secret name. Default: 'keydb_password' |  | MaxLength: 100 <br /> |
| `moodleConfigSessionRedisPrefix` _string_ | MoodleConfigSessionRedisPrefix defines prefix for redis session. Default: '' |  | MaxLength: 100 <br /> |
| `moodleConfigSessionRedisSerializerUseIgbinary` _boolean_ | MoodleConfigSessionRedisSerializerUseIgbinary whether igbinary is used for redis session. Default: false |  |  |
| `moodleConfigSessionRedisCompressor` _[SessionRedisCompressor](#sessionrediscompressor)_ | MoodleConfigSessionRedisCompressor defines redis session compresor |  | Enum: [none gzip zstd] <br /> |
| `moodleRedisMucStorePrefix` _string_ | MoodleRedisMucStorePrefix defines prefix for redis MUC store. Default: '' |  | MaxLength: 100 <br /> |
| `moodleRedisMucStoreSerializer` _integer_ | MoodleRedisMucStoreSerializer defines serializer for redis MUC store. Default: 1 |  | Maximum: 2 <br />Minimum: 1 <br /> |
| `moodleRedisMucStoreCompressor` _integer_ | MoodleRedisMucStoreCompressor defines compressor for redis MUC store. Default: 0 |  | Maximum: 1 <br />Minimum: 0 <br /> |
| `routineStatusCrNotify` _[RoutineStatusCrNotify](#routinestatuscrnotify)_ | RoutineStatusCrNotify specification using ansible URI module |  |  |
| `routineStatusCrNotifyTermination` _[RoutineStatusCrNotify](#routinestatuscrnotify)_ | RoutineStatusCrNotifyTermination specification using ansible URI module |  |  |


#### NfsSpec



NfsSpec defines the desired state of Nfs



_Appears in:_
- [LMSMoodleSpec](#lmsmoodlespec)
- [LMSMoodleTemplateSpec](#lmsmoodletemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `ganeshaImage` _string_ | GaneshaImage defines image for ganesha server container |  | MaxLength: 255 <br /> |
| `ganeshaPvcDataSize` _string_ | GaneshaPvcDataSize defines ganesha server storage size |  | MaxLength: 20 <br />MinLength: 2 <br /> |
| `ganeshaPvcDataStorageAccessMode` _[StorageAccessMode](#storageaccessmode)_ | GaneshaPvcDataStorageAccessMode defines ganesha server storage access modes |  | Enum: [ReadWriteOnce ReadOnlyMany ReadWriteMany] <br /> |
| `ganeshaPvcDataStorageClassName` _string_ | GaneshaPvcDataStorageClassName defines ganesha server storage class |  | MaxLength: 63 <br />MinLength: 2 <br /> |
| `ganeshaResourceRequests` _boolean_ | GaneshaResourceRequests whether ganesha resource requests are added. Default: true |  |  |
| `ganeshaResourceRequestsCpu` _string_ | GaneshaResourceRequestsCpu set ganesha resource requests cpu |  | MaxLength: 20 <br /> |
| `ganeshaResourceRequestsMemory` _string_ | GaneshaResourceRequestsMemory set ganesha resource requests memory |  | MaxLength: 20 <br /> |
| `ganeshaResourceLimits` _boolean_ | GaneshaResourceLimits whether ganesha resource limits are added. Default: false |  |  |
| `ganeshaResourceLimitsCpu` _string_ | GaneshaResourceLimitsCpu set ganesha resource limits cpu |  | MaxLength: 20 <br /> |
| `ganeshaResourceLimitsMemory` _string_ | GaneshaResourceLimitsMemory set ganesha resource limits memory |  | MaxLength: 20 <br /> |
| `ganeshaTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | GaneshaTolerations defines any tolerations for Ganesha server pods. |  |  |
| `ganeshaNodeSelector` _string_ | GaneshaNodeSelector defines any node labels selectors for Ganesha pods. |  |  |
| `ganeshaAffinity` _string_ | GaneshaAffinity defines any affinity rules for Ganesha pods. |  |  |
| `ganeshaExportUserid` _integer_ | GaneshaExportUserid defines export folder userid |  |  |
| `ganeshaExportGroupid` _integer_ | GaneshaExportGroupid defines export folder groupid |  |  |
| `ganeshaExportMode` _string_ | GaneshaExportMode defines folder permissions mode |  | Pattern: `[0-7]{4}` <br /> |
| `ganeshaPvcDataAutoexpansion` _boolean_ | GaneshaPvcDataAutoexpansion enables autoexpansion |  |  |
| `ganeshaPvcDataAutoexpansionIncrementGib` _integer_ | GaneshaPvcDataAutoexpansionIncrementGib defines Gib to increment |  |  |
| `ganeshaPvcDataAutoexpansionCapGib` _integer_ | GaneshaPvcDataAutoexpansionCapGib defines limit for autoexpansion increments |  |  |
| `ganeshaExtraBlockConfig` _string_ | GaneshaExtraBlockConfig contains extra block in ganesha server ganesha config |  |  |
| `ganeshaConfLogLevel` _string_ | GaneshaConfLogLevel defines nfs log level. Default: EVENT |  | Enum: [NULL FATAL MAJ CRIT WARN EVENT INFO DEBUG MID_DEBUG M_DBG FULL_DEBUG F_DBG] <br /> |
| `ganeshaVpaSpec` _string_ | GaneshaVpaSpec set ganesha horizontal pod autoscaler spec |  |  |
| `ganeshaNetpolOmit` _boolean_ | GaneshaNetpolOmit whether to omit default network policy for ganesha. Default: true |  |  |
| `ganeshaNetpolIngressIpblock` _string_ | GaneshaNetpolIngressIpblock defines ingress ip block for ganesha default network policy |  |  |
| `ganeshaNetpolEgressIpblock` _string_ | GaneshaNetpolEgressIpblock defines egress ip block for ganesha default network policy |  |  |


#### PostgresMode

_Underlying type:_ _string_

PostgresMode describes mode postgres runs

_Validation:_
- Enum: [standalone readreplicas]

_Appears in:_
- [PostgresSpec](#postgresspec)



#### PostgresSpec



PostgresSpec defines the desired state of Postgres



_Appears in:_
- [LMSMoodleSpec](#lmsmoodlespec)
- [LMSMoodleTemplateSpec](#lmsmoodletemplatespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `postgresMode` _[PostgresMode](#postgresmode)_ | PostgresMode describes mode postgres runs |  | Enum: [standalone readreplicas] <br /> |
| `postgresExtraConfig` _string_ | PostgresExtraConfig contains extra postgres config |  |  |
| `postgresSize` _integer_ | PostgresSize defines postgres number of replicas |  |  |
| `postgresImage` _string_ | PostgresImage defines image for postgres container |  | MaxLength: 255 <br /> |
| `postgresUpgrade` _boolean_ | PostgresUpgrade defines whether postgres upgrade is enabled |  |  |
| `postgresPvcDataSize` _string_ | PostgresPvcDataSize defines postgres storage size |  | MaxLength: 20 <br />MinLength: 2 <br /> |
| `postgresPvcDataStorageAccessMode` _[StorageAccessMode](#storageaccessmode)_ | PostgresPvcDataStorageAccessMode defines postgres storage access modes |  | Enum: [ReadWriteOnce ReadOnlyMany ReadWriteMany] <br /> |
| `postgresPvcDataStorageClassName` _string_ | PostgresPvcDataStorageClassName defines postgres storage class |  | MaxLength: 63 <br />MinLength: 2 <br /> |
| `postgresPvcDataAutoexpansion` _boolean_ | PostgresPvcDataAutoexpansion enables autoexpansion |  |  |
| `postgresPvcDataAutoexpansionIncrementGib` _integer_ | PostgresPvcDataAutoexpansionIncrementGib defines Gib to increment |  |  |
| `postgresPvcDataAutoexpansionCapGib` _integer_ | PostgresPvcDataAutoexpansionCapGib defines limit for autoexpansion increments |  |  |
| `postgresResourceRequests` _boolean_ | PostgresResourceRequests whether postgres resource requests are added. Default: true |  |  |
| `postgresResourceRequestsCpu` _string_ | PostgresResourceRequestsCpu set postgres resource requests cpu |  | MaxLength: 20 <br /> |
| `postgresResourceRequestsMemory` _string_ | PostgresResourceRequestsMemory set postgres resource requests memory |  | MaxLength: 20 <br /> |
| `postgresResourceLimits` _boolean_ | PostgresResourceLimits whether postgres resource limits are added. Default: false |  |  |
| `postgresResourceLimitsCpu` _string_ | PostgresResourceLimitsCpu set postgres resource limits cpu |  | MaxLength: 20 <br /> |
| `postgresResourceLimitsMemory` _string_ | PostgresResourceLimitsMemory set postgres resource limits memory |  | MaxLength: 20 <br /> |
| `postgresTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | PostgresTolerations defines any tolerations for Postgres pods. |  |  |
| `postgresNodeSelector` _string_ | PostgresNodeSelector defines any node labels selectors for Postgres pods. |  |  |
| `postgresAffinity` _string_ | PostgresAffinity defines any affinity rules for Postgres pods. |  |  |
| `postgresVpaSpec` _string_ | PostgresVpaSpec set postgres horizontal pod autoscaler spec |  |  |
| `postgresNetpolOmit` _boolean_ | PostgresNetpolOmit whether to omit default network policy for postgres. Default: true |  |  |
| `postgresNetpolIngressIpblock` _string_ | PostgresNetpolIngressIpblock defines ingress ip block for postgres default network policy |  |  |
| `postgresNetpolEgressIpblock` _string_ | PostgresNetpolEgressIpblock defines egress ip block for postgres default network policy |  |  |
| `postgresReadreplicasSize` _integer_ | PostgresReadreplicasSize defines postgres readreplicas number of replicas |  |  |
| `postgresReadreplicasPvcDataSize` _string_ | PostgresReadreplicasPvcDataSize defines postgres readreplicas storage size |  | MaxLength: 20 <br />MinLength: 2 <br /> |
| `postgresReadreplicasPvcDataStorageAccessMode` _[StorageAccessMode](#storageaccessmode)_ | PostgresReadreplicasPvcDataStorageAccessMode defines postgres readreplicas storage access modes |  | Enum: [ReadWriteOnce ReadOnlyMany ReadWriteMany] <br /> |
| `postgresReadreplicasPvcDataStorageClassName` _string_ | PostgresReadreplicasPvcDataStorageClassName defines postgres readreplicas storage class |  | MaxLength: 63 <br />MinLength: 2 <br /> |
| `postgresReadreplicasPvcDataAutoexpansion` _boolean_ | PostgresReadreplicasPvcDataAutoexpansion enables autoexpansion |  |  |
| `postgresReadreplicasPvcDataAutoexpansionIncrementGib` _integer_ | PostgresReadreplicasPvcDataAutoexpansionIncrementGib defines Gib to increment |  |  |
| `postgresReadreplicasPvcDataAutoexpansionCapGib` _integer_ | PostgresReadreplicasPvcDataAutoexpansionCapGib defines limit for autoexpansion increments |  |  |
| `postgresReadreplicasResourceRequests` _boolean_ | PostgresReadreplicasResourceRequests whether postgres readreplicas resource requests are added. Default: true |  |  |
| `postgresReadreplicasResourceRequestsCpu` _string_ | PostgresReadreplicasResourceRequestsCpu set postgres readreplicas resource requests cpu |  | MaxLength: 20 <br /> |
| `postgresReadreplicasResourceRequestsMemory` _string_ | PostgresReadreplicasResourceRequestsMemory set postgres readreplicas resource requests memory |  | MaxLength: 20 <br /> |
| `postgresReadreplicasResourceLimits` _boolean_ | PostgresReadreplicasResourceLimits whether postgres readreplicas resource limits are added. Default: false |  |  |
| `postgresReadreplicasResourceLimitsCpu` _string_ | PostgresReadreplicasResourceLimitsCpu set postgres readreplicas resource limits cpu |  | MaxLength: 20 <br /> |
| `postgresReadreplicasResourceLimitsMemory` _string_ | PostgresReadreplicasResourceLimitsMemory set postgres readreplicas resource limits memory |  | MaxLength: 20 <br /> |
| `postgresReadreplicasTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | PostgresReadreplicasTolerations defines any tolerations for PostgresReadreplicas pods. |  |  |
| `postgresReadreplicasNodeSelector` _string_ | PostgresReadreplicasNodeSelector defines any node labels selectors for PostgresReadreplicas pods. |  |  |
| `postgresReadreplicasAffinity` _string_ | PostgresReadreplicasAffinity defines any affinity rules for PostgresReadreplicas pods. |  |  |
| `postgresReadreplicasVpaSpec` _string_ | PostgresReadreplicasVpaSpec set postgres readreplicas horizontal pod autoscaler spec |  |  |
| `pgbouncerExtraConfig` _string_ | PgbouncerExtraConfig contains extra pgbouncer config |  |  |
| `pgbouncerResourceRequests` _boolean_ | PgbouncerResourceRequests whether pgbouncer resource requests are added. Default: true |  |  |
| `pgbouncerResourceRequestsCpu` _string_ | PgbouncerResourceRequestsCpu set pgbouncer resource requests cpu |  | MaxLength: 20 <br /> |
| `pgbouncerResourceRequestsMemory` _string_ | PgbouncerResourceRequestsMemory set pgbouncer resource requests memory |  | MaxLength: 20 <br /> |
| `pgbouncerResourceLimits` _boolean_ | PgbouncerResourceLimits whether pgbouncer resource limits are added. Default: false |  |  |
| `pgbouncerResourceLimitsCpu` _string_ | PgbouncerResourceLimitsCpu set pgbouncer resource limits cpu |  | MaxLength: 20 <br /> |
| `pgbouncerResourceLimitsMemory` _string_ | PgbouncerResourceLimitsMemory set pgbouncer resource limits memory |  | MaxLength: 20 <br /> |
| `pgbouncerTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | PgbouncerTolerations defines any tolerations for Pgbouncer pods. |  |  |
| `pgbouncerNodeSelector` _string_ | PgbouncerNodeSelector defines any node labels selectors for Pgbouncer pods. |  |  |
| `pgbouncerAffinity` _string_ | PgbouncerAffinity defines any affinity rules for Pgbouncer pods. |  |  |
| `pgbouncerVpaSpec` _string_ | PgbouncerVpaSpec set pgbouncer horizontal pod autoscaler spec |  |  |
| `pgbouncerNetpolOmit` _boolean_ | PgbouncerNetpolOmit whether to omit default network policy for pgbouncer. Default: true |  |  |
| `pgbouncerNetpolIngressIpblock` _string_ | PgbouncerNetpolIngressIpblock defines ipblock for pgbouncer default network policy |  |  |
| `pgbouncerNetpolEgressIpblock` _string_ | PgbouncerNetpolEgressIpblock defines egress ip block for pgbouncer default network policy |  |  |
| `pgbouncerReadonlyExtraConfig` _string_ | PgbouncerReadonlyExtraConfig contains extra pgbouncer readonly config |  |  |
| `pgbouncerReadonlyResourceRequests` _boolean_ | PgbouncerReadonlyResourceRequests whether pgbouncer readonly resource requests are added. Default: true |  |  |
| `pgbouncerReadonlyResourceRequestsCpu` _string_ | PgbouncerReadonlyResourceRequestsCpu set pgbouncer readonly resource requests cpu |  | MaxLength: 20 <br /> |
| `pgbouncerReadonlyResourceRequestsMemory` _string_ | PgbouncerReadonlyResourceRequestsMemory set pgbouncer readonly resource requests memory |  | MaxLength: 20 <br /> |
| `pgbouncerReadonlyResourceLimits` _boolean_ | PgbouncerReadonlyResourceLimits whether pgbouncer readonly resource limits are added. Default: false |  |  |
| `pgbouncerReadonlyResourceLimitsCpu` _string_ | PgbouncerReadonlyResourceLimitsCpu set pgbouncer readonly resource limits cpu |  | MaxLength: 20 <br /> |
| `pgbouncerReadonlyResourceLimitsMemory` _string_ | PgbouncerReadonlyResourceLimitsMemory set pgbouncer readonly resource limits memory |  | MaxLength: 20 <br /> |
| `pgbouncerReadonlyTolerations` _[Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26.6/#toleration-v1-core) array_ | PgbouncerReadonlyTolerations defines any tolerations for PgbouncerReadonly pods. |  |  |
| `pgbouncerReadonlyNodeSelector` _string_ | PgbouncerReadonlyNodeSelector defines any node labels selectors for PgbouncerReadonly pods. |  |  |
| `pgbouncerReadonlyAffinity` _string_ | PgbouncerReadonlyAffinity defines any affinity rules for PgbouncerReadonly pods. |  |  |
| `pgbouncerReadonlyVpaSpec` _string_ | PgbouncerReadonlyVpaSpec set pgbouncer readonly horizontal pod autoscaler spec |  |  |


#### RoutineStatusCrNotify



RoutineStatusCrNotify specification using ansible URI module



_Appears in:_
- [MoodleSpec](#moodlespec)

| Field | Description | Default | Validation |
| --- | --- | --- | --- |
| `url` _string_ | HTTP or HTTPS URL in the form (http\|https)://host.domain[:port]/path |  |  |
| `statusCode` _integer array_ | StatusCode A list of valid, numeric, HTTP status codes that signifies success of the request. |  |  |
| `method` _string_ | Method The HTTP method of the request or response. |  | Enum: [GET POST PUT PATCH DELETE] <br /> |
| `uuid` _[RoutineStatusCrNotifyUUID](#routinestatuscrnotifyuuid)_ | UUID used when notifying status to an endpoint |  | MaxLength: 36 <br />MinLength: 36 <br /> |
| `headers` _[RoutineStatusCrNotifyHeaders](#routinestatuscrnotifyheaders)_ | Headers used when notifying status to an endpoint |  |  |
| `jwtSecretEnvName` _string_ | JwtSecretEnvName environment variable name that holds secret to generate jwt tokens |  |  |


#### RoutineStatusCrNotifyHeaders



RoutineStatusCrNotifyHeaders used when notifying status to an endpoint



_Appears in:_
- [RoutineStatusCrNotify](#routinestatuscrnotify)



#### RoutineStatusCrNotifyUUID

_Underlying type:_ _string_

RoutineStatusCrNotifyUUID used when notifying status to an endpoint

_Validation:_
- MaxLength: 36
- MinLength: 36

_Appears in:_
- [RoutineStatusCrNotify](#routinestatuscrnotify)



#### SessionRedisCompressor

_Underlying type:_ _string_

SessionRedisCompressor describes Moodle redis session compresor

_Validation:_
- Enum: [none gzip zstd]

_Appears in:_
- [MoodleSpec](#moodlespec)



#### StorageAccessMode

_Underlying type:_ _string_

StorageAccessMode describes storage access modes

_Validation:_
- Enum: [ReadWriteOnce ReadOnlyMany ReadWriteMany]

_Appears in:_
- [KeydbSpec](#keydbspec)
- [MoodleSpec](#moodlespec)
- [NfsSpec](#nfsspec)
- [PostgresSpec](#postgresspec)



