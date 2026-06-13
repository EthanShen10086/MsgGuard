#!/bin/bash
set -euo pipefail
cd "$(dirname "$0")/../.."
bash deploy/mtls/gen-certs.sh
export MTLS_ADMIN_REQUIRED=true
export MTLS_CLIENT_HEADER=X-Client-Cert-Subject
docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.mtls.yml up -d --build
echo "Tier4 mTLS up. API https://localhost:8443 (client cert required for /api/v1/admin/*)"
echo "Site http://localhost:8088/support/"
