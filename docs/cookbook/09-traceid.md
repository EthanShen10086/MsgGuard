# TraceID 追溯

## 对应原始需求
TraceID 追溯

## 涉及文件
services/gateway/cmd/server/main.go

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
curl -v localhost:8080/health 2>&1 | grep -i x-request
```
**期望输出：** X-Request-ID header

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
