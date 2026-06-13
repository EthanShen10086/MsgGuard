# App Store Metadata Draft — MsgGuard 1.0

> 复制到 App Store Connect。产品 ID 与 `Products.storekit` / `StoreManager` 一致。

## App Information

| Field | zh-Hans | en-US |
|-------|---------|-------|
| **Name** | MsgGuard 短信卫士 | MsgGuard SMS Filter |
| **Subtitle** | 本地 AI 垃圾短信拦截 | On-device AI spam SMS filter |
| **Primary Category** | Utilities | Utilities |
| **Secondary Category** | Productivity | Productivity |
| **Content Rights** | Does not contain third-party content | — |

## Description

### zh-Hans（≤4000 字符）

MsgGuard 是一款注重隐私的混合 AI 垃圾短信过滤器，对标「熊猫吃短信」体验，默认在设备本地完成识别，不上传短信原文。

**核心功能**
- Message Filter Extension：系统级来电/短信过滤，骚扰信息直接进垃圾箱
- 三层混合引擎：关键词规则 → 朴素贝叶斯 → Core ML，可离线工作
- 用户反馈闭环：误拦/漏拦一键标注，帮助本地模型持续改进
- Pro 订阅：云端 AI 二次确认（可选）、高级规则、详细统计

**隐私承诺**
- 默认不上传短信内容；云端 AI 需 Pro 且用户主动开启
- 分析事件仅含事件名与 TraceID，不含消息正文
- 详见 https://msgguard.app/privacy

### en-US

MsgGuard is a privacy-first hybrid AI spam SMS filter. Classification runs on your device by default—message bodies are never uploaded unless you opt in to Cloud AI (Pro).

**Features**
- iOS Message Filter Extension
- Hybrid pipeline: rules → Bayes → Core ML, works offline
- Feedback loop for false positives/negatives
- Pro: optional cloud AI, advanced rules, detailed stats

Privacy: https://msgguard.app/privacy

## Keywords

| Locale | Keywords (100 chars max) |
|--------|--------------------------|
| zh-Hans | 垃圾短信,短信拦截,骚扰短信,短信过滤,熊猫吃短信,防骚扰,短信卫士 |
| en-US | spam SMS,SMS filter,junk messages,text blocker,message filter |

## URLs

| Field | URL |
|-------|-----|
| Privacy Policy | https://msgguard.app/privacy |
| Support | https://msgguard.app/support |
| Marketing | https://msgguard.app |

## In-App Purchases

| Product ID | Type | Reference Name |
|------------|------|----------------|
| `com.ethanshen.msgguard.pro.monthly` | Auto-renewable | Pro Monthly |
| `com.ethanshen.msgguard.pro.yearly` | Auto-renewable | Pro Yearly |

**Subscription Group:** MsgGuard Pro (`MG000001-0001`)

## Privacy Nutrition Labels (App Store Connect)

| Data Type | Linked to User | Used for Tracking | Purpose |
|-----------|----------------|-------------------|---------|
| Product Interaction | No | No | Analytics (opt-in events) |
| Crash Data | No | No | Diagnostics (opt-in) |
| User ID (device UUID) | No | No | Analytics correlation |

**Not collected by default:** SMS content, contacts, location, photos.

## Screenshots Checklist

- [ ] 6.7" — Dashboard + filter stats
- [ ] 6.7" — Settings + Extension toggle
- [ ] 6.7" — Subscription / Pro features
- [ ] 6.1" — same set
- [ ] iPad (if supported later)

## Review Notes

See [review-notes.md](./review-notes.md).

## Export Compliance

- Uses encryption: **Yes** (HTTPS only) → ERN exempt standard encryption
- Message Filter Extension: explain on-device classification in Review Notes
