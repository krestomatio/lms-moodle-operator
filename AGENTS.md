# Repository guidance

Follow this repository's documentation and any more specific instructions for
the files being changed.

## Repository scope

- This Go meta-operator turns an `LMSMoodle` plus `LMSMoodleTemplate` into a
  namespace, network policies, and Moodle/Postgres/KeyDB/NFS custom resources.
- API contracts live in `api/lms/v1alpha1/`; reconciliation, dependency ordering,
  finalization, suspension, and conditions live in `internal/controller/lms/`.
  `kio-web-app-api` is the primary producer of `LMSMoodle` resources.
- Preserve CRD group/kinds, merge precedence, names/labels, finalizers, status
  conditions, and child-operator API compatibility. Coordinate dependency version
  changes with every child operator and deployment pin.

## Validation

- Initialize `hack/mk`, then run `make test` and `make lint`. Add focused controller
  tests for reconciliation and lifecycle changes.
- When API markers/types change, regenerate manifests/code/docs and run
  `make bundle`; review CRD, RBAC, samples, dependency metadata, and generated diffs.
- E2E, deployment, image push, and release targets require an isolated cluster or
  explicit external authorization.
