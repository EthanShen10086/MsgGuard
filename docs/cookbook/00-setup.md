# 环境搭建

## 对应原始需求
环境搭建

## 涉及文件
apps/ios/setup.sh, ml/requirements.txt

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cd apps/ios && bash setup.sh
cd ml && pip install -r requirements.txt
```
**期望输出：** setup completes without error

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
