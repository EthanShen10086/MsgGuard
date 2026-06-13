#!/bin/bash
# Tier 1 with MongoDB driver — postgres still runs for side-by-side migration tests.
set -euo pipefail
cd "$(dirname "$0")/../.."
export CONFIG_PATH=deploy/config.tier1.mongodb.yaml
export DATABASE_DRIVER=mongodb
export DATABASE_DSN="mongodb://msgguard:msgguard@mongodb:27017/?authSource=admin"
docker compose -f deploy/docker-compose.yml up -d --build
echo "Tier1 (MongoDB) up. curl http://localhost:8080/health"
