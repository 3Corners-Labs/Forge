# ADR-0021: Permit a guided resumable M1 delivery session

Status: Accepted
Owner: Forge Product, Experience, Delivery, and Architecture
Date: 2026-07-23

## Context

ADR-0009 established `forge delivery start` as the public M1 boundary and
required its initial submission not to imply plan approval, worker authority,
or publication authority. ADR-0020 subsequently established a Requirements
Proposal and explicit promotion boundary. The accepted M1 allowed-journey
oracle version 2.0.0 now requires the public experience to continue from an
explicitly promoted proposal through separately approved planning, bounded
execution, independent verification, separate publication approval, and an
idempotent contract-double pull request.

Treating the initial free-text submission as authority would violate the
Experience and Intent boundary. Requiring users to leave the public journey and
call internal capabilities would violate the public-interface assurance
boundary. The journey therefore needs one guided coordinator without collapsing
its independent authority gates.

## Decision

`forge delivery start` MAY coordinate a resumable guided M1 session through the
existing public capability contracts. Its initial submission creates only a
digest-bound Requirements Proposal. The session SHALL expose and require three
separate human actions:

1. promotion of the exact current Requirements Proposal digest;
2. approval of the exact Delivery Plan through Approval Record v2; and
3. approval of publication of the exact independently verified candidate
   digest through a separate Approval Record v2.

The coordinator SHALL NOT manufacture scope, policy, approval, execution,
verification, evidence, or publication authority. It resolves every current
record from its owning capability and presents opaque resumable references
rather than treating terminal input, a session checkpoint, or a channel message
as canonical state.

Long-running work MAY continue after the applicable explicit gate and SHALL
provide immediate acknowledgement and stage progress. Cancellation, expiry,
changed digests, stale checkpoints, unresolved consequential requirements,
failed verification, unavailable evidence, and ambiguous provider state fail
closed through their owning contracts.

For M1, publication is restricted to the deterministic credential-free
source-control contract double. Live GitHub credentials, ambient developer
credentials, merge, review approval, release, deployment, production data, and
production effects remain forbidden. The final result states that human review
and merge remain outside Forge authority.

This decision supersedes only ADR-0009's stopping restriction. ADR-0009's
command name, server-derived scope, normalized results, and prohibition on
authority bypass remain Accepted.

## Consequences

- A user can experience the complete M1 trust path without invoking internal
  services or hidden test endpoints.
- The coordinator adds no new canonical owner; capability records remain the
  source of truth.
- A later Workbench or Webex adapter can use the same session and capability
  contracts.
- Durable multi-process resume remains dependent on the accepted capability
  persistence adapters and opaque server-side references, not client files as
  authority.

## Validation

Evidence SHALL prove:

- raw input stops at a Requirements Proposal until its exact digest is
  explicitly promoted;
- plan and publication approvals are distinct, exact, current, human decisions;
- execution starts only after plan approval and publication starts only after
  independent verification, evidence readiness, Policy authority, and separate
  publication approval;
- plan approval cannot be replayed for publication;
- the contract double creates at most one branch and pull request for one
  deterministic operation, including ambiguous-result reconciliation;
- the denied cross-scope journey fails before repository discovery;
- every wait state is keyboard and screen-reader legible and emits immediate
  progress;
- live credentials, merge, release, deployment, and production remain
  unreachable.

## Accountable review

Sijbren Beukenkamp accepted this decision as proposed through the authenticated
collaboration boundary on 2026-07-23.

## Supersession

This decision specializes ADR-0005, ADR-0009, ADR-0016, ADR-0019, and ADR-0020.
It does not change capability ownership, authorize a new provider, or weaken an
existing human gate.
