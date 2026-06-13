# Cookbook 31 — StoreKit 2 + Entitlements

## 架构

```
StoreManager (StoreKit 2)
    ↓ syncFromStoreKit
EntitlementManager (EntitlementProviding)
    ↓ Keychain
AppState.isPro → Settings / Cloud LLM gate
```

## 本地测试

1. Xcode → Scheme → Run → Options → StoreKit Configuration → `Products.storekit`
2. 运行 App，进入 Settings → Upgrade to Pro
3. 使用 StoreKit 测试购买 / Restore Purchases

## 验收

```bash
cd apps/ios && xcodegen generate
xcodebuild -scheme MsgGuard-iOS -destination 'platform=iOS Simulator,name=iPhone 16' build
# BUILD SUCCEEDED

# 启动后 EntitlementManager 从 Transaction.currentEntitlements 恢复
# Settings 中 Cloud LLM 在未订阅时 disabled
```

## App Store Connect

| Product ID | Type |
|------------|------|
| `com.ethanshen.msgguard.pro.monthly` | Auto-renewable |
| `com.ethanshen.msgguard.pro.yearly` | Auto-renewable |

Subscription Group ID: `MG000001-0001`（与 Products.storekit 一致）
