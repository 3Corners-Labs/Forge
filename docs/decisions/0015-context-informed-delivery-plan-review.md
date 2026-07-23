# ADR-0015: Establish context-informed delivery planning and plan review

Status: Accepted
Owner: Forge Delivery, Approval, and Architecture
Date: 2026-07-23

## Context

The bounded M1 path now implements authenticated intent capture, exact scope resolution, root task delegation, Policy-owned authority for context assembly, and a bounded Context Package manifest. The next outcome is a reviewed delivery plan that can later become the subject of separately issued worker authority.

The Accepted architecture currently contains an ordering contradiction:

- ADR-0010 and the trustworthy vertical slice place the bounded Context Package before the reviewed Delivery Plan; and
- the runtime workflow places `Planned` and generic `Authorized` before `ContextReady`.

Planning before context is available prevents the reviewer from seeing source revisions, omissions, conflicts, freshness, classification, and other facts that can materially change the proposed work. Treating the existing `Planned` state as an undocumented preliminary plan would create a second meaning for one state and would not bind the later execution proposal to the reviewed facts.

The architecture also distinguishes plan review from publication approval but does not yet define their exact subjects. Reusing one approval for both actions would permit approval replay across different risks and lifecycle stages. Letting Delivery record its own consequential approval would cross the Approval capability boundary. Letting an Approval Record grant execution authority would cross the Policy boundary.

## Decision

Adopt a context-informed, digest-bound M1 planning sequence. Context assembly precedes the authoritative Delivery Plan. Plan review and publication approval are separate Approval lifecycles bound to different actions and subject digests.

### Revised first-slice lifecycle

The first-slice delivery lifecycle becomes:

~~~text
Captured
  -> ContextAuthorizing
  -> ContextReady
  -> Planned
  -> AwaitingPlanApproval
  -> PlanApproved
  -> ExecutionAuthorizing
  -> Executing
  -> Verifying
  -> AwaitingPublicationApproval
  -> Publishing
  -> Completed

Any active state -> Cancelling -> Cancelled
Any active state -> Failed -> Recovering -> prior safe state or FailedFinal
~~~

The state names describe durable delivery progress; they do not themselves grant authority:

- `ContextAuthorizing` means Delivery is coordinating exact Identity and Policy prerequisites for context assembly.
- `ContextReady` means Delivery has validated and recorded an immutable Context Package reference and manifest digest for the run.
- `Planned` means Delivery has recorded one canonical plan bound to the current intent and Context Package.
- `AwaitingPlanApproval` means Approval has opened a review for the exact plan digest.
- `PlanApproved` means Delivery has validated and recorded a current permitted plan-review Approval Record.
- `ExecutionAuthorizing` means later implementation may request worker-specific delegation and Policy authority. It does not inherit context authority or plan approval as execution authority.
- `AwaitingPublicationApproval` is the separate human publication gate for the independently verified candidate digest.

`Authorized` and `AwaitingApproval` are removed as ambiguous generic state names. An implementation SHALL name the protected action and approval subject explicitly.

Every transition uses the current delivery-run version, a scope-bound idempotency key, and optimistic concurrency. A committed transition and its minimized outbox entry are atomic. Capability calls may occur outside the transaction, but their results become workflow facts only after Delivery revalidates their immutable identifiers, digests, scope, lifecycle, and expected predecessor state.

### Context acquisition before planning

The Context Package is assembled under the existing context-only root delegation and Execution Envelope. Delivery records only:

- Context Package identifier and manifest digest;
- Execution Envelope and Policy Decision references;
- source revision set, omissions, conflicts, classification, and expiry needed for review;
- correlation and evidence references; and
- the delivery-run version against which the context was accepted.

Delivery does not copy unrestricted source content or reinterpret context authority. A missing, expired, mutated, cross-scope, stale, or conflicting package fails closed according to policy and the declared plan profile. Context mutation cannot create a plan transition.

If intent, repository binding, relevant source revision, Context Package manifest, classification, applicable policy, or mandatory obligation changes, the recorded plan and any plan approval become invalid. Recovery returns to the earliest safe state that can reproduce current context; it never silently reuses the prior approval.

### Canonical Delivery Plan

Intent and Delivery owns the versioned Delivery Plan and its canonical digest. For bounded M1 the plan records:

- contract version, plan identifier, delivery-run identifier, and plan version;
- immutable tenant, workspace, project, repository, and source-revision references;
- captured-intent version and Context Package identifier plus manifest digest;
- objective, constraints, and acceptance criteria;
- an ordered acyclic graph of bounded steps with stable identifiers and dependencies;
- affected capabilities and declared change paths;
- expected outputs and candidate artefact kinds;
- assumptions, known omissions, conflicts, risks, limitations, and stop conditions;
- declared model, tool, secret, destination, network, and filesystem requirements;
- time, concurrency, attempt, compute, token, and cost budgets;
- verification intent and mandatory guardrail references;
- rollback and manual-recovery instructions; and
- evidence correlation and creation time.

Every authority-, scope-, input-, execution-, assurance-, cost-, and review-relevant field contributes to the canonical digest. Presentation-only projections and reviewer comments do not alter the plan digest unless they change a normative plan field.

The first M1 planner is deterministic and profile-bound. It may produce only the narrow Forge change-preparation graph permitted by ADR-0010. It cannot add models, tools, secrets, destinations, network access, paths, capabilities, budgets, or actions beyond trusted profile constraints. General planning, model-generated plans, dynamic graph expansion, and multi-agent scheduling remain excluded.

Delivery validates that the plan references the exact current intent and Context Package before committing `Planned`. Plan construction is not authorization, verification, execution, evidence acceptance, or approval.

### Human plan review

Approval owns the plan-review lifecycle:

- action: `forge.delivery.execute-plan`;
- subject type: `delivery-plan`;
- subject digest: the canonical Delivery Plan digest;
- environment: `developer-trial`;
- scope: the exact tenant, workspace, project, repository, and delivery run;
- evidence: the referenced intent, Context Package, assumptions, omissions, conflicts, risks, budgets, verification intent, limitations, and rollback;
- expiry: no later than the earliest applicable principal, delegation, context, policy, plan-profile, or Developer Trial bound.

Experience presents that exact review subject through the public CLI boundary. Immediate acknowledgement and a visible waiting state follow ADR-0013; review duration is not forced into an interaction-latency target.

The authenticated human who initiated the bounded M1 run MAY review its plan because the plan is a proposal for execution, not an independently produced candidate verification or publication approval. The approval record still requires exact actor attribution, active authority, supported conditions, expiry, and policy-defined eligibility. This M1 exception does not authorize self-approval of candidate verification, publication, release, merge, deployment, or policy mutation.

An approval decision does not grant worker execution authority. After a permitted plan approval is recorded, Policy must separately evaluate worker execution from trusted current facts and issue an exact worker Execution Envelope under ADR-0012. Approval cannot widen the plan, delegation, envelope, Context Package, or policy outcome.

Denied, expired, revoked, indeterminate, mismatched, or conditionally unsatisfied plan approval cannot transition the run to `PlanApproved`. Reuse of a plan approval for a different plan digest, action, run, repository, environment, or publication request fails closed.

### Capability boundaries

- `forge.capability.delivery` owns the Delivery Plan, canonical digest, durable run state, optimistic concurrency, idempotency, and outbox transition.
- `forge.capability.context` owns Context Package assembly and reproduction; it neither constructs the plan nor approves it.
- `forge.capability.identity` owns principal and delegation lifecycle.
- `forge.capability.policy` evaluates protected actions and later issues Execution Envelopes; it neither owns the plan nor records human approval.
- `forge.capability.approval` owns plan-review requests and Approval Records.
- `forge.capability.experience` presents the plan and records authenticated human input through public contracts; it owns no canonical delivery or approval state.
- `forge.capability.verification` later owns the Assurance Plan and candidate verification. The Delivery Plan may reference verification intent but cannot select itself as the sole oracle.

The initial implementation remains within the modular control plane. This decision provides no evidence for service extraction.

### Iteration boundary

The next bounded implementation MAY stop at `PlanApproved` and prove the context-to-plan-review path without starting a worker. It may add:

- Context Package reproduction sufficient to validate the manifest and declared source revisions;
- Delivery-owned context recording, deterministic plan construction, digesting, and transitions through `PlanApproved`;
- Approval-owned plan review with an in-process Developer Trial adapter;
- public-contract projections and deterministic tests needed for review; and
- architecture, performance, evidence, and implementation-baseline updates.

That implementation SHALL NOT add child delegation, worker-specific authority, worker execution, repository mutation, verification, publication approval, source-control publication, live providers, production storage, or production credentials.

Child or subsequent task delegation, worker-envelope scope, worker isolation enforcement, credential resolution, and execution recovery require their own accepted design and separately activated implementation scope.

## Alternatives considered

- **Keep `Planned -> Authorized -> ContextReady`:** rejected because the reviewed plan would not be bound to the context facts that determine scope, omissions, risk, and feasibility.
- **Treat `Planned` as a preliminary context request:** rejected because it overloads the canonical Delivery Plan and leaves the final execution proposal without an explicit reviewed state.
- **Assemble context and plan concurrently:** rejected for the bounded M1 path because plan completeness and review depend on the exact context manifest and its omissions.
- **Let Delivery record a Boolean plan approval:** rejected because Approval owns consequential human approval lifecycle, actor attribution, conditions, expiry, and replay prevention.
- **Reuse plan approval for publication:** rejected because plan execution and publication act on different subjects, risks, evidence, and lifecycle times.
- **Let the Approval Record authorize execution directly:** rejected because approval records human judgment while Policy owns authorization and Execution Envelope issuance.
- **Implement planning, child delegation, worker authority, and isolation together:** rejected because ADR-0014 explicitly limits the current root delegation to context preparation and the worker boundary has separate isolation and recovery stop conditions.

## Consequences

- Reviewers see the exact context revisions, omissions, conflicts, risk, budgets, verification intent, limitations, and rollback before approving execution.
- The lifecycle gains additional explicit states and one separate human gate, increasing state-machine and evidence work while removing ambiguous authorization semantics.
- Plan, approval, execution authority, candidate verification, and publication approval have distinct owners and digests.
- Context or plan changes invalidate dependent approvals deterministically.
- The initial deterministic planner remains narrow and testable without introducing a model or general orchestration engine.
- Worker implementation cannot begin until later authority and isolation design is accepted.
- Accepted runtime, capability-contract, vertical-slice, data-flow, roadmap explanation, and journey documents require a consistency update when this ADR is accepted.

## Validation

The architecture reconciliation SHALL demonstrate:

- one canonical lifecycle order across ADR-0010, runtime workflow, vertical slice, first-slice contracts, data flow, roadmap, and journey assurance;
- distinct state and action names for context authorization, plan approval, execution authorization, and publication approval;
- no state name or Approval Record being represented as Policy authority;
- capability ownership consistent with the accepted module and trust boundaries; and
- architecture checks detecting reintroduction of the contradictory ordering.

The bounded implementation SHALL provide deterministic positive and negative tests proving:

- Context Package identity, manifest digest, scope, source revisions, lifecycle, omissions, conflicts, classification, and policy references are validated before planning;
- stable plan canonicalization and digest changes for every relevant input, scope, step, dependency, requirement, budget, risk, assurance, limitation, and rollback mutation;
- rejection of cyclic graphs, duplicate or unknown steps, broadened paths or capabilities, undeclared models, tools, secrets, destinations, network or filesystem access, and unbounded resources;
- optimistic-concurrency conflict, atomic transition/outbox behaviour, idempotent replay, cancellation, and recovery to the earliest safe state;
- exact plan-review action, subject digest, scope, actor, environment, conditions, expiry, and evidence binding;
- denial of cross-scope, stale-context, changed-plan, expired, revoked, indeterminate, unsupported-condition, and publication-replay attempts;
- plan approval cannot issue or widen delegation, Execution Envelope, Context Package, worker authority, verification, evidence, or publication;
- immediate CLI acknowledgement and explicit waiting state for human review; and
- focused benchmarks and resource bounds under the applicable M1 and ADR-0013 validation envelope.

## Accountable review

Sijbren Beukenkamp accepted this decision as drafted through the authenticated collaboration boundary on 2026-07-23, acting as the temporary accountable Product, Architecture, Security and Privacy, Delivery, Approval, Experience, Quality, and Engineering reviewer established for bounded M1.

Acceptance authorizes only the architecture reconciliation. Functional implementation remains blocked until this Accepted decision reaches protected `main` and a separate bounded M1 implementation iteration is activated.

## Supersession

If accepted, this decision specializes ADR-0010, ADR-0012, ADR-0013, and ADR-0014 for context-informed planning. It revises the ordering and ambiguous state names in the Accepted runtime workflow while preserving the authority/context separation and the prohibition on worker execution without separately issued authority.
