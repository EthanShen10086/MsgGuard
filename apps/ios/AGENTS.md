# MsgGuard Agent Guide

## Quick Start

```bash
cd apps/ios && bash setup.sh
open MsgGuard.xcodeproj
```

## Architecture

- `AppState` = global truth
- `FilterEngine` = classification (unit tested)
- `BlocklistStore` = App Group persistence
- Never log SMS body in production

## Conventions

- MGError for all user-facing errors
- All strings in Localizable.strings (en + zh-Hans)
- Conventional Commits
