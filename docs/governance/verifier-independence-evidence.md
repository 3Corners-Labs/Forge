# Verifier independence: attack results

Status: Draft
Owner: Forge Engineering and Architecture
Verified: 2026-07-25

## Question

OWASP AISVS 1.0 Appendix C, requirement AC.8.1, states that agents cannot
approve, merge, sign, or deploy artifacts that they themselves generated. Market
research found no shipped agent product enforcing this structurally.

Forge claims a mandatory deciding independent check. This exercise
attacks that claim rather than restating it.

## Result in one line

**The boundary holds structurally, and is substantively empty for the kind of
change Forge actually makes today.**

## The structural boundary holds

Three attacks were run against the verifier engine directly, bypassing the
pipeline so the contract layer had to defend itself.

| Attack | Expected | Result |
|---|---|---|
| Candidate supplies only its own controls, all marked consequential | Result must not pass | **Held.** Rejected at plan binding: `BindPlan` refuses any plan without at least one control that is `Consequential && !CandidateOwned` |
| Candidate invents a control and declares it Forge-owned | Verifier must refuse | **Held.** `ErrUnsupportedControl`. The verifier only implements controls compiled into its own binary |
| Candidate writes outside the declared paths | Path scope must fail | **Held.** Outcome `failed` |

Two independent enforcement points, not one:

- `BindPlan` refuses to construct a plan lacking a deciding independent check.
- `BindResult` refuses to mark a result `passed` unless some finding is
  consequential, not candidate-owned, and passed.

Neither can be satisfied by anything the candidate supplies, because the
verifier rejects control identifiers it does not implement. This is stronger
than any of the eight agent products surveyed, none of which enforce
generator/verifier separation at all.

## The deciding check is substantively empty

Forge implements exactly three Forge-owned controls:
`forge.candidate-integrity`, `forge.declared-path-scope`, `forge.go-format`.

Running them against a documentation-only candidate gives outcome `passed`, with
all three consequential controls passing. Examined individually, in that
scenario **none of them is capable of failing**:

| Control | Why it cannot fail on a documentation change |
|---|---|
| `forge.candidate-integrity` | `bundle.Validate()` is already enforced at verifier entry; a bundle that would fail here never reaches evaluation |
| `forge.declared-path-scope` | The lifecycle already requires `slices.Equal(assurancePlan.ChangedPaths, candidate.ChangedPaths)`, so the control compares the candidate against its own declaration |
| `forge.go-format` | The candidate contains no Go files, so the check is vacuous |

The third attack showed that `forge.declared-path-scope` *can* fail when the
plan and bundle disagree. That state is unreachable through the real pipeline,
because the lifecycle enforces equality before verification runs. Write-path
scope is genuinely enforced — but at execution, by the Execution Envelope, not
by the verifier.

So the verifier is independent, mandatory, and correct, and for a documentation
change it certifies nothing. Had Forge's pipeline processed this session's own
documentation commits, its deciding check would have passed vacuously on
every one of them.

This is the difference between satisfying AC.8.1 in form and in substance.

## Recommended remedy

Forge already owns controls that genuinely fail. `scripts/check_architecture.py`
runs 25 checks over document metadata, navigation coverage, decision-register
consistency, link resolution, schema sets, module boundaries and more. During
this session it rejected a real defect: a document committed without its
navigation entry.

The proposal is therefore not to invent a new check but to promote the ones
that exist:

1. Add `forge.architecture-fitness` as a Forge-owned consequential control that
   executes `scripts/check_architecture.py` inside the verifier.
2. Add `forge.go-test` and `forge.go-vet` as Forge-owned consequential controls
   for candidates touching Go.
3. Require the verification profile to select at least one control that is
   capable of failing for the affected paths, rather than merely present. A
   plan whose consequential controls are all vacuous for the changed file types
   SHOULD be refused.

Point 3 is the durable fix. Without it, adding controls only postpones the same
weakness to the next file type Forge learns to change.

## Remediation applied

All three recommendations were implemented, with one deliberate deviation.

**Deviation.** The recommendation was to execute `scripts/check_architecture.py`
inside the verifier. That was rejected on inspection: the verifier image is
`FROM scratch` with no shell, interpreter or toolchain, and adding them would
enlarge the attack surface of the one component whose entire purpose is being
trustworthy. The controls are instead compiled into the verifier binary, so the
image stays minimal and a candidate cannot alter a control's rule.

**Applicability, the durable fix.** `Applicable(controlID, changedPaths)` reports
whether a control is capable of failing for the paths under change. The two
integrity controls are permanently inapplicable, by design and with the reason
recorded in the code. `Verify` now refuses any plan whose deciding
independent checks are all inapplicable, returning `ErrNoApplicableCheck`.
Forge fails closed rather than certifying what it cannot check.

**Two controls capable of failing.** `forge.document-metadata` requires every
canonical document the candidate touches to declare a valid lifecycle status and
an owner, because authority in this repository is resolved by status.
`forge.markdown-links` requires every relative link in a changed Markdown file to
resolve, against either another candidate file or the independently
reconstructed checkout; external URLs and fragments are ignored.

**Boundary correction.** The first attempt placed the control identifiers in the
verifier package and referenced them from the control plane. The module boundary
check rejected it: `control-plane-cannot-import-verifier`. The identifiers now
live in the shared verification contract, which both sides may depend on. The
fitness checks caught this, unprompted, during implementation.

### Measured effect

| Scenario | Before | After |
|---|---|---|
| Documentation change, clean | certified, no control able to fail | refused under the old profile; passes on its merits under the new one |
| Documentation change, broken link | would have been certified | `failed` |
| Go change | `forge.go-format` applicable | unchanged |

The shipped M1-C profile carried only the two integrity controls, so before this
change the live journey certified strictly nothing. It now carries five.

Tests cover the refusal, both new controls in their passing and failing forms,
and link resolution against source, sibling candidate files, external URLs and
fragments.

## Provenance and the AC.8.1 invariant

Candidate provenance was recorded nowhere, so nothing could require that a
reviewer differ from the producer. That is now closed.

`AssurancePlan.ProducedBy` is mandatory and is set from the workload principal
Policy issued for the worker phase. The candidate cannot assert it.
`Result.VerifiedBy` is mandatory and is set by the verifier from its own
identity. `Result.Validate` refuses any result whose verifier is the producing
principal.

The rule lives in the contract rather than in a control, because a profile can
omit a control but cannot omit an invariant. This mirrors how the deciding-check
rule is already enforced.

**The invariant is dormant today, and that is stated rather than disguised.**
Forge's verifier is deterministic and therefore never the producing principal,
so the check cannot currently fire. It becomes load-bearing the moment
verification becomes model-assisted, which is precisely when it would be
expensive to retrofit. It is asserted by test now.

What has immediate value is the recording itself: every assurance plan now
names the principal that produced the candidate, which the
[delivery measures](ai-delivery-experiment.md) need and which OWASP AISVS
AC.10.1 requires.

### Contract mode does not verify

Investigating provenance surfaced a separate honesty problem. The default
contract-mode runner records a passing finding for every declared control
**without evaluating any of them**. It is a stub.

This is not an oversight but a consequence of the module boundary: the control
plane may not import the verifier deployment unit, so the real controls can only
execute in operational mode, where the verifier runs as a separate immutable
image. Contract mode exercises the trust path and the state machine and proves
nothing about the candidate.

Its results are now stamped `forge.verifier.contract-double.v1`, distinct from
the real verifier identity, so evidence can never conflate the two. Anyone
reading a result can tell whether anything was actually checked.

## Not established

- The verifier runs Forge-owned code, but nothing attests *which build* of the
  verifier ran. Binding the verifier image digest into `VerifiedBy` would close
  that and is the natural next step.
- Model identity for the ADR-0023 change-preparation path is not yet part of
  `ProducedBy`; the workload principal is recorded, not the model that proposed
  the change.
- These attacks target the contract and engine layers. The Docker-isolated
  verifier deployment path was not attacked.
- Contract mode remains unverified by construction, as described above.

## Reproduction

The attacks were first run through a throwaway harness that imported the
production engine and contracts directly. That harness is not retained, because
every scenario it exercised is now a permanent test:

| Scenario | Test |
|---|---|
| Plan whose independent controls cannot fail is refused | `TestVerifierRefusesPlanWithoutApplicableCheck` |
| Unknown control identifier is refused | covered by `ErrUnsupportedControl` in `evaluate` |
| Document without lifecycle status or owner fails | `TestDocumentMetadataControlFailsWithoutLifecycleStatus` |
| Unresolvable Markdown link fails | `TestMarkdownLinksControlFailsOnUnresolvableLink` |
| Links resolve against source and sibling candidates | `TestMarkdownLinksControlResolvesAgainstSourceAndCandidate` |
| Producer cannot verify its own candidate | `TestProducerCannotVerifyItsOwnCandidate` |
| Plan without recorded provenance is refused | `TestPlanRequiresRecordedProvenance` |

Run `go test ./internal/verifier/... ./internal/contracts/verification/...`.

A finding that only a discarded script can demonstrate is not evidence. Where an
attack is worth recording, it belongs in the suite that runs on every change.
