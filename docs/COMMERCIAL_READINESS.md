# Commercial Readiness

Honest status: **Implemented** / **In Progress** / **Planned**

## Startup (Tier 0–1)

| Item | Status |
|------|--------|
| Seed dataset + benchmark gate in CI | Implemented |
| docker-compose single VPS | Implemented |
| App Store metadata draft | Implemented |
| Privacy policy + static site | Implemented |
| Support page | Implemented |
| Pricing page | In Progress |
| Status page | In Progress |

## Growth (Tier 2–3)

| Item | Status |
|------|--------|
| PostgreSQL + Redis wired | Implemented |
| Helm staging environment | Implemented |
| Shadow mode metrics endpoint | Implemented |
| Grafana dashboards | Implemented |
| Data flywheel cron script | Implemented |
| Full-stack Helm (gateway + deps + site) | In Progress |

## Enterprise (Tier 4)

| Item | Status |
|------|--------|
| Multi-replica HPA | Implemented |
| mTLS (Caddy + gateway admin) | Implemented |
| Cross-region backup script | Implemented |
| Audit logging on feedback | Implemented |
| A/B via feature flags | Implemented |
| Runbooks (LLM/DB/model rollback) | Implemented |
| Auth/RBAC + admin API | Implemented |
| Product metrics flywheel | Implemented |
| OIDC / SSO for admin | Implemented (configure `OIDC_*` at deploy) |
| WAF / DDoS edge | Documented + Helm ingress rate-limit (`docs/security/EDGE_WAF.md`) |

## Architecture (Phase 2)

| Item | Status |
|------|--------|
| `pkg/app/container.go` DI | Implemented |
| Store drivers (memory/postgres/mongodb) | Implemented |
| `pkg/httpauth` shared RBAC | Implemented |
| iOS EntitlementProviding | Implemented |
| `scripts/verify.sh` acceptance | Implemented |

## Phase 2.5 — Product Wiring

| Item | Status |
|------|--------|
| iOS L2 CoreML OTA | Implemented |
| OTP protection + FPR gate | Implemented |
| Pro entitlements gating | Implemented |
| Quota + flags admin API | Implemented |
| NATS flywheel worker | Implemented |
| Tier matrix + values-prod | Implemented |
| App Store ship assets | Implemented |

## Phase 2.6 — Long-term

| Item | Status |
|------|--------|
| Shadow Prometheus + Grafana | Implemented |
| Rules CDN ETag/304 | Implemented |
| Multi-locale model publish | Implemented |
| iCloud autoSync + CallDirectory | Implemented |
| cert-manager Helm templates | Implemented |

## Phase 4 — Web Admin & Site

| Item | Status |
|------|--------|
| Admin SPA (Dashboard, Feedback, Models, Flags, Quota) | Implemented |
| pricing.html | Implemented |
| Expanded support links | Implemented |

## Phase 7 — Commercial platform (2026-06)

| Item | Status |
|------|--------|
| CI: `staging` / `release` / `ios-beta` / `security` workflows | Implemented |
| Production AUTH_SECRET startup gate | Implemented |
| Redis token revocation | Implemented (requires `REDIS_URL`) |
| App Attest middleware stub | Implemented |
| App Store Server API client + webhook | Implemented (Apple creds at deploy) |
| S3 model registry overlay | Implemented (requires `MODEL_S3_BUCKET`) |
| Flywheel debounce (`FLYWHEEL_MIN_*`) | Implemented |
| OpenAPI entitlements/auth/webhooks | Implemented |
| Helm prod security defaults | Implemented |

**Still requires deploy-time secrets (not in repo):** `PROD_AUTH_SECRET`, `APPLE_*`, `KUBE_CONFIG_PROD`, ASC for TestFlight, real S3/AWS SDK upload wiring.

## Phase 5 — SRE & Privacy

| Item | Status |
|------|--------|
| SLO document | Implemented |
| Privacy DELETE endpoint | Implemented |
| Prometheus alert rules template | Implemented |

## Phase 6 — Analytics v2

| Item | Status |
|------|--------|
| Event taxonomy v2 schema doc | In Progress |
| Subscription funnel in aggregate_metrics | In Progress |

## Switching Guide

Only change `deploy/config.*.yaml` and run matching `deploy/tiers/*.sh` — no business code changes.
