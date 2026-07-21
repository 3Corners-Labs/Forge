# F0 final acceptance manifest

Status: Accepted
Owner: Forge Architecture Board
Version: 0.1

## Purpose

Bind the complete F0 evidence set to one immutable repository baseline and record the accountable decision authorizing the bounded M1 implementation.

## Decision

F0 is Accepted. G0 through G10 are Accepted for F0, required local and protected-CI checks passed, the implementation baseline remains current, and residual risks and deferred evidence have owners.

Together with Accepted [ADR-0010](../decisions/0010-authorize-bounded-m1-implementation.md), this manifest authorizes the bounded M1 implementation. That authority becomes operationally effective when:

1. this manifest and ADR-0010 remain Accepted together;
2. the acceptance records continue to bind reviewed revision `773312e`;
3. the Accepted revision reaches protected `main` through the required checks and review path.

## Immutable F0 baseline

| Field | Value |
|---|---|
| Repository | `3Corners-Labs/Forge` |
| Protected branch | `main` |
| F0 baseline commit | `93efcfcae8e946310ebd07294d744f38a3083b0e` |
| F0 baseline tree | `eeea9e5e0d95534cebc8ac253672d00b818027c8` |
| Source PR | `https://github.com/3Corners-Labs/Forge/pull/10` |
| Reviewed PR head | `3aed6e8ba00b6fc2f5df7e6aa407ca553bb631ec` |
| Authorization review revision | `773312e` |
| Merge method | Squash through protected branch |
| Merge time | 2026-07-21T12:25:48Z |

The baseline commit is the single candidate revision containing the complete accepted gate set. Earlier gate revisions preserve the review history and are not competing implementation baselines.

## Gate evidence manifest

| Gate | Accepted scope | Evidence | Outcome |
|---|---|---|---|
| G0 | Developer bootstrap safety | [G0 evidence](g0-development-safety-evidence.md), candidate `9db0afb` | Accepted for F0 |
| G1 | Governance authority and lifecycle | ADR-0006, ADR-0008, accepted authority set | Accepted for F0 |
| G2 | Architecture metamodel | Accepted metamodel and ADR-0002 | Accepted for F0 |
| G3 | Bounded first slice | Accepted vertical slice and ADR-0007 | Accepted for F0 |
| G4 | Non-overlapping trust contracts | Trust schemas and ADR-0003 | Accepted for F0 |
| G5 | First-slice capability contracts | Accepted capability contracts | Accepted for F0 |
| G6 | Threat and privacy readiness | [G6 evidence](g6-threat-privacy-evidence.md), candidate `7ff31e3` | Accepted for F0 |
| G7 | Deployment evidence | [G7 evidence](g7-deployment-evidence.md), candidate `e3c6355` | Accepted for F0 |
| G8 | Required architecture fitness | [G8 corrected evidence](g8-architecture-fitness-evidence.md), candidate `ea8c94f` | Accepted for F0 |
| G9 | User intent assurance | [G9 evidence](g9-user-intent-assurance-evidence.md), candidate `ee0eff9` | Accepted for F0 |
| G10 | Offensive assurance safety | [G10 evidence](g10-offensive-assurance-safety-evidence.md), candidate `0784c62` | Accepted for F0 |

Every gate remains limited to its Before-M1 evidence. Acceptance does not claim that M1-exit or Before-Trial tests have passed.

## Verification evidence

The F0 candidate passed:

~~~text
python3 scripts/check_architecture.py
task verify
go test ./... -race -coverprofile=<temporary-path>
docker compose -f deploy/compose/docker-compose.yml config --quiet
~~~

Protected PR #10 recorded required `validate` and `verify` checks as passing, one effective code-owner approval, no unresolved review thread, no conflict, and an exact-head protected squash merge. At acceptance review, the commands SHALL be rerun for the ADR/manifest candidate; prior success is not treated as permanent evidence.

## Authorized M1 envelope

Subject to Accepted ADR-0010, M1 is limited to:

- the public `forge delivery start` boundary from ADR-0009;
- one Developer Trial tenant, workspace, project, and Forge repository binding per run;
- the accepted intent-to-pull-request path and participating capability contracts;
- a modular control plane and separately enforced worker boundary;
- deterministic implementation by default, with any model adapter requiring separate provider eligibility and accountable approval;
- one source-control adapter or contract double behind the accepted publication contract;
- the accepted allowed and denied G9 scenarios;
- evidence needed to evaluate the declared M1 exit criteria.

M1 authority is implementation authority, not production authority, release authority, capability-maturity promotion, or evidence that the target behaviour already exists.

## Explicit exclusions

The following remain unauthorized:

- merge, release, production deployment, or production-policy mutation;
- production credentials, production data, customer data, or production support claims;
- active or autonomous offensive security testing;
- broad multi-agent workforce coordination or self-expanding autonomy;
- general model routing, unrestricted tools, or unreviewed external providers;
- billing, certification, accreditation, marketplace, and SaaS control-plane functionality;
- Workplace-specific behaviour;
- service extraction, new deployment topology, or capability-boundary changes without an Accepted ADR.

## Deferred evidence

### Required at M1 exit

- implemented public CLI allowed and denied journeys;
- identity, scope, policy, envelope, context, worker, secret, cancellation, revocation, digest, approval, evidence, and publication negative tests;
- dependency and control-plane/worker isolation checks against the implemented modules;
- idempotent ambiguous-publication reconciliation;
- interrupted-run recovery without duplicate publication;
- evidence integrity and evidence-unavailability behaviour;
- implementation-baseline update with exact source and evidence.

### Required before Trial

- clean-host Developer profile conformance;
- integrated isolation, deletion, cleanup, restore, telemetry, and evidence-disclosure tests;
- repeatability, negative authorization, privacy, cost, duration, and quarantine measurements;
- approved identity, retention, provider, evidence-integrity, recovery, compatibility, and support decisions for the claimed Trial envelope.

## Residual risks

| Risk | Owner | Review trigger or expiry |
|---|---|---|
| One person temporarily represents Product, Architecture, Security/Privacy, Engineering, and Operations | Sijbren Beukenkamp | Review before M1 exit and replace with distinct accountable roles before any production-shaped claim |
| Worker isolation is a contract, not an implemented control | Forge Engineering and Security Architecture | Stop before executing untrusted repository content; resolve during M1 |
| Cryptographic architecture and evidence signing remain Proposed | Forge Security Architecture | Resolve before claiming cryptographic assurance or production integrity |
| Identity, source-control, secret, and persistence providers remain unselected | Respective capability owners | Provider choice requires eligibility, contract evidence, and an ADR when materially constraining |
| Developer Trial uses development dependencies and credentials | Forge Platform Engineering | Never use outside the Developer profile; revalidate at M1 exit and Before Trial |
| Production topology, retention, recovery, capacity, and support envelopes remain unresolved | Forge Architecture, Privacy, and Operations | Must be resolved before any production-shaped profile becomes Supported |
| CODEOWNERS uses two repository collaborator identities; organizational independence is not established | Forge Architecture Board | Review before M1 exit and whenever collaborator access changes |

## Accountable review record

| Role | Reviewer | Outcome | Date |
|---|---|---|---|
| Delivery owner | Sijbren Beukenkamp | Accept | 2026-07-21 |
| Product | Sijbren Beukenkamp, temporary | Accept | 2026-07-21 |
| Architecture | Sijbren Beukenkamp, temporary | Accept | 2026-07-21 |
| Security and Privacy | Sijbren Beukenkamp, temporary | Accept | 2026-07-21 |
| Engineering | Sijbren Beukenkamp, temporary | Accept | 2026-07-21 |
| Operations | Sijbren Beukenkamp, temporary | Accept | 2026-07-21 |

Sijbren Beukenkamp accepted this manifest and ADR-0010 as drafted at revision `773312e` on 2026-07-21 and stated: `F0 Accepted; bounded M1 implementation authorized under ADR-0010`. The AI contributor prepared the manifest and verification evidence but did not provide accountable acceptance.

## Outcome

**F0 Accepted; bounded M1 implementation authorized under ADR-0010.** The authority becomes operationally effective when this Accepted packet reaches protected `main`. All explicit exclusions, stop conditions, residual risks, and deferred evidence remain binding.
