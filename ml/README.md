# MsgGuard ML Pipeline

End-to-end data → train → benchmark → publish loop.

## Quick Start

```bash
cd ml
pip install -r requirements.txt
make data      # collect + clean + split
make train     # bayes + coreml
make benchmark # F1 gate
make infer TEXT="免费贷款无抵押"
```

## Layout

- `datasets/seed/` — committed sample data (474+ rows)
- `datasets/benchmark/` — fixed regression test sets
- `pipeline/` — collect, clean, label, merge, split
- `train/` — bayes + coreml training
- `export/` — Core ML export
- `benchmark/` — offline metrics + CI gate
- `flywheel/` — user feedback ingest + retrain

See `docs/cookbook/18-data-collection.md` through `22-benchmark.md`.
