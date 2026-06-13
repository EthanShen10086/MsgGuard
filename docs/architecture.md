# MsgGuard Architecture

## Overview

MsgGuard is a hybrid AI spam SMS filter with three layers:

1. **iOS Platform App** — App Store product, onboarding, rules, feedback, subscription
2. **Message Filter Extension** — On-device L0-L2 classification (<15ms)
3. **Go Backend** — Rules/model CDN, optional cloud LLM (defer), feedback with TraceID

## iOS Targets

| Target | Role |
|--------|------|
| MsgGuard-iOS | Main SwiftUI app (platform carrier) |
| MessageFilterExtension | IdentityLookup filter engine |
| MsgGuardWidget | Today blocked count widget |

## SPM Packages

- `SharedModels` — Types, AppConstants, MGLogger
- `DesignSystem` — Theme, UserMode, components
- `FilterEngine` — L0 heuristics, L1 Bayes, L2 Core ML stub
- `BlocklistStore` — App Group persistence actor

## Backend Services

| Service | Port | Role |
|---------|------|------|
| gateway | 8080 | API entry, middleware stack, classify defer |
| rules | 8081 | Rule bundle distribution |
| classify | 8082 | Stub (logic in gateway) |
| model | 8083 | Core ML metadata CDN |
| feedback | 8084 | Stub (logic in gateway) |

## Data Flow

```
SMS → Extension → FilterEngine (L0→L1→L2) → Messages.app
                      ↓ opt-in
                 iOS defer → Gateway → LLM
Main App → BlocklistStore → App Group → Extension reads rules
```

## Observability

- iOS: MGLogger, AnalyticsManager, PerformanceMonitor (OSSignpost)
- Backend: RequestID → OpenTelemetry → Jaeger, Prometheus /metrics
