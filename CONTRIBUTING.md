# Contributing

1. Start from an issue.
2. Use an RFC for significant proposals and an ADR for accepted decisions.
3. Create a short-lived branch from `main`.
4. Implement the smallest coherent change.
5. Run `task verify`.
6. Open a pull request using the repository template.
7. Merge only after required checks and review succeed.

Use Conventional Commits, for example:

```text
feat(trust): add TLS policy validation
fix(compose): correct OpenBao healthcheck
docs(architecture): document certificate provider contract
```

Code, tests, contracts, diagrams, and ADR/RFC documentation must be updated in the same pull request when they describe the same change.
