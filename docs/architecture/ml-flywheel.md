# ML Flywheel Architecture

**Last updated:** 2026-06-17

## Loop

```
User feedback / samples → feedback store → NATS trigger → flywheel worker
    → label / merge datasets → train → benchmark gate → export CoreML
    → model service publish → iOS OTA
```

## Pipeline (`ml/`)

| Stage | Script / Makefile |
|-------|-------------------|
| Collect | `pipeline/collect_*.py`, `sources.yaml` |
| Clean / label | `pipeline/clean.py`, `label.py` |
| Train Bayes | `train/train_bayes.py` |
| Train CoreML | `train/train_coreml.py` |
| Benchmark | `benchmark/run_benchmark.py` — FPR gate |
| Export iOS | `export/export_bayes_ios.py`, `verify_coreml_export.py` |
| Product metrics | `product/aggregate_metrics.py` |

## Gates

- CI runs benchmark; `gate_passed=False` blocks merge
- Adversarial OTP FPR tracked separately (`OTPGuard` + benchmark baselines)

## Scheduling

- `deploy/k8s/cron-retrain.yaml` — K8s CronJob
- `ml/flywheel/schedule_retrain.sh` — manual / VPS cron

## Shadow Mode

Gateway `/api/v1/classify/shadow` compares local vs cloud; Prometheus at `/metrics/shadow`.
