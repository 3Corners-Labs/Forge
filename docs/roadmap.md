# Forge implementation roadmap

Status: Accepted
Owner: Forge Product and Architecture
Version: 0.6

## Purpose

This roadmap evolves Forge from the minimum platform capable of preparing a trustworthy change to itself into a bounded means for delivering and maintaining 3Corners applications, beginning with Workplace Orchestrator.

[ADR-0024](decisions/0024-rebaseline-forge-as-bounded-means.md) rebaselined this roadmap. Forge is no longer developed as a product for independent sale. M1 through M3 are committed scope; M4 through M9 were demoted to evidence label Vision and are recorded as direction under [Vision direction (M4–M9)](#vision-direction-m4m9).

It describes target delivery sequencing, not implemented state or calendar commitments. Capability maturity remains separate from document status. A milestone becomes complete only through the stated acceptance evidence and accountable human review.

## Delivery rules

- Milestones deliver end-to-end outcomes, not disconnected service skeletons.
- Capability is pulled by a demonstrated delivery need, not pushed by completeness. Deferred scope is recorded with the need that would reinstate it.
- The initial implementation remains a modular control plane with isolated workers under accepted ADR-0001.
- A capability becomes a separate deployment unit only through measured evidence and an Accepted ADR.
- Models, tools, frameworks, providers, and protocols remain replaceable adapters.
- Every milestone preserves identity, policy, bounded context, evidence, human accountability, and separation of duties.
- Capabilities remain Research until a Trial is explicitly authorized. Trial evidence may justify promotion to Supported or Reference only through a separate accountable lifecycle decision.
- Production support requires a declared conformance envelope, not merely successful deployment.
- A later milestone SHALL NOT weaken evidence or controls accepted by an earlier milestone.

## Complexity scale

| Rating | Meaning |
|---|---|
| Medium | Bounded change within a small number of modules and one trust boundary |
| Large | Multiple capability contracts, persistent state, or one external integration boundary |
| Extra Large | Cross-plane delivery with isolation, security, recovery, and operational evidence |
| Programme | Multiple products or deployment profiles with organizational and operational change |

Complexity is relative and intentionally not converted into dates until team size, skills, delivery capacity, and support expectations are known.

## Entry conditions

Milestone 1 functional implementation begins only when the remaining applicable gates in the [architecture readiness register](governance/architecture-readiness.md) permit it. Preparatory architecture, contract tests, safety hardening, and fitness-function work may proceed to close those gates.

The accepted implementation boundary is defined by:

- [ADR-0006](decisions/0006-accept-minimum-architecture-baseline.md);
- [Trustworthy intent-to-PR vertical slice](architecture/vertical-slice.md);
- [First-slice capability contracts](architecture/first-slice-capability-contracts.md);
- [First-slice data flow and privacy analysis](architecture/first-slice-data-flow.md);
- [Forge platform threat model](architecture/threat-model.md).

## F0 — Foundation readiness gate

F0 is a pre-implementation gate, not a product milestone. M1 remains the first implementation milestone and the first minimum platform outcome.

G0 through G10 are Accepted for F0 at protected-main baseline `93efcfcae8e946310ebd07294d744f38a3083b0e`. Sijbren Beukenkamp accepted the [F0 final acceptance manifest](governance/f0-acceptance-manifest.md) and [ADR-0010](decisions/0010-authorize-bounded-m1-implementation.md) at review revision `773312e`. Bounded M1 authority becomes operationally effective when the Accepted packet reaches protected `main`.

### Objective

Preserve and complete the architecture-foundation work required to begin M1 without confusing accepted design, implementation maturity, and verification evidence.

### Deliverables and current status

| Foundation work | Current status | Evidence | Required disposition |
|---|---|---|---|
| Accept the constitution, minimum capability boundary, trust model, and foundational ADRs | Canonical target baseline Accepted | [ADR-0006](decisions/0006-accept-minimum-architecture-baseline.md) and [ADR-0008](decisions/0008-accept-canonical-documentation-baseline.md) | Keep unimplemented capabilities at Research; authorize Trials incrementally when a milestone requires them |
| Close the mandatory readiness gates | Complete for pre-M1 scope | [Architecture readiness register](governance/architecture-readiness.md), [F0 readiness plan](governance/f0-readiness-plan.md), and gate evidence | Preserve the accepted evidence and complete the separate final authorization decision |
| Reconcile document status, capability maturity, and evidence labels | Complete for the canonical baseline | Accepted [documentation lifecycle](governance/document-lifecycle.md), ADR-0008, and lifecycle fitness check | Maintain the separation without treating Accepted design as implemented capability |
| Accept or revise the metamodel, trust-object, and deployment-profile decisions | Complete for F0 | ADR-0002, ADR-0003, ADR-0004, and accepted G7 evidence | Do not treat profile composition as production support |
| Define non-functional requirements and deployment profiles | Complete for F0 | First-slice validation hypotheses, Accepted deployment architecture, and G7 evidence | Validate M1-exit and Before-Trial evidence in their assigned phases |
| Establish documentation, contract, dependency, isolation, and architecture fitness checks | Complete for F0 | Required protected CI and accepted G8 evidence | Apply them to implemented M1 modules; M2 expands them into supply-chain assurance |
| Define accepted journey, exploratory scenario, and offensive-test rules-of-engagement contracts | Complete for F0 | Accepted ADR-0005, G9 scenarios, and G10 safety contract | Execute M1 journeys at M1 exit; broad active testing remains excluded from M1 |
| Replace development defaults that could be mistaken for production posture | Complete for F0 | Accepted G0 evidence, loopback checks, and runtime-profile rejection | Preserve the restriction and validate the implemented M1 public and worker paths |
| Review and accept or revise the product masterplan | Accepted | Accepted [Forge mega masterplan](foundation/masterplan.md) and ADR-0008 | ADR-0007 remains authoritative for implementation sequencing and the bounded M1 scope |
| Complete source-content migration without loss | In progress | [Content migration map](governance/content-migration.md) and source map | Preserve provenance, finish concept-level verification, and prevent source removal until migration rules pass |
| Establish capability owners, quality scenarios, threat coverage, and the ADR register | Partially complete for M1 | Accepted first-slice contracts, threat model, ADR register, and validation hypotheses | Assign implementation and operational owners and close missing quality scenarios before authorizing each Trial |

### Dependencies

- Accepted constitution, documentation lifecycle, metamodel, ADR-0001 through ADR-0008, and first-slice architecture.
- Accountable Product, Architecture, Security/Privacy, Engineering, and Operations review.
- Current repository and implementation evidence from the implementation baseline and CI.

### Acceptance criteria

- The architecture readiness register explicitly authorizes M1 functional implementation; every mandatory gate is `Accepted for F0` and none remains Open, In progress, or Ready for review.
- Minimum dependency, isolation, schema, source-integrity, documentation, developer-profile, journey, and safety checks pass in required CI.
- Trial decisions identify identity, repository credentials, retention, evidence integrity, worker isolation, public-interface journeys, and deployment assumptions.
- The masterplan and roadmap either agree or retain an explicit accepted decision explaining each remaining difference.
- Source content remains discoverable with provenance, and no source is removed before its migration disposition is verified.
- Capability owners and accountable reviewers are named for the M1 Trial.

### Estimated complexity

**Large.** F0 is primarily architecture, governance, contract, safety, and validation work, but it crosses every trust boundary used by M1.

### Risks

- Treating document acceptance as proof that controls are implemented.
- Turning readiness into an unbounded documentation programme rather than closing the minimum M1 gaps.
- Silently deferring mandatory checks to later milestones.
- Losing source intent during consolidation.
- Allowing one temporary multi-role reviewer arrangement to become permanent without review.

## Milestone progression

| Milestone | Outcome | Capability posture at exit | Status |
|---|---|---|---|
| M1 | Forge prepares one trustworthy change to Forge | Minimum self-building kernel | Committed |
| M2 | Forge builds and verifies reproducibly, to the extent delivery requires | Demand-driven engineering foundation | Committed |
| M3 | Forge routes across deterministic, local, and external intelligence | Sovereign model routing | Committed |
| M4–M9 | Recorded direction, not committed scope | See [Vision direction](#vision-direction-m4m9) | Vision |

Workplace Orchestrator delivery begins after M1 exit rather than at the former M6. Its feasibility work may proceed in parallel, because it tests a product-existence assumption rather than a Forge capability.

## Sequencing decision

The masterplan originally places broad read-only observation and advisory capability before Forge prepares changes to itself. This roadmap deliberately starts with the narrower accepted intent-to-pull-request kernel required by the user and ADR-0006.

M1 includes only the repository inspection, bounded context, planning, and explanation needed to prepare one governed Forge change. [ADR-0007](decisions/0007-sequence-self-building-before-full-advisory.md) records this distinction so the roadmap does not silently contradict the masterplan.

The general Software System Twin and multi-audience advisory product were the former M4. [ADR-0024](decisions/0024-rebaseline-forge-as-bounded-means.md) demoted them to Vision: Forge builds the systems it explains, so explaining a system it did not build is not yet a delivery need.

---

## M1 — Minimum self-building kernel

### Objective

Enable Forge to transform one authenticated human request into one reviewable pull request against the Forge repository through a complete, independently verifiable trust path. This is the minimum platform capable of building itself.

Forge prepares changes only. It cannot merge, release, change production policy, approve itself, or acquire production credentials.

M1 assembles the bounded Context Package before constructing the canonical Delivery Plan. Human plan review binds that exact plan digest but does not grant worker authority. Policy separately authorizes execution, and publication later requires a different human approval bound to the independently verified candidate digest, as established by ADR-0015.

### Deliverables

- Modular control-plane foundation with enforced capability and dependency boundaries.
- Minimal tenant, workspace, project, membership, and repository binding for a Developer Trial.
- Normalized human and task-agent identity with explicit, expiring delegation.
- Executable validation for Policy Decision, Execution Envelope, Context Package, Approval Record, and Evidence Record.
- Deterministic policy decision and enforcement at retrieval, worker, secret, verification, approval, evidence, and publication gates.
- Durable delivery-run state machine with optimistic concurrency, cancellation, bounded retry, checkpoints, transactional outbox, and manual recovery.
- Purpose-bound repository context assembly with source revision, provenance, omissions, conflicts, and reproducible manifest.
- One isolated worker able to prepare a bounded Forge repository change using a narrow change-preparation contract. The implementation MAY use deterministic transformation or one approved model adapter; general model routing is excluded.
- Independent candidate verification bound to an immutable digest.
- Human publication approval bound to action, repository, environment, digest, conditions, and expiry.
- Idempotent source-control adapter that creates a branch and pull request and reconciles ambiguous provider results.
- Append-oriented evidence records with implemented integrity verification.
- Public CLI workflow explaining intent, assumptions, changes, checks, cost, risk, limitations, approval, and rollback.
- One accepted allowed journey and one denied journey through the same CLI boundary.

### Dependencies

- F0 accepted as complete and the architecture readiness register explicitly authorizes M1 functional implementation.
- Accepted constitution, metamodel, first-slice architecture, trust-object separation, capability contracts, threat model, and data-flow analysis.
- Approved identity, repository credential, evidence-integrity, retention, and worker-isolation Trial decisions.
- Developer profile capable of running deterministic tests without production credentials.

### Acceptance criteria

- Every exit criterion in the accepted vertical-slice architecture passes.
- Cross-tenant, expired-authority, changed-digest, unknown-obligation, and approval-replay tests fail closed.
- Context mutation cannot add authority, tools, models, destinations, classifications, or secret references.
- Worker loss preserves committed delivery state, revokes task credentials, and permits bounded recovery.
- Candidate-controlled code and tests cannot alter the independent acceptance oracle.
- One source-control timeout is reconciled without duplicate branch or pull-request creation.
- Evidence mutation is detected and required evidence unavailability blocks publication.
- Forge prepares a real, reviewable improvement to Forge on an isolated branch; an accountable human remains the only publication approver.
- The implementation baseline is updated with exact source and verification evidence.

### Remaining delivery tranches

The remaining M1 work is delivered as a small number of end-to-end tranches rather than one pull request per capability or control:

| Tranche | Outcome | Affected capabilities | Decision dependencies |
|---|---|---|---|
| M1-A Context to approved plan | Reproduce the exact context, construct the deterministic Delivery Plan, and obtain one human plan approval through `PlanApproved` | Context, Delivery, Policy, Approval, Experience | ADR-0015 and ADR-0016 |
| M1-B Approved plan to verified candidate | Persist and resume the run, autonomously prepare one bounded candidate in an isolated worker, independently verify its digest, and commit required evidence | Delivery, Identity, Policy, Execution, Verification, Evidence | ADR-0017 and ADR-0018 |
| M1-C Verified candidate to reviewable pull request | Guide free text into an explicitly promoted requirements proposal, obtain separate publication approval, reconcile one GitHub branch and pull request idempotently, and expose the allowed and denied public CLI journeys | Intent, Delivery, Policy, Approval, Evidence, Source Control, Experience | ADR-0016, ADR-0018, ADR-0019, and ADR-0020 |
| M1 exit | Exercise recovery, ambiguity, performance, leakage, isolation, and journey evidence and prepare one real Forge pull request | all participating M1 capabilities | every applicable Accepted M1 decision and exit contract |

One human plan approval opens the bounded autonomous preparation and verification window for the unchanged plan. One later human publication approval opens only the exact idempotent branch-and-pull-request operation. Ordinary declared steps, controls, checkpoints, progress, bounded retries, and exact reconciliation do not introduce micro-approval gates.

An implementation detail does not require a new ADR merely because it affects several modules. Work returns to architecture review only for a new or changed trust boundary, capability owner, durable provider, external destination, credential model, deployment unit, autonomous power, production claim, or contradiction with an Accepted decision.

### Estimated complexity

**Extra Large.** The feature surface is deliberately narrow, but it crosses identity, policy, state, isolation, repository, evidence, approval, and external side-effect boundaries.

### Risks

- Attempting the full capability catalogue instead of one complete path.
- Weak local isolation creating false security confidence.
- Circular trust if candidate code controls its verifier or policy.
- Provider credentials becoming broader or longer-lived than the task.
- Evidence integrity being reduced to an append-only storage claim.
- The temporary concentration of accountable roles reducing independent challenge.

---

## M2 — Demand-driven engineering and supply-chain foundation

### Objective

Make Forge builds, verification results, and promoted artefacts reproducible and
attributable to the degree that Workplace Orchestrator delivery actually
requires. Scope is pulled by delivery need, not pushed by completeness.

[ADR-0024](decisions/0024-rebaseline-forge-as-bounded-means.md) reduced this
milestone from a full supply-chain programme to demand-driven scope. Capability
that no delivery has requested is not built here.

### Deliverables

- Versioned developer-environment blueprint and pinned toolchain.
- Enforced source-module dependency rules derived from the accepted metamodel.
- Reproducible build producing artefacts identified by immutable digest.
- SBOM, dependency and license assessment, and signature or attestation for
  artefacts that leave the developer environment.
- Build-once promotion with rollback metadata.
- Secret, source, dependency, and configuration checks on every applicable build.
- Backup, restore, and interrupted-build rehearsal for the Developer Trial.

Deferred until a delivery need exists: ephemeral preview environments, full
provenance attestation chains, image and IaC scanning, and separate telemetry
and evidence pipelines. Each remains recorded in the Vision direction.

### Dependencies

- M1 candidate digest, verification, approval, publication, and evidence
  semantics.
- An accepted supply-chain decision covering builder identity, signing trust,
  artefact formats, and revocation, taken when the first artefact must leave the
  developer environment.

### Acceptance criteria

- Two clean builds from the same reviewed source and declared inputs produce the
  expected reproducibility result; permitted variance is documented.
- Changing one artefact byte causes digest or signature verification failure.
- A mutable tag cannot authorize promotion.
- Untrusted pull-request code cannot access signing material.
- Rollback deploys a previously verified artefact.
- Backup and restore reproduce declared state without silently resurrecting
  logically deleted data.
- Every deferred deliverable is recorded with the delivery need that would
  trigger it.

### Estimated complexity

**Medium.** Reduced from Large by ADR-0024. The remaining work is deterministic
and bounded to the developer environment.

### Risks

- Chasing bit-for-bit reproducibility before defining permitted variance.
- Deferral becoming permanent omission because no trigger was recorded.
- Signing keys being exposed to untrusted builds.
- Security scanning volume producing unactionable noise.

---

## M3 — Demand-driven intelligence and sovereign model routing

### Objective

Allow Forge to use deterministic methods, locally hosted models, and external
models through stable capability contracts, with sufficient routing to satisfy
the [ADR-0024](decisions/0024-rebaseline-forge-as-bounded-means.md) sovereignty
tiers. Policy, state, budgets, and verification remain deterministic.

The sovereignty tiers move locally hosted model support forward from the
previous plan. This milestone delivers the measured local path; it does not
deliver a general model registry.

### Deliverables

- Versioned Task Request contract and one Agent Resource definition.
- Model capability record covering evaluation, classification eligibility,
  residency, latency, availability, and cost for each eligible model.
- Policy-driven router across deterministic logic, one locally hosted model or
  small language model, and one external provider.
- Per-task recording of the model class that proved sufficient, as required by
  the ADR-0024 sovereignty tiers.
- One complete Forge delivery path executed against a locally hosted model, with
  measured quality, latency, and cost compared to the external path.
- Structured-output validation, bounded retry, health-aware fallback,
  cancellation, and cost accounting.
- Evaluation with versioned datasets, quality thresholds, and a promotion
  lifecycle.
- Per-request and aggregate token, time, and financial budgets.

Deferred until a delivery need exists: general tool registry and broker, MCP
edge adapters, indexed retrieval, and multi-provider routing beyond the two
paths above.

### Dependencies

- M1 authority, context, worker, evidence, and approval contracts.
- M2 supply-chain controls for models, adapters, and evaluation datasets.
- Accepted intelligence, provider-data-handling, and budget decisions.
- A locally hosted model environment the accountable owner and a community
  contributor can both run.

### Acceptance criteria

- Replacing the external model provider requires adapter and configuration
  changes only, not domain changes.
- One supported task runs through deterministic, locally hosted, and external
  paths with comparable evidence.
- Restricted content is blocked from an ineligible provider before disclosure.
- Invalid structured output cannot reach a side-effecting operation.
- Budget exhaustion stops execution or uses only a pre-authorized fallback.
- Model, prompt, context, policy, and evaluation versions are reproducible from
  the execution record.
- Forge remains administrable, inspectable, and recoverable without external
  model connectivity.
- Sovereignty tier three is satisfied by measured evidence before any
  community-maintainability claim is published.

### Estimated complexity

**Large.** Reduced from Extra Large by ADR-0024, but the locally hosted path,
quality variance, and cost governance remain genuinely difficult.

### Risks

- Local model quality proving insufficient for real delivery, invalidating the
  community-maintainability claim.
- Lowest-common-denominator abstractions hiding useful provider capabilities.
- Prompt injection crossing from retrieved data into actions.
- Cost and fallback loops causing denial of wallet.
- Evaluation datasets leaking repository information.

---

## Vision direction (M4–M9)

[ADR-0024](decisions/0024-rebaseline-forge-as-bounded-means.md) demoted the
following to evidence label **Vision**. They record intended direction. They are
**not** a committed plan, **not** a dependency for current work, and **not** a
justification for building capability ahead of demonstrated need.

The previous roadmap presented these as sequenced milestones with acceptance
criteria. That presentation implied commitments the available capacity could not
honour. The detail is preserved in Git history and may be recovered if a
delivery need reinstates any of it.

| Direction | Intent | Reinstated when |
|---|---|---|
| Software System Twin and advisory mode | Ingest, visualize, and explain an existing system with source, freshness, and confidence | A delivery need requires explaining a system Forge did not build |
| Collaborative multi-role agent workforce | Specialist agents with typed hand-offs, challenge reviews, and autonomy per risk class | Single-path delivery measurably limits throughput, per the AI delivery experiment measures |
| Application Environment Factory | Generated application and development environments from accepted intent | A second application is commissioned |
| Production platform core | One production-shaped, customer-controlled self-hosted profile with SLOs, recovery, and support | A customer contracts for supported operation |
| Enterprise deployment breadth | SaaS, dedicated, restricted, and disconnected profile conformance | Supported production operation exists and demand justifies breadth |
| Generality and bounded autonomy | A second unrelated application, conformance suites, and expanded autonomy on measured evidence | Workplace Orchestrator is delivered and its measures support expansion |

Reinstating any direction as committed scope requires an accepted decision, not
a roadmap edit.

## Cross-milestone success measures

Forge measures outcomes rather than generated volume:

- lead time from intent to independently verified preview or pull request;
- percentage of material changes with complete intent-to-evidence traceability;
- human correction and rejection rate for generated proposals;
- accepted journey coverage, freshness, repeatability, and regression escape rate;
- security coverage, false-positive rate, validated remediation time, and residual-risk disclosure;
- provider replacement effort and deployment-profile conformance;
- rollback, backup, restore, and incident-recovery success;
- model, tool, compute, financial, energy, and human-attention cost per accepted outcome;
- accessibility, comprehension, trust, and end-user outcome measures;
- SLO achievement and time to safe recovery.

Lines of generated code, number of agents, model size, and number of deployment services are not success measures.

## Roadmap governance

- Each milestone requires an accountable owner, accepted scope, risk review, measurable exit evidence, and explicit human completion decision.
- Milestone completion changes relevant capabilities from Research to Trial, Supported, or Reference only when their individual conformance evidence permits it.
- New trust boundaries, deployment units, public contracts, production profiles, or authority levels require ADR review.
- Failed assumptions trigger scope reduction, additional evidence, or architecture revision; they do not justify silent control bypass.
- The roadmap SHALL be reviewed after every milestone, material incident, failed assurance exercise, provider portability failure, or major product-boundary change.
