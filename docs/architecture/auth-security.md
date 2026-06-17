# Auth & Security Architecture

**Last updated:** 2026-06-17

## Auth Layers

| Layer | Mechanism | Scope |
|-------|-----------|-------|
| Device token | HMAC JWT via `/api/v1/auth/device` | Feedback, analytics, entitlements from iOS |
| App Attest | `APP_ATTEST_REQUIRED` + `X-App-Attest` header | Optional device route hardening |
| Token revocation | Redis `msgguard:revoked:*` when `REDIS_URL` set | Logout / compromise response |
| Admin Bearer | `/api/v1/auth/token` (bootstrap gated in prod) | Admin API, feedback list |
| Model RBAC | Bearer + `models:write` | Model register |
| mTLS | Caddy client cert or `X-Client-Cert-Subject` | Admin paths in Tier 4 |

## Production Gates

- `AUTH_SECRET` validated at startup — **fails in `MSGGUARD_ENV=production`** if weak/missing
- `AUTH_BOOTSTRAP_ENABLED=false` in prod disables open token issuance
- `MODEL_DOWNLOAD_AUTH_REQUIRED=true` protects model downloads
- `MTLS_ADMIN_REQUIRED=true` for admin routes behind ingress

## RBAC Roles

Defined in `pkg/adapters/memory/auth.go` + voxera-kit authorizer:

- `admin` — metrics, flags, quota, shadow stats
- `ml_engineer` — feedback read, model register
- Device scope — single-device ingest

## Privacy

- PII redaction middleware on gateway (phone, email patterns)
- `DELETE /api/v1/privacy/me?device_id=` — analytics erasure (GDPR-style)
- Audit log on feedback create (voxera-kit audit)

## References

- [Threat model](../security/THREAT_MODEL.md)
- [Cookbook: Auth RBAC](../cookbook/27-auth-rbac.md)
- `pkg/httpauth/` — shared middleware
