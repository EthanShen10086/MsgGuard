# Product Iteration SOP

## Weekly Review (every Monday)

1. Run `python ml/product/aggregate_metrics.py`
2. Open Grafana product dashboard
3. Check guardrails in [METRICS.md](./METRICS.md)

## Decision Rules

| Signal | Threshold | Action |
|--------|-----------|--------|
| 误杀率 FPR | > 2% | Add adversarial samples; review keyword_allow rules |
| Shadow 分歧率 | > 15% | Compare LLM vs local; tune confidence threshold |
| 订阅转化率 | < 5% | A/B onboarding copy; review paywall timing |
| Extension P99 | > 15ms | Profile FilterEngine layers; disable L2 if needed |
| feedback 投诉 | spike > 2x | Hotfix rules bundle; pause model rollout |

## Model Release Gate
1. `make benchmark` gate_passed=True
2. Shadow 7-day disagree rate stable
3. Admin sign-off via audit log

## Roadmap Update
After each sprint, update [ROADMAP.md](./ROADMAP.md) data-driven section with metrics snapshot.
