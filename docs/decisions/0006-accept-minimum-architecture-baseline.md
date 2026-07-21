# ADR-0006: Accept the minimum architecture baseline

Status: Accepted
Owner: Forge Architecture Board
Date: 2026-07-21

## Context

Forge cannot begin functional capability implementation without a minimum authoritative baseline. The repository prepared a bounded acceptance packet covering constitutional guardrails, lifecycle semantics, the architecture metamodel, five foundational decisions, the intent-to-pull-request slice, its capability contracts, and its threat, data-flow, and privacy boundaries.

The exact reviewed Proposed package is Git commit `f2809ef0e0da698db74e563251286cc159278865`.

Sijbren Beukenkamp explicitly accepted the package on 2026-07-21 and stated that, temporarily, he represents all required accountable roles. This temporary role concentration is recorded rather than treated as implicit separation of duties.

## Decision

Accept the following documents and decisions at the reviewed revision:

- `docs/foundation/constitution.md` version 1.1;
- `docs/governance/document-lifecycle.md`;
- `docs/architecture/metamodel.md` version 0.1;
- ADR-0001 through ADR-0005;
- `docs/architecture/vertical-slice.md` version 0.1;
- `docs/architecture/first-slice-capability-contracts.md` version 0.1;
- `docs/architecture/threat-model.md` version 0.1;
- `docs/architecture/first-slice-data-flow.md` version 0.1.

This acceptance authorizes architecture, contract-test, and implementation work for the bounded trustworthy intent-to-pull-request slice only after the remaining readiness conditions applicable to that work are satisfied. It does not accept the full masterplan or capability catalogue, claim that a capability is implemented, support a production deployment profile, or authorize merge or production release automation.

## Accountable review record

| Role | Reviewer | Outcome | Date |
|---|---|---|---|
| Product | Sijbren Beukenkamp | Accept | 2026-07-21 |
| Architecture | Sijbren Beukenkamp | Accept | 2026-07-21 |
| Security and Privacy | Sijbren Beukenkamp | Accept | 2026-07-21 |
| Engineering | Sijbren Beukenkamp | Accept | 2026-07-21 |
| Operations | Sijbren Beukenkamp | Accept | 2026-07-21 |

The organization permits this temporary multi-role representation for the bootstrap phase. The review record SHALL be updated or superseded when separate accountable owners are appointed.

## Consequences

- The accepted constitution now has the highest repository architecture authority under `/AGENTS.md`.
- ADR-0001 through ADR-0005 govern their declared decisions.
- The accepted metamodel, first-slice boundary, contracts, threat model, and data-flow analysis govern work within that slice.
- Broader architecture remains Proposed and cannot silently expand implementation scope.
- Readiness gates requiring implemented controls, operational evidence, or profile validation remain open or in progress.
- One person currently represents all review roles, reducing organizational challenge independence. Candidate implementation self-attestation and self-approval prohibitions remain unchanged.

## Validation

- Architecture CI SHALL verify document statuses, navigation, links, contract schemas, source integrity, and developer-profile constraints.
- Implementation work SHALL demonstrate the conformance evidence named by the accepted slice documents.
- The readiness register SHALL continue to distinguish accepted design from validated implementation.
- This decision SHALL be reviewed when accountable roles are separated, the first Trial begins, the slice boundary changes, or a production-shaped profile is proposed.

## Supersession

None. A later decision may replace the temporary reviewer-role arrangement without reopening unrelated accepted architecture.
