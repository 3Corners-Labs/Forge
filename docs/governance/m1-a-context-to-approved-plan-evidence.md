# M1-A context-to-approved-plan evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-23

## Bounded outcome

This combined M1-A tranche implements one coherent control-plane path from an already captured Delivery run to one exact human-approved Delivery Plan:

~~~text
Captured
  -> ContextAuthorizing
  -> ContextReady
  -> Planned
  -> AwaitingPlanApproval
  -> PlanApproved
~~~

The states describe recorded progress only. `PlanApproved` is not worker authority. This tranche stops before child delegation, worker Execution Envelope issuance, worker execution, repository mutation, candidate verification, publication approval, source-control publication, merge, release, deployment, live providers, production credentials, or production operation.

## Why this combined tranche

ADR-0015 requires the authoritative plan to follow an exact reproducible Context Package so that omissions, conflicts, classifications, source revisions, risks, budgets, verification intent, limitations, and rollback are visible before human review. ADR-0016 requires that review to use a separate Approval-owned request and scope-bound v2 record rather than a Delivery Boolean or reused publication approval.

Combining Context reproduction, planning, plan review, and the `PlanApproved` transition proves the first complete human decision point without adding micro-approvals between deterministic reversible steps.

## Implemented capability scope

### Context

Context resolves a capability-owned manifest by tenant and package identifier, verifies its digest and lifecycle, rechecks the Tenant-derived run binding, and compares current authorized source metadata with the exact recorded source revisions, locators, classifications, content digests, and freshness facts.

The result reports exact reproduction, existing and newly observed omissions, changed source identifiers, stale source identifiers, and checked source revisions. It returns no source content and cannot add authority.

### Delivery

Delivery provides:

- explicit transitions through `ContextAuthorizing`, `ContextReady`, `Planned`, `AwaitingPlanApproval`, and `PlanApproved`;
- expected-version optimistic concurrency and scope-bound idempotency;
- atomic state and minimized outbox updates in the concurrency-safe in-process Developer Trial adapter;
- exact Context Package, manifest, Execution Envelope scope, Policy reference, classification, source-revision, omission, conflict, and expiry recording;
- a deterministic, ordered acyclic Delivery Plan whose profile is resolved from a trusted defensive registry rather than accepted from the caller;
- canonical SHA-256 identity over every scope, input, graph, path, capability, resource, assurance, risk, limitation, stop, rollback, recovery, expiry, and budget field;
- rejection of path traversal, undeclared capabilities or paths, cycles, duplicates, unbounded resources, and any model, tool, secret, destination, or network requirement in M1-A;
- an Approval subject projection containing only immutable plan-review facts; and
- transition to `PlanApproved` only after resolving the current Approval-owned v2 record and validating exact expected facts.

The in-process adapter proves state ownership, atomicity, replay, concurrency, and outbox semantics but is not durable across process loss. ADR-0018 reserves PostgreSQL durability for M1-B.

### Policy

Policy evaluates `forge.approval.decide` at the approval decision point for the exact Approval Request digest. It binds the active human, current Delivery plan, Tenant scope, purpose, classification, risk, Developer Trial environment, policy versions, obligations, and expiry.

The Policy Decision establishes eligibility to decide. It is not the human decision and does not create an Approval Record.

### Approval

Approval provides:

- a closed `forge.approval-request.v1` schema and strict decoder;
- a closed, separately versioned `forge.approval-record.v2` schema and strict decoder while leaving Approval Record v1 unchanged and readable;
- canonical collection ordering, deterministic request identity, and deterministic record integrity;
- subject resolution through a port rather than Delivery persistence access;
- Policy decision resolution through a port rather than Policy runtime access;
- exact intended-human, scope, environment, purpose, classification, risk, action, subject, digest, evidence, condition, policy, and expiry validation;
- fresh Delivery-owned subject resolution both when the request opens and immediately before the human decision;
- atomic tenant-scoped request and decision idempotency;
- one current decision branch per request;
- immutable revocation records referencing the prior approved record; and
- current-state resolution that rejects stale branches.

An Approval Request is not approval. An approved record is not execution authority. Neither object can widen a plan, delegation, envelope, context package, path, capability, budget, model, tool, secret, destination, or expiry.

### Experience

The Approval Request and Delivery plan expose the immutable identifiers, source revisions, omissions, conflicts, risks, budgets, verification intent, limitations, rollback, and evidence references needed by a later CLI or Workbench projection. No public UI is added in M1-A; the implementation baseline continues to state that the current CLI does not expose plan review.

## Capability boundaries

- `forge.capability.context` owns Context Package assembly and reproduction.
- `forge.capability.delivery` owns the plan, plan digest, workflow state, optimistic concurrency, idempotency, and outbox.
- `forge.capability.policy` owns approver-eligibility decisions.
- `forge.capability.approval` owns requests, human decision records, revocation chains, idempotency, and current-state resolution.
- `forge.capability.experience` will render authorized projections but owns no canonical state.

The shared trust-contract module owns serialization, strict decoding, structural validation, canonicalization, digest and integrity checks, and expected-fact helpers. It owns no runtime state and grants no authority.

Architecture fitness registers Approval as a Trial control-plane module with dependencies only on Identity and shared failure and trust contracts. Delivery may consume the shared trust contract but cannot import Approval or Policy runtime modules or construct Approval Records or Execution Envelopes.

## Governing architecture

The implementation follows the Accepted Constitution, trustworthy vertical slice, first-slice capability contracts, first-slice data flow and privacy analysis, ADR-0015, ADR-0016, and the current implementation baseline.

No accepted capability boundary, authority ordering, deployment unit, provider decision, or production posture changes in this tranche.

## Verification evidence

Deterministic tests cover:

- closed strict Approval Request v1 and Approval Record v2 decoding, unknown-field rejection, and explicit v1/v2 separation;
- canonical digest stability under permitted collection ordering and mutation detection for request, scope, subject, and record fields;
- cross-environment, changed-plan, unsupported-condition, evidence, expiry, decision, and expected-fact rejection;
- the complete atomic Delivery lifecycle through `PlanApproved`, state/version conflicts, idempotent replay, and one minimized outbox event per committed transition;
- exact Context Package, envelope scope, source revision, classification, omission, conflict, and expiry binding before planning;
- deterministic plan graph and digest validation, normative budget invalidation, path traversal denial, and prohibition of undeclared external resources;
- denial of caller-controlled or unregistered planner profiles and direct caller-supplied Approval Requests;
- exact and changed Context reproduction plus cross-run manifest rejection;
- Approval request, permit, rejection, evidence and condition mismatch, idempotency conflict, cross-tenant resolution, immutable revocation, stale-record rejection, and fork prevention;
- Policy permit binding to exact approver, request digest, run, scope, purpose, classification, risk, and environment, including cross-scope denial; and
- end-to-end composition of Delivery and Approval through shared immutable subject and capability-owned request/record resolver ports without runtime-module imports; and
- architecture fitness for the schema set, module dependencies, lifecycle order, and prohibitions on Delivery-owned envelopes or approvals.

Repository verification ran:

- `python3 scripts/check_architecture.py`;
- `go test ./...`;
- `go test -race ./...`;
- `go test ./... -cover`;
- `python3 -m unittest scripts/test_autonomy_grant.py`;
- `go vet ./...`;
- `go build ./...`;
- `go run golang.org/x/vuln/cmd/govulncheck@v1.1.4 ./...`; and
- focused `go test -run '^$' -bench ... -benchmem` commands.

The focused packages reported 72.9% Approval coverage, 76.7% Context coverage, 74.6% shared trust-contract coverage, 81.3% Delivery coverage, and 88.9% Policy coverage. Coverage percentages describe this candidate run only and are not an accepted assurance threshold.

## Performance observations

Focused in-process benchmarks on an Apple M1 Pro (`darwin/arm64`, eight benchmark workers) measured:

| Operation | Time | Allocated bytes | Allocations |
|---|---:|---:|---:|
| Approval Request canonicalization and integrity | 6.284 µs/op | 4,083 B/op | 24 allocs/op |
| Deterministic Delivery Plan construction | 10.559 µs/op | 6,872 B/op | 37 allocs/op |
| Current Approval Record resolution and integrity verification | 4.404 µs/op | 2,739 B/op | 19 allocs/op |
| Concurrent current Approval resolution with eight benchmark workers | 2.297 µs/op | 2,745 B/op | 19 allocs/op |
| Current immutable revocation resolution | 3.962 µs/op | 2,867 B/op | 19 allocs/op |

Two 2,000-sample focused observations measured Approval resolution at p50 3.208 µs, p95 5.792 µs, and p99 24.959 µs, and deterministic plan construction at p50 7.292 µs, p95 11.625 µs, and p99 41.416 µs. The concurrent benchmark corresponds to approximately 435,000 in-process resolutions per second under that isolated benchmark configuration.

These observations indicate that the implemented in-process controls are not a perceptible UX bottleneck in isolation. They do not establish sustained capacity, a maximum contention envelope, PostgreSQL latency, provider latency, noisy-neighbour isolation, crash recovery, or production capacity. ADR-0013 performance and scalability gates therefore remain active for later tranches.

## Explicit exclusions

This tranche does not implement PostgreSQL or MinIO persistence, crash recovery, scheduling, cancellation, child delegation, worker authority, workers, models, tools, secrets, repository mutation, candidate artefacts, independent candidate verification, publication approval, evidence persistence, GitHub integration, pull-request publication, merge, release, deployment, production profiles, public CLI wiring, or a support commitment.

## Authorized delivery evidence

Delivery used grant `forge.autonomy-grant.m1-bounded-delivery-v1`, protected revision `7c09c6084d761cbe0c4005a0ecf8954b69328625`, digest `sha256:1e9b2e8030d5cb7229a92909ac8099cae209c70ea7b663d7a659a4533acff0fa`, task `m1-a-context-to-approved-plan`, and branch `codex/m1-a-context-to-approved-plan`.

The human-approved amendment added only `mkdocs.yml` to register this evidence document in canonical navigation. Merge to protected `main` remains a mandatory human gate.

## Architecture decision assessment

No new ADR is proposed. The tranche implements ADR-0015 and ADR-0016 inside the accepted modular control plane.

A new ADR SHALL precede another approval authority model, general or model-generated planning, dynamic graph expansion, public experience semantics beyond the accepted CLI boundary, production signing or key custody, movement of Approval ownership, service extraction, or any attempt to treat plan approval as worker or publication authority.
