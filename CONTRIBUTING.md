# Contributing to MsgGuard

## Branch strategy

| Branch | Purpose | Merges to | Deploy |
|--------|---------|-----------|--------|
| `feature/*` | New work | `staging` via PR | — |
| `staging` | Pre-production integration | `main` after QA | Auto → staging K8s |
| `main` | Release-ready | — | Tag `v*.*.*` → production |
| `hotfix/*` | Urgent fixes | `main` + `staging` | Tag patch release |

**Do not push features directly to `main`.** Open a PR; CI must pass.

## PR checklist

- [ ] `go build` / `swift test` / relevant tests pass locally
- [ ] OpenAPI updated if API changed (`docs/api/openapi.yaml`)
- [ ] Architecture doc updated if boundaries changed (`docs/architecture/`)
- [ ] No secrets in code; use GitHub Environments for deploy
- [ ] `COMMERCIAL_READINESS.md` status updated if shipping surface changed

## Local development

```bash
make verify          # acceptance script
cd services/gateway && AUTH_SECRET=dev-secret-goes-here-32chars-min go run ./cmd/server
cd apps/admin && npm run dev
```

## Release process

See [docs/RELEASE.md](docs/RELEASE.md).

## Environments (GitHub)

Configure repository **Environments**: `staging`, `production`, `testflight` with secrets documented in `docs/RELEASE.md`.
