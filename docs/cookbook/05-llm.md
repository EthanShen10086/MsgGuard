# 云端 LLM L3

## 对应原始需求
云端 LLM L3

## 涉及文件
services/gateway/internal/handler/classify.go

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
curl -X POST localhost:8080/api/v1/classify -d '{"body":"free gift"}'
```
**期望输出：** JSON action field

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
