# M1 live model discovery trial readiness evidence

Status: Accepted
Owner: Forge Product, Experience, Intelligence, Security, Privacy, Quality, and Engineering
Verified: 2026-07-24

## Iteration outcome

This iteration makes the ADR-0022 conversational path runnable with one exact,
separately approved OpenAI Developer Trial activation. It does not activate a
provider or make a live request.

`forge delivery start --message` can now load a closed activation JSON file,
bind it to a separately supplied SHA-256 digest, configure the bounded OpenAI
Responses adapter, acknowledge the request immediately, and project each
non-authoritative discovery turn with its exact activation, model, prompt,
schema, token, latency, and calculated cost receipt.

The live path remains fail-closed when the activation is absent, malformed,
expired, replaced, symlinked, outside its endpoint or classification scope,
over budget, or different from the separately approved digest. Credentials
remain environment references and are resolved only immediately before an
eligible provider request.

## Why

The prior M1 implementation proved the authority chain with contract doubles
but left the real conversational virtual Product Manager unreachable from the
public CLI. Forge needs end-user discovery evidence before hardening a product
experience that may be wrong. The new runtime seam allows that learning without
making a model authoritative or pulling general M3 routing and M5 multi-agent
coordination into M1.

The activation also records inference effort and token pricing. Forge can
therefore compare quality, privacy, latency, and estimated cost rather than
selecting a tool or model on reputation alone.

## Applicable architecture

- The engineering Constitution requires human authority, untrusted
  intelligence, explicit capability contracts, bounded execution, responsible
  resource use, independent trust, and governed evolution.
- The Intelligence architecture requires policy and data eligibility before
  quality, locality, latency, availability, and resource selection, then the
  smallest sufficient model.
- The Security architecture requires authorization before secret access and
  provider disclosure, late-bound secrets, schema validation, allow-listed
  egress, and digest-bound approvals.
- The Deployment architecture limits this work to the connected, local
  Developer Trial profile. It establishes no Supported or production profile.
- ADR-0020 keeps interaction input non-authoritative until exact Requirements
  Proposal promotion.
- ADR-0021 preserves separate proposal, plan, verification, publication, and
  merge boundaries.
- ADR-0022 permits exactly this narrow provider-replaceable,
  non-authoritative discovery path and requires a separate provider activation
  before the first live request.

## Capability effects

- `forge.capability.experience` owns the CLI conversation, immediate progress,
  questions, draft projection, and presentation of the Requirements Proposal.
- `forge.capability.intelligence` owns activation loading, eligibility,
  OpenAI invocation, prompt and model binding, limits, normalized failure, and
  usage receipts.
- `forge.capability.intent` remains the only owner of promoted canonical
  intent.
- Identity and Policy-derived scope remains ahead of repository discovery,
  credential resolution, and provider invocation.
- Delivery and Approval retain the existing exact proposal, plan, and
  publication gates.

No capability imports another capability's authority, no module dependency
changes, and no new deployment unit is introduced.

## Activation and invocation

The activation file conforms to
`schemas/forge/intelligence/v1/provider-activation.schema.json`. It contains no
credential value. Its `credentialRef` names an environment variable such as
`env:FORGE_OPENAI_API_KEY`.

An accountable human reviews the exact file, its data posture, endpoint,
model, inference effort, prompt, retention, price assumptions, limits, expiry,
credential scope, evaluation evidence, and failure behaviour. The human then
approves the exact SHA-256 digest through the authenticated collaboration
boundary. The file cannot approve itself: Forge requires that separately
approved digest at runtime and compares it before constructing an eligible
advisor.

The Developer Trial invocation is:

```text
forge delivery start \
  --workspace fixture-workspace \
  --project forge \
  --repository forge \
  --message "<synthetic product request>" \
  --provider-activation "<approved activation.json>" \
  --provider-activation-digest "sha256:<approved exact digest>" \
  --output json
```

The credential is injected into the named environment only for the local
process. It SHALL NOT be written into the activation, command line, repository,
prompt, event stream, error, or evidence.

The activation is loaded from one regular file with a 64 KiB limit, strict
unknown-field rejection, one-object decoding, exact byte digest, open-time file
identity check, closed model and classification allow-lists, exact adapter and
prompt versions, endpoint-to-region binding, expiry, turn, input, output,
token, latency, and cost limits.

The OpenAI adapter disables provider application-state storage, uses a strict
closed response schema, sets the approved reasoning effort, applies an output
token cap, rejects redirects, disables environment-proxy routing, and discards
provider error bodies. The cost receipt is calculated from actual response
token counts and the approved activation rates. It is an estimate for
selection and budget enforcement, not a substitute for provider billing
reconciliation.

## Model evaluation

`assurance/evaluations/discovery/v1/plan.json` defines five synthetic cases for:

- ambiguity and professional requirement discovery;
- prompt injection and authority claims;
- privacy and data minimization;
- UX responsiveness, performance, cost, and operations; and
- reversible mockup learning.

The controlled candidates are `gpt-5.6-luna`, `gpt-5.6-terra`, and
`gpt-5.6-sol` under the same adapter, prompt, schema, medium reasoning effort,
and budgets. Current OpenAI guidance describes Luna as cost-sensitive, Terra as
the intelligence-and-cost balance, and Sol as the flagship for complex work.
Terra is the provisional reference candidate; the accepted selection rule
still chooses the lowest-cost model that passes every deterministic safety,
provenance, privacy, performance, cost, and human-quality threshold.

Each model requires a separate exact activation so results cannot be
misattributed. A live evaluation uses only synthetic data and cannot proceed
from this plan alone.

Official provider references:

- [OpenAI model catalogue](https://developers.openai.com/api/docs/models)
- [OpenAI GPT-5.6 guidance](https://developers.openai.com/api/docs/guides/latest-model)
- [OpenAI structured outputs](https://developers.openai.com/api/docs/guides/structured-outputs)
- [OpenAI API data controls](https://developers.openai.com/api/docs/guides/your-data)

## Automated verification

Tests cover:

- exact activation parsing, byte-level digest changes, expiry, closed fields,
  non-regular files, symlink rejection, and open-time identity;
- separate runtime digest comparison before an advisor becomes available;
- exact provider, endpoint, region, adapter, prompt, model, reasoning,
classification, credential reference, pricing, and budget validation;
- atomic per-session and per-project/day token and financial reservation,
  measured-use reconciliation, and concurrency enforcement before invocation;
- credential resolution only at invocation;
- storage-disabled, structured, output-token-bounded OpenAI requests;
- redirect and unapproved endpoint rejection before unsafe disclosure;
- response schema, provenance, receipt, activation, token, latency, and cost
  validation;
- measured cost calculation from provider token counts;
- synthetic evaluation-plan safety thresholds;
- activation-option propagation through the public CLI;
- preserved conversational proposal, plan, verification, publication, and
  denied-before-discovery journeys; and
- the existing deterministic no-model path.

The full repository verification result is regenerated for the final candidate
and recorded on its pull request.

Three local Apple M1 Pro contract-double benchmark runs, including activation,
session, project/day, concurrency, receipt, and proposal validation, measured
approximately 4.9 to 5.1 microseconds per turn, 4,289 bytes, and 14 allocations.
This is a bounded Forge-overhead observation, not a provider-latency or
production-capacity claim.

## Explicit exclusions and residual risk

No provider activation, API credential, API call, customer data, production
data, repository content, live GitHub publication, merge, release, deployment,
or production effect is included.

The current Developer Trial uses a synthetic local identity and an
authenticated human approval outside the process. It does not yet provide a
durable signed provider-activation store, billing reconciliation, general
model registry, dynamic routing, fallback provider, durable conversation,
Workbench, Webex channel, multiple virtual Product Managers, Agent Broker,
tools, or production secret broker.

The CLI is still a linear JSON event stream with readable prompts rather than
the planned rich conversational terminal. The live trial is intended to
produce the feedback needed for that UX round.

## ADR assessment

No new ADR is proposed. ADR-0022 already accepted the external provider seam,
activation gate, non-authoritative model role, evaluation obligation, and
public conversational boundary. This iteration implements that decision
without changing capability ownership, authority, milestone scope, deployment
topology, or accepted journey outcomes.
