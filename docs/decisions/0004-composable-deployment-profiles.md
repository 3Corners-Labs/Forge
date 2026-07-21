# ADR-0004: Model deployment profiles as compositions

Status: Accepted
Owner: Forge Platform and Security Architecture
Date: 2026-07-20

## Context

The masterplan and deployment document used different profile lists. Terms such as self-hosted or SaaS do not by themselves define connectivity, tenancy, operator, assurance or provider ownership.

## Decision

Define profiles as compositions of independent connectivity, tenancy, operations, assurance and provider-ownership dimensions. Maintain a small named catalogue for product communication. A profile becomes Supported only after its measurable conformance envelope passes.

Developer is the only current Trial profile. No production profile is Supported.

## Alternatives

- One profile per possible combination: rejected because combinations grow exponentially.
- Only SaaS and self-hosted: rejected because both hide material trust and operational differences.
- One mandatory Kubernetes topology: rejected because compact and disconnected deployments may require different supported topologies.

## Consequences

Profiles share contracts, source, build definitions and conformance semantics but may use different signed deployment-unit sets. Product documentation must state exact supported dimensions.

## Validation

- A profile manifest resolves every required dimension.
- Export/import and contract tests run without Forge SaaS.
- Each Supported profile declares SLO, capacity, RTO/RPO, upgrade, rollback, residency, key and support obligations.

## Supersession

None.
