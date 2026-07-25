# G8 architecture fitness evidence

Status: Proposed
Owner: Forge Architecture and Engineering
Version: 0.1

## Objective

Make Forge architecture, dependency, isolation, assurance, documentation, and source-integrity rules executable and required in CI without claiming that the future worker boundary is implemented.

## Applicable architecture

- Accepted constitution and repository operating contract.
- Accepted ADR-0001 modular platform and isolated-worker decision.
- Accepted architecture metamodel and first-slice capability contracts.
- Accepted threat model and ADR-0005 assurance boundaries.
- Proposed F0 readiness plan.

## Module and isolation rules

`architecture/module-boundaries.json` records the current source modules, allowed internal dependencies, deployment-unit mapping, maturity, and planned worker modules. Planned paths remain explicitly Research and absent; their presence while marked planned fails validation.

The executable rules deny:

- control-plane imports of worker implementation;
- worker imports of control-plane implementation;
- shared contracts importing either runtime;
- undeclared or unknown internal module dependencies.

Representative fixtures under `architecture/testdata/` prove failures for control-to-worker, worker-to-control, and unknown-module imports. The checker also scans imports in every currently implemented Go module.

## Required CI controls

The architecture workflow:

- runs on every pull request and every push to `main`, without path filters;
- has read-only repository permissions and a ten-minute timeout;
- runs the dependency-free architecture checker, metadata validation, Markdown-link validation, and strict documentation build;
- is covered, together with the checker and architecture manifests, by explicit Architecture and Engineering CODEOWNERS rules.

The checker verifies these workflow controls and the main CI commands for formatting, vetting, race-enabled tests with coverage, and a reproducible build. It also executes trust-schema, profile, journey-contract, offensive-safety, source-integrity, documentation, capability-ID, navigation, link, and developer-profile checks.

## Failure recovery

When a required fitness check fails:

1. reproduce with `python3 scripts/check_architecture.py` and `task verify` from the exact candidate revision;
2. identify the named rule and offending artefact; do not disable, path-filter, rename, or weaken the check to obtain a pass;
3. correct the dependency, ownership, contract, documentation, or fixture and rerun the complete suite;
4. if the accepted architecture must change, stop implementation and obtain a reviewed ADR before changing the oracle;
5. if CI infrastructure fails without a candidate defect, preserve the failed run, retry only within normal GitHub controls, and escalate to Forge Engineering and Architecture;
6. use a reviewed revert of the offending candidate as the safe rollback. Do not bypass the required check or use administrative merge privileges.

## External enforcement

An initial read on 2026-07-21 returned `404 Branch not protected`. Sijbren Beukenkamp then explicitly authorized the required protection, and the repository owner applied it through the GitHub branch-protection API.

The verified `main` settings require strict, up-to-date `verify` and `validate` checks; one approving review with code-owner review; resolved conversations; enforcement for administrators; and denial of force pushes and deletion. These settings use the normal protected path and do not authorize administrative bypass.

PR #10 exposed that the inherited CODEOWNERS file referenced organization teams that do not exist. The candidate replaces those aliases with the two actual writable repository collaborators, `cyberbeuk` and `3c-unified-git`, and the architecture checker now rejects reintroduction of unverified `@3corners/` team aliases. This makes required ownership effective rather than documentary.

## Verification commands

~~~text
python3 scripts/check_architecture.py
task verify
gh api repos/3corners-labs/forge/branches/main/protection
~~~

## Deferred evidence

The M1-exit candidate runs the dependency and isolation rules against the
implemented control plane, shared contracts, worker, and verifier modules.
Local architecture checks pass without a capability-boundary exception. The
real Docker journey additionally proves separate deployment images and no
control-plane implementation inside either runtime. Before Trial SHALL pass
all required checks from a clean checkout in the protected workflow.

## Review record

| Field | Value |
|---|---|
| Gate | G8 Architecture fitness functions |
| Gate status | Accepted for F0 |
| Candidate revision | ea8c94f |
| Delivery owner | Sijbren Beukenkamp, temporary |
| Accountable reviewer | Sijbren Beukenkamp under ADR-0006 and ADR-0008 |
| Review outcome | Corrective evidence accepted for F0 on 2026-07-21 |
| External evidence | GitHub `main` protection verified on 2026-07-21; strict required checks `verify` and `validate` |
| Required decision | Complete |
| Residual risks | Worker modules do not yet exist; repository administrators retain platform-level ownership; the second collaborator's organizational independence is not established |

Sijbren Beukenkamp accepted the corrected pre-M1 G8 evidence at revision `ea8c94f` on 2026-07-21. M1-exit and Before-Trial evidence remains assigned to its declared later phase.
