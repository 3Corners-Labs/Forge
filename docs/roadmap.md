# Forge implementation roadmap

Status: Accepted
Owner: Forge Product and Architecture
Version: 0.4

## Purpose

This roadmap evolves Forge from the minimum platform capable of preparing a trustworthy change to itself into a complete Enterprise AI Engineering Platform.

It describes target delivery sequencing, not implemented state or calendar commitments. Capability maturity remains separate from document status. A milestone becomes complete only through the stated acceptance evidence and accountable human review.

## Delivery rules

- Milestones deliver end-to-end outcomes, not disconnected service skeletons.
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

| Milestone | Outcome | Capability posture at exit |
|---|---|---|
| M1 | Forge prepares one trustworthy change to Forge | Minimum self-building kernel |
| M2 | Forge builds and verifies reproducibly | Engineering and supply-chain foundation |
| M3 | Forge uses models and tools through governed contracts | Governed Intelligence and Tool Platform |
| M4 | Forge understands and explains existing systems | Software System Twin and advisory mode |
| M5 | Specialist agents collaborate under explicit authority | Forge development-team application |
| M6 | Forge builds its first business application | Application Environment Factory and Workplace lighthouse |
| M7 | One production-shaped customer-controlled profile is Supported | Production platform core |
| M8 | SaaS, dedicated, restricted, and disconnected profiles conform | Enterprise deployment and governance platform |
| M9 | Forge proves generality and earns bounded autonomy | Published enterprise capability and support envelope |

## Sequencing decision

The Proposed masterplan originally places broad read-only observation and advisory capability before Forge prepares changes to itself. This roadmap deliberately starts with the narrower accepted intent-to-pull-request kernel required by the user and ADR-0006.

M1 includes only the repository inspection, bounded context, planning, and explanation needed to prepare one governed Forge change. M4 remains responsible for the general Software System Twin and multi-audience advisory product. [ADR-0007](decisions/0007-sequence-self-building-before-full-advisory.md) records this distinction so the roadmap does not silently contradict the masterplan.

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

## M2 — Reproducible engineering and supply-chain foundation

### Objective

Make every Forge build, verification result, development environment, and promoted artefact within the declared Trial envelope reproducible, attributable, and recoverable before broadening intelligent behaviour.

### Deliverables

- Versioned developer-environment blueprint and pinned toolchain.
- Enforced source-module dependency rules derived from the accepted metamodel.
- Isolated build pipeline with locked dependencies and deny-by-default credential access.
- OCI-compatible artefacts identified by immutable digest.
- SBOM, build provenance, vulnerability and license assessment, and signature or attestation.
- Build-once promotion workflow with revocation and rollback metadata.
- Ephemeral preview environment with synthetic identities, deterministic fixtures, reset, and cleanup.
- Separate telemetry and evidence pipelines with controlled promotion from telemetry to evidence.
- Secret, source, dependency, filesystem, image, configuration, and IaC checks on every applicable build.
- Backup, restore, migration, rollback, and interrupted-build rehearsal for the Developer Trial.
- Expansion of the minimum F0 architecture, documentation, dependency, and isolation fitness checks into full supply-chain and release assurance.

### Dependencies

- M1 candidate digest, verification, approval, publication, and evidence semantics.
- Accepted software supply-chain decision covering builder identity, signing trust, artefact formats, and revocation.
- Approved secret and signing-key handling for CI and local development.

### Acceptance criteria

- Two clean builds from the same reviewed source and declared inputs produce the expected reproducibility result; any permitted variance is documented and signed.
- Changing one artefact byte causes digest or signature verification failure.
- A mutable tag cannot authorize promotion.
- Untrusted pull-request code cannot access signing material or production credentials.
- The same verified digest moves through compatible test and preview environments without rebuilding.
- Rollback deploys a previously verified artefact and compatible data state.
- Backup and restore reproduce declared state without resurrecting logically deleted data silently.
- CI publishes complete evidence for every required gate and rejects stale or missing provenance.

### Estimated complexity

**Large.** Most work is deterministic, but supply-chain identity, signing, reproducibility, recovery, and CI isolation introduce high assurance requirements.

### Risks

- Chasing perfect bit-for-bit reproducibility before defining permitted platform variance.
- CI supply-chain actions becoming an unreviewed execution dependency.
- Signing keys being exposed to untrusted builds.
- Security scanning volume producing unactionable noise.
- Preview and recovery automation becoming more complex than the product kernel.

---

## M3 — Governed Intelligence and Tool Platform

### Objective

Allow Forge to use deterministic methods, local models, external models, and tools through stable capability contracts while policy, state, budgets, and verification remain deterministic.

### Deliverables

- Versioned Task Request and Agent Resource contracts.
- Model registry recording capability, version, evaluation, classification eligibility, residency, retention, latency, availability, energy, and cost.
- Policy-driven router across deterministic logic, one local model, and one external provider.
- Structured-output validation, bounded retry, health-aware fallback, cancellation, and cost accounting.
- Evaluation service with versioned datasets, quality thresholds, regressions, and promotion lifecycle.
- Context assembly supporting exact retrieval first and optional indexed retrieval with authorization before ranking.
- Tool registry and broker with schemas, risk class, side effects, dry-run, idempotency, health, and evidence.
- One native tool adapter and one MCP-compatible edge adapter behind the same Forge authorization contract.
- Short-lived secret-reference resolution inside the worker boundary.
- Per-request, project, tenant, and organization token, compute, time, concurrency, and financial budgets.
- Provider and tool usage explanations generated from structured routing and evidence records.

### Dependencies

- M1 authority, context, worker, evidence, and approval contracts.
- M2 supply-chain controls for models, tools, adapters, evaluation datasets, and executable dependencies.
- Accepted intelligence, tool, provider-data-handling, budget, and secret-management ADRs.
- Approved data-classification and provider-eligibility policy.

### Acceptance criteria

- Replacing the external model provider requires adapter and configuration changes only, not domain changes.
- One supported task runs through deterministic, local-model, and external-model paths with comparable evidence.
- Restricted test content is blocked from an ineligible provider before disclosure.
- Invalid structured output cannot reach a side-effecting tool.
- Budget exhaustion stops execution or uses only a pre-authorized fallback.
- A tool or MCP server cannot widen its actions, destination, credentials, or evidence access.
- Model, prompt, context, tool, policy, and evaluation versions are reproducible from the execution record.
- Forge remains administrable, inspectable, and recoverable without model connectivity.

### Estimated complexity

**Extra Large.** The work combines probabilistic evaluation, provider portability, privacy, cost governance, tool side effects, and failure handling.

### Risks

- Lowest-common-denominator abstractions hiding useful provider capabilities.
- Model quality variance making deterministic acceptance difficult.
- Prompt injection crossing from retrieved data into tool actions.
- Cost and fallback loops causing denial of wallet.
- MCP or gateway products being mistaken for authorization boundaries.
- Evaluation datasets leaking customer or repository information.

---

## M4 — Software System Twin and advisory Forge

### Objective

Enable Forge to ingest, understand, visualize, and explain an existing software system before it gains broader write authority.

### Deliverables

- Authorized ingestion of repository artefacts, accepted documents, decisions, contracts, tests, dependencies, releases, and selected runtime evidence.
- Federated knowledge relationship model that preserves capability-owned canonical records.
- Rebuildable Software System Twin projections with source, version, freshness, confidence, conflicts, and authorization scope.
- Code, contract, architecture, user-journey, data-flow, dependency, deployment, and evidence views.
- Impact-analysis and contradiction-detection services.
- Read-only product, architecture, security, engineering, and operations advisory roles.
- Minimal Workbench supporting audience-specific explanations and drill-down to evidence.
- Governed memory promotion, retention, deletion, revocation, and staleness handling.
- Import and export of portable project knowledge without a hidden dependency on Forge SaaS.

### Dependencies

- M3 governed intelligence and context assembly.
- M2 provenance and artefact identity.
- M1 scope, identity, policy, and evidence controls.
- Accepted knowledge, memory, projection, and Workbench contracts.

### Acceptance criteria

- The Twin can be deleted and rebuilt from canonical artefacts and evidence without losing source truth.
- Every displayed relationship exposes source, version, freshness, confidence, and access scope.
- A removed or superseded source is excluded, marked stale, or shown with its lifecycle status as required.
- Cross-project and cross-tenant retrieval tests return no unauthorized existence or content signals.
- Advisory agents identify known seeded contradictions and clearly distinguish evidence from inference.
- Different audience modes explain the same system state consistently without changing canonical data.
- Forge analyzes its own repository in read-only mode and produces a reviewable change-impact plan.

### Estimated complexity

**Extra Large.** The central challenge is trustworthy projection and explanation across heterogeneous sources without creating a second source of truth.

### Risks

- The knowledge graph becoming a central transactional database.
- Stale projections being presented as current architecture.
- Over-collection of source, prompts, personal data, or telemetry.
- Model-generated relationships being mistaken for accepted facts.
- Visualization quality masking missing evidence.

---

## M5 — Collaborative Forge development team

### Objective

Turn Forge from a single governed change path into a coordinated application-development team whose specialist agents can discover, design, implement, assure, and prepare reversible changes under explicit authority.

### Deliverables

- Agent registry with immutable identity, owner, purpose, accepted tasks, tools, context policy, budgets, evaluations, and maximum autonomy.
- Product, architecture, UX, engineering, quality, end-user assurance, security, delivery, operations, and support role resources.
- Durable delivery graph with typed hand-offs, challenge reviews, escalation, cancellation, compensation, and manual recovery.
- Autonomy levels 0–3 applied per action and risk class rather than globally.
- Workbench views for backlog, plan, dependencies, agent authority, progress, evidence, approvals, and recovery.
- Accepted and exploratory journey lifecycle with synthetic identities, public-interface execution, accessibility evidence, and visible quarantine.
- Independent verification and separation-of-duty enforcement across agent roles.
- Reproducible development and preview environment generation for Forge changes.
- Passive application scanning and bounded authenticated security journeys for previews.
- Safe self-development controls keeping the candidate version distinct from trusted policy, verification, and promotion control.

### Dependencies

- M4 System Twin and explanations.
- M3 agent, intelligence, context, model, tool, evaluation, and budget contracts.
- M2 previews and supply-chain foundation.
- M1 durable workflow, approval, evidence, and source-control publication.
- Accepted agentic-delivery, journey-assurance, experience, and passive-security Trial specifications.
- Accepted M1 journey contracts from F0, expanded here into multi-role, interrupted, accessibility, exploratory, and quarantine behaviour.

### Acceptance criteria

- Forge prepares a non-trivial Forge change using at least product, architecture, engineering, and independent quality roles.
- Every hand-off preserves intent, scope, authority, context provenance, unresolved risks, and evidence.
- No agent approves its own architecture, suppresses its own finding, changes an accepted expected outcome, or merges its own candidate.
- An interrupted multi-step delivery resumes or recovers without duplicate side effects.
- One accepted allowed journey, one denied journey, one interrupted journey, and one accessibility journey pass through the public interface.
- Exploratory scenarios remain non-blocking until accountable review.
- The candidate Forge version cannot weaken the external branch, policy, verification, or approval controls evaluating it.
- Human reviewers can understand what changed, why, how it was verified, residual risk, cost, and rollback.

### Estimated complexity

**Programme.** This milestone introduces multi-role product behaviour, human experience, durable coordination, probabilistic collaboration, and circular-trust resistance.

### Risks

- Agent collaboration multiplying plausible errors and cost rather than increasing quality.
- Role names becoming simulated personalities without real capability or authority boundaries.
- Human reviewers being overwhelmed by generated artefacts and approvals.
- Candidate self-development influencing its own verifier or evidence.
- Journey automation changing expected outcomes to match implementation.
- Accessibility and UX being deferred behind backend orchestration work.

---

## M6 — Application Environment Factory and Workplace lighthouse

### Objective

Prove that Forge can build and operate a real business application without embedding that application’s domain architecture into Forge core.

### Deliverables

- Stable, licensed, version-controlled import of the Workplace Orchestrator source baseline with hashes and provenance.
- Product workspace containing accepted intent, personas, journeys, domain decisions, privacy requirements, and product-owned architecture.
- Application and development-environment blueprints generated through Forge capability contracts.
- Workplace Software System Twin with role-specific product, architecture, UX, security, engineering, and operations views.
- Thin Workplace vertical slice covering tenant identity, one resource, availability or hold, accessible interaction, audit, and one provider-neutral adapter.
- Independent tests, threat evidence, journey assurance, SBOM, provenance, signatures, and deployment manifests.
- SaaS-like and connected self-hosted reference deployments using the same signed application artefact where profile compatibility permits.
- Demonstrated rollback, backup, restore, provider outage, key-service outage, and runtime explanation.

### Dependencies

- M5 usable Forge development-team application.
- Confirmed Workplace source ownership, license, authorship, and immutable import location.
- M2 application supply chain and environment blueprints.
- M3 provider-neutral model, tool, and adapter contracts where intelligence is used.
- Accepted Workplace product baseline and lighthouse-specific ADRs.

### Acceptance criteria

- Every Workplace lighthouse exit criterion passes.
- Original product intent traces to accepted requirements, decisions, code, tests, evidence, and deployment.
- Workplace domain logic has no direct dependency on a named model, KMS, identity broker, event broker, or collaboration vendor.
- Tenant context cannot be caller-overridden and negative isolation tests pass.
- The core journey meets its declared accessibility target without depending exclusively on 3D interaction.
- SaaS-like and self-hosted profiles run the same signed application artefact or document an accepted profile-specific deployment-unit difference.
- Forge-specific modules contain no workplace reservation, occupancy, navigation, sensor, or collaboration-domain logic.
- End users, architects, engineers, security reviewers, and operators receive evidence-backed explanations.

### Estimated complexity

**Programme.** This is the first full product commission and tests Forge across product, architecture, engineering, security, UX, deployment, and operations.

### Risks

- Workplace requirements becoming hidden Forge platform requirements.
- Source provenance or licensing remaining unresolved.
- Lighthouse scope overwhelming the still-maturing Forge platform.
- Employee, location, occupancy, or presence data creating elevated privacy risk.
- A successful single product creating false confidence in platform generality.

---

## M7 — Production platform core

### Objective

Make Forge supportable for one production-shaped, customer-controlled connected self-hosted profile before attempting full SaaS and sovereign breadth.

### Deliverables

- Accepted connected self-hosted profile manifest with exact connectivity, tenancy, operations, assurance, and provider-ownership dimensions.
- Production identity federation, workload identity, delegation revocation, operator access, and break-glass procedures.
- Tenant-aware transactional storage, object storage, queues or outbox, caches, indexes, evidence, telemetry, secrets, and cryptographic domains.
- Versioned Crypto Service Port with rotation, rewrap, revocation, recovery, provider health, and at least two conforming provider adapters before portability is claimed.
- Supported installation, configuration, migration, upgrade, rollback, backup, restore, and uninstall contracts.
- SLOs, capacity envelopes, admission control, quotas, noisy-neighbour controls, RTO/RPO, alerts, runbooks, and incident response.
- Signed release bundles, trust-root updates, vulnerability response, support bundles, and revocation.
- Active DAST and vulnerability-regression Trial against isolated, resettable non-production targets under approved rules of engagement.
- Customer export and clean-environment import with semantic comparison.

### Dependencies

- M6 application and Forge operational evidence.
- M2 supply-chain and recovery foundation.
- Accepted production deployment, identity, cryptography, data, retention, support, and security-testing ADRs.
- Accepted F0 offensive-assurance safety contract, expanded here into active DAST and vulnerability-regression operation.
- Named accountable production, security, privacy, and operations owners separated where policy requires.

### Acceptance criteria

- The profile’s published capacity, latency, availability, RTO/RPO, upgrade, rollback, and support scenarios pass under the declared validation envelope.
- Tenant, project, operator, cache, queue, index, telemetry, evidence, secret, and key negative-isolation suites pass.
- Key rotation, provider outage, rewrap, restore, wrong-tenant denial, revocation, and destructive-action dual control pass.
- A clean operator installs, upgrades, backs up, restores, rolls back, exports, and imports without undocumented Forge access.
- Production mode rejects developer credentials, debug exporters, unsafe port bindings, mutable release identities, and unverified extensions.
- Security-test scope, redirect and DNS-rebinding denial, rate limits, kill switch, cleanup, and finding validation pass.
- Residual risks, unsupported combinations, and operational responsibilities are published.

### Estimated complexity

**Programme.** Production support combines platform engineering, security, cryptography, data lifecycle, reliability, documentation, and customer operations.

### Risks

- Supporting production operations before the team can sustain incident and upgrade obligations.
- Key-management portability claims exceeding tested provider semantics.
- Self-hosted variability making support and conformance uneconomical.
- Active security testing escaping its authorized target or causing disruption.
- Unquantified scale leading to premature service extraction.

---

## M8 — Enterprise deployment and governance platform

### Objective

Extend the production core into Forge SaaS, dedicated SaaS, restricted self-hosted, and disconnected self-hosted profiles with shared contracts and explicit enterprise governance.

### Deliverables

- Forge SaaS pooled or partitioned multi-tenant control plane with regional data controls and governed operator access.
- Dedicated SaaS topology for customers requiring stronger isolation.
- Restricted and disconnected bundles with local identity, models, tools, telemetry, updates, revocation snapshots, and evidence export.
- Fleet management for version, policy, health, capacity, vulnerability, upgrade, and support state.
- Organization policy packs, entitlements, budgets, quotas, evidence packs, retention, legal hold, and customer-controlled provider choices.
- Extension lifecycle covering publisher identity, signatures, permissions, isolation, provenance, compatibility, revocation, rollback, and data cleanup.
- Portable project, policy, workflow, evidence, configuration, migration, and supported artefact formats.
- Compliance and assurance profiles that compose controls without forking core semantics.
- Customer-managed and customer-controlled key options with tested shared-responsibility models.
- Product, legal, community, enterprise, SaaS, and support model recorded through accepted decisions.

### Dependencies

- M7 Supported connected self-hosted profile and operating model.
- M6 application evidence proving profile contracts against a real product.
- Accepted commercial, tenancy, residency, operator-access, offline-update, extension, and profile-conformance decisions.
- Operational organization capable of SaaS security, privacy, support, incident, and fleet responsibilities.

### Acceptance criteria

- SaaS tenant isolation, operator access, regional routing, deletion, support access, abuse protection, backup, and restore tests pass.
- A supported project exports from SaaS and imports into a clean supported self-hosted environment with semantic equivalence.
- A disconnected installation verifies, installs, operates, updates, rolls back, and exports evidence without runtime Forge SaaS access.
- A revoked extension, publisher, model, tool, signing identity, or offline bundle is rejected according to policy.
- Supported profiles publish conformance matrices, SLOs, capacity limits, residency, key ownership, recovery, upgrade, and support boundaries.
- Foundational identity, isolation, cryptography, audit, accessibility, portability, backup, and safe-upgrade qualities are not made optional commercial add-ons.
- Profile-specific deployment units do not fork application domain source or policy semantics.

### Estimated complexity

**Programme.** This milestone introduces multi-tenant SaaS operations, offline distribution, fleet management, enterprise governance, and commercial product boundaries.

### Risks

- Profile combinations exceeding the team’s operational capacity.
- SaaS convenience creating hidden dependencies for self-hosted customers.
- Commercial packaging weakening foundational safety or portability.
- Offline updates and revocation becoming stale or unverifiable.
- Extensions becoming an authority-escalation or supply-chain path.
- Jurisdiction-specific legal obligations being mistaken for one universal compliance profile.

---

## M9 — Generality, continuous improvement, and bounded autonomy

### Objective

Prove that Forge is a general Enterprise AI Engineering Platform rather than a Forge- or Workplace-specific system, then increase autonomy only where measured evidence supports it.

### Deliverables

- A second unrelated application developed through the same platform contracts and deployment profiles.
- Comparative portability, capability-fit, cost, quality, and coupling analysis across Forge, Workplace, and the second application.
- Governed runtime and user-feedback ingestion producing product hypotheses and improvement proposals rather than automatic truth.
- Bounded dependency maintenance, vulnerability remediation, test repair, cost optimization, incident diagnosis, and recovery assistance.
- Autonomy level 4 for specifically accepted operational actions with independent policy, evidence, rollback, and emergency stop.
- Human-led penetration-test and red-team intake with validated remediation and regression scenarios.
- Autonomous security-testing Trial only after external scope enforcement, kill switch, target denial, isolation, evidence integrity, and independent finding validation pass.
- Capability and profile conformance suites suitable for external implementation comparison.
- Standards, certification, or accreditation proposals kept organizationally separate from product self-attestation.
- Continuous measurement of user outcomes, escaped defects, recovery, review corrections, security coverage, provider replacement effort, cost, and sustainability.

### Dependencies

- M8 supported enterprise profiles and governance.
- A genuinely unrelated second application with independent product ownership.
- Longitudinal operational, quality, security, cost, and user-outcome evidence.
- Accepted autonomy, learning, maintenance, incident, offensive-assurance, and conformance decisions.
- Independent accountable owners and qualified human assurance functions.

### Acceptance criteria

- The second application uses existing Forge contracts without introducing hidden Workplace-specific semantics into platform core.
- Provider and deployment substitutions remain bounded to adapters and profile configuration.
- Product or runtime learning cannot silently change accepted intent, policy, memory, scenarios, or release gates.
- Every autonomous operational action has exact scope, risk class, evidence, rollback or compensation, and emergency termination.
- The candidate system cannot approve, certify, or provide the sole evidence for its own trustworthiness.
- Measured autonomy reduces lead time or operational toil without degrading security, quality, review correction rate, cost, accessibility, or recovery outcomes.
- Conformance reports disclose unsupported capabilities, failed checks, evidence freshness, confidence, and residual risk.
- Forge meets the published Enterprise AI Engineering Platform capability and support envelope across at least two unrelated applications and the supported deployment profiles.

### Estimated complexity

**Programme.** Completion depends on multiple products, sustained production evidence, independent assurance, organizational maturity, and controlled expansion of authority.

### Risks

- Generality being claimed from superficial differences between applications.
- Feedback loops poisoning knowledge, memory, evaluation, or accepted intent.
- Automation optimizing visible metrics while degrading user outcomes or safety.
- Autonomous maintenance or security testing exceeding authority.
- Certification becoming circular self-attestation or a distraction from product quality.
- Increased autonomy outpacing incident response and human oversight capacity.

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
