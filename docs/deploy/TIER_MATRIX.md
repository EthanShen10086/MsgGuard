# Tier Deployment Matrix

Switch **only** config + tier script — no business code changes.

| Tier | Script | DB Driver | Static Site | mTLS | Helm |
|------|--------|-----------|-------------|------|------|
| 0 Local | `tier0-local.sh` | memory | — | — | — |
| 1 Compose | `tier1-compose.sh` | postgres | `:8088` caddy-site | — | — |
| 1 MongoDB | `tier1-mongodb.sh` | mongodb | `:8088` | — | — |
| 2 GPU | `tier2-compose-gpu.sh` | postgres | — | — | — |
| 3 Helm Staging | `tier3-helm-staging.sh` | postgres (values) | external | optional | staging |
| 3 Helm MongoDB | `tier3-helm-mongodb.sh` | mongodb overlay | external | optional | staging+mongo |
| 4 mTLS Compose | `tier4-mtls-compose.sh` | postgres | `:8088` | Caddy :8443 | — |
| 4 Helm Prod | `tier4-helm-prod.sh` | see values-prod | external | values-mtls overlay | prod |

## Environment variables (Gateway)

| Variable | Purpose |
|----------|---------|
| `DATABASE_DRIVER` | `postgres` / `mongodb` / `memory` |
| `DATABASE_DSN` | Connection string |
| `AUTH_SECRET` | HMAC token signing |
| `MTLS_ADMIN_REQUIRED` | Admin path client cert / header check |
| `MTLS_CLIENT_HEADER` | e.g. `X-Client-Cert-Subject` behind Caddy |

## Verify

```bash
make verify
docker compose -f deploy/docker-compose.yml ps
curl http://localhost:8088/support/
```
