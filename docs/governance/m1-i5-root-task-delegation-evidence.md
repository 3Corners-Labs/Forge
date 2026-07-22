# M1 iteration 5 root-task-delegation evidence

Status: Proposed
Owner: Forge Engineering and Architecture
Version: 0.1
Verified: 2026-07-22

## Bounded outcome

This iteration implements the accepted ADR-0014 root authority bootstrap. Identity canonicalizes one bounded root `DelegationRequest`; Policy evaluates `forge.delegation.create` at the execution decision point against its exact digest and the verified human, task agent, Tenant scope, captured Delivery run, classification, risk, and Developer Trial environment; Identity then validates the same active permit before minting authority.

The result is one immutable, expiring, revocable, Identity-owned root delegation for `forge.context.assemble`. It is not an Execution Envelope and grants no retrieval, execution, tool, model, approval, publication, merge, release, deployment, policy-mutation, or production authority.

## Implemented scope

Identity provides:

- a versioned immutable request with canonical collections and a deterministic digest over every authority-relevant field while excluding correlation and idempotency metadata;
- exact M1 bounds: distinct verified synthetic human and agent subjects, one Developer Trial scope and captured run, purpose `forge.delivery.m1-change`, action `forge.context.assemble`, one repository and its classification, no destination or parent delegation, and `context-preparation-only` autonomy;
- exact active parent-Policy-Decision and enforced-obligation verification before issuance;
- expiry narrowed by request, both principals, parent decision, and configured Developer Trial maximum lifetime;
- atomic tenant-scoped idempotent creation, conflict on changed authority under the same key, and prevention of replay after revocation;
- Identity-owned revocation and disclosure-safe resolution of the one-record root chain; and
- defensive immutable projections and normalized failure contracts.

Policy provides a separate versioned root-delegation evaluation that validates the human, agent, Tenant scope, Delivery run, and canonical request. Matching includes action, request digest, run, purpose, classification, risk, environment, subject, tenant, workspace, and project. It returns the existing canonical `PolicyDecision`; it never creates or stores delegation authority.

## Affected capabilities

- `forge.capability.identity`: owns the request contract, canonical digest, delegation issuance, lifecycle, idempotency, revocation, and resolution.
- `forge.capability.policy`: owns the exact `forge.delegation.create` decision, obligations, outcome, reasons, and policy lifecycle.

Tenant and Delivery facts remain immutable inputs owned by their existing capabilities. Identity imports only the shared trust contract to verify its parent Policy Decision and does not import the Policy runtime. Policy consumes Identity facts but cannot mint a delegation. The module graph therefore remains acyclic.

## Governing architecture

The implementation follows the Constitution; ADR-0003, ADR-0010, ADR-0012, ADR-0013, and ADR-0014; Identity and Tenancy; Policy and Authorization; Security; Runtime and Workflow; the first-slice capability contracts and data flow; and the bounded Developer Trial implementation baseline.

ADR-0014 explicitly permits Identity to depend on the shared trust contract for parent-decision validation. The registered module boundary and fitness check permit that exact dependency while continuing to reject an Identity-to-Policy-runtime import.

## Verification evidence

Deterministic tests cover:

- canonical digest stability independent of non-authority replay metadata and digest changes for every M1 field that can vary within the bounded contract;
- rejection of changed purpose, action, destination, environment, autonomy, child authority, duplicates, and other broadened requests;
- exact decision point, human, action, request-digest, outcome, lifecycle, and obligation enforcement;
- independent active synthetic human and task-agent verification and subject-kind binding;
- exact Tenant, workspace, project, repository, run, purpose, classification, risk, and environment binding;
- permit, deny, indeterminate, no-match, mismatched, expired, unverified, cross-binding, unknown-adapter, and cancellation failures;
- expiry narrowing, defensive copies, atomic concurrent idempotency, changed-authority conflict, revocation, post-revocation replay denial, and one-record chain resolution; and
- architecture fitness for declared dependencies and the existing control-plane/worker isolation rules.

Repository verification runs formatting, architecture checks, all Go tests, vetting, the race detector, coverage, build, and focused delegation benchmarks. Measured benchmark results are recorded on the candidate revision before review; they are observations for this in-process Trial implementation, not a production capacity claim.

Focused in-process benchmarks on an Apple M1 Pro (`darwin/arm64`, eight benchmark workers) measured:

| Operation | Time | Allocated bytes | Allocations |
|---|---:|---:|---:|
| Request canonicalization and digest | 1.401 µs/op | 1,217 B/op | 9 allocs/op |
| Policy delegation evaluation | 7.543 µs/op | 5,269 B/op | 70 allocs/op |
| Identity permit verification and creation | 6.540 µs/op | 4,444 B/op | 40 allocs/op |
| Idempotent creation replay | 4.417 µs/op | 3,027 B/op | 26 allocs/op |
| Active root-chain resolution | 0.421 µs/op | 512 B/op | 7 allocs/op |
| Idempotent revocation lookup/update | 0.369 µs/op | 96 B/op | 6 allocs/op |
| Revoked-authority lookup and denial | 0.210 µs/op | 144 B/op | 4 allocs/op |

These microsecond-scale observations indicate that the required in-process authority check is not currently a perceptible UX bottleneck in isolation. They are not complete ADR-0013 envelope validation: they exclude tail-latency distributions, load saturation, network, external provider, durable-store, and failure-recovery latency and therefore do not authorize extrapolation to a distributed or production topology.

## Explicit exclusions

This iteration does not implement child or renewed delegation, parent-delegation narrowing, provider credentials, workload credentials, durable or distributed revocation, offline authority, policy administration, durable policy or delegation storage, an Execution Envelope issuer, Context Package assembly, retrieval, workers, models, tools, approvals, publication, merge, release, deployment, production profiles, or CLI wiring.

## Authorized delivery evidence

Delivery used grant `forge.autonomy-grant.m1-bounded-delivery-v1`, protected revision `6c6ff370736f7628ac552e8dbfb6b8ff6976bcc9`, digest `sha256:1e9b2e8030d5cb7229a92909ac8099cae209c70ea7b663d7a659a4533acff0fa`, task `m1-iteration-5-root-task-delegation`, and branch `codex/m1-iteration-5-root-task-delegation`.

The human-approved amendment added only `mkdocs.yml` to register this evidence document in canonical navigation. Merge to the protected branch remains a mandatory human gate.

## Architecture decision assessment

No new ADR is proposed. This iteration implements ADR-0014 and the issuance prerequisite established by ADR-0012 without changing capability ownership, trust-object schemas, deployment topology, or public interfaces. A new ADR SHALL precede child delegation, another parent-authority model, expanded autonomy, durable or distributed lifecycle topology, a new shared trust object, or movement of delegation ownership outside Identity.
