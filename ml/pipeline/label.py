#!/usr/bin/env python3
"""Semi-auto label pass-through for unlabeled rows."""
import csv
from pathlib import Path

ROOT = Path(__file__).parent.parent
INP = ROOT / "data" / "processed" / "all.csv"
OUT = ROOT / "data" / "labeled"
OUT.mkdir(parents=True, exist_ok=True)


def main():
    rows = list(csv.DictReader(INP.open(encoding="utf-8")))
    out = OUT / "labeled.csv"
    with out.open("w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=rows[0].keys())
        w.writeheader()
        w.writerows(rows)
    print(f"label: {len(rows)} rows -> {out}")


if __name__ == "__main__":
    main()
