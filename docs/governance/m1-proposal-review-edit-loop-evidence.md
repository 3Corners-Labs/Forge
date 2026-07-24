# M1 Requirements Proposal review and edit evidence

Status: Proposed
Owner: Forge Product, Experience, Intent, Quality, and Engineering
Verified: 2026-07-24

## Iteration outcome

This iteration adds a bounded conversational review loop between Requirements
Proposal creation and promotion. A human can:

- promote the exact current proposal;
- replace the outcome, problem, or one existing item in an allowlisted
  structured field; or
- reject the proposal without creating canonical Intent or Delivery authority.

Each accepted edit creates a new immutable proposal version with a new digest,
the previous digest, a required human change explanation, and explicit
human-correction provenance. The prior version remains available as superseded
history. It is never overwritten or promoted as current.

## Why

The synthetic live discovery trial showed that a model can produce an otherwise
useful proposal with one or two semantically wrong requirements. Regenerating
the complete proposal spends time and provider budget, can lose good content,
and can introduce new errors. A precise human correction is the smaller and
more trustworthy operation.

The loop also makes the authority transition visible: editing remains proposal
refinement, rejection creates no authority, and promotion alone creates
canonical Intent.

## Applicable architecture

- The Constitution requires explicit provenance, human agency, untrusted model
  output, minimum necessary authority, auditable transitions, and responsible
  resource use.
- ADR-0020 defines a versioned, digest-bound Requirements Proposal, permits
  ordinary draft refinement without micro-approval, and requires promotion of
  only the exact current digest.
- ADR-0021 permits the CLI to coordinate the resumable guided journey through
  existing capability contracts. It does not allow the CLI to own canonical
  records or authority.
- ADR-0022 keeps model assistance non-authoritative and requires deterministic
  capabilities to remain usable without a model.
- The Experience architecture owns proposal lifecycle and role-appropriate
  interaction. The Intent architecture owns the explicit transition to
  canonical product intent.

## Capability effects

- `forge.capability.experience` validates the closed edit request, creates the
  new immutable version, preserves superseded history, records the human
  correction, and resolves only the current version index for promotion.
- `forge.capability.intent` is unchanged. Its existing exact-current checks now
  reject a superseded digest and accept only the revised digest.
- The CLI remains an Experience adapter and coordinator. It presents
  `promote`, `edit`, and `reject`, but does not mutate storage or grant
  authority itself.

Delivery, Policy, Approval, Intelligence, Execution, Verification, Evidence,
and Source Control are unchanged. No module dependency, deployment unit,
public schema, accepted journey oracle, or authority boundary changes.

## Automated verification

Tests prove:

- a list-item correction retains the proposal ID and creates version 2 with a
  new digest and the exact previous digest;
- the previous digest-bound version is retained unchanged while current
  resolution moves atomically to the revision;
- the correction and its provenance are included in the promotion subject;
- stale digests, unknown fields, and invalid list indexes fail closed;
- Intent rejects the superseded digest and accepts the exact revised digest;
- the conversational CLI can edit and then promote the returned revision;
- rejection stops without calling promotion; and
- the existing deterministic and model-assisted guided journeys still reach
  their previous bounded outcomes.

The final candidate regenerates repository tests, race tests, architecture
fitness, vulnerability scanning, and CI evidence.

## Operational behaviour and rollback

The edit operation is local and deterministic. It does not invoke a model,
tool, repository provider, worker, verifier, or publication adapter. The first
slice supports scalar replacement for outcome and problem, plus replacement of
one existing item in a closed set of structured text fields. It intentionally
does not add, remove, reorder, or edit structured questions.

Rollback removes the review/edit coordinator and revision operation. Existing
version-1 proposals and exact promotion remain compatible. No credential,
external data, release, deployment, or production recovery is required.

## Limitations

- Revision history uses the in-memory Developer Trial store and does not
  survive process termination.
- The CLI currently renders JSON events plus accessible line prompts; a richer
  terminal interface remains a later Experience iteration.
- Adding or removing list items and resolving structured questions require
  separate typed operations rather than overloading replacement.
- Rejection is an interaction outcome, not yet a durable Experience record.

## ADR assessment

No ADR is proposed. ADR-0020 already requires versioned proposal refinement and
exact-current promotion; ADR-0021 already assigns coordination to the CLI; and
ADR-0022 already keeps model output non-authoritative. This iteration
implements those accepted decisions without changing architecture.
