# macOS Mail Extension

MsgGuard filters **macOS Mail.app** messages using the same `HybridFilterEngine` as the iOS SMS extension.

## Platform note

iOS does **not** expose a system Mail filter API. This target is **macOS only** (MailKit `MEExtension`).

## Setup

```bash
cd apps/macos
xcodegen generate
xcodebuild -scheme MailExtension -destination 'platform=macOS' build
```

Enable in **Mail → Settings → Extensions → MsgGuard Mail**.

## Behavior

- Reads `FilterConfig` + models from App Group `group.com.ethanshen.msgguard`
- Classifies `subject + preview` via L0/L1/L2
- Spam-like messages → **Move to Trash** (`MEAction.moveToTrash`)

## Requirements

- macOS 14+
- Mail extension capability (Apple Developer portal)
- Shared App Group with iOS app for rules/model sync
