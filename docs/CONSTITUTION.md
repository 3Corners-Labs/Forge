---
title: Engineering Constitution
project: Forge
document: CONSTITUTION.md
version: 1.0-draft
status: Draft
---

# Forge Engineering Constitution

> **Historical Draft:** this document is preserved as consolidation source. The current Accepted constitution is [`foundation/constitution.md`](foundation/constitution.md); this historical document is not normative.

## Purpose

The Constitution establishes the architectural values that every contributor—human or AI—MUST follow. These principles are intentionally stable and are expected to change only under exceptional circumstances.

---

# Principle 1 — Human Accountability

Humans remain accountable for all architectural, security, legal and business decisions.

AI systems MAY assist, recommend and automate, but SHALL NOT become the authoritative decision maker.

---

# Principle 2 — Repository as the Source of Truth

The Git repository is the canonical source of engineering knowledge.

Architecture, policies, ADRs, documentation and source code SHALL be version controlled.

Conversation history SHALL NEVER supersede repository content.

---

# Principle 3 — Everything as Code

Forge treats engineering artefacts as version-controlled assets.

This includes:

- Knowledge as Code
- Architecture as Code
- Context as Code
- Policy as Code
- Infrastructure as Code
- Documentation as Code

---

# Principle 4 — Vendor Independence

The platform SHALL depend on capabilities rather than vendors.

Every external dependency SHOULD be replaceable through stable abstractions.

No AI provider is considered mandatory.

---

# Principle 5 — Open Standards First

Whenever practical, Forge SHALL adopt open and widely supported standards.

Examples include:

- OpenAPI
- OpenTelemetry
- OAuth 2.0
- OpenID Connect
- Model Context Protocol (MCP)

---

# Principle 6 — Security by Design

Security SHALL be incorporated during architecture and implementation rather than added afterwards.

Identity, authorization, auditing and encryption are platform capabilities.

---

# Principle 7 — Privacy by Design

Only the minimum information required to complete a task SHOULD be processed.

Personal and sensitive information SHALL be handled according to applicable legislation and organizational policy.

---

# Principle 8 — Context as Code

Context SHALL be engineered rather than improvised.

Context packages SHOULD be:

- deterministic;
- reproducible;
- auditable;
- minimal.

---

# Principle 9 — Knowledge as Code

Knowledge is an engineering asset.

Engineering knowledge SHALL be:

- version controlled;
- peer reviewed;
- searchable;
- reusable.

---

# Principle 10 — Audit by Default

Every significant platform action SHOULD produce evidence.

Evidence MAY include:

- policy evaluations;
- execution logs;
- architectural decisions;
- model selections;
- test results.

---

# Principle 11 — Economic Intelligence

Forge SHALL select the least expensive solution that satisfies functional, security and governance requirements.

Capability takes precedence over vendor branding.

---

# Principle 12 — Token & Compute Stewardship

Tokens, compute, storage, energy and human attention are valuable resources.

Solutions SHOULD minimize unnecessary consumption while preserving quality.

---

# Principle 13 — Continuous Validation

Engineering quality SHALL be continuously verified through automated and human validation.

Testing extends beyond software to documentation, policies and AI behaviour.

---

# Principle 14 — Explainability

Architectural decisions SHOULD be understandable by engineers.

When practical, systems SHOULD explain:

- why a decision was made;
- which policies applied;
- which evidence supported the outcome.

---

# Principle 15 — Sustainable Engineering

Forge prioritizes maintainability over short-term convenience.

Design decisions SHOULD favor simplicity, modularity and long-term operability.

---

# Compliance

All repository artefacts SHOULD align with this Constitution.

When conflicts arise:

1. Constitution
2. ADRs
3. Architecture
4. Implementation

The higher-order document prevails.

---

_End of CONSTITUTION.md_
