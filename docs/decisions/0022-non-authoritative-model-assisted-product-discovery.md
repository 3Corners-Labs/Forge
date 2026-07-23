# ADR-0022: Permit non-authoritative model-assisted product discovery

Status: Accepted
Owner: Forge Product, Experience, Intelligence, Security, Privacy, and Architecture
Date: 2026-07-23

## Context

ADR-0020 establishes guided product discovery in which free text, chat messages,
mockups, retrieved content, and model output are untrusted interaction input.
It requires an exact Requirements Proposal to be promoted by an authenticated
human before canonical intent or consequential delivery work exists. ADR-0021
lets the public M1 journey continue through later independent gates without
giving its coordinator authority owned by another capability.

The deterministic M1-C trial proves that authority separation, proposal
promotion, plan approval, verification, and publication approval can work. It
does not provide the adaptive product discovery needed for Forge to behave as a
professional virtual Product Manager. Scripted questions cannot reliably
recognize ambiguity, contradictions, stakeholder needs, product-area
dependencies, UX uncertainty, or the next question that most reduces delivery
risk.

The Accepted Agentic Delivery architecture already defines Product, UX,
Architecture, Engineering, Quality, Security, Delivery, and Operations agent
roles. The Accepted Intelligence architecture already assigns provider-neutral
model and method selection to `forge.capability.intelligence`. The roadmap,
however, sequences general model routing in M3 and broad specialist-agent
collaboration in M5. ADR-0010 permits at most one model adapter during M1 only
after provider eligibility, classification, budget, retention, and accountable
approval, and explicitly excludes broad multi-agent orchestration.

Forge therefore needs a narrow trial that tests real adaptive discovery now
without silently importing the later general Intelligence or Workforce
platforms into M1.

## Decision

Permit one non-authoritative, provider-replaceable model-assisted product
discovery path in the Developer Trial profile.

The path SHALL remain:

~~~text
untrusted interaction
  -> Experience-owned interaction draft
  -> bounded discovery request
  -> Intelligence-owned model adapter
  -> schema-valid discovery-turn proposal
  -> deterministic validation and Experience projection
  -> Requirements Proposal
  -> explicit human promotion of the exact proposal digest
  -> canonical Intent and the existing governed delivery journey
~~~

This decision pulls forward only the minimum Intelligence responsibility needed
to evaluate the conversational product experience. It does not authorize
general model routing, model-generated execution, tools, repository access,
dynamic agent delegation, autonomous approval, or broad multi-agent
orchestration.

### Virtual Product Manager role

The virtual Product Manager is a versioned behavioural role, not a new authority
boundary or capability owner. It MAY:

- clarify outcomes, users, workflows, constraints, success measures, and
  acceptance criteria;
- identify contradictions, assumptions, hypotheses, missing stakeholders, and
  consequential unresolved questions;
- propose functional, UX, accessibility, privacy, security, performance,
  operational, and reversibility considerations;
- propose bounded exploratory mockups or learn-and-discard experiments;
- propose changes to a Requirements Proposal and explain the reason and
  uncertainty; and
- summarize requirement coverage from structured traceability evidence.

It SHALL NOT:

- resolve identity, scope, repository binding, policy, provider eligibility, or
  authority;
- promote a Requirements Proposal or change canonical Intent;
- approve a plan, candidate, publication, release, deployment, policy, or
  architecture decision;
- invoke a tool, inspect a repository, mutate state outside its draft, or create
  a Delivery run;
- assert that implementation satisfies a requirement without independent
  Verification-owned evidence; or
- hide uncertainty, invent provenance, or represent its own narrative as proof.

Prompt instructions such as role, tone, and discovery method guide model
behaviour but are not security controls. Forge-owned schemas, policy decisions,
budgets, capability boundaries, lifecycle validation, and human gates remain
the controls.

### Discovery contract

`forge.capability.experience` owns the interaction session, user-visible
conversation, draft projection, and presentation of the resulting proposal.
`forge.capability.intent` remains the only owner of promoted canonical intent.
`forge.capability.intelligence` owns model eligibility, adapter invocation,
model and prompt-version references, bounded retry, refusal handling, usage
receipt, and normalized provider failure.

Experience invokes Intelligence through a versioned
`DiscoveryTurnRequest`. The request SHALL contain only server-resolved,
purpose-minimized data:

- session, correlation, tenant, workspace, project, human, and task-agent
  references;
- the versioned virtual Product Manager role-profile reference;
- current interaction-draft and Requirements-Proposal references and digests;
- direct user statements, declared classifications, and authorized context
  excerpts with provenance;
- product-area references when already known;
- unresolved questions, assumptions, hypotheses, constraints, and exclusions;
- the allowed response schema and requested discovery purpose; and
- token, latency, cost, retry, and expiry limits.

It SHALL NOT contain a credential value, authority object, Approval Record,
secret, unrestricted transcript, undeclared attachment, whole repository, or
provider-controlled continuation identifier as canonical state.

The adapter returns a `DiscoveryTurnProposal` with:

- a user-facing response;
- zero or more prioritized material questions with the consequence of remaining
  unanswered;
- proposed direct-statement, inference, assumption, hypothesis, constraint,
  acceptance-criterion, and product-area changes;
- proposed exploratory artefacts and the hypothesis each would test;
- uncertainty and limitation declarations;
- input provenance references;
- model, adapter, role-profile, prompt, and response-schema versions;
- normalized refusal or incomplete-result state; and
- measured token, latency, and cost usage.

The output is untrusted even when a provider guarantees JSON-schema
conformance. Forge SHALL validate schema, closed enums, identifiers, sizes,
classification, provenance references, budgets, and prohibited fields before
Experience can use it. A valid output remains a proposal and grants no
authority.

### Provider and data eligibility

No provider becomes eligible merely because this ADR is Accepted. Before the
first live request, an accountable human SHALL approve one exact Developer Trial
provider activation that records:

- provider, endpoint region, adapter version, and exact allowed model
  identifiers;
- eligible purposes and classifications;
- retention, training-use, residency, subprocessor, and deletion posture;
- prompt and response storage behaviour;
- credential reference, scope, rotation, revocation, and spend limit;
- request, session, project, and daily token, latency, concurrency, and
  financial budgets;
- approved fallback and failure behaviour; and
- evaluation evidence for discovery quality, prompt injection, privacy,
  structured output, refusal, latency, and cost.

The first trial SHALL use synthetic or explicitly approved non-production data.
Restricted content, secrets, repository source, customer data, production data,
attachments, tools, and provider-side conversation storage are forbidden.
Requests SHALL disable optional provider application-state storage where the
provider supports that control. Provider failure pauses or degrades discovery;
it cannot corrupt session state or bypass proposal promotion.

OpenAI is the preferred first adapter candidate for evaluation, not an eligible
provider by this decision alone. Its current API documentation states that API
data is not used for training unless a customer opts in, while default abuse
monitoring may retain prompts and responses for up to 30 days. Zero Data
Retention and Modified Abuse Monitoring require separate eligibility and
approval. The activation record SHALL reflect the actual Forge API project
configuration rather than relying on a general provider statement.

### Agent Broker and multiple virtual Product Managers

The target architecture MAY assign virtual Product Managers to different
product areas and route bounded work to Product, UX, Architecture, Engineering,
Security, Quality, and other specialist agents through
`forge.capability.workforce`.

The first trial authorizes one active virtual Product Manager role per discovery
session. Its contracts SHALL use explicit role-profile, task-agent, product-area,
and parent-task references and SHALL NOT encode a global singleton Product
Manager. This preserves later composition without implementing it early.

A future Agent Broker belongs to Workforce coordination and uses Identity-owned
delegation, Policy-owned decisions, Intelligence-owned model selection,
Tool-owned brokerage, and Delivery-owned work state through their public
contracts. It SHALL route work but cannot mint authority, widen a task,
self-select an approval, or let one agent become the sole verifier of its own
consequential output.

General Agent Broker implementation, dynamic role discovery, multiple active
virtual Product Managers, model-to-model hand-off, tool selection, child
delegation, and specialist-agent execution remain Research under this decision.
Their implementation requires the applicable M3 and M5 contracts and an
Accepted decision if milestone sequencing or trust boundaries change.

### Conversational experience

The public boundary remains `forge delivery start`. The conversational CLI
SHOULD:

- accept natural language and provide immediate acknowledgement;
- ask only material questions and state safe reversible assumptions;
- distinguish user statements, Forge inferences, assumptions, hypotheses, and
  unresolved questions;
- allow proposal revision and exploratory mockup loops without micro-approval;
- show meaningful human actions instead of requiring digest transcription;
- keep digests, model identifiers, usage, and provenance available as technical
  details;
- support pause, cancel, resume, keyboard use, screen readers, and a linear-log
  projection; and
- label simulation, contract-double, and live provider effects unambiguously.

CLI, a later Workbench, and a later Webex adapter SHALL use the same Experience,
Discovery, Requirements Proposal, and promotion contracts.

## Alternatives considered

- **Keep scripted discovery until M3:** rejected because Forge would validate a
  false product experience and could discover foundational interaction flaws
  only after hardening the wrong CLI.
- **Treat a role prompt as the guardrail:** rejected because prompts and model
  output are untrusted and cannot enforce authority, privacy, budget, or
  capability boundaries.
- **Let the model write canonical requirements directly:** rejected because it
  bypasses provenance, exact human promotion, and Intent ownership.
- **Implement the full Agent Broker and several virtual Product Managers now:**
  rejected because it expands M1 into the M3 and M5 platforms before the narrow
  discovery contract has user evidence.
- **Bind Experience directly to one provider SDK:** rejected because provider
  concerns belong behind Intelligence and would contaminate the channel and
  product contracts.
- **Use only a local model:** rejected as a universal choice because quality,
  hardware, latency, privacy, and deployment profiles differ; a local adapter
  remains a later eligible option behind the same contract.
- **Let the virtual Product Manager certify its own requirement coverage:**
  rejected by the independent-trust principle; it may explain
  Verification-owned evidence but cannot replace it.

## Consequences

- Forge can test a genuinely adaptive virtual Product Manager before completing
  the general Intelligence and Workforce milestones.
- The trial introduces one new external-processing and credential boundary,
  but keeps it disabled until an exact provider activation passes accountable
  review.
- Experience gains a conversational path without owning model selection or
  canonical product state.
- Intelligence gains only a narrow Trial contract and adapter; it does not yet
  claim the M3 router, registry, evaluation service, or Supported maturity.
- Contract fields support future product-area ownership and role composition
  without implementing multiple active agents.
- Provider unavailability affects discovery quality but not administration,
  inspection, policy, recovery, or existing deterministic journeys.
- The implementation must add model-quality, privacy, prompt-injection,
  performance, cost, accessibility, and degradation evidence.

## Validation

Evidence SHALL prove:

- raw interaction and every model result remain non-authoritative until exact
  human proposal promotion;
- malformed, oversized, unknown-field, cross-scope, stale, and
  provenance-invalid model outputs fail closed;
- prompt injection cannot add authority, tools, models, destinations, secrets,
  repository access, execution, or approval;
- ineligible classifications and destinations are blocked before provider
  disclosure;
- secrets, Restricted data, repository source, and unrestricted transcripts do
  not enter provider requests, logs, errors, or evidence;
- refusal, timeout, rate limiting, malformed output, budget exhaustion, and
  provider unavailability preserve resumable deterministic session state;
- exact model, adapter, role profile, prompt, schema, input references, usage,
  and provider receipt are attributable without retaining unrestricted prompt
  content as evidence;
- model-assisted proposals distinguish statements, inferences, assumptions,
  hypotheses, questions, and experiments and cannot invent valid provenance;
- conversation feedback is immediate, material gates remain distinct, and the
  linear projection is keyboard and screen-reader legible;
- the provider adapter can be replaced by a contract double without changing
  Experience, Intent, Policy, Approval, Delivery, Verification, or Evidence
  contracts;
- one virtual Product Manager can explain requirement-to-evidence coverage but
  cannot create or alter Verification-owned results; and
- the deterministic no-model M1 journey remains available for inspection,
  recovery, regression, and provider-outage testing.

## Accountable review

Sijbren Beukenkamp accepted this decision as drafted through the authenticated
collaboration boundary on 2026-07-23. Acceptance authorizes only the bounded
architecture above. The first live provider request still requires the separate
exact provider activation described in this decision.

## Supersession

If Accepted, this decision specializes ADR-0003, ADR-0005, ADR-0010, ADR-0013,
ADR-0020, and ADR-0021. It supersedes only ADR-0010's M1 sequencing restriction
to permit one non-authoritative model-assisted discovery adapter and the
corresponding conversational Experience trial.

It does not supersede capability ownership, human proposal promotion, plan or
publication approval, independent verification, the no-production boundary,
M3 general model and tool governance, or M5 broad specialist-agent
collaboration.

## Provider references

- [OpenAI API data controls](https://developers.openai.com/api/docs/guides/your-data)
- [OpenAI structured model outputs](https://developers.openai.com/api/docs/guides/structured-outputs)
