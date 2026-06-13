# 分级部署切换

## 对应原始需求
分级部署切换

## 涉及文件
deploy/tiers/

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
ls deploy/tiers/
```
**期望输出：** tier0-4 scripts

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
