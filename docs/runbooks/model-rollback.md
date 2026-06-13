# Model Rollback Runbook

## When
Benchmark gate fails or production误杀率上升

## Steps
1. Check `ml/benchmark/reports/latest/report.json`
2. List versions: `curl localhost:8083/api/v1/models/latest`
3. Register previous version: `python ml/flywheel/publish_to_backend.py` with old artifacts
4. iOS: force sync via app restart → `ModelUpdateService.checkAndUpdate()`
5. Verify Extension uses previous modelVersion in App Group config
