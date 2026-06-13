# Commercial Readiness

## Startup (Tier 0–1)
- [ ] Seed dataset + benchmark gate in CI
- [ ] docker-compose single VPS
- [ ] App Store metadata draft
- [ ] Privacy policy (no raw SMS upload default)

## Growth (Tier 2–3)
- [ ] PostgreSQL + Redis wired
- [ ] Helm staging environment
- [ ] Shadow mode metrics
- [ ] Grafana dashboards
- [ ] Data flywheel cron

## Enterprise (Tier 4)
- [ ] Multi-replica HPA
- [ ] mTLS (optional)
- [ ] Cross-region backup (`deploy/ops/backup.sh`)
- [ ] Audit logging
- [ ] A/B experiment via featureflag
- [ ] Runbooks for LLM/DB/model rollback

## Switching Guide
Only change `deploy/config.*.yaml` and run matching `deploy/tiers/*.sh` — no business code changes.
