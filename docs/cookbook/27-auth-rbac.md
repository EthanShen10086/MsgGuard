# Auth & RBAC

## 对应原始需求
管理员身份权限、API 访问控制

## 涉及文件
- `pkg/adapters/memory/auth.go` — HMAC JWT-like tokens + RBAC
- `services/gateway/internal/handler/admin.go` — token 签发、metrics
- `services/gateway/internal/middleware/auth_prod.go` — production bootstrap gate
- `services/admin/cmd/server/main.go` — standalone quota / feature flags
- `apps/admin/` — admin web SPA (Phase 4)

## 动手验收
```bash
# 获取 admin token (dev / AUTH_BOOTSTRAP_ENABLED=true only)
curl -X POST localhost:8080/api/v1/auth/token \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"admin1","roles":["admin"]}'

# 无 token 访问 feedback list → 401
curl localhost:8080/api/v1/feedback

# 带 token 访问
TOKEN=<paste>
curl localhost:8080/api/v1/feedback -H "Authorization: Bearer $TOKEN"
```
**期望输出：** 401 without token; JSON array with token

## Production Auth Notes

| Setting | Dev default | Production |
|---------|-------------|------------|
| `MSGGUARD_ENV` | unset | `production` |
| `AUTH_BOOTSTRAP_ENABLED` | true (non-prod) | **false** — disable `/api/v1/auth/token` |
| `AUTH_SECRET` | dev default | **rotate** — strong random secret in K8s |
| `AUTH_DEVICE_TOKEN_ENABLED` | true | true — iOS device tokens |
| `MODEL_DOWNLOAD_AUTH_REQUIRED` | false | **true** — protect model downloads |
| `MTLS_ADMIN_REQUIRED` | false | true in Tier 4 — admin paths need client cert |

### Token issuance in production

1. Do **not** expose bootstrap token endpoint publicly.
2. Issue admin tokens via internal ops (VPN + mTLS) or future OIDC (Phase 10).
3. Admin web (`apps/admin`) expects pasted Bearer token — no bootstrap in prod build.

### Device vs admin

- **Device token** (`POST /api/v1/auth/device`): mobile clients; scoped to feedback/analytics ingest.
- **Admin Bearer**: `admin` or `ml_engineer` roles for list/metrics/flags.

### mTLS (Tier 4)

```bash
# Behind Caddy with client cert forwarded as header
MTLS_ADMIN_REQUIRED=true
MTLS_CLIENT_HEADER=X-Client-Cert-Subject
```

Admin paths: `/api/v1/admin/*`

## Debug 指南
- 403 → 检查 role 是否含 admin/ml_engineer
- 403 on `/auth/token` in prod → expected; use pre-issued token
- 设置 `AUTH_SECRET` 环境变量（生产必改）

## 扩展阅读
- `docs/api/openapi.yaml` securitySchemes
- `docs/architecture/auth-security.md`
- `docs/security/THREAT_MODEL.md`
