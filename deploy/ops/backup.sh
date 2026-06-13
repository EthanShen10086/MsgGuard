#!/bin/bash
# Backup PostgreSQL for MsgGuard
set -euo pipefail
DSN="${DATABASE_DSN:-postgres://msgguard:msgguard@localhost:5432/msgguard}"
OUT="${1:-./backup-$(date +%Y%m%d).sql}"
pg_dump "$DSN" > "$OUT"
echo "Backup written to $OUT"
