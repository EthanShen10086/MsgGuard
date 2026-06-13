# Commercial Readiness

## Startup (Tier 0–1)
- [x] Seed dataset + benchmark gate in CI
- [x] docker-compose single VPS
- [ ] App Store metadata draft
- [ ] Privacy policy (no raw SMS upload default)

## Growth (Tier 2–3)
- [x] PostgreSQL + Redis wired (via Container config)
- [x] Helm staging environment
- [x] Shadow mode metrics endpoint
- [x] Grafana dashboards (msgguard + product)
- [x] Data flywheel cron script

## Enterprise (Tier 4)
- [x] Multi-replica HPA (Helm template)
- [ ] mTLS (optional)
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
