# ADR-0020: Establish guided intent and exploratory feedback boundaries

Status: Accepted
Owner: Forge Product, Experience, Intent, and Architecture
Date: 2026-07-23

## Context

Forge must learn what a person is trying to achieve before it spends worker,
model, tool, or human attention on implementation. A free-text request is useful
conversation input but may be incomplete, contradictory, technically
prescriptive without explaining the outcome, or unsafe to treat as authority.
Making the message itself canonical intent would let ambiguous chat text trigger
consequential work. Requiring a complete specification before Forge can respond
would instead reproduce a heavyweight requirements process and prevent rapid
learning.

Users also need low-cost ways to explore behaviour through sketches, live
mockups, prototypes, and iterative learn-and-discard loops. Those artefacts are
valuable evidence of learning but are not automatically accepted requirements,
architecture, verification or production software.

ADR-0009 fixes the M1 public CLI boundary. The Accepted Experience architecture
requires Forge to explain what the user asked, what Forge understood and
assumed, alternatives, limitations and evidence. The Accepted first-slice
contracts already assign user input to Experience, canonical intent and delivery
state to Intent and Delivery, and provider receipts to a source-control adapter.
This decision closes the lifecycle between those boundaries without selecting a
model, chat provider or Workbench framework.

The capability catalogue currently lacks a stable identifier for the already
accepted Source Control owner used by M1-C and ADR-0019. That omission must be
corrected before implementation so the publication adapter is not hidden inside
Delivery, Tool or Experience.

## Decision

Forge uses a guided, channel-neutral interaction lifecycle:

~~~text
untrusted input
  -> interaction draft
  -> requirements proposal
  -> clarification or exploration
  -> explicit human promotion
  -> canonical intent
  -> plan, execution, verification and publication gates
~~~

### Interaction input is not authority

CLI text, Workbench input, chat messages, Webex messages, GitHub issues,
attachments, transcripts, repository content, retrieved content, mockups and
model output SHALL be treated as untrusted interaction input.

Experience MAY acknowledge, normalize and structure that input as an
Interaction Draft. An Interaction Draft SHALL NOT:

- establish tenant, workspace, project or repository scope;
- approve a plan, candidate, publication, release or deployment;
- grant a model, tool, secret, destination or worker authority;
- change accepted requirements, architecture, policy or expected outcomes; or
- trigger consequential execution or external mutation.

Authenticated identity and Tenant-owned selectors are resolved independently.
Conversation channels transport interaction; they do not become a source of
truth or an authorization boundary.

### Requirements Proposal

Experience produces a versioned, digest-bound Requirements Proposal containing
at least:

- the user outcome and affected users;
- the problem or opportunity as Forge currently understands it;
- functional behaviour and visible UX expectations;
- accessibility, privacy, security and operational considerations known at the
  time;
- constraints, exclusions and acceptance criteria;
- assumptions and supporting evidence references;
- unresolved material questions and the consequence of each unanswered
  question;
- proposed experiments or mockups;
- provenance, lifecycle, version, freshness and correlation identifiers; and
- an explanation of what changed from the previous proposal.

The proposal distinguishes observed facts, direct user statements, Forge
inferences and unresolved hypotheses. A model-generated narrative alone is not
evidence.

Forge SHOULD ask only questions whose answers can materially change user
outcomes, scope, architecture, data handling, risk, cost, accessibility,
acceptance or reversibility. It SHOULD state a safe proposed assumption when
progress can continue reversibly and expose that assumption for correction.
Unresolved questions whose declared consequence is consequential SHALL block
promotion or the affected later gate.

### Promotion to canonical intent

Only an authenticated human acting within resolved scope MAY explicitly promote
one exact Requirements Proposal digest. Promotion creates a new canonical
Intent-owned version and Delivery run through the normal `SubmitIntent` and
`CaptureIntent` boundary.

Promotion is not plan approval, worker authority, candidate acceptance or
publication approval. Later material feedback creates a new proposal and, when
promoted, a new intent version that invalidates derived plans and authority as
required. Editing chat history or a source issue never mutates canonical intent.

For M1-C the promotion interaction is deterministic and explicit through the
public CLI. Forge SHALL NOT imply model understanding when no approved model
adapter participated.

### Exploration and learn-and-discard loops

Forge MAY create bounded Exploration Artefacts such as wireframes, clickable
mockups, synthetic previews, alternative flows or disposable prototypes before
canonical implementation.

Every Exploration Artefact SHALL declare:

- lifecycle `exploratory`;
- the proposal digest and hypothesis it explores;
- synthetic or explicitly authorized data;
- expiry and cleanup behaviour;
- known limitations and prohibited uses; and
- that it is not accepted requirements, production software or verification
  evidence.

Discarding an exploration SHALL NOT delete promoted intent or required evidence.
Keeping or promoting learning from an exploration requires an explicit human
action that creates a new versioned Requirements Proposal or another record
owned by the applicable capability. Artefact bytes, chat history or visual
appearance cannot silently become canonical architecture or expected outcomes.

### Experience and channel evolution

The M1 interface remains `forge delivery start` under ADR-0009. The CLI and any
later Workbench or chat adapter SHALL invoke the same public Experience and
capability contracts. A channel MAY resume a draft through a safe opaque
reference, but server-side resolution determines the human, scope, current
proposal and permitted actions.

Webex Messaging is the intended first external chat candidate, but its bot
identity, webhook, card, OAuth, tenant mapping, attachment, retention and
availability choices require a separate provider decision when implementation
is near. This ADR does not bind Forge core contracts to Webex.

### Source Control capability registration

Register `forge.capability.source-control` as the stable Platform capability for
provider-neutral source-control publication and reconciliation. It owns
publication operation identity, provider operation state, provider receipts and
reconciled remote references.

Source Control SHALL NOT own canonical intent, Delivery state, Policy Decisions,
Approval Records, Verification Results, Evidence Records, repository scope or
credentials. ADR-0019 remains the Developer Trial GitHub adapter decision.

## Alternatives considered

- **Treat the first free-text message as canonical intent:** rejected because
  ambiguity, prompt injection and accidental wording would gain authority.
- **Require a complete specification before Forge responds:** rejected because
  it prevents iterative discovery and moves professional clarification work
  back to the user.
- **Let chat history or GitHub issues become the source of truth:** rejected
  because edits, deletion, retention, identity and threading semantics are
  channel-controlled.
- **Treat a successful mockup as accepted design or verification:** rejected
  because prototypes optimize learning and may omit architecture, accessibility,
  security, data, failure and operational behaviour.
- **Ask for human confirmation after every clarification or reversible step:**
  rejected because explicit proposal promotion and later consequential gates
  are sufficient; ordinary draft refinement does not require micro-approval.
- **Place GitHub publication inside Delivery or Tool:** rejected because
  provider operation state and reconciliation have a distinct accepted owner
  and replacement boundary.

## Consequences

- Forge can feel conversational while keeping authority explicit and reviewable.
- Users see what Forge thinks it understood before implementation begins.
- Requirements may evolve through fast disposable experiments without silently
  changing canonical state.
- Experience must maintain draft and proposal lifecycle, provenance,
  explanations and invalidation projections.
- M1-C gains additional deterministic contract and journey tests but no model or
  Webex dependency.
- A later model-assisted requirements facilitator must obey the same proposal
  and promotion boundary.
- A later Workbench can improve presentation without becoming a second source of
  truth.

## Validation

Evidence SHALL prove:

- unknown, cross-scope and unauthenticated interaction references fail closed
  before repository discovery;
- raw text, chat messages, attachments, mockups and model output cannot directly
  create a Delivery run or authority object;
- proposal digests change when normative content changes and remain stable under
  permitted canonical ordering;
- facts, user statements, inferences, assumptions, hypotheses and unresolved
  questions remain distinguishable;
- consequential unresolved questions block promotion while reversible
  assumptions remain visible and correctable;
- only the exact current proposal digest can be promoted by the authenticated
  human within server-resolved scope;
- changed promoted intent invalidates derived plans, approvals and later
  authority;
- exploratory artefacts expose lifecycle, hypothesis, data posture, expiry,
  cleanup, limitations and prohibited-use labels;
- discarded explorations cannot remove promoted intent or required evidence;
- CLI allowed and denied journeys use the same public Experience boundary; and
- a source-control contract double can replace the GitHub adapter without
  changing Experience, Intent, Delivery, Policy, Approval, Verification or
  Evidence domain contracts.

## Supersession

This decision specializes ADR-0003, ADR-0005, ADR-0009, ADR-0013, ADR-0015,
ADR-0016 and ADR-0019. It corrects the capability-catalogue omission for the
already accepted Source Control owner. It does not select a model, Webex
provider, Workbench framework, production interaction store or general
application-generation contract.
