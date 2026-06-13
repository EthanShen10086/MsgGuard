#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/../.."
export CONFIG_PATH=deploy/config.tier1.yaml
docker compose -f deploy/docker-compose.yml up -d --build
echo "Tier1 up. curl http://localhost:8080/health"
