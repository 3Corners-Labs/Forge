# Engineering Constitution v0.1

## Principles

1. Humans remain accountable for architecture, security, releases and consequential changes.
2. AI agents create reviewable artifacts; they never bypass protected branches or approval gates.
3. Privacy, security and accessibility are release criteria, not optional reviews.
4. Prefer open standards and replaceable components over vendor-specific coupling.
5. Start modular; split deployment units only when scale, ownership or risk justifies it.
6. Every externally observable behavior has an automated test and operational signal.
7. Documentation, schemas and architecture decisions are versioned with the code.
8. Community users receive a secure and useful baseline; security is not an enterprise-only feature.
9. Build once and promote immutable artifacts across environments and release channels.
10. Claims are labeled as proven, validated, proposed or experimental.

## Pull requests

- No direct commits to the protected default branch.
- At least one human approval for normal changes.
- Security-sensitive and architectural changes require CODEOWNER approval.
- Tests, linting, dependency review and secret scanning must pass.

## AI programming agents

- Agents receive least-privilege context packages and tools.
- Source code or customer data is routed only to providers allowed by data-classification policy.
- Agent output is treated as untrusted until validated by deterministic tooling and human review.
- Model, prompt, context sources, tools and costs are auditable.
