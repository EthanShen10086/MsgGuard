# 可观测 Debug

## 对应原始需求
可观测 Debug

## 涉及文件
deploy/prometheus/prometheus.yml

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
open http://localhost:16686
```
**期望输出：** Jaeger UI

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
