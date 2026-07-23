# M1-C guided delivery to reviewable pull-request evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-23

## Bounded outcome

This combined M1-C tranche makes the first complete Forge Developer Trial
journey executable through the public CLI:

~~~text
untrusted requirements input
  -> Requirements Proposal
  -> explicit proposal promotion
  -> Delivery Plan
  -> separate plan approval
  -> bounded candidate preparation
  -> independent verification
  -> separate publication approval
  -> one idempotent contract-double pull request
~~~

The result is reviewable learning evidence. It is not a real GitHub pull
request and grants no merge, release, deployment, production, or live
credential authority.

## Why this combined tranche

A partial internal kernel could validate contracts but could not answer the
product question of how a person experiences Forge. ADR-0020 requires free text
to become a proposal rather than authority. ADR-0021 permits the public CLI to
coordinate the existing capability contracts while retaining three distinct
human decisions. Combining those boundaries with M1-B execution,
verification, evidence, and M1-C contract-double publication produces one
feedback loop without claiming live-provider readiness.

## Governing architecture

The Accepted Constitution, trustworthy vertical slice, first-slice capability
contracts, first-slice data flow, experience, policy, evidence, security, user
journey assurance, and quality architecture govern this tranche. ADR-0009,
ADR-0010, ADR-0013, ADR-0016, ADR-0017, ADR-0018, ADR-0019, ADR-0020, and
ADR-0021 apply.

ADR-0021 is the only architectural change in this tranche. It supersedes only
ADR-0009's former stopping restriction and leaves all capability ownership,
authority, provider, no-merge, no-release, no-deployment, and no-production
constraints intact.

## Implemented capability scope

### Experience and Intent

Experience owns an immutable Requirements Proposal that binds tenant,
principal, session, direct statements, assumptions, hypotheses, open questions,
experiments, and an exact canonical digest. Intent permits promotion only by
the current authenticated human against the exact current proposal. Material
open questions block promotion. Promotion then invokes Delivery's existing
intent-capture contract; neither Experience nor Intent manufactures scope or
execution authority.

### Approval, Policy, and Delivery

Approval adds a separate publication-review action for a
`verified-candidate`. It resolves the current Delivery-owned candidate subject,
including verification and evidence references, both when review opens and
when the human decision is recorded.

Policy evaluates publication eligibility and alone issues the bounded
publication Execution Envelope. Delivery exposes a publication subject only
from `AwaitingPublicationApproval`; it does not approve or publish.

Plan approval and publication approval have distinct actions, subject types,
digests, request identifiers, Policy decisions, records, and idempotency keys.
One cannot be replayed as the other.

### Source Control

Source Control owns provider-neutral branch and pull-request publication. It
re-resolves the current Delivery subject, Approval record, Policy decision,
Execution Envelope, Evidence stream, and Tenant-derived repository binding
before any provider mutation. The exact candidate, resulting tree, changed-path
manifest, base branch and revision, provider, minimized title and body, and
required candidate evidence must match the trusted publication profile.

The M1 provider is a credential-free deterministic contract double. Operation
and branch identities derive from the exact bounded subject. Reconciliation
handles transport ambiguity before retry and rejects conflicting remote state.
Repeated exact publication returns the same operation and pull request rather
than creating duplicates. The provider exposes no merge operation.

### Public CLI experience

`forge delivery start` reads a closed bounded JSON requirements input, emits
immediate newline-delimited JSON progress, and requires exactly:

1. `promote <proposal-digest>`;
2. `approve-plan <plan-digest>`; and
3. `approve-publication <candidate-digest>`.

The default Developer Trial journey composes synthetic identity, Tenant-derived
scope, deterministic Context and Delivery planning, Policy and Approval,
bounded contract execution, independent contract verification, Evidence, and
Source Control. The final projection states that human review and merge remain
required.

The synthetic denied observer is rejected before repository-binding discovery.
Caller-provided workspace, project, repository, digests, and terminal input
remain untrusted selectors or decisions rather than canonical authority.

## Journey oracle

The allowed M1 oracle is superseded to version 2.0.0. Success now requires the
complete guided flow through one reviewable contract-double pull request.
Forbidden outcomes include raw text becoming authority, execution before plan
approval, plan-approval replay, publication before verification or evidence,
live GitHub credentials, duplicate publication, merge, release, deployment,
and production effects.

The denied oracle remains fail-closed before repository discovery. The
architecture checker fixes both expectations and the executable candidate
state so a later candidate cannot silently weaken or shorten the journey.

## Tests and verification

Deterministic tests cover:

- proposal canonicalization, digest mutation, direct-statement separation,
  tenant isolation, current-proposal resolution, material-question blocking,
  human-only promotion, and exact promotion replay;
- separate publication request and Approval Record v2 semantics, current
  candidate resolution, evidence binding, expiry, actor eligibility, and
  rejection of plan-approval replay;
- exact publication Policy evaluation and envelope issuance;
- Source Control denial for wrong approval, missing evidence, changed
  candidate, tree or paths, cross-repository authority, unreviewed metadata or
  base revision, conflicting provider state, and invalid metadata;
- idempotent replay and ambiguity reconciliation before and after branch or
  pull-request creation;
- all three exact public human actions, immediate stage progress, refusal of a
  wrong proposal digest, and `mergePermitted: false`;
- the complete default Developer Trial journey to one reviewable
  contract-double pull request; and
- the denied cross-scope journey stopping before repository discovery.

Candidate verification SHALL include the architecture checker, full Go test
suite, race detector, Go vet, build, and vulnerability scan. Exact observations
are revision-specific and SHALL be attached to the pull request checks.

## Performance and scalability posture

The public command writes progress before each potentially longer stage. Pure
digest, policy, approval, and idempotency checks remain deterministic and
in-process for this trial. Expensive preparation and verification occur only
after explicit plan approval; publication validation occurs only after
independent verification and separate approval.

The contract-double journey is not capacity evidence for PostgreSQL, MinIO,
GitHub, Webex, model providers, concurrent projects, large repositories, or
production workloads. ADR-0013 remains the gate for tool selection and measured
p50, p95, p99, throughput, resource, cost, and failure envelopes. Required
security, privacy, authority, integrity, and evidence checks are never sampled
away for performance.

## Explicit exclusions

This tranche does not implement free-form conversation, live mockups, Webex
Messaging, model or agent orchestration, live repository reads, live GitHub
credentials, durable cross-process session resume, production isolation,
release, merge, deployment, production data, or supported operation.

The in-process contract worker and verifier prove orchestration and trust
ordering, not the process or container isolation already demonstrated
separately by M1-B.

## Authorized delivery evidence

Delivery used grant `forge.autonomy-grant.m1-bounded-delivery-v1`, protected
revision `d96e9fb5596e30eb47cdff5b8659c3e421d7b865`, task
`m1-c-guided-delivery-to-reviewable-pr`, and branch
`codex/m1-c-guided-delivery-to-reviewable-pr`.

The activated scope includes ADR-0021, the M1-C capability code, public CLI,
accepted journey oracle, architecture fitness, implementation baseline, this
evidence document, and canonical navigation. Protected-default-branch landing
remains a mandatory exact-PR and exact-head human gate.

## Architecture decision assessment

No additional ADR is proposed. The implementation realizes Accepted ADR-0021.

A new ADR SHALL precede live provider credentials, another public channel such
as Webex Messaging, durable coordinator ownership, a different approval model,
new capability ownership, service extraction, merge or release authority,
production deployment, or weakening an accepted journey oracle.
