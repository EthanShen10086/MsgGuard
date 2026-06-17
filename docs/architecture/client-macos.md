# macOS Client Architecture

**Last updated:** 2026-06-17

## Targets

| Target | Role |
|--------|------|
| MailHost | Host app for MailKit extension distribution |
| MailExtension | MailKit message filter (RFC822 body parsing) |

## Status

**In Progress** — extension skeleton, TestFlight/notarize lanes in Fastfile; not App Store primary SKU.

## Flow

```
Mail.app → MailExtension → shared FilterEngine patterns → allow/junk
                ↓ optional
           Gateway classify (same API as iOS defer)
```

## Key Files

- `apps/macos/MailExtension/`
- `apps/macos/MailHost/`
- `docs/product/MAIL_EXTENSION.md`
