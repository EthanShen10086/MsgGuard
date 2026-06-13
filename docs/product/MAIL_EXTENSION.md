# macOS Mail Extension

MsgGuard filters **macOS Mail.app** messages using the same `HybridFilterEngine` as the iOS SMS extension.

## Platform note

iOS does **not** expose a system Mail filter API. This target is **macOS only** (MailKit `MEExtension`).

## Setup

```bash
cd apps/macos
xcodegen generate
xcodebuild -scheme MsgGuardMailHost -destination 'platform=macOS' build
```

Enable in **Mail → Settings → Extensions → MsgGuard Mail**.

## Classification flow

1. First pass: if `rawData` is unavailable → `MEMessageActionDecision.invokeAgainWithBody`
2. Second pass: parse RFC822 plain text + subject, run `HybridFilterEngine`
3. Spam → `MEMessageAction.moveToTrash`

## Distribution

- **TestFlight / Mac App Store**: `bundle exec fastlane mac mail_testflight`
- **Developer ID + notarization**: `bundle exec fastlane mac mail_notarize`

See [MACOS_MAIL_ASC.md](../app-store/MACOS_MAIL_ASC.md).

## Behavior

- Reads `FilterConfig` + models from App Group `group.com.ethanshen.msgguard`
- Classifies `subject + preview` via L0/L1/L2
- Spam-like messages → **Move to Trash** (`MEAction.moveToTrash`)

## Requirements

- macOS 14+
- Mail extension capability (Apple Developer portal)
- Shared App Group with iOS app for rules/model sync
