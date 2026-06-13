# Database Failure Runbook

## Symptoms
- Feedback API 500 errors
- Gateway logs: postgres connection refused

## Steps
1. Check postgres: `docker compose ps postgres`
2. Fallback: set `database.driver: memory` in config (data loss for new feedback)
3. Restore from backup: `deploy/ops/backup.sh`
4. Re-point `DATABASE_DSN` and restart gateway

## Prevention
- Weekly cron: `deploy/ops/backup.sh /backups/msgguard.sql`
- Staging restore drill monthly
