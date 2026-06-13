# 日常维护

## 对应原始需求
日常维护

## 涉及文件
ml/flywheel/schedule_retrain.sh

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
bash ml/flywheel/schedule_retrain.sh
```
**期望输出：** retrain pipeline runs

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
