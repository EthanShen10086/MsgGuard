#!/usr/bin/env bash
# Sync rules bundle to CDN / object storage (R2, S3, or local static dir).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
VERSION="${1:-1.0.0}"
OUT="${CDN_OUT:-$ROOT/deploy/site/rules}"
mkdir -p "$OUT/$VERSION"
curl -sf "${GATEWAY_URL:-http://localhost:8080}/api/v1/rules/latest" -o "$OUT/$VERSION/rules.json"
cp "$OUT/$VERSION/rules.json" "$OUT/latest.json"
echo "Synced rules v$VERSION to $OUT"
