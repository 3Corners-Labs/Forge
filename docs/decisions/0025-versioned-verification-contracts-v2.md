# ADR-0025: Introduce versioned verification contracts v2

Status: Accepted
Owner: Forge Architecture, Verification, and Engineering
Date: 2026-07-25

## Context

Provenance recording and the self-certification invariant required by OWASP
AISVS AC.8.1 added two mandatory fields to the verification contracts:
`producedBy` on the Assurance Plan and `verifiedBy` on the Verification Result.

They were added in place, to `forge.assurance-plan.v1` and
`forge.verification-result.v1`. Both are closed accepted schemas.

[ADR-0016](0016-scope-bound-approval-record-v2.md) already rejected exactly this
approach for the Approval Record:

> Adding required fields to `forge.approval-record.v1` would silently break a
> closed accepted schema.

The break is not theoretical. `DecodePlan` and `DecodeResult` reject unknown
fields. The verifier is a separate deployment unit built as an immutable image,
so a verifier image predating the change cannot decode a plan carrying
`producedBy` at all. An operator running a stale image sees a decode failure
rather than a version mismatch. Contract versioning exists to make that
difference visible.

The defect was found by an audit of the decision register against current
intent, requested after the change had already been committed.

## Decision

Introduce `forge.assurance-plan.v2` and `forge.verification-result.v2`, and
restore the v1 schemas to their accepted shape.

The v2 contracts carry every v1 field plus:

- `producedBy` on the Assurance Plan: the workload principal Policy issued for
  the worker phase that produced the candidate. Mandatory. Server-derived; a
  candidate cannot assert it.
- `verifiedBy` on the Verification Result: the principal that evaluated the
  plan. Mandatory. Set by the verifier from its own identity.

A result whose `verifiedBy` equals the plan's `producedBy` SHALL be refused.
That rule remains a contract invariant rather than a control, because a
verification profile may omit a control but cannot omit an invariant.

v1 remains readable for historical validation. It SHALL NOT satisfy a newly
implemented verification gate, mirroring the Approval Record v1 treatment.

Both schema files are retained under `schemas/forge/trust/v1/`, following the
`approval-record-v2.schema.json` precedent: the directory names the trust
contract family, not the version of each contract in it.

## Clarifications to earlier decisions

Two ambiguities surfaced in the same audit. Neither changes a decision; both
prevent a decision from being misread.

### The contributor harness is not bounded autonomous delivery

[ADR-0024](0024-rebaseline-forge-as-bounded-means.md) section 2 recognises the
contributor harness as the primary delivery mechanism. That refers to the
operating contract, context router, decision register, fitness checks, and an
external coding assistant working **under human review at every decision**.

It does **not** refer to, activate, or make routine the bounded autonomous
delivery mode established by
[ADR-0011](0011-bounded-autonomous-contributor-delivery.md). That mode remains
inactive by default and still requires a validated grant and a per-task human
activation. Recognising the harness as the delivery mechanism grants no
autonomy.

### ADR-0007 step 4 is Vision, not plan

[ADR-0007](0007-sequence-self-building-before-full-advisory.md) lists the general
Software System Twin and advisory product as step 4. ADR-0024 demoted that
milestone to evidence label Vision. Step 4 records intended direction and is not
committed scope.

## Alternatives considered

- **Amend v1 in place and record the exception:** rejected because ADR-0016
  settled the question for the same class of artefact, and because the decode
  failure between deployment units is real rather than hypothetical.
- **Bump only the Assurance Plan:** rejected because the result carries the
  matching new field and the two are validated against each other.
- **Move the invariant into a control instead of the contract:** rejected for
  the reason recorded above; a profile can omit a control.
- **Defer until Workplace delivery begins:** rejected because carrying a broken
  trust contract into a new repository is the worst moment to move technical
  debt.

## Consequences

- Any stored v1 Assurance Plan or Verification Result no longer satisfies a
  current verification gate. Developer Trial state is local and may be discarded.
- The verifier image SHALL be rebuilt before operational mode is used, and a
  stale image now fails with a version mismatch rather than a decode error.
- The trust schema set check gains two entries.
- Contract changes of this class are cheap while no external consumer exists and
  expensive afterwards. Workplace Orchestrator will be the first external
  consumer, which is why this is settled before it starts.

## Validation

Evidence SHALL prove:

- a v1 plan or result is refused by a current gate while remaining decodable for
  historical validation;
- a plan without `producedBy` and a result without `verifiedBy` fail closed;
- a result whose verifier equals the producer is refused;
- the verifier and control plane agree on the contract version, and a mismatch
  is reported as such;
- the trust schema set, JSON validity, and closed-object checks pass.

## Accountable review

Sijbren Beukenkamp accepted this decision on 2026-07-25 through the
authenticated collaboration boundary, after the drafting contributor reported
the defect in its own work and proposed this remedy. The acceptance covered the
v2 contracts and the two clarifications as described above.

The drafting contributor recorded the acceptance; it did not make it.

## Supersession

This decision specializes ADR-0013 and ADR-0017 by versioning the artefacts they
govern, and follows the precedent set by ADR-0016.

It supersedes only the in-place amendment of `forge.assurance-plan.v1` and
`forge.verification-result.v1` introduced earlier on the same branch. It does not
supersede capability ownership, the blocking-check rule, the applicability
rule, isolated execution, or any approval or publication restriction.
