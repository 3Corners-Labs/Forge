# G9 user intent assurance evidence

Status: Proposed
Owner: Forge Product, Experience, and Quality Architecture
Version: 0.1

## Objective

Define independently owned, reproducible allowed and denied M1 journey oracles through the public Forge CLI without claiming that the M1 delivery runner exists or that either journey has executed.

## Applicable architecture

- Accepted ADR-0005 separates User Journey Assurance from implementation and generic verification.
- Accepted ADR-0009 establishes `forge delivery start` as the M1 public CLI boundary.
- Accepted first-slice architecture, capability contracts, and data-flow analysis define the trust path, normalized failures, side-effect limits, and privacy gates.
- Accepted User Journey Assurance architecture defines scenario authority, evidence, retry, and quarantine requirements.
- Proposed F0 readiness plan assigns execution evidence to M1 exit and measured quality evidence to Before Trial.

## Candidate scenarios

| Scenario | Role | Required result | Lifecycle |
|---|---|---|---|
| `forge.journey.m1.allowed-intent-capture` | Permitted synthetic developer | Capture exactly one bounded intent and stop before execution or publication | Reviewed |
| `forge.journey.m1.denied-intent-capture` | Unauthorized synthetic observer | Return `forbidden` before repository discovery with no existence disclosure or side effect | Reviewed |

Both scenarios invoke `forge delivery start` and bind explicit persona, entitlements, denied capabilities, Developer Trial profile, synthetic fixture, steps, deterministic expected outcomes, forbidden outcomes, redacted evidence, single-attempt retry, and accountable quarantine handling.

The scenario owner is Forge Product, Experience, and Quality Architecture. The execution-plan implementation owner is Forge Engineering. Repository ownership requires Architecture and Engineering review for changes under `assurance/`; candidate code cannot silently treat its observed behaviour as the accepted oracle.

## Execution plan

`assurance/journeys/m1/trial-execution-plan.json` binds:

- the public CLI entry point;
- each scenario and synthetic identity profile;
- each fixture by SHA-256 digest;
- an immutable application-artifact digest at execution;
- a non-candidate-controlled immutable runner digest;
- the journey-result schema;
- separate Before M1, M1-exit, and Before-Trial evidence.

The plan state is `not-executable-before-m1`. The current iteration adds no journey runner, hidden interface, repository access, worker, model, tool, source-control side effect, or bypass.

## Structural verification

`check_journey_contracts` rejects:

- a missing or extra M1 scenario;
- implementation ownership of the oracle;
- a non-CLI or non-ADR-0009 entry point;
- use of a profile other than Developer Trial;
- missing fixtures or changed fixture digests;
- non-deterministic M1 outcomes, missing forbidden outcomes, retry masking, or incomplete quarantine and redaction policies;
- identical allowed and denied roles;
- a denial that does not occur before repository discovery;
- an unbound application, runner, identity, fixture, or result contract;
- enabling the runner during F0.

The required architecture workflow now watches `assurance/**`, and CODEOWNERS makes scenario changes visible to architecture and engineering reviewers.

## Verification commands

~~~text
python3 scripts/check_architecture.py
task verify
~~~

## Deferred evidence

M1 exit SHALL execute both accepted scenarios through the implemented public CLI and bind immutable result records. Before Trial SHALL measure repeatability, negative authorization, evidence privacy, duration, cost, and quarantine handling. This F0 evidence does not claim either phase has passed.

## Review record

| Field | Value |
|---|---|
| Gate | G9 User intent assurance |
| Gate status | Ready for review |
| Candidate revision | Commit containing this record |
| Delivery owner | Sijbren Beukenkamp, temporary |
| Accountable reviewer | Sijbren Beukenkamp under ADR-0006 and ADR-0008; review pending |
| Required decision | Accept or reject both reviewed scenario oracles, ADR-0009, and pre-M1 G9 evidence at the exact committed revision |
| Residual risks | The runner and public command are not implemented; one person fills multiple review roles; clean protected-CI evidence remains pending |

G9 SHALL remain Ready for review until the accountable human promotes both scenario lifecycles to `accepted` and records `Accepted for F0` at the reviewed revision.
