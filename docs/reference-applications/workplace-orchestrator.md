# Workplace Orchestrator lighthouse application

Status: Accepted
Owner: Workplace Orchestrator Product

> **Revised by [ADR-0024](../decisions/0024-rebaseline-forge-as-bounded-means.md).** Workplace delivery begins after M1 exit rather than at the former M6, proceeds in a **separate repository**, and is subject to the sovereignty tiers recorded below. Market positioning and the load-bearing product assumption are recorded separately in [Workplace Orchestrator positioning](workplace-orchestrator-positioning.md).

## Relationship to Forge

Workplace Orchestrator is the first business application commissioned to the Forge development team after Forge reaches its controlled change-preparation stage. It is simultaneously:

- a real product with independent users, outcomes, privacy duties, domain logic, and roadmap;
- the first product proof that Forge can turn a substantial architecture baseline into working software;
- a stress test for agentic product, architecture, UX, engineering, security, delivery, and operations collaboration;
- a validation of identical signed artefacts across SaaS and self-hosted profiles.

Workplace does not become part of the Forge platform core. It consumes Forge capabilities through versioned contracts in the same way that future applications will.

## Repository boundary

Workplace Orchestrator SHALL be developed in a separate repository with its own lifecycle, license, and decision register.

Forge modules SHALL NOT contain workplace reservation, occupancy, navigation, sensor, or collaboration-domain logic. Workplace SHALL NOT depend on unreleased Forge runtime internals.

The separation is not only a capability boundary. It is the only honest test of whether Forge can deliver into a codebase it did not generate, which is the question the [AI delivery experiment](../governance/ai-delivery-experiment.md) measures.

## Sovereignty tiers

The product is positioned for organizations reducing hyperscaler dependency. That obligation reaches into how the product is built, not only how it runs.

| Tier | Requirement |
|---|---|
| Workplace runtime | SHALL function completely without any external model provider. Model-assisted features are optional enhancements that degrade to deterministic behaviour. |
| Forge build-time, accountable owner | MAY use an external cloud model. This is a tool choice by the accountable engineer, not a product dependency. |
| Forge build-time, community contributor | SHALL be operable against a locally hosted model or small language model for at least one complete delivery path, before any claim that the application is community-maintainable. |

No claim of community maintainability, sovereignty, or hyperscaler independence SHALL be published before tier-three evidence exists.

## Existing product baseline

The current source set contains:

- document hierarchy and architecture governance;
- product masterplan and consolidation delta;
- logical, data, trust, deployment, and cryptographic architecture;
- technology selection method and candidate radar;
- shared Intelligence Platform reference architecture;
- Engineering Foundation and AI Software Factory;
- ADR register, roadmap, glossary, and documentation governance.

The source set SHALL be imported as immutable versioned input before consolidation. Original text, status, authorship, and relationships remain traceable. Forge-generated revisions are new versions rather than silent replacements.

## Product-owned concerns

Workplace retains ownership of:

- workplace intent, need, capability, resource, allocation, reservation, hold, and check-in;
- campus, building, floor, zone, geometry, wayfinding, and accessible alternatives;
- occupancy, presence, wireless, BLE/UWB, sensors, and privacy boundaries;
- collaboration, calendar, room-device, Cisco/Webex, and facility adapters;
- decision constraints, scoring, optimization, and product explanations;
- workplace-specific legal, retention, consent, and employee-trust requirements.

## Forge capabilities exercised

| Forge capability | Workplace proof |
|---|---|
| Product and knowledge | Import and reconcile the existing masterplan without loss |
| Software System Twin | Visualize buildings-domain architecture, user journeys, data, deployment, and runtime evidence |
| Agentic workforce | Product, architecture, UX, engineering, test, security, privacy, accessibility, and operations hand-offs |
| Intelligence Platform | Provider-neutral decision assistance, context, agents, tools, and evaluation |
| Execution platform | Reproducible builds, tests, previews, adapters, and controlled environments |
| Trust and cryptography | Tenant isolation, workload identity, envelope encryption, key lifecycle, audit, and recovery |
| Delivery | Signed build-once artefacts, release channels, SaaS/self-hosted deployment, and rollback |
| Operations | OpenTelemetry, SLOs, diagnostics, support bundles, backup, and restore |

## First governed journey

1. Import the ten source documents with hashes, status, and provenance.
2. Build a product knowledge graph and flag contradictions, assumptions, missing evidence, and open decisions.
3. Let product, architecture, UX, privacy, security, and operations agents propose a consolidated product baseline.
4. Require human acceptance of product boundaries and consequential decisions.
5. Generate a Software System Twin and role-specific explanations.
6. Generate the reproducible development, test, and preview environment.
7. Implement a thin vertical slice covering tenant identity, one workplace resource, availability or hold, an accessible user journey, audit, and one provider-neutral adapter boundary.
8. Produce independent tests, threat evidence, SBOM, provenance, signatures, and deployment evidence.
9. Promote the same artefact to a SaaS-like and self-hosted reference profile.
10. Demonstrate rollback, backup/restore, and explanation to an end user, architect, engineer, and operator.

## Lighthouse exit criteria

- Original intent is traceable to accepted requirements, decisions, code, and tests.
- Domain logic contains no direct dependency on a named model, KMS, identity broker, event broker, or collaboration vendor.
- Tenant context cannot be caller-overridden and negative isolation tests pass.
- The core journey meets its declared accessibility target without relying exclusively on 3D interaction.
- Provider outage, KMS outage, rollback, and restore have bounded tested behaviour.
- SaaS and self-hosted use the same signed application artefact.
- The Workbench explains behaviour and residual risk to each target audience.

Generality is assessed after the Workplace lighthouse through the separate second-application milestone. It is not a Workplace product exit criterion.

## Source snapshot register

These hashes identify the source snapshot assessed for this Forge baseline. The files still require a stable version-controlled Workplace product location.

| Source | SHA-256 |
|---|---|
| `00-document-set-overview.md` | `e74ec7cae8871a75d3211cc8c22331970dba9de1eff1ab032e9de5f1ee7c59f4` |
| `01-masterplan-v0.2-source.md` | `0f9bb0ce6297727f91b529840ec3195082377355e48a4d1f5297235686bc2076` |
| `01-masterplan-v0.3-consolidation.md` | `cf7bcfe3d56fd077df654bbd0ca4a2568da871ff425b66a40ce0e39af82cf80f` |
| `02-logical-cryptographic-v0.1-source.md` | `13fe22e5b485b5fb444083de2fd9c2a170285738a4eb73696d2fe3ab2ab8f535` |
| `03-technology-selection-v0.1-source.md` | `0813a4180a5ebcf50b31160fe7e057e4d7100d6fd8b19002078290ea3a690bad` |
| `04-intelligence-platform-reference-architecture-v0.1.md` | `fbc2d242066496fb14575054ca0c20a6cfb563ffbc7429456e489dc63b7b3651` |
| `05-engineering-foundation-ai-software-factory-v0.1.md` | `bbd443c3dfaa2a427838b9e0e0900aef5571dd405f91425be50c774a7f29868a` |
| `06-adr-register-and-roadmap.md` | `ba6bcf121d2ae5bbf5927defb5925f17d8ad1e9793aae49c4fdde495c7565a2f` |
| `07-shared-glossary.md` | `42b8a624ac3eeb486bf8e36cf351711e55d7db2e8d2679b5948e82e445c231f0` |
| `08-documentation-governance.md` | `63c42aea2e1104bb7c56a75de8c8f82ac29787c698fefe0e8f6db8d7b4560cf2` |
