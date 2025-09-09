# LMS Moodle Operator

LMS Moodle Operator is a [meta operator](https://github.com/cncf/tag-app-delivery/blob/main/operator-whitepaper/v1/Operator-WhitePaper_v1-0.md#operator-of-operators) to automate the deployment and management of Moodle instances in Kubernetes, refered as `LMSMoodle` resources. It handles the full stack required to run them: Postgres, Keydb, NFS-Ganesha, and Moodle. Those components have their own Operators. In addition, a `LMSMoodleTemplate` custom resource is provided. It is like a template that can be reused when creating a `LMSMoodle` resource.

**Key Technologies:**

* Kubernetes
* Operator SDK
* Moodle
* Postgres
* Keydb
* NFS Ganesha

## Prerequisites

* [Moodle Operator](https://github.com/krestomatio/moodle-operator): For automated provisioning and configuration of the web layer: php-fpm, nginx, cronjob, among other resources.
* [Postgres Operator](https://github.com/krestomatio/postgres-operator) (Optional): For deployment and management of PostgreSQL database.
* [NFS Operator](https://github.com/krestomatio/nfs-operator) (Optional): For deployment and management of shared storage using NFS Ganesha.
* [Keydb Operator](https://github.com/krestomatio/keydb-operator) (Optional): For deployment and management of compatible Redis database: Keydb.

## Installation

> **Important Note:** This LMS Moodle Operator is currently in **Beta** stage. Proceed with caution in production deployments.

To install this this operator along **all** its required and optional prerequisites, follow these steps:

1. **Install Operator:**
```bash
kubectl apply -k https://github.com/krestomatio/lms-moodle-operator/config/operators?ref=v0.6.3
```

2. **Configure a LMSMoodleTemplate:**
- Download and modify [this `LMSMoodleTemplate` sample](https://raw.githubusercontent.com/krestomatio/lms-moodle-operator/v0.6.3/config/samples/lms_v1alpha1_lmsmoodletemplate.yaml) file to define a lms moodle template. A `LMSMoodleTemplate` can be use by one or many `LMSMoodle` resources as spec template.
```bash
curl -sSL 'https://raw.githubusercontent.com/krestomatio/lms-moodle-operator/v0.6.3/config/samples/lms_v1alpha1_lmsmoodletemplate.yaml' -o lms_v1alpha1_lmsmoodletemplate.yaml
# modify lms_v1alpha1_lmsmoodletemplate.yaml
```

3. **Deploy the LMSMoodleTemplate:**
- Deploy Moodle `LMSMoodleTemplate` using the modified configuration:
```bash
kubectl apply -f lms_v1alpha1_lmsmoodletemplate.yaml
```

4. **Configure a LMSMoodle:**
> **Note:** `LMSMoodle` resource specification has precedence over `LMSMoodleTemplate` specification.
- Download and modify [this sample](https://raw.githubusercontent.com/krestomatio/lms-moodle-operator/v0.6.3/config/samples/lms_v1alpha1_lmsmoodle.yaml) file to reflect your specific `LMSMoodle` stack configuration options. This file defines the desired state for your instance and all its layers handle by the operators. Note that it references a `LMSMoodleTemplate` resource by its name. The `LMSMoodleTemplate` resoure in the previous step.
```bash
curl -sSL 'https://raw.githubusercontent.com/krestomatio/lms-moodle-operator/v0.6.3/config/samples/lms_v1alpha1_lmsmoodle.yaml' -o lms_v1alpha1_lmsmoodle.yaml
# modify lms_v1alpha1_lmsmoodle.yaml
```

5. **Deploy the LMSMoodle:**
- Deploy a Moodle `LMSMoodle` using the modified configuration:
```bash
kubectl apply -f lms_v1alpha1_lmsmoodle.yaml
```

6. **Monitor Logs:**
- Track the LMS Moodle Operator logs for insights into the deployment process:
```bash
kubectl -n lms-moodle-operator-system logs -l control-plane=controller-manager -c manager -f
```

- Monitor the status of your deployed `LMSMoodle` instance:
```bash
kubectl get -f lms_v1alpha1_lmsmoodle.yaml -w
```

## Uninstall

1. **Delete LMSMoodle:**
```bash
# Caution: This step leads to data loss. Proceed with caution.
kubectl delete -f lms_v1alpha1_lmsmoodletemplate.yaml
```

2. **Delete LMSMoodleTemplate:**
```bash
# Caution: This step leads to data loss. Proceed with caution.
kubectl delete -f lms_v1alpha1_lmsmoodle.yaml
```

3. **Uninstall the Operator:**
```bash
kubectl delete -k https://github.com/krestomatio/lms-moodle-operator/config/operators?ref=v0.6.3
```

## Configuration

LMSMoodleTemplate and `LMSMoodle` custom resources (CRs) can be configure via their spec field: check [API Reference](api.md) for the respective documentation.

## Contributing

* Report bugs, request enhancements, or propose new features using GitHub issues.
