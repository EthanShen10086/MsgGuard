# API Versioning Policy

MsgGuard HTTP API 当前前缀：**`/api/v1`**。

## Principles

1. **URL 版本化** — 破坏性变更通过新前缀（如 `/api/v2`）发布，旧版本保留至少 **6 个月** 弃用期。
2. **OpenAPI 为契约** — [`openapi.yaml`](./openapi.yaml) 与实现同步；CI 未来接入 Spectral lint。
3. **响应头** — 弃用端点返回 `Deprecation: true` 与 `Sunset: <RFC 7231 date>`（Gateway 中间件待实现）。

## Non-breaking changes (same `/api/v1`)

- 新增可选 JSON 字段
- 新增端点
- 放宽校验（接受原先 400 的输入）
- 性能优化、bug 修复

## Breaking changes (require `/api/v2` or major bump)

- 删除或重命名字段/端点
- 更改认证方式
- 更改错误码语义
- 更改分页/排序默认行为

## Deprecation process

| Week | Action |
|------|--------|
| 0 | ADR + CHANGELOG + OpenAPI `deprecated: true` |
| 1 | Admin 公告 + metrics 标注 `api_version=v1` |
| 4 | 文档与客户端 SDK 迁移指南 |
| 26 | 移除 v1 端点（或返回 410） |

## Client guidance

- iOS/macOS：随 App Store 版本绑定 API 能力；OTA 规则/模型与 API 版本独立。
- Admin SPA：与 gateway 同仓库发布，无长期多版本并存需求。
- 外部集成：使用 `User-Agent` + `X-Client-Version` 请求头便于服务端遥测。

## Related

- [`docs/adr/README.md`](../adr/README.md)
- [`docs/RELEASE.md`](../RELEASE.md)
