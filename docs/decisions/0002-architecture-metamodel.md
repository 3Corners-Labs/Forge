# ADR-0002: Adopt one architecture metamodel

Status: Accepted
Owner: Forge Architecture
Date: 2026-07-20

## Context

Forge documentation describes planes, platforms, capabilities, modules, workers, services, providers and profiles. Without a common metamodel these views can become competing decompositions and encourage premature service boundaries.

## Decision

Adopt the concepts and mapping rules in `docs/architecture/metamodel.md`. Capabilities own semantics; modules own code boundaries; deployment units own runtime operation; planes group concerns; profiles compose deployment requirements.

Control, Trust, Knowledge, Execution, Experience and Operations are logical planes. Agentic Workforce and Intelligence span planes through governed capabilities and receive no implicit authority.

## Alternatives

- Treat each capability as a microservice: rejected because scale and ownership boundaries are unproven.
- Use only deployment diagrams: rejected because logical ownership would become topology-dependent.
- Preserve multiple undocumented views: rejected because contradictions would remain invisible.

## Consequences

Documentation and future code dependency rules need stable capability IDs. One modular process may host several capabilities, but direct cross-capability persistence remains prohibited.

## Validation

- Every first-slice component maps to a plane, capability, module and deployment unit.
- Architecture fitness tests detect duplicate capability IDs and forbidden module dependencies.
- A topology change does not require changing capability semantics.

## Supersession

None.
