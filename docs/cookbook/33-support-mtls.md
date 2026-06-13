# Cookbook 33 — Support Site + mTLS (Enterprise)

## Support 页面

| 文件 | URL |
|------|-----|
| `deploy/site/support/index.html` | https://msgguard.app/support |
| `docs/legal/SUPPORT.md` | 法务/客服参考 |

本地预览：

```bash
python3 -m http.server 8888 --directory deploy/site
open http://localhost:8888/support/
```

iOS Settings / Help 已链接 `https://msgguard.app/support`。

## mTLS 架构

```
Client (admin cert) → Caddy mTLS :8443 → Gateway :8080
                         ↓ X-Client-Cert-Subject
                    /api/v1/admin/* 需 header + Bearer token
```

公开 API（classify/feedback/health）**不**强制 mTLS。

## 动手验收

```bash
# 1. 生成证书
bash deploy/mtls/gen-certs.sh

# 2. 启动 mTLS compose
bash deploy/tiers/tier4-mtls-compose.sh

# 3. Gateway 中间件单测
cd pkg/httpauth && go test ./...

# 4. Admin 无 cert → 401（经 Caddy）
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/token \
  -H 'Content-Type: application/json' -d '{"roles":["admin"]}' | python3 -c "import sys,json;print(json.load(sys.stdin)['access_token'])")
curl -sk https://localhost:8443/api/v1/admin/metrics/summary -H "Authorization: Bearer $TOKEN"
# 期望 401（无 client cert）

# 5. Helm mTLS values
helm template msgguard deploy/helm/msgguard -f deploy/helm/msgguard/values-mtls.yaml | rg MTLS_ADMIN
```

## Helm 生产

```bash
helm upgrade msgguard deploy/helm/msgguard \
  -f deploy/helm/msgguard/values-prod.yaml \
  -f deploy/helm/msgguard/values-mtls.yaml
```

见 `docs/runbooks/mtls-rotation.md` 轮换流程。
