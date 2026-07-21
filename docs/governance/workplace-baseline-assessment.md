# Workplace Orchestrator baseline assessment

Status: Accepted
Owner: Forge Architecture

## Purpose

This register records how the ten documents in the Workplace Orchestrator architecture baseline v0.3 enter Forge. Forge itself is built first as the application development team; Workplace Orchestrator is that team's first business application. Reusable platform patterns shape Forge, while workplace-specific product intent and domain architecture remain owned by the Workplace product workspace.

## Disposition terms

- Adopt: usable as a Forge principle with only naming changes.
- Adapt: valuable pattern requiring Forge-specific scope or challenge.
- Defer: retain as candidate or decision backlog pending evidence.
- Product-owned: preserved as Workplace Orchestrator input but not embedded in the Forge platform core.

## Source mapping

| Source document | Forge disposition | Integrated destination |
|---|---|---|
| `00-document-set-overview.md` | Adopt single-responsibility document hierarchy; adapt lifecycle statuses | Masterplan, documentation lifecycle |
| `01-masterplan-v0.2-source.md` | Adopt cross-product principles; retain workplace vision, users, domain, privacy constraints, journeys, and roadmap as product-owned input | Forge masterplan plus Workplace product workspace |
| `01-masterplan-v0.3-consolidation.md` | Adopt separation of product, shared platform, Intelligence Platform, and Engineering Foundation; adopt same-artifact and useful-community principles | Masterplan product structure and deployment |
| `02-logical-cryptographic-v0.1-source.md` | Adopt reusable trust and crypto platform patterns; retain workplace data classifications, flows, and quality scenarios as product-owned architecture | Forge security architecture plus Workplace product workspace |
| `03-technology-selection-v0.1-source.md` | Adopt evidence-led selection, contract-first technology, outbox progression, build-once promotion, supply-chain, and technology-radar method; defer specific products to ADR/PoC | Masterplan governance and decision backlog |
| `04-intelligence-platform-reference-architecture-v0.1.md` | Adopt capability gateway, router, Agent Runtime, Tool Gateway, Context Platform, registry, evaluation, observability, TaskRequest, AgentDefinition, deployment profiles, and bounded multi-model collaboration | Masterplan and intelligence architecture |
| `05-engineering-foundation-ai-software-factory-v0.1.md` | Adopt controlled agent resources, role separation, context packages, quality gates, release channels, and implementation sequence; expand UX/product/operations workforce | Masterplan agentic workforce and roadmap |
| `06-adr-register-and-roadmap.md` | Adapt reusable INT/ENG/DEP/CRY decisions; retain WO decisions in the Workplace decision register | Forge and Workplace decision backlogs |
| `07-shared-glossary.md` | Adopt generic platform terms; retain workplace terms in the product glossary | Forge glossary plus Workplace product glossary |
| `08-documentation-governance.md` | Adopt evidence labels, review dimensions, required capability artefacts, and change lifecycle | Masterplan governance and documentation lifecycle |

## Architecture challenges applied

The source baseline is not copied uncritically:

1. Kubernetes remains a strong production reference candidate, not a requirement for every compact self-hosted Forge profile.
2. Named products such as Keycloak, OpenBao, PostgreSQL, NATS, Argo CD, Harbor, and Kyverno remain Trial candidates until Forge-specific quality scenarios and operational rehearsals pass.
3. Forge needs a broader workforce and product model than programming agents alone, including discovery, product, UX, operations, support, and end-user explanation.
4. Workplace's Digital Building pattern becomes a Software System Twin derived from engineering and runtime evidence.
5. Cryptographic patterns are adopted, but final algorithms, key domains, service topology, and KMS profiles require a dedicated Forge threat model and ADRs.
6. Community and Enterprise separation remains a product/legal decision; foundational security, accessibility, portability, recovery, and audit cannot be premium-only.
7. The workplace-specific domain model, booking concurrency, occupancy, indoor navigation, Cisco/Webex adapters, and sensor analytics remain mandatory product inputs and lighthouse scenarios, but are not Forge core capabilities.

## Provenance limitation

The source documents currently exist outside this repository and are not yet versioned Forge product artefacts. Before implementation, ownership, licensing, version, and a stable Workplace product workspace SHALL be confirmed. The original files SHALL be imported without destructive rewriting so Forge can preserve provenance and demonstrate controlled consolidation.
