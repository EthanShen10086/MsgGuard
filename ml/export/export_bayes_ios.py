#!/usr/bin/env python3
"""Export iOS-compatible Bayes JSON from locale training CSV."""
import argparse
import json
import re
import sys
from pathlib import Path

import pandas as pd

ROOT = Path(__file__).parent.parent
sys.path.insert(0, str(ROOT))


def tokenize(text: str) -> list[str]:
    return [t for t in re.split(r"\W+", text.lower()) if len(t) > 1]


def load_locale_df(locale: str) -> pd.DataFrame:
    from locale_utils import data_locales_for, normalize_locale, output_dir

    tag = normalize_locale(locale)
    train = ROOT / "data" / "processed" / "train.csv"
    fallback = ROOT / "data" / "processed" / "all.csv"
    path = train if train.exists() else fallback
    if not path.exists():
        raise SystemExit("Run: make data")
    df = pd.read_csv(path)
    if "locale" in df.columns:
        allowed = set(data_locales_for(tag))
        df = df[df["locale"].isin(allowed)]
    if df.empty:
        raise SystemExit(f"No rows for locale {tag}")
    return df


def export_bayes_ios(locale: str) -> Path:
    from locale_utils import normalize_locale, output_dir

    tag = normalize_locale(locale)
    df = load_locale_df(tag)
    word_counts: dict[str, dict[str, int]] = {"spam": {}, "ham": {}}
    category_counts: dict[str, int] = {"spam": 0, "ham": 0}
    total = 0
    for _, row in df.iterrows():
        label = row["label"]
        bucket = "spam" if label in ("spam", "phishing", "promotion") else "ham"
        category_counts[bucket] += 1
        total += 1
        for token in tokenize(str(row["text"])):
            word_counts[bucket][token] = word_counts[bucket].get(token, 0) + 1
    payload = {
        "wordCounts": word_counts,
        "categoryCounts": category_counts,
        "totalDocuments": total,
    }
    out = output_dir(ROOT, tag) / "bayes_model.json"
    out.write_text(json.dumps(payload, ensure_ascii=False, indent=2), encoding="utf-8")
    return out


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--locale", default="zh-Hans")
    args = parser.parse_args()
    path = export_bayes_ios(args.locale)
    print(json.dumps({"path": str(path), "locale": args.locale}, indent=2))


if __name__ == "__main__":
    main()
