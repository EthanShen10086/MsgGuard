#!/usr/bin/env bash
# Production Helm with mTLS overlay
set -euo pipefail
cd "$(dirname "$0")/../.."
helm upgrade --install msgguard deploy/helm/msgguard \
  -f deploy/helm/msgguard/values-prod.yaml \
  -f deploy/helm/msgguard/values-mtls.yaml \
  -n msgguard-prod --create-namespace
echo "Tier4 Helm prod (mTLS) deployed"
