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
