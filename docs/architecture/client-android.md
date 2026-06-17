# Android Client Architecture (Placeholder)

**Last updated:** 2026-06-17

**Status:** Planned — no `apps/android/` scaffold yet.

## Intended Scope

- SMS filtering via AndroidId / default SMS app constraints (TBD)
- Shared rule bundle format from gateway rules CDN
- On-device TFLite or ONNX model (parity with iOS L1/L2)
- Analytics + feedback via same `/api/v1/*` contracts

## Dependencies

- Backend platform stable (gateway auth, rules CDN, model registry)
- Legal/privacy copy aligned with iOS (`docs/legal/PRIVACY.md`)

## Open Questions

- Play Store policy for SMS access vs. carrier partnerships
- Whether to ship Android-only free tier vs. cross-platform Pro
