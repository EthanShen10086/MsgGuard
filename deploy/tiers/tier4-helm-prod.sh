#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/../.."
helm upgrade --install msgguard deploy/helm/msgguard -f deploy/helm/msgguard/values-prod.yaml -n msgguard-prod --create-namespace
echo "Tier4 Helm prod deployed"
