# Contributing

1. Start from an issue.
2. Use an RFC for significant proposals and an ADR for accepted decisions.
3. Create a short-lived branch from `main`.
4. Implement the smallest coherent change.
5. Sign off every commit (`git commit -s`), see below.
6. Run `task verify`.
7. Open a pull request using the repository template.
8. Merge only after required checks and review succeed.

Use Conventional Commits, for example:

```text
feat(trust): add TLS policy validation
fix(compose): correct OpenBao healthcheck
docs(architecture): document certificate provider contract
```

Code, tests, contracts, diagrams, and ADR/RFC documentation must be updated in the same pull request when they describe the same change.

## Licence and sign-off

Forge is licensed under the Apache License, Version 2.0 (`LICENSE`, `NOTICE`); contributions are accepted under the same licence. By signing off a commit you certify the Developer Certificate of Origin 1.1 (`DCO`) for it. Use `git commit -s`, which adds:

```text
Signed-off-by: Your Name <you@example.com>
```

The email address must match the commit author; a GitHub no-reply address is fine. The `DCO` check fails a pull request with an unsigned commit; fix it with `git commit --amend -s` or `git rebase --signoff main` and force-push. See ADR-0026.
