# M1-B approved-plan-to-verified-candidate evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-23

## Bounded outcome

This combined M1-B tranche implements the internal path from one exact approved Delivery Plan to one independently verified candidate:

~~~text
PlanApproved
  -> ExecutionAuthorizing
  -> Executing
  -> Verifying
  -> AwaitingPublicationApproval
~~~

`AwaitingPublicationApproval` records that the candidate passed its bound Assurance Plan. It is not publication approval and grants no source-control, merge, release, deployment, or production authority. M1-C remains responsible for the separate candidate-digest approval and GitHub publication adapter.

## Why this combined tranche

ADR-0017 requires bounded child authority, a closed worker request, an isolated worker, a separate verifier deployment unit, and independent consequential controls. ADR-0018 requires restart-safe capability-owned state and append-oriented evidence before autonomy can continue safely across process loss.

Implementing those contracts, runtimes, transitions, adapters, and tests together proves one coherent boundary. Splitting them would temporarily make candidate production look usable without independent verification or durable evidence.

## Governing architecture

The Accepted Constitution, trustworthy vertical slice, first-slice capability contracts, first-slice data flow and privacy analysis, runtime and workflow architecture, ADR-0010, ADR-0012, ADR-0013, ADR-0017, ADR-0018, and the implementation baseline govern this tranche.

No accepted authority order, capability owner, provider boundary, publication gate, or production posture changes.

## Implemented capability scope

### Identity and Policy

Identity owns a separate worker-phase delegation request, delegation lifecycle, workload principal binding, expiry, idempotency, and revocation. It does not widen the existing root task delegation.

Policy evaluates the exact approved plan, current Approval Record v2, current Context manifest, Tenant-derived scope, human, task agent, worker principal, action, risk, classification, environment, and requested bounds. Policy alone issues the worker Execution Envelope. The envelope contains no model, tool, secret, destination, or network authority and permits only `forge.worker.prepare-change`.

PostgreSQL adapters persist worker delegations and worker envelopes in capability-owned schemas.

### Delivery

Delivery owns the M1-B lifecycle states, optimistic versions, transition idempotency, candidate references, verification references, and transactional outbox. It accepts a candidate only when the worker request, envelope, delegation, approved plan, current approval, Context manifest, source revision, scope, changed paths, bundle digest, tree digest, result, and artefact references match exactly.

It enters `AwaitingPublicationApproval` only after a current Assurance Plan and independent passing Verification Result match the exact plan and candidate digests. The PostgreSQL adapter commits each state transition and its minimized outbox event atomically.

### Execution and worker

Execution constructs a closed digest-bound worker request from trusted plan and envelope facts. The request can select only the Forge-owned deterministic file-replacement transformation and cannot name a command, executable, model, tool, secret, destination, or network endpoint.

The worker is a separate `forge-worker` deployment unit. Its runtime rejects traversal, symlink escapes, source-digest changes, undeclared paths, duplicate paths, expired authority, and output-limit breaches. The Docker adapter pins an image digest and requests no network, a read-only root filesystem, all capabilities dropped, no privilege escalation, a non-root user, bounded processes, memory and CPU, bounded mounts, and cleanup proof.

### Verification

Verification owns immutable Assurance Plans and Verification Results. Every passing plan requires at least one consequential control not owned by the candidate. Candidate-owned evidence can be recorded but cannot be the sole publication oracle.

The separate `forge-verifier` deployment unit uses a read-only fresh checkout and imports shared contracts, never the worker or control-plane implementation. The current Forge-owned controls verify candidate integrity, exact changed-path scope, and Go formatting. PostgreSQL persists Assurance Plans and results with digest-conflict denial.

### Evidence

Evidence owns append-only records, tenant/project/run chain heads, serializable compare-and-swap, idempotency, canonical SHA-256 integrity, previous-record linkage, stream verification, and mutation or fork detection.

MinIO stores bounded candidate or verification bytes under server-derived tenant/project/run and digest keys. Uploads use a temporary object, byte re-read, SHA-256 verification, immutable-key promotion, and PostgreSQL metadata. Object-prefix traversal is rejected. Evidence is distinct from telemetry.

### Restart-safe resolution

Developer Trial PostgreSQL adapters now exist for Delivery, Identity worker delegation, Policy worker envelopes, Context manifests, Approval requests and records, Tenant scope facts, Verification plans and results, Evidence records and chain heads, and artefact metadata. Each adapter owns its schema and is consumed through capability ports; no module reads another capability's tables.

These are replaceable Developer Trial adapters. They do not prove production durability, backup, restore, RTO, RPO, high availability, multi-host safety, or supported capacity.

## Contracts and boundaries

Closed schemas and strict Go decoders cover:

- Worker Request v1;
- Candidate Bundle v1;
- Worker Result v1;
- Assurance Plan v1; and
- Verification Result v1.

Architecture fitness registers control plane, shared-contract, worker, and verifier deployment units. Negative fixtures prove that the control plane cannot import worker or verifier runtimes, the worker cannot import the control plane, and the verifier cannot import either the control plane or worker.

Workers and verifiers receive no PostgreSQL or MinIO credentials through their closed contracts. The control plane invokes them through ports, allowing the current Docker adapters to be replaced without changing capability ownership.

## Verification evidence

Deterministic tests cover:

- strict unknown-field rejection and digest mutation for worker and verification contracts;
- path traversal, undeclared write scope, source-digest changes, symlinks, duplicate files, and output bounds;
- the complete Delivery lifecycle from `PlanApproved` to `AwaitingPublicationApproval`, including exact authority, candidate, assurance, evidence, version, idempotency, and outbox binding;
- independent-oracle requirements and rejection of candidate-owned evidence as the sole consequential pass;
- fresh-checkout source matching, candidate integrity, declared path scope, Go formatting, unsupported independent controls, and empty-workspace denial;
- required Docker isolation arguments, digest-pinned images, bounded resources, non-root execution, read-only verifier checkout, strict outputs, and cleanup proof;
- append-only evidence integrity, record mutation, idempotency mutation, previous-digest linkage, and fork prevention;
- PostgreSQL evidence-chain continuation after service reconstruction;
- MinIO upload, digest-addressed promotion, durable metadata, and object-prefix traversal denial; and
- positive module boundaries plus negative control-plane, worker, verifier, and undeclared-dependency fixtures.

The dedicated M1 durable-verification workflow provisions synthetic PostgreSQL and MinIO services and executes the persistence and artefact integration tests. The general CI workflow continues to run formatting, vet, race tests, coverage, and builds.

Candidate verification ran:

- `python3 scripts/check_architecture.py`;
- `python3 -m unittest scripts/test_autonomy_grant.py`;
- `go test -race ./...`;
- PostgreSQL and MinIO integration tests for Delivery, Evidence, and Verification;
- `go vet ./...`;
- `go build ./...`;
- separate scratch-image builds for `forge-worker` and `forge-verifier`; and
- `go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...`.

The vulnerability scan found no called vulnerabilities. It reported vulnerable symbols in indirect packages or required modules that the candidate does not call; this remains dependency-monitoring input rather than a claim that those modules contain no known vulnerability.

## Performance and scalability posture

The worker and verifier are bounded background operations, so they do not sit on an interactive UX response path. Admission, authority, state transition, and status projection remain control operations governed by ADR-0013. Expensive transformation and independent verification may take longer but remain bounded by plan, envelope, process, memory, output, attempt, and expiry limits.

This tranche adds indexes for run identifiers, state/version, idempotency, outbox readiness, evidence chain heads, record identifiers, candidate digests, and artefact digests. It does not yet establish a sustained throughput, repository-size, concurrent-run, database-contention, MinIO-capacity, or recovery envelope. Those measurements remain required before a broader Trial claim.

## Explicit exclusions

This tranche does not implement a public CLI path, model-generated changes, arbitrary commands or tools, networked workers, secret delivery, dynamic worker profiles, general repository mutation, cancellation orchestration, automatic outbox dispatch, orphan-object reconciliation, retention execution, production cryptography, signatures, backup or restore, publication approval, GitHub publication, pull-request creation, merge, release, deployment, production profiles, or a support commitment.

The Docker tests verify the constructed isolation contract. Host-level adversarial escape assurance and production sandbox selection remain later security work.

## Authorized delivery evidence

Delivery used grant `forge.autonomy-grant.m1-bounded-delivery-v1`, protected revision `448386982e360a03faab393d9c33328f96e504e7`, digest `sha256:1e9b2e8030d5cb7229a92909ac8099cae209c70ea7b663d7a659a4533acff0fa`, task `m1-b-approved-plan-to-verified-candidate`, and branch `codex/m1-b-approved-plan-to-verified-candidate`.

Human-approved narrow amendments added Approval, Context, and Tenant only for restart-safe resolution of existing facts; added `architecture/testdata/forbidden-undeclared-module.json` only to preserve its fail-closed expectation after Evidence registration; and added `mkdocs.yml` only to register this document in canonical navigation.

Merge to protected `main` remains a mandatory human gate.

## Architecture decision assessment

No new ADR is proposed. This tranche implements ADR-0017 and ADR-0018 within their accepted boundaries.

A new ADR SHALL precede arbitrary worker commands, model or tool authority, network or secret access, another sandbox trust model, candidate-controlled consequential verification, shared capability tables, production persistence or object storage, production signing or key custody, or any attempt to treat a passing Verification Result as publication approval.
