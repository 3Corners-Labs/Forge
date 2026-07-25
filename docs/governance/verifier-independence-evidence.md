# Verifier independence: attack results

Status: Draft
Owner: Forge Engineering and Architecture
Verified: 2026-07-25

## Question

OWASP AISVS 1.0 Appendix C, requirement AC.8.1, states that agents cannot
approve, merge, sign, or deploy artifacts that they themselves generated. Market
research found no shipped agent product enforcing this structurally.

Forge claims a mandatory independent consequential oracle. This exercise
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

- `BindPlan` refuses to construct a plan lacking an independent consequential
  control.
- `BindResult` refuses to mark a result `passed` unless some finding is
  consequential, not candidate-owned, and passed.

Neither can be satisfied by anything the candidate supplies, because the
verifier rejects control identifiers it does not implement. This is stronger
than any of the eight agent products surveyed, none of which enforce
generator/verifier separation at all.

## The independent oracle is substantively empty

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
documentation commits, its independent oracle would have passed vacuously on
every one of them.

This is the difference between satisfying AC.8.1 in form and in substance.

## Recommended remedy

Forge already owns controls that genuinely fail. `scripts/check_architecture.py`
runs 25 checks over document metadata, navigation coverage, decision-register
consistency, link resolution, schema sets, module boundaries and more. During
this session it rejected a real defect: a document committed without its
navigation entry.

The proposal is therefore not to invent an oracle but to promote the ones that
exist:

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
recorded in the code. `Verify` now refuses any plan whose independent
consequential controls are all inapplicable, returning `ErrNoApplicableOracle`.
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

## Not established

- **Agent identity is not bound anywhere.** Nothing records which model or
  harness produced a candidate, so nothing can require that a future
  model-assisted reviewer differs from the generator. The cross-harness rule
  ("one model may not review its own work") is absent, and it will be needed
  before any model-assisted verification is added.
- The verifier runs Forge-owned code, but nothing attests *which build* of the
  verifier ran. Binding a verifier image digest into the result would close
  that.
- These attacks target the contract and engine layers. The Docker-isolated
  verifier deployment path was not attacked.

## Reproduction

The attack harness is `dist/verifierproof/main.go`, in the gitignored `dist/`
directory. It imports the production `internal/verifier` engine and the
production contracts, so the behaviour under test is the shipped behaviour.
