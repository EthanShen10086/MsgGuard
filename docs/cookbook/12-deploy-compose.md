# 单机 Docker 部署

## 对应原始需求
单机 Docker 部署

## 涉及文件
deploy/docker-compose.yml

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
./deploy/tiers/tier1-compose.sh
```
**期望输出：** curl localhost:8080/health -> ok

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
