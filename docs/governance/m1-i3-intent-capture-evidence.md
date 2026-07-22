# M1 iteration 3 intent-capture evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-22

## Bounded outcome

This iteration implements the next authorized M1 processing step after identity and scope admission: one directly authenticated human can capture one bounded change objective, its constraints, and its acceptance criteria as an immutable delivery run in `Captured` state.

The outcome stops before policy evaluation, Execution Envelope issuance, plan recording, worker scheduling, context retrieval, approval, evidence persistence, or source-control activity.

## Implemented scope

The Intent and Delivery implementation provides:

- versioned capture request and result contracts with required correlation and idempotency identifiers;
- explicit direct-human admission with a required empty delegation chain for this pre-agent step;
- receiving-capability validation of the active normalized principal and Tenant-minted scope;
- exact repository and classification binding through an immutable Tenant-owned record;
- capability-owned objective, constraints, acceptance criteria, purpose, actor, scope, classification, state, and optimistic-concurrency version;
- immutable delivery-run projections with defensive collection copies;
- request fingerprints that exclude correlation metadata but bind every authority- or intent-relevant field;
- atomic scope-bound replay or conflict behaviour for idempotency keys;
- a concurrency-safe in-memory Developer Trial store behind a replaceable persistence port; and
- normalized disclosure-safe invalid, forbidden, expired, conflict, indeterminate, and dependency-unavailable outcomes.

Tenant receiving contracts now validate an active principal against a minted scope and validate a repository binding against that exact scope. Repository bindings with an unknown classification fail closed.

## Affected capabilities

- `forge.capability.intent`: owns the captured objective, constraints, and acceptance criteria.
- `forge.capability.delivery`: owns the initial immutable delivery-run identity, `Captured` state, version, and idempotent creation command.
- `forge.capability.tenant`: exposes only validation of its own immutable scope and repository-binding facts to a receiving capability.

Identity still owns principal verification. Tenant still owns membership and repository binding. Delivery does not evaluate policy, issue authority, authenticate callers, discover repositories, or execute work.

## Governing architecture

The implementation follows the engineering constitution; ADR-0001, ADR-0003, ADR-0007, ADR-0009, and ADR-0010; the accepted capability model; Identity and Tenancy architecture; Runtime and Durable Workflow architecture; first-slice capability contracts and data flow; Agentic Delivery; and the Developer Trial security and privacy decisions.

The accepted sequence places intent capture after authentication and scope resolution and before deterministic policy evaluation. The public `forge delivery start` command remains unimplemented because this iteration does not yet provide the complete Experience-to-policy boundary required by ADR-0009.

## Verification evidence

Deterministic tests cover:

- successful capture from a direct synthetic human with authoritative scope and repository binding;
- immutable run fields and defensive collection copies;
- exact replay for the same scope, actor, idempotency key, and intent despite a new correlation identifier;
- safe conflict for changed intent under an existing idempotency key;
- concurrent capture committing exactly one run identity;
- unsupported contract version and missing correlation, idempotency, delegation, objective, constraint, or acceptance input;
- rejection of workload or delegated actors in this bounded direct-human step;
- expired principal, principal/scope mismatch, unminted scope, and unminted binding before persistence;
- rejection of bindings minted by a different Tenant authority even when their visible scope identifiers match;
- rejection of a repository binding minted for another repository selector in the same project;
- unknown repository classification denial;
- store conflict, unavailability, malformed response, and diagnostic non-disclosure; and
- no persistence call after any precondition failure.

Repository verification additionally runs architecture fitness, Go formatting, vetting, all tests, the race detector, coverage, and the build.

## Explicit exclusions

This iteration does not implement a durable database adapter, delivery transitions after `Captured`, optimistic-concurrency mutation, transactional outbox, intent editing, policy evaluation, policy decision issuance, Execution Envelopes, agent delegation, planning, cancellation, recovery, context assembly, workers, credentials, verification, approvals, evidence storage, source-control publication, journey execution, or the public M1 CLI.

The in-memory store demonstrates atomic command semantics only. Process loss discards its records, so it is not durability or recovery evidence.

## Architecture decision assessment

No new ADR is proposed. The implementation realizes the already accepted `CaptureIntent` responsibility inside the modular control plane and preserves the existing capability ownership and trust sequence. An ADR SHALL precede selection of a durable storage topology, a changed public CLI contract, a new authority boundary, or a separate delivery deployment unit.
