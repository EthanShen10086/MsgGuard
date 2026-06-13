# 数据流水线总览

## 对应原始需求
数据流水线总览

## 涉及文件
ml/Makefile

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cd ml && make help
```
**期望输出：** targets listed

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
