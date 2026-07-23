# Runtime and durable workflow architecture

Status: Accepted
Owner: Forge Delivery and Execution Architecture
Version: 0.2

## Purpose

Separate durable delivery coordination from isolated, potentially untrusted execution. The control plane owns business state; workers execute bounded tasks and cannot become the source of delivery truth.

## Delivery-run state

The first slice uses action- and subject-specific states:

~~~text
Captured -> ContextAuthorizing -> ContextReady -> Planned
         -> AwaitingPlanApproval -> PlanApproved
         -> ExecutionAuthorizing -> Executing -> Verifying
         -> AwaitingPublicationApproval -> Publishing -> Completed

Any active state -> Cancelling -> Cancelled
Any active state -> Failed -> Recovering -> prior safe state or FailedFinal
~~~

Context is assembled before the canonical Delivery Plan so that source revisions, omissions, conflicts, classification, and lifecycle facts are part of the reviewed proposal. `PlanApproved` records a valid human decision on the exact plan digest; it is not worker authority. `ExecutionAuthorizing` requires separately evaluated Policy authority before execution. `AwaitingPublicationApproval` is a later, separate human gate bound to the independently verified candidate digest.

The generic `Authorized` and `AwaitingApproval` names are prohibited because they obscure the protected action and approval subject.

Transitions are commands evaluated against current state, policy where applicable, idempotency key and expected version. Every committed transition produces an outbox entry for evidence and downstream events. Cross-capability results become delivery facts only after their identifiers, digests, scope, lifecycle, and predecessor state have been revalidated.

## Control-plane responsibilities

- own delivery-run state and optimistic concurrency version;
- record the Context Package reference and manifest digest before constructing the canonical Delivery Plan;
- own the Delivery Plan and its canonical digest without treating the plan as authority;
- schedule work only after policy and prerequisite checks;
- issue or reference the Execution Envelope;
- enforce deadlines, cancellation and approval waits;
- reconcile ambiguous provider outcomes;
- checkpoint retryable progress;
- preserve manual recovery instructions.

The control plane does not execute repository code, retain plaintext secrets, record its own Approval Records or trust worker-reported success without verification.

## Worker responsibilities

- validate envelope integrity, expiry, scope and parent lineage;
- create an ephemeral isolated workspace;
- resolve only authorized short-lived credential references;
- enforce filesystem, network, resource and time limits;
- execute declared steps and return structured results;
- clean credentials and workspace according to retention policy;
- stop safely on cancellation, revocation or budget exhaustion.

Workers do not mutate durable workflow state directly. They return results correlated with the task and candidate digest.

## Idempotency and delivery semantics

Commands and external side effects have caller-supplied or server-derived idempotency keys scoped to tenant, action and target. Internal events assume at-least-once delivery; consumers record processed event IDs or implement naturally idempotent state transitions.

Publication adapters reconcile remote state before retrying an ambiguous request. A duplicate pull-request request must return the existing remote reference or an explicit conflict, not create another PR silently.

## Cancellation and recovery

Cancellation propagates from delivery run to queued tasks, workers, tool calls and provider requests where supported. Non-cancellable external actions are declared before execution and require stronger approval or compensating action.

Recovery distinguishes:

- retry same task with the same candidate inputs;
- retry from a prior checkpoint;
- compensate a completed external action;
- accept an externally completed result after reconciliation;
- require manual intervention.

Retry budgets include attempts, elapsed time and cost. Infinite retry loops are prohibited.

## Scaling and partitioning

Durable delivery state is partitioned by immutable tenant and delivery-run ID. Queues are partitioned by workload and risk class with per-tenant concurrency and cost limits. Scheduling preserves fairness and supports admission rejection before expensive provider work.

Service extraction requires evidence that a module needs distinct scale, isolation, availability, recovery, ownership or technology.

## Initial conformance scenarios

- Two concurrent commands cannot commit conflicting state versions.
- A changed Context Package or Delivery Plan invalidates its dependent plan approval and later authority.
- Plan approval cannot be replayed as worker authority or publication approval.
- Worker loss does not lose committed delivery state or create duplicate publication.
- Cancellation revokes task credentials and blocks subsequent side effects.
- A mismatched tenant, task or candidate digest is rejected.
- Provider timeout is reconciled before retry.
- Manual recovery can proceed without the failed worker or model provider.
