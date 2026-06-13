#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/../.."
docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.gpu.yml up -d --build
echo "Tier2 GPU compose up"
