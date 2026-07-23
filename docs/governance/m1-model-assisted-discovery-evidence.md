# M1 model-assisted product discovery evidence

Status: Accepted
Owner: Forge Product, Experience, Intelligence, Security, Privacy, and Engineering
Verified: 2026-07-23

## Outcome

This iteration implements the non-live portion of ADR-0022. Forge can accept a
natural-language product request through `forge delivery start`, obtain bounded
adaptive discovery turns through a provider-neutral contract, ask material
questions, convert a complete discovery draft into a Requirements Proposal, and
continue through the existing separate proposal, plan, and publication gates.

The default CLI has no eligible model provider. A live call remains impossible
until an accountable human activates one exact provider record. Tests use a
credential-free contract double; the OpenAI adapter is tested only through an
in-process HTTP transport double.

## Architecture and capability boundaries

- `forge.capability.experience` owns the conversational turn and public
  projection. It does not select a provider or own canonical Intent.
- `forge.capability.intelligence` owns provider activation, adapter invocation,
  role/prompt/model/schema references, resource usage, and normalized provider
  failure.
- `forge.capability.intent` remains the only owner of promoted canonical
  intent.
- Identity and Tenant scope are resolved before repository discovery and model
  disclosure. The denied cross-scope journey reaches neither.
- The virtual Product Manager is a versioned behavioural role. It receives no
  tool, repository, approval, execution, publication, release, deployment, or
  production authority.
- The new shared discovery contract imports no runtime capability. Experience
  and Intelligence depend on that contract rather than each other's runtime.

## Implemented contracts

`DiscoveryTurnRequest` binds the session, correlation, server-resolved scope,
human and task-agent references, role, purpose, classification, direct
statements, current draft, prior questions, provenance, limits, and expiry.

The request and `DiscoveryTurnProposal` pair separates direct statements from
the response, complete or evolving requirement draft, inferences, assumptions,
hypotheses, material questions, uncertainty, limitations, refusal,
incompleteness, provenance, resource use, and provider receipt. Validation
rejects unknown
versions, scope mismatch, invented provenance, invalid questions, incomplete
final drafts, missing receipts, and token, latency, cost, input, or output
budget overruns.

The OpenAI adapter:

- permits only the declared Responses API endpoint hosts;
- resolves a credential reference at invocation time and excludes its value
  from request bodies, errors, and evidence;
- sends `store: false`;
- requests a strict closed JSON Schema;
- adds provider identity and usage from the provider response rather than
  trusting model-generated receipt fields; and
- rejects redirects so an authorization header cannot cross the approved
  endpoint boundary; and
- discards provider response bodies on HTTP failure.

Intelligence rejects request limits above the provider activation budgets and
rejects any response receipt whose provider, model, adapter, or prompt version
does not match that activation. It also enforces the activation's per-session
turn limit before adapter invocation and bounds the complete provider wire
request and response, not only the embedded discovery object.

## Experience and human attention

The deterministic requirements-file journey preserves exact digest-entry
regression coverage. The conversational path keeps the same exact bound
subjects but presents meaningful `yes` or `no` decisions while Forge retains
the digest internally.

The review projection preserves outcome, users, problem, functional and UX
behaviour, accessibility, privacy, security, operations, constraints,
exclusions, acceptance criteria, inferences, assumptions, hypotheses,
non-material questions, experiments, and product-area references. These fields
remain bound into the exact Requirements Proposal digest rather than being
reduced to a model-authored summary.

Clarification, assumptions, revisions, mockup proposals, progress, tests,
bounded repair, retry, and verification do not create micro-approval gates.
Human attention remains required only for exact Requirements Proposal
promotion, exact plan approval, exact publication approval, and the existing
architecture, scope, irreversible-effect, control-weakening, ambiguity, and
budget stop conditions.

## Verification

Automated evidence covers:

- closed request and proposal validation;
- complete-final-draft and material-question behaviour;
- expiry, role, provenance, size, token, latency, and cost rejection;
- purpose and classification denial before adapter invocation;
- normalized provider failure without sensitive response disclosure;
- OpenAI endpoint allow-listing before credential resolution;
- redirect rejection at the provider boundary;
- storage-disabled strict structured-output requests;
- provider-activation budget and receipt binding;
- per-session turn and complete wire-size budget enforcement;
- credential exclusion from request bodies and errors;
- Experience-owned scope, task-agent, and proposal validation;
- multi-turn conversational conversion into the existing Requirements
  Proposal;
- lossless discovery-draft conversion and complete proposal review projection;
- meaningful human confirmations bound internally to current immutable
  subjects;
- the full contract-double conversational path through a reviewable simulated
  pull request; and
- denied cross-scope discovery stopping before repository discovery and model
  invocation.

Repository verification runs formatting, architecture checks, all Go tests,
vetting, race detection, build, documentation links, navigation, decision
registration, capability identifiers, module boundaries, schemas, and existing
M1 journey controls.

Two in-process Intelligence validation and contract-double invocation benchmark
runs on the local Apple M1 Pro measured approximately 4.2 to 8.5 microseconds
per proposal with approximately 4 KiB and 12 allocations per operation. This is a
developer-machine baseline for deterministic Forge overhead, not a
provider-latency or production-capacity claim. Live-provider latency, quality,
cost, concurrency, and degradation remain explicit trial measurements.

The repository-pinned `govulncheck` reports no vulnerabilities reachable from
Forge code. It also reports two vulnerabilities in imported packages and 21 in
required modules whose affected symbols are not called. CI retains the same
scan so changed reachability fails the delivery gate.

## Explicit exclusions and residual risks

This iteration does not activate a provider, make a live model call, commit a
provider credential, implement a durable conversation store, invoke tools,
read repository source, retain unrestricted transcripts, implement a model
registry or router, coordinate several agents, implement an Agent Broker,
create several active virtual Product Managers, use Webex, call GitHub, merge,
release, deploy, access production, or claim Supported intelligence.

The CLI projection remains line-oriented JSON plus prompts rather than the
planned rich terminal experience. Provider quality, real latency, cost,
retention, residency, refusal behaviour, prompt-injection resistance, and
adaptive discovery usefulness remain unvalidated until the first exact
provider activation and human trial.

## ADR assessment

ADR-0022 was accepted before implementation and covers the new external-model,
credential, Intelligence, and Experience boundaries. No additional ADR is
proposed for this iteration. A new decision remains required before general
model routing, tool-enabled models, dynamic Agent Broker implementation,
multiple active virtual Product Managers, model-generated execution, or
another authority model.
