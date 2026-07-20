# Deployment architecture

Status: Proposed
Owner: Forge Platform Engineering
Version: 0.2

## Profile dimensions

A deployment profile composes independent dimensions instead of creating a new product fork for every combination:

| Dimension | Values under consideration |
|---|---|
| Connectivity | Connected, restricted, disconnected |
| Tenancy | Single, pooled multi-tenant, partitioned multi-tenant, dedicated |
| Operations | Developer, customer, partner, Forge |
| Assurance | Baseline, enterprise, sovereign, regulated |
| Provider ownership | Forge-managed, customer-selected, customer-controlled |

## Named profile catalogue

No production profile is currently Supported.

| Profile | Composition | Maturity | Intended validation |
|---|---|---|---|
| Developer | Connected; single workspace; developer-operated; baseline | Trial | Safe local contribution and first vertical slice |
| Connected self-hosted | Connected; single or partitioned tenant; customer-operated | Research | Customer-controlled production-shaped reference |
| Restricted self-hosted | Restricted; partitioned or dedicated; customer-operated | Research | Regulated private-provider operation |
| Disconnected self-hosted | Disconnected; dedicated; customer-operated; sovereign | Research | Offline installation, update and evidence export |
| Forge SaaS | Connected; pooled or partitioned; Forge-operated | Research | Managed multi-tenant operation |
| Dedicated Forge SaaS | Connected or restricted; dedicated; Forge-operated | Research | Stronger enterprise isolation |

“Self-hosted” without its connectivity, tenancy, provider, assurance and operations dimensions is not a sufficient conformance claim.

## Shared logical contract

Profiles share:

- capability semantics and public contract schemas;
- policy and evidence semantics;
- source and reproducible build definitions;
- OCI-compatible artefact and manifest formats;
- configuration, migration and backup formats where portable;
- identity, storage, telemetry and provider ports;
- conformance tests and security defaults.

Profiles MAY use different signed deployment-unit sets where topology requires it. “Same artefact” means an immutable build-once artefact is promoted unchanged between environments that consume that artefact; it does not require every profile to use one binary or topology.

## Shared logical components

- experience clients and public API gateway;
- policy-driven control plane;
- durable workflow and project state;
- knowledge, artefact and context stores;
- isolated worker pools;
- model, tool, source-control and deployment adapters;
- audit evidence sink;
- telemetry and operational cost pipeline.

Audit evidence and operational telemetry have separate retention, integrity and access contracts.

## Portability and exit contract

Projects, policies, workflow definitions, architecture decisions, evidence manifests and supported artefact metadata SHALL use portable versioned formats. Import and export SHALL work without Forge SaaS.

Profile specifications additionally define treatment of identities, secret references, encrypted data, provider-bound keys, extensions, large artefacts, deleted records and unsupported capabilities. Portability claims require a tested export, clean-environment import and semantic comparison.

## Scalability model

The control plane partitions by immutable tenant and workspace identifiers. Workers scale independently by workload and risk class. Provider adapters implement quotas, deadlines, backoff, circuit breaking and bulkheads. Durable workflows support idempotency, cancellation, retry, checkpoint and recovery.

Before a profile becomes Supported it SHALL define:

- tenant, project, user and concurrent-run validation ranges;
- API, policy, workflow, context, evidence and telemetry volumes;
- queue fairness, admission control and noisy-neighbour limits;
- SLOs, RTO/RPO, backup and restore tests;
- upgrade, rollback and schema-compatibility windows;
- residency, retention, support and provider-failure behaviour;
- cost, compute and storage envelopes.

## Extraction criteria

A module may become a separate deployment unit only when an Accepted ADR provides measurable evidence for at least one of:

- materially different scaling characteristics;
- stronger isolation or regulatory boundary;
- independent availability or recovery objective;
- independently owned release cadence;
- technology requirement that cannot safely coexist.

The ADR also defines operational cost, contract, data ownership, migration, failure behaviour and rollback.
