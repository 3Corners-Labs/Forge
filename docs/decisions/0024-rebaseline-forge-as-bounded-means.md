# ADR-0024: Rebaseline Forge as a bounded means for sovereign application delivery

Status: Accepted
Owner: Forge Product and Architecture
Date: 2026-07-25

## Context

The accepted corpus was written on the assumption that Forge is the product.
Enduring principle 10 states that platforms outlive applications and that no
lighthouse application may become the hidden platform architecture. The
[masterplan](../foundation/masterplan.md) requires Forge to exist as a usable
development-team application before accepting the Workplace commission. The
[roadmap](../roadmap.md) places Workplace Orchestrator at M6 of M9, behind two
Programme-sized milestones.

On 2026-07-25 the accountable owner stated a different intent through the
authenticated collaboration boundary:

- Forge is the **means**, not the product. The objective is to learn whether an
  AI development team, under human accountability, can build **and maintain** a
  real application without a large human engineering organization.
- Workplace Orchestrator is the vehicle for that learning and is intended as an
  open-source community edition.
- The target context is Dutch public-sector organizations moving away from
  United States hyperscaler dependency toward open-source and self-hosted
  operation.
- Available human capacity is deliberately minimal, and the same person remains
  the sole accountable Product, Architecture, Security, Privacy, Engineering,
  and Operations reviewer.

Two verified facts constrain any plan built on that intent.

First, the delivery evidence contradicts the assumed mechanism. Every material
change to this repository — 40 commits and approximately 19,000 lines of
non-test Go — was produced by the **contributor harness** governed by
[ADR-0011](0011-bounded-autonomous-contributor-delivery.md): `/AGENTS.md`, the
context router, the decision register, autonomy grants, and the architecture
fitness checks, executed by an external coding assistant under human review.
The Forge runtime has produced zero real changes. Its source-control adapter is
not live-configured and no model provider is active by default, as recorded in
the [implementation baseline](../governance/implementation-baseline.md). The
working delivery mechanism is the governance layer, not the Go kernel.

Second, sovereignty is now a product requirement rather than a deployment
preference, and it reaches back into Forge itself. An open-source application
promoted as free from hyperscaler dependency, but maintainable only through a
proprietary external model subscription held by one person, reproduces the
dependency it claims to remove one layer higher.

Continuing under the accepted sequencing would defer the stated objective behind
platform work that the available capacity cannot complete, while the mechanism
that already delivers remains unrecognized by the architecture.

## Decision

Rebaseline Forge as a **bounded means for sovereign application delivery**.

### 1. Purpose

Forge's purpose is to deliver and maintain 3Corners applications, beginning with
Workplace Orchestrator, under human accountability and with verifiable evidence.
Forge is not developed as a product for independent sale. Reusability remains a
desirable property; it is no longer a governing constraint.

Enduring principle 10 ("Platforms outlive applications") is demoted to evidence
label **Vision**. It states an aspiration, not a current design constraint.
Forge MAY be shaped by Workplace requirements where that shortens delivery, so
long as the coupling is recorded and reversible.

### 2. Delivery mechanism

The contributor harness is recognized as the **primary delivery mechanism** for
the current stage, not merely a bootstrap workaround. It comprises the operating
contract, context router, decision register, autonomy grants, fitness checks,
journey assurance, and an external coding assistant under human review.

The Forge runtime SHALL be extended only where a Workplace delivery need
demonstrates value. Runtime capability is no longer built ahead of demand.
Capabilities already implemented remain at their recorded maturity; none is
withdrawn.

This does not weaken any control. Separation of duties, independent
verification, evidence, and human publication authority continue to apply to
both mechanisms.

### 3. Milestone scope

Milestones M4 through M9 are demoted to evidence label **Vision**. They record
intended direction and SHALL NOT be treated as a committed plan, a dependency
for current work, or a justification for building capability ahead of need.

M1 is completed to its accepted exit criteria, because it is the only way to
learn whether the trust path holds against reality, and because an unfinished
kernel is worth less than a demonstrated one.

M2 and M3 are reduced to **demand-driven scope**: only the supply-chain and
intelligence capability that Workplace delivery actually requires.

Workplace Orchestrator moves from M6 to the **next delivery objective** after
M1 exit. Its feasibility spike MAY begin immediately and in parallel, because it
tests a product-existence assumption rather than a Forge capability.

### 4. Application boundary

Workplace Orchestrator SHALL be developed in a separate repository with its own
lifecycle, license, and decision register. Forge modules SHALL NOT contain
workplace domain logic, and Workplace SHALL NOT depend on unreleased Forge
runtime internals.

This preserves the one part of principle 10 that still earns its cost: a
separate repository is the only honest test of whether the harness can deliver
into a codebase it did not generate.

### 5. Sovereignty tiers

Three distinct tiers replace the previous undifferentiated treatment of model
dependency:

| Tier | Requirement |
|---|---|
| Workplace runtime | SHALL function completely without any external model provider. Model-assisted features are optional enhancements that degrade to deterministic behaviour. |
| Forge build-time, accountable owner | MAY use an external cloud model. This is a tool choice by the accountable engineer, not a product dependency. |
| Forge build-time, community contributor | SHALL be operable against a locally hosted model or small language model for at least one complete delivery path, before any claim that the application is community-maintainable. |

The third tier makes model-capability routing load-bearing earlier than M3
assumed. It does not authorize the general M3 model registry. Every Forge task
SHALL record which model class was sufficient for it, so that the locally
hosted path is measured rather than asserted.

No claim of community maintainability, sovereignty, or hyperscaler independence
SHALL be published before tier-three evidence exists.

### 6. Approval proportionality

The current gate model assumes a review organization that does not exist. A
single accountable owner with minimal available hours is a bottleneck in their
own control system, which degrades review quality rather than improving it.

Consequential gates SHALL be preserved without change: architecture decisions,
scope expansion, protected-branch landing, control weakening, accepted-oracle
changes, production effects, destructive actions, and publication of an exact
verified candidate.

Routine gates MAY be satisfied by deterministic checks with evidence rather than
synchronous human attention, provided the evidence remains reviewable after the
fact and any failure fails closed. Reducing a routine gate SHALL be recorded;
silent bypass remains prohibited.

The concentration of accountable roles in one person remains an accepted,
recorded risk. It is mitigated by deterministic verification and by explicit
adversarial challenge, not by pretending independent review exists.

### 7. Learning measurement

The central question — whether an AI development team can build and maintain a
real application under human accountability — SHALL be measured against declared
criteria, recorded in
[AI delivery experiment measurement](../governance/ai-delivery-experiment.md).
Without falsifiable criteria the exercise produces an impression rather than a
result.

## Alternatives considered

- **Continue the accepted Forge-first sequence:** rejected because it places the
  stated objective behind two Programme-sized milestones that the declared
  capacity cannot complete, and because it keeps building runtime capability
  that no consumer has requested.
- **Abandon the Forge runtime and use the harness alone:** rejected because the
  runtime already encodes a trust path that the harness only approximates, and
  because abandoning it before M1 exit destroys the evidence that would tell us
  whether that trust path works.
- **Build Workplace inside the Forge repository:** rejected because it removes
  the only honest test of harness generality and permanently entangles domain
  logic with platform code.
- **Keep principle 10 as a binding constraint:** rejected because it forces
  premature generalization for a second application that does not exist and is
  not funded.
- **Defer the sovereignty requirement to a later milestone:** rejected because
  it is the product's differentiator, and because retrofitting model
  independence after the architecture assumes a cloud model is substantially
  more expensive than designing for it.
- **Remove approval gates entirely to match available capacity:** rejected
  because the consequential gates are the product thesis; removing them would
  make the experiment answer a different and less interesting question.

## Consequences

- The accepted masterplan, principles, and roadmap require revision. Those
  revisions are consequences of this decision and SHALL NOT precede its
  acceptance.
- The [roadmap structure check](../../scripts/check_architecture.py) currently
  requires nine milestones in order. It SHALL be amended alongside the roadmap
  revision, as a recorded change rather than a silent relaxation.
- Forge gains an explicit consumer. Capability decisions become answerable to a
  concrete delivery need instead of an imagined platform market.
- Measuring model-class sufficiency per task adds recording overhead to every
  Forge delivery, and is the price of the community-maintainability claim.
- Demoting M4 through M9 removes the appearance of a committed multi-year plan.
  This is intended: the previous roadmap implied commitments the capacity could
  not honour.
- The single-reviewer risk becomes explicit and measured rather than noted and
  deferred.

## Validation

Evidence SHALL prove:

- Workplace Orchestrator exists in a separate repository with an independent
  decision register, and Forge modules contain no workplace domain logic;
- the M1 exit criteria are met against a real repository operation rather than a
  contract double, or the failure is recorded with equal precision;
- at least one complete Forge delivery path executes against a locally hosted
  model, with measured quality, latency, and cost compared to the cloud path,
  before any community-maintainability claim is published;
- every Forge delivery records the model class that proved sufficient;
- each reduced routine gate has a deterministic check, retained evidence, and a
  demonstrated fail-closed path;
- the declared learning criteria are recorded before delivery begins and
  evaluated afterwards, including the criteria that would falsify the approach;
- no accepted control, separation of duty, or independent oracle is weakened by
  this rebaseline.

## Accountable review

Sijbren Beukenkamp accepted this decision as drafted through the authenticated
collaboration boundary on 2026-07-25, acting as the temporary accountable
Product, Architecture, Security, Privacy, Engineering, and Operations reviewer
established for bounded M1.

The acceptance explicitly covered the three points on which the drafting
contributor requested a separate judgement: completing M1 before Workplace
delivery begins (section 3), the locally hosted model requirement for community
contributors (section 5, tier three), and the falsification criteria (section 7).

This decision revises the framing under which the accepted corpus was written.
The drafting contributor recorded the acceptance; it did not make it, and no
autonomy grant authorizes a contributor to accept an architecture decision.

The following changes are now authorized and SHALL be delivered as reviewable
changes:

1. revision of [`docs/roadmap.md`](../roadmap.md) and the roadmap structure
   check;
2. revision of [`docs/foundation/masterplan.md`](../foundation/masterplan.md)
   sections 4 and 18, and the evidence label of principle 10 in
   [`docs/foundation/principles.md`](../foundation/principles.md);
3. revision of
   [`docs/reference-applications/workplace-orchestrator.md`](../reference-applications/workplace-orchestrator.md)
   to reflect the separate repository and the sovereignty tiers;
4. an update to the
   [implementation baseline](../governance/implementation-baseline.md) recording
   the changed framing.

## Supersession

If Accepted, this decision specializes
[ADR-0006](0006-accept-minimum-architecture-baseline.md),
[ADR-0007](0007-sequence-self-building-before-full-advisory.md),
[ADR-0008](0008-accept-canonical-documentation-baseline.md),
[ADR-0010](0010-authorize-bounded-m1-implementation.md), and
[ADR-0011](0011-bounded-autonomous-contributor-delivery.md).

It supersedes only:

- the masterplan's Forge-first sequencing, so that Workplace delivery follows M1
  exit rather than M5;
- the treatment of enduring principle 10 as a binding design constraint; and
- the roadmap's presentation of M4 through M9 as committed scope.

It does not supersede the constitution, capability ownership, authority and
context separation, isolated execution, independent verification, exact
publication approval, evidence integrity, or the no-merge, no-release,
no-deployment, and no-production restrictions.
