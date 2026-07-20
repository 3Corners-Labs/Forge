# Forge Principles source map

Status: Proposed
Owner: Forge Architecture
Source ID: `forge.principles.manuscript.v0.1`

## Purpose

This map makes the complete *3Corners Forge Principles* manuscript selectively retrievable. It is a navigation and migration artefact, not a substitute for the source.

Later architectural insights are preserved in separate [v0.1](3corners-forge-principles-addendum-v0.1.md) and [v0.2](3corners-forge-principles-addendum-v0.2.md) addenda. The addenda complement but never rewrite the original chapter structure or each other. Addendum v0.2 records end-user intent assurance and bounded continuous security testing.

## Source assessment

The manuscript contains a coherent architecture philosophy and a substantial draft standard. It also contains editorial conversation, mixed Dutch and English, repeated transition text, a repeated Chapter 26 Knowledge Specification, no Chapter 27, and proposals that have not been accepted. These characteristics are preserved as provenance and SHALL NOT be silently normalized in the source snapshot.

## Structural map

| Chapters | Source theme | Canonical destination |
|---|---|---|
| 1–4 | Problem, Engineering Operating System, vision, philosophy | `foundation/vision.md`, `foundation/principles.md`, `foundation/masterplan.md` |
| 5–13 | Capability, context, memory, knowledge, runtime, policy, intelligence, evidence | `architecture/capabilities.md`, `architecture/knowledge.md`, `architecture/agentic-delivery.md`, `architecture/intelligence.md` |
| 14–21 | Models, tools, identity, security, observability, economics, governance, architecture | relevant documents under `architecture/` and `governance/` |
| 22–42 | Reference architecture, canonical specifications, standards, adoption, measurement and evolution | architecture specifications, schemas, fitness functions, ADRs and roadmap |
| 43–50 | Core, enterprise, sovereign, government, air-gapped, agent, edge and composed profiles | `architecture/deployment.md` and future profile specifications |
| 51–54 | Conformance, certification, accreditation and standards governance | future assurance and standards-governance specifications |
| 55 | Manifesto | `foundation/principles.md` and product narrative |

## Chapter catalogue

| Chapter | Lines in original | Topic | Primary retrieval purpose |
|---:|---:|---|---|
| 1 | 11–248 | The AI Engineering Problem | Why probabilistic systems need engineering discipline |
| 2 | 249–550 | Engineering Operating Systems | Platform thesis and stable capability contracts |
| 3 | 551–830 | Vision | Long-term product and human-accountability intent |
| 4 | 831–1222 | Design Philosophy | Twelve enduring principles |
| 5 | 1223–1605 | Capability Model | Capability landscape, contracts and maturity |
| 6 | 1606–1962 | Context Engineering | Minimal sufficient, governed context |
| 7 | 1963–2366 | Memory Engineering | Memory scopes, ownership and lifecycle |
| 8 | 2367–2807 | Knowledge Engineering | Knowledge governance, provenance and repository truth |
| 9 | 2808–3202 | Runtime Engineering | Deterministic execution around probabilistic capability |
| 10 | 3203–3593 | Policy Engineering | Declarative, explainable policy before execution |
| 11 | 3594–3992 | Intelligence Engineering | Governed and measurable system intelligence |
| 12 | 3993–4369 | Evidence Engineering | Evidence as a first-class trust artefact |
| 13 | 4370–4708 | Capability Engineering | Boundaries, contracts, taxonomy and lifecycle |
| 14 | 4709–5081 | Model Engineering | Replaceable models, routing and evaluation |
| 15 | 5082–5488 | Tool Engineering | Governed tools, contracts and MCP positioning |
| 16 | 5489–5844 | Identity Engineering | Human, workload, agent and delegated identity |
| 17 | 5845–6207 | Security Engineering | Zero Trust and AI-specific threat model |
| 18 | 6208–6559 | Observability Engineering | Context lineage, tracing and feedback |
| 19 | 6560–6902 | Economics Engineering | Cost, value, sustainability and compute stewardship |
| 20 | 6903–7311 | Governance Engineering | Ownership, ADRs and controlled evolution |
| 21 | 7312–7708 | Architecture Engineering | Boundaries, constraints and invariants |
| 22 | 7709–8103 | Forge Reference Architecture | Layering, request flow and conformance |
| 23 | 8104–8478 | Capability Specification Framework | Required specification structure and identifiers |
| 24 | 8479–8873 | Context Package Specification | Canonical context object and lifecycle |
| 25 | 8874–9289 | Memory Specification | Canonical memory object and promotion |
| 26a | 9290–9654 | Knowledge Specification | Canonical governed knowledge object |
| 26b | 9655–10024 | Repeated Knowledge Specification | Deduplicate equivalent requirements in retrieval results while preserving both source occurrences |
| 28 | 10025–10424 | Policy Specification | Deterministic decisions, obligations and evidence |
| 29 | 10425–10798 | Model Specification | Provider-neutral model capability |
| 30 | 10799–11235 | Tool Specification | Tool discovery, risk and execution contract |
| 31 | 11236–11613 | Identity Specification | Identity chains and federation |
| 32 | 11614–12009 | Security Specification | Trust boundaries and continuous security |
| 33 | 12010–12382 | Evidence Specification | Evidence object, lifecycle and integrity |
| 34 | 12383–12789 | Observability Specification | Telemetry and operational visibility |
| 35 | 12790–13162 | Governance Specification | Governance decisions and ownership |
| 36 | 13163–13535 | Architecture Specification | Architecture objects and conformance |
| 37 | 13536–13875 | Specification Framework | Standardized specification lifecycle |
| 38 | 13876–14174 | Standards Governance | Ownership and controlled standards evolution |
| 39 | 14175–14470 | Adoption and Transformation | Incremental organizational adoption |
| 40 | 14471–14802 | Engineering Measurement | Quality and maturity measurement |
| 41 | 14803–15120 | Continuous Optimization | Evidence-driven improvement |
| 42 | 15121–15415 | Architectural Evolution | Compatibility and governed evolution |
| 43 | 15416–15695 | Core Profile | Minimum conforming capability set |
| 44 | 15696–15965 | Enterprise Profile | Enterprise integration and governance |
| 45 | 15966–16254 | Sovereign Profile | Sovereignty, residency and control |
| 46 | 16255–16536 | Government Profile | Public-sector accountability and assurance |
| 47 | 16537–16805 | Air-Gapped Profile | Disconnected operation and supply chain |
| 48 | 16806–17103 | AI Agent Profile | Agent identity, authority and lifecycle |
| 49 | 17104–17393 | Edge Profile | Constrained distributed execution |
| 50 | 17394–17637 | Profile Composition | Independent composable profile dimensions |
| 51 | 17638–17900 | Conformance Framework | Objective, evidence-based conformance |
| 52 | 17901–18207 | Certification Framework | Independent certification lifecycle |
| 53 | 18208–18696 | Accreditation Framework | Governance of certification authorities |
| 54 | 18697–19021 | Forge Standards Governance | Evolution and ownership of the standard |
| 55 | 19022–19299 | Forge Manifesto | Enduring public narrative and ten decision principles |

## Promotion rules

A manuscript statement becomes canonical only when:

1. its complete surrounding rationale has been reviewed;
2. contradictions with current Forge material are identified;
3. one canonical destination is selected;
4. normative strength is chosen deliberately;
5. acceptance and evidence requirements are defined;
6. the migration is recorded in `governance/content-migration.md`.
