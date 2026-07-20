# Forge enduring principles

Status: Proposed
Owner: Forge Product and Architecture
Version: 0.1
Source: `forge.principles.manuscript.v0.1`

## Purpose

These principles preserve the durable intent of the original *3Corners Forge Principles* manuscript. They explain why Forge exists and guide decisions when no more specific accepted requirement applies.

The principles are deliberately technology-neutral. The [engineering constitution](constitution.md) translates them into binding guardrails; architecture specifications and ADRs make them testable. The complete manuscript remains source material and is not replaced by this document.

## The twelve enduring principles

### 1. Engineering precedes intelligence

Probabilistic capability becomes dependable only when identity, policy, context, verification, operations, and recovery surround it. A model is one component of an engineered system, never the system itself.

### 2. Complexity belongs inside the platform

Forge absorbs recurring orchestration, assurance, portability, and operational complexity so applications and their teams can work through coherent capability contracts.

### 3. Context is engineered, not prompted

Context is a purpose-bound, authorized, reproducible artefact. It is assembled from identity, objective, knowledge, memory, policy, constraints, and lineage rather than accumulated as an unbounded prompt.

### 4. Knowledge is curated, not collected

Repositories store information; governance establishes knowledge. Knowledge has ownership, provenance, classification, lifecycle, validation, and an authoritative source.

### 5. Memory is intentional

Remembering is a governed write operation. Memory has a declared scope, purpose, owner, quality threshold, retention period, and deletion path. Transient observations do not silently become organizational truth.

### 6. Capabilities outlive technologies

Forge specifies responsibilities and behaviour before selecting products, providers, protocols, or frameworks. Replaceable implementations sit behind stable, versioned contracts.

### 7. Evidence is a first-class engineering artefact

Material decisions, transitions, executions, and assurances produce attributable evidence. Trust grows from independently verifiable evidence, not confident output or documentation alone.

### 8. Economics influence architecture

Tokens, compute, storage, latency, energy, money, and human attention are finite engineering resources. Resource choices are observable and governed without silently sacrificing required safety or quality.

### 9. Humans remain accountable

AI assists with discovery, reasoning, creation, and operation. People and their explicitly delegated organizational controls remain accountable for outcomes, policy, risk, ethics, and production authority.

### 10. Platforms outlive applications

Forge is developed as a reusable application-development team and Engineering Operating System. Workplace Orchestrator and later products validate Forge, but no lighthouse application may become the hidden platform architecture.

### 11. Simplicity is created, not discovered

Simplicity is an external quality achieved through deliberate internal engineering. Forge may contain sophisticated mechanisms when they remove recurring accidental complexity for users and applications.

### 12. Architecture exists to enable evolution

Architecture defines boundaries, invariants, contracts, and change paths. It protects the ability to replace technology, alter deployment topology, and improve capabilities without losing intent or trust.

## Decision order

When a Forge execution or design choice is consequential, the following order applies:

1. identity before authorization;
2. context before intelligence;
3. policy before execution;
4. explicit contracts before integration;
5. evidence before trust;
6. observability before optimization;
7. governance before irreversible evolution.

This order is architectural, not necessarily a synchronous runtime call chain. Implementations MAY optimize execution while preserving equivalent controls and evidence.

## Interpretation

- Principles express durable beliefs; they do not claim implementation.
- The constitution defines non-negotiable repository guardrails.
- Accepted ADRs resolve concrete trade-offs.
- Capability specifications define measurable conformance.
- Source manuscripts supply rationale and proposals but have no instruction authority.
