# MsgGuard ML Pipeline

End-to-end data → train → benchmark → publish loop.

## Quick Start

```bash
cd ml
pip install -r requirements.txt
make data                 # collect + clean + split
make train-all-locales    # zh-Hans + en-US bayes/coreml + iOS bayes JSON
make export-all-locales   # Core ML + bayes_model.json per locale
make benchmark            # per-locale F1 gates
make infer TEXT="免费贷款无抵押"
```

## Per-locale training

| Locale | Data filter | Output dir |
|--------|-------------|------------|
| `zh-Hans` | CSV `locale` in `zh`, `zh-Hans` | `output/zh-Hans/` |
| `en-US` | CSV `locale` in `en`, `en-US` | `output/en-US/` |

```bash
make train-locale LOCALE=zh-Hans
make export-locale LOCALE=en-US
MODEL_LOCALES=zh-Hans,en-US python3 flywheel/publish_to_backend.py
```

Benchmark gates are defined in `benchmark/baselines.yaml` under `locales:`.

## Layout

- `datasets/seed/` — committed sample data (474+ rows)
- `datasets/benchmark/` — fixed regression test sets (zh/en/adversarial)
- `locale_utils.py` — BCP-47 mapping + output paths
- `pipeline/` — collect, clean, label, merge, split
- `train/` — bayes + coreml training (`--locale`)
- `export/` — Core ML + iOS bayes JSON export
- `benchmark/` — offline metrics + per-locale CI gate
- `flywheel/` — user feedback ingest + `retrain-all`

See `docs/cookbook/18-data-collection.md` through `22-benchmark.md`.
