# Commercial Architecture

**Last updated:** 2026-06-17

## Tiers

See [TIER_MATRIX.md](../deploy/TIER_MATRIX.md). Switch config + tier script only.

| Phase | Tier | Capability |
|-------|------|------------|
| Startup | 0–1 | Compose, seed ML, static site |
| Growth | 2–3 | Postgres/Mongo, Helm, Grafana, flywheel |
| Enterprise | 4 | mTLS, HPA, backup, audit, A/B flags |

## Monetization

- **Free:** on-device L0–L2, basic stats
- **Pro (StoreKit):** custom rules, advanced stats, cloud LLM opt-in
- Entitlements via `EntitlementProviding` + Keychain persistence

## Static Site (`deploy/site/`)

- `/` — landing
- `/privacy/` — policy
- `/support/` — FAQ + setup
- `/pricing.html` — tier comparison (Phase 4)
- `/status.html` — status stub (Phase 5)

## App Store

- Metadata: `docs/app-store/`
- Fastlane: `apps/ios/fastlane/`
- ASC checklist: `docs/app-store/ASC_CHECKLIST.md`

## Readiness Tracking

Honest status per item: [COMMERCIAL_READINESS.md](../COMMERCIAL_READINESS.md)
