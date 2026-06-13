# Cookbook 32 — Helm MongoDB + Privacy Site

## Helm MongoDB 切换

仅改 values，无需改业务代码：

```bash
# 渲染检查
helm template msgguard deploy/helm/msgguard \
  -f deploy/helm/msgguard/values-staging.yaml \
  -f deploy/helm/msgguard/values-mongodb.yaml \
  | rg DATABASE_DRIVER
# 期望: mongodb

# 部署 staging（需集群内 MongoDB + Redis）
bash deploy/tiers/tier3-helm-mongodb.sh
```

`values-mongodb.yaml` 覆盖：
- `database.driver: mongodb`
- `database.dsn` — MongoDB 连接串
- `auth.secret` — 与 model 服务一致的 HMAC secret

## 隐私政策站点

| 文件 | 用途 |
|------|------|
| `docs/legal/PRIVACY.md` | 完整政策（App Store / 法务） |
| `deploy/site/privacy/index.html` | 静态页 → https://msgguard.app/privacy |

Caddy 托管：将 `deploy/site` 挂载到 `/var/www/msgguard`。

```bash
# 本地预览
python3 -m http.server 8888 --directory deploy/site
open http://localhost:8888/privacy/
```

## App Store

- 元数据草稿：`docs/app-store/metadata.md`
- 审核说明：`docs/app-store/review-notes.md`

## 验收

```bash
test -f docs/legal/PRIVACY.md
test -f deploy/site/privacy/index.html
helm template msgguard deploy/helm/msgguard -f deploy/helm/msgguard/values-mongodb.yaml | rg -q mongodb
```
