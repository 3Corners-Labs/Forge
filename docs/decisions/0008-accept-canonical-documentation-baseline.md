# ADR-0008: Accept the canonical documentation baseline

Status: Accepted
Owner: Forge Architecture Board
Date: 2026-07-21

## Context

Forge had an Accepted minimum architecture for the bounded intent-to-pull-request slice, but most product, target-architecture, governance, reference, and roadmap documents remained Proposed. That limited authority was intentional while contradictions were being reconciled. It now obscures the reviewed product intent and makes canonical target documents appear no more authoritative than unresolved proposals.

Sijbren Beukenkamp authorized promotion of every canonical document judged ready on 2026-07-21. Under the temporary multi-role authority recorded in ADR-0006, this authorization represents Product, Architecture, Security and Privacy, Engineering, and Operations review. The exact reviewed source baseline is Git commit `83a962e2d3b2e8ef453094a23f6057cd76fbb6dc`.

Acceptance remains separate from capability maturity, implementation evidence, readiness-gate status, or production support. An Accepted target document may govern a Research capability and may require later ADRs or validation before a Trial or Supported profile is authorized.

## Decision

Accept the canonical documentation baseline in these groups:

- product and direction: `docs/foundation/vision.md`, `docs/foundation/principles.md`, `docs/foundation/masterplan.md`, and `docs/roadmap.md`;
- architecture orientation and capability boundaries: `docs/architecture/overview.md`, `capabilities.md`, `identity-and-tenancy.md`, `policy-and-authorization.md`, `runtime-and-workflow.md`, `security.md`, `data-and-privacy.md`, `deployment.md`, `evidence-and-observability.md`, `experience.md`, `knowledge.md`, `intelligence.md`, `quality-and-operations.md`, `user-journey-assurance.md`, `continuous-security-testing.md`, `software-supply-chain.md`, `extensibility.md`, and `agentic-delivery.md`;
- governance and navigation: `docs/README.md`, `docs/context/README.md`, `docs/decisions/README.md`, `docs/governance/architecture-readiness.md`, `implementation-baseline.md`, `delivery.md`, `content-migration.md`, and `workplace-baseline-assessment.md`;
- reference intent: `docs/reference/glossary.md` and `docs/reference-applications/workplace-orchestrator.md`.

The acceptance change also:

- reconciles the masterplan status in the roadmap while preserving ADR-0007 as the implementation-sequencing authority;
- updates the security architecture to reflect that ADR-0003 is already Accepted;
- replaces the implementation-readiness context-pack dependency on the fulfilled acceptance packet with ADR-0006 and ADR-0007;
- marks `docs/governance/minimum-architecture-acceptance.md` Superseded by ADR-0006.

Keep `docs/architecture/cryptography.md` Proposed. Its required quality scenarios have not been validated, and its final key hierarchy and provider decisions remain open. Keep the ADR template Draft. Source manuscripts retain instruction authority `none` and are not promoted by this decision.

## Consequences

- The accepted documents are normative within their stated target, governance, or reference scope.
- The architecture-readiness register remains authoritative that Forge is not yet ready for M1 functional implementation.
- All unimplemented capabilities remain Research unless a separate lifecycle decision and objective evidence promote them.
- No production profile becomes Supported, and no technology candidate becomes selected, through this decision.
- ADR-0007 governs the implementation order where the masterplan's product-level stages differ from the roadmap.
- The fulfilled minimum acceptance packet remains available as decision history without competing with ADR-0006.
- Cryptographic architecture cannot become normative until its acceptance prerequisites are reviewed or explicitly revised through a later decision.

## Accountable review record

| Role | Reviewer | Outcome | Date |
|---|---|---|---|
| Product | Sijbren Beukenkamp | Accept listed scope | 2026-07-21 |
| Architecture | Sijbren Beukenkamp | Accept listed scope | 2026-07-21 |
| Security and Privacy | Sijbren Beukenkamp | Accept listed scope; exclude cryptography | 2026-07-21 |
| Engineering | Sijbren Beukenkamp | Accept listed scope | 2026-07-21 |
| Operations | Sijbren Beukenkamp | Accept listed scope | 2026-07-21 |

The AI contributor prepared the review and consistency changes but did not supply accountable human acceptance.

## Residual risks

- F0 gates G0 and G6 through G10 remain incomplete and continue to block M1 functional implementation.
- Cryptographic conformance, key hierarchy, provider selection, recovery, and destructive-key controls remain unresolved.
- Source-content migration and Workplace source import remain incomplete.
- One person temporarily represents all accountable review roles, reducing independence.
- Broad target acceptance can be mistaken for implemented capability; the implementation baseline and maturity labels remain mandatory controls against that error.

## Validation

- Documentation and architecture checks SHALL pass with the promoted statuses and new decision register entry.
- The implementation baseline SHALL continue to enumerate unimplemented capabilities and unsupported profiles.
- Readiness work SHALL close every mandatory F0 gate before M1 functional implementation begins.
- Any later contradiction with the accepted baseline requires an explicit change and, when architecturally significant, a superseding ADR.

## Supersession

This decision expands the accepted documentation scope established by ADR-0006. It does not supersede ADR-0006 or ADR-0007.
