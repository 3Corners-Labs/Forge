---
title: Repository Map
version: 1.0
status: Proposed
owner: Engineering
---

# Repository Map

This map describes documentation responsibilities without asserting implementation structure.

| Area | Responsibility |
|---|---|
| `README.md` | Concise entry point, scope, and maturity |
| `SYSTEM.md` | Instructions for AI-assisted engineering |
| `docs/standards/` | Canonical engineering constitution and standards |
| `docs/ARCHITECTURE.md` | Normative logical architecture after acceptance |
| `docs/SYSTEM_OVERVIEW.md` | Accessible architectural orientation |
| `docs/CAPABILITY_CATALOG.md` | Capability responsibilities and maturity |
| `docs/SECURITY.md` | Security architecture baseline |
| `docs/adr/` | Accepted architecture decisions and index |
| `docs/rfc/` | Proposals under review and index |
| `docs/architecture/decisions/` | Supporting decision material; accepted decisions belong in `docs/adr/` |
| `docs/architecture/diagrams/` | Reusable diagrams whose source is reviewable |
| `docs/architecture/threat-models/` | Approved threat models and data-flow analyses |
| `docs/architecture/fitness-functions/` | Automated architectural constraints |

## Ownership rules

- Every normative document declares status and owner.
- Only accepted documents are authoritative.
- Duplicate normative content is replaced with a link to the canonical source.
- Significant architectural changes create or update an ADR.
- Diagrams and prose change together.
