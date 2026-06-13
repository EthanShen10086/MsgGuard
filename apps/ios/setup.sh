#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"

echo "==> Installing tools (if missing)"
command -v xcodegen >/dev/null || brew install xcodegen
command -v swiftlint >/dev/null || brew install swiftlint
command -v swiftformat >/dev/null || brew install swiftformat

echo "==> Building Swift packages"
for pkg in SharedModels DesignSystem FilterEngine BlocklistStore; do
  echo "Building $pkg..."
  (cd "Packages/$pkg" && swift build && swift test)
done

echo "==> Generating Xcode project"
xcodegen generate

echo "==> Done. Open MsgGuard.xcodeproj"
