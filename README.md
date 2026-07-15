# 3Corners Forge

Forge is the reusable engineering foundation for 3Corners products. It provides a portable local platform, a Go CLI, documentation-as-code, CI guardrails and a future integration point for online and offline programming agents.

## Quick start

```bash
go run ./cmd/forge doctor
go run ./cmd/forge up
```

Local services:

- PostgreSQL: `localhost:5432`
- Keycloak: `http://localhost:8081`
- OpenBao: `http://localhost:8200`
- MinIO: `http://localhost:9001`
- Grafana: `http://localhost:3000`
- Prometheus: `http://localhost:9090`
- Mailpit: `http://localhost:8025`

Development credentials are intentionally insecure and must never be reused outside local development.

## Verification

```bash
task verify
```
