# Forge implementation baseline

Status: Accepted
Owner: Forge Engineering and Architecture
Version: 1.0
Verified: 2026-07-23

## Purpose

This document separates verified repository state from target architecture. It is the canonical current-state orientation until an implemented capability catalogue replaces it.

The presence of a document, schema, development dependency, or candidate technology does not prove that the corresponding Forge capability is implemented or supported.

## Verified current state

Forge currently provides a small developer bootstrap:

| Area | Implemented behaviour | Evidence | Maturity |
|---|---|---|---|
| CLI entry point | `init`, `doctor`, `up`, `down`, `version` commands; `up` and `down` validate the Developer profile before invoking Docker | `cmd/forge`, `internal/cli` and tests | Trial |
| Configuration | Strict parsing plus fail-closed Developer bootstrap validation for runtime mode, profile, and disabled intelligence | `internal/config` and tests | Trial |
| Prerequisite diagnostics | Detects Git, Go, Docker, and Task executables | `internal/doctor` | Trial |
| Local platform control | Invokes the developer Docker Compose profile | `internal/platform` | Trial |
| Developer environment | Loopback-bound PostgreSQL, Keycloak, OpenBao, MinIO, OpenTelemetry Collector, Prometheus, Grafana, and Mailpit | `deploy/compose` | Trial |
| Trust contract candidates | Closed JSON Schemas for Execution Envelope, Context Package, Policy Decision, Approval Record v1, Approval Request v1, Approval Record v2, and Evidence Record | `schemas/forge/trust/v1` | Research |
| Executable trust contracts | Strict Go decoding, structural validation, exact authority binding, child-envelope narrowing, canonical digest and integrity binding, Approval v1/v2 version separation, expected-fact approval validation, evidence-subject binding, and normalized disclosure-safe failures | `internal/contracts/trust`, `internal/contracts/failure`, and tests | Trial |
| Developer Trial identity admission | Provider-neutral assertion verification, audience and lifecycle validation, immutable normalized principal, and disclosure-safe provider failure normalization; only deterministic test providers exist | `internal/identity` and tests | Trial |
| Developer Trial scope admission | Single-test-tenant workspace membership resolution, server-derived immutable project scope, delayed repository binding discovery, and disclosure-safe cross-scope denial; storage and binding adapters remain test doubles | `internal/tenant` and tests | Trial |
| Developer Trial intent capture | Versioned direct-human intent capture, immutable `Captured` delivery run, Tenant-derived repository/classification binding, scope-bound atomic idempotency, and a concurrency-safe in-memory store behind a replaceable port | `internal/delivery` and tests | Trial |
| Developer Trial policy decision | Deterministic fail-closed evaluation of trusted versioned rules against exact Identity, Tenant, Delivery, action, purpose, classification, risk, and environment facts; deny precedence, mandatory obligations, bounded expiry, reason codes, and digest-bound Policy Decisions | `internal/policy`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial plan-review eligibility | Policy-owned `forge.approval.decide` evaluation bound to the exact Approval Request digest, active human, Delivery plan, Tenant scope, purpose, classification, risk, Developer Trial environment, obligations, policy versions, and expiry | `internal/policy` and tests | Trial |
| Developer Trial root task delegation | Immutable canonical root-delegation request, digest-bound `forge.delegation.create` Policy evaluation, Identity-owned bounded issuance, atomic tenant-scoped idempotency, expiry narrowing, revocation, and disclosure-safe one-record chain resolution for one synthetic task agent | `internal/identity`, `internal/policy`, and tests | Trial |
| Developer Trial root Execution Envelope | Policy-owned re-evaluation and decision ledger, live Identity delegation resolution, server-derived limits and expiry, canonical integrity binding, atomic idempotency, and one root envelope limited to context assembly | `internal/policy`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial context manifest assembly and reproduction | Authority and obligation enforcement before source-metadata discovery; exact tenant/project/repository validation; bounded deterministic Context Package manifests; explicit omissions; integrity and atomic replay; plus capability-owned manifest reproduction that reports changed, missing, or stale source metadata without returning content | `internal/context`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial context-informed planning | Explicit in-process transitions from `Captured` through `ContextAuthorizing`, `ContextReady`, `Planned`, `AwaitingPlanApproval`, and `PlanApproved`; optimistic concurrency, scope-bound idempotency, atomic minimized outbox entries, bounded deterministic acyclic plans, canonical plan digests, and exact invalidation checks | `internal/delivery` and tests | Trial |
| Developer Trial plan approval | Approval-owned immutable request and scope-bound Approval Record v2 lifecycle with canonical digests, exact human and Policy eligibility, evidence and condition binding, atomic idempotency, immutable revocation chains, fork prevention, and current-record resolution | `internal/approval`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial worker authority | Separate worker-phase delegation, workload principal binding, current Approval/Context/plan revalidation, Policy-owned bounded worker envelope issuance, expiry, idempotency, revocation, and PostgreSQL adapters | `internal/identity`, `internal/policy`, and tests | Trial |
| Developer Trial candidate preparation | Closed digest-bound worker contracts, deterministic file replacement, path/source/symlink/output controls, isolated Docker invocation, candidate bundle and tree digests, and Delivery transitions through `Executing` and `Verifying` | `internal/contracts/worker`, `internal/execution`, `internal/worker`, `cmd/forge-worker`, `internal/delivery`, and tests | Trial |
| Developer Trial independent verification | Closed Assurance Plan and result contracts, mandatory independent consequential oracle, separate verifier deployment unit, fresh-checkout integrity/path/Go-format controls, durable result adapter, and Delivery transition to `AwaitingPublicationApproval` | `internal/contracts/verification`, `internal/verification`, `internal/verifier`, `cmd/forge-verifier`, and tests | Trial |
| Developer Trial durable control state | Capability-owned PostgreSQL adapters for Tenant, Identity worker delegation, Delivery lifecycle and outbox, Policy worker envelopes, Context manifests, Approval records, Verification records, and Evidence; adapters remain unwired from the public CLI | capability `postgres_store.go` files and integration tests | Trial |
| Developer Trial evidence and artefacts | Append-only digest-linked Evidence streams with fork detection; bounded MinIO temporary upload, byte verification, immutable digest promotion, and PostgreSQL metadata | `internal/evidence`, dedicated durable-verification workflow, and tests | Trial |
| Assurance contract candidates | Closed JSON Schemas for journey scenarios and results | `schemas/forge/assurance/v1` | Research |
| M1 journey assurance candidates | Reviewed allowed and denied public-CLI scenarios, synthetic fixtures, and an execution plan that remains disabled until M1 | `assurance/journeys/m1` | Research |
| Offensive assurance safety candidate | Closed rules-of-engagement schema, disabled contract-only M1 configuration, digest-bound target/configuration, and negative safety fixtures | `schemas/forge/security/v1`, `security/assurance` | Research |
| Architecture module boundaries | Current module graph, planned worker boundary, negative dependency fixtures, and CI-control validation | `architecture/module-boundaries.json`, `architecture/testdata`, `scripts/check_architecture.py` | Trial |
| M1 threat and privacy candidates | Complete priority-threat mappings, phase-assigned negative contracts, Developer Trial decisions, and synthetic leakage fixtures | `security/threat-controls`, `security/privacy`, `security/fixtures` | Research |
| Deployment profile candidates | Closed profile schema plus Developer Trial and connected self-hosted reference manifests with phased conformance plans | `schemas/forge/deployment/v1`, `profiles` | Research |
| Architecture fitness checks | Validates document metadata, F0 governance, capability identifiers, schema sets, source integrity, navigation, local links, and developer-only Compose risks and constraints | `scripts/check_architecture.py` | Trial |
| Contributor delivery delegation | Closed, integrity-bound autonomy-grant schema; protected-revision validator; candidate non-interference, lifecycle, scope, budget, and mandatory-gate tests; inactive-by-default reviewed M1 grant | `.forge/autonomy`, `schemas/forge/contributor/v1`, `scripts/validate_autonomy_grant.py`, and tests | Trial |

The default configuration sets `intelligence.enabled` to `false`. The implemented M1 capabilities are currently callable Go modules and replaceable Developer Trial adapters only; they are not wired into the CLI and provide no live identity or repository provider. The current CLI does not authenticate users, capture intent through a public command, evaluate policy, issue execution envelopes through a public path, retrieve repository content, invoke models or tools, coordinate agents, expose plan review, publish pull requests, or operate the durable M1 path. PostgreSQL, MinIO, worker, verifier, and evidence modules therefore prove bounded internal behaviour, not a usable public Forge workflow.

Contributor delivery delegation is repository governance, not Forge product
runtime functionality. It becomes active only after a separate authenticated
human activation bound to an eligible protected grant revision and exact task.

## Target capabilities not yet implemented

All capabilities in `architecture/capabilities.md` remain Research unless a later evidence record says otherwise. In particular, the repository does not yet implement:

- durable tenant, workspace, project, membership, repository-binding, or quota administration beyond the bounded Developer Trial facts and PostgreSQL resolver;
- identity federation, credential acquisition, general child delegation, distributed workload identity, or distributed revocation beyond the bounded M1 worker lifecycle;
- durable policy sourcing, policy administration, durable decision or envelope persistence, revocation distribution, or distributed authorization enforcement beyond the bounded in-process Developer Trial evaluator and root issuer;
- scheduling, cancellation orchestration, automatic outbox dispatch, or automated recovery beyond durable checkpoint and pending-outbox resolution;
- repository-content retrieval, ranking, indexing, embedding, or memory lifecycle beyond bounded source-metadata manifest assembly and reproduction;
- agent registry, workforce coordination, model routing, or tool brokerage;
- production-grade sandbox assurance, dynamic workers, models, tools, secrets, or network authority;
- publication approval or candidate-digest approval;
- evidence signing, production non-repudiation, retention execution, or supported backup and restore;
- Workbench, Software System Twin, public product API, or journey runner;
- release, production deployment, SaaS tenancy, or supported self-hosted operation.

## Deployment posture

The Compose environment is a Developer Trial asset only. Its static development credentials, development-mode identity and secret services, version-tagged images, debug telemetry exporter, and local volumes are prohibited as production evidence.

No production profile is Supported. No production SLA, compatibility guarantee, recovery objective, or capacity envelope has been accepted.

## Verification baseline

As of the verified date:

- repository architecture checks pass;
- Go tests pass;
- configuration and CLI developer-safety paths have substantive automated test coverage;
- the worktree baseline is the revision recorded by Git, not this document alone.

Verification results are observations, not permanent claims. They SHALL be regenerated for each candidate revision.

## Claim rules

Contributors SHALL:

- label target behaviour as Proposed, Research, or Trial as applicable;
- link implementation claims to source and deterministic verification;
- update this baseline in the same change that materially alters implemented scope;
- avoid describing Compose dependencies as implemented Forge capabilities;
- avoid using `enterprise-ready`, `production-ready`, `secure`, `scalable`, `portable`, or `supported` without an accepted measurable envelope and evidence.

## Update triggers

Update this baseline when a capability becomes executable, a profile changes maturity, a public contract is introduced, a trust boundary changes, or verification evidence materially changes.
