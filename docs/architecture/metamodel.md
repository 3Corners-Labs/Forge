# Forge architecture metamodel

Status: Proposed
Owner: Forge Architecture
Version: 0.1

## Purpose

This metamodel defines how Forge's principles, planes, capabilities, modules, deployment units, providers, profiles, and artefacts relate. It prevents conceptual views from becoming competing architectures.

## Concepts

| Concept | Definition | Must not be assumed to be |
|---|---|---|
| Principle | Durable belief guiding trade-offs | A testable implementation contract |
| Guardrail | Non-negotiable boundary in an Accepted constitution | A product feature |
| Plane | Logical grouping of concerns with related trust or operational characteristics | A deployable service |
| Capability | Durable responsibility with owned semantics and contracts | A process, team, database, or product |
| Module | Enforced source-code ownership and dependency boundary | Automatically a capability or service |
| Deployment unit | Independently executable and operable software unit | Automatically independently owned data |
| Provider adapter | Replaceable binding between a Forge contract and external technology | An authorization boundary |
| Profile | Composition of connectivity, tenancy, operations, assurance and provider requirements | A fork of application source |
| Canonical artefact | Versioned record owned by one capability | A generated index, dashboard, or model context |
| Projection | Rebuildable view over canonical artefacts and evidence | A second source of truth |

## Logical planes

Forge uses six logical planes:

1. **Experience** — interaction, review, explanation and visualization.
2. **Control** — intent, delivery state, orchestration, budgets and approvals.
3. **Knowledge** — authorized discovery, context assembly, memory and derived relationships.
4. **Trust** — identity, policy, risk, verification, evidence and security controls.
5. **Execution** — isolated work, tools, models, builds, tests and environment operations.
6. **Operations** — deployment, runtime health, telemetry, recovery, support and learning.

The Agentic Workforce and Intelligence Platform span planes through capabilities. They are not privileged planes and do not own authorization.

## Capability-to-plane mapping

| Capability family | Primary plane | Critical dependencies |
|---|---|---|
| Tenant and workspace | Control | Identity, Policy, Evidence |
| Intent and portfolio | Control | Knowledge, Experience |
| Delivery orchestration | Control | Policy, Execution, Evidence |
| Identity and delegation | Trust | Tenant, Policy |
| Policy decision | Trust | Identity, Classification, Configuration |
| Risk and verification | Trust | Evidence, Delivery, Security |
| Approval and release eligibility | Trust | Policy, Verification, Evidence |
| Audit evidence | Trust | All side-effecting capabilities |
| Knowledge relationships | Knowledge | Capability-owned canonical artefacts |
| Context assembly | Knowledge | Identity, Policy, Knowledge, Memory |
| Memory lifecycle | Knowledge | Policy, Classification, Evidence |
| Agent registry and coordination | Control | Identity, Policy, Intelligence, Tools |
| Model routing | Execution | Policy, Evaluation, Budget |
| Tool brokerage | Execution | Identity, Policy, Sandbox, Evidence |
| Isolated work environments | Execution | Policy, Secrets, Operations |
| Application environment factory | Execution | Delivery, Supply chain, Operations |
| Workbench and System Twin views | Experience | Capability APIs and projections |
| Telemetry and reliability | Operations | Deployment units and providers |
| Release and runtime operations | Operations | Supply chain, Policy, Evidence |

## Mapping rules

- A capability has one accountable owner even when implemented in multiple modules.
- A module may implement multiple small capabilities only when dependency rules preserve their contracts.
- A deployment unit may host multiple modules; this does not authorize shared-table coupling.
- A capability becomes a separate deployment unit only through an Accepted ADR with measured drivers.
- Cross-plane interaction uses explicit commands, queries, events, or evidence references.
- Trust-plane decisions are enforced by the receiving capability; trust is not inferred from network location.
- Projections expose source, version, freshness, confidence and authorization scope.

## Initial deployment mapping

The first slice uses a modular control-plane process plus isolated worker processes:

| Deployment unit | Initial modules | Isolation reason |
|---|---|---|
| Forge control plane | workspace, delivery, identity adapter, policy decision, context manifest, evidence metadata | Transactional coordination and rapid domain learning |
| Forge worker | repository adapter, tool broker, build/test runner | Untrusted source and executable workload isolation |
| Evidence sink | append-oriented evidence records | Integrity and write-path separation |
| Provider adapters | model, source control and optional external tools | Quota, failure and data-egress boundary |

This mapping is a Trial hypothesis. It does not authorize a microservice decomposition or shared unrestricted persistence.
