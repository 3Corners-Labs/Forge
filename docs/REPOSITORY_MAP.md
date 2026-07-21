---
title: Repository Map
version: 1.0
status: Proposed
owner: Engineering
---

# Repository Map

> **Historical migration source:** this map predates the current `docs/decisions/` and documentation navigation. Use [`docs/README.md`](README.md) for current documentation organization and the [implementation baseline](governance/implementation-baseline.md) for implemented repository state.

This map describes documentation responsibilities without asserting implementation structure.

| Area | Responsibility |
|---|---|
| `README.md` | Concise entry point, scope, and maturity |
| `SYSTEM.md` | Instructions for AI-assisted engineering |
| `docs/foundation/` | Proposed vision, principles, masterplan, and constitutional review target |
| `docs/architecture/` | Proposed target architecture and capability contracts |
| `docs/governance/` | Readiness, delivery, lifecycle, migration, and acceptance controls |
| `docs/decisions/` | Proposed and future Accepted architecture decisions and index |
| `docs/adr/` | Historical ADR location pending migration |
| `docs/rfc/` | Proposals under review and index |
| `docs/context/` | Task-specific authoritative context routing |
| `docs/reference/` | Shared terminology |
| `docs/reference-applications/` | Product dossiers that validate but do not define Forge core |
| `docs/source-material/` | Immutable, non-normative source manuscripts and retrieval maps |

## Ownership rules

- Every normative document declares status and owner.
- Only accepted documents are authoritative.
- Duplicate normative content is replaced with a link to the canonical source.
- Significant architectural changes create or update an ADR.
- Diagrams and prose change together.
