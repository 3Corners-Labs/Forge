# M1 discovery resilience evidence

Status: Proposed
Owner: Forge Product, Experience, Intelligence, Quality, and Engineering
Verified: 2026-07-24

## Iteration outcome

This iteration preserves an active conversational discovery checkpoint when a
provider returns a normalized bounded-retry failure. The public CLI emits an
`awaiting-retry` event, explains that current answers remain available, and
requires one meaningful retry decision. A retry reuses the same session ID,
direct statements, draft, prior questions, and provenance references. It does
not repeat questions that the human already answered.

For a separately activated live provider, Experience now uses the exact
validated per-request input, output, token, latency, and cost limits from the
provider activation. The former fixed twenty-second Developer Trial latency is
retained only as the contract-double default.

The checkpoint is an in-process Experience projection, not canonical Intent or
authority. Durable multi-process resume still requires a capability-owned
persistence adapter and opaque server-side reference.

## Why

The first live multi-turn Terra trial produced one valid discovery turn and
then timed out while refining the draft. The CLI returned the normalized
failure but terminated, losing the running coordinator state and forcing the
human to repeat answers in a new process. A second attempt with a consolidated
request hit the same fixed limit.

This behaviour violated the intended degradation outcome: provider failure
must pause discovery without corrupting or discarding the deterministic
session. Raising only the activation budget could not help because Experience
replaced the approved latency with its own fixed value.

## Applicable architecture

- Constitution sections 5, 11, 16, 18, and 19 govern untrusted intelligence,
  responsible resource use, bounded execution, recoverability, and governed
  evolution.
- ADR-0021 permits the public CLI to coordinate a resumable session while
  capability-owned records remain canonical and explicitly excludes
  client-checkpoint authority.
- ADR-0022 requires provider timeout and unavailability to preserve resumable
  deterministic session state. Experience owns interaction state;
  Intelligence owns eligibility, provider invocation, budgets, receipts, and
  normalized failure.
- The first-slice capability contracts make `RetryBounded` the only retry
  posture used here. Budget exhaustion and new-authority failures remain
  fail-closed.
- The Intelligence architecture requires approved latency and resource budgets
  to participate in selection and routing rather than being silently replaced.

## Capability effects

- `forge.capability.experience` preserves and presents the non-authoritative
  interaction checkpoint and asks the human whether to spend another approved
  provider turn.
- `forge.capability.intelligence` continues to enforce the exact activation,
  per-turn and aggregate budgets, concurrency, expiry, and normalized retry
  classification.
- Intent, Delivery, Policy, Approval, Verification, Evidence, and
  Source Control are unchanged and receive no model or checkpoint authority.

No module dependency, deployment unit, public schema, accepted journey oracle,
or authority boundary changes.

## Live trial observation

The approved synthetic multi-turn trial used:

- activation digest
  `sha256:7a01a9f07021cd10adc054fabd1d176c6ad7a2a6c7a43611b4d98fdd2d460e37`;
- model `gpt-5.6-terra`;
- medium reasoning effort;
- a twenty-second per-turn latency budget; and
- synthetic Internal content with provider state storage and training disabled.

The first turn completed in 11,456 milliseconds using 805 input tokens and 977
output tokens, with an estimated cost of 16,668 microdollars. It asked three
material questions about acceptance, access, and mockup scope. The human's
existing Forge requirements supplied answers covering separate acceptance
criteria, independent validation and approval, SSO, RBAC with future
attribute-level conditions, and lightweight in-browser plus linked mockups.

The next refinement and one consolidated restart both reached the exact
twenty-second deadline. No valid usage receipt was available for either timed
out request, so their final provider token use and cost remain unknown and are
not represented as zero. The fourth approved turn was intentionally not spent.

This is evidence that the tested Terra/medium combination does not yet satisfy
the accepted evaluation plan's twenty-second p95 threshold for this case. The
threshold and model-selection oracle are unchanged by this iteration.

## Automated verification

Tests prove:

- a live activation's exact limits configure the Experience request;
- the no-activation contract-double path retains bounded defaults;
- a `RetryBounded` normalized failure emits an accessible pause event;
- retry preserves session, statements, draft, prior questions, and provenance;
- an answered question is not presented again;
- a failure requiring new authority is neither retried nor reclassified;
- proposal promotion, plan approval, verification, publication approval, and
  the denied-before-discovery journey remain unchanged; and
- provider eligibility, budgets, expiry, concurrency, receipts, privacy, and
  schema validation remain owned by Intelligence.

The final candidate regenerates repository tests, race tests, architecture
fitness, vulnerability scanning, and CI evidence.

## Operational behaviour and rollback

Retry is never automatic after a provider timeout because the provider may have
processed chargeable tokens without returning a receipt. The user sees the
safe normalized failure and chooses whether to consume another already
authorized turn. Intelligence still blocks the retry when turn, session,
project, cost, concurrency, or activation limits are exhausted.

Rollback reverts the CLI retry projection and restores the fixed Developer
Trial defaults. It does not require provider, credential, data, release,
deployment, or production recovery.

## Limitations

- The checkpoint survives retry only while the CLI process remains active.
- Process termination, host loss, or restart still requires a future
  capability-owned durable Experience adapter.
- No automatic provider fallback, dynamic model routing, billing
  reconciliation, Workbench, Webex adapter, or Agent Broker is introduced.
- A longer latency activation is not evidence that a model passes the accepted
  model-selection performance threshold.

## ADR assessment

No new ADR is proposed. ADR-0021 already defines resumable coordination and
prohibits client-held authority. ADR-0022 already assigns interaction state to
Experience, provider budgets and normalized failure to Intelligence, and
requires timeout recovery. This iteration implements those accepted decisions
without changing architecture.
