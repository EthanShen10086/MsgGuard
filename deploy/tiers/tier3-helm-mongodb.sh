#!/bin/bash
# Helm staging with MongoDB driver (requires MongoDB + Redis in cluster).
set -euo pipefail
cd "$(dirname "$0")/../.."
helm upgrade --install msgguard deploy/helm/msgguard \
  -f deploy/helm/msgguard/values-staging.yaml \
  -f deploy/helm/msgguard/values-mongodb.yaml \
  -n msgguard-staging --create-namespace
echo "Tier3 Helm staging (MongoDB) deployed"
