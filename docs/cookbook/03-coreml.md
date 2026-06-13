# Core ML L2

## 对应原始需求
Core ML L2

## 涉及文件
ml/export/export_coreml.py

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cd ml && make export
```
**期望输出：** coreml_export.json created

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
