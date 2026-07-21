# Architecture readiness register

Status: Proposed
Owner: Forge Architecture Board
Version: 0.2

## Decision

Forge is **not implementation-ready** beyond safety hardening of the existing development bootstrap and architecture-as-code work.

Functional capability implementation may begin only after the mandatory gates below are accepted. This register turns the enterprise architecture review into a traceable remediation plan.

## Strengths to preserve

- Engineering before intelligence.
- Human accountability and explicit approval boundaries.
- Models, tools, and retrieved content are untrusted by default.
- Capability-first and provider-neutral architecture.
- Modular-monolith-first implementation posture.
- Context, knowledge, memory, evidence, economics, and operations as engineering disciplines.
- Same logical contracts for SaaS and self-hosted deployment.
- Source preservation and selective context routing.

## Mandatory readiness gates

| Gate | Required outcome | Evidence | Status |
|---|---|---|---|
| G0 Development safety | Existing bootstrap is explicitly developer-only, loopback-bound, and rejected for non-developer profiles | Automated configuration and Compose checks | In progress |
| G1 Governance | One lifecycle taxonomy and an accepted minimum authority set | [ADR-0006](../decisions/0006-accept-minimum-architecture-baseline.md) and accepted authority set | Accepted |
| G2 Metamodel | Plane, capability, module, deployment unit, provider and profile relationships are unambiguous | Accepted metamodel and ADR-0002 | Accepted |
| G3 First slice | One bounded intent-to-PR journey with explicit exclusions and quality hypotheses | Accepted vertical-slice architecture | Accepted |
| G4 Trust contracts | Execution Envelope, Context Package, Policy Decision, Approval Record and Evidence Record are versioned and non-overlapping | Schemas and accepted ADR-0003 | Accepted design baseline |
| G5 Capability contracts | First-slice capabilities define ownership, data, commands, events, errors, security and NFRs | [Accepted first-slice capability contracts](../architecture/first-slice-capability-contracts.md) | Accepted |
| G6 Threat and privacy | Tenant isolation, data classification, misuse, privacy and supply-chain threats have controls and tests | Accepted [threat model](../architecture/threat-model.md) and [first-slice data-flow review](../architecture/first-slice-data-flow.md); implementation evidence pending | In progress |
| G7 Deployment evidence | Developer and one production-shaped reference profile have explicit conformance requirements | Profile specification and test plan | Open |
| G8 Fitness functions | Architecture, schema, dependency, isolation, source integrity and documentation rules run in CI | Passing required checks | In progress |
| G9 User intent assurance | Accepted multi-role scenarios have independent ownership, public-interface semantics, expected outcomes and reproducible evidence | Scenario/result schemas, ADR-0005 and Trial test plan | Proposed |
| G10 Offensive assurance safety | Active and autonomous security tests have rules of engagement, external scope enforcement, impact controls, kill switch and finding governance | Threat review, safety conformance plan and ADR-0005 | Proposed |

Only G0 safety hardening and work required to close G1–G10 are authorized before readiness acceptance.

## Findings and disposition

| Finding | Required disposition |
|---|---|
| No Accepted architecture baseline | Resolved by ADR-0006 for the bounded first slice; broader future-state documents remain Proposed |
| Conflicting lifecycle vocabularies | Separate document status, technology maturity and evidence labels |
| Missing capability contracts | Specify only capabilities needed for the first vertical slice, then expand incrementally |
| Execution Envelope and Context Package overlap | Define separate authority and information objects with one-way derivation |
| Overloaded Control and Assurance boundaries | Split logical responsibilities without prematurely creating network services |
| Evidence and telemetry combined | Separate lifecycle, storage and integrity contracts; allow deliberate promotion |
| Central knowledge graph overreach | Use a federated, rebuildable projection over capability-owned records |
| Inconsistent deployment profile catalogues | Adopt one composable profile model and maturity matrix |
| Unquantified scalability | Record validation envelopes and measurable quality scenarios per supported profile |
| Extension trust incomplete | Add publisher, permission, provenance, revocation and isolation governance |
| Old and proposed documentation both appear canonical | Keep source accessible but visually and structurally separate from current navigation |
| Workplace exit depends on a later second application | Move second-application proof to the Forge generality milestone |
| Journey tests lack an owner for accepted intent | Add User Journey Assurance and prevent an agent or implementation from redefining expected outcomes |
| Automated pentesting is unspecified and potentially unsafe | Layer checks by risk and require external scope, rules of engagement, kill switch, validation, and human assurance |

## Implementation boundary

The first functional implementation SHALL NOT attempt to build the full capability catalogue. It SHALL prove this path:

~~~text
Human intent
  -> authenticated workspace
  -> authorized execution envelope
  -> bounded context package
  -> reviewed delivery plan
  -> isolated repository change
  -> independent verification
  -> accepted role-based journey evidence where a public interface is affected
  -> human approval
  -> pull request and evidence record
  -> role-appropriate explanation
~~~

Production deployment, autonomous merge, billing, certification, accreditation, a public extension marketplace, and Workplace-specific functionality are outside this boundary.

## Review ownership

Human acceptance is required from product, architecture, security/privacy, engineering, and operations owners. An AI agent may prepare evidence and identify contradictions but cannot mark its own architecture Accepted.

The [minimum architecture acceptance packet](minimum-architecture-acceptance.md) identifies the smallest coherent human decision set. The [implementation baseline](implementation-baseline.md) prevents target architecture from being confused with verified repository state while that review remains open.
