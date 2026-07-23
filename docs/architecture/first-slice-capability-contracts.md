# First-slice capability contracts

Status: Accepted
Owner: Forge Product and Architecture
Version: 0.2

## Scope

This document defines the minimum logical contracts for the trustworthy intent-to-pull-request vertical slice. It refines the participating-capability table in `vertical-slice.md`; it does not claim implementation.

Merge, production release, autonomous remediation, billing, broad workforce coordination, and Workplace-specific behaviour remain excluded.

## Shared contract rules

Every command, query, event, and result carries:

- schema and contract version;
- immutable tenant, workspace, project, and correlation scope where applicable;
- authenticated actor and delegation reference;
- idempotency key for state change or external side effect;
- classification and evidence handling metadata;
- normalized success or error outcome.

Identifiers used for authorization are server-derived or resolved from capability-owned records. Prompts, retrieved content, repositories, models, tools, and caller-supplied metadata cannot authoritatively set tenant, project, principal, action, approval, or policy scope.

### Normalized error contract

Errors SHALL contain a stable code, safe message, correlation identifier, retry classification, and optional structured details that have passed disclosure policy. Initial codes are:

| Code | Meaning | Default retry posture |
|---|---|---|
| `invalid_request` | Schema, state, or invariant failure | Do not retry unchanged |
| `unauthenticated` | Identity assertion is absent or invalid | Reauthenticate |
| `forbidden` | Policy denied the exact action | Do not retry unchanged |
| `indeterminate` | Required policy or authority could not be resolved | Fail closed |
| `conflict` | Version, idempotency, or remote-state conflict | Reconcile first |
| `expired` | Envelope, delegation, or approval expired | Obtain new authority |
| `budget_exhausted` | Time, token, compute, cost, or attempt limit reached | New approval or budget required |
| `dependency_unavailable` | Required provider or capability is unavailable | Bounded retry if authorized |
| `verification_failed` | Candidate did not satisfy declared criteria | Correct candidate |
| `evidence_unavailable` | Required evidence cannot be committed | Block consequential action |
| `manual_recovery_required` | Automated reconciliation is unsafe | Human action required |

Unknown errors fail closed at authority, retrieval, secret, mutation, disclosure, and publication gates.

## Canonical ownership

| Capability | Canonical records | Records it must not own |
|---|---|---|
| Tenant and workspace | tenants, workspaces, projects, memberships, repository bindings | external identity lifecycle, policy decisions |
| Identity and delegation | normalized principals, assertions, delegation chains, revocation | membership, final authorization |
| Intent and delivery | intent, acceptance criteria, delivery plan, run state, candidate reference | policy, worker state, approval |
| Policy decision | policy versions, decisions, obligations, reason codes | execution, approval, business state |
| Context assembly | Context Package manifest and assembly result | canonical source, authority, durable memory |
| Tool brokerage | tool definitions, invocation requests, normalized receipts | tool implementation, policy semantics |
| Isolated execution | ephemeral workspace state and structured execution result | durable delivery state, approval |
| Verification | findings, criteria results, candidate verification result | candidate mutation, release approval |
| Approval | approval request and immutable Approval Record | candidate creation, policy evaluation |
| Audit evidence | Evidence Records and integrity relationships | domain state, unrestricted telemetry |
| Source-control adapter | provider operation receipt and reconciled remote reference | Forge policy, delivery state |
| Experience | user input and authorized projections | canonical delivery, policy, or evidence data |

## Commands, queries, and events

### Tenant and workspace

- Query `ResolveWorkspaceScope(principal, requestedWorkspace, requestedProject, repository)` returns immutable scope or a denial.
- Query `GetRepositoryBinding(scope)` returns an authorized immutable repository binding.
- Event `WorkspaceMembershipChanged` invalidates affected sessions and pending authority according to policy.
- Cross-tenant negative queries SHALL return no existence-revealing content.

### Identity and delegation

- Command `VerifyIdentity(assertion, audience)` returns a normalized principal and assurance metadata.
- Command `CreateDelegation(parentAuthority, subject, purpose, actions, resources, expiry)` creates strictly narrower authority.
- Command `RevokeDelegation(delegationId, reason)` records revocation and triggers propagation.
- Query `ResolveDelegationChain(id)` returns a disclosure-safe verified chain.
- Events include `PrincipalRevoked`, `DelegationIssued`, and `DelegationRevoked`.

### Intent and delivery

- Command `CaptureIntent(scope, actor, objective, constraints, acceptanceCriteria)` creates one delivery run.
- Command `RecordContext(runId, expectedVersion, contextPackageId, manifestDigest, evidenceRefs)` binds the exact validated Context Package before planning.
- Command `RecordPlan(runId, expectedVersion, plan)` records a canonical context-informed delivery graph and digest.
- Command `TransitionDelivery(runId, expectedVersion, transition, evidenceRefs)` performs one valid state transition and writes an outbox record atomically.
- Command `RecordCandidate(runId, expectedVersion, digest, artefactRefs)` binds the candidate without accepting it.
- Query `GetDeliveryRun(runId)` returns an authorized role-specific view.
- Events reflect committed transitions and never precede durable state.

### Policy decision

- Command `EvaluatePolicy(subject, action, resource, scope, purpose, classification, risk, environment)` returns the canonical Policy Decision schema.
- Decisions are `permit`, `deny`, or `indeterminate`; indeterminate is never permit.
- Mandatory unknown obligations cause denial.
- Policy versions, input digests, expiry, obligations, and reason codes make the decision reproducible without copying sensitive input into evidence.

### Context assembly

- Command `AssembleContext(envelopeId, objective, authorizedSources, limits)` returns a Context Package manifest.
- Query `ReproduceContext(packageId)` resolves the same authorized source revisions or reports explicit omissions and staleness.
- Assembly authorizes before discovery and content return.
- Package mutation cannot add actions, destinations, tools, models, secrets, or classifications.

### Tool brokerage

- Query `ResolveTool(toolId, version)` returns identity, schema, risk, side effects, health, and compatibility metadata.
- Command `InvokeTool(envelopeId, policyDecisionId, toolId, arguments, idempotencyKey)` validates authority, arguments, destination, limits, and obligations before dispatch.
- Results include normalized outcome, provider receipt, cost/resource use, and evidence references.
- The broker reconciles ambiguous side effects before retry.

### Isolated execution

- Command `StartWorkerTask(envelopeId, planStep, sourceRevision, contextPackageId)` starts one bounded task.
- Command `CancelWorkerTask(taskId, reason)` terminates or marks non-cancellable work for compensation.
- Result `WorkerTaskResult` contains input digests, output artefact references, candidate digest, checks performed, resource use, and cleanup result.
- Workers cannot directly transition delivery state, approve, publish, or resolve unrestricted credentials.

### Verification

- Command `VerifyCandidate(candidateDigest, criteria, independentRunner)` evaluates the immutable candidate.
- Result `VerificationResult` records each criterion, finding, tool and data version, independence properties, and evidence references.
- A changed candidate digest invalidates the result.
- Candidate-controlled tests cannot be the sole oracle for consequential criteria.

### Approval

- Command `RequestApproval(actor, action, subjectDigest, environment, evidenceRefs, expiry)` opens a review without granting authority.
- Command `DecideApproval(requestId, approver, decision, conditions)` produces an Approval Record.
- Approval is valid only for the bound actor or role, action, digest, environment, conditions, and time window.
- Plan review uses action `forge.delivery.execute-plan`, subject type `delivery-plan`, and the exact canonical plan digest. It cannot authorize execution or publication.
- Publication approval is a separate request bound to the independently verified candidate digest.
- The candidate creator cannot satisfy a required independent approval role.

### Audit evidence

- Command `AppendEvidence(record)` validates classification, references, retention, and integrity before durable append.
- Query `VerifyEvidence(recordId)` verifies subject digests and integrity relationships independently of the originating capability.
- Query `ExportEvidence(scope, purpose, policyDecision)` is a separately authorized disclosure.
- Evidence stores references and digests rather than unrestricted source, prompt, secret, or personal content.

### Source-control adapter

- Command `PublishPullRequest(envelopeId, policyDecisionId, approvalRecordId, candidateDigest, repositoryBinding, branch, metadata, idempotencyKey)` performs the only included external publication.
- Before mutation, the adapter verifies exact digest, approval, repository binding, branch constraints, and required evidence readiness.
- On timeout or ambiguous response it queries remote state using the idempotency key before retry.
- Result contains immutable provider receipt, branch reference, pull-request reference, and observed remote digest.

### Experience

- Command `SubmitIntent` maps authenticated user input to `CaptureIntent`; it cannot set authoritative scope.
- Commands `ReviewPlan` and `DecideApproval` expose the exact subject, context revisions, omissions, conflicts, assumptions, checks, limitations, cost, risk, and rollback.
- Queries return role-appropriate projections with source, version, freshness, and authorization scope.
- The first interface MAY be a CLI, but the same public boundary SHALL be usable by the accepted journey scenario.

## State and consistency boundaries

- Delivery state uses optimistic concurrency and transactional outbox semantics.
- Delivery follows `Captured -> ContextAuthorizing -> ContextReady -> Planned -> AwaitingPlanApproval -> PlanApproved -> ExecutionAuthorizing`; later execution and publication states preserve their separately named gates.
- Context identity or manifest changes invalidate the plan; plan digest changes invalidate plan approval and later authority.
- Internal events are at least once; consumers are idempotent.
- External publication is idempotent and reconciled against provider state.
- Evidence required for publication is committed before or atomically coordinated with the publication transition.
- No transaction spans capability-owned stores or an external provider.
- Manual recovery is a designed result when safe automatic reconciliation is impossible.

## Security and privacy obligations

- Authenticate before membership resolution; authorize before retrieval or action.
- Derive tenant and project scope on the server and bind it through storage, caches, queues, workers, telemetry, and evidence.
- Resolve short-lived secret references only inside the authorized worker boundary.
- Deny worker egress by default and permit exact destinations only.
- Exclude secrets, unrestricted source, prompts, personal data, and tool output from routine telemetry.
- Treat repository content, context items, provider output, tool output, generated code, and generated tests as untrusted.
- Re-evaluate policy when a relevant input or protected resource changes.

## Initial quality hypotheses

The validation envelope and performance hypotheses in `vertical-slice.md` apply. In addition:

- identity, policy, context, approval, and publication denial paths fail closed;
- cancellation revokes resolvable task credentials within the declared propagation bound;
- every committed delivery transition survives worker loss;
- the same idempotency key cannot create two pull requests;
- evidence integrity verification detects subject or reference mutation;
- context reproduction exposes missing or changed sources rather than silently substituting content;
- a contract-test source-control double can replace the first provider without domain changes.

These are Proposed hypotheses until tests and an accepted support envelope establish otherwise.

## Conformance evidence required before Trial

- schema validation for every public object;
- state-transition and optimistic-concurrency tests;
- cross-tenant, expired-authority, changed-digest, and unknown-obligation negative tests;
- context authority-separation tests;
- context-informed plan ordering, digest invalidation, and plan-approval replay tests;
- worker isolation, cancellation, quota, and egress tests;
- independent verification and self-approval denial tests;
- publication idempotency and ambiguous-result reconciliation tests;
- evidence integrity, disclosure, and unavailable-sink tests;
- one allowed and one denied accepted journey through the same public interface.
