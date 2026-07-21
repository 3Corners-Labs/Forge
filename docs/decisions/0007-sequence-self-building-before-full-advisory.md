# ADR-0007: Sequence the self-building kernel before full advisory capability

Status: Accepted
Owner: Forge Product and Architecture
Date: 2026-07-21

## Context

The Proposed mega masterplan sequences a broad read-only “Forge observes and advises” stage before Forge prepares changes to itself. The implementation-roadmap request explicitly requires the roadmap to start with the minimum platform capable of building itself.

The initial roadmap rewrite at commit `01827ed7c695f09ae7d49c8a10e2dce082a3ae9a` moved change preparation to M1 and the general Software System Twin to M4, but did not preserve the original foundation work or explain this sequencing difference. That made the roadmap appear to conflict silently with the masterplan and placed some readiness deliverables after the milestone that depends on them.

Sijbren Beukenkamp approved correcting the roadmap and, under the temporary multi-role authority recorded in ADR-0006, approved the sequencing decision on 2026-07-21.

## Decision

Use the following implementation order:

1. Complete F0, the non-functional foundation readiness gate.
2. Implement M1 as the minimum trustworthy self-building kernel that prepares one Forge pull request.
3. Expand reproducibility, supply-chain, intelligence, and tool capabilities in M2 and M3.
4. Implement the general Software System Twin and multi-audience advisory product in M4.

M1 includes only the read-only repository inspection, authorized context assembly, bounded planning, and evidence-backed explanation necessary for one governed Forge change. It does not implement the general knowledge projection, multi-system Twin, or broad advisory workforce assigned to M4.

F0 preserves the former Horizon 0 and masterplan M0 work with explicit status and evidence. Mandatory readiness controls cannot be deferred to a later milestone merely because their mature platform implementation belongs there.

## Alternatives considered

- **Implement the full advisory product before any self-change:** rejected for the implementation roadmap because it delays the accepted intent-to-pull-request proof and conflicts with the requested starting outcome.
- **Move the full Software System Twin into M1:** rejected because it expands the first slice beyond ADR-0006 and increases knowledge, UX, and projection risk before the trust path is proven.
- **Leave the conflict implicit:** rejected because Proposed masterplan text would appear to define a different delivery order without rationale or authority.
- **Delete the original foundation and advisory intent:** rejected because roadmap restructuring must preserve provenance and product intent.

## Consequences

- The accepted first-slice trust path remains the first implementation outcome.
- M1 must provide enough inspection and explanation to prepare a trustworthy change, but cannot claim the M4 advisory or System Twin capabilities.
- F0 becomes an explicit pre-M1 gate and is not counted as a product implementation milestone.
- The Proposed masterplan must be reconciled with this decision before it is accepted.
- Architecture checks, journey contracts, and offensive-assurance safety contracts have minimum F0 deliverables even though later milestones expand them.
- The roadmap remains Proposed as a whole; this ADR accepts only the sequencing and boundary decision.

## Validation

- Every former Horizon 0 and masterplan M0 responsibility maps to F0 or an explicitly justified later milestone.
- M1 acceptance evidence contains no claim of a general Software System Twin or broad advisory workforce.
- M4 can be implemented without changing M1 authority, delivery, evidence, or source-control contracts.
- The roadmap and readiness register agree on which controls must exist before M1.
- Masterplan acceptance review resolves its stage ordering against this ADR.

## Supersession

None.
