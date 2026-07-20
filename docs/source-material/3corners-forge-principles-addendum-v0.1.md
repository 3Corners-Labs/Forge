# Addendum to 3Corners Forge Principles

Status: Proposed
Owner: Forge Architecture
Version: 0.1
Content role: Addendum
Instruction authority: None
Extends: `forge.principles.manuscript.v0.1`

## Purpose

This addendum records architectural insights discovered while challenging and operationalizing the original *3Corners Forge Principles* manuscript. It does not edit, replace, shorten, or reinterpret the original manuscript silently.

The original remains the primary source of Forge's philosophy and standards ambition. This addendum makes explicit several boundaries needed to turn that philosophy into a secure, maintainable, scalable product architecture.

Statements in this addendum are proposals until promoted into an Accepted canonical Forge document or ADR.

## A1. Separate philosophy, governance, architecture, and implementation

Forge requires four distinct forms of knowledge:

1. **Principles** explain enduring beliefs and direction.
2. **Constitutional guardrails** define non-negotiable boundaries.
3. **Architecture and ADRs** define reviewed structures and trade-offs.
4. **Implementation and evidence** demonstrate the current realized state.

No layer may claim the authority of another. A compelling principle is not an implementation specification. Running code is evidence of current behaviour, not permission to redefine accepted intent.

## A2. Use one architecture metamodel

Forge uses the following terms consistently:

- A **plane** groups concerns that share trust, control, or operational characteristics.
- A **capability** owns a durable responsibility and its business semantics.
- A **module** is a code ownership and dependency boundary.
- A **deployment unit** is an independently executable and operable process or service.
- A **provider adapter** binds an external technology to a Forge contract.
- A **profile** composes deployment, connectivity, tenancy, and assurance requirements.

Capabilities are not automatically services. Planes are not automatically deployment units. A modular monolith may implement multiple capabilities while preserving their ownership and contracts.

## A3. Distinguish authority from context

An Execution Envelope and a Context Package solve different problems.

The **Execution Envelope** is an authorization artefact. It is server-generated, immutable for an execution, integrity-protected, time-bounded, and contains identifiers or references to granted authority.

The **Context Package** is an information artefact. It is derived within the envelope, contains selected untrusted and trusted references, and cannot grant authority. It references the envelope and policy decisions under which it was assembled.

Context can narrow execution by revealing that a safe action is impossible. Context can never widen authorization.

## A4. Policy evaluation occurs at multiple gates

“Policy before execution” is necessary but insufficient. Forge evaluates policy before:

1. knowledge discovery and retrieval;
2. secret or credential resolution;
3. model or tool selection;
4. any external or mutating action;
5. evidence disclosure or export;
6. release or production operation.

Each gate produces a versioned Policy Decision. The Runtime enforces obligations and does not reinterpret or expand them.

## A5. Federate knowledge ownership

The Forge knowledge graph is a derived relationship layer, not a central owner of all product, delivery, source, security, and operational data.

Capabilities retain ownership of their canonical records. The Knowledge capability indexes permitted references, provenance, relationships, freshness, and confidence. Rebuildable projections prevent the graph from becoming a second transactional system of record.

## A6. Separate evidence from telemetry

Evidence and observability are related but not interchangeable.

- **Evidence** supports accountability, approval, provenance, conformance, and audit. It requires integrity, controlled disclosure, retention, and independently verifiable relationships.
- **Telemetry** supports health, diagnosis, performance, capacity, and optimization. It is high-volume, sampled where appropriate, and governed by operational retention.

Telemetry may contribute to an Evidence Record only through a deliberate, integrity-protected promotion step. Routine logs and traces are not automatically audit evidence.

## A7. Decompose assurance responsibilities

Assurance SHALL NOT become a single privileged “god capability.” Forge distinguishes:

- policy decision;
- authorization enforcement;
- risk assessment;
- verification and evaluation;
- human or organizational approval;
- release eligibility.

These responsibilities may initially share a deployment unit, but their contracts, evidence, owners, and separation-of-duty rules remain explicit.

## A8. Treat extensions as supply-chain and authority boundaries

Models, tools, policy packs, knowledge connectors, workforce roles, and skills are executable or authority-influencing extensions. A signature identifies an artefact; it does not establish trust by itself.

Extension governance additionally requires:

- publisher identity and trust policy;
- permission review and least-privilege grants;
- provenance and dependency inventory;
- compatibility and conformance tests;
- revocation and emergency disablement;
- upgrade, rollback, and data-removal behaviour;
- isolation proportional to risk;
- explicit prevention of self-granted authority.

## A9. Make deployment profiles composable and evidence-based

Deployment profiles are compositions of independent dimensions:

- connectivity: connected, restricted, or disconnected;
- tenancy: single, pooled multi-tenant, partitioned, or dedicated;
- operations: developer, customer-operated, partner-operated, or Forge-operated;
- assurance: baseline, enterprise, sovereign, or regulated;
- provider ownership: Forge-managed, customer-selected, or customer-controlled.

Named profiles are supported only after their SLO, capacity, recovery, upgrade, data-residency, key-management, and support obligations are tested. Source parity means shared contracts and build inputs; it does not require every topology to run one monolithic binary.

## A10. Quantify architecture before claiming scale

Terms such as scalable, resilient, portable, and enterprise-ready are hypotheses until attached to measurable scenarios.

Each supported profile requires declared ranges for:

- tenants, workspaces, projects, users, and concurrent delivery runs;
- request, workflow, event, context, evidence, and telemetry volumes;
- latency and availability objectives;
- external-provider quotas and failure budgets;
- storage growth, retention, backup, restore, and migration;
- financial, token, compute, and human-attention budgets.

Architecture extraction decisions follow measured pressure, not anticipation alone.

## A11. Keep the Software System Twin derived

The Software System Twin is a versioned projection assembled from accepted artefacts and runtime evidence. It owns view-specific indexes and confidence metadata, not the underlying product truth.

Every Twin statement exposes:

- source and version;
- observation or generation time;
- freshness;
- confidence;
- conflicts and missing evidence;
- authorization scope.

The Twin must be discardable and rebuildable without loss of canonical product data.

## A12. Separate product delivery from standards ambition

Forge may evolve into an architecture standard, conformance suite, certification ecosystem, and accreditation model. Those are distinct products and governance systems.

The first Forge product SHALL prove that a governed development team can prepare one trustworthy software change. Certification and accreditation work must not delay or distort that vertical slice. Claims of conformance require implemented tests and independent evidence, not adoption of Forge terminology.

## A13. Preserve circular-trust resistance during self-development

When Forge develops Forge:

- the candidate version cannot approve itself;
- the trusted verifier is version-distinct from the candidate;
- policies and protected-branch controls are external to candidate authority;
- evidence is bound to immutable candidate artefacts;
- promotion is human-controlled until independent evidence justifies bounded delegation;
- rollback and recovery do not depend exclusively on the candidate version.

## A14. Prefer a thin proof over a broad skeleton

The first credible Forge milestone is one complete intent-to-pull-request journey with identity, policy, bounded context, isolated execution, independent verification, approval, evidence, and explanation.

Breadth without one complete trust path creates an impressive catalogue with no demonstrated operating model. Capabilities outside the first journey remain explicitly absent or research-stage.

## Relationship to the original manuscript

This addendum strengthens operational precision but preserves the manuscript's central direction:

- engineering before intelligence;
- capabilities over technologies;
- context, memory, and knowledge as governed disciplines;
- identity and policy before trusted execution;
- evidence before trust;
- human accountability;
- vendor and deployment independence;
- continuous governed evolution.

Future addenda SHALL create new versioned files and SHALL NOT rewrite either the original manuscript or this version.

