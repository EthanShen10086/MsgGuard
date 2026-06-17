# Architecture Decision Records

MsgGuard 重大架构决策记录。格式参考 [ciphera architecture-decisions](https://github.com/EthanShen10086/ciphera)。

## ADR-001: On-device vs Cloud LLM Defer

**Status:** Accepted  
**Date:** 2026-06-01

### Context

短信分类需在扩展进程内低延迟完成（P99 &lt; 15ms），同时支持云端 LLM 提升准确率。用户隐私与离线场景要求默认不将正文上传。

### Decision

采用 **L1 设备 Bayes/CoreML → L2 启发式 → L3 云端 LLM（可选、Pro、feature flag）** 分层。Gateway `/api/v1/classify/defer` 在 LLM 不可用时返回启发式结果而非 5xx。

### Consequences

- iOS/macOS 扩展不依赖网络即可过滤大部分垃圾短信。
- 云端路径受 quota、熔断器、shadow mode 约束，便于 SRE 观测与回滚。
- 需维护 per-locale 模型 OTA 与 benchmark gate。

---

## ADR-002: OIDC Admin SSO vs Device Bearer Token

**Status:** Accepted  
**Date:** 2026-06-10

### Context

Admin SPA 需企业 SSO；iOS/Android 设备 API 需无交互长期凭证。单一 JWT 方案无法同时满足人机与机机场景。

### Decision

- **Admin：** OIDC Authorization Code（`pkg/httpauth/oidc_provider.go`），邮箱白名单或域名后缀；签发短期 admin JWT。
- **Device / API：** HMAC Bearer token（`pkg/adapters/memory/auth.go` 或 Redis 吊销），RBAC 角色 `device`、`ml_engineer`、`admin`。
- 生产环境 `AUTH_SECRET` ≥ 32 字符；OIDC 无 allowlist 时拒绝登录。

### Consequences

- Admin 与设备凭证生命周期分离，降低 token 泄露影响面。
- Helm `values-prod.yaml` 需配置 OIDC 与 `ADMIN_EMAILS`。
- 文档见 `docs/cookbook/27-auth-rbac.md`。

---

## ADR-003: Model OTA + Canary Rollout

**Status:** Accepted  
**Date:** 2026-06-12

### Context

ML 模型需频繁迭代；全量推送可能导致客户端 FPR 回归。需与 release 流水线联动。

### Decision

1. Model 服务注册版本 + checksum；iOS 按 locale 拉取 `.mlmodel` + featurizer JSON。
2. S3 overlay（可选）存储灰度包；`ml/scripts/canary_rollout.py` 在 staging 按百分比切换。
3. Release workflow 在 ML benchmark gate 通过后才构建镜像；Prometheus `shadow_*` 指标监控 disagree 率。

### Consequences

- 模型回滚见 `docs/runbooks/model-rollback.md`。
- Canary 依赖 `STAGING_API_URL` / `RELEASE_ADMIN_TOKEN` secrets，未配置时跳过。
- Shadow disagree &gt; 15% 触发告警（`deploy/prometheus/alerts.yml`）。
