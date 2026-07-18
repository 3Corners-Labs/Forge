# Forge

Forge is an enterprise engineering foundation for governed AI capabilities.

The repository currently defines the architectural and governance baseline. It does not claim that the described platform capabilities have already been implemented or validated.

## Status

| Area | Status |
|---|---|
| Engineering principles | Proposed baseline |
| Reference architecture | Proposed |
| Capability contracts | To be defined |
| Security architecture | Proposed baseline |
| Production readiness | Not assessed |

Status terms used throughout the documentation:

- **Proposed**: documented for review; not yet approved.
- **Accepted**: approved through the architecture decision process.
- **Validated**: supported by repeatable evidence.
- **Implemented**: available in a delivered system; this does not imply production readiness.

## Start here

1. [Engineering Constitution](docs/standards/engineering-constitution.md)
2. [System Overview](docs/SYSTEM_OVERVIEW.md)
3. [Architecture](docs/ARCHITECTURE.md)
4. [Capability Catalog](docs/CAPABILITY_CATALOG.md)
5. [Security Architecture](docs/SECURITY.md)
6. [Engineering Workflow](docs/TASK_WORKFLOW.md)
7. [Architecture decisions](docs/adr/index.md)
8. [Requests for comments](docs/rfc/index.md)

## Scope

Forge separates reusable platform capabilities from product-specific behavior. The baseline covers governance, identity, policy, context, execution, integration, evidence, and operations.

Product features, provider-specific behavior, and deployment topology are outside the normative baseline until accepted through an ADR.

## Decision rule

When documents conflict, the order of authority is:

1. accepted Engineering Constitution;
2. accepted ADRs;
3. accepted Architecture;
4. capability contracts;
5. operational guidance.

Draft or proposed material never overrides an accepted decision.

## Contribution

Changes begin with an explicit task or RFC. Architecture, security, public contracts, and deployment-boundary changes require architecture review before implementation. Human approval remains mandatory for consequential decisions.
