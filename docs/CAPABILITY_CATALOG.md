---
title: Capability Catalog
version: 1.0
status: Proposed
owner: Architecture
---

# Capability Catalog

The catalog describes platform responsibilities, not products or AI task examples.

| Capability | Responsibility | State | Depends on |
|---|---|---|---|
| Experience contracts | Stable external requests and responses | Proposed | Identity context, execution |
| Identity context | Authenticate principals and delegation | Proposed | Enterprise identity source |
| Policy decision | Decide authorization and conditions | Proposed | Identity, classification, policy data |
| Data classification | Determine handling and disclosure constraints | Proposed | Organizational policy |
| Context assembly | Build minimal, attributable context packages | Proposed | Policy, knowledge sources |
| Execution orchestration | Coordinate bounded workflows | Proposed | Context, registry, policy |
| Model and tool registry | Describe eligible capabilities and constraints | Proposed | Configuration, policy |
| Tool authorization | Authorize exact external side effects | Proposed | Identity, policy, resource metadata |
| Human approval | Manage accountable approval and expiry | Proposed | Identity, policy, evidence |
| Validation and evaluation | Validate inputs, outputs, and quality evidence | Proposed | Execution, policy |
| Durable knowledge | Store reusable knowledge with provenance | Proposed | Classification, access control |
| Execution state | Store transient checkpoints and cancellation state | Proposed | Execution |
| Audit evidence | Preserve tamper-resistant decision history | Proposed | All consequential capabilities |
| Operational telemetry | Provide logs, metrics, traces, and alerts | Proposed | All runtime capabilities |
| Cost and quota control | Bound tokens, calls, time, and spend | Proposed | Registry, execution, policy |
| Secrets and key management | Issue and rotate scoped credentials | Proposed | Workload identity |

## Contract requirements

Before a capability can be marked accepted, it must define:

- accountable owner;
- inputs, outputs, and error semantics;
- identity and authorization requirements;
- data classification and retention;
- state ownership and consistency;
- telemetry and audit events;
- service objectives and failure behavior;
- versioning and compatibility rules;
- deterministic acceptance tests.
