# AI 后端客户端协同

## 对应原始需求
AI 后端客户端协同

## 涉及文件
docs/api/openapi.yaml

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
curl localhost:8083/api/v1/models/latest
```
**期望输出：** version JSON

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
