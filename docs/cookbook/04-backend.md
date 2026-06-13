# Go 后端 Gateway

## 对应原始需求
Go 后端 Gateway

## 涉及文件
services/gateway/cmd/server/main.go

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cd services/gateway && go build ./cmd/server
```
**期望输出：** build succeeds

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
