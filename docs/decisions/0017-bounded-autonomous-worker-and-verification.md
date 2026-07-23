# ADR-0017: Establish the bounded autonomous worker and independent verifier

Status: Accepted
Owner: Forge Execution, Verification, Security, and Architecture
Date: 2026-07-23

## Context

M1 must turn an approved context-informed plan into one candidate change without allowing the task agent, repository, generated content, worker, or candidate-controlled tests to become its own authority or sole acceptance oracle. ADR-0014 currently permits only one context-preparation delegation. ADR-0015 stops the next implementation at `PlanApproved` and requires separately designed worker authority and isolation.

Trust is an enabling platform property, not an objective to maximize independently of delivery. Requiring a human decision for every file operation, check, retry, or state transition would prevent Forge from functioning as an AI-supported development team. Conversely, allowing the worker to infer new actions from repository content or its own output would turn autonomy into self-granted authority.

The Developer Trial security decision requires a separate process with bounded filesystem, network, time, resource, credential, cancellation, and cleanup scope. The current profile already requires Docker for local operation, but no accepted decision assigns Docker a worker-isolation role or defines the boundary between candidate preparation and independent verification.

## Decision

Use two short-lived, separately authorized Docker-isolated execution roles for the Developer Trial:

1. one deterministic change-preparation worker; and
2. one independent verifier operating from an immutable candidate digest.

Docker is a profile-specific isolation adapter, not a capability contract or production topology commitment.

### Autonomy posture

After one valid human approval of the exact Delivery Plan, Forge MAY autonomously perform every deterministic step declared by that unchanged plan within its separately issued execution authority. It does not request additional human approval for:

- exact read-only repository inspection already declared by the plan;
- deterministic file transformations within declared paths;
- bounded internal retries that preserve identical inputs and side-effect identity;
- checkpoints, progress reporting, cancellation handling, cleanup, and evidence emission;
- independent controls already selected by the accepted Assurance Plan; or
- recovery to a previously committed safe checkpoint when no external outcome is ambiguous.

Human attention returns only when:

- a normative plan, context, scope, destination, action, classification, risk, budget, tool, model, secret, or candidate digest changes;
- a mandatory guardrail fails, expires, becomes stale, or is indeterminate;
- a side effect or recovery outcome cannot be reconciled safely;
- policy requires a named human decision; or
- publication of the verified candidate is requested.

Controls are concentrated at trust and side-effect boundaries rather than repeated inside every pure transformation. Cheap deterministic admission checks run before execution. Independent verification runs in parallel where dependencies and budgets permit. Exact valid evidence may be reused under ADR-0013. No authority, isolation, secret, integrity, or approval check is sampled away.

### Phase-specific worker authority

After `PlanApproved`, Identity MAY issue a second phase-specific root delegation to one verified synthetic worker workload principal. It is a sibling of the context-preparation delegation, not a child that widens the context action. Its parent authority is a new exact Policy Decision under the bootstrap pattern of ADR-0014.

The worker-delegation request is bound to:

- the initiating human, task agent, and worker workload identities;
- exact tenant, workspace, project, repository, and delivery run;
- current Delivery Plan identifier and digest;
- current Context Package identifier and manifest digest;
- current plan Approval Request and Approval Record v2 identifiers and digests;
- immutable source revision and declared change paths;
- action `forge.worker.prepare-change`;
- purpose `forge.m1.prepare-change`;
- environment `developer-trial`;
- classification, risk, destinations, filesystem, network, model, tool, secret, time, resource, concurrency, attempt, compute, token, and cost bounds; and
- expiry no later than every parent fact and the Developer Trial maximum.

Policy evaluates action `forge.delegation.create` against the canonical request digest. Identity owns issuance, lifecycle, revocation, and resolution. This specializes ADR-0014 from one context root delegation to multiple non-overlapping phase delegations; it does not introduce general child delegation, delegation inheritance, renewal, offline authority, or dynamic expansion.

Policy then issues a worker Execution Envelope under ADR-0012 for action `forge.worker.prepare-change`. The envelope references the exact worker delegation, plan approval, plan digest, context, source revision, policy decisions, declared paths, resource bounds, and evidence correlation. Possession of the plan approval or context envelope cannot substitute for worker authority.

### Worker contract

The worker receives immutable identifiers and digests plus an ephemeral authorized checkout. Its versioned request declares:

- worker envelope and delegation identifiers;
- plan step identifiers and dependency completion digests;
- source revision, Context Package, and plan identifiers and digests;
- allowed read and write paths;
- deterministic transformation identifier and version;
- required input artefact digests;
- resource, duration, attempt, process, filesystem, and output bounds;
- cancellation and cleanup correlation; and
- expected structured output contract.

The first M1 transformation is deterministic and Forge-owned. It may create or modify only declared repository paths and may not invoke a model, general tool broker, shell supplied by repository content, package installer, network destination, secret, or production credential. Repository content and embedded instructions remain data.

The worker returns a structured result containing input digests, changed-path manifest, patch or candidate artefact references, resulting tree and candidate digest, resource use, bounded logs, cleanup result, and evidence references. It cannot transition Delivery state, approve, verify, publish, resolve unrestricted credentials, alter policy, or write control-plane storage.

### Developer Trial isolation

The launcher creates one fresh Docker container per task with:

- a pinned Forge-owned worker image reference and declared runtime version;
- non-root user, read-only root filesystem, dropped Linux capabilities, no privileged mode, no host socket, and bounded process count;
- one ephemeral workspace containing only the authorized checkout and output mount;
- default-deny network using Docker network isolation;
- explicit CPU, memory, process, filesystem, output-size, and wall-clock limits;
- no inherited environment credentials and no secret mount for the deterministic first worker;
- cancellation by control-plane request plus hard deadline termination;
- deterministic cleanup of container, workspace, temporary artefacts, and credential handles; and
- a structured cleanup outcome that can enter manual recovery.

The launcher adapter verifies effective configuration before start and records safe evidence. Unsupported Docker hosts, unavailable enforcement, image mismatch, cleanup failure, or inability to prove the declared boundary blocks execution. Docker Desktop and host kernels differ; M1 evidence names the tested host envelope and makes no production portability claim.

### Independent verification

Verification owns a separate immutable Assurance Plan derived under ADR-0013 from the candidate digest, affected capabilities and paths, accepted quality envelope, risk, classification, and mandatory guardrails.

The verifier:

- starts in a fresh container and workspace without the worker's writable filesystem;
- receives the immutable source revision, candidate artefact, candidate digest, acceptance criteria, and Assurance Plan;
- applies the candidate to a fresh checkout and recomputes its digest;
- runs Forge-owned contract, architecture, unit, integration, journey, security, privacy, and performance controls applicable to the change;
- treats candidate-controlled tests as evidence inputs but never as the sole consequential oracle;
- cannot mutate the candidate or Delivery plan; and
- returns versioned findings and a Verification Result bound to the exact candidate digest, control versions, environment, and evidence.

Candidate-controlled code executes only when an applicable control explicitly requires it and only inside the verifier isolation envelope. A candidate that changes after verification invalidates the result and every dependent approval.

### Cancellation, retry, and recovery

Delivery owns durable scheduling and state. Workers and verifiers are disposable:

- cancellation prevents new starts, signals active containers, revokes resolvable phase authority, and blocks later publication;
- retry uses the same immutable inputs and side-effect identity within declared attempt, elapsed-time, and resource budgets;
- worker loss never commits candidate success without a complete validated result;
- verifier loss never converts missing evidence into pass;
- cleanup failure is visible and blocks unsafe reuse;
- ambiguous external state is not present in deterministic preparation or verification; if introduced later it requires reconciliation; and
- manual recovery is required when authority, cleanup, candidate identity, or evidence cannot be proven.

### Capability and deployment boundaries

- Delivery schedules and records durable state but does not execute repository code.
- Identity owns worker principal and delegation lifecycle.
- Policy owns delegation permits and worker Execution Envelope issuance.
- Execution owns disposable workspace and structured task results.
- Verification owns Assurance Plans, independent controls, findings, and Verification Results.
- Evidence owns immutable evidence records, not telemetry or domain state.
- the worker deployment unit cannot import control-plane runtime modules; it depends only on shared worker and trust contracts.
- the control plane cannot import worker implementation.

The worker and verifier MAY share one Forge-owned image initially, but they run as separate containers with different entry points, mounts, authority, and writable state. Separate images or deployment units require evidence or a later ADR.

### Iteration boundary

After this decision and its persistence dependencies are Accepted, a bounded implementation MAY progress from `PlanApproved` through worker preparation, candidate recording, `Verifying`, and a verified candidate outcome. It SHALL NOT publish a branch or pull request, use a live source-control credential, merge, release, deploy, invoke a model, enable a general tool broker, or claim production isolation.

## Alternatives considered

- **Approve every worker step separately:** rejected because unchanged low-risk steps already reviewed in the immutable plan should execute autonomously inside one envelope.
- **Give the task agent one broad M1 delegation:** rejected because long-lived cross-phase authority increases blast radius and weakens phase-specific revocation.
- **Create a child delegation from the context-only delegation:** rejected because a child cannot add the worker action absent from its parent.
- **Run the worker in the control-plane process:** rejected because repository and candidate behaviour must not share durable authority and state.
- **Use an ordinary host subprocess only:** rejected because the initial Developer Trial requires enforceable filesystem, network, resource, process, and cleanup boundaries across supported developer hosts.
- **Use Kubernetes immediately:** rejected because the Developer Trial needs one local worker and has no evidence justifying orchestration complexity.
- **Let the worker verify its own output:** rejected because it creates a circular oracle.
- **Prevent all candidate-controlled execution:** rejected because meaningful verification may need builds or tests; isolation and independent control selection address the risk.

## Consequences

- One human plan decision enables a meaningful autonomous preparation and verification window.
- Phase-specific authority reduces blast radius at the cost of additional deterministic Policy and Identity operations.
- Docker becomes a measured Developer Trial adapter and prerequisite, not a universal deployment choice.
- Separate worker and verifier runs consume more compute but preserve independent evidence.
- The first transformation remains deterministic; agentic/model-assisted transformation requires M3 contracts or a separately approved M1 adapter.
- Host-specific isolation limits and performance must be measured and explained.

## Validation

Evidence SHALL prove:

- exact phase-delegation and envelope binding to identities, scope, plan, approval, context, source revision, paths, action, environment, limits, and expiry;
- denial after plan, context, approval, policy, identity, source, classification, risk, or budget mutation;
- architecture fitness for control-plane/worker dependency and deployment-unit isolation;
- non-root, dropped-capability, read-only-root, no-socket, no-network, bounded-resource, timeout, cancellation, and cleanup enforcement;
- path traversal, symlink escape, environment leakage, credential inheritance, process escape, egress, fork-bomb, oversized-output, timeout, and cleanup negative tests;
- candidate digest reproduction from immutable inputs;
- independent Assurance Plan selection and verification;
- candidate-controlled tests are not the sole oracle and cannot change expected outcomes;
- worker/verifier loss, bounded retry, revocation, cancellation, checkpoint, and manual-recovery behaviour;
- progress remains visible while background work runs;
- p50, p95, p99, throughput, cold start, saturation, CPU, memory, storage, failure, and cleanup measurements under ADR-0013; and
- no worker, verifier, task agent, repository content, model, or tool can approve or publish.

## Accountable review

Sijbren Beukenkamp accepted this decision as drafted through the authenticated collaboration boundary on 2026-07-23, acting as the temporary accountable Architecture, Security and Privacy, Identity, Policy, Delivery, Execution, Verification, Quality, Operations, and Engineering reviewer established for bounded M1.

Functional implementation remains blocked until the decision reaches protected `main` and a separate bounded grant is activated.

## Supersession

If accepted, this decision specializes ADR-0012, ADR-0013, ADR-0014, and ADR-0015 for phase-specific worker authority and independent verification. It does not authorize M3 model or tool capabilities, production isolation, publication, release, or deployment.
