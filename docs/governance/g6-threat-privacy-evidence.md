# G6 threat and privacy evidence

Status: Proposed
Owner: Forge Security and Privacy Architecture
Version: 0.1

## Objective

Provide complete pre-M1 threat, privacy, control, test-contract, evidence, and ownership coverage for the trustworthy intent-to-pull-request slice without claiming that unimplemented runtime controls have passed.

## Applicable architecture

- Accepted security architecture and first-slice threat model.
- Accepted identity and tenancy, policy and authorization, first-slice data flow, supply-chain, User Journey Assurance, and Continuous Security Testing architecture.
- Accepted ADR-0003 trust-object separation and ADR-0005 assurance boundaries.
- Accepted Developer Trial deployment-profile contract and Proposed F0 readiness plan.

The Proposed cryptographic architecture is informative only. This evidence does not accept a key hierarchy, signing provider, encryption provider, or recovery design.

## Threat-control plan

`security/threat-controls/m1.json` maps every priority threat in the accepted first-slice threat model:

- prompt injection;
- confused deputy and cross-tenant leakage;
- malicious repository and tool abuse;
- approval replay and poisoned knowledge;
- provider leakage and evidence tampering;
- denial of wallet;
- supply-chain compromise;
- candidate self-attestation;
- operator misuse;
- scenario-oracle capture and synthetic-user leakage;
- offensive-testing scope escape;
- untrusted scanner or template.

Every mapping has a named owner, controls, enforcement points, one or more negative test contracts, disclosure-safe expected results, evidence obligations, and an explicit phase. The architecture checker requires the exact threat set and rejects missing, duplicate, unowned, untestable, or unphased entries.

## Developer Trial decisions

`security/privacy/developer-trial-decisions.json` resolves the pre-M1 position for:

- identity assurance and tenant topology;
- repository credentials and worker isolation;
- context classification and retention;
- evidence integrity and telemetry minimization;
- cancellation and revocation;
- source-control destination.

Each topic declares the decision, enforcement point, M1 evidence, and Before-Trial evidence. The profile remains developer-only: it allows synthetic local identities and contract-test repository/source-control adapters, but no production credentials, support claim, isolation claim, retention promise, cryptographic assurance claim, or live publication bypass.

The file separately preserves unresolved production-shaped constraints for identity, tenancy, locations, worker technology, retention/legal hold, evidence keys, telemetry, revocation, source control, recovery, and operator access.

## Fixtures and disclosure safety

Versioned fixtures cover:

- a synthetic Restricted secret and personal-data marker using the reserved `.invalid` domain;
- a cross-scope request with a no-existence-disclosure result;
- expired authority;
- malicious untrusted repository/context/tool/model content;
- a changed digest and mutable-tag-only artefact.

The seeded marker is forbidden from telemetry, routine evidence, model providers, tool output, and pull-request descriptions. Fixtures under G9 and G10 are reused by reference rather than copied into a competing assurance source.

## Evidence phases

Before-M1 checks validate the structure and integrity of supply-chain, independent-runner, journey-oracle, offensive-scope, and disabled-scanner contracts. M1-exit contracts cover cross-scope authority, context, secrets, cancellation, revocation, evidence, budgets, tools, workers, and publication. Before-Trial contracts cover integrated tenant isolation, operator access, deletion, cleanup, restore, telemetry, and evidence disclosure.

Only the Before-M1 structural evidence is represented as passing in F0.

## Verification commands

~~~text
python3 scripts/check_architecture.py
task verify
~~~

## Review record

| Field | Value |
|---|---|
| Gate | G6 Threat and privacy readiness |
| Gate status | Accepted for F0 |
| Candidate revision | 7ff31e3 |
| Delivery owner | Sijbren Beukenkamp, temporary |
| Accountable reviewer | Sijbren Beukenkamp under ADR-0006 and ADR-0008 |
| Review outcome | Accepted for F0 on 2026-07-21 |
| Required decision | Complete |
| Residual risks | Runtime controls are not implemented; cryptographic design remains Proposed; production constraints remain unresolved; one person fills multiple review roles |

Sijbren Beukenkamp accepted the pre-M1 G6 evidence at revision `7ff31e3` on 2026-07-21. Acceptance does not promote cryptography, a production profile, or any runtime capability.
