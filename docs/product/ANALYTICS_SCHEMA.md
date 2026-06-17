# Analytics Event Taxonomy v2

**Last updated:** 2026-06-17  
**Endpoint:** `POST /api/v1/analytics`  
**Schema version:** `2.0`

## Envelope

```json
{
  "name": "event_name",
  "device_id": "uuid-v4",
  "props": {
    "schema_version": "2.0",
    "platform": "ios",
    "app_version": "1.2.0",
    "locale": "zh-Hans"
  }
}
```

Server adds: `id`, `trace_id` (from `X-Request-ID`), `timestamp`.

## Core Events

| Event | When | Required props |
|-------|------|----------------|
| `app_launched` | Cold start | `platform`, `app_version` |
| `onboarding_step` | Each onboarding screen | `step`, `step_index` |
| `onboarding_completed` | Onboarding finished | `duration_sec` |
| `filter_completed` | Extension classified message | `layer`, `action`, `latency_ms` |
| `feedback_submitted` | User submits feedback | `category`, `has_trace_id` |
| `purchase_started` | Paywall shown | `product_id`, `source` |
| `purchase_completed` | StoreKit success | `product_id`, `is_trial` |
| `purchase_restored` | Restore purchases | `product_id` |
| `subscription_funnel` | Funnel step | `step` (see below) |
| `cloud_llm_toggled` | Settings change | `enabled` |
| `model_ota_applied` | OTA success | `locale`, `model_version` |
| `rules_sync` | Rule sync result | `status`, `etag` |
| `error` | Client error | `code`, `domain` |

## Subscription Funnel Steps (`subscription_funnel`)

| step | Description |
|------|-------------|
| `paywall_view` | Subscription screen viewed |
| `product_selected` | User tapped a product |
| `purchase_attempt` | StoreKit purchase initiated |
| `purchase_success` | Transaction verified |
| `purchase_cancel` | User cancelled |
| `purchase_failed` | StoreKit error |
| `restore_attempt` | Restore tapped |
| `restore_success` | Entitlement restored |

Aggregate in `ml/product/aggregate_metrics.py` → `subscription_funnel` block (Phase 6).

## v1 → v2 Migration

| v1 name | v2 name | Notes |
|---------|---------|-------|
| `app_launched` | `app_launched` | add `schema_version` |
| `onboarding_completed` | `onboarding_completed` | unchanged |
| `filter_completed` | `filter_completed` | add `layer` |
| `purchase_completed` | `purchase_completed` | add `product_id` |
| — | `subscription_funnel` | new funnel events |

Legacy events without `schema_version` treated as v1 in aggregations.

## PII Policy

- No SMS body, phone numbers, or contact names in `props`
- `device_id` is random UUID — used for privacy DELETE only
- Optional `user_id` hash (planned) — never raw Apple ID

## Admin / ML Usage

- Dashboard counts via `CountByName` on gateway analytics store
- Weekly report pulls admin summary + funnel placeholders
- See [METRICS.md](METRICS.md) for north-star definitions
