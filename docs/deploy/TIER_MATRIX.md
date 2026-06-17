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

## Full-Stack Helm (Phase 4+)

The Helm chart at `deploy/helm/msgguard` can deploy the **full backend stack** in one release:

- Gateway Deployment + Service + Ingress
- ConfigMap from `values.yaml` / overlays (`values-mongodb.yaml`, `values-mtls.yaml`, `values-prod.yaml`)
- HPA, cert-manager Certificate (optional)
- Prometheus scrape annotations (when metrics enabled)

**Not yet in chart:** admin SPA static assets, Caddy site bundle — serve `deploy/site/` and `apps/admin/dist/` via external CDN or add an optional `site` subchart (Planned).

```bash
# Staging full stack
helm upgrade --install msgguard deploy/helm/msgguard \
  -f deploy/helm/msgguard/values.yaml \
  -f deploy/helm/msgguard/values-mongodb.yaml

# Production overlay
helm upgrade --install msgguard deploy/helm/msgguard \
  -f deploy/helm/msgguard/values-prod.yaml \
  -f deploy/helm/msgguard/values-mtls.yaml
```

Pair with external static hosting for `/privacy/`, `/support/`, `/pricing.html`, `/status.html`.

## Environment variables (Gateway)

| Variable | Purpose |
|----------|---------|
| `DATABASE_DRIVER` | `postgres` / `mongodb` / `memory` |
| `DATABASE_DSN` | Connection string |
| `AUTH_SECRET` | HMAC token signing |
| `AUTH_BOOTSTRAP_ENABLED` | Allow `/api/v1/auth/token` (disable in prod) |
| `MTLS_ADMIN_REQUIRED` | Admin path client cert / header check |
| `MTLS_CLIENT_HEADER` | e.g. `X-Client-Cert-Subject` behind Caddy |

## Verify

```bash
make verify
docker compose -f deploy/docker-compose.yml ps
curl http://localhost:8088/support/
```
