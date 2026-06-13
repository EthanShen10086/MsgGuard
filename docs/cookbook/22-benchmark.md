# Benchmark 成本

## 对应原始需求
Benchmark 成本

## 涉及文件
ml/benchmark/run_benchmark.py

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cd ml && make benchmark
```
**期望输出：** gate_passed=True

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
