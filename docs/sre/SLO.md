# MsgGuard SLO Definitions

**Last updated:** 2026-06-17

## Service Level Objectives

| Service | SLI | SLO | Window |
|---------|-----|-----|--------|
| Gateway availability | Successful `/health` + non-5xx on classify/feedback | 99.5% | 30d |
| Gateway latency (classify defer) | p99 < 2s (excludes LLM upstream) | 95% of hours | 7d |
| Gateway latency (LLM defer) | p99 < 8s end-to-end | 90% of hours | 7d |
| Rules CDN | `/api/v1/rules/latest` 200 + ETag | 99.9% | 30d |
| Model CDN | `/api/v1/models/latest` 200 | 99.5% | 30d |
| iOS extension | Filter P99 < 15ms on-device | 99% sessions | 7d (client) |

## Error Budget

- **Availability:** 0.5% monthly downtime ≈ 3.6 hours/month for gateway
- Burn alerts when >50% budget consumed in 7 days (see `deploy/prometheus/alerts.yml`)

## Measurement

| SLI | Source |
|-----|--------|
| HTTP availability | Prometheus `up{job="msgguard-gateway"}` |
| Error rate | `rate(http_requests_total{status=~"5.."}[5m])` |
| Latency | voxera-kit metrics histogram on gateway |
| Shadow disagree rate | `/metrics/shadow` |
| Extension latency | iOS OSSignpost / Analytics |

## Non-Goals (v1)

- Multi-region active-active
- 99.99% for LLM upstream (provider SLA excluded)
- Android client SLOs (not shipped)

## Runbooks

| Alert | Runbook |
|-------|---------|
| GatewayDown | Restart compose / `kubectl rollout restart`; check DB DSN |
| HighErrorRate | Check LLM circuit breaker, DB connectivity, recent deploy |
| HighShadowDisagreeRate | Review shadow stats admin endpoint; pause cloud_llm flag |

## Review

Revisit SLOs when Phase 7 production auth + Phase 12 mature product targets are active.
