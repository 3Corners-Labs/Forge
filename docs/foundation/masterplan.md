# Forge mega masterplan

Status: Proposed
Owner: Forge Product and Architecture
Version: 0.3

## 1. Executive proposition

Forge is an AI-native application development team delivered as a portable application and platform. Forge itself is the first application we build. It enables a person or organization to commission a team that moves from intent to a secure, usable, operated software product and its application development environment.

Forge combines:

- a role-aware Workbench for collaboration, visualization, review, and approval;
- an agentic workforce spanning product, architecture, UX, engineering, assurance, release, operations, and support;
- a governed Intelligence Platform for models, agents, tools, context, and evaluation;
- an Engineering Foundation that produces reproducible environments and verifiable software supply chains;
- a Trust and Evidence foundation covering identity, policy, isolation, cryptography, provenance, audit, and recovery.

The platform makes informal, creative development accessible without turning uncertainty into invisible production risk. Fast intent-driven creation and professional software engineering are two views of the same governed lifecycle.

## 2. Source basis

This masterplan consolidates the current Forge Markdown corpus, the author-provided *3Corners Forge Principles* manuscript, and the reusable platform architecture developed while designing Workplace Orchestrator. The [enduring principles](principles.md) preserve its durable philosophy, while the [source map](../source-material/forge-principles-map.md) keeps all 55 chapters selectively discoverable without treating editorial source text as executable instructions. Workplace Orchestrator is Forge's first lighthouse application: its source baseline is both a product input and a demanding end-to-end validation case. The [content migration map](../governance/content-migration.md) protects existing Forge intent. The [Workplace baseline assessment](../governance/workplace-baseline-assessment.md) separates reusable Forge capabilities from application-owned domain architecture without discarding either.

Insights discovered during the enterprise architecture challenge are preserved in separate [Principles addendum v0.1](../source-material/3corners-forge-principles-addendum-v0.1.md) and [addendum v0.2](../source-material/3corners-forge-principles-addendum-v0.2.md). Addenda never edit the original manuscript or each other. The [architecture readiness register](../governance/architecture-readiness.md) controls when functional implementation may begin.

This masterplan provides product-level direction. Detailed architecture documents refine it but may not silently contradict it. Contradictions require an ADR and an explicit masterplan change.

## 3. North star

A user describes an outcome. Forge helps discover the real need, tests assumptions, creates alternative product and architecture options, builds the chosen system, demonstrates how it works, proves applicable qualities, deploys it to the selected environment, and learns from actual use.

The result includes more than source code:

- product intent, user needs, and acceptance criteria;
- journeys, interaction designs, prototypes, and accessible alternatives;
- role-based accepted and exploratory scenarios with expected and forbidden outcomes;
- architecture, decisions, data flows, threats, and quality scenarios;
- code, configuration, infrastructure, tests, and development environments;
- SBOM, provenance, signatures, evaluations, and approval evidence;
- deployment, observability, recovery, support, and improvement plans.

## 4. Product boundaries

### Forge is

- an application and application-environment factory;
- a coordination system for humans and specialist agents;
- a policy-controlled execution and evidence platform;
- a visual and conversational learning environment;
- portable across Forge SaaS and self-hosted deployment;
- able to develop Forge itself under stricter self-modification controls.

### Forge first, Workplace second

[Workplace Orchestrator](../reference-applications/workplace-orchestrator.md) is the first business application commissioned to the Forge team. Its existing masterplan, privacy and security requirements, logical architecture, cryptographic design, technology research, and product roadmap become structured input to Forge rather than being merged into Forge platform code.

Before accepting that commission, Forge must exist as a usable development-team application: it can ingest knowledge, coordinate roles, prepare branches and pull requests, run independent checks, explain its work, and operate within human-controlled authority. Workplace then proves that the team can ingest substantial prior thinking, challenge it, create a governed product workspace, produce a Software System Twin, build a reproducible application environment, and deliver a secure vertical slice for both SaaS and self-hosted operation.

## 4A. Bootstrap strategy

Forge has a bootstrap paradox: the future development team cannot be required to create its own first working version. We resolve this through explicit stages and unchanged human authority.

### Stage 0: Human-built foundation

Humans use conventional engineering tools and bounded coding assistants to establish the constitution, product model, architecture, repository controls, reproducible environment, identity, policy, evidence, and a minimal Workbench. Generated work is treated as normal untrusted contribution.

### Stage 1: Forge observes and advises

Forge imports its own repository, builds a Software System Twin, explains architecture and code, identifies contradictions, proposes plans, and performs read-only reviews. It has no write or execution authority.

### Stage 2: Forge prepares changes to Forge

Task-scoped agents can create isolated workspaces, branches, code, tests, documentation, and pull requests. Independent checks and human review remain mandatory. Forge cannot approve its own architecture, weaken its own controls, merge protected branches, change production policy, or obtain production credentials.

### Stage 3: Forge operates as the development team

Product, architecture, UX, engineering, assurance, release, operations, and support roles coordinate through the Workbench. Bounded reversible actions can be automated according to evidence and risk. Humans retain outcome, policy, budget, and production accountability.

### Stage 4: Forge builds Workplace Orchestrator

The team accepts the Workplace product baseline as its first business commission and delivers the governed journey defined in the reference-application dossier.

### Stage 5: Prove generality

Forge develops a second unrelated application. This exposes accidental coupling to Workplace and validates the platform contracts.

At every stage, the system version under review is distinct from the trusted version performing checks. Promotion requires evidence external to the candidate change, preventing circular self-attestation.

### Forge is not

- an unbounded autonomous loop with production authority;
- a single model, prompt library, IDE plug-in, or workflow framework;
- a replacement for accountable product and risk ownership;
- a claim that generated code is correct because it compiles;
- a mandatory microservice or Kubernetes topology for every installation;
- a hidden dependency on Forge SaaS for self-hosted customers.

## 5. Primary users and experiences

| User | Desired outcome | Forge experience |
|---|---|---|
| Creator | Turn an idea into a working application | Guided intent, prototype, explanation, preview, and safe publication |
| Product professional | Validate value and coordinate scope | Outcomes, hypotheses, research, journeys, backlog, and delivery confidence |
| Architect | Control boundaries and qualities | Capability map, decisions, scenarios, threats, dependencies, and fitness results |
| UX professional | Create inclusive effective experiences | Research evidence, flows, prototypes, design system, accessibility, and analytics |
| Engineer | Understand and evolve software safely | Reproducible environment, code impact, tests, traces, reviews, and rollback |
| Security or risk owner | Constrain and evidence risk | Data flows, policy, threat model, findings, approvals, and audit |
| Operator | Run reliable products | Releases, SLOs, telemetry, incidents, cost, capacity, and recovery |
| End user | Understand and trust product behaviour | Plain-language behaviour, data use, limitations, support, and feedback |
| Platform administrator | Govern Forge itself | Tenancy, identity, policies, models, tools, budgets, providers, and fleet health |

## 6. Product structure

Forge separates product concerns from reusable platform concerns.

The [architecture metamodel](../architecture/metamodel.md) is authoritative for the distinction between logical planes, capabilities, code modules, deployment units, providers and profiles once Accepted.

~~~text
Forge Workbench and Public Interfaces
                 |
       Delivery Control Plane
        /        |         \
Knowledge   Agentic Workforce   Assurance
        \        |         /
    Intelligence Platform
                 |
    Isolated Execution Platform
                 |
Source, Build, Cloud, Runtime, Tool and Model Adapters
                 |
      Evidence and Operations
~~~

### 6.1 Forge Workbench

The Workbench supports conversation, visual exploration, role-specific views, editable artefacts, previews, reviews, and approvals. It makes simple work approachable while permitting expert inspection.

### 6.2 Delivery Control Plane

The control plane owns projects, intent, work graphs, lifecycle state, budgets, approvals, and orchestration. It uses durable, deterministic state transitions and cannot delegate authorization to a model.

### 6.3 Knowledge Platform

The Knowledge Platform manages versioned artefacts, a project knowledge graph, context packages, provenance, classification, retrieval, memory scopes, and retention. Repository files remain portable first-class sources.

### 6.4 Agentic Workforce

The workforce provides governed roles and challenge relationships across the lifecycle. Agents create reviewable outputs through task-scoped authority. They never gain authority merely because they produced a convincing explanation.

### 6.5 Intelligence Platform

Applications request capabilities rather than named models. Policy selects deterministic logic, optimization, classical ML, local SLM, self-hosted model, private endpoint, or external model based on eligibility and evidence.

### 6.6 Execution Platform

Builds, tests, previews, browser sessions, code execution, migrations, and operational tools run in isolated environments with bounded identity, network, filesystem, secret, compute, time, and cost.

### 6.7 Assurance, Evidence, and Operations

Assurance enforces policy, safety, quality gates, user journey assurance, bounded security testing, approval, and release eligibility. Evidence captures traceability, evaluation, provenance, telemetry, audit, cost, and actual outcomes under privacy and retention controls.

## 7. The Forge Software System Twin

Forge maintains a living, evidence-backed model of every product it creates. This Software System Twin links intended behaviour to the implemented and running system.

It includes:

- outcomes, requirements, journeys, and acceptance evidence;
- system context, capabilities, modules, APIs, events, and dependencies;
- data ownership, classification, lineage, residency, and retention;
- source artefacts, tests, supply-chain evidence, and releases;
- deployment topology, configuration relationships, runtime traces, and SLOs;
- known risks, incidents, costs, and improvement opportunities.

The twin is not a manually maintained diagram or a new source of truth. It is a versioned projection generated from accepted artefacts and runtime evidence. Every view exposes its source, version, freshness, and confidence.

The Workbench uses the twin to explain the system differently to an end user, UX designer, architect, engineer, security reviewer, or operator.

## 8. Agentic workforce operating model

### 8.1 Role families

- Discovery and product agents clarify users, value, constraints, and success.
- Architecture agents propose boundaries and deliberately challenge quality and complexity.
- UX agents research, prototype, test accessibility, and explain interaction decisions.
- Engineering agents prepare small implementation changes and development environments.
- Test and evaluation agents construct independent evidence and failure analysis.
- End-user assurance agents execute accepted and exploratory multi-role scenarios through public interfaces without redefining expected outcomes.
- Security, privacy, and compliance agents identify risk and propose controls.
- Autonomous security-testing agents, when explicitly enabled, operate as externally constrained offensive components under rules of engagement.
- Delivery and release agents assemble evidence and coordinate reversible promotion.
- Operations and support agents observe outcomes and prepare bounded remediation.

### 8.2 Forge Agent Resource

Every executable agent is a versioned resource containing:

- immutable identity, owner, purpose, and accepted task types;
- model-routing and fallback policy;
- allowed tools, repositories, environments, and scopes;
- context, data-classification, retention, and network policy;
- token, compute, time, concurrency, and financial budgets;
- required output schema, validations, approvals, and evaluation suite;
- maximum autonomy and lifecycle status.

Credentials are short-lived and task-scoped. Repository writes are branch-scoped. Engineering agents do not receive production credentials.

### 8.3 Separation of duties

No agent accepts its own consequential architecture, suppresses its own security finding, weakens its own quality gate, or merges its own protected change. Independent deterministic validation or accountable human review is required according to risk.

### 8.4 Autonomy

Autonomy is granted per action and risk class: observe, advise, prepare, execute reversible actions, or perform bounded operations. Destructive data actions, production release, policy change, and authority expansion require explicit safeguards and normally human approval.

## 9. Intent-to-operation lifecycle

~~~text
Intent
  -> Discovery and evidence
  -> Product outcomes and acceptance criteria
  -> Architecture and threat decisions
  -> UX journeys and prototypes
  -> Plan and isolated implementation
  -> Independent verification and evaluation
  -> Multi-role user journey and adversarial assurance
  -> Security, privacy and operational readiness
  -> Approval and immutable artifact promotion
  -> Deployment, observation and support
  -> Measured learning back to intent
~~~

The lifecycle is a graph: work can branch, iterate, pause for evidence, or return to an earlier decision. Forge supports methodology profiles by composing controls such as work-in-progress limits, review policy, release cadence, stage gates, SRE practices, and compliance evidence.

## 10. Stable intelligence contracts

### 10.1 Task request

An intelligence request identifies capability, purpose, tenant and project context, data classification, execution profile, required quality and latency, budget class, tools, response schema, and evidence correlation. Tenant context is derived from verified identity and cannot be overridden by the caller.

### 10.2 Agent definition

Agent definitions are versioned software artefacts with accepted tasks, tools, context policy, routing policy, limits, approvals, evaluations, and owner.

### 10.3 Tool contract

Tools expose versioned JSON-compatible schemas, risk class, required identity, side effects, dry-run behaviour, idempotency, policy hooks, audit events, health, and error semantics. MCP MAY be supported at the edge; it is not an internal authority boundary.

### 10.4 Evaluation before promotion

Models, prompts, agents, tools, and context strategies move through Research, Trial, and Accepted states using application-specific golden sets, regressions, security and privacy tests, tool-use tests, and cost/latency comparisons.

Accepted user outcomes are versioned independently of implementation and exercised through public interfaces. Exploratory agent findings become candidate scenarios or product hypotheses until reviewed. Automated offensive testing is selected by risk and constrained by external scope, impact, target, data, and kill-switch controls; it complements qualified human penetration testing.

## 11. Knowledge and context

Forge creates bounded Context Packages instead of dumping an entire repository into a model. Every item has source, version, owner, classification, permission, freshness, and purpose.

Exact and structured retrieval are preferred where they solve the task. Search, code indexes, vectors, and graphs are derived mechanisms. Authorization occurs before ranking or generation.

Durable memory writes are explicit governed actions. Execution, delivery-run, project, organization, and approved product-learning scopes have different ownership and retention.

## 12. Security, privacy, and cryptography

### 12.1 Trust posture

Users, agents, models, prompts, retrieved content, repositories, dependencies, tool output, web content, and inbound events are untrusted by default. Network location alone never creates trust.

### 12.2 Identity and tenant context

Human, workload, device, and agent identities remain distinct. Delegation narrows authority. Tenant and project context are immutable, server-derived request properties enforced in policy, data, events, cache, telemetry, knowledge, and cryptography.

### 12.3 Data and key separation

Data, wrapped data keys, key-encryption keys, provider metadata, secrets, backups, and audit evidence are separate assets with distinct ownership and recovery policy.

### 12.4 Cryptographic architecture

Forge uses layered protection: infrastructure encryption, platform encryption, transport security, and selective application-level encryption. Classified application data MAY use envelope encryption with tenant- and purpose-specific key domains and authenticated associated data binding tenant, object, field, and crypto version.

Domain code uses a versioned Crypto Service Port rather than provider SDKs. Provider adapters can support native cloud KMS, KMIP-compatible services, or KMS/HSM combinations. OpenBao is a Trial reference candidate, never an implicit domain dependency.

Crypto agility requires explicit provider alias, logical key reference, key version, algorithm policy, purpose, and envelope version. Rotation, rewrap, revocation, destruction, outage, backup, restore, lost-key, and provider-switch scenarios are designed and tested.

Destructive key actions require impact preview, dual control, delay, recovery validation, and tamper-evident evidence.

### 12.5 Secrets and supply chain

Secrets are short-lived, scoped, injected at execution time, and excluded from prompts and telemetry. Software artefacts are built once, identified by immutable digest, supplied with SBOM and provenance, signed, and promoted without rebuilding.

## 13. Data, API, event, and workflow architecture

- Each capability owns its canonical data and schema evolution.
- Cross-capability writes use commands or owned contracts, never shared-table shortcuts.
- Public APIs are versioned and schema-described.
- Asynchronous events use versioned envelopes with tenant, event ID, time, schema, correlation, causation, trace context, and idempotency semantics.
- Delivery is assumed at least once; consumers are idempotent.
- Transactional outbox is the initial pattern when state change and event publication form one business action.
- Durable workflows define cancellation, deadlines, retries, compensation, checkpoints, and manual recovery.
- A broker, cache, vector database, or service mesh is introduced only after measured need and an ADR.

## 14. Deployment and product profiles

The same signed software artefacts and contracts support:

| Profile | Connectivity | Operations | Typical use |
|---|---|---|---|
| Developer | Local with optional external services | Individual developer | Fast safe contribution |
| Connected self-hosted | Policy-controlled connectivity | Customer or partner | Customer-controlled production |
| Restricted self-hosted | Approved private services | Customer | Regulated environment |
| Disconnected self-hosted | No runtime egress | Customer | Sovereign or air-gapped environment |
| Forge SaaS | Connected | Forge | Managed multi-tenant service |
| Dedicated Forge SaaS | Connected with isolated resources | Forge | Stronger enterprise isolation |

OCI artefacts, portable configuration schemas, identity and storage contracts, telemetry semantics, migrations, tests, and security defaults are shared. Deployment profiles bind providers and topology; they do not fork application source.

Kubernetes is a strong production reference candidate, not a universal product assumption. The deployment contract SHALL also permit a compact supported profile if it can satisfy security, recovery, observability, and upgrade conformance.

## 15. Community and commercial direction

The licensing and edition model remains an explicit product and legal ADR. Any community distribution must remain independently useful and secure. Identity, tenant isolation, cryptography interfaces, auditability, accessibility, backup/restore, signed updates, provider interfaces, and safe upgrades are baseline qualities.

Commercial value SHOULD focus on managed SaaS, support, fleet operations, certified integrations, advanced governance, evidence packs, scale profiles, dedicated topology, long-term support, and managed lifecycle—not on withholding foundational safety.

## 16. Quality and operations

Every major capability requires:

1. problem statement and measurable quality scenarios;
2. logical design, owned data, and trust/data flows;
3. security, privacy, misuse, and accessibility analysis;
4. API, event, error, and migration contracts;
5. ADRs and credible alternatives;
6. PoC or test evidence;
7. SLO, capacity, cost, backup, recovery, and operational design;
8. implementation, rollout, rollback, and support plan.

User-facing capabilities additionally require accepted multi-role scenarios covering successful, denied, interrupted, accessibility, privacy, and recovery behaviour. Security-sensitive capabilities require layered per-build assurance, versioned abuse cases, and rules of engagement for active testing.

Forge correlates delivery intent, agent work, model and tool routing, artefacts, releases, runtime behaviour, cost, and user outcomes without indiscriminately storing prompts, code, personal data, or secrets.

## 17. Evidence and architecture governance

### 17.1 Lifecycle status

Document authority follows the repository documentation lifecycle:

Draft -> Proposed -> Accepted -> Superseded or Deprecated -> Archive

Technology and implementation maturity is separate:

Research -> Trial -> Supported -> Reference -> Deprecated -> Retired

An Accepted architecture document can describe a Research capability. A Reference implementation cannot redefine an Accepted contract.

### 17.2 Evidence labels

- Standard: supported by an applicable published standard.
- Proven: demonstrated in comparable production use.
- Validated: tested in the Forge reference implementation.
- Hypothesis: plausible but untested.
- Vision: intended future direction.

Document status, implementation maturity, and evidence are different: an accepted decision can still describe Research capabilities and contain explicit hypotheses that require later validation.

### 17.3 Technology selection

Technology follows capability need, quality scenarios, standards, candidate comparison, security and license review, vertical PoC, operational rehearsal, ADR, and only then supported baseline.

Weighted evaluation considers architectural fit, operational simplicity, security/privacy, open governance and licensing, skills and support, portability, measured performance, and developer experience.

## 18. Roadmap

> **Implementation sequencing:** accepted [ADR-0007](../decisions/0007-sequence-self-building-before-full-advisory.md) places the bounded self-building kernel before the general advisory product. The stages below retain the product-level capability progression; the canonical implementation sequence and F0 readiness gate are defined in the [implementation roadmap](../roadmap.md).

### M0: Coherent architecture baseline

- preserve the Accepted constitution and review, revise, or accept this masterplan;
- complete source-content migration without loss;
- establish capability owners, quality scenarios, threat model, and ADR register;
- automate documentation, contract, and architecture fitness checks.

### M1: Reproducible engineering foundation

- create development-environment blueprints and hermetic toolchains;
- implement CI assurance, dependency and secret scanning;
- establish scenario contracts, synthetic identities, preview reset, and safe passive application scanning;
- produce OCI artefacts, SBOM, provenance, signatures, and release channels;
- establish telemetry, local test data, backup, and recovery rehearsal.

The first functional proof within this foundation is the [trustworthy intent-to-PR vertical slice](../architecture/vertical-slice.md). Capabilities outside that path remain explicitly Research.

### M2: Forge observes and advises

- create an authenticated workspace and immutable execution envelope;
- capture intent and visualize an initial Software System Twin;
- assemble a provenance-rich Context Package;
- let planning and architecture agents prepare a reviewed delivery graph;
- explain the current Forge repository without write authority.

### M3: Forge prepares Forge changes

- let a coding and test agent create one branch and pull request in isolation;
- apply policy, independent verification, approval, and complete evidence.
- implement capability-based routing across deterministic logic, one external model, and one local model;
- demonstrate native and MCP tool mediation;
- demonstrate structured output, evaluation, provider failover, privacy filters, and budget enforcement;
- replace a provider without application-code changes.

### M4: Forge development team application

- coordinate product, architecture, UX, engineering, assurance, release, operations, and support roles;
- provide role-specific Workbench views, human approvals, and separation of duties;
- execute accepted user scenarios through public interfaces and propose new exploratory scenarios without changing intent autonomously;
- orchestrate risk-based source, supply-chain, DAST, vulnerability-regression, and adversarial testing with coverage disclosure;
- generate reproducible application development and preview environments;
- use Forge to prepare improvements to Forge without circular self-approval.

### M5: Build Workplace Orchestrator

- import the Workplace Orchestrator baseline as a versioned product workspace;
- generate its application and development/test environment from intent and accepted decisions;
- support product, architecture, UX, engineering, quality, and security collaboration;
- provide live previews, accessibility evidence, deployment, rollback, and runtime explanation;
- deliver one production-shaped Workplace Orchestrator vertical slice;

### M6: Self-hosted and SaaS parity

- deploy identical signed artefacts to SaaS and connected self-hosted profiles;
- validate restricted and disconnected bundles;
- verify upgrades, rollback, backup, restore, key recovery, tenant isolation, and provider portability;
- publish supported conformance and SLO profiles.

### M7: Prove generality and bounded autonomy

- validate a second unrelated application to prove Forge is not coupled to the workplace domain;
- connect product and runtime learning to governed improvement proposals;
- automate maintenance, dependency updates, incident diagnosis, and optimization within policy;
- validate separation of duties and safe self-development of Forge;
- increase autonomy only where outcome and safety evidence supports it.

## 19. Initial architecture decision backlog

| Class | Decision | Required evidence |
|---|---|---|
| Product | Open/community, enterprise, and SaaS model | Legal, ecosystem, economics, and safety analysis |
| Architecture | Modular control plane and worker isolation | Dependency and failure-isolation tests |
| Intelligence | Shared capability gateway and task contract | Two-application, multi-provider PoC |
| Agents | Workflow-first execution and Agent Resource | Bounded multi-step change with independent review |
| Tools | Native contract with MCP at edge | Native/MCP authorization and failure comparison |
| Knowledge | Files plus derived knowledge graph | Rebuild, provenance, authorization, and staleness PoC |
| Data | Transactional baseline and outbox evolution | Load, replay, migration, and recovery tests |
| Security | Execution identity and tenant isolation | Automated negative and adversarial tests |
| Assurance | User journey scenarios and expected outcomes | Two-role public-interface journey with reproducible evidence |
| Security testing | Layered automated and human assurance | Bounded preview scan, target denial, kill switch and finding validation |
| Crypto | Envelope encryption and Crypto Service Port | Rotation, outage, restore, and two-provider PoC |
| Deployment | OCI portability and supported orchestrators | SaaS/self-hosted conformance matrix |
| Delivery | Build-once digest promotion | Signed promotion and rollback rehearsal |
| Experience | Workbench and Software System Twin | Multi-role comprehension and usability evaluation |

## 20. Success measures

Forge will measure:

- lead time from intent to independently verified preview;
- percentage of material changes with complete intent-to-evidence traceability;
- escaped defects, security findings, accessibility failures, and recovery success;
- human review time and rate of rejected or corrected agent proposals;
- provider replacement effort and self-hosted/SaaS conformance;
- user comprehension of system behaviour and decision rationale;
- accepted intent coverage by role, journey freshness, regression escape rate, scenario flake rate, and exploratory finding yield;
- security coverage, false-positive rate, mean time to validated finding and verified remediation, and residual-risk disclosure;
- model, compute, human-attention, financial, and sustainability cost per accepted outcome;
- SLO achievement and mean time to safe recovery.

Velocity or generated lines of code are not primary success measures.

## 21. Principal risks

- Vision breadth can prevent delivery of a coherent first slice.
- An attractive Workbench can mask weak evidence or unsafe authority.
- Agent collaboration can multiply cost and plausible errors without independent checks.
- Platform abstractions can become lowest-common-denominator complexity.
- Supporting many deployment profiles can exceed a small team's operational capacity.
- Self-development can create circular trust and change-control failure.
- Knowledge retention can violate privacy or preserve poisoned context.
- Synthetic-user evidence can capture sensitive data or create false confidence when expected results are weak.
- Autonomous security testing can exceed scope, mutate data, disrupt targets, or legitimize incomplete scanner coverage.
- Commercial pressure can weaken community safety or portability.

Each roadmap milestone SHALL retire named risks with measured evidence. Capabilities that cannot yet meet this condition remain explicitly Proposed or Trial.
