#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/.."
make ios-test 2>/dev/null || echo "iOS: run cd apps/ios && bash setup.sh"
cd services/gateway && CONFIG_PATH=../../deploy/config.yaml go run ./cmd/server
