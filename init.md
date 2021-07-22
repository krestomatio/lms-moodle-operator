# How this project was bootstrap
## References
[Operator Framework SDK Tutorial](https://v1-9-x.sdk.operatorframework.io/docs/building-operators/golang/tutorial/)

[OKD Operators Tutorial](https://docs.okd.io/latest/operators/operator_sdk/golang/osdk-golang-tutorial.html)

[Using kubebuilder (same base that Operator SDK unde the hood) to build a CronJob](https://book.kubebuilder.io/cronjob-tutorial/cronjob-tutorial.html)
## Requeriments
* git
* go version 1.16
* docker version 17.03+.
* kubectl and access to a Kubernetes cluster of a compatible version.

## Init operator
```bash
# New go operator project
mkdir kio-operator
cd kio-operator
operator-sdk init \
    --domain=app.krestomat.io \
    --project-name "kio-operator" \
    --repo github.com/krestomatio/kio-operator

export GO111MODULE=on
operator-sdk edit --multigroup=true
# git commit

# Add Site api in M4e group
operator-sdk create api \
  --group=m4e \
  --version=v1alpha1 \
  --kind=Site \
  --resource \
  --controller
```
