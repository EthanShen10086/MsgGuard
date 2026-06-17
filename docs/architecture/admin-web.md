# Admin Web Architecture

**Last updated:** 2026-06-17

**Status:** In Progress (Phase 4)

## App

`apps/admin/` — Vite + React SPA served statically or via `npm run dev`.

## Pages

| Route | API | Purpose |
|-------|-----|---------|
| Dashboard | `GET /api/v1/admin/metrics/summary` | Event counts, feedback, shadow stats |
| Feedback | `GET /api/v1/feedback` | Flywheel review queue |
| Models | `GET /api/v1/models/latest` | Current model metadata |
| Flags | `GET/POST /api/v1/admin/flags` | Feature flags |
| Quota | `GET/POST /api/v1/admin/quota/whitelist` | LLM quota whitelist |

## Auth Flow

1. Admin enters token (or app fetches via bootstrap in dev only)
2. Token stored in `sessionStorage`
3. All requests: `Authorization: Bearer <token>`

Production: disable bootstrap; issue tokens via secure ops channel. Optional mTLS at ingress.

## Config

| Env | Default | Description |
|-----|---------|-------------|
| `VITE_GATEWAY_URL` | `http://localhost:8080` | Gateway base URL |

## Build

```bash
cd apps/admin && npm install && npm run build
# dist/ can be served by Caddy or any static host
```

## Future

- OIDC / SSO instead of manual Bearer paste
- Embedded in Helm chart as ConfigMap + sidecar or separate ingress
- Read-only role for support staff
