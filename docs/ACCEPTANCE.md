# MsgGuard Acceptance Checklist

逐项验收，每项含命令与期望输出。

## WS0 — Repository
- [ ] `git clone https://github.com/EthanShen10086/MsgGuard.git` 成功

## ML Pipeline
- [ ] `cd ml && make data` → `data/processed/all.csv` 400+ rows
- [ ] `cd ml && make train` → `output/bayes_pipeline.joblib` 存在
- [ ] `cd ml && make benchmark` → `gate_passed=True`
- [ ] `cd ml && make infer TEXT="免费贷款无抵押"` → `"label": "spam"` 或 `1`

## Backend
- [ ] `cd services/gateway && go build ./cmd/server` 成功
- [ ] `curl localhost:8080/health` → `ok`
- [ ] `curl -X POST localhost:8080/api/v1/feedback -H 'Content-Type: application/json' -d '{"body":"test","label":"ham"}'` → JSON id
- [ ] `curl localhost:8080/metrics` → Prometheus text

## iOS
- [ ] `cd apps/ios && bash setup.sh && xcodebuild -scheme MsgGuard-iOS -destination 'platform=iOS Simulator,name=iPhone 16' build` → BUILD SUCCEEDED

## Deploy
- [ ] `./deploy/tiers/tier1-compose.sh` → services up
- [ ] `helm template msgguard deploy/helm/msgguard` → Deployment YAML

## Phase 2 — Architecture & Product Flywheel
- [x] `pkg/app/container.go` exists; gateway uses `app.NewContainer`
- [x] `curl -X POST localhost:8080/api/v1/auth/token` returns access_token
- [x] `GET /api/v1/feedback` without token → 401
- [x] `python ml/product/aggregate_metrics.py` → `ml/product/reports/weekly.json`
- [x] iOS App Group: `analytics.jsonl`, `crash_reporter.installed` after launch
- [x] `docs/COMMERCIAL_READINESS.md` 存在
- [x] Cookbook 00–29 每章有验收命令

## Phase 2.1 — MongoDB / Model Auth / Entitlements
- [x] `pkg/adapters/mongodb` — feedback/rules/analytics Port 实现
- [x] `DATABASE_DRIVER=mongodb` 切换（见 Cookbook 30）
- [x] Model 服务 `/api/v1/models/register` 自身 RBAC（401/201）
- [x] iOS `EntitlementProviding` + Keychain 持久化
- [x] `make verify` / `bash scripts/verify.sh` 全量验收

## Phase 2.2 — Compose MongoDB + StoreKit
- [x] `deploy/docker-compose.yml` 含 MongoDB 服务
- [x] `make tier1-mongodb` / `deploy/tiers/tier1-mongodb.sh`
- [x] `StoreManager` — 购买 / 恢复 / Transaction 监听
- [x] `Products.storekit` 本地 StoreKit 测试配置
- [x] Cloud LLM 需 Pro  entitlement（Settings 禁用提示）
- [x] Cookbook 31

## Phase 2.3 — App Store / Privacy / Helm MongoDB
- [x] `docs/app-store/metadata.md` + `review-notes.md`
- [x] `docs/legal/PRIVACY.md` + static `deploy/site/privacy/index.html`
- [x] Helm `values-mongodb.yaml` + `DATABASE_DRIVER` in ConfigMap
- [x] `make tier3-mongodb` / `deploy/tiers/tier3-helm-mongodb.sh`
- [x] Cookbook 32

## Phase 2.4 — Support + mTLS
- [x] `deploy/site/support/index.html` + `docs/legal/SUPPORT.md`
- [x] `deploy/mtls/gen-certs.sh` + Caddy mTLS + `docker-compose.mtls.yml`
- [x] Gateway `MTLS_ADMIN_REQUIRED` middleware (`pkg/httpauth/clientcert.go`)
- [x] Helm `values-mtls.yaml` + ingress client-auth annotations
- [x] `make tier4-mtls` / Cookbook 33

## Phase 2.5 — Gap Closure (Product + Commercial Wiring)
- [x] CoreMLClassifier loads App Group `.mlmodelc`; ModelUpdateService OTA + SHA256
- [x] OTPGuard + adversarial FPR benchmark gate
- [x] Pro entitlements: customRules / advancedStats / cloudLLM gating
- [x] Gateway QuotaStore + FeatureFlag wired; admin quota/flags API
- [x] NATS flywheel worker + feedback trigger
- [x] Tier matrix doc + values-prod + caddy-site in compose
- [x] ASC checklist + Fastfile + screenshots guide

## Phase 2.6 — P2/P3 Long-term Items
- [x] Shadow mode Prometheus metrics + admin stats + Grafana provisioning
- [x] Rules CDN: ETag/304, version routes, iOS incremental sync, sync-rules.sh
- [x] Multi-locale model publish + iOS dynamic locale OTA paths
- [x] iCloud autoSync (NSUbiquitousKeyValueStore) + Pro gating
- [x] CallDirectory extension skeleton + Settings guidance
- [x] cert-manager Helm Certificate templates (prod values)
- [x] flywheel worker binary removed from git + .gitignore
