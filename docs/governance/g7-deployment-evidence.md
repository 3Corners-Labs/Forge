# G7 deployment evidence

Status: Proposed
Owner: Forge Platform Engineering
Version: 0.1

## Objective

Define explicit conformance requirements for the Developer Trial profile and one production-shaped connected self-hosted reference profile without claiming production implementation or support.

## Applicable architecture

- Accepted ADR-0001 modular platform.
- Accepted ADR-0004 composable deployment profiles.
- Accepted architecture metamodel, deployment architecture, security architecture, quality and operations architecture, and first-slice validation envelope.
- Accepted readiness register and Proposed F0 readiness evidence plan.

## Profile artefacts

| Artefact | Maturity | Purpose |
|---|---|---|
| profiles/developer-trial.profile.json | Trial profile contract; implementation incomplete | Complete M1 Developer Trial envelope and phase-specific conformance plan |
| profiles/connected-self-hosted-reference.profile.json | Research | Production-shaped customer-operated target and test plan without support claims |
| schemas/forge/deployment/v1/profile.schema.json | Research contract candidate | Closed machine-readable deployment-profile shape |

Both profiles resolve connectivity, tenancy, operations, assurance, and provider-ownership dimensions. They use the same logical contract while declaring different operators, assurance expectations, provider ownership, implementation state, and unsupported claims.

## Required control ownership

Each profile assigns ownership, assumptions, and Before M1, M1 exit, and Before Trial tests for:

- identity;
- storage;
- repository access;
- worker isolation;
- network policy;
- secrets;
- audit evidence;
- telemetry;
- backup and restore;
- upgrade and rollback;
- capacity and recovery.

The manifests describe logical deployment requirements. They do not select Kubernetes, a cloud, an identity provider, a secret provider, a database topology, or a production worker-isolation technology.

## Executable conformance plan

The architecture checker validates:

- the closed profile schema;
- the exact initial profile set and maturity;
- unique profile and test identifiers;
- every required dimension;
- every required control owner and assumption;
- non-empty test assignments for all three evidence phases;
- commands and expected results for profile-level conformance tests;
- explicit unsupported claims;
- absence of a Supported production profile.

The Developer Trial before-M1 tests use current commands. M1-exit and Before-Trial commands are declared test contracts and must not be represented as passing until their corresponding functionality exists.

## Unsupported claims

This evidence does not support any claim that:

- a production deployment exists;
- the connected self-hosted reference is Supported;
- a production SLA, RTO, RPO, capacity envelope, residency promise, or compatibility window is accepted;
- provider portability, SaaS conformance, restricted connectivity, or disconnected operation has been validated;
- the current Compose bootstrap implements the complete Developer Trial profile.

## M1-exit candidate evidence

The Developer Trial Compose profile now supplies the PostgreSQL and MinIO
services exercised by the operational public journey. Worker and verifier
images build reproducibly from separate deployment units and are invoked by
immutable local image ID with explicit resource and network restrictions.
This is a local Developer Trial composition only; it does not complete the
profile's backup, restore, upgrade, rollback, capacity, recovery, identity,
secret-provider, or connected self-hosted conformance requirements.

## Verification commands

~~~text
python3 scripts/check_architecture.py
task verify
~~~

## Review record

| Field | Value |
|---|---|
| Gate | G7 Deployment evidence |
| Gate status | Accepted for F0 |
| Candidate revision | e3c6355 |
| Delivery owner | Sijbren Beukenkamp, temporary |
| Accountable reviewer | Sijbren Beukenkamp under ADR-0006 and ADR-0008 |
| Review outcome | Accepted for F0 on 2026-07-21 |
| Residual risks | Profile tests beyond the existing bootstrap are planned; topology and provider selections remain undecided; role concentration reduces independent challenge |
| Required decision | Complete |

Sijbren Beukenkamp accepted the pre-M1 G7 evidence at revision `e3c6355` on 2026-07-21. M1-exit and Before-Trial evidence remains assigned to its declared later phase. A later topology or provider selection requires normal architecture review and an ADR when it materially constrains Forge.
