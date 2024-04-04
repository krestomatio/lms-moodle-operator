# Kio Operator

Kio operator is a meta operator to automate deployment and management of Moodle instances in Kubernetes, refered as `Sites`. It handles the stack required to run them: postgres, keydb, nfs, and moodle. Those components have their own operators. In addition, a Flavor is like a template to reuse when creating Sites resources.

**Key Technologies:**

* Kubernetes
* Ansible Operator SDK
* Moodle
* Postgres
* Keydb
* NFS Ganesha

## Prerequisites

* [Moodle Operator](https://github.com/krestomatio/moodle-operator): For automated provisioning and configuration of the web layer: php-fpm, nginx, cronjob, among other resources.
* [Postgres Operator](https://github.com/krestomatio/moodle-operator) (Optional): For deployment and management of PostgreSQL database.
* [NFS Operator](https://github.com/krestomatio/moodle-operator) (Optional): For deployment and management of shared storage using NFS Ganesha.
* [Keydb Operator](https://github.com/krestomatio/moodle-operator) (Optional): For deployment and management of compatible Redis database: Keydb.

## Installation

> **Important Note:** This Kio Operator is currently in **Beta** stage. Proceed with caution in production deployments.

To install this this operator along **all** its required and optional prerequisites, follow these steps:

1. **Install Operator:**
```bash
# All operators
kubectl apply -k https://github.com/krestomatio/kio-operator/config/operators?ref=v0.3.45
# Only Kio Operator
# kubectl apply -k https://github.com/krestomatio/kio-operator/config/default?ref=v0.3.44
```

2. **Configure a Flavor:**
- Download and modify [this Flavor sample](https://raw.githubusercontent.com/krestomatio/kio-operator/v0.3.44/config/samples/m4e_v1alpha1_flavor.yaml) file to define a Site flavor or template. A Flavor can be use by one or many Sites resources as spec template.
```bash
curl -sSL 'https://raw.githubusercontent.com/krestomatio/kio-operator/v0.3.44/config/samples/m4e_v1alpha1_flavor.yaml' -o m4e_v1alpha1_flavor.yaml
# modify m4e_v1alpha1_flavor.yaml
```

3. **Deploy the Flavor:**
- Deploy Moodle Flavor using the modified configuration:
```bash
kubectl apply -f m4e_v1alpha1_flavor.yaml
```

4. **Configure a Site:**
> **Note:** Site resource specification has precedence over Flavor specification.
- Download and modify [this sample](https://raw.githubusercontent.com/krestomatio/kio-operator/v0.3.44/config/samples/m4e_v1alpha1_site.yaml) file to reflect your specific Site stack configuration options. This file defines the desired state for your instance and all its layers handle by the operators. Note that it references a Flavor resource by its name. The Flavor resoure in the previous step.
```bash
curl -sSL 'https://raw.githubusercontent.com/krestomatio/kio-operator/v0.3.44/config/samples/m4e_v1alpha1_site.yaml' -o m4e_v1alpha1_site.yaml
# modify m4e_v1alpha1_site.yaml
```

5. **Deploy the Site:**
- Deploy a Moodle Site using the modified configuration:
```bash
kubectl apply -f m4e_v1alpha1_site.yaml
```

6. **Monitor Logs:**
- Track the Kio Operator logs for insights into the deployment process:
```bash
kubectl -n kio-operator-system logs -l control-plane=controller-manager -c manager -f
```

- Monitor the status of your deployed Site instance:
```bash
kubectl get -f m4e_v1alpha1_site.yaml -w
```

## Uninstall

1. **Delete Site:**
```bash
# Caution: This step leads to data loss. Proceed with caution.
kubectl delete -f m4e_v1alpha1_flavor.yaml
```

2. **Delete Flavor:**
```bash
# Caution: This step leads to data loss. Proceed with caution.
kubectl delete -f m4e_v1alpha1_site.yaml
```

3. **Uninstall the Operator:**
```bash
# All operators
kubectl delete -k https://github.com/krestomatio/kio-operator/config/operators?ref=v0.3.44
# Only Kio Operator
# kubectl delete -k https://github.com/krestomatio/kio-operator/config/default?ref=v0.3.44
```

## Configuration

Flavor and Site custom resources (CRs) can be configure via their spec field: check [API Reference](api.md) for the respective documentation.

## Contributing

* Report bugs, request enhancements, or propose new features using GitHub issues.
