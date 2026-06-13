# 兜底降级

## 对应原始需求
兜底降级

## 涉及文件
services/gateway/internal/handler/classify.go

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
unset QWEN_API_KEY && curl classify
```
**期望输出：** heuristic fallback

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
