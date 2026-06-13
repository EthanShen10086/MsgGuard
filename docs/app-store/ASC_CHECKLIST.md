# App Store Connect Checklist

## App Record
- [ ] Create app `com.ethanshen.msgguard` in ASC
- [ ] Upload metadata from [metadata.md](./metadata.md)
- [ ] Privacy URL: https://msgguard.app/privacy
- [ ] Support URL: https://msgguard.app/support

## Subscriptions
- [ ] Create group `MsgGuard Pro` (`MG000001-0001`)
- [ ] Add `com.ethanshen.msgguard.pro.monthly`
- [ ] Add `com.ethanshen.msgguard.pro.yearly`
- [ ] Match [Products.storekit](../../apps/ios/App/Shared/Store/Products.storekit)

## Compliance
- [ ] Privacy nutrition labels (see metadata.md)
- [ ] Export compliance: standard encryption only
- [ ] Message Filter Extension review notes: [review-notes.md](./review-notes.md)

## Build
- [ ] Archive via Xcode or `fastlane beta`
- [ ] Upload to TestFlight
- [ ] Internal testing pass on Extension + IAP

## Screenshots
See [SCREENSHOTS.md](./SCREENSHOTS.md)
