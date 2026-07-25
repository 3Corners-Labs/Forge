# M1 exit operational journey evidence

Status: Proposed
Owner: Forge Product, Experience, Context, Delivery, Identity, Policy, Execution, Verification, Evidence, Approval, Source Control, Platform, Security, Privacy, Quality, and Engineering
Verified: 2026-07-25

## Iteration outcome

This candidate makes the accepted M1 journey operational through the public
`forge delivery start` command. The default remains the fast,
credential-free contract mode. An explicit `--runtime docker` mode now
performs the same guided journey with:

- exact, purpose-limited source bytes admitted by Context;
- a non-root, network-disabled, read-only isolated worker image;
- an independently reconstructed, read-only verifier checkout and separate
  non-root verifier image;
- capability-owned PostgreSQL state;
- digest-addressed MinIO proposal and candidate artefacts; and
- a PostgreSQL-backed contract-double publication provider that reconciles
  idempotently after control-plane restart.

The observed allowed journey reached exactly one reviewable contract-double
pull request. It retained separate exact proposal-promotion, plan-approval,
and publication-approval actions. No live GitHub credential, merge, release,
deployment, production effect, or provider activation was used.

## Why

The prior public journey proved the authority chain with in-process runners
and memory stores. It did not prove that a user could run Forge with the
accepted isolation and durability choices. This tranche closes that gap
without broadening model, repository, publication, or production authority.

The implementation also adds `forge doctor --operational`. It checks the
Docker daemon and the presence of the required PostgreSQL and MinIO
configuration without displaying configuration values.

## Applicable architecture

- The Engineering Constitution requires explicit human authority, untrusted
  intelligence, bounded execution, independent verification, evidence,
  privacy, and responsible resource use.
- ADR-0010 bounds M1 and forbids production effects.
- ADR-0013 requires risk-proportional assurance without weakening
  consequential publication blockers.
- ADR-0017 requires isolated deterministic materialization and independent
  verification.
- ADR-0018 selects PostgreSQL control state and MinIO artefacts for the
  Developer Trial.
- ADR-0019 keeps live GitHub publication separately configured and approved.
- ADR-0020 and ADR-0021 preserve the guided, resumable human authority gates.
- ADR-0022 and ADR-0023 keep model discovery and change preparation
  non-authoritative.

This candidate implements accepted choices. It introduces no new capability,
trust boundary, provider commitment, or production topology. No ADR is
required.

## Capability effects

| Capability | Operational responsibility |
|---|---|
| Context | Catalogues two declared repository sources at an exact revision and returns only purpose-limited bytes and digests |
| Delivery | Owns the plan, lifecycle, declared paths, candidate state, and durable run records |
| Identity | Persists root and worker delegations behind Identity-owned stores |
| Policy | Persists exact envelopes and keeps deterministic decisions authoritative |
| Execution | Reconstructs an ephemeral write workspace and launches the immutable worker with closed inputs and resource limits |
| Verification | Independently reconstructs source and launches a separate immutable verifier with a read-only checkout |
| Evidence | Persists append-only records and digest-addressed bounded artefacts |
| Approval | Persists separate plan and publication approvals |
| Source Control | Persists publication operations and reconciles the contract double without duplicate pull requests |
| Experience | Exposes immediate progress and exact human actions through the public CLI |
| Tool | Reports operational prerequisite readiness without resolving or printing secrets |

No capability reads another capability's database tables directly. Shared
contracts remain implementation-free, and the control plane does not import
worker or verifier implementation.

## Implemented operational controls

- Docker images must be immutable OCI digests or immutable local image IDs;
  mutable tags fail closed.
- Worker and verifier run as UID/GID `65532`, with no network, no
  capabilities, no privilege escalation, a read-only root filesystem,
  bounded CPU, memory, processes, duration, and output, and an isolated
  no-exec temporary filesystem.
- The Go runtime is limited to one logical processor while the container
  retains enough bounded PIDs to start safely.
- The worker sees only admitted write-path source. The verifier receives a
  separately reconstructed read-only source view and the candidate bundle.
- Bind mounts use purpose-built project-local ephemeral directories. The
  repository itself is never mounted into either runtime.
- Source paths, symlinks, digests, output size, changed-file count, candidate
  digest, tree digest, plan digest, and result integrity are checked before
  promotion.
- Temporary worker and verifier material is removed and cleanup is asserted.
- PostgreSQL migrations are capability-owned and fail closed.
- MinIO artefacts are size-bounded, digest-verified on write and read, and
  recorded through Evidence metadata.
- Contract-double branch and pull-request state survives provider
  reconstruction and repeated publication remains idempotent.

## Generated tests

Automated tests cover:

- repository-catalog path, revision, digest, classification, and
  purpose-limitation checks;
- independent worker and verifier source reconstruction, traversal and stale
  digest denial;
- immutable-image validation and exact Docker isolation arguments;
- a real immutable worker image execution;
- the complete public operational journey using real Docker, PostgreSQL, and
  MinIO;
- fixed-UID worker execution on Linux with explicit launcher-readable and
  removable artefact modes inside an otherwise private project-local runtime
  root, plus bounded, non-sensitive preparation failure codes;
- PostgreSQL Evidence continuity after service reconstruction;
- MinIO digest-addressed write and verified read;
- PostgreSQL contract-double reconciliation without duplicate publication;
- operational prerequisite diagnostics;
- denied cross-scope admission before repository discovery; and
- module-boundary and architecture fixtures.

The local candidate passed:

~~~text
go test ./...
go test -race ./...
go vet ./...
go build ./...
python3 -m unittest scripts/test_autonomy_grant.py
python3 scripts/check_architecture.py
go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...
docker run --rm --network none --mount type=bind,src=$PWD,dst=/docs,readonly \
  squidfunk/mkdocs-material:9.7.7 build --strict --site-dir /tmp/forge-site
~~~

The first vulnerability scan found two reachable dependency vulnerabilities.
The candidate upgraded `github.com/jackc/pgx/v5` to `v5.9.2` and
`golang.org/x/text` to `v0.39.0`; the repeated scan reported zero reachable
vulnerabilities. It still reported one known vulnerability in an imported
package and 21 in required modules on paths that Forge does not call. These
remain dependency-monitoring input, not a claim that the dependency graph is
vulnerability-free.

The exact committed candidate revision must still pass the protected-CI suite.

## Boundaries and next product step

This is an operational Developer Trial, not a production platform. It does
not implement general repository indexing, cross-process conversational
resume, federated users, live GitHub publication, merge, release, deployment,
production isolation, backup/restore assurance, or a supported capacity
envelope.

After this candidate is protected-CI green, the next product tranche is to
use this Forge path for one real Forge self-change and then begin the first
Workspace Orchestration vertical slice through the conversational Product
Manager flow. Live provider and live GitHub effects remain separately
approved actions.
