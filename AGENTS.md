# Repository guidance

Before working, read these optional instruction files in order,
resolving these paths from this repository's root:

1. `.agents/organization/AGENTS.md`
2. `.agents/workspace/AGENTS.md`

Read each file if it resolves to a readable regular file.
Skip absent files; report broken or unreadable links.
Read each resolved file only once to avoid duplicate loading and cycles.

Resolve references inside imported files relative to their real target
directory after following symlinks.

Apply organization guidance, then workspace guidance, then the repository
instructions below. More specific applicable instructions take precedence.

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
