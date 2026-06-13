# App Store 上架

## 对应原始需求
App Store 上架

## 涉及文件
docs/app-store/metadata.md

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
open docs/app-store/metadata.md
```
**期望输出：** metadata checklist present

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
