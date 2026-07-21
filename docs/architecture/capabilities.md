# Capability model

Status: Accepted
Owner: Forge Architecture
Version: 0.3

Capabilities describe durable responsibilities. They are neither deployment units nor automatically separate services. The [architecture metamodel](metamodel.md) maps them to logical planes and initial modules.

## Capability catalogue

| Stable ID | Domain | Capability | Responsibility | Initial maturity |
|---|---|---|---|---|
| `forge.capability.tenant` | Organization | Tenant and workspace | Ownership, isolation, membership and quotas | Research |
| `forge.capability.entitlement` | Commercial | Plans and entitlements | Packaging, usage, billing evidence and customer entitlements | Research |
| `forge.capability.identity` | Identity | Identity and delegation | Human, workload and agent identity; scoped delegation | Research |
| `forge.capability.intent` | Product | Intent and portfolio | Outcomes, users, hypotheses and priorities | Research |
| `forge.capability.delivery` | Delivery | Work orchestration | Plans, durable state, dependencies, cancellation and recovery | Research |
| `forge.capability.knowledge` | Knowledge | Knowledge relationships | Provenance-rich, authorized and rebuildable relationships over owned artefacts | Research |
| `forge.capability.context` | Knowledge | Context assembly | Purpose-bound selection and reproducible Context Packages | Research |
| `forge.capability.memory` | Knowledge | Memory lifecycle | Scoped promotion, retrieval, retention and deletion | Research |
| `forge.capability.architecture` | Architecture | Decision management | Constraints, ADRs, fitness functions and dependency rules | Research |
| `forge.capability.experience` | Experience | Workbench and visualization | Conversational, graphical and expert views | Research |
| `forge.capability.twin` | Experience | Software System Twin | Rebuildable system projections with freshness and confidence | Research |
| `forge.capability.workforce` | Workforce | Agent registry and coordination | Roles, skills, hand-offs, supervision and evaluation | Research |
| `forge.capability.intelligence` | Intelligence | Model and method routing | Eligibility, quality, privacy, latency and cost selection | Research |
| `forge.capability.tool` | Tools | Tool registry and brokerage | Discovery, authorization, validation and invocation | Research |
| `forge.capability.execution` | Execution | Isolated work environments | Sandboxes, builds, tests, previews and controlled egress | Research |
| `forge.capability.environment` | Environments | Application environment factory | Reproducible development, test, preview and production blueprints | Research |
| `forge.capability.policy` | Trust | Policy decision | Deterministic permit, deny and obligation decisions | Research |
| `forge.capability.verification` | Trust | Risk and verification | Risk assessment, quality gates and independent evaluation | Research |
| `forge.capability.journey-assurance` | Trust | User journey assurance | Role-based intent scenarios, public-interface execution, expected outcomes and journey evidence | Research |
| `forge.capability.security-testing` | Trust | Continuous security testing | Risk-based source, supply-chain, DAST, vulnerability-regression and adversarial assurance | Research |
| `forge.capability.approval` | Trust | Approval and release eligibility | Human or organizational authorization bound to immutable artefacts | Research |
| `forge.capability.evidence` | Trust | Audit evidence | Integrity-protected decisions, lineage, attestations and disclosures | Research |
| `forge.capability.telemetry` | Operations | Observability | Metrics, logs, traces, profiles and operational cost | Research |
| `forge.capability.release` | Operations | Release operations | Promotion, deployment, rollback and runtime configuration | Research |
| `forge.capability.reliability` | Operations | Reliability and incident response | SLOs, capacity, continuity, recovery and incident lifecycle | Research |
| `forge.capability.feedback` | Customer success | Feedback and support | Product feedback, cases, adoption and service communication | Research |
| `forge.capability.extension` | Platform | Configuration and extensions | Configuration, schemas, compatibility, installation and migrations | Research |
| `forge.capability.supplychain` | Platform | Software supply chain | Dependencies, provenance, SBOM, signing and immutable promotion | Research |

Commercial, customer-success, certification and accreditation concerns are outside the first implementation slice even though they remain part of the long-term product landscape.

## Capability contract

Before functional implementation, every participating capability SHALL define:

- stable ID and accountable owner;
- purpose and explicit non-responsibilities;
- commands, queries, events, outputs and normalized errors;
- canonical data ownership and classification;
- consistency, transaction and idempotency expectations;
- identity, authorization, privacy and audit requirements;
- availability, latency, throughput, capacity and recovery hypotheses;
- versioning, compatibility and migration policy;
- maturity and support envelope;
- conformance tests and required evidence.

The first-slice boundary is defined in [Trustworthy intent-to-PR vertical slice](vertical-slice.md), with its detailed logical contracts in [First-slice capability contracts](first-slice-capability-contracts.md). Other capabilities remain Research and SHALL NOT be represented as implemented.

User Journey Assurance is specified separately from generic Verification because it owns accepted role-based product scenarios and expected outcomes, not test-runner infrastructure or release approval. Continuous Security Testing is separate because offensive tooling requires stricter scope, safety and finding-governance contracts. Both may initially run on shared execution and verification infrastructure.

## Execution Envelope

Every agent or tool execution receives one integrity-protected Execution Envelope containing authority, not secret values or retrieved content:

- envelope ID, schema version, issuer, issue time and expiry;
- tenant, workspace, project and delivery-run identifiers;
- initiating principal and delegated agent identity;
- declared purpose, action set and risk class;
- referenced Policy Decisions and Approval Records;
- approved classifications and external destinations;
- eligible model, tool and secret-reference identifiers;
- filesystem, network, token, compute, cost, time and concurrency limits;
- idempotency key and evidence correlation identifier.

The envelope is server-derived and immutable for one execution. Delegation may create a child envelope with strictly narrower authority and explicit parent lineage. A Context Package references an envelope but cannot grant or expand authority.
