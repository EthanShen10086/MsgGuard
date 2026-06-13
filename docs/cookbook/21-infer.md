# 推理脚本

## 对应原始需求
推理脚本

## 涉及文件
ml/infer/infer_bayes.py

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cd ml && make infer TEXT='免费贷款'
```
**期望输出：** JSON label spam

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
