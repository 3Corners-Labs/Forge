# Forge architecture

Forge has three layers:

1. **Developer experience** — CLI, devcontainers, templates and documentation.
2. **Engineering runtime** — local services, telemetry and artifact infrastructure.
3. **Automation** — CI/CD, policy gates, release automation and programming agents.

The current v0.1 implementation intentionally uses Docker Compose for the golden-path local runtime. Kubernetes packaging follows after the contracts and bootstrap flow are validated.
