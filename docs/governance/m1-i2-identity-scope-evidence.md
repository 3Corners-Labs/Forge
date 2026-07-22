# M1 iteration 2 identity and scope evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-22

## Bounded outcome

This iteration implements the first admission segment of the authorized M1 path: a provider assertion becomes an immutable normalized principal, and that principal may resolve one Developer Trial workspace, project, and repository binding only when capability-owned membership permits it.

The outcome stops before intent capture, policy evaluation, delivery-run creation, context retrieval, worker execution, approval, evidence persistence, or source-control activity.

## Implemented scope

The Identity module implements:

- a provider-neutral assertion-verification boundary;
- versioned verification requests with correlation identifiers;
- audience, required-field, principal-kind, issue-time, and expiry validation;
- immutable normalized principal facts with no retained credential;
- receiving-capability lifecycle revalidation; and
- normalized, disclosure-safe authentication and provider-unavailability failures.

The Tenant module implements:

- one logical synthetic tenant per Developer Trial directory;
- versioned scope and repository-binding queries with correlation identifiers;
- capability-owned workspace, project, and membership definitions;
- untrusted workspace, project, and repository selectors with no tenant override;
- membership resolution before project or repository binding return;
- immutable scope values that callers cannot construct from exported fields;
- repository binding discovery only after a valid scope has been minted;
- binding verification against the authoritative tenant, workspace, project, and repository selector; and
- indistinguishable disclosure-safe denial for unavailable and unauthorized scopes.

Both modules use replaceable interfaces. Only deterministic test providers and binding sources exist in this iteration.

## Affected capabilities

- `forge.capability.identity`: assertion normalization and active-principal validation only.
- `forge.capability.tenant`: bounded membership, scope, and repository-binding resolution only.

Policy does not evaluate these records in this iteration. Tenant does not authenticate a principal, and Identity does not own membership or final authorization.

## Governing architecture

The implementation follows the engineering constitution; ADR-0001, ADR-0003, ADR-0007, ADR-0009, and ADR-0010; the accepted capability model; Identity and Tenancy architecture; first-slice capability contracts and data flow; security architecture; and accepted Developer Trial threat and privacy decisions.

The public `forge delivery start` command remains unimplemented. ADR-0009 requires that command to capture intent and return a delivery-run reference through the real public boundary. Adding a command that only resolves identity or scope would misrepresent the implemented state.

## Verification evidence

Deterministic tests cover:

- valid synthetic human assertion normalization;
- missing credentials and invalid provider assertions;
- provider unavailability without diagnostic disclosure;
- unsupported request-version denial before provider or repository-adapter calls;
- wrong audience, future issue time, and expired assertion denial;
- principal expiry revalidation by the receiving capability;
- successful membership-derived scope and repository binding;
- absence of a caller-authoritative tenant selector;
- non-member and unknown-scope responses with the same safe code and message;
- denial before repository binding discovery;
- rejection of a scope not minted by the Tenant module;
- binding absence, adapter unavailability, and authoritative-scope mismatch; and
- expiry before repository binding resolution.

Repository verification additionally runs architecture fitness, Go formatting, vetting, all tests, and the build.

## Explicit exclusions

This iteration does not implement identity-provider selection, federation protocol, interactive login, reusable credentials, delegation, revocation propagation, durable membership storage, quotas, cross-adapter tenant isolation, repository credentials, repository content discovery, policy evaluation, intent capture, delivery state, the public M1 CLI, telemetry, or Evidence Records.

## Architecture decision assessment

No new ADR is proposed. The implementation realizes already accepted capability responsibilities inside the modular control plane and uses the explicitly permitted Developer Trial contract-double posture. An ADR SHALL precede selection of a live identity provider, federation protocol, durable tenancy topology, repository provider, credential model, or changed authority boundary.
