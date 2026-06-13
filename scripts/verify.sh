#!/usr/bin/env bash
# MsgGuard acceptance verification — run from repo root.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

PASS=0
FAIL=0

ok()   { echo "✓ $1"; PASS=$((PASS + 1)); }
fail() { echo "✗ $1"; FAIL=$((FAIL + 1)); }

echo "=== Go build ==="
for svc in gateway model admin rules classify feedback flywheel; do
  if (cd "services/$svc" && go build ./...); then
    ok "services/$svc"
  else
    fail "services/$svc"
  fi
done
for pkg in app adapters httpauth; do
  if (cd "pkg/$pkg" && go build ./...); then
    ok "pkg/$pkg"
  else
    fail "pkg/$pkg"
  fi
done

echo "=== ML benchmark ==="
BENCH_LOG="$(cd ml && make benchmark 2>&1)" || true
if echo "$BENCH_LOG" | rg -q 'gate_passed=True'; then
  ok "ml benchmark gate"
else
  fail "ml benchmark gate"
fi

echo "=== Product metrics ==="
if python3 ml/product/aggregate_metrics.py && test -f ml/product/reports/weekly.json; then
  ok "aggregate_metrics.py"
else
  fail "aggregate_metrics.py"
fi

echo "=== Gateway smoke (auth + model auth) ==="
GWPID=""
MODELPID=""
GW_PORT="${GW_PORT:-18080}"
cleanup() {
  [[ -n "$GWPID" ]] && kill "$GWPID" 2>/dev/null || true
  [[ -n "$MODELPID" ]] && kill "$MODELPID" 2>/dev/null || true
}
trap cleanup EXIT

(cd services/model && AUTH_SECRET=msgguard-dev-secret PORT=8083 go run ./cmd/server) &
MODELPID=$!
sleep 2

(cd services/gateway && AUTH_SECRET=msgguard-dev-secret PORT="$GW_PORT" CONFIG_PATH=../../deploy/config.yaml go run ./cmd/server) &
GWPID=$!
sleep 3

GW="http://localhost:${GW_PORT}"
if curl -sf "$GW/health" | grep -q ok; then
  ok "GET /health"
else
  fail "GET /health"
fi

TOKEN=$(curl -sf -X POST "$GW/api/v1/auth/token" \
  -H 'Content-Type: application/json' \
  -d '{"user_id":"verify","roles":["admin"]}' | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])")
if [[ -n "$TOKEN" ]]; then
  ok "POST /api/v1/auth/token"
else
  fail "POST /api/v1/auth/token"
fi

CODE=$(curl -s -o /dev/null -w '%{http_code}' "$GW/api/v1/feedback")
if [[ "$CODE" == "401" ]]; then
  ok "GET /api/v1/feedback → 401"
else
  fail "GET /api/v1/feedback → expected 401 got $CODE"
fi

REG_CODE=$(curl -s -o /dev/null -w '%{http_code}' -X POST http://localhost:8083/api/v1/models/register \
  -H 'Content-Type: application/json' -d '{"version":"0.0.1","locale":"zh-Hans","checksum":"x"}')
if [[ "$REG_CODE" == "401" ]]; then
  ok "model register without token → 401"
else
  fail "model register without token → expected 401 got $REG_CODE"
fi

REG_OK=$(curl -s -o /dev/null -w '%{http_code}' -X POST http://localhost:8083/api/v1/models/register \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"version":"0.0.1","locale":"zh-Hans","checksum":"x"}')
if [[ "$REG_OK" == "201" ]]; then
  ok "model register with token → 201"
else
  fail "model register with token → expected 201 got $REG_OK"
fi

echo "=== iOS build ==="
if (cd apps/ios && xcodegen generate >/dev/null && sleep 2 && \
    xcodebuild -scheme MsgGuard-iOS -destination 'platform=iOS Simulator,name=iPhone 16' clean build 2>&1 | rg -q 'BUILD SUCCEEDED'); then
  ok "iOS BUILD SUCCEEDED"
else
  fail "iOS build"
fi

echo "=== Helm + legal ==="
if test -f docs/legal/PRIVACY.md && test -f deploy/site/privacy/index.html; then
  ok "privacy policy files"
else
  fail "privacy policy files"
fi
if command -v helm >/dev/null 2>&1; then
  if helm template msgguard deploy/helm/msgguard -f deploy/helm/msgguard/values-mongodb.yaml 2>/dev/null | rg -q 'mongodb'; then
    ok "helm mongodb values"
  else
    fail "helm mongodb values"
  fi
elif rg -q 'driver: mongodb' deploy/helm/msgguard/values-mongodb.yaml; then
  ok "helm mongodb values (values file)"
else
  fail "helm mongodb values"
fi
if test -f docs/app-store/metadata.md && test -f docs/app-store/review-notes.md; then
  ok "app store metadata draft"
else
  fail "app store metadata draft"
fi
if test -f docs/app-store/ASC_CHECKLIST.md && test -f fastlane/Fastfile; then
  ok "app store ship checklist"
else
  fail "app store ship checklist"
fi
if test -f docs/deploy/TIER_MATRIX.md; then
  ok "tier matrix doc"
else
  fail "tier matrix doc"
fi

echo "=== Support + mTLS ==="
if test -f docs/legal/SUPPORT.md && test -f deploy/site/support/index.html; then
  ok "support page"
else
  fail "support page"
fi
if test -x deploy/mtls/gen-certs.sh && test -f pkg/httpauth/clientcert.go; then
  ok "mtls assets"
else
  fail "mtls assets"
fi
if test -f deploy/grafana/provisioning/datasources/prometheus.yml; then
  ok "grafana provisioning"
else
  fail "grafana provisioning"
fi
if test -f deploy/helm/msgguard/templates/cert-manager/certificate.yaml; then
  ok "cert-manager helm template"
else
  fail "cert-manager helm template"
fi
if (cd pkg/httpauth && go test ./...); then
  ok "httpauth mTLS tests"
else
  fail "httpauth mTLS tests"
fi

echo ""
echo "Results: $PASS passed, $FAIL failed"
if [[ "$FAIL" -gt 0 ]]; then
  exit 1
fi
echo "All checks passed."
