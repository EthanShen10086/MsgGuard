# On-Call Runbook

MsgGuard 生产值班与升级路径。**Last updated:** 2026-06-17

## Ownership

| Area | Primary | Backup |
|------|---------|--------|
| Gateway / Helm | @EthanShen10086 | TBD |
| ML pipeline / flywheel | @EthanShen10086 | TBD |
| iOS / macOS clients | @EthanShen10086 | TBD |

见 [`.github/CODEOWNERS`](../../.github/CODEOWNERS)。

## Escalation

1. **L1 — 告警响应（15 min）**  
   查阅 Prometheus/Grafana → 对应 runbook（`docs/runbooks/`）→ 必要时重启 Pod / 回滚 Helm release。

2. **L2 — 服务负责人（30 min）**  
   若 SLO burn-rate 告警持续 &gt; 1h 或数据丢失风险，联系 Primary owner。

3. **L3 — 安全事件**  
   疑似数据泄露、认证绕过、WAF 绕过 → 按 [`SECURITY.md`](../../SECURITY.md) 报告流程，同时暂停 admin OIDC 发行（rotate `AUTH_SECRET`）。

## Paging (placeholder)

配置以下之一并写入 GitHub Environment `production` secrets：

- **PagerDuty** — `PAGERDUTY_ROUTING_KEY`
- **Opsgenie** — `OPSGENIE_API_KEY`

Alertmanager 路由示例见 `deploy/prometheus/alertmanager.yml`（待与云监控对接）。

## Common alerts

| Alert | First action | Runbook |
|-------|--------------|---------|
| `GatewayDown` | `kubectl rollout restart deploy/msgguard-gateway` | `docs/runbooks/db-failure.md` |
| `HighErrorRate` | 检查 LLM 熔断、DB DSN | `docs/runbooks/llm-outage.md` |
| `SLOErrorBudgetBurnFast` | 暂停 deploy；检查最近变更 | `docs/sre/SLO.md` |
| `ShadowDisagreeRateHigh` | 关闭 `cloud_llm` flag | `docs/runbooks/model-rollback.md` |

## Communication

- **内部：** Slack/飞书 #msgguard-oncall（占位）
- **用户：** 状态页 / `docs/legal/SUPPORT.md` 联系渠道
- **事后：** 48h 内完成 incident 摘要（severity、timeline、action items）

## SLO reference

见 [`docs/sre/SLO.md`](./SLO.md)。错误预算消耗 &gt; 50% / 7d 时冻结非紧急发布。
