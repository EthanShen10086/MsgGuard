# 可插拔切换

## 对应原始需求
可插拔切换

## 涉及文件
pkg/ports/

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
grep driver deploy/config.yaml
```
**期望输出：** postgres/memory drivers

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
