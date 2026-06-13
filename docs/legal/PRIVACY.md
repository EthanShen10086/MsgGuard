# MsgGuard Privacy Policy

**Effective date:** 2026-06-13  
**Contact:** privacy@msgguard.app

---

## 中文摘要

MsgGuard（「短信卫士」）默认在您的 iPhone **本地**识别垃圾短信，**不上传短信原文**。仅在您订阅 Pro 并主动开启「云端 AI」时，才会将单条短信文本发送至我们的 API 进行分类。您可随时在设置中关闭。

---

## English

### 1. Who we are

MsgGuard ("we") provides an on-device SMS filtering app for iOS with an optional cloud AI tier for Pro subscribers.

### 2. What we collect

| Data | Default | When collected | Purpose |
|------|---------|----------------|---------|
| SMS message body | **Not collected** | Only if Pro + Cloud AI enabled | Spam/ham classification |
| Filter result (spam/ham) | On-device only | Always local | Block or allow message |
| User feedback label | Hash + label optional | When you submit feedback | Improve models |
| Analytics events | Event name, trace ID | App usage | Product improvement |
| Crash reports | Opt-in | If enabled | Stability |
| Device identifier | Random UUID | First launch | Analytics dedup |
| Purchase status | Via Apple StoreKit | Subscription | Entitlement |

We do **not** sell personal data.

### 3. On-device processing (default)

By default, all filtering uses rules, Bayes, and Core ML models stored in the App Group container. No message content leaves your device.

### 4. Cloud AI (opt-in, Pro)

When you enable **Cloud AI** in Settings:
- Single message text is sent over HTTPS to `api.msgguard.app`
- Requests include a trace ID for support; logs are PII-redacted
- You can disable this at any time

### 5. Data retention

- Feedback: 90 days rolling (server)
- Analytics: aggregated, 30–90 days
- Crash reports: 30 days (if opted in)
- Local models/rules: until you delete the app

### 6. Your rights

Depending on jurisdiction you may request access, correction, or deletion of server-side data tied to your device ID or trace ID. Email privacy@msgguard.app.

### 7. Children

MsgGuard is not directed at children under 13.

### 8. Changes

We will post updates at https://msgguard.app/privacy and update the effective date.

---

## 中文完整版

### 1. 适用范围

本政策适用于 MsgGuard iOS 应用及可选后端 API（`api.msgguard.app`）。

### 2. 本地处理（默认）

- 短信内容在设备上完成分类，不上传原文。
- 规则、贝叶斯模型、Core ML 模型存储于 App Group，仅本机与扩展共享。

### 3. 云端 AI（可选，Pro）

开启后单条短信经 HTTPS 发送至服务器仅用于分类；可随时关闭。

### 4. 反馈与分析

- 反馈默认仅存标签与哈希，不含全文（除非您明确提交样本文本）。
- 分析事件不含短信正文。

### 5. 联系我们

privacy@msgguard.app
