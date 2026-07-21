# ADR-0010: Authorize the bounded M1 implementation

Status: Proposed
Owner: Forge Architecture Board
Date: 2026-07-21

## Context

Forge deliberately separated architecture acceptance, pre-implementation readiness, implementation authority, capability maturity, and production support. ADR-0006 accepted the minimum architecture boundary; ADR-0007 sequenced the minimum self-building kernel first; ADR-0009 established the public M1 CLI entry point. None of those decisions alone authorized functional implementation.

G0 through G10 are Accepted for F0 and are consolidated in protected `main` commit `93efcfcae8e946310ebd07294d744f38a3083b0e`, tree `eeea9e5e0d95534cebc8ac253672d00b818027c8`. The [F0 final acceptance manifest](../governance/f0-acceptance-manifest.md) binds the evidence, owners, verification, residual risks, and deferred obligations.

The remaining decision is whether those controls are sufficient to begin the bounded M1 implementation without implying production readiness or permitting the implementation to redefine its own authority.

## Decision

Once this ADR and the F0 final acceptance manifest are both Accepted at the same protected-main revision, authorize implementation of **M1 — Minimum self-building kernel** within the boundary below.

### Authorized outcome

Forge MAY implement one complete path that transforms one authenticated human request into one reviewable pull request against the Forge repository:

~~~text
authenticated human
  -> resolved Developer Trial workspace and repository scope
  -> captured intent and acceptance criteria
  -> deterministic Policy Decision and Execution Envelope
  -> bounded Context Package
  -> reviewed delivery plan
  -> separately isolated worker prepares one candidate
  -> independent verification binds the candidate digest
  -> human publication approval binds action, repository, environment, digest, conditions, and expiry
  -> idempotent branch and pull-request publication
  -> integrity-checked Evidence Records and role-appropriate explanation
~~~

The public entry point SHALL be `forge delivery start` as defined by ADR-0009. The existing `up` and `down` commands SHALL preserve their Developer-only validation.

### Authorized capability work

Implementation MAY add only the M1 responsibilities of:

- Tenant and Workspace;
- Identity and Delegation;
- Intent and Delivery;
- Policy Decision;
- Context Assembly;
- Tool Brokerage only where needed for the narrow worker contract;
- Isolated Execution;
- Verification;
- Approval;
- Audit Evidence;
- one replaceable source-control adapter;
- Experience through the public CLI and accepted journey interface.

These are capability responsibilities, not authorization to create separate services. The control plane remains modular-monolith-first under ADR-0001. Untrusted or resource-intensive repository work runs outside the control-plane process behind the accepted worker contract and isolation rules.

### Authorized operating envelope

- Developer Trial profile only.
- One synthetic test tenant, workspace, project, and Forge repository binding per run.
- Synthetic or explicitly approved non-production identities and data.
- Contract doubles before live provider integration where practical.
- Deterministic transformation by default.
- At most one model adapter only after separate provider eligibility, classification, budget, retention, and accountable approval; general model routing remains excluded.
- One source-control provider adapter after credential, destination, idempotency, ambiguity, and evidence tests pass.
- No publication without separate human approval of the exact candidate digest.
- Required architecture, dependency, isolation, journey, security, and CI fitness functions remain non-bypassable.

### Explicitly unauthorized

M1 SHALL NOT:

- merge a pull request, release software, deploy to production, or mutate production policy;
- use production credentials, production data, customer data, or production targets;
- claim Supported capability maturity, production readiness, SLA, RTO, RPO, capacity, residency, or compatibility;
- enable active or autonomous offensive testing;
- create broad multi-agent orchestration, self-expanding delegation, or autonomous approval;
- implement the Software System Twin, general advisory mode, broad model/tool routing, marketplace, billing, certification, or accreditation;
- implement Workplace-specific features;
- change an accepted capability, trust, deployment, authority, evidence, or public-interface boundary without normal review and an Accepted superseding ADR;
- treat generated code, generated tests, retrieved content, a model, a tool, or the candidate itself as the sole consequential oracle.

## Delivery controls

Every M1 iteration SHALL:

1. identify the bounded outcome and affected capabilities;
2. load the applicable accepted architecture and gate evidence;
3. preserve canonical ownership and cross-plane contracts;
4. implement the smallest end-to-end increment through public contracts;
5. add deterministic positive and negative tests proportional to risk;
6. update the implementation baseline and user-facing explanation;
7. record any material architecture change in an ADR before implementation;
8. use protected pull requests with required checks and accountable human approval.

No increment may claim a later milestone, Trial exit, Supported maturity, or production posture.

## Stop conditions

Implementation SHALL stop and return to accountable architecture review when:

- an accepted F0 assumption or gate result is contradicted;
- worker isolation cannot enforce the declared filesystem, process, network, secret, time, compute, cost, or cleanup boundary;
- a provider or framework choice materially binds a capability contract or deployment profile;
- a new external data destination, trust boundary, data class, credential type, production action, or autonomous capability is proposed;
- required CI, branch protection, CODEOWNERS, scenario ownership, kill switch, or evidence controls are missing or bypassable;
- scope expands beyond one bounded intent-to-pull-request journey;
- safe rollback or disclosure-safe evidence cannot be maintained.

## M1 exit evidence

M1 is complete only when the roadmap and vertical-slice exit criteria pass, including:

- accepted allowed and denied journeys through the same real public CLI;
- cross-scope, expired-authority, changed-digest, unknown-obligation, approval-replay, context-mutation, secret-leakage, cancellation, revocation, evidence, and publication negative tests;
- implemented control-plane/worker dependency and isolation conformance;
- interrupted-run recovery without duplicate publication;
- ambiguous source-control reconciliation without duplicate branch or pull request;
- independent verification and human approval bound to the exact candidate digest;
- evidence mutation detection and publication blocking when required evidence is unavailable;
- a real reviewable Forge improvement prepared on an isolated branch;
- an updated implementation baseline with exact source, checks, exclusions, and residual risks;
- accountable human M1-exit review.

Passing M1 does not authorize merge, release, production deployment, or promotion to Supported.

## Rollback and recovery

- Revert incomplete increments through the protected pull-request path.
- Preserve the existing Developer bootstrap and its safety validation until replacements pass equivalent evidence.
- Cancel running work, revoke resolvable task credentials, clean ephemeral worker state, and block publication when authority or evidence becomes invalid.
- Preserve durable delivery and minimized evidence needed for reconciliation; do not preserve unrestricted source, prompts, secrets, or personal data as rollback material.
- Use manual recovery when automatic reconciliation cannot prove safety or idempotency.

## Alternatives

- Begin M1 from accepted architecture alone: rejected because design acceptance is not implementation authority.
- Authorize the full capability catalogue: rejected because it violates the bounded vertical slice and increases circular trust risk.
- Implement advisory/System Twin capability before self-building: rejected by ADR-0007.
- Permit production deployment after the first successful journey: rejected because production conformance belongs to later milestones and profiles.
- Let implementation choose its own safety or journey oracle: rejected by ADR-0005 and G8–G10.

## Consequences

- Engineering may begin functional M1 work only after this Proposed decision becomes Accepted with the F0 manifest.
- M1 can deliver vertical increments across several capabilities, but cannot bypass their ownership or collapse their canonical records.
- The accepted F0 contracts become implementation constraints and M1-exit test obligations.
- Capability maturity remains Research unless a separate evidence-based lifecycle decision promotes a bounded Trial.
- The temporary concentration of accountable roles remains visible and must be revisited before M1 exit and any production-shaped claim.

## Accountable review

| Role | Reviewer | Outcome | Date |
|---|---|---|---|
| Product | Sijbren Beukenkamp, temporary | Pending | — |
| Architecture | Sijbren Beukenkamp, temporary | Pending | — |
| Security and Privacy | Sijbren Beukenkamp, temporary | Pending | — |
| Engineering | Sijbren Beukenkamp, temporary | Pending | — |
| Operations | Sijbren Beukenkamp, temporary | Pending | — |

Acceptance SHALL name the exact reviewed Git revision and explicitly approve both this decision and the F0 final acceptance manifest. The AI contributor prepared the draft but cannot provide accountable acceptance.

## Validation

- ADR and manifest statuses SHALL match and remain Proposed until accountable acceptance.
- Architecture checks SHALL verify the immutable F0 baseline, complete gate set, required exclusions, and absence of a premature implementation-ready claim.
- The authorization revision SHALL pass `task verify` and protected `validate` and `verify` checks.
- M1 implementation SHALL demonstrate every exit criterion without redefining the accepted oracle.

## Supersession

This decision activates the implementation sequence established by ADR-0007 for M1 only. It does not supersede ADR-0001 through ADR-0009.
