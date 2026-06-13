#!/usr/bin/env python3
"""Merge external + user feedback datasets."""
import csv
from pathlib import Path

ROOT = Path(__file__).parent.parent
PROCESSED = ROOT / "data" / "processed"
USER = ROOT / "data" / "labeled" / "user"


def main():
    rows = list(csv.DictReader((PROCESSED / "all.csv").open(encoding="utf-8")))
    if USER.exists():
        for path in USER.glob("*.csv"):
            for row in csv.DictReader(path.open(encoding="utf-8")):
                row["source"] = "user_feedback"
                rows.append(row)
    seen = {}
    for r in rows:
        seen[r.get("content_hash", r["text"])] = r
    merged = list(seen.values())
    out = PROCESSED / "merged.csv"
    with out.open("w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=["text", "label", "source", "locale", "content_hash"])
        w.writeheader()
        w.writerows(merged)
    print(f"merge: {len(merged)} rows -> {out}")


if __name__ == "__main__":
    main()
