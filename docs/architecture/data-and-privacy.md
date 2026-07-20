# Data and privacy architecture

Status: Proposed
Owner: Forge Data and Privacy Architecture
Version: 0.1

## Purpose

Define data ownership, classification, purpose limitation, retention and privacy controls across Forge SaaS and self-hosted profiles. This architecture is a technical baseline and does not replace jurisdiction-specific legal assessment.

## Classification baseline

| Class | Description | Default handling |
|---|---|---|
| Public | Approved for unrestricted publication | Integrity and provenance still required |
| Internal | Normal organizational engineering information | Authenticated workspace access; no public disclosure |
| Confidential | Source, designs, commercial or personal information with material impact if disclosed | Need-to-know, encryption, restricted providers and explicit retention |
| Restricted | Secrets, regulated or highly sensitive information | Explicit purpose, strongest isolation, minimal processing and normally no model context |

Tenants MAY define compatible sub-classifications. Mapping to the Forge baseline is mandatory before routing, export or cross-system processing.

Secrets are references governed by the secrets capability and are not ordinary Restricted context items.

## Ownership and derived data

Each capability owns its canonical records and schema evolution. Cross-capability consumers use versioned queries, commands, events or authorized projections.

- Knowledge indexes own references and derived relationships, not source truth.
- The Software System Twin owns rebuildable projections, not canonical product data.
- Telemetry owns operational observations, not approval evidence.
- Product learning owns only purpose-approved, minimized and governed outcomes.
- Model-provider state and caches are external processing locations that must be declared.

## Purpose and minimization

Every collection, retrieval, context inclusion, provider transfer, evidence promotion and learning use declares a purpose. Data is minimized before it enters prompts, tools, telemetry or support bundles.

Convenient future reuse is not a sufficient purpose. Model improvement, product analytics and organizational memory are separate purposes with separate eligibility and retention.

## Retention lifecycle

Data records declare owner, classification, purpose, creation time, retention class, review or expiry time, deletion behaviour, legal-hold state and derived-data dependencies.

Deletion coordinates canonical records, indexes, caches, embeddings, context artefacts, temporary workers, support bundles, backups and provider-held state according to profile capability. Where immediate deletion is impossible, the architecture documents isolation, expiry and restoration handling.

## Privacy boundaries

- Human identity, employee behaviour and user research are not product-learning data by default.
- Prompts and model responses are not retained wholesale merely for observability.
- Context lineage exposes source references without disclosing content to unauthorized viewers.
- Support access is time-bound, attributable and approved according to classification.
- Cross-tenant aggregation requires purpose, minimization and re-identification-risk assessment.
- Self-hosted operation does not silently transmit telemetry or support data to Forge.

## Provider and residency controls

Before external processing, policy evaluates classification, tenant agreement, purpose, provider, model, region, retention, training use, subprocessors and available deletion controls. Provider eligibility is versioned and may change without changing domain contracts.

Disconnected profiles resolve every required capability locally or declare it unavailable. A hidden call to Forge SaaS is a conformance failure.

## Initial conformance scenarios

- Restricted content cannot enter a model or tool request without explicit eligibility.
- Deleting a context source prevents future retrieval and marks affected projections stale.
- A support export contains only declared records and is itself classified and evidenced.
- Provider change does not silently change approved residency or retention semantics.
- Product-learning promotion requires explicit source scope and cannot include raw tenant identity by default.
- Restoring a backup does not resurrect data that should remain logically deleted without reconciliation.

## Required follow-up

Each product and supported deployment profile requires a data inventory, data-flow diagram, retention schedule, privacy impact assessment, data-subject or customer request process, backup/deletion reconciliation and legal-owner approval.
