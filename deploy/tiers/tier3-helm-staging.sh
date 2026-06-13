#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/../.."
helm upgrade --install msgguard deploy/helm/msgguard -f deploy/helm/msgguard/values-staging.yaml -n msgguard-staging --create-namespace
echo "Tier3 Helm staging deployed"
