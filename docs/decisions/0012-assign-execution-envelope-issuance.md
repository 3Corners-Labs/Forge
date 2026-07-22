# ADR-0012: Assign Execution Envelope issuance to the Policy trust boundary

Status: Accepted
Owner: Forge Trust and Delivery Architecture
Date: 2026-07-22

## Context

M1 has implemented human identity admission, Tenant-minted scope, immutable intent capture, and deterministic Policy Decisions. The next trust-path step requires an immutable Execution Envelope before Context Assembly may discover or retrieve repository content.

The Accepted architecture agrees on the separation and sequence but does not assign issuance with sufficient precision:

- ADR-0003 defines the Execution Envelope as the immutable authorization artefact and keeps it separate from Policy Decisions and Context Packages;
- the vertical slice says Policy grants a task-scoped Execution Envelope;
- the first-slice data inventory names a Delivery and trust issuance boundary;
- Runtime says the control plane issues or references the envelope; and
- the capability contract says Delivery does not own authorization and Policy does not own execution.

Leaving this unresolved would allow Delivery to become an implicit policy authority, invite a new undocumented capability, or let Context Assembly accept a caller-constructed envelope. Any of those outcomes would violate the accepted capability or trust boundaries.

The current executable envelope contract also requires both an initiating subject and an agent identity with an explicit delegation chain. Forge has not yet implemented task-agent identity or delegation issuance. Filling those fields with caller input or a synthetic constant would bypass Identity ownership.

## Decision

Assign deterministic root and child Execution Envelope issuance to the Policy trust boundary within the control plane.

Policy owns the authority-minting decision and SHALL:

- consume only active Identity-minted principal and delegation facts, Tenant-minted scope, immutable Delivery-owned run facts, an applicable permit Policy Decision, and trusted Developer Trial profile constraints;
- verify the exact decision point, subject, action, resource, policy version, validity window, and issuance-time obligations before granting authority;
- derive tenant, workspace, project, delivery-run, principal, agent, delegation, purpose, risk, actions, classifications, destinations, model IDs, tool IDs, secret-reference IDs, limits, validity, idempotency, and evidence-correlation fields server-side;
- refuse deny, indeterminate, expired, mismatched, unsupported-obligation, unverified-delegation, or authority-expanding input;
- emit the existing canonical Execution Envelope contract as an immutable, digest-bound authority object containing references rather than secret values or retrieved content; and
- require every child envelope to preserve parent lineage and be strictly narrower than its active parent.

Delivery owns orchestration and durable delivery state. It MAY request issuance, store or reference the resulting envelope ID, and transition a run only after successful issuance. Delivery SHALL NOT construct envelope authority fields, reinterpret a Policy Decision, weaken obligations, mint identity or delegation, or treat envelope possession as permission for an unrelated action.

Identity remains the canonical owner of human, workload, and task-agent identity and delegation. Execution Envelope implementation SHALL wait until Identity can supply a verified, expiring task-agent identity and delegation chain for the requested purpose, actions, resource, classifications, destinations, and autonomy limit. Policy and Delivery SHALL NOT substitute an implementation constant or caller assertion for that prerequisite.

Context Assembly, workers, tools, verification, approval, evidence, and source-control adapters are enforcement points. They validate the envelope and applicable Policy Decision before their protected action; they cannot issue, widen, or reinterpret authority.

The initial modular-monolith implementation MAY place issuance in the existing `internal/policy` module. This placement does not create a new capability or service and does not make Policy responsible for execution, durable delivery state, approval, or enforcement performed by another capability.

If accepted, this decision clarifies ownership only. It does not independently authorize functional work, child delegation, or a wider M1 outcome. Implementation remains bounded by ADR-0010, the applicable M1 iteration authorization, protected review, and the Developer Trial envelope.

## Alternatives considered

- **Let Delivery issue the envelope:** rejected because Delivery would become an authorization authority despite its explicit non-responsibility, and a future orchestration change could silently widen authority.
- **Create a separate Authority capability or service:** rejected because the accepted capability catalogue defines no such capability and M1 has no evidence requiring another deployment unit.
- **Return a Policy Decision and Execution Envelope as one combined object:** rejected because ADR-0003 requires separate authority objects with distinct lifecycles and references.
- **Let Context Assembly or the worker construct the envelope:** rejected because an enforcement point that processes untrusted content or code cannot mint or expand its own authority.
- **Use a fixed Developer Trial agent identity without delegation:** rejected because environment convenience does not override Identity ownership, attribution, expiry, or delegation requirements.

## Consequences

- The trust boundary that converts a permit decision into executable authority becomes explicit and testable.
- Delivery can coordinate issuance without importing or owning Policy semantics, avoiding a circular module dependency.
- Policy gains authority issuance responsibility but does not gain execution, workflow-state, approval, or provider responsibilities.
- Task-agent identity and narrowing delegation become the next implementation prerequisite before envelope issuance.
- Envelope fields and limits remain server-derived; public interfaces cannot submit authoritative scope or delegation values.
- Context Assembly can depend on one canonical authority object without copying authority into the Context Package.
- Revocation, expiry, policy changes, or changed authority inputs require denial or re-evaluation according to the applicable gate; an existing envelope is never mutated in place.
- No production support, external policy engine, signing-key topology, provider binding, or service extraction is selected by this decision.

For SaaS and self-hosted profiles, the logical owner remains the same. Profile-specific signing, key custody, revocation distribution, persistence, and recovery require their own accepted conformance decisions before a production-shaped claim.

## Validation

Implementation SHALL provide deterministic positive and negative tests proving:

- only an active exact permit decision can authorize issuance;
- Identity-minted task-agent identity and delegation are required and cannot be replaced by caller values;
- subject, scope, delivery run, resource, purpose, action, classification, risk, environment, policy version, and evidence correlation are bound exactly;
- validity and all limits are server-derived, bounded by upstream authority, and cannot be widened;
- unknown or unenforced issuance-time obligations fail closed;
- deny, indeterminate, expired, revoked, mismatched, replayed, and mutated inputs produce disclosure-safe failures;
- envelope identity and integrity change when an authority-relevant input changes;
- no envelope contains secrets or retrieved content;
- child envelopes preserve lineage and are strictly narrower; and
- Delivery, Context Assembly, workers, tools, and adapters cannot mint or widen envelopes.

Architecture fitness SHALL preserve the control-plane module boundary and prevent Delivery from importing Policy or constructing trust authority. M1 journey evidence SHALL demonstrate authorization before context discovery and retrieval through the same public path.

Sijbren Beukenkamp accepted this decision as drafted through the authenticated collaboration boundary on 2026-07-22, acting as the temporary accountable Architecture, Security, Identity, Delivery, and Engineering reviewer established for bounded M1. The decision becomes operationally normative when it reaches protected `main` through the required review path.

Execution Envelope implementation remains blocked until Identity provides the required task-agent identity and delegation prerequisite and a separate bounded M1 iteration is activated.

## Supersession

This decision clarifies ADR-0003 and ADR-0010. It does not supersede them or change the accepted capability catalogue, public M1 boundary, deployment profile, or authorization outcome.
