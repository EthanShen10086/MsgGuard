# App Store Review Notes

## Demo Account (if backend review needed)

Not required — core filtering works fully offline. Optional API token for analytics QA:

```
POST https://api.msgguard.app/api/v1/auth/token
Body: {"roles":["admin"]}
```

## Message Filter Extension

MsgGuard uses Apple's `ILMessageFilterExtension`. All classification runs on-device via:
1. Keyword/heuristic rules (user-editable)
2. On-device Naive Bayes model (trained locally from user samples)
3. Optional Core ML model downloaded from our CDN (metadata only in network calls)

**No SMS body is sent to our servers unless the user enables "Cloud AI" (Pro, opt-in).**

## How to Enable Extension

Settings → Messages → Unknown & Spam → MsgGuard → Enable

## In-App Purchase Testing

- StoreKit Configuration file included: `apps/ios/App/Shared/Store/Products.storekit`
- Sandbox tester recommended for subscription review

## Contact

support@msgguard.app
