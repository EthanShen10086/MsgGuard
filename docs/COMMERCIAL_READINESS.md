# Commercial Readiness

## Startup (Tier 0–1)
- [x] Seed dataset + benchmark gate in CI
- [x] docker-compose single VPS
- [x] App Store metadata draft (`docs/app-store/metadata.md`)
- [x] Privacy policy (`docs/legal/PRIVACY.md` + `deploy/site/privacy/`)
- [x] Support page (`docs/legal/SUPPORT.md` + `deploy/site/support/`)

## Growth (Tier 2–3)
- [x] PostgreSQL + Redis wired (via Container config)
- [x] Helm staging environment
- [x] Shadow mode metrics endpoint
- [x] Grafana dashboards (msgguard + product)
- [x] Data flywheel cron script

## Enterprise (Tier 4)
- [x] Multi-replica HPA (Helm template)
- [x] mTLS (Caddy + gateway admin middleware, Helm values-mtls)
- [x] Cross-region backup script (`deploy/ops/backup.sh`)
- [x] Audit logging (voxera-kit audit on feedback)
- [x] A/B experiment via featureflag (admin service)
- [x] Runbooks for LLM/DB/model rollback
- [x] Auth/RBAC + admin API
- [x] Product metrics flywheel (METRICS.md + aggregate_metrics.py)

## Switching Guide
Only change `deploy/config.*.yaml` and run matching `deploy/tiers/*.sh` — no business code changes.

## Architecture (Phase 2)
- `pkg/app/container.go` — unified DI / IoC root
- All stores wired via config: feedback, rules, analytics, cache, queue, model registry
- `pkg/adapters/mongodb` — alternate DB driver (Cookbook 30)
- `pkg/httpauth` — shared RBAC middleware (gateway + model service)
- iOS `EntitlementProviding` — ScrollCap-aligned subscription entitlements
- `scripts/verify.sh` — one-shot acceptance runner

## Phase 2.5 — Product + Commercial Wiring
- [x] iOS L2 CoreML load + OTA (`CoreMLClassifier`, `ModelUpdateService`, App Group install)
- [x] OTP protection (`OTPGuard`) + adversarial FPR benchmark gate
- [x] Pro entitlements gating (custom rules, advanced stats, cloud LLM)
- [x] Gateway `QuotaStore` + `FeatureFlag` on classify; admin quota/flags API (gateway + admin service)
- [x] NATS flywheel worker (`services/flywheel`) + feedback trigger publish
- [x] Tier matrix (`docs/deploy/TIER_MATRIX.md`), `values-prod.yaml`, caddy-site in compose
- [x] App Store ship assets (`ASC_CHECKLIST.md`, `SCREENSHOTS.md`, `fastlane/Fastfile`)
