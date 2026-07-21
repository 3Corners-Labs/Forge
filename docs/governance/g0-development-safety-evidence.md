# G0 development safety evidence

Status: Proposed
Owner: Forge Engineering and Platform Engineering
Version: 0.1

## Objective

Demonstrate that the existing local bootstrap is explicitly Developer Trial only, exposes services only on loopback, and refuses non-developer configuration before invoking Docker.

This evidence covers the pre-M1 G0 scope only. It does not support a production profile, prove container isolation, or authorize M1 functional implementation.

## Applicable architecture

- Accepted engineering constitution sections 4, 10, 13, 16, 18, and 19.
- Accepted ADR-0004 deployment-profile composition.
- Accepted security and deployment architecture.
- Accepted readiness register and Proposed F0 readiness evidence plan.

## Implemented controls

| Control | Enforcement | Evidence |
|---|---|---|
| Runtime mode is compose | Config.ValidateDeveloperBootstrap | Configuration unit tests |
| Runtime profile is developer | Config.ValidateDeveloperBootstrap | Non-developer profile matrix and production-shaped fixture |
| Ungoverned intelligence is disabled | Config.ValidateDeveloperBootstrap | Negative intelligence.enabled=true test |
| forge up validates before side effects | Shared CLI action path | CLI side-effect denial test |
| forge down validates before side effects | Shared CLI action path | CLI side-effect denial test |
| Missing configuration fails before side effects | Shared CLI action path | CLI missing-file denial test |
| Published service ports use loopback | Architecture fitness check | Compose port scan |
| Development risks remain explicit | Architecture fitness check | Static credential, Keycloak dev mode, OpenBao dev mode, debug exporter, and local-volume markers |
| Production-shaped configuration is rejected | CLI test fixture | internal/cli/testdata/production-shaped.yaml |

## Developer limitations

The Compose environment contains static development credentials, Keycloak development mode, OpenBao development mode and root token, MinIO and Grafana development credentials, a debug telemetry exporter, version-tagged images, and local named volumes. These properties are permitted only for the Developer Trial bootstrap and are prohibited as production evidence.

Direct invocation of deploy/compose/docker-compose.yml is an explicit developer operation outside the Forge CLI profile contract. It does not create or simulate a supported non-developer profile.

The bootstrap does not implement authenticated Forge users, tenant isolation, workload identity, governed intelligence, worker isolation, production secret handling, durable evidence, backup/restore guarantees, or a production support envelope.

## Operation and recovery

- Start through forge up with runtime.mode: compose, runtime.profile: developer, and intelligence.enabled: false.
- Stop through forge down using the same validated configuration.
- A failed start is recovered by correcting configuration or prerequisites and retrying; the CLI does not silently change profiles.
- forge down preserves named development volumes. Volume deletion is a separate destructive operator decision and is not performed by Forge.
- If Docker reports an ambiguous local result, the operator inspects docker compose -f deploy/compose/docker-compose.yml ps before retrying.

## Verification commands

~~~text
go test ./internal/config ./internal/cli
go test ./... -race
python3 scripts/check_architecture.py
task verify
~~~

## Evidence assignment

Before M1 evidence in this document is implemented and reproducible. M1 must later prove that its public CLI and worker paths cannot bypass the Developer profile restriction. Before Trial, the complete Developer profile conformance suite must pass on a clean supported host.

## Review record

| Field | Value |
|---|---|
| Gate | G0 Development safety |
| Gate status | Ready for review |
| Candidate revision | 9db0afb |
| Delivery owner | Sijbren Beukenkamp, temporary |
| Accountable reviewer | Sijbren Beukenkamp under ADR-0006 and ADR-0008; review pending |
| Residual risks | Direct Compose invocation remains a developer operation; Docker/container isolation is not claimed; clean-host conformance remains Before Trial evidence |
| Required decision | Accept or reject G0 at candidate revision 9db0afb |

G0 SHALL remain Ready for review until the accountable human review records Accepted for F0.
