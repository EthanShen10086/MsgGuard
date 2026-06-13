# MsgGuard Product Metrics

## North Star
**Effective Block Rate** — 真垃圾短信被正确拦截 / 用户收到的总垃圾短信（设备端 + 用户确认）

## Guardrail Metrics
| Metric | Target | Source |
|--------|--------|--------|
| False Positive Rate (误杀率) | <= 2% | benchmark + user feedback |
| False Negative Rate (漏杀率) | <= 15% | shadow mode + samples |
| Extension P99 latency | < 15ms | FilterEngine Signpost |
| User complaint rate | < 0.5% | feedback API |

## Business Metrics
| Metric | Description |
|--------|-------------|
| Subscription conversion | onboarding → purchase_completed |
| Sample feeding rate | sample_submitted / MAU |
| LLM opt-in rate | cloudLLMEnabled users / MAU |
| Rule sync success | sync without error |

## Event Taxonomy (iOS → `/api/v1/analytics`)
- `app_launched`, `onboarding_completed`, `filter_completed`, `feedback_submitted`, `purchase_completed`, `error`
