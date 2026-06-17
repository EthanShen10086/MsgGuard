# Release & Deployment Guide

## Environments

| Environment | Trigger | Workflow | Required secrets |
|-------------|---------|----------|------------------|
| CI | PR / push `main` | `ci.yml` | — |
| Staging | push `staging` | `deploy-staging.yml` | `KUBE_CONFIG`, `STAGING_AUTH_SECRET`, optional `STAGING_API_URL` |
| Production | tag `v*.*.*` | `release.yml` | `KUBE_CONFIG_PROD`, `PROD_AUTH_SECRET`, optional `RELEASE_ADMIN_TOKEN` |
| TestFlight | manual | `ios-beta.yml` | `APP_STORE_CONNECT_API_KEY_*`, `FASTLANE_*` |

## Production release (backend)

1. Merge `staging` → `main` after QA sign-off.
2. Tag: `git tag v1.2.3 && git push origin v1.2.3`
3. `release.yml` runs: ML benchmark gate → Docker push to `ghcr.io` → Helm validate → optional canary (5%) → prod deploy.

## Canary model rollout

After staging deploy, or from release workflow:

```bash
export GATEWAY_URL=https://staging-api.msgguard.app
export ADMIN_TOKEN=<ml_engineer or admin token>
python ml/scripts/canary_rollout.py --gateway "$GATEWAY_URL" --token "$ADMIN_TOKEN" --percentage 10
```

Promote to 100% via Admin → Flags or `canary_rollout.py --percentage 100`.

Rollback: `POST /api/v1/admin/models/rollback` or runbook `docs/runbooks/model-rollback.md`.

## iOS / macOS TestFlight

GitHub Actions → **iOS Beta (TestFlight)** → choose `ios` or `macos-mail`.

Requires App Store Connect API key in environment `testflight`.

## Production configuration (ops-only)

Set in Helm `values-prod.yaml` or secrets — **not in git**:

| Variable | Value |
|----------|-------|
| `MSGGUARD_ENV` | `production` |
| `AUTH_SECRET` | ≥32 random chars |
| `AUTH_BOOTSTRAP_ENABLED` | `false` |
| `MODEL_DOWNLOAD_AUTH_REQUIRED` | `true` |
| `APPLE_ISSUER_ID` / `APPLE_KEY_ID` / `APPLE_PRIVATE_KEY` | App Store Server API |
| `GOOGLE_SAFE_BROWSING_API_KEY` | Optional threat intel |

## Branch protection (GitHub settings)

Recommended for `main` and `staging`:

- Require PR before merge
- Require status checks: `go-gateway`, `ml-benchmark`, `swift-packages`
- Require 1 approval for `main`
