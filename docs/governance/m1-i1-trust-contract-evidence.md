# M1 iteration 1 trust-contract evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-21

## Decision context

ADR-0010 authorizes the bounded M1 implementation. This iteration establishes executable, provider-neutral representations of the trust objects that later M1 delivery behavior must consume without weakening their authority or capability boundaries.

## Implemented scope

The iteration implements:

- strict JSON decoding for Policy Decision, Execution Envelope, Context Package, Approval Record, and Evidence Record v1;
- rejection of unknown fields and trailing JSON data;
- structural and lifecycle validation aligned with the accepted closed schemas;
- exact policy-decision binding to decision point, subject, action, and resource;
- active-window and granted-action checks for execution envelopes;
- child-envelope checks that prohibit scope, action, resource-list, validity-window, or limit expansion;
- exact Context Package binding to one active Execution Envelope without granting authority through context;
- exact approval binding to action, subject type, and subject digest;
- exact evidence binding to subject type, identifier, and digest;
- a versioned, disclosure-safe normalized failure contract; and
- module-boundary declarations that keep the contract packages in the `shared-contract` deployment unit.

## Affected capabilities

This iteration provides shared contract evidence for Policy, Intent and Delivery, Context, Approval, Evidence, and Verification. It does not claim those capabilities are operational. Capability ownership, persistence, issuance, retrieval, orchestration, execution, publication, and production use remain outside this iteration.

## Governing architecture

The implementation is constrained by the engineering constitution, ADR-0001, ADR-0003, ADR-0007, ADR-0010, the trustworthy first vertical slice, the first-slice capability contracts, the first-slice data flow, the security architecture, and the accepted G6 and G8 evidence.

The shared packages define data contracts and fail-closed checks only. They do not become a second policy engine, context service, approval service, evidence store, or delivery runtime.

## Verification evidence

The iteration is verified by:

- unit tests for valid instances of every trust object;
- negative tests for unknown or authority-shaped Context Package fields and trailing data;
- policy denial, indeterminate, mismatch, expiry, and unknown-obligation tests;
- envelope action, expiry, and child-authority expansion tests;
- Context Package envelope and expiry-binding tests;
- approval action, digest, revocation, and expiry tests;
- evidence subject-mismatch tests;
- normalized failure contract tests;
- repository architecture fitness checks; and
- the repository verification task.

## Explicit exclusions

This iteration does not implement authentication, policy evaluation, envelope issuance or revocation, context retrieval, secret resolution, model or tool execution, durable workflow, cryptographic verification, evidence persistence, approval UX, repository mutation, pull-request publication, or release.

## Architecture decision assessment

No new ADR is proposed. The implementation realizes accepted contracts within the already accepted shared-contract boundary. A new ADR SHALL be proposed before changing capability ownership, authority semantics, schema compatibility, cryptographic guarantees, deployment isolation, or provider neutrality.
