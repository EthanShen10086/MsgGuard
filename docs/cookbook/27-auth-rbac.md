# Auth & RBAC

## 对应原始需求
管理员身份权限、API 访问控制

## 涉及文件
- `pkg/adapters/memory/auth.go` — HMAC JWT-like tokens + RBAC
- `services/gateway/internal/handler/admin.go` — token 签发、metrics
- `services/admin/cmd/server/main.go` — quota 白名单、feature flags

## 动手验收
```bash
# 获取 admin token
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

## Debug 指南
- 403 → 检查 role 是否含 admin/ml_engineer
- 设置 `AUTH_SECRET` 环境变量（生产必改）

## 扩展阅读
- `docs/api/openapi.yaml` securitySchemes
