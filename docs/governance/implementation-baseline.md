# Forge implementation baseline

Status: Accepted
Owner: Forge Engineering and Architecture
Version: 0.5
Verified: 2026-07-22

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
| Trust contract candidates | Closed JSON Schemas for Execution Envelope, Context Package, Policy Decision, Approval Record, and Evidence Record | `schemas/forge/trust/v1` | Research |
| Executable trust contracts | Strict Go decoding, structural validation, exact authority binding, child-envelope narrowing, digest binding, evidence-subject binding, and normalized disclosure-safe failures for the five trust objects | `internal/contracts/trust`, `internal/contracts/failure`, and tests | Trial |
| Developer Trial identity admission | Provider-neutral assertion verification, audience and lifecycle validation, immutable normalized principal, and disclosure-safe provider failure normalization; only deterministic test providers exist | `internal/identity` and tests | Trial |
| Developer Trial scope admission | Single-test-tenant workspace membership resolution, server-derived immutable project scope, delayed repository binding discovery, and disclosure-safe cross-scope denial; storage and binding adapters remain test doubles | `internal/tenant` and tests | Trial |
| Developer Trial intent capture | Versioned direct-human intent capture, immutable `Captured` delivery run, Tenant-derived repository/classification binding, scope-bound atomic idempotency, and a concurrency-safe in-memory store behind a replaceable port | `internal/delivery` and tests | Trial |
| Assurance contract candidates | Closed JSON Schemas for journey scenarios and results | `schemas/forge/assurance/v1` | Research |
| M1 journey assurance candidates | Reviewed allowed and denied public-CLI scenarios, synthetic fixtures, and an execution plan that remains disabled until M1 | `assurance/journeys/m1` | Research |
| Offensive assurance safety candidate | Closed rules-of-engagement schema, disabled contract-only M1 configuration, digest-bound target/configuration, and negative safety fixtures | `schemas/forge/security/v1`, `security/assurance` | Research |
| Architecture module boundaries | Current module graph, planned worker boundary, negative dependency fixtures, and CI-control validation | `architecture/module-boundaries.json`, `architecture/testdata`, `scripts/check_architecture.py` | Trial |
| M1 threat and privacy candidates | Complete priority-threat mappings, phase-assigned negative contracts, Developer Trial decisions, and synthetic leakage fixtures | `security/threat-controls`, `security/privacy`, `security/fixtures` | Research |
| Deployment profile candidates | Closed profile schema plus Developer Trial and connected self-hosted reference manifests with phased conformance plans | `schemas/forge/deployment/v1`, `profiles` | Research |
| Architecture fitness checks | Validates document metadata, F0 governance, capability identifiers, schema sets, source integrity, navigation, local links, and developer-only Compose risks and constraints | `scripts/check_architecture.py` | Trial |

The default configuration sets `intelligence.enabled` to `false`. The executable trust contracts validate supplied objects but do not own or operate the capabilities that create, store, revoke, or distribute them. Identity, scope admission, and intent capture are currently callable Go modules only; they are not wired into the CLI and provide no live identity provider, durable tenant or delivery store, or repository provider. The current CLI does not authenticate users, capture intent through a public command, issue execution envelopes, retrieve context, invoke models or tools, coordinate agents, prepare repository changes, verify candidates, request approvals, publish pull requests, or persist evidence.

## Target capabilities not yet implemented

All capabilities in `architecture/capabilities.md` remain Research unless a later evidence record says otherwise. In particular, the repository does not yet implement:

- durable tenant, workspace, project, membership, repository-binding, or quota management beyond the bounded in-memory Developer Trial resolver;
- identity federation, credential acquisition, delegation, workload identity issuance, or revocation beyond assertion normalization and lifecycle checks;
- policy evaluation, decision issuance, revocation distribution, or distributed authorization enforcement beyond local contract checks;
- durable delivery-run state, transitions beyond in-memory `Captured`, transactional outbox, scheduling, cancellation, or recovery;
- authorized knowledge retrieval, Context Package assembly, or memory lifecycle;
- agent registry, workforce coordination, model routing, or tool brokerage;
- isolated worker enforcement for filesystem, network, secrets, time, compute, or cost;
- independent candidate verification or digest-bound approval;
- append-oriented integrity-protected evidence storage;
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
