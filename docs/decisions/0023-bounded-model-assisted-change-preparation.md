# ADR-0023: Permit bounded model-assisted change preparation for M1 exit

Status: Accepted
Owner: Forge Product, Intelligence, Execution, Verification, Security, Privacy, and Architecture
Date: 2026-07-24

## Context

Forge M1 now provides a public guided path from free text through adaptive
product discovery, an explicitly promoted Requirements Proposal, a reviewed
Delivery Plan, bounded worker authority, independent verification, separate
publication approval, and an idempotent source-control contract double.

The implemented Developer Trial worker does not yet prepare a meaningful change
from the promoted requirements. Its deterministic profile writes one
predefined synthetic Markdown file. That proves the trust path but not Forge's
central product hypothesis: that Forge can help build software from understood
intent.

ADR-0010 permits at most one model adapter in M1 after separate provider,
classification, budget, retention, and accountable approval. ADR-0017 chose a
deterministic first worker and requires a separately approved M1 adapter or the
later M3 contracts before model-assisted transformation. ADR-0022 pulled
forward one non-authoritative discovery adapter but explicitly withheld model
repository access, execution, tools, and change preparation.

Waiting until the complete M2 and M3 platforms exist would postpone the
essential end-user learning loop until after Forge has hardened an artificial
worker experience. Giving the discovery model a repository or letting it write
the candidate directly would instead collapse Intelligence, Context,
Execution, Delivery, and Verification boundaries.

Forge therefore needs one narrow M1-exit change-preparation adapter that tests
real implementation value without importing general routing, tools, or a
multi-agent workforce.

## Decision

Permit one optional, replaceable, non-authoritative model-assisted change
preparation path in the Developer Trial profile. M1 retains one eligible
external model provider integration. Discovery and change preparation use
separate purpose-specific contracts, activations, budgets and receipts behind
that same provider edge; neither activation authorizes the other purpose.

The path SHALL remain:

~~~text
promoted canonical Intent
  -> bounded reproducible Context Package
  -> deterministic Delivery Plan
  -> human approval of the exact plan digest
  -> Policy-issued model-proposal authority
  -> purpose-minimized model request through Intelligence
  -> schema-valid untrusted Change Proposal
  -> separately Policy-issued worker authority
  -> deterministic validation and isolated materialization by Execution
  -> immutable candidate digest
  -> independent Verification-owned Assurance Plan and result
  -> separate human publication approval
  -> approved source-control adapter
~~~

The model proposes candidate bytes. It does not receive or create authority,
apply changes, invoke tools, execute repository code, approve a plan, verify a
candidate, publish, merge, release, deploy, or alter accepted requirements,
architecture, policy, tests, or expected journey outcomes.

### Capability ownership

- `forge.capability.intent` owns promoted product intent and acceptance
  criteria.
- `forge.capability.delivery` owns the canonical plan, declared change paths,
  lifecycle, budgets, candidate reference, and recovery state.
- `forge.capability.context` owns the authorized reproducible context manifest
  and purpose-minimized source projections.
- `forge.capability.identity` owns human, task-agent, and workload identity and
  phase-specific delegation lifecycle.
- `forge.capability.policy` decides whether model disclosure and change
  preparation are permitted and issues the exact Execution Envelope.
- `forge.capability.intelligence` owns provider eligibility, activation,
  adapter invocation, structured-output admission, measured usage, normalized
  provider failure, and the model receipt.
- `forge.capability.execution` owns isolated deterministic validation and
  materialization of the admitted Change Proposal into candidate bytes.
- `forge.capability.verification` independently selects and evaluates the
  Assurance Plan against the immutable candidate digest.
- `forge.capability.evidence` records minimized attributable references and
  integrity evidence without retaining unrestricted prompts, source, or model
  responses by default.
- `forge.capability.approval` and `forge.capability.source-control` retain their
  existing separate publication responsibilities.
- `forge.capability.experience` presents progress, cost, uncertainty,
  limitations, candidate explanation, approvals, and recovery state.

No capability may read another capability's persistence or reinterpret a model
proposal as authority.

After the exact plan approval, Policy MAY admit the declared model-proposal and
materialization phases without another human micro-approval. It issues
non-overlapping authority for:

- action `forge.intelligence.propose-change`, bound to the task agent, provider
  destination, Context Package, plan, allowed source projections,
  classification, purpose and model budgets; and
- action `forge.worker.prepare-change`, bound to the isolated workload,
  admitted Change Proposal digest, source revision, plan, paths and worker
  limits.

Neither action is implied by plan approval, discovery activation, context
authority or the other phase's envelope. A changed material fact requires
deterministic re-evaluation, not routine human approval.

### Change-preparation contract

Intelligence receives a versioned `ChangePreparationRequest` containing only
server-resolved and policy-admitted data:

- correlation, tenant, workspace, project, repository, run, plan, context, and
  task-agent references;
- exact source revision, Context Package manifest digest, Delivery Plan digest,
  Requirements Record reference, and acceptance-criteria references;
- the allowlisted read paths, write paths, and existing content digests;
- bounded source excerpts or complete files selected by Context for this
  purpose and classification;
- applicable accepted architecture excerpts and immutable control references;
- required change types, output contract, test intent, exclusions, and
  unresolved limitations;
- exact provider activation, model, role profile, prompt, schema, retry,
  expiry, token, latency, compute, and financial limits; and
- declared classification, destination, retention, and training-use posture.

It SHALL NOT contain:

- a credential value, Approval Record, Policy Decision, Execution Envelope, or
  other authority object;
- unrestricted repository content, whole-repository archives, unrelated
  history, ambient chat, hidden memory, environment variables, secrets, or
  production/customer data;
- repository, destination, branch, base, path, tool, model, or budget choices
  controlled by the requester or model;
- executable shell, package-install, network, merge, release, deployment, or
  policy-mutation authority; or
- a provider continuation identifier used as canonical Forge state.

The adapter returns a versioned `ChangeProposal` containing:

- exact request, source, plan, context, role, prompt, schema, model, and
  activation references;
- one bounded set of file replacements or patches, each bound to an allowed
  path and expected source-content digest;
- a requirements-and-plan rationale for each change;
- proposed test changes clearly labelled as candidate-controlled evidence;
- explicit assumptions, uncertainty, omissions, and unsupported requirements;
- normalized incomplete, refusal, or unsafe-output state; and
- measured token, latency, financial-cost, and provider-receipt metadata.

Unknown fields, undeclared paths, stale source digests, oversized output,
unsupported encodings, binary or symlink targets, prohibited instructions,
missing provenance, budget overrun, invalid patch structure, or mismatched
request references fail closed before materialization.

### Deterministic materialization

Execution receives only an admitted Change Proposal plus the existing exact
worker authority. It materializes the proposal in the ADR-0017 isolated worker
using Forge-owned deterministic logic.

The materializer SHALL:

- start from the exact approved immutable source revision;
- revalidate plan, context, delegation, envelope, proposal, path, content
  digest, size, file-count, output, and expiry bindings;
- apply changes only to declared paths without following symlinks;
- run no model-supplied command and interpret no repository instruction as
  authority;
- produce the changed-path manifest, tree digest, candidate digest, resource
  use, cleanup outcome, and safe evidence references; and
- discard or quarantine incomplete output without transitioning Delivery to
  candidate success.

The model provider credential remains inside the Intelligence adapter. It is
never mounted into the worker or verifier.

### Independent verification and human authority

Model-generated code, tests, explanations, and self-review are candidate input,
not verification. Verification derives its Assurance Plan independently from
accepted intent, architecture, risk, changed capabilities and paths, the
Application Quality Envelope, and Forge-owned controls.

Applicable architecture, contract, unit, integration, journey, security,
privacy, accessibility, performance, leakage, recovery, and source-integrity
checks remain publication blockers. Candidate-controlled tests may add evidence
but cannot replace a consequential Forge-owned oracle.

Plan approval authorizes neither the model call nor publication by itself.
Policy separately admits the exact model request and worker operation.
Publication still requires a later human Approval Record v2 bound to the exact
verified candidate digest. Forge remains unable to merge.

### Provider activation and resource selection

No discovery activation automatically permits change preparation. Before the
first live preparation request, an accountable human SHALL approve one exact
Developer Trial activation for purpose `forge.m1.prepare-change` recording:

- provider, adapter, endpoint/region, model identifiers and role-profile
  versions;
- eligible classifications and exact external destination;
- retention, training use, residency, deletion, provider storage and
  subprocessor posture;
- credential reference, project, scope, rotation, revocation and spend limit;
- per-request and aggregate input, output, token, latency, attempt, concurrency,
  compute and financial budgets;
- approved failure, retry, cancellation and no-fallback behaviour; and
- evaluation evidence for patch quality, structured output, source-digest
  binding, prompt injection, data leakage, refusal, latency, cost and
  reproducibility.

Forge selects the smallest evaluated model sufficient for the declared change.
Provider or model prestige is not an eligibility criterion. An expensive model
is used only when measured quality requirements justify it within budget.

Provider failure preserves the deterministic Delivery checkpoint. Bounded
retry may reuse only the same exact request and activation. A changed plan,
context, source, path, provider, model, budget, destination, classification, or
prompt version requires new admission and invalidates the earlier proposal.

### M1 boundary

This decision pulls forward only one model-assisted Engineering preparation
role for one approved Forge change. It does not implement:

- general model registry or dynamic routing;
- model-selected tools, shell, package installation, network access, MCP, or
  secret use;
- multiple engineering agents, virtual Product Manager hand-off orchestration,
  architecture agents, UX agents, agent debate, or general Agent Broker;
- autonomous plan, requirements, architecture, policy, approval, expected
  outcome, verification, publication, merge, release or deployment changes;
- production data, credentials, target, profile or maturity; or
- Workplace-specific implementation.

A deterministic no-model change-preparation path remains available for
regression, recovery, provider outage and adapter conformance.

## Alternatives considered

- **Keep the predefined synthetic transformation through M1 exit:** rejected
  because it proves controls but not that Forge can prepare a useful change
  from understood intent.
- **Wait for all M2 and M3 capabilities:** rejected because the first real
  product feedback would arrive after substantial platform hardening.
- **Let the discovery model write the candidate:** rejected because discovery
  has no repository, execution, or plan authority and uses a different purpose
  and data activation.
- **Give the model a shell or coding-agent toolset:** rejected because M1 has no
  general Tool Broker, dynamic action, secret, egress, or package-install
  contract.
- **Apply a returned patch in the control plane:** rejected because candidate
  materialization belongs in the isolated Execution boundary.
- **Require human approval of every proposed file:** rejected because the human
  already approves the exact plan and paths; independent verification and
  separate publication approval provide the consequential gates.
- **Use one fixed frontier model for every change:** rejected because quality,
  cost, latency, privacy and resource proportionality require measured
  eligibility and smallest-sufficient selection.

## Consequences

- Forge can test its central application-builder hypothesis during M1 exit.
- One additional external data-processing purpose and model request becomes
  possible only after exact activation.
- Source excerpts may leave the local environment when classification and
  activation permit it, increasing minimization, provider, retention, leakage
  and cost obligations.
- Deterministic materialization and independent verification preserve the
  existing authority sequence.
- Failures can pause at a durable checkpoint without repeating product
  discovery or human plan review when immutable inputs remain exact.
- The narrow adapter informs M3 model registry, routing, evaluation and Tool
  contracts without prematurely claiming those capabilities.

## Validation

Evidence SHALL prove:

- model output cannot add authority, repository, source revision, paths, tools,
  models, destinations, classifications, secrets, budgets, checks or actions;
- only Context-authorized, purpose-minimized source with exact provenance and
  content digests reaches an eligible provider;
- malformed, oversized, stale, cross-scope, unknown-field, symlink, binary,
  executable-instruction and prompt-injection fixtures fail closed;
- provider credentials never enter request content, worker/verifier
  environments, candidate bytes, logs, errors, telemetry or Evidence Records;
- exact replay, timeout, rate limit, refusal, incomplete output, cancellation,
  budget exhaustion and provider unavailability preserve or recover the
  correct Delivery state without duplicate candidates;
- deterministic materialization reproduces the candidate digest from the exact
  source and Change Proposal;
- changed paths and resulting tree exactly match the admitted proposal and
  worker envelope;
- candidate-controlled tests and model explanations are never the sole
  consequential oracle;
- independent verification detects seeded requirement, architecture, security,
  privacy, accessibility, performance and journey failures applicable to the
  change;
- the provider adapter can be replaced by a contract double without changing
  Intent, Delivery, Context, Policy, Execution, Verification, Approval,
  Evidence, Source Control or Experience domain contracts;
- p50, p95 and p99 latency, token use, financial cost, output size, failure
  rate, retry, cancellation and quality are measured against the declared
  Developer Trial envelope;
- one allowed public journey prepares a useful real Forge candidate, while
  cross-scope, stale-authority and denied-provider journeys stop before source
  disclosure or mutation; and
- the verified candidate can proceed only through ADR-0019 publication
  authority and remains unable to merge.

## Accountable review

Sijbren Beukenkamp accepted this decision as drafted through the authenticated
collaboration boundary on 2026-07-24, acting as the temporary accountable
Product, Intelligence, Execution, Verification, Security, Privacy, and
Architecture reviewer established for bounded M1.

This acceptance authorizes only the bounded architecture above.
Runtime implementation remains blocked until:

1. ADR-0023 reaches protected `main` as Accepted;
2. the separately activated contributor task covers the affected capabilities
   and implementation paths; the candidate does not modify this Accepted ADR;
3. a separate bounded implementation task is activated against that protected
   revision; and
4. each live provider request receives its own exact purpose-specific
   activation.

The first live GitHub publication additionally remains blocked until the
ADR-0019 GitHub App permission and protected-branch denial review is complete.

## Supersession

If Accepted, this decision specializes ADR-0003, ADR-0005, ADR-0010, ADR-0013,
ADR-0017, ADR-0018, ADR-0019, ADR-0020, ADR-0021, and ADR-0022.

It supersedes only:

- ADR-0010's single-purpose interpretation of the one-provider M1 model
  allowance, so the same one eligible provider edge may expose separately
  activated discovery and change-proposal contracts;
- ADR-0017's deterministic-first restriction, to permit one optional
  model-assisted Change Proposal before deterministic isolated
  materialization; and
- ADR-0022's discovery-only model restriction, to permit a separately
  activated purpose `forge.m1.prepare-change`.

It does not supersede capability ownership, exact plan approval,
phase-specific Policy and Identity authority, isolated execution, independent
verification, exact candidate publication approval, provider eligibility,
evidence integrity, no-merge, no-release, no-deployment, no-production,
no-general-routing, no-general-tools, or no-broad-multi-agent restrictions.
