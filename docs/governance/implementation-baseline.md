# Forge implementation baseline

Status: Accepted
Owner: Forge Engineering and Architecture
Version: 1.0
Verified: 2026-07-25

## Purpose

This document separates verified repository state from target architecture. It is the canonical current-state orientation until an implemented capability catalogue replaces it.

The presence of a document, schema, development dependency, or candidate technology does not prove that the corresponding Forge capability is implemented or supported.

## Framing change

Accepted [ADR-0024](../decisions/0024-rebaseline-forge-as-bounded-means.md) rebaselined Forge as a bounded means for delivering 3Corners applications rather than a product developed for independent sale. Three consequences affect how this baseline is read:

1. **The contributor harness is the recognized primary delivery mechanism.** Every material change recorded below was produced by the [ADR-0011](../decisions/0011-bounded-autonomous-contributor-delivery.md) contributor harness under human review, not by the Forge runtime. The runtime has produced no real repository change: its source-control adapter is not live-configured and no model provider is active by default. This remains true as of the verified date.
2. **Runtime capability is extended only on demonstrated delivery need.** Absence of a capability below is no longer evidence of a gap to be closed on schedule. Roadmap milestones M4 through M9 carry evidence label Vision.
3. **Sovereignty tiers now apply.** No claim of community maintainability or hyperscaler independence may be published before at least one complete Forge delivery path is measured against a locally hosted model.

The verified state below is unchanged by that decision. Nothing was implemented, withdrawn, or promoted by it.

## Verified current state

Forge currently provides a small developer bootstrap:

| Area | Implemented behaviour | Evidence | Maturity |
|---|---|---|---|
| CLI entry point | `init`, `doctor`, `up`, `down`, `version`, `delivery start`, `workbench serve`, and `workbench chat`; both Workbench surfaces use the same local API and retain exact proposal-promotion, plan-approval, and publication-approval actions | `cmd/forge`, `internal/cli` and tests | Trial |
| Forge Workbench | Loopback-only API-first Developer Trial with an embedded responsive web cockpit, conversational terminal client, live agent and dependency projection, decision inbox, activity stream, candidate/test/evidence inspection, OpenAPI description, and browser accessibility journey; state remains in-process and projects the existing guided M1 capability composition | `internal/experience/workbench`, `internal/cli`, `assurance/journeys/workbench`, and tests | Trial |
| Configuration | Strict parsing plus fail-closed Developer bootstrap validation for runtime mode, profile, and disabled intelligence | `internal/config` and tests | Trial |
| Prerequisite diagnostics | Detects Git, Go, Docker, and Task executables; operational mode also tests the Docker daemon and required durable-runtime configuration without displaying values | `internal/doctor` | Trial |
| Local platform control | Invokes the developer Docker Compose profile | `internal/platform` | Trial |
| Developer environment | Loopback-bound PostgreSQL, Keycloak, OpenBao, MinIO, OpenTelemetry Collector, Prometheus, Grafana, and Mailpit | `deploy/compose` | Trial |
| Trust contract candidates | Closed JSON Schemas for Execution Envelope, Context Package, Policy Decision, Approval Record v1, Approval Request v1, Approval Record v2, and Evidence Record | `schemas/forge/trust/v1` | Research |
| Executable trust contracts | Strict Go decoding, structural validation, exact authority binding, child-envelope narrowing, canonical digest and integrity binding, Approval v1/v2 version separation, expected-fact approval validation, evidence-subject binding, and normalized disclosure-safe failures | `internal/contracts/trust`, `internal/contracts/failure`, and tests | Trial |
| Developer Trial identity admission | Provider-neutral assertion verification, audience and lifecycle validation, immutable normalized principal, and disclosure-safe provider failure normalization; only deterministic test providers exist | `internal/identity` and tests | Trial |
| Developer Trial scope admission | Single-test-tenant workspace membership resolution, server-derived immutable project scope, delayed repository binding discovery, and disclosure-safe cross-scope denial; storage and binding adapters remain test doubles | `internal/tenant` and tests | Trial |
| Developer Trial intent capture | Versioned direct-human intent capture, immutable `Captured` delivery run, Tenant-derived repository/classification binding, scope-bound atomic idempotency, and a concurrency-safe in-memory store behind a replaceable port | `internal/delivery` and tests | Trial |
| Developer Trial policy decision | Deterministic fail-closed evaluation of trusted versioned rules against exact Identity, Tenant, Delivery, action, purpose, classification, risk, and environment facts; deny precedence, mandatory obligations, bounded expiry, reason codes, and digest-bound Policy Decisions | `internal/policy`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial plan-review eligibility | Policy-owned `forge.approval.decide` evaluation bound to the exact Approval Request digest, active human, Delivery plan, Tenant scope, purpose, classification, risk, Developer Trial environment, obligations, policy versions, and expiry | `internal/policy` and tests | Trial |
| Developer Trial root task delegation | Immutable canonical phase-specific context and model-proposal delegation requests, digest-bound `forge.delegation.create` Policy evaluation, Identity-owned bounded issuance, atomic tenant-scoped idempotency, expiry narrowing, revocation, and disclosure-safe one-record chain resolution for one synthetic task agent | `internal/identity`, `internal/policy`, and tests | Trial |
| Developer Trial root Execution Envelope | Policy-owned re-evaluation and decision ledger, live Identity delegation resolution, server-derived limits and expiry, canonical integrity binding, atomic idempotency, and one root envelope limited to context assembly | `internal/policy`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial context manifest assembly and reproduction | Authority and obligation enforcement before source-metadata discovery; exact tenant/project/repository validation; bounded deterministic Context Package manifests; explicit omissions; integrity and atomic replay; plus capability-owned manifest reproduction that reports changed, missing, or stale source metadata without returning content | `internal/context`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial context-informed planning | Explicit in-process transitions from `Captured` through `ContextAuthorizing`, `ContextReady`, `Planned`, `AwaitingPlanApproval`, and `PlanApproved`; optimistic concurrency, scope-bound idempotency, atomic minimized outbox entries, bounded deterministic acyclic plans, canonical plan digests, and exact invalidation checks | `internal/delivery` and tests | Trial |
| Developer Trial plan approval | Approval-owned immutable request and scope-bound Approval Record v2 lifecycle with canonical digests, exact human and Policy eligibility, evidence and condition binding, atomic idempotency, immutable revocation chains, fork prevention, and current-record resolution | `internal/approval`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial worker authority | Separate worker-phase delegation, workload principal binding, current Approval/Context/plan and admitted Change Proposal digest revalidation, Policy-owned bounded worker envelope issuance, mandatory `proposal-admitted` obligation, expiry, idempotency, revocation, and PostgreSQL adapters | `internal/identity`, `internal/policy`, and tests | Trial |
| Developer Trial candidate preparation | Closed digest-bound Change Proposal and worker contracts; deterministic conversion of admitted proposal bytes; exact purpose-limited source reconstruction; path/source/symlink/output controls; immutable non-root network-disabled Docker invocation; candidate bundle and tree digests; cleanup proof; and Delivery transitions through `Executing` and `Verifying` | `internal/contracts/change`, `internal/contracts/worker`, `internal/intelligence`, `internal/execution`, `internal/worker`, `cmd/forge-worker`, `internal/delivery`, and tests | Trial |
| Developer Trial independent verification | Closed Assurance Plan and result contracts, mandatory independent consequential oracle, separate immutable non-root verifier deployment unit, independently reconstructed read-only source, durable result adapter, cleanup proof, and Delivery transition to `AwaitingPublicationApproval`. Forge-owned controls cover candidate integrity, declared path scope, Go format, canonical document metadata, and Markdown link resolution; the verifier additionally refuses any plan whose independent consequential controls are all incapable of failing for the changed paths. Controls are compiled into the `FROM scratch` verifier image rather than executed as subprocesses | `internal/contracts/verification`, `internal/verification`, `internal/verifier`, `cmd/forge-verifier`, [verifier independence evidence](verifier-independence-evidence.md), and tests | Trial |
| Developer Trial Requirements Proposal | Digest-bound structured proposal that separates direct human statements from assumptions, hypotheses, open questions, and experiments; conversational review can produce immutable human-sourced field-correction revisions or reject without authority; superseded digests fail closed and exact authenticated-human promotion of the current revision is required before Delivery captures intent | `internal/experience`, `internal/intent`, `internal/cli`, and tests | Trial |
| Developer Trial model-assisted discovery | Closed provider-neutral discovery request and proposal contracts; Experience-owned conversational turns using exact validated activation limits; in-process checkpoint preservation and an explicit accessible retry decision after normalized bounded-retry failure; strict regular-file activation loading bound to a separately approved exact digest; Intelligence-owned eligibility, pricing, budgets, adapter invocation, activation-bound receipt and disclosure-safe failures; strict OpenAI Responses adapter with storage, redirects, and environment-proxy routing disabled; measured token, latency, and estimated-cost receipts; synthetic model-selection plan; and preserved contract-double journey. No provider is active by default | `internal/contracts/discovery`, `internal/intelligence`, `internal/experience`, `internal/cli`, `schemas/forge/intelligence/v1`, `assurance/evaluations/discovery`, and tests | Trial |
| Developer Trial model-assisted change preparation | Optional ADR-0023 path with a distinct exact provider activation; phase-specific task-agent delegation; Policy-issued `forge.intelligence.propose-change` authority; closed request and proposal schemas; purpose-minimized source and acceptance bindings; strict OpenAI structured output without tools; activation-bound usage receipt; deterministic contract-double fallback; proposal-digest-bound worker authority; and explicit assumptions, uncertainty, omissions, cost, candidate, and verification projection. No change provider is active by default | `internal/contracts/change`, `internal/intelligence`, `internal/identity`, `internal/policy`, `internal/execution`, `internal/cli`, `schemas/forge/intelligence/v1`, `assurance/evaluations/change-preparation`, and tests | Trial |
| Developer Trial publication approval | Separate Approval Record v2 action for the exact verified candidate, current evidence references, Tenant-derived scope, environment, purpose, classification, risk, Policy eligibility, lifecycle, and expiry | `internal/approval`, `internal/policy`, `internal/contracts/trust`, and tests | Trial |
| Developer Trial source-control publication | Provider-neutral exact-subject, tree, path, repository, base, metadata-profile, current Approval, Policy, and required Evidence validation; deterministic operation and branch identity; bounded ambiguity reconciliation; a credential-free idempotent contract double; and an ADR-0019 GitHub App adapter limited to Git data and draft pull-request endpoints with exact repository, credential reference, candidate files, marker reconciliation, and no merge/release/deploy/admin surface. The GitHub adapter is not live-configured | `internal/sourcecontrol` and tests | Trial |
| Guided M1 public journey | Developer Trial composition from untrusted requirements input through proposal promotion, plan approval, bounded candidate preparation, independent verification, publication approval, and one reviewable contract-double pull request; default contract mode stays in-process, while explicit Docker mode uses exact local source, isolated runtimes, and durable stores; denied cross-scope identity fails before repository discovery | `internal/cli`, `assurance/journeys/m1`, and tests | Trial |
| Developer Trial durable control state | The operational public journey wires capability-owned PostgreSQL adapters for Experience, Intent, Delivery lifecycle and outbox, Identity delegations, Policy envelopes, Context manifests, Approval records, Verification records, Evidence, and Source Control; automatic cross-process journey resume remains unimplemented | capability `postgres_store.go` files, `internal/cli/durable_runtime.go`, and integration tests | Trial |
| Developer Trial evidence and artefacts | Append-only digest-linked Evidence streams with fork detection; bounded MinIO temporary upload, byte verification, immutable digest promotion, verified reads, PostgreSQL metadata, and operational proposal/candidate artefact wiring | `internal/evidence`, `internal/cli`, dedicated durable-verification workflow, and tests | Trial |
| Assurance contract candidates | Closed JSON Schemas for journey scenarios and results | `schemas/forge/assurance/v1` | Research |
| M1 journey assurance candidates | Accepted allowed oracle v2.0.0 for the complete guided contract-double path, unchanged denied-before-discovery oracle, synthetic fixtures, and an executable M1 candidate plan | `assurance/journeys/m1` | Trial |
| Offensive assurance safety candidate | Closed rules-of-engagement schema, disabled contract-only M1 configuration, digest-bound target/configuration, and negative safety fixtures | `schemas/forge/security/v1`, `security/assurance` | Research |
| Architecture module boundaries | Current module graph, planned worker boundary, negative dependency fixtures, and CI-control validation | `architecture/module-boundaries.json`, `architecture/testdata`, `scripts/check_architecture.py` | Trial |
| M1 threat and privacy candidates | Complete priority-threat mappings, phase-assigned negative contracts, Developer Trial decisions, and synthetic leakage fixtures | `security/threat-controls`, `security/privacy`, `security/fixtures` | Research |
| Deployment profile candidates | Closed profile schema plus Developer Trial and connected self-hosted reference manifests with phased conformance plans | `schemas/forge/deployment/v1`, `profiles` | Research |
| Architecture fitness checks | Validates document metadata, F0 governance, capability identifiers, schema sets, source integrity, navigation, local links, and developer-only Compose risks and constraints | `scripts/check_architecture.py` | Trial |
| Contributor delivery delegation | Closed, integrity-bound autonomy-grant schema; protected-revision validator; candidate non-interference, lifecycle, scope, budget, and mandatory-gate tests; inactive-by-default reviewed M1 grant | `.forge/autonomy`, `schemas/forge/contributor/v1`, `scripts/validate_autonomy_grant.py`, and tests | Trial |

The default configuration sets `intelligence.enabled` to `false`. The public
`forge delivery start` command now accepts either the deterministic requirements
file or a natural-language message through the same bounded M1 capabilities.
The discovery and change-preparation paths remain unavailable until an
accountable human approves the applicable purpose-specific provider activation
and digest. When supplied to the public command, Forge can configure the
bounded OpenAI adapter without resolving its environment credential until an
eligible invocation. No activation or credential is committed. Contract tests
compose the same paths with credential-free model and HTTP transport doubles.
The delivery journey uses synthetic Developer Trial identity and separate
human gates. Contract mode uses in-process runners and memory stores.
Operational mode uses immutable Docker worker and verifier images,
capability-owned PostgreSQL stores, MinIO artefacts, and a durable
credential-free source-control contract double. Both remain bounded learning
workflows, not live repository workflows.

The command does not federate a real identity, perform general repository
retrieval, invoke a tool, coordinate an agent workforce, call GitHub, create a
real branch or pull request, merge, release, deploy, or operate production.
Operational mode admits only two declared local text sources and requires
explicit exact revision and immutable runtime bindings. Durable records survive
store reconstruction, but automatic cross-process conversational resume remains
unimplemented. The optional live discovery and change-preparation seams prove no
provider quality, latency, availability, billing, retention, residency, or
production fitness until separately approved synthetic trials record that
evidence.

Contributor delivery delegation is repository governance, not Forge product
runtime functionality. It becomes active only after a separate authenticated
human activation bound to an eligible protected grant revision and exact task.

## Target capabilities not yet implemented

All capabilities in `architecture/capabilities.md` remain Research unless a later evidence record says otherwise. In particular, the repository does not yet implement:

- durable tenant, workspace, project, membership, repository-binding, or quota administration beyond the bounded Developer Trial facts and PostgreSQL resolver;
- identity federation, credential acquisition, general child delegation, distributed workload identity, or distributed revocation beyond the bounded M1 worker lifecycle;
- durable policy sourcing, policy administration, durable decision or envelope persistence, revocation distribution, or distributed authorization enforcement beyond the bounded in-process Developer Trial evaluator and root issuer;
- scheduling, cancellation orchestration, automatic outbox dispatch, or automated recovery beyond durable checkpoint and pending-outbox resolution;
- repository-content retrieval, ranking, indexing, embedding, or memory lifecycle beyond two declared purpose-limited exact-source projections;
- agent registry, workforce coordination, general model routing, or tool brokerage beyond the separate disabled discovery and change-preparation purposes at the same provider edge;
- production-grade sandbox assurance, dynamic workers, models, tools, secrets, or network authority;
- evidence signing, production non-repudiation, retention execution, or supported backup and restore;
- durable cross-process Workbench resume, multi-user Workbench administration, Webex Messaging, Software System Twin, or a supported remote product API;
- live GitHub App activation, installation-token brokering, pull-request merge, or review-provider integration beyond the unconfigured bounded publication adapter;
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
