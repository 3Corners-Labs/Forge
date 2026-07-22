# M1 authority-to-context slice evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-22

## Bounded outcome

This combined delivery tranche implements two sequential M1 increments behind separate capability boundaries:

1. Policy issues one integrity-bound root Execution Envelope after re-evaluating the exact context action and resolving one active Identity-owned root delegation.
2. Context validates the envelope, Policy-owned decision, principals, scope, run, live delegation, obligations, and budgets before it calls a source-metadata adapter, then emits a minimal reproducible Context Package manifest.

The tranche stops before repository-content retrieval, context reproduction, ranking, indexing, embeddings, models, tools, planning, workers, candidate creation, verification, approval, evidence persistence, publication, merge, release, deployment, or production operation.

## Implemented Policy scope

Policy provides:

- an in-process Trial decision ledger containing defensive copies of decisions minted by the existing evaluator;
- disclosure-safe decision resolution so enforcement capabilities need not trust caller-constructed `PolicyDecision` values;
- a versioned root-envelope issuance request bound to verified human and agent principals, Tenant scope, immutable captured Delivery run, delegation reference, expected Policy Decision ID, and idempotency key;
- exact re-evaluation of `forge.context.assemble` at the retrieval gate under trusted Developer Trial environment, risk, policy-version, obligation, lifetime, and limit constraints;
- live disclosure-safe root-delegation resolution through the Identity boundary before every issuance or replay;
- server-derived scope, principals, delegation chain, purpose, actions, classifications, empty destinations/models/tools/secrets/approvals, policy references, expiry, and zero-external-cost limits;
- validity narrowed by both principals, delegation, Policy Decision, configured maximum lifetime, and the remaining envelope duration;
- deterministic envelope ID and SHA-256 integrity binding over every serialized field except the derived ID and integrity record; and
- tenant-scoped atomic idempotency, defensive replay, conflict on changed trusted authority, and normalized persistence failures.

The issuer does not accept a caller-supplied Policy Decision, risk, environment, actions, classifications, destinations, models, tools, secrets, cost, concurrency, or expiry.

## Implemented Context scope

Context provides:

- a versioned assembly request referencing capability-minted principals, scope, run, delegation, envelope, and Policy Decision ID;
- integrity and active-action validation of the exact Execution Envelope;
- live Identity delegation resolution and Policy-owned decision resolution without importing Policy runtime implementation;
- trusted obligation enforcement and exact subject, agent, tenant, workspace, project, delivery run, repository, purpose, risk, action, classification, and evidence-correlation binding;
- authority completion before any source adapter call, including denial paths that prove the adapter was not invoked;
- a provider-neutral metadata-discovery port that carries authority-derived scope and bounded requested source IDs but no credentials or content bytes;
- rejection of cross-scope, unrequested, duplicate, downgraded-classification, malformed, or unsupported-authority metadata;
- deterministic item ordering, explicit requested-source omissions, metadata-size and item-count budgets, and expiry bounded by the envelope and profile;
- deterministic Context Package ID and manifest digest with mutation verification; and
- tenant-scoped atomic idempotency that replays without rediscovery and conflicts before discovery when the request or trusted profile changes.

The resulting Context Package is information metadata only. It references exactly one envelope and cannot add authority.

## Affected capabilities and boundaries

- `forge.capability.policy` owns evaluation, decision lifecycle, and root Execution Envelope issuance under ADR-0012.
- `forge.capability.context` owns authorization enforcement at its boundary and the Context Package manifest.

Identity remains the sole delegation and revocation owner. Tenant remains the scope and repository-binding owner. Delivery remains the run and objective owner. The shared trust-contract module calculates and verifies canonical digests but grants no authority.

The new Context runtime module depends on Delivery, Identity, Tenant, failure contracts, and trust contracts. It cannot import Policy runtime implementation. Its end-to-end external test may use Policy through public contracts; architecture fitness now correctly evaluates runtime source files rather than treating test-only imports as deployed dependencies. Delivery remains unable to import Policy or construct or integrity-bind Execution Envelopes.

## Governing architecture

The implementation follows the Constitution; ADR-0003, ADR-0010, ADR-0012, ADR-0013, and ADR-0014; Identity and Tenancy; Policy and Authorization; Knowledge, Context, and Memory; Security; Runtime and Workflow; the vertical slice; and the first-slice capability contracts and data flow.

These sources require authority to remain separate from context, authorization before discovery or retrieval, server-derived scope and grants, active delegation, deterministic policy, fail-closed obligations, minimum relevant context, explicit provenance and classification, reproducible manifests, and no authority expansion through retrieved content.

## Verification evidence

Deterministic tests cover:

- Policy-owned decision lookup, defensive copies, and rejection of forged decision references;
- exact live delegation, principal, scope, run, repository, classification, action, policy-version, risk, environment, obligation, expiry, and evidence-correlation binding;
- permit issuance plus denial for deny, indeterminate, expired, revoked, malformed, unsupported-obligation, mismatched, and unverified inputs;
- envelope limit narrowing, empty external authority fields, integrity mutation detection, defensive replay, idempotency conflict, and normalized store failures;
- Context denial before discovery for mutated envelopes, forged decisions, unverified identity, revoked delegation, and unenforced obligations;
- cross-tenant, cross-workspace, cross-project, cross-repository, classification-downgrade, unrequested-source, unsupported-authority, and malformed-digest metadata rejection;
- deterministic minimal manifest creation, explicit omissions, metadata budget failure, manifest mutation detection, defensive replay without rediscovery, and conflict before rediscovery;
- an end-to-end synthetic path through identity, Tenant scope, intent capture, Policy decisions, root delegation, Execution Envelope, source-metadata discovery, and Context Package; and
- architecture fitness preventing runtime Context-to-Policy coupling and Delivery envelope construction.

Repository verification runs formatting, architecture fitness, all Go tests, vetting, race detection, coverage, build, and focused benchmarks.

## Performance observations

Focused in-process benchmarks on an Apple M1 Pro (`darwin/arm64`, eight benchmark workers) measured:

| Operation | Time | Allocated bytes | Allocations |
|---|---:|---:|---:|
| Root Execution Envelope creation | 50.438 µs/op | 16,892 B/op | 152 allocs/op |
| Root Execution Envelope replay | 53.685 µs/op | 16,073 B/op | 144 allocs/op |
| Authorized Context Package assembly with metadata discovery | 47.503 µs/op | 15,327 B/op | 130 allocs/op |
| Authorized Context Package replay without rediscovery | 31.187 µs/op | 9,055 B/op | 81 allocs/op |

These microbenchmarks indicate that the combined in-process authority-to-context path is not currently a perceptible UX bottleneck in isolation and is far below the M1 local policy-decision hypothesis. They are not complete ADR-0013 validation: they do not provide p50/p95/p99 distributions, saturation, external repository latency, durable-store behaviour, noisy-neighbour behaviour, or failure-recovery capacity and therefore support no production or distributed-topology claim.

## Explicit exclusions

This tranche does not implement child envelopes, child delegation, envelope or decision revocation distribution, durable policy/decision/envelope/context storage, signing keys, external policy engines, repository-content retrieval, context reproduction, ranking, indexing, embedding, memory, models, tools, credentials, workers, plans, candidates, verification, approvals, evidence storage, source-control publication, CLI wiring, production profiles, or support commitments.

## Authorized delivery evidence

Delivery used grant `forge.autonomy-grant.m1-bounded-delivery-v1`, protected revision `03a0d06ae506e39595ee7677fe7b713042d84d03`, digest `sha256:1e9b2e8030d5cb7229a92909ac8099cae209c70ea7b663d7a659a4533acff0fa`, task `m1-authority-to-context-slice`, and branch `codex/m1-authority-to-context-slice`.

The human-approved amendment added only `mkdocs.yml` to register this evidence document in canonical navigation. Merge to protected `main` remains a mandatory human gate.

## Architecture decision assessment

No new ADR is proposed. The tranche realizes ADR-0012 root issuance and ADR-0003 authority/context separation within the accepted modular-monolith topology. It registers the already accepted Context capability as a Trial module without changing ownership or adding a deployment unit.

A new ADR SHALL precede child-envelope issuance semantics beyond the accepted narrowing contract, durable or distributed decision/envelope/context topology, cryptographic signing or key custody, repository provider selection, content retrieval or indexing architecture, model or tool destinations, or movement of issuance outside Policy.
