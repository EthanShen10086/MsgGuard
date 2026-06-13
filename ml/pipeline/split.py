#!/usr/bin/env python3
"""Stratified train/val/test split."""
import csv
import json
from collections import defaultdict
from pathlib import Path
import random

ROOT = Path(__file__).parent.parent
INP = ROOT / "data" / "processed" / "all.csv"
OUT = ROOT / "data" / "processed"
REG = ROOT / "data" / "registry"


def main():
    rows = list(csv.DictReader(INP.open(encoding="utf-8")))
    by_label: dict[str, list] = defaultdict(list)
    for r in rows:
        by_label[r["label"]].append(r)
    random.seed(42)
    splits = {"train": [], "val": [], "test": []}
    for label, group in by_label.items():
        random.shuffle(group)
        n = len(group)
        t = int(n * 0.8)
        v = int(n * 0.1)
        splits["train"].extend(group[:t])
        splits["val"].extend(group[t : t + v])
        splits["test"].extend(group[t + v :])
    for name, data in splits.items():
        path = OUT / f"{name}.csv"
        with path.open("w", newline="", encoding="utf-8") as f:
            w = csv.DictWriter(f, fieldnames=rows[0].keys())
            w.writeheader()
            w.writerows(data)
        print(f"split: {name}={len(data)} -> {path}")
    REG.mkdir(parents=True, exist_ok=True)
    manifest = {
        "version": "1.0.0",
        "total": len(rows),
        "splits": {k: len(v) for k, v in splits.items()},
        "labels": {k: len(v) for k, v in by_label.items()},
    }
    (REG / "manifest.json").write_text(json.dumps(manifest, indent=2))


if __name__ == "__main__":
    main()
