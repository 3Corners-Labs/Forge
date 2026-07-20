# Identity, delegation, and tenancy

Status: Proposed
Owner: Forge Identity and Security Architecture
Version: 0.1

## Purpose

Define how Forge knows who is acting, for which tenant and workspace, through which delegation, and within which isolation boundary. Authentication proves an identity assertion; it does not by itself grant project access or execution authority.

## Responsibilities

`forge.capability.identity` owns:

- normalized human, workload and agent principals;
- verification of identity-provider assertions;
- delegation chains, expiry and revocation;
- workload identity issuance references;
- authentication assurance metadata.

`forge.capability.tenant` owns:

- tenants, workspaces, projects and membership;
- immutable scope identifiers;
- repository and environment bindings;
- tenant topology and quota assignments;
- lifecycle, export and deletion coordination.

Neither capability evaluates the final action policy. Policy consumes their verified facts.

## Canonical identities

Forge distinguishes:

- human principal;
- service or workload principal;
- agent identity;
- provider identity;
- extension publisher identity;
- emergency or break-glass identity.

An agent identity is never an alias for the initiating human. Evidence preserves both and the complete delegation chain.

## Tenant context

Tenant, workspace and project identifiers are derived server-side from authenticated membership and selected resources. Callers may request a scope but cannot authoritatively set or override it through prompts, context documents, tool output or provider metadata.

Scope travels through:

- API and event envelopes;
- workflow state and idempotency keys;
- caches, indexes and queues;
- object and relational storage;
- secret and key references;
- telemetry and Evidence Records;
- worker credentials and network policy.

Every storage and transport adapter SHALL demonstrate negative cross-tenant tests.

## Delegation

Delegation is explicit, time-bounded and narrowing. A delegation records issuer, subject, agent, permitted purposes, actions, resources, classifications, destinations, maximum autonomy, issue and expiry time, revocation source and parent delegation.

Delegation cannot:

- add authority absent from its parent;
- survive parent revocation or expiry;
- convert contextual instructions into permissions;
- authorize self-approval or policy modification implicitly;
- be reused for a different tenant, environment or purpose.

## Lifecycle and events

Principal lifecycle events include verified, assurance-changed, suspended and revoked. Tenant lifecycle events include created, membership-changed, export-requested, deletion-scheduled, suspended and closed.

Revocation propagates to active sessions, unissued envelopes and resolvable credentials. Already-running high-risk work is cancelled or re-evaluated according to policy.

## Initial quality and security scenarios

- A caller-provided tenant ID conflicting with verified membership is rejected before retrieval.
- An agent cannot exchange a task delegation for broader credentials.
- Revocation prevents new tool, secret and publication decisions within the declared propagation objective.
- Evidence identifies human initiator, agent, workload and provider separately.
- A project export contains no other tenant's identifiers, data, keys or telemetry.
- Operator access in a future SaaS profile is separately attributable, approved and visible to audit.

## Open decisions

- Identity federation protocols and reference provider.
- Tenant isolation topology for each deployment profile.
- Revocation propagation and offline-operation semantics.
- Break-glass ownership, duration and notification.

These require ADRs before a production-shaped profile becomes Supported.
