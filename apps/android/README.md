# MsgGuard Android (SMS Filter)

Android does not expose an iOS-style Message Filter extension. This module sketches the closest production path:

1. **Default SMS app role** (API 29+) — `RoleManager.ROLE_SMS` lets MsgGuard receive `SMS_RECEIVED` and classify before notifying the user.
2. **On-device heuristics + Bayes** — mirror iOS `FilterEngine` locally for offline latency.
3. **Gateway defer** — unknown messages call `POST /api/v1/classify/defer` with a device token from `POST /api/v1/auth/device`.

## Components

| File | Purpose |
|------|---------|
| `FilterService.kt` | Stub `Service` that would run classification on inbound SMS |
| `network/GatewayClient.kt` | Device auth + classify API client |

## Setup

1. Set `MSGGUARD_API_BASE` in `local.properties` or build config (default `http://10.0.2.2:8080` for emulator).
2. Request SMS role in the main activity (not included in this skeleton).
3. Register `FilterService` in `AndroidManifest.xml` when wiring receivers.

## Privacy

Message bodies are sent to your gateway only when local rules are inconclusive — match iOS defer behavior and disclose in Play Store data safety form.

## Enterprise

Work-profile deployments can inject `tenant_id` on analytics/feedback payloads; see `docs/enterprise/OVERVIEW.md`.
