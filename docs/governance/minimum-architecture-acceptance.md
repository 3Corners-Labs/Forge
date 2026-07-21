# Minimum architecture acceptance packet

Status: Proposed
Owner: Forge Architecture Board
Version: 0.1

## Purpose

This packet defines the smallest coherent set of human decisions needed before functional implementation of the trustworthy intent-to-pull-request slice may begin. It does not accept itself or any referenced document.

Reviewers accept, revise, or reject each decision explicitly. Silence, repository presence, merged wording, or AI recommendation is not acceptance.

## Proposed acceptance scope

| Item | Decision requested | Why it is required |
|---|---|---|
| `foundation/constitution.md` | Accept as the initial constitutional guardrails | Establishes human authority, trust, bounded execution, and evidence rules |
| `governance/document-lifecycle.md` | Accept lifecycle, maturity, and evidence vocabularies | Prevents Proposed architecture from being represented as implemented |
| ADR-0001 | Accept modular-control-plane-first posture | Prevents premature distributed architecture |
| ADR-0002 | Accept the architecture metamodel | Makes plane, capability, module, deployment unit, provider, and profile unambiguous |
| ADR-0003 | Accept authority/context separation | Establishes the first-slice trust-object boundary |
| ADR-0004 | Accept composable deployment profiles | Prevents ambiguous SaaS and self-hosted claims |
| ADR-0005 | Accept separate journey and security assurance capabilities | Preserves product intent and bounds offensive automation |
| `architecture/vertical-slice.md` | Accept scope, exclusions, failure behaviour, and exit criteria | Creates one bounded implementation target |
| `architecture/first-slice-capability-contracts.md` | Accept logical ownership and contract baseline | Closes the minimum G5 design gap before implementation |
| `architecture/threat-model.md` | Accept first-slice threat scope, subject to named follow-ups | Makes security assumptions and tests explicit |
| `architecture/first-slice-data-flow.md` | Accept first-slice data movement, minimization, and privacy obligations | Establishes the initial G6 review boundary |

The broader masterplan, full capability catalogue, future deployment profiles, Workbench, agentic workforce, and Workplace lighthouse remain Proposed or Research. They need not be accepted wholesale to begin the first slice.

## Cross-document decisions requiring confirmation

Reviewers SHALL confirm that:

1. Forge is an Engineering Operating System and governed development-team platform, not an LLM application or orchestration framework.
2. The first implementation is a modular control plane with isolated workers, not a microservice programme.
3. Execution Envelope, Policy Decision, Approval Record, Context Package, and Evidence Record remain non-overlapping objects.
4. Deterministic components own authority and durable state; models and tools do not.
5. No production profile is Supported by this acceptance.
6. Named products and protocols remain replaceable candidates unless separately decided.
7. Human approval remains mandatory for publication in the first slice.
8. Accepted user outcomes cannot be changed by the implementation or agent under test.
9. Active security testing is excluded from the first implementation except for safe contract and isolation tests.
10. Acceptance authorizes implementation of the bounded slice only, not the full capability catalogue.

## Required review roles

Acceptance requires accountable human review from:

- Product: scope, user outcome, exclusions, and journey ownership.
- Architecture: metamodel, boundaries, contracts, and evolution path.
- Security and privacy: threat scope, policy gates, isolation, data handling, and evidence.
- Engineering: feasibility, testability, state, failure, and recovery semantics.
- Operations: operability hypotheses, evidence availability, and recovery ownership.

One person MAY hold multiple accountable roles when organizational policy permits, but the record SHALL make that explicit. AI agents cannot satisfy an accountable role.

## Decision record template

The acceptance change SHALL record:

| Field | Required value |
|---|---|
| Decision ID | Stable repository identifier |
| Exact revision | Git commit digest reviewed |
| Items accepted | Paths and versions |
| Items revised or excluded | Paths, sections, and rationale |
| Product reviewer | Identity, outcome, date |
| Architecture reviewer | Identity, outcome, date |
| Security/privacy reviewer | Identity, outcome, date |
| Engineering reviewer | Identity, outcome, date |
| Operations reviewer | Identity, outcome, date |
| Residual risks | Owned, dated, and linked |
| Implementation authorization | Exact slice, expiry or review condition |

Approval SHOULD be captured through an ADR or governance decision that updates referenced document statuses in the same reviewable change. Status changes SHALL NOT be prepared by the candidate implementation itself after review.

## Readiness after acceptance

Acceptance of this packet closes only the decision portion of G1–G5. Implementation still depends on:

- closing G0 developer-safety checks;
- assigning and approving the profile-specific retention, provider, and privacy decisions identified under G6;
- defining a production-shaped reference-profile test plan under G7 without claiming support;
- maintaining passing fitness functions under G8;
- accepting the scenario and result contracts for G9;
- keeping active and autonomous offensive operation outside the slice under G10.

Functional work SHALL start with contract tests and one thin end-to-end path. It SHALL NOT broaden scope merely because a future capability appears in the catalogue.
