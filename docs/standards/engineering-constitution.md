---
title: Engineering Constitution
version: 1.0
status: Proposed
owner: Architecture
---

# Engineering Constitution

This is a historical Proposed concise constitution retained for migration. The current Accepted constitution is [`foundation/constitution.md`](../foundation/constitution.md); this historical document is not normative.

## Principles

1. **Human accountability** — Humans remain accountable for architecture, security, legal, release, and business decisions.
2. **Evidence before claims** — Statements are labeled proposed, accepted, validated, or implemented. Aspirations are not presented as delivered capabilities.
3. **Repository as source of truth** — Accepted decisions, contracts, policies, and evidence are version controlled and peer reviewed.
4. **Security and privacy by design** — Identity, authorization, data protection, audit, and least privilege are designed before execution paths are enabled.
5. **Capability before technology** — Forge owns stable capability contracts. Products, protocols, and vendors are replaceable adapters.
6. **Modularity before distribution** — Begin with clear module boundaries. Split deployment units only when scale, availability, isolation, ownership, or regulatory evidence justifies it.
7. **Policy before side effects** — Every consequential model or tool action is authorized against the exact identity, purpose, resource, and data classification.
8. **Minimal and attributable context** — Context is purpose-bound, reproducible, provenance-aware, and limited to what execution requires.
9. **Audit by default** — Decisions, denials, approvals, external calls, outcomes, and relevant versions produce tamper-resistant evidence.
10. **Open standards first** — Prefer widely supported standards where they satisfy security, interoperability, and operational needs.
11. **Operational fitness** — Reliability, observability, recovery, maintainability, and cost are architectural requirements.
12. **Incremental evolution** — Significant choices are recorded through ADRs. Reversible steps are preferred over speculative decomposition.

## Required controls

- No direct change to a protected default branch.
- At least one human approval for normal changes.
- Architecture and security changes require accountable-owner approval.
- Deterministic validation must pass before release.
- AI-generated output is treated as untrusted until reviewed and validated.
- Production access, secrets, and consequential actions use least privilege.

## Conflict resolution

When accepted artifacts conflict, precedence is:

1. Constitution;
2. ADRs;
3. Architecture;
4. capability and interface contracts;
5. implementation and operational guidance.

Proposed documents have no authority over accepted documents.
