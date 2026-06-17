# iOS Client Architecture

**Last updated:** 2026-06-17

## Targets

| Target | Role |
|--------|------|
| MsgGuard-iOS | Main SwiftUI app — onboarding, rules, feedback, subscription |
| MessageFilterExtension | IdentityLookup filter engine (L0–L2, &lt;15ms) |
| MsgGuardWidget | Today blocked count widget |
| CallDirectoryExtension | Call directory skeleton (Phase 2.6) |

## SPM Packages

- `SharedModels` — types, AppConstants, MGLogger
- `DesignSystem` — theme, UserMode, Elder Mode
- `FilterEngine` — L0 heuristics, L1 Bayes, L2 Core ML, OTPGuard
- `BlocklistStore` — App Group persistence actor

## Classification Pipeline

```
SMS → Extension → FilterEngine (L0 → L1 → L2) → Messages.app
                      ↓ Pro + opt-in
                 defer → Gateway → LLM
```

## Sync & OTA

- Rules: incremental sync via ETag/304 from gateway rules CDN
- Models: per-locale Bayes + CoreML OTA via model service; SHA256 verify
- iCloud KV: optional auto-sync for blocklist (Pro)

## Key Files

- `apps/ios/Packages/FilterEngine/` — classification layers
- `apps/ios/App/Shared/Networking/` — APIClient, CloudSyncService
- `apps/ios/App/Shared/Analytics/` — AnalyticsManager
