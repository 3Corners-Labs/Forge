# Existing-content migration map

Status: Accepted
Owner: Forge Architecture

## Purpose

The current Forge Markdown corpus contains substantial valuable thinking as well as duplication and contradictions. No existing document will be removed until its concepts have been mapped, consolidated, reviewed, and recoverably preserved in Git history.

## Migration rule

For each source document:

1. inventory its unique concepts and normative statements;
2. assign each concept to one canonical destination;
3. reconcile contradictions through an ADR or explicit review decision;
4. verify that links, rationale, and intent remain discoverable;
5. mark the source Superseded with a destination link;
6. remove or archive it only in a later, separately reviewable change.

## Document mapping

| Existing source | Canonical destination | Preserved intent | Current disposition |
|---|---|---|---|
| `/README.md` | `/README.md`, vision, architecture, quality, roadmap | Vision, lifecycle, trust, extensibility, quality, security, accessibility, observability, stack and roadmap | Retain until section-level verification |
| `/SYSTEM.md` | architecture overview, security, deployment, quality | Planes, capability flow, lifecycle, trust, resilience, scalability and constraints | Retain; remove internal duplicate only after review |
| `docs/CONSTITUTION.md` | `foundation/constitution.md` | Human accountability, everything-as-code, openness, security, privacy, context, knowledge, audit, economics, stewardship and sustainability | Primary consolidation source |
| `docs/standards/engineering-constitution.md` | `foundation/constitution.md` | Concise principles, controls and conflict resolution | Reconcile with full constitution |
| `docs/ENGINEERING_INTENT.md` | `foundation/vision.md` and roadmap | Mission, outcomes, priorities, trade-offs, anti-goals and improvement | Retain until concept check |
| `docs/DESIGN_PHILOSOPHY.md` | vision and experience | Human-centred simplicity and progressive disclosure | Consolidate |
| `docs/SYSTEM_OVERVIEW.md` | `architecture/overview.md` | Modular posture, planes, lifecycle and trust | Primary overview source |
| `docs/ARCHITECTURE.md` | all `architecture/` documents | Capabilities, context, memory, policy, events, runtime, resilience, deployment, data, providers and governance | Split by concern; duplicated sections excluded only after verification |
| `docs/CAPABILITY_CATALOG.md` | `architecture/capabilities.md` | AI-supported delivery tasks | Expand into enterprise capability model |
| `docs/SECURITY.md` | `architecture/security.md` | Assets, threats, controls, environments and required artefacts | Preserve and deepen |
| `docs/PROJECT_MEMORY.md` | `architecture/knowledge.md` | Project memory as governed knowledge | Preserve and deepen |
| `docs/INTELLIGENCE_MODEL.md` | `architecture/intelligence.md` | Intelligence as governed capability | Preserve and deepen |
| `docs/MODEL_SELECTION_POLICY.md` | `architecture/intelligence.md` | Policy, performance, cost and sustainability routing | Preserve and deepen |
| `docs/LLM_BOOTSTRAP.md` | `architecture/intelligence.md` | Provider-independent bootstrap | Preserve and deepen |
| `docs/TOKEN_STEWARDSHIP.md` | `architecture/intelligence.md` | Efficient, visible token use | Preserve and deepen |
| `docs/COMPUTE_STEWARDSHIP.md` | `architecture/intelligence.md` | Purposeful measurable compute use | Preserve and deepen |
| `docs/TASK_WORKFLOW.md` | `governance/delivery.md` | Branch, review, verification and definition of done | Preserve and deepen |
| `docs/AGENTS.md` | `/AGENTS.md` | Human/AI contributor contract | Replace only after approval |
| `docs/REPOSITORY_MAP.md` | `docs/README.md` and documentation lifecycle | Discoverability, ownership and authority | Preserve and update after moves |
| `docs/ENGINEERING_LANGUAGE.md` | `reference/glossary.md` | English canonical language and precise terminology | Preserve |
| `docs/adr/*` | `docs/decisions/` | ADR template and decision history | Migrate without losing history |
| `docs/rfc/*` | future proposal workflow and decisions | Repository strategy proposal | Retain pending decision |
| `docs/CODE_OF_CONDUCT.md` | governance, unchanged | Community conduct | Retain as-is |
| Author-provided *3Corners Forge Principles* | `foundation/principles.md`, constitution, architecture, profiles and future conformance specifications | Complete Engineering Operating System philosophy, capability theory, draft specifications, profiles, conformance, certification and manifesto | Preserve source snapshot; migrate by chapter using the source map |
| *3Corners Forge Principles Addendum v0.1* | metamodel, trust contracts, deployment, readiness and future ADRs | Insights from the enterprise architecture challenge without modifying the original manuscript | Preserve as Proposed addendum; promote only through review |
| *3Corners Forge Principles Addendum v0.2* | user journey assurance, continuous security testing, capability model and ADR-0005 | End-user intent preservation and bounded automated adversarial assurance | Preserve as Proposed addendum; promote only through review |

## Cross-cutting concept checklist

The consolidation is incomplete until it demonstrably covers repository-as-truth, everything-as-code, knowledge-as-code, context-as-code, intentional memory, human accountability, identity before trust, policy before execution, evidence before trust, vendor independence, open standards, privacy, economic intelligence, token and compute stewardship, sustainability, continuous validation, localization, accessibility, observability, resilience, extensibility, progressive delivery, deployment profiles, objective conformance, and governed standards evolution.

The [Forge Principles source map](../source-material/forge-principles-map.md) is the chapter-level migration ledger for the author-provided manuscript. Its duplicate Chapter 26 and absent Chapter 27 are source anomalies, not permission to discard content or silently renumber the archive.
