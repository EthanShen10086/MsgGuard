# 配置切换

## 对应原始需求
配置切换

## 涉及文件
deploy/config.yaml

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cat deploy/config.yaml
```
**期望输出：** database/cache sections present

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
