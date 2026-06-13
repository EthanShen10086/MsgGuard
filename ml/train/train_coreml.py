#!/usr/bin/env python3
"""Train sklearn pipeline for Core ML export (optionally per locale)."""
import argparse
import json
import sys
from pathlib import Path
from typing import Optional

import joblib
import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import f1_score
from sklearn.model_selection import train_test_split
from sklearn.pipeline import Pipeline

ROOT = Path(__file__).parent.parent
sys.path.insert(0, str(ROOT))
from locale_utils import data_locales_for, normalize_locale, output_dir  # noqa: E402

DATA = ROOT / "data" / "processed" / "train.csv"
FALLBACK = ROOT / "data" / "processed" / "all.csv"


def load_data(locale: Optional[str]) -> pd.DataFrame:
    path = DATA if DATA.exists() else FALLBACK
    if not path.exists():
        raise SystemExit("Run: make data")
    df = pd.read_csv(path)
    if locale and "locale" in df.columns:
        allowed = set(data_locales_for(locale))
        df = df[df["locale"].isin(allowed)]
        if df.empty:
            raise SystemExit(f"No training rows for locale {locale}")
    df["label"] = df["label"].apply(lambda x: 1 if x in ("spam", "phishing", "promotion") else 0)
    return df


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--locale", default=None)
    args = parser.parse_args()
    tag = normalize_locale(args.locale) if args.locale else None
    df = load_data(tag)
    out = output_dir(ROOT, tag) if tag else ROOT / "output"
    out.mkdir(parents=True, exist_ok=True)

    X_train, X_test, y_train, y_test = train_test_split(
        df["text"], df["label"], test_size=0.2, random_state=42, stratify=df["label"]
    )
    pipe = Pipeline([
        ("tfidf", TfidfVectorizer(max_features=5000, ngram_range=(1, 2))),
        ("clf", LogisticRegression(max_iter=500)),
    ])
    pipe.fit(X_train, y_train)
    preds = pipe.predict(X_test)
    f1 = f1_score(y_test, preds)
    joblib.dump(pipe, out / "coreml_pipeline.joblib")
    (out / "coreml_metrics.json").write_text(json.dumps({"f1": f1, "locale": tag or "all"}, indent=2))
    print(json.dumps({"locale": tag or "all", "f1": round(f1, 4), "rows": len(df)}))


if __name__ == "__main__":
    main()
