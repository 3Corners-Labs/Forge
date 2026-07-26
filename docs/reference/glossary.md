# Glossary

Status: Accepted
Owner: Forge Architecture

- Agent: an identity-bound software actor executing a role within delegated authority.
- Agentic workforce: coordinated specialist agents and humans spanning the software lifecycle.
- Accepted scenario: a reviewed, versioned user journey whose expected and forbidden outcomes may contribute to release eligibility.
- Autonomous security testing: offensive testing in which software selects or adapts targets, techniques, or exploitation steps within externally enforced rules of engagement.
- Capability: a durable business responsibility independent of deployment topology.
- Deciding independent check: a control that Forge owns, that may determine the outcome, and that is capable of failing for the paths a candidate changes. Verification is refused when a plan has none, because certifying a candidate against checks that cannot fail is empty assurance. Earlier documents call this an "independent consequential oracle"; that phrase overstated a bounded mechanism and is no longer used in new writing.
- Delivery run: a traceable instance that transforms intent into artefacts and evidence.
- Evidence: structured, attributable proof of decisions, execution, validation, and outcomes.
- End-User Assurance Agent: a governed role that executes and explores scenarios as least-privilege personas through public application interfaces without owning product intent or release approval.
- Execution envelope: immutable constraints and authority supplied to an agent or tool invocation.
- Test oracle: the kind of authority that decides whether an observed outcome is correct, recorded per expected outcome as deterministic, reference, statistical, human, or semantic-supporting. This is the established testing term and is retained; it is distinct from a deciding independent check, which is one concrete control.
- Forge Workbench: the role-aware interface for collaboration, visualization, explanation, and approval.
- Intent-driven development: progressive conversion of human outcomes into governed delivery artefacts.
- Exploratory scenario: a non-authoritative journey used to discover behaviour, risks, or product opportunities; it becomes release-gating only after accountable review.
- Plane: a logical grouping of capabilities with a shared architectural concern.
- Provider adapter: a replaceable boundary around an external model, tool, or platform.
- Rules of engagement: a versioned authorization contract defining scope, targets, permitted techniques, limits, stop conditions, evidence handling, cleanup, and accountability for active security testing.
- User Journey Assurance: the capability that links intent, personas, accepted scenarios, public-interface execution, expected outcomes, and reproducible journey evidence.
