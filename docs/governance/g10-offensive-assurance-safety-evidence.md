# G10 offensive assurance safety evidence

Status: Proposed
Owner: Forge Security and Quality Architecture
Version: 0.1

## Objective

Make the pre-M1 security-assurance boundary machine-readable and fail closed while keeping every active and autonomous offensive runner disabled.

## Applicable architecture

- Accepted ADR-0005 separates Continuous Security Testing and treats autonomous offensive components as untrusted.
- Accepted Continuous Security Testing architecture requires externally enforced scope, limits, kill switch, evidence, cleanup, and incident handling.
- Accepted threat model identifies offensive scope escape and untrusted scanner/template threats.
- Accepted security architecture and F0 readiness plan deny production and third-party targeting by default.

## Contract and safe instance

`schemas/forge/security/v1/rules-of-engagement.schema.json` defines a closed contract covering:

- requester, approver, operator, and emergency contact;
- purpose and threat references;
- environment, ownership proof, authorization window, and target/configuration digests;
- exact targets, exclusions, and production/third-party denial;
- permitted and prohibited techniques;
- rate, concurrency, duration, cost, payload, and impact limits;
- kill switch, fail-closed behaviour, stop conditions, and termination bound;
- external deterministic scope enforcement and generated-content restrictions;
- evidence classification, recipients, retention, redaction, and coverage disclosure;
- cleanup and incident handling.

`security/assurance/m1.rules-of-engagement.json` is an authorized, contract-only M1 record. Its target set is one repository-local, non-network fixture. Its runner configuration is disabled, non-autonomous, externally scoped, unable to accept generated configuration, and held behind an engaged kill switch. Both referenced configurations are SHA-256 bound.

The permitted techniques are limited to schema validation, configuration analysis, dependency metadata analysis, and isolation contract tests. Active network requests, request rate, cost units, and payload bytes are all zero.

## Negative conformance tests

The deterministic architecture checker applies representative malicious or unsafe overrides and requires these failures:

| Fixture | Required denial |
|---|---|
| `production-target.json` | Production target |
| `third-party-target.json` | Third-party target |
| `scope-expansion.json` | Target outside the digest-bound set |
| `generated-activation.json` | Generated content enabling a runner |
| `kill-switch-override.json` | Disabled kill switch or fail-open behaviour |

The negative fixtures are inputs to the external checker; they are never merged into an executable configuration. Changes under `security/assurance/**` run the required architecture workflow and require Architecture and Engineering ownership review.

## Verification commands

~~~text
python3 scripts/check_architecture.py
task verify
~~~

## Deferred evidence

The M1-exit candidate gives neither the public CLI nor worker, verifier, model,
Tool, or source-control paths an offensive-assurance action. Both runtime
containers have no network and accept closed plan-controlled inputs; the
architecture checker continues to require the disabled contract-only rules of
engagement. Before Trial, any proposal for active security testing requires
separate human authorization, isolated resettable targets, external network
enforcement, emergency termination testing, and independent validation. No
such Trial is authorized here.

## Review record

| Field | Value |
|---|---|
| Gate | G10 Offensive assurance safety |
| Gate status | Accepted for F0 |
| Candidate revision | 0784c62 |
| Delivery owner | Sijbren Beukenkamp, temporary |
| Accountable reviewer | Sijbren Beukenkamp under ADR-0006 and ADR-0008 |
| Review outcome | Accepted for F0 on 2026-07-21 |
| Required decision | Complete |
| Residual risks | Runtime enforcement does not yet exist; protected-CI configuration still requires G8 review; one person fills multiple review roles |

Sijbren Beukenkamp accepted the rules of engagement and pre-M1 G10 evidence at revision `0784c62` on 2026-07-21. Acceptance does not authorize an active or autonomous security-testing runner.
