# 贝叶斯 L1

## 对应原始需求
贝叶斯 L1

## 涉及文件
ml/train/train_bayes.py

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cd ml && make train
```
**期望输出：** f1 > 0.5 in output

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
