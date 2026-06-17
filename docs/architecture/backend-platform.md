# Backend Platform Architecture

**Last updated:** 2026-06-17

## Services

| Service | Port | Role |
|---------|------|------|
| gateway | 8080 | API entry, middleware, classify/defer, admin, feedback, analytics |
| rules | 8081 | Rule bundle distribution |
| classify | 8082 | Stub (logic in gateway) |
| model | 8083 | Core ML / Bayes metadata CDN |
| feedback | 8084 | Stub (logic in gateway) |
| admin | 8085 | Standalone quota/flags (optional; duplicated on gateway) |
| flywheel | — | NATS worker for retrain triggers |

## DI Root

`pkg/app/container.go` wires stores, auth, quota, flags, LLM, queue from `deploy/config.*.yaml`.

## Store Drivers

| Driver | Package |
|--------|---------|
| memory | `pkg/adapters/memory` |
| postgres | `pkg/adapters/postgres` |
| mongodb | `pkg/adapters/mongodb` |

Switch via `DATABASE_DRIVER` + tier script — no business code changes.

## Middleware Stack (Gateway)

Recovery → RequestID → Logging → Metrics → SecurityHeaders → LoadShed → RateLimit → Timeout → PIIRedact → (Tracing) → (mTLS admin)

## API Surface

See `docs/api/openapi.yaml`. Admin paths under `/api/v1/admin/*` require Bearer + `admin` role; production may require mTLS.

## Deploy

- Tier 0–1: docker-compose
- Tier 3–4: Helm (`deploy/helm/msgguard`)
- Static site: Caddy `:8088` or external CDN
