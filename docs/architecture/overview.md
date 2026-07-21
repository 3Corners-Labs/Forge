# Architecture overview

Status: Proposed
Owner: Forge Architecture

## Scope

This document describes the target logical architecture. The current repository only contains an early CLI and local development bootstrap; target capabilities SHALL NOT be represented as implemented until verified in code and tests. The current evidence boundary is maintained in the [implementation baseline](../governance/implementation-baseline.md).

## Architectural style

Forge is a modular, policy-driven control plane with independently executable workers. The initial implementation SHOULD be a modular monolith plus isolated worker processes. Deployment units may be extracted when required by scale, isolation, ownership, or availability.

The [architecture metamodel](metamodel.md) defines the relationship between planes, capabilities, modules, deployment units, providers, profiles, and artefacts. The diagrams below are a logical concern view, not a service topology.

~~~mermaid
flowchart LR
    Human[Human collaborators] --> Experience[Experience plane]
    System[External systems] --> Gateway[API and event gateway]
    Experience --> Control[Forge control plane]
    Gateway --> Control
    Control --> Knowledge[Knowledge and context]
    Control --> Workforce[Agentic workforce]
    Control --> Trust[Trust plane]
    Workforce --> Execution[Isolated execution plane]
    Execution --> Providers[Model, tool and platform adapters]
    Control --> Operations[Operations plane]
    Execution --> Trust
    Execution --> Operations
    Trust --> Experience
    Operations --> Experience
~~~

## Planes

### Experience plane

Provides conversation, workspaces, backlogs, visualizations, reviews, approvals, and role-specific explanations. It is not a privileged bypass around platform policy.

### Control plane

Owns intent, projects, delivery runs, orchestration, state transitions, budgets, approvals, and capability selection. It contains deterministic workflow logic and never delegates authorization to a model.

### Knowledge plane

Manages versioned project knowledge, retrieval, provenance, data classification, context assembly, and retention.

### Agentic workforce

Provides role-based agents for product, architecture, UX, engineering, end-user assurance, security, delivery, operations, and support. Roles are capability bundles governed by policy, not simulated personalities with implicit authority.

### Execution plane

Runs builds, tests, previews, migrations, browser sessions, and tool calls inside isolated, time-bounded environments with controlled network and secret access.

### Trust plane

Provides distinct capabilities for identity, policy decision, authorization enforcement, risk, verification, user journey assurance, continuous security testing, approval, release eligibility, security, and audit evidence. These capabilities may initially share deployment infrastructure but retain separate contracts and separation-of-duty rules.

### Operations plane

Operates releases and running products, and captures telemetry, cost, incidents, recovery, and operational outcomes. Telemetry becomes audit evidence only through an explicit integrity-protected promotion step.

## Dependency rules

- Domain capabilities depend on contracts, not provider SDKs.
- User-facing clients depend on public APIs, never internal persistence.
- Model and tool adapters do not make policy decisions.
- Workers receive a bounded execution envelope and cannot expand their own authority.
- Evidence is emitted through a stable contract and SHALL NOT become a synchronous dependency for every read path.
- Cross-capability data access occurs through explicit contracts or owned events.

## Deployment invariants

Self-hosted and SaaS deployments share capability contracts, artefact formats, policy semantics, and migration paths. They may differ in tenancy topology, managed operations, identity federation, and available provider adapters.
