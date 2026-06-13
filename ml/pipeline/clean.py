#!/usr/bin/env python3
"""Clean and merge seed datasets into data/processed/all.csv."""
import csv
import hashlib
import re
from pathlib import Path

ROOT = Path(__file__).parent.parent
SEED = ROOT / "datasets" / "seed"
OUT = ROOT / "data" / "processed"
OUT.mkdir(parents=True, exist_ok=True)

PHONE_RE = re.compile(r"1[3-9]\d{9}")
ID_RE = re.compile(r"\d{17}[\dXx]")


def redact_pii(text: str) -> str:
    text = PHONE_RE.sub("[PHONE]", text)
    text = ID_RE.sub("[ID]", text)
    return text.strip()


def load_seed() -> list[dict]:
    rows = []
    for path in sorted(SEED.glob("*.csv")):
        with path.open(encoding="utf-8") as f:
            for row in csv.DictReader(f):
                text = redact_pii(row["text"])
                if len(text) < 4:
                    continue
                label = row["label"].strip().lower()
                if label not in ("spam", "ham", "promotion", "phishing"):
                    label = "spam" if label in ("1", "true") else "ham"
                h = hashlib.sha256(text.encode()).hexdigest()[:16]
                locale = "zh" if any("\u4e00" <= c <= "\u9fff" for c in text) else "en"
                rows.append({
                    "text": text, "label": label, "source": path.stem,
                    "locale": locale, "content_hash": h,
                })
    seen = set()
    deduped = []
    for r in rows:
        if r["content_hash"] in seen:
            continue
        seen.add(r["content_hash"])
        deduped.append(r)
    return deduped


def main():
    rows = load_seed()
    out_path = OUT / "all.csv"
    with out_path.open("w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=["text", "label", "source", "locale", "content_hash"])
        w.writeheader()
        w.writerows(rows)
    print(f"clean: wrote {len(rows)} rows -> {out_path}")


if __name__ == "__main__":
    main()
