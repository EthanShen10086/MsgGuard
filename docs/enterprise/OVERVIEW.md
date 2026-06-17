# Enterprise & SMB Roadmap

MsgGuard ships consumer-first (iOS Message Filter, optional cloud LLM). Enterprise tiers add tenant isolation, SSO, and managed rollout without forking the core filter engine.

## Tiers

| Tier | Audience | Highlights |
|------|----------|------------|
| **SMB** | Small teams, clinics | Shared gateway, device tokens, basic analytics export |
| **Enterprise** | Regulated orgs | `tenant_id` on feedback/analytics, OIDC SSO, mTLS admin, private model registry |

## Phase 7 stubs (current)

- **Tenant column** — `tenant_id` optional on `feedback` and `analytics_events` (Postgres migration `001_add_tenant_id.sql`). Handlers accept `tenant_id` in JSON; empty means consumer/default tenant.
- **OIDC** — `pkg/httpauth/oidc_provider.go` authorization-code flow for admin; set `OIDC_ISSUER`, `OIDC_CLIENT_ID`, `OIDC_CLIENT_SECRET`, `OIDC_ADMIN_EMAILS` or `OIDC_ADMIN_DOMAIN`. Production: `OIDC_ENFORCE_ADMIN=true` disables bootstrap token.
- **Admin** — Model promote/rollback and canary flags via `/api/v1/admin/models/*` and `ml/scripts/canary_rollout.py`.

## Near-term roadmap

1. **Identity** — OIDC authorization code flow for admin console; map `sub` → `tenant_id` claim.
2. **Data plane** — Row-level filters on feedback/analytics list APIs; per-tenant S3 prefixes for model artifacts.
3. **Policy** — Tenant-scoped rule bundles and blocklists synced from gateway.
4. **Compliance** — Audit log export, data retention jobs, BYOK for LLM keys.
5. **Android enterprise** — Work profile SMS filter + MDM config profile for gateway URL and tenant header.

## Deployment notes

- Enable `MTLS_ADMIN_REQUIRED` for admin APIs in production enterprise installs.
- Run Postgres migrations before enabling multi-tenant ingestion.
- Keep consumer and enterprise traffic on separate gateway replicas when mixing SSO and device-token auth.
