# Programming Agent Instructions

- Work only from a scoped issue, RFC, or explicit human instruction.
- Use a feature branch and draft pull request; never push directly to `main`.
- Human approval is mandatory before merge.
- Do not access production systems or secrets.
- Do not weaken tests, scanners, security, privacy, accessibility, or architecture gates.
- Update code, tests, contracts, architecture, and documentation together.
- Run `task verify` and report all commands and limitations.
- Changes to `architecture/**`, `adrs/**`, `rfcs/**`, `internal/trust/**`, `deploy/**`, `.github/workflows/**`, or public contracts require explicit architecture review.
