# App Store 上架

## 对应原始需求
App Store 上架

## 涉及文件
- docs/app-store/metadata.md
- docs/app-store/review-notes.md
- docs/legal/PRIVACY.md
- deploy/site/privacy/index.html

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
test -f docs/app-store/metadata.md && test -f docs/legal/PRIVACY.md
open deploy/site/privacy/index.html
```
**期望输出：** metadata + privacy files present

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
- [Cookbook 32 — Helm / Privacy](./32-helm-privacy-appstore.md)
