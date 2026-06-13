# GPU 训练

## 对应原始需求
GPU 训练

## 涉及文件
deploy/k8s/gpu-training-job.yaml

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
kubectl apply -f deploy/k8s/gpu-training-job.yaml --dry-run=client
```
**期望输出：** job valid

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
