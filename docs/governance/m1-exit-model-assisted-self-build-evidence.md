# M1 exit model-assisted self-build evidence

Status: Accepted
Owner: Forge Product, Intelligence, Identity, Policy, Execution, Verification, Source Control, Experience, Security, Privacy, Quality, and Engineering
Verified: 2026-07-25

## Iteration outcome

This iteration replaces the fixed M1 candidate template boundary with the
ADR-0023 change-preparation contract and adds the ADR-0019 GitHub App provider
implementation. It does not activate an external model, GitHub credential, or
live repository mutation.

The public Developer Trial journey now constructs one digest-bound
`ChangePreparationRequest` after exact plan approval. With no live activation,
the contract double derives a reviewable candidate from the promoted
Requirements Proposal instead of writing fixed content. With a separately
approved `forge.m1.prepare-change` activation, the same journey can obtain a
strict structured `ChangeProposal` through Intelligence.

Model output remains information. Identity mints a phase-specific task-agent
delegation; Policy issues exact model-proposal authority; the worker delegation
and decision bind the admitted proposal digest; Execution converts only its
declared UTF-8 replacements into the closed worker contract; and Verification
derives and evaluates its own Assurance Plan. Publication still requires a
later human approval for the exact verified candidate.

## Why

The prior journey proved the trust chain but could not test Forge's central
product hypothesis because the worker always wrote the same Markdown. This
slice introduces meaningful requirements-to-change behaviour without granting
repository, tool, shell, verification, publication, or merge authority to a
model.

The GitHub App adapter completes the bounded provider seam needed for a later
real pull-request trial. Keeping live credentials and effects behind separate
activations lets the implementation and negative contracts mature before
external cost or repository mutation occurs.

## Applicable architecture

- The Engineering Constitution requires human authority, untrusted
  intelligence, explicit capability contracts, bounded execution,
  independent verification, responsible resource use, and evidence.
- ADR-0010 bounds M1 external providers and prohibits production effects.
- ADR-0013 keeps assurance proportional while consequential guardrails remain
  publication blockers.
- ADR-0017 requires deterministic isolated materialization and independent
  verification.
- ADR-0019 defines the dedicated GitHub App permissions, endpoint allow-list,
  reconciliation, protected-branch denial, and no-merge boundary.
- ADR-0020 and ADR-0021 preserve proposal promotion, plan approval, and
  publication approval as separate human actions.
- ADR-0022 keeps model output non-authoritative.
- ADR-0023 explicitly permits the bounded model-assisted change-preparation
  path implemented here.

No architecture is changed by this implementation, so no new ADR is required.

## Capability effects

- Context supplies exact context and source bindings without provider
  authority.
- Delivery owns the plan, declared paths, lifecycle, candidate, and recovery
  state.
- Identity owns separate task-agent and worker delegations.
- Policy owns model-proposal and materialization decisions and envelopes.
- Intelligence owns activation eligibility, provider invocation, structured
  output admission, measured usage, and receipts.
- Execution owns deterministic conversion and the isolated worker port.
- Verification remains independent from model and candidate-controlled tests.
- Evidence retains digest references without unrestricted prompts or source.
- Approval retains separate human publication authority.
- Source Control owns idempotent provider publication and reconciliation.
- Experience reports proposal digest, provider, model, usage, uncertainty,
  omissions, candidate digest, verification digest, and publication state.

The shared `change-contracts` module has no runtime dependencies. Intelligence
does not import Worker, Execution does not import the worker implementation,
and Source Control does not import Approval or Delivery persistence.

## Implemented controls

The change request and proposal contracts:

- bind tenant, workspace, project, repository, run, task agent, plan, context,
  requirements, source revision, purpose, classification, expiry, budgets,
  allowed paths, acceptance criteria, test intent, controls, and source
  digests;
- reject unknown fields at the OpenAI structured-output edge and reject
  undeclared paths, stale source digests, invalid UTF-8, NUL bytes, duplicates,
  oversized input/output, budget overruns, invented references, refusal, and
  incomplete output before materialization;
- carry no credential value, Approval Record, Policy Decision, Execution
  Envelope, shell, tool, branch, merge, release, deployment, or production
  authority; and
- record activation, model, adapter, role, prompt, schema, provider request,
  tokens, latency, and estimated cost.

The OpenAI adapter disables provider storage, requests one strict JSON schema,
forbids redirects and environment-proxy routing, resolves the credential only
at the eligible invocation, exposes no tool field, and discards provider error
bodies.

The GitHub App provider:

- accepts one configured repository and credential reference;
- uses only Git references, commits, blobs, trees, and pull-request endpoints;
- creates a deterministic branch and draft pull request;
- binds commit and pull-request state to the operation marker, base revision,
  candidate digest, and tree digest;
- treats mutation ambiguity as a reconciliation requirement; and
- has no merge, release, deployment, Actions, secrets, hooks, collaborator, or
  administration method or allowed endpoint.

## Verification

Automated tests cover:

- canonical request and proposal digests and defensive cloning;
- exact source, plan, context, activation, receipt, path, budget, expiry, and
  classification binding;
- phase-specific Identity delegation and Policy obligation enforcement;
- cross-scope, stale-source, refusal, incomplete, malformed, oversized,
  mismatched-authority, and provider-failure denial;
- absence of model tool authority and credential leakage from request bodies;
- proposal-digest-bound worker authority and dependency evidence;
- worker path, symlink, source-digest, output, cleanup, candidate, and tree
  controls;
- independent Assurance Plan and verification-result validation;
- exact GitHub App endpoint allow-list, short-lived token reference,
  deterministic commit marker, draft PR reconciliation, foreign repository
  denial, and explicit merge/release/deployment denial;
- the complete existing allowed and denied public M1 journeys; and
- module-boundary fixtures proving that Intelligence cannot import Worker and
  shared Change contracts cannot import runtime capabilities.

`assurance/evaluations/change-preparation/v1/plan.json` records the synthetic
quality and safety cases. The accepted M1 journey oracle is not changed.

The implementation was verified with:

- `go test ./...`;
- `go test -race ./...`;
- `go vet ./...`;
- `go build ./...`;
- `python3 scripts/check_architecture.py`; and
- strict JSON parsing of both new assurance records.

The local environment did not contain MkDocs. Canonical navigation, metadata,
and link checks passed through `scripts/check_architecture.py`; the repository
CI remains responsible for its pinned `mkdocs-material==9.7.7` strict build.

## Current limitation and next gate

This candidate is merge-ready implementation evidence, not final live M1 exit
evidence. The following remain separate human gates:

1. approve one exact `forge.m1.prepare-change` provider activation and run the
   synthetic quality trial;
2. review the ADR-0019 GitHub App repository permissions and protected-branch
   denial evidence;
3. approve one exact verified candidate for a real draft pull request; and
4. review and merge that pull request through the normal protected branch.

Live merge, release, deployment, production, general agent routing, tools,
multi-agent coordination, and unrestricted repository access remain forbidden.
