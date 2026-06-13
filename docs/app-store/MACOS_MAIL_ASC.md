# macOS Mail Extension — App Store Connect Checklist

## Apple Developer Portal
- [ ] Enable **Mail Extension** capability for App ID `com.ethanshen.msgguard.mailhost.filter`
- [ ] Create App ID `com.ethanshen.msgguard.mailhost` (host app embeds extension)
- [ ] Enable **App Groups** (`group.com.ethanshen.msgguard`) on both IDs
- [ ] Create provisioning profiles:
  - Mac App Store: host + extension
  - Developer ID (optional notarized direct): host + extension

## App Store Connect
- [ ] Create macOS app record `com.ethanshen.msgguard.mailhost`
- [ ] Category: Utilities
- [ ] Privacy policy URL: https://msgguard.app/privacy
- [ ] Review notes: extension uses on-device HybridFilterEngine; may move spam to Trash

## Build & Upload
```bash
cd apps/macos && xcodegen generate
# TestFlight (Mac App Store)
bundle exec fastlane mac mail_testflight
# Or Developer ID + notarization
bundle exec fastlane mac mail_notarize
```

## Environment variables (CI / local)
| Variable | Purpose |
|----------|---------|
| `FASTLANE_APPLE_ID` | Apple ID email |
| `FASTLANE_TEAM_ID` | Team ID |
| `FASTLANE_APPLE_APPLICATION_SPECIFIC_PASSWORD` | Notary / upload |
| `MACOS_HOST_PROFILE` | Host provisioning profile name |
| `MACOS_MAIL_EXT_PROFILE` | Extension provisioning profile name |

## User enablement
Mail → Settings → Extensions → enable **MsgGuard Mail**

See also [MAIL_EXTENSION.md](../product/MAIL_EXTENSION.md).
