# F0 readiness evidence and ownership plan

Status: Accepted
Owner: Forge Architecture Board
Version: 0.1

## Purpose

This plan defines how Forge closes the F0 foundation-readiness gate without confusing accepted design, pre-implementation evidence, M1 implementation evidence, or Trial promotion evidence.

F0 authorizes only safety hardening, architecture-as-code, contract tests, test harnesses, governance decisions, and fitness functions needed to make M1 safe to start. It does not authorize the M1 functional trust path, claim that a capability is implemented, or promote a capability beyond Research.

## Authority and scope

The Accepted constitution, ADR-0001 through ADR-0009, architecture readiness register, first-slice architecture, capability contracts, threat model, data-flow analysis, and roadmap govern this plan.

This plan clarifies evidence timing and ownership. It does not remove a control or weaken an accepted exit criterion. If closing a gate requires changing an accepted capability boundary, trust invariant, implementation sequence, or assurance obligation, work stops until an ADR records the decision.

## Gate status vocabulary

Gate progress uses only these values:

| Status | Meaning |
|---|---|
| Open | Required work has not started or has no accepted delivery plan |
| In progress | Authorized readiness work or evidence preparation is underway |
| Ready for review | Every pre-M1 evidence item is present and independently reproducible, but accountable acceptance is pending |
| Accepted for F0 | Accountable reviewers accepted the gate for its pre-M1 scope at an exact revision |

Gate status is not document status, capability maturity, evidence strength, or implementation support. An Accepted document may govern an In progress gate and a Research capability.

## Evidence phases

Every requirement belongs to exactly one phase:

| Phase | Purpose | Effect on authorization |
|---|---|---|
| Before M1 | Proves that design, safety constraints, ownership, test contracts, and existing bootstrap controls are sufficient to begin implementation | Required for `Accepted for F0` |
| M1 exit | Proves that the implemented intent-to-pull-request path conforms to its accepted contracts | Required to complete M1; cannot be claimed during F0 |
| Before Trial | Proves that the integrated implementation operates within the declared Developer Trial envelope | Required for Trial promotion; cannot be claimed from design alone |

Evidence that depends on functionality not yet authorized for implementation SHALL be assigned to M1 exit or Before Trial. F0 must still accept its test contract, fixture design, enforcement point, owner, and required outcome before M1 starts.

## Ownership model

| Gate | Delivery owner role | Accountable review roles |
|---|---|---|
| G0 | Forge Engineering and Platform Engineering | Engineering, Security, Operations |
| G1 | Forge Architecture Board | Product, Architecture, Security/Privacy, Engineering, Operations |
| G2 | Forge Architecture | Architecture, Engineering |
| G3 | Forge Product and Architecture | Product, Architecture, Security/Privacy, Engineering, Operations |
| G4 | Forge Security and Knowledge Architecture | Architecture, Security/Privacy, Engineering |
| G5 | Forge Product and Architecture | Product, Architecture, Security/Privacy, Engineering, Operations |
| G6 | Forge Security and Privacy Architecture | Security/Privacy, Architecture, Engineering, Operations |
| G7 | Forge Platform Engineering | Architecture, Security/Privacy, Engineering, Operations |
| G8 | Forge Architecture and Engineering | Architecture, Security/Privacy, Engineering |
| G9 | Forge Product, Experience, and Quality Architecture | Product, Security/Privacy, Engineering |
| G10 | Forge Security and Quality Architecture | Security/Privacy, Architecture, Engineering, Operations |

Sijbren Beukenkamp is the temporary named delivery owner and accountable human reviewer for G0 through G10 under explicit instruction and the multi-role authority recorded by ADR-0006 and ADR-0008. That temporary concentration does not create independent separation of duties and remains a recorded residual risk until distinct owners and reviewers are appointed.

An AI contributor may prepare evidence, tests, contradictions, and recommendations. It cannot accept a gate, assign itself an accountable role, or provide the sole evidence for its own consequential output.

## Gate evidence assignments

### G0 — Development safety

Before M1:

- configuration and CLI tests reject non-developer runtime profiles before Docker or another side effect is invoked;
- Compose validation proves loopback-only bindings and identifies developer-only services, credentials, exporters, and storage;
- one production-shaped negative fixture fails closed;
- Developer Trial limitations and recovery instructions are documented;
- an individual delivery owner and accountable reviewers are recorded.

M1 exit:

- the M1 CLI and worker cannot route around the Developer profile restriction.

Before Trial:

- the complete Developer Trial profile conformance suite passes on a clean supported host.

### G1 through G5 — Accepted design baseline

Before M1:

- accepted documents and ADRs remain internally consistent and pass architecture fitness checks;
- no status promotion is treated as implementation evidence;
- changes to the accepted boundary use normal lifecycle review and an ADR where required.

M1 exit and Before Trial:

- implementation and Trial evidence demonstrate conformance without redefining the accepted contracts.

### G6 — Threat and privacy readiness

Before M1:

- every M1 threat maps to a control, enforcement point, negative test contract, evidence obligation, and owner;
- Developer Trial decisions cover identity assurance, tenant topology, repository credentials, worker isolation, context classification, retention, evidence integrity, telemetry minimization, cancellation, revocation, and source-control destination;
- seeded-secret and personal-data fixtures plus disclosure-safe expected results are versioned;
- tests applicable to the existing bootstrap and schema contracts pass;
- tests requiring M1 functionality are explicitly assigned to M1 exit or Before Trial rather than claimed as complete.

M1 exit:

- cross-scope, authority, context, secret, cancellation, revocation, evidence, and publication tests pass against the implemented public contracts.

Before Trial:

- integrated cross-tenant, deletion, cleanup, restore, telemetry, and evidence-disclosure suites pass within the Developer Trial profile.

### G7 — Deployment evidence

Before M1:

- the Developer Trial profile specification declares every profile dimension, deployment unit, provider assumption, trust boundary, resource envelope, and unsupported property;
- one production-shaped reference profile has a specification and executable conformance test plan while remaining Research;
- identity, storage, repository, isolation, network, secret, evidence, telemetry, backup, restore, upgrade, rollback, capacity, and recovery assumptions have owners;
- no production support claim is made.

M1 exit:

- M1 operates within the accepted Developer Trial envelope.

Before Trial:

- the Developer Trial profile passes its conformance suite; the production-shaped profile remains Research until separately validated.

### G8 — Architecture fitness functions

Before M1:

- required CI enforces document metadata, lifecycle baseline, navigation, links, ADR registration, schemas, source integrity, capability IDs, and developer-profile restrictions;
- dependency and control-plane/worker isolation rules have executable checks with representative failing fixtures;
- journey-contract and offensive-safety rules have executable structural checks;
- required checks cannot be silently skipped by changing generated or candidate-controlled files;
- CI ownership and failure-recovery instructions are documented.

M1 exit:

- the implemented modules and worker boundary pass dependency and isolation checks.

Before Trial:

- required fitness checks pass from a clean checkout in the protected CI workflow.

### G9 — User intent assurance

Before M1:

- one allowed and one denied M1 scenario are versioned and owned independently from the implementation;
- each scenario defines persona, entitlements, public CLI boundary, fixtures, expected and forbidden outcomes, deterministic assertions, evidence handling, retry, and quarantine policy;
- changing an accepted outcome requires accountable review and visible history;
- the Trial execution plan binds scenario, application revision, runner, identity, fixture, and result schema.

M1 exit:

- both scenarios execute through the same public CLI available to a real user, and candidate code cannot rewrite their oracle.

Before Trial:

- repeatability, negative authorization behaviour, evidence privacy, cost, duration, and quarantine handling are measured.

### G10 — Offensive assurance safety

Before M1:

- a versioned rules-of-engagement contract defines ownership, exact targets, exclusions, permitted and prohibited techniques, rate and impact limits, stop conditions, evidence, cleanup, and incident handling;
- production and third-party targets are denied by default;
- deterministic tests prove that generated content cannot enable an active runner, broaden target scope, or override a kill switch;
- M1 permits only safe contract, configuration, dependency, and isolation tests;
- no active or autonomous offensive runner is enabled.

M1 exit:

- M1 authority, worker, and tool contracts preserve the prohibition and cannot reinterpret the rules of engagement.

Before Trial:

- any later active-testing Trial requires separate authorization, isolated targets, external scope enforcement, emergency termination, and independent validation.

## Gate acceptance record

Every gate acceptance SHALL record:

- gate identifier and exact scope;
- exact Git revision reviewed;
- evidence paths, commands, results, and immutable digests where applicable;
- individual delivery owner;
- accountable reviewers, outcomes, and dates;
- exclusions and evidence assigned to M1 exit or Before Trial;
- residual risks, owners, review dates, and expiry conditions;
- explicit `Accepted for F0` outcome.

A gate cannot be accepted from document presence, a passing candidate-controlled test alone, or an AI recommendation.

## F0 final acceptance

F0 is complete only when:

1. G0 through G10 are `Accepted for F0` at one candidate revision;
2. every required local and protected-CI check passes for that revision;
3. the implementation baseline is current;
4. residual risks and downstream evidence assignments have owners;
5. source migration remains traceable and the masterplan, roadmap, and accepted decisions remain reconciled;
6. individual delivery owners and accountable reviewers are named for the M1 Trial;
7. Product, Architecture, Security/Privacy, Engineering, and Operations accept the complete evidence manifest;
8. a human-approved ADR explicitly authorizes the bounded M1 implementation.

Until then, the architecture readiness register continues to deny M1 functional implementation.
