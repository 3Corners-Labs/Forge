# Bounded autonomous contributor delivery

Status: Accepted
Owner: Forge Architecture and Engineering
Version: 1.0

## Purpose

This procedure implements ADR-0011 for contributors working on the Forge
repository. It allows one authenticated human to delegate routine, reversible
delivery work for one task while preserving architecture, scope, landing,
release, production, destructive-action, policy, evidence, and budget gates.

It is not a Forge runtime delegation mechanism and does not implement Identity,
Delegation, Policy Decision, Execution Envelope, Approval Record, or Evidence
Record capabilities.

## Current grant

The first grant is
`forge.autonomy-grant.m1-bounded-delivery-v1`. It is stored at the canonical
path `.forge/autonomy/grants/forge.autonomy-grant.m1-bounded-delivery-v1.json`
and is constrained to M1 work in `3Corners-Labs/Forge` on `codex/m1-*` branches.

The grant is inactive by default. Its `authorized` lifecycle status means only
that it may be considered for activation after accountable review and protected
landing. A merged file, identifier, digest, or copied phrase does not activate
it. The grant expires at the timestamp recorded in the canonical JSON and
cannot be used after expiry or revocation.

The initial grant permits no new external paid spend. Platform-managed Codex
usage remains governed by the authenticated product boundary and is not an
external purchase authority delegated by this grant.

## Candidate-independent validation

Before asking the human to activate a task, the contributor SHALL:

1. fetch or otherwise verify the protected `origin/main` state;
2. select a full protected revision containing the grant, schema, validator,
   operating contract, and tests;
3. load `scripts/validate_autonomy_grant.py` from that same protected revision,
   not from the candidate worktree;
4. supply the exact task, future `codex/m1-*` branch, capability IDs, paths,
   requested routine actions, budgets, repository, and terminal state; and
5. require the validator result `eligible-for-human-activation`.

One reproducible validator extraction pattern is:

~~~bash
protected_revision=<full-protected-main-sha>
validator_path=$(mktemp)
git show "$protected_revision:scripts/validate_autonomy_grant.py" > "$validator_path"
python3 "$validator_path" \
  --repo-root . \
  --grant-id forge.autonomy-grant.m1-bounded-delivery-v1 \
  --protected-revision "$protected_revision" \
  --repository 3Corners-Labs/Forge \
  --task m1-iteration-4-example \
  --branch codex/m1-iteration-4-example \
  --capability forge.capability.policy \
  --path internal/delivery/example.go \
  --action edit-task-scope \
  --action run-verification
~~~

The contributor performs this validation; the human does not need to run the
command. Temporary validator material SHALL be removed through the
environment's safe temporary-file lifecycle.

The validator fails closed for unknown schema versions, malformed or
integrity-mismatched grants, non-protected revisions, candidate substitution,
wrong repository, wrong base, wrong task, wrong branch, wrong capability,
wrong path, undeclared action, expiry, revocation, budget overrun, missing
mandatory gates, or an unsupported terminal state.

## Human activation

After successful candidate-independent validation, the contributor gives the
accountable human one fully populated sentence in this form:

> Activate bounded autonomous delivery grant `<grant-id>` at protected revision
> `<full-sha>` for task `<task-id>` on branch `<branch>`, affecting capabilities
> `<capability-ids>` and paths `<paths>`, with terminal state
> `pull-request-merge-ready`.

The human activates by returning that exact sentence through the current
authenticated collaboration boundary. The activation is valid only for those
values. Silence, earlier approval, a similar phrase for another task, or an
activation bound to another revision or branch grants nothing.

The contributor reports the grant ID, protected revision, task, branch,
capabilities, paths, start time, limits, and terminal state when acknowledging
activation. Each pull request produced under the grant records the same values.

## Autonomous operating envelope

Within the exact activation, the contributor may load relevant context, create
the named feature branch, edit scoped paths, run checks, repair in-scope
failures, create ordinary commits, push without force, open and update a draft
pull request, address non-architectural review feedback, mark the pull request
ready, and monitor it until merge-ready.

Progress updates are informational. The contributor does not wait for a reply
while actions remain within the grant, limits, architecture, task intent, and
repository protections.

The mandatory human gate `protected-default-branch-landing` always remains in
force. After all current checks and reviews pass, the contributor requests a
fresh landing confirmation naming the exact pull request and current head SHA.

## Revocation and safe recovery

The accountable owner may revoke the active grant by saying:

> Revoke bounded autonomous delivery grant `<grant-id>` for task `<task-id>`.

Revocation, expiry, cancellation, terminal-state arrival, scope uncertainty,
authority or evidence loss, budget exhaustion, or another declared stop
condition prevents new delegated actions at the next enforceable boundary. The
contributor preserves the reviewable worktree and minimized evidence, reports
the last completed action, and waits for separate recovery authority.

The contributor never modifies the grant or operating contract from the active
candidate branch. Renewal or widening requires normal accountable review,
protected landing, and a later activation against the new protected revision.

## Verification evidence

The implementation is verified by:

- the closed JSON Schema in `schemas/forge/contributor/v1`;
- the dependency-free validator in `scripts/validate_autonomy_grant.py`;
- positive and negative tests in `scripts/test_autonomy_grant.py`;
- architecture fitness in `scripts/check_architecture.py`;
- explicit CODEOWNERS coverage for authority-defining files; and
- required pull-request checks for architecture, CI, vulnerability scanning,
  and CodeQL.

The negative suite covers candidate-controlled widening, integrity mutation,
missing human gates, protected-path inclusion, expiry, revocation, wrong scope,
wrong lifecycle, budget overrun, and unknown or unreachable protected
revisions.

## Rollback

Before activation, revert the implementation through the protected
pull-request path; no task authority or production state exists. During an
active task, revoke the activation first, stop new delegated actions, preserve
reviewable evidence, and then revert through a separately reviewed pull
request. Never delete evidence or weaken controls to simulate rollback.

## Limitations

- The grant delegates contributor workflow only and proves neither correctness
  nor product capability maturity.
- Authentication of the activation remains the authenticated collaboration
  boundary; repository code does not authenticate the human.
- Platform permission prompts, branch protection, CODEOWNERS, required checks,
  and external tool controls remain independently enforceable.
- The initial grant permits one contributor at a time and terminates at
  `pull-request-merge-ready`.
