# ADR-0016: Introduce a scope-bound Approval Record v2

Status: Accepted
Owner: Forge Approval, Security, and Architecture
Date: 2026-07-23

## Context

ADR-0015 requires human plan review to bind action `forge.delivery.execute-plan` to the exact Delivery Plan digest, tenant, workspace, project, repository, delivery run, Developer Trial environment, evidence, conditions, actor eligibility, and expiry. It separates that decision from later Policy-issued worker authority and from publication approval.

The Accepted Approval Record v1 contract cannot prove those invariants independently. It contains the approver, action, subject type and digest, decision, time window, conditions, optional delegation reference, and evidence correlation identifier. It does not contain:

- tenant, workspace, project, repository, or delivery-run scope;
- environment or classification;
- the Approval Request whose exact review subject was presented;
- the Policy Decisions that established approver eligibility;
- the evidence references reviewed by the human;
- an issuer or immutable decision-chain reference; or
- integrity binding over the complete record.

ADR-0003 already requires approval replay against another environment to fail, but v1 cannot distinguish environments. Encoding scope or environment only inside the subject digest would force every consumer to resolve an external subject before it can reject cross-scope replay and would leave the Approval Record itself unable to explain its declared boundary.

Adding required fields to `forge.approval-record.v1` would silently break a closed accepted schema. Treating Approval-owned storage as sufficient integrity would weaken the portable trust contract at capability boundaries. A versioned replacement is required before plan approval becomes executable.

## Decision

Introduce `forge.approval-record.v2` and a matching digest-bound `forge.approval-request.v1`. Approval Record v1 remains readable for compatibility and historical validation but cannot satisfy a newly implemented M1 plan-approval or publication gate.

### Approval Request

Approval owns the immutable Approval Request and its canonical digest. The request is a review subject, not human approval and not execution authority.

The request records:

- contract version, request identifier, request digest, and evidence correlation identifier;
- requesting actor and intended approver subject or eligible role references;
- exact tenant, workspace, project, repository, and delivery-run scope;
- environment, purpose, classification, and risk class;
- protected action, subject type, subject identifier, and subject digest;
- required evidence references and the versions or digests needed to reproduce the review projection;
- registered condition identifiers the approver may impose;
- requested expiry and the maximum decision lifetime;
- Policy Decision references required to decide the request; and
- creation time, issuer, and current request lifecycle.

For the first implementation:

- action is exactly `forge.delivery.execute-plan`;
- subject type is exactly `delivery-plan`;
- the subject identifier and digest resolve to the current Delivery-owned canonical plan;
- scope is the one Tenant-derived Developer Trial scope of that plan;
- environment is `developer-trial`;
- classification and risk cannot be lower than the current plan, repository, and Context Package facts;
- evidence includes the captured intent, Context Package identity and manifest digest, source revisions, omissions, conflicts, assumptions, constraints, risks, budgets, verification intent, limitations, and rollback references;
- the initiating authenticated human is the named intended approver permitted by ADR-0015;
- condition identifiers come from a trusted Developer Trial registry; and
- no model, tool, repository content, prompt, or caller-provided projection can change authoritative request facts.

Approval resolves the subject from a Delivery-owned port before creating the request. It does not import or read Delivery persistence. The resolved projection carries only the immutable identifiers, digests, scope, lifecycle, classification, risk, review references, and expiry required for the Approval contract.

The canonical request digest binds every authority-, subject-, scope-, environment-, evidence-, condition-, lifecycle-, and expiry-relevant field. Presentation layout and disclosure-safe explanatory text are not authority. If a normative fact changes, a new request and new eligibility decision are required.

### Approver eligibility

Policy evaluates action `forge.approval.decide` for the exact Approval Request digest as resource. Trusted input includes the active human principal, Tenant-derived scope, environment, purpose, classification, risk, protected action, subject type, and request expiry.

Approval accepts only a current Policy-owned permit whose:

- decision point is the plan-review approval gate;
- subject is the exact active human approver;
- action is `forge.approval.decide`;
- resource is the canonical Approval Request digest;
- scope, environment, classification, purpose, and risk facts match the request;
- policy versions and expiry remain valid; and
- mandatory obligations and conditions are recognized and reported as enforced.

Policy does not create the Approval Request, decide for the human, or mint the Approval Record. The human decision is not a Policy Decision. A permit to decide does not imply an approved outcome.

### Approval Record v2

Approval Record v2 records:

- schema version, immutable record identifier, issuer, and evidence correlation identifier;
- Approval Request identifier and canonical request digest;
- approver identifier and optional delegation reference;
- exact tenant, workspace, project, repository, and delivery-run scope;
- environment, purpose, classification, and risk class;
- protected action, subject type, subject identifier, and subject digest;
- decision: `approved`, `rejected`, or `revoked`;
- decision time and expiry;
- imposed registered condition identifiers;
- Policy Decision identifiers used for approver eligibility;
- evidence references reviewed for the decision;
- optional previous Approval Record identifier for an immutable decision chain; and
- integrity algorithm, digest, and optional signature reference.

The integrity digest uses deterministic SHA-256 canonicalization and binds every serialized field except the digest itself. The record identifier is derived from or separately bound to the canonical record content. A signature remains optional for the in-process Developer Trial; provider selection, signing keys, rotation, and production verification require the accepted cryptography and deployment decisions.

Collections are canonical, duplicate-free, and contain no empty values. Unknown condition or obligation identifiers fail closed. Conditions constrain use; they cannot widen the request, plan, Policy Decision, delegation, Execution Envelope, destination, model, tool, secret, path, budget, or expiry.

The record expiry is no later than the earliest of:

- Approval Request expiry;
- human principal expiry;
- approver-eligibility Policy Decision expiry;
- Delivery Plan and Context Package validity;
- applicable delegation or policy lifetime where required; and
- the Developer Trial maximum approval lifetime.

An approval record never outlives a changed or invalidated subject merely because its timestamp remains current. Consumers resolve the current request, subject digest, Approval decision chain, and required eligibility decisions at the protected gate.

### Immutable lifecycle and revocation

Approval Records are immutable. Rejection and approval create one record for the request. Revocation creates a new v2 record whose previous-record reference identifies the current approved record and whose request, scope, environment, action, subject, and digest remain exact.

Approval owns current-state resolution and rejects:

- two active decision branches for one request;
- revocation of an unknown, expired, rejected, already revoked, cross-scope, or mismatched record;
- reuse of an idempotency key for different request or decision input; and
- a stale record when a later valid decision-chain record exists.

Creation and revocation are atomic and tenant-scoped. Resolver results are defensive, disclosure-safe projections from the Approval-owned store. Caller-provided Approval Record values are not canonical state.

The initial adapter is concurrency-safe and in-process. It proves ownership, atomicity, idempotency, chain resolution, and mutation detection but does not claim production durability.

### Gate validation

A receiving capability validates Approval Record v2 against explicit expected facts:

- current Approval-owned resolution of the record and decision chain;
- schema version, issuer, integrity, request identifier, and request digest;
- approver and eligibility Policy Decisions;
- tenant, workspace, project, repository, delivery run, environment, purpose, classification, and risk;
- action, subject type, subject identifier, and subject digest;
- approved decision, active lifecycle, conditions, evidence references, and expiry.

For the ADR-0015 implementation, Delivery may transition only from `AwaitingPlanApproval` to `PlanApproved` after those validations pass for `forge.delivery.execute-plan`. Delivery records the Approval Record identifier and request digest as workflow facts. It does not copy the record, reinterpret policy, mint approval, or treat `PlanApproved` as worker authority.

The same v2 contract may later represent publication approval with another exact action and candidate subject, but publication implementation remains unauthorized in this iteration. An approval for one action or subject can never satisfy another gate.

### Version compatibility

- `forge.approval-record.v1` and its schema remain immutable and readable.
- V1 may validate its own historical structure but SHALL be rejected by new M1 plan, execution, and publication gates that require exact scope and environment.
- V2 uses a separate schema identifier, decoder, validator, canonicalizer, and authorizer.
- No implicit up-conversion from v1 to v2 is allowed. Migration requires reconstruction from authoritative request, subject, scope, eligibility, and evidence sources followed by a new human decision.
- Existing trust-object names do not imply that v1 and v2 have equal assurance.

### Capability and module boundaries

- `forge.capability.approval` owns Approval Requests, Approval Records, idempotency, decision chains, revocation, and current-state resolution.
- `forge.capability.delivery` owns the plan and run state. It supplies an immutable subject projection through a port and consumes only validated Approval-owned resolution through a port and shared trust contract.
- `forge.capability.policy` owns eligibility Policy Decisions. Approval validates them through a Policy-owned resolver and does not import Policy runtime implementation.
- `forge.capability.identity` owns principal lifecycle and delegation.
- `forge.capability.tenant` owns scope and repository binding.
- `forge.capability.experience` renders the exact request and sends authenticated human input; it owns no Approval state.
- the shared trust-contract module owns versioned serialization, canonicalization, structural validation, integrity validation, and expected-fact authorization helpers without owning Approval state or policy.

The initial Approval module remains part of the modular control plane. No service extraction is authorized.

### Bounded implementation authorization

After this decision is Accepted on protected `main`, a separately activated M1 iteration MAY implement:

- Approval Request v1 and Approval Record v2 schemas and shared trust validation;
- a bounded Approval module with in-process store, subject and Policy resolver ports, plan-review request, decide, revoke, and resolve operations;
- Context reproduction and Delivery transitions from the existing captured run through `ContextReady`, `Planned`, `AwaitingPlanApproval`, and `PlanApproved`;
- deterministic plan construction under ADR-0015;
- Policy rules required only for exact plan-review approver eligibility;
- architecture fitness, focused benchmarks, evidence, documentation, and implementation-baseline updates.

It SHALL NOT implement worker or child delegation, worker Execution Envelopes, execution, tools, models, secrets, repository mutation, candidate verification, publication approval, source-control publication, live providers, production storage, production credentials, or a production maturity claim.

## Alternatives considered

- **Add required fields to Approval Record v1:** rejected because it silently changes a closed accepted schema and breaks exact version semantics.
- **Bind scope and environment only inside the plan digest:** rejected because the Approval Record cannot independently reject or explain cross-scope and cross-environment replay.
- **Store scope only on the Approval Request:** rejected because a detached record could be replayed without resolving and matching its request at every gate.
- **Let Delivery store a review Boolean:** rejected because it crosses Approval ownership and loses actor, eligibility, conditions, expiry, evidence, integrity, and revocation.
- **Let Policy produce the Approval Record:** rejected because Policy decides eligibility and protected-action authorization, not human judgment.
- **Trust caller-supplied Policy Decisions or Approval Records:** rejected because public values are untrusted and cannot replace capability-owned resolution.
- **Mutate an approved record to revoked:** rejected because mutation destroys the reviewed historical decision and weakens evidence lineage.
- **Require signatures in the first in-process Trial:** rejected because integrity and owner-controlled resolution can prove the bounded contract without prematurely selecting production key infrastructure.
- **Implement only through `Planned` and defer approval:** viable as a smaller increment, but rejected as the preferred sequence because v2 resolves a known accepted-contract gap and allows one coherent context-to-`PlanApproved` tranche.

## Consequences

- Approval becomes independently scope-, environment-, subject-, evidence-, eligibility-, condition-, and lifecycle-bound.
- One generic v2 contract can support distinct future approval actions without sharing approvals between them.
- V1 compatibility is explicit but weaker records cannot silently pass stronger gates.
- Approval introduces a new capability-owned state module and resolver boundary while remaining in the modular control plane.
- Immutable revocation creates additional records and resolution work but preserves decision history and detects forks.
- Every protected consumer performs more validation. ADR-0013 requires focused latency and allocation evidence rather than removing required checks for speed.
- Production signing, durable topology, external identity providers, retention schedules, and publication remain open or later decisions.

## Validation

Architecture and implementation evidence SHALL prove:

- closed v1 remains unchanged and v2 has an independent schema identifier and decoder;
- strict request and record decoding rejects unknown fields, invalid versions, duplicates, empty identifiers, invalid time windows, unsupported conditions, and malformed integrity;
- canonical request and record digests are stable under permitted collection ordering and change for every relevant field mutation;
- Policy-owned eligibility resolution binds exact approver, request digest, decision point, scope, environment, purpose, classification, risk, policy versions, obligations, and expiry;
- Approval-owned subject resolution rejects stale, changed, cross-scope, cross-repository, cross-run, cross-environment, or invalid plans before review;
- atomic tenant-scoped creation, idempotent replay, conflict on changed input, immutable revocation, fork detection, and disclosure-safe current-state resolution;
- rejection of mismatched action, subject type, subject identifier, subject digest, actor, request, evidence, condition, policy reference, scope, environment, purpose, classification, risk, decision, expiry, or integrity;
- v1 cannot satisfy plan, execution, or publication gates requiring v2;
- plan approval cannot issue or widen delegation, an Execution Envelope, worker authority, verification, evidence acceptance, publication, merge, release, deployment, or policy mutation;
- Delivery cannot construct Approval Records and Approval cannot mutate Delivery state;
- exact negative tests for cross-tenant, cross-project, cross-repository, cross-run, cross-environment, changed-plan, stale-context, expired, revoked, indeterminate, unsupported-condition, evidence-mismatch, and approval-replay attempts; and
- focused p50, p95, p99, throughput, allocation, contention, and revocation-resolution measurements under the applicable M1 and ADR-0013 validation envelope.

## Accountable review

Sijbren Beukenkamp accepted this decision as drafted through the authenticated collaboration boundary on 2026-07-23, acting as the temporary accountable Architecture, Security and Privacy, Approval, Identity, Policy, Delivery, Experience, Quality, and Engineering reviewer established for bounded M1.

Acceptance authorizes only this contract architecture. Functional implementation remains blocked until the Accepted decision reaches protected `main` and a separate bounded M1 implementation iteration is activated.

## Supersession

If accepted, this decision specializes ADR-0003 and ADR-0015 by adding a stronger versioned Approval contract. It does not alter Execution Envelope issuance under ADR-0012, risk-proportional assurance under ADR-0013, or root context delegation under ADR-0014.
