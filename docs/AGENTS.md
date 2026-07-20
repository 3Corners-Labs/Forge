---
title: AGENTS
project: Forge
version: 0.1-draft
status: Draft
---

# AGENTS.md

## Purpose

This document defines the operational contract for all human and AI agents contributing to the Forge repository.

Every AI agent interacting with this repository SHALL comply with this document unless a more specific document explicitly overrides it.

---

# 1. Core Principles

An agent SHALL:

- Follow the Engineering Constitution.
- Treat the Git repository as the Source of Truth.
- Preserve architectural consistency.
- Prefer modifying existing documentation over creating duplicates.
- Produce deterministic, reviewable artefacts.
- Explain significant design decisions.

---

# 2. Repository Authority

Agents SHALL use repository knowledge in the following order:

1. Constitution
2. Architecture
3. ADRs
4. Engineering Handbook
5. Source Code
6. Tests

Conversation history SHALL NOT override repository content.

---

# 3. Engineering Behavior

Agents SHOULD:

- make the smallest safe change;
- avoid unnecessary complexity;
- prefer open standards;
- preserve backwards compatibility where practical;
- document architectural decisions.

Agents SHALL NOT:

- invent undocumented architecture;
- silently change public interfaces;
- bypass policy;
- embed secrets;
- duplicate business logic.

---

# 4. Documentation Rules

Every significant change SHOULD update the relevant documentation.

Documentation SHALL be written in English.

Markdown is the preferred documentation format.

Mermaid SHOULD be used for architecture visualizations.

---

# 5. Coding Standards

Agents SHALL:

- write readable code;
- prefer composition over inheritance;
- avoid vendor lock-in;
- include meaningful tests where appropriate;
- follow language-specific style guides.

---

# 6. AI-specific Rules

AI agents SHALL:

- identify assumptions;
- separate facts from proposals;
- minimize token usage;
- reuse repository knowledge;
- avoid unnecessary reasoning repetition.

Human approval remains mandatory for production changes.

---

# 7. Pull Request Expectations

Each proposed contribution SHOULD include:

- summary;
- rationale;
- affected components;
- documentation updates;
- testing evidence.

---

# 8. Definition of Success

A successful AI contribution:

- improves the repository;
- is understandable by humans;
- is reproducible;
- is testable;
- is auditable;
- aligns with the Forge architecture.

---

End of AGENTS.md (Draft)
