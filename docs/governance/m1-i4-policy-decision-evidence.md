# M1 iteration 4 policy-decision evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-22

## Bounded outcome

This iteration implements the next authorized M1 processing step after intent capture: deterministic evaluation of trusted, versioned policy against an authenticated principal, Tenant-minted scope, immutable captured delivery run, requested action, purpose, classification, risk, and Developer Trial environment.

The outcome is a contract-valid `PolicyDecision` with permit, deny, or indeterminate status, mandatory obligations, reason codes, bounded validity, and a digest-bound identity. The iteration stops before Execution Envelope issuance, context assembly, policy enforcement, agent delegation, planning, worker scheduling, tool or model invocation, approval, or source-control activity.

## Implemented scope

The Policy implementation provides:

- a versioned `EvaluatePolicy` request contract for the single bounded `forge.context.assemble` action in the Developer Trial environment;
- receiving-capability validation of the active normalized principal, Tenant-minted scope, and immutable `Captured` delivery run;
- exact subject, resource, delivery-run, purpose, classification, tenant, workspace, project, risk, action, and environment matching;
- deterministic rule ordering and conflict handling, with matching deny rules taking precedence over permits;
- fail-closed indeterminate outcomes when no policy applies or mandatory obligation values conflict;
- fail-closed denial when a permit contains an obligation the receiving runtime does not support;
- canonical obligation ordering, stable reason codes, and immutable defensive policy copies;
- decision validity bounded by both configured policy TTL and principal expiry;
- digest-bound decision identities covering the evaluated authority and context facts, matching rules, outcome, obligations, policy version, and lifecycle; and
- normalized disclosure-safe invalid, forbidden, indeterminate, and dependency-unavailable failures.

Policy definitions are trusted constructor input in this iteration. The implementation does not claim a policy authoring, storage, distribution, or administration capability.

## Affected capability

- `forge.capability.policy`: owns policy evaluation, decision outcomes, obligations, reason codes, policy-version binding, and decision lifecycle.

Identity continues to own the principal. Tenant continues to own scope and repository binding. Intent and Delivery continue to own the captured objective and delivery run. The trust-contract module defines the shared `PolicyDecision` shape. Policy does not mint execution authority, own delivery state, perform enforcement, or execute work.

## Governing architecture

The implementation follows the engineering constitution; ADR-0001, ADR-0003, ADR-0007, ADR-0009, and ADR-0010; the accepted capability model; Policy and Authorization architecture; Security architecture; Identity and Tenancy architecture; first-slice capability contracts and data flow; and the Developer Trial security and privacy decisions.

Those documents require deterministic policy evaluation, explicit deny precedence, mandatory obligations, safe reason codes, version binding, fail-closed unknowns, and separation between a Policy Decision and an Execution Envelope. The registered module depends only on Delivery, Identity, Tenant, failure contracts, and trust contracts in the control plane.

## Verification evidence

Deterministic tests cover:

- a permit decision bound to the exact subject, action, resource, delivery run, purpose, classification, scope, risk, and environment;
- stable decision identity and canonical obligation ordering across repeated evaluation;
- validity narrowed to principal expiry;
- defensive copying of policy definitions and obligation parameters;
- mandatory and ordinary deny precedence over permits;
- fail-closed unknown mandatory obligations;
- no-applicable-policy and conflicting-obligation indeterminate outcomes;
- no-match outcomes for every variable policy match dimension, plus rejection of unsupported action and environment values before matching;
- invalid request contract, context, action, environment, risk, principal, scope, and delivery-run facts;
- dependency, clock, and cancellation failure normalization; and
- malformed policy-set, rule, match-list, and obligation rejection.

Repository verification additionally runs architecture fitness, Go formatting, vetting, all tests, the race detector, coverage, and the build.

## Explicit exclusions

This iteration does not implement a policy language, external policy engine, durable policy store, policy administration API, policy distribution, decision persistence, revocation channel, distributed enforcement point, Execution Envelope issuer, Context Package assembly, knowledge retrieval, agent delegation, planning, worker isolation, credentials, tools, models, verification, approvals, evidence storage, source-control publication, journey execution, or the public M1 CLI.

A permit decision is not execution authority. No caller may treat it as an Execution Envelope or bypass a future enforcement point.

## Authorized delivery evidence

Delivery used grant `forge.autonomy-grant.m1-bounded-delivery-v1`, protected revision `424fdc7424736f3c018114b1b318b9c84125dfff`, digest `sha256:1e9b2e8030d5cb7229a92909ac8099cae209c70ea7b663d7a659a4533acff0fa`, task `m1-iteration-4-policy-decision`, and branch `codex/m1-iteration-4-policy-decision`.

The human-approved amendment added only `mkdocs.yml` to register this evidence document in canonical navigation. Merge to the protected branch remains a mandatory human gate.

## Architecture decision assessment

No new ADR is proposed. The implementation realizes the already accepted `EvaluatePolicy` responsibility inside the modular control plane and preserves the established capability, deployment, and trust boundaries. An ADR SHALL precede a new policy authority, external policy engine, policy language, durable policy topology, changed public contract, separate deployment unit, or any design that combines a Policy Decision with execution authority.
