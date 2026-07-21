# Trustworthy intent-to-PR vertical slice

Status: Accepted
Owner: Forge Product and Architecture
Version: 0.1

## Objective

Prove that Forge can turn one human request into one reviewable pull request through a complete, independently verifiable trust path.

The slice validates Forge as a development-team application before broadening the workforce or accepting the Workplace Orchestrator commission.

## Included journey

1. An authenticated human selects one tenant, workspace, project and repository.
2. Forge records intent, constraints and acceptance criteria.
3. Policy grants a task-scoped Execution Envelope.
4. Knowledge retrieval produces a bounded Context Package under that envelope.
5. Forge proposes a delivery plan and obtains required approval.
6. One isolated worker creates a branch and prepares a bounded repository change.
7. Independent checks evaluate the candidate change.
8. A human approves publication of the immutable candidate digest.
9. Forge pushes the branch and opens a pull request.
10. Forge records an Evidence Record and explains outcome, limitations and rollback.

Merge, production release, autonomous remediation, multi-agent optimization and Workplace-specific code are excluded.

## Participating capability contracts

The detailed command, query, event, ownership, error, consistency, and conformance baseline is defined in [First-slice capability contracts](first-slice-capability-contracts.md).

| Capability | Owns | Receives | Produces | Explicitly does not own |
|---|---|---|---|---|
| Tenant and workspace | workspace membership, project/repository binding | authenticated principal | immutable tenant/workspace/project context | external identity lifecycle |
| Identity and delegation | principal, workload and agent identity; delegation chain | identity-provider assertion | verified principal and scoped delegation | project membership or policy outcome |
| Intent and delivery | intent, acceptance criteria, delivery-run state | human request and workspace context | delivery plan, state transitions and candidate digest | authorization or tool execution |
| Policy decision | evaluated rules and obligations | identity, action, resource, classification and risk | versioned Policy Decision | execution or approval |
| Context assembly | Context Package manifest | envelope, authorized sources and objective | bounded, reproducible context references | canonical source content or authority |
| Tool brokerage | tool metadata and invocation mediation | envelope and validated tool request | normalized result and invocation evidence | tool implementation or policy decision |
| Isolated execution | ephemeral workspace and execution state | envelope, candidate source revision and plan | candidate artefacts, test results and digest | durable project state or approval |
| Verification | independent test and review results | immutable candidate digest and declared criteria | findings and verification evidence | mutation of the candidate |
| Approval | approval lifecycle | actor, action, candidate digest, evidence and expiry | versioned Approval Record | candidate creation or policy evaluation |
| Audit evidence | append-oriented Evidence Records | referenced decisions, transitions and attestations | integrity-protected evidence manifest | operational telemetry or domain state |
| Source-control adapter | branch and pull-request interaction | approved, idempotent publication request | remote references and provider receipt | Forge policy or delivery state |
| Experience | conversation and role-specific projections | capability APIs and evidence | human request, review and explanation | canonical delivery or evidence data |

## Trust-object relationships

~~~mermaid
flowchart LR
    Principal --> PolicyRequest
    Workspace --> PolicyRequest
    PolicyRequest --> PolicyDecision
    PolicyDecision --> Envelope[Execution Envelope]
    Envelope --> Context[Context Package]
    Context --> Plan[Delivery Plan]
    Envelope --> Worker[Isolated Worker]
    Plan --> Worker
    Worker --> Candidate[Candidate Digest]
    Candidate --> Verification
    Verification --> Approval[Approval Record]
    Approval --> Publication[PR Publication]
    Publication --> Evidence[Evidence Record]
    PolicyDecision --> Evidence
    Envelope --> Evidence
    Verification --> Evidence
    Approval --> Evidence
~~~

- A Context Package references exactly one Execution Envelope.
- A Policy Decision may authorize envelope issuance or a later action gate.
- An Approval Record is bound to one immutable subject digest and action.
- An Evidence Record references objects; it does not copy secrets or unrestricted source content.
- A tool, model, repository or context item cannot expand the envelope.

Initial schemas live under `schemas/forge/trust/v1/` and remain Proposed until ADR-0003 is Accepted.

## Initial validation envelope

These figures are bounded hypotheses for the first slice, not production commitments:

| Dimension | Initial validation range |
|---|---|
| Tenancy | one test tenant, one workspace and one project per run |
| Repository | Git repository up to 100,000 tracked files and 5 GiB checkout size |
| Concurrency | 5 active delivery runs and 2 worker executions on one developer host |
| Context manifest | at most 2,000 referenced items and 10 MiB serialized metadata; content is fetched selectively |
| Delivery duration | 30 minutes expected, 2 hours hard execution expiry unless renewed |
| Policy decision | local p95 below 100 ms under 20 decisions/second |
| Control operations | local p95 below 500 ms excluding provider and repository latency |
| Evidence | at most 10,000 records per delivery run; append p95 below 250 ms locally |
| Recovery | no loss of committed delivery transitions; interrupted worker can be safely retried or manually recovered |
| Cost | declared per-run model, tool, compute and time budgets; no unbounded fallback |

Failure to meet a hypothesis triggers measurement and architecture review, not premature service extraction.

## Mandatory failure behaviour

| Failure | Required behaviour |
|---|---|
| Identity or workspace cannot be verified | Deny before retrieval or execution |
| Policy engine unavailable | Fail closed for retrieval, secrets, mutation and publication |
| Context source unavailable | Mark missing evidence; stop or continue only if the plan remains valid within policy |
| Context contains embedded instructions | Treat as data; never widen authority |
| Worker timeout or crash | Preserve durable control state, revoke task credentials and allow bounded retry |
| Tool or model exceeds quota | Stop or use a pre-authorized fallback; record decision and cost |
| Candidate changes after verification | Invalidate verification and approval |
| Approval expires or digest differs | Deny publication |
| Source-control publication is ambiguous | Reconcile idempotently before retrying |
| Evidence sink unavailable | Block consequential publication unless policy explicitly permits deferred evidence with a durable local outbox |

## Security and privacy minimum

The complete initial inventory, processing gates, trust boundaries, minimization rules, and privacy follow-ups are defined in [First-slice data flow and privacy analysis](first-slice-data-flow.md).

- Tenant and workspace identifiers are server-derived and immutable within the run.
- Credentials are short-lived references resolved only inside the authorized worker boundary.
- Worker filesystem and network access are allow-listed and time-bounded.
- Source, prompts, secrets and personal data are excluded from routine telemetry.
- Context retrieval is authorized before ranking or content return.
- Candidate code and generated tests cannot be their own sole verification oracle.
- The worker cannot approve, push, merge or change policy without separate grants.
- Evidence disclosure applies classification and retention policy.

## Exit criteria

The slice is complete only when:

- trust schemas validate every run object;
- negative tenant, digest, expiry and authority tests pass;
- one interrupted run is recovered without duplicate publication;
- one source-control provider can be replaced by a contract-test double;
- context can be reproduced from its manifest and source revisions;
- independent verification and human approval bind the exact candidate digest;
- the Workbench or CLI explains intent, assumptions, changes, checks, cost, risk and rollback;
- one versioned scenario exercises that same public Workbench or CLI boundary and cannot alter its own expected outcome;
- no claim is made that excluded capabilities are implemented.
