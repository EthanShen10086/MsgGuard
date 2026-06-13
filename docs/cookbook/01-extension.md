# Message Filter Extension

## 对应原始需求
Message Filter Extension

## 涉及文件
apps/ios/MessageFilterExtension/MessageFilterExtension.swift

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
cd apps/ios && xcodebuild -scheme MsgGuard-iOS -destination 'platform=iOS Simulator,name=iPhone 16' build
```
**期望输出：** BUILD SUCCEEDED

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
