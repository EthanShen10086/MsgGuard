# LLM Outage Runbook

## Symptoms
- `/api/v1/classify` returns heuristic results only
- Logs: `circuit open` or LLM provider errors

## Steps
1. Check `curl localhost:8080/health`
2. Verify API keys: `QWEN_API_KEY`, `DEEPSEEK_API_KEY`
3. Set `features.cloud_llm: false` in config → restart gateway
4. Extension continues local L0–L2 filtering (no user impact for non-defer)
5. Monitor Jaeger for failed spans

## Rollback
Revert config and restart. Circuit breaker auto-closes after 30s.
