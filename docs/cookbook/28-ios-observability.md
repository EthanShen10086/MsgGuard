# iOS Observability

## 对应原始需求
对齐 ScrollCap：Analytics JSONL、CrashReporter、TraceID、延迟初始化

## 涉及文件
- `apps/ios/App/Shared/Analytics/AnalyticsManager.swift`
- `apps/ios/App/Shared/Analytics/CrashReporter.swift`
- `apps/ios/App/Shared/Networking/NetworkClient.swift`
- `apps/ios/App/Shared/MsgGuardApp.swift`

## 动手验收
```bash
cd apps/ios && xcodegen generate && xcodebuild -scheme MsgGuard-iOS -destination 'platform=iOS Simulator,name=iPhone 16' build
# 启动 App 后检查 App Group:
# analytics.jsonl, crash_reporter.installed 存在
```
**期望输出：** BUILD SUCCEEDED; sentinel files in App Group container

## Debug 指南
- Analytics 不上报 → 检查 gateway `/api/v1/analytics`
- TraceID → APIClient.lastTraceID 写入 analytics props

## 扩展阅读
- ScrollCap REQUIREMENTS.md Analytics 章节
