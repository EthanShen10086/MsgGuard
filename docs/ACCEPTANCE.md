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

## Docs
- [ ] `docs/api/openapi.yaml` 存在
- [ ] `docs/COMMERCIAL_READINESS.md` 存在
- [ ] Cookbook 00–26 每章有验收命令
