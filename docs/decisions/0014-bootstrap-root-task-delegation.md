# ADR-0014: Bootstrap root task delegation with digest-bound policy authority

Status: Accepted
Owner: Forge Identity and Security Architecture
Date: 2026-07-22

## Context

M1 has implemented normalized human identity, Tenant-minted scope, captured intent, and deterministic Policy Decisions. ADR-0012 requires Identity-minted task-agent identity and delegation before Policy may issue an Execution Envelope.

The Accepted Identity contract defines `CreateDelegation(parentAuthority, subject, purpose, actions, resources, expiry)` but does not define the root `parentAuthority`. This creates a bootstrap problem:

- Identity owns delegation but cannot decide final action policy;
- Policy may issue an Execution Envelope only after Identity supplies delegation;
- an Execution Envelope therefore cannot be the parent of the first task delegation;
- a human identity alone proves attribution but not task authority; and
- the current Policy Decision contract does not directly expose every delegation bound such as scope, classification, destination, environment, and maximum autonomy.

Using caller assertions, a fixed task agent, or an unbound Policy Decision would let authority-relevant fields escape the policy decision. Introducing a second grant object would add another trust lifecycle before evidence shows that it is necessary.

## Decision

Use an existing canonical permit `PolicyDecision` as the parent authority for one root task delegation, with its resource bound to the digest of an immutable Identity-owned `DelegationRequest`.

### Delegation request

Identity owns the versioned `DelegationRequest` contract and its canonical digest. The request is not authority. It is untrusted requested authority until Policy permits its exact digest and Identity mints a delegation.

The request records:

- contract version, idempotency key, and evidence correlation identifier;
- initiating human subject and verified task-agent subject;
- tenant, workspace, project, repository, and delivery-run identifiers;
- purpose, actions, resources, classifications, and destinations;
- environment and maximum-autonomy level;
- requested expiry; and
- optional parent delegation ID, which SHALL be absent for the bounded M1 root delegation.

The canonical digest binds the contract version and every authority-relevant field using a deterministic representation. The idempotency key and correlation metadata do not grant authority and remain outside the authority digest, but stay on the request and resulting record for replay control and evidence. Collection ordering is canonical and duplicate or empty values are invalid.

For the first M1 implementation, the request is limited to:

- one verified synthetic human initiator and one separately verified synthetic task-agent principal;
- the exact Tenant-minted Developer Trial project and Forge repository binding of one captured delivery run;
- purpose `forge.delivery.m1-change`;
- action `forge.context.assemble` only;
- the repository resource and its actual classification;
- no external destination, model, tool, secret reference, approval, publication, release, merge, deployment, policy mutation, or production action;
- environment `developer-trial`; and
- maximum autonomy `context-preparation-only`.

Extending these values is a later policy and iteration decision, not an implementation default.

### Root authorization sequence

The root sequence is:

1. Identity verifies the active human and task-agent principals independently and preserves their distinct subject IDs and kinds.
2. The control-plane caller submits the requested delegation bounds with Tenant-minted scope and immutable Delivery-owned run facts. Caller input remains non-authoritative.
3. Identity canonicalizes the requested bounds as a `DelegationRequest`; receiving capabilities can validate and recompute its digest without accessing credentials.
4. Policy evaluates action `forge.delegation.create` at decision point `execution` against the active human, exact request digest as resource, Tenant scope, delivery purpose, classification, risk, and Developer Trial environment.
5. Policy returns the existing canonical `PolicyDecision`. Deny or indeterminate creates no delegation.
6. Identity receives the same request and Policy Decision, recomputes the digest, and verifies an active exact permit for decision point `execution`, initiating human subject, action `forge.delegation.create`, and that digest.
7. Identity verifies that every issuance-time obligation is supported and reported as enforced before minting.
8. Identity mints an immutable root delegation whose parent-authority reference is the Policy Decision ID.

Policy does not mint or store the delegation. Identity does not reinterpret policy or infer authority from the human principal. Delivery may coordinate the calls and retain a delegation reference but cannot construct authoritative request facts, change the digest, or mint the record.

### Delegation lifecycle

The root delegation records the architecture-required issuer, human subject, agent, purpose, actions, resources, classifications, destinations, maximum autonomy, issue and expiry time, revocation source, parent-authority reference, scope, environment, and evidence correlation.

Its expiry SHALL be no later than the earliest of:

- requested expiry;
- human-principal expiry;
- task-agent-principal expiry;
- parent Policy Decision expiry; and
- the Developer Trial maximum delegation lifetime.

Creation is atomic and idempotent within immutable tenant scope. Reuse of an idempotency key with a different request digest fails safely. Identity owns revocation state and disclosure-safe chain resolution. Revocation or expiry prevents new envelope issuance and downstream protected actions at their next enforcement point.

The M1 record contains one root delegation only. Child delegation, parent-delegation narrowing, offline authority, delegation renewal, provider credentials, and distributed revocation propagation require later separately accepted design and evidence.

The delegation is not an Execution Envelope and does not itself permit context retrieval, worker execution, tool invocation, approval, or publication. Policy remains responsible for issuing an exact Execution Envelope under ADR-0012; downstream enforcement validates both the delegation reference and the envelope appropriate to its action.

### Module boundary

The initial modular-monolith implementation MAY extend `internal/identity` to depend on the shared trust contract solely to validate the parent `PolicyDecision`. Policy already consumes Identity facts. Identity SHALL NOT import the Policy runtime module, so no circular runtime dependency is introduced.

This decision clarifies root authority bootstrap only. It does not authorize functional implementation, an external identity or policy provider, a production profile, child delegation, broader autonomy, or a new trust-object schema.

## Alternatives considered

- **Use the initiating human principal as parent authority:** rejected because authentication and attribution do not grant task actions, scope, classification, destination, or autonomy.
- **Use the Execution Envelope as parent authority:** rejected because ADR-0012 requires delegation before envelope issuance, creating a circular dependency.
- **Let Identity evaluate task policy itself:** rejected because Identity owns delegation lifecycle but not final action authorization.
- **Let Delivery construct or mint delegation:** rejected because Delivery owns workflow state and explicitly does not own authorization or identity.
- **Introduce a new Delegation Grant trust object:** rejected because the existing Policy Decision can bind the canonical request digest without adding another cross-capability authority object and lifecycle.
- **Reuse a Policy Decision whose resource is only the repository:** rejected because classification, destination, environment, autonomy, agent, scope, run, and expiry mutations would not be exactly bound.
- **Use a fixed Developer Trial agent or delegation constant:** rejected because development convenience cannot replace verified identity, attribution, lifecycle, digest binding, or revocation.
- **Implement root and child delegation together:** rejected because M1 requires only the minimum root task path and child authority adds lineage, narrowing, revocation, and recovery risk.

## Consequences

- Root delegation has an unambiguous parent authority without requiring an Execution Envelope or a new trust-object type.
- Policy authorizes the exact requested bounds by digest while Identity remains the canonical delegation owner.
- The module graph remains acyclic: Policy can consume Identity contracts, while Identity consumes only the shared Policy Decision contract rather than Policy implementation.
- A second Policy evaluation adds bounded latency and evidence but protects the authority bootstrap. ADR-0013 requires it to meet the applicable M1 performance envelope rather than being silently removed for speed.
- The request and delegation contracts need canonicalization, idempotency, lifecycle, revocation, and disclosure-safe validation tests.
- One root delegation is sufficient to unlock bounded Execution Envelope issuance for context assembly; later actions require their own explicit policy and delegation scope.
- No identity federation, external credential issuance, provider selection, production revocation objective, or Supported maturity is implied.

## Validation

Implementation SHALL provide deterministic positive and negative tests proving:

- distinct active human and task-agent principals with the required kinds and synthetic Developer Trial profile;
- canonical request digest stability independent of input collection order;
- digest changes for every authority-relevant subject, agent, scope, run, purpose, action, resource, classification, destination, environment, autonomy, and expiry mutation;
- exact active permit validation for decision point, human subject, `forge.delegation.create`, request digest, and obligations;
- denial for deny, indeterminate, expired, mismatched, unsupported-obligation, unverified-agent, cross-scope, wrong-run, wrong-environment, and broadened input;
- expiry bounded by the request, both principals, parent decision, and Developer Trial limit;
- atomic idempotent replay and conflict on changed authority under the same key;
- revocation preventing resolution for new downstream authority without revealing sensitive reason or existence information;
- immutable defensive projections and no credential, secret, prompt, source content, or unrestricted policy input in the request or record;
- no child delegation, self-approval, publication, merge, release, deployment, policy mutation, or production authority; and
- architecture fitness preventing Identity from importing Policy runtime implementation or Delivery from minting delegation.

Performance evidence SHALL measure canonicalization, permit verification, creation, replay, revocation lookup, and chain resolution under the applicable M1 validation envelope. The implementation remains Trial and in-process until evidence justifies another topology.

Sijbren Beukenkamp accepted this decision as drafted through the authenticated collaboration boundary on 2026-07-22, acting as the temporary accountable Architecture, Security and Privacy, Identity, Policy, Delivery, and Engineering reviewer established for bounded M1. The decision becomes operationally normative when it reaches protected `main` through the required review path.

Delegation implementation remains blocked until this decision lands and a separate bounded M1 implementation iteration is activated.

## Supersession

This decision clarifies the root `parentAuthority` required by the Accepted Identity contract and ADR-0012. It does not supersede ADR-0003, ADR-0010, ADR-0012, or ADR-0013.
