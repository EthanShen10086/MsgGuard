#!/usr/bin/env python3
"""Train sklearn pipeline for Core ML export."""
import json
from pathlib import Path

import joblib
import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import f1_score
from sklearn.model_selection import train_test_split
from sklearn.pipeline import Pipeline

ROOT = Path(__file__).parent.parent
DATA = ROOT / "data" / "processed" / "train.csv"
FALLBACK = ROOT / "data" / "processed" / "all.csv"
OUTPUT = ROOT / "output"
OUTPUT.mkdir(exist_ok=True)


def main():
    path = DATA if DATA.exists() else FALLBACK
    df = pd.read_csv(path)
    df["label"] = df["label"].apply(lambda x: 1 if x in ("spam", "phishing", "promotion") else 0)
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
    joblib.dump(pipe, OUTPUT / "coreml_pipeline.joblib")
    (OUTPUT / "coreml_metrics.json").write_text(json.dumps({"f1": f1}, indent=2))
    print(json.dumps({"f1": round(f1, 4)}))


if __name__ == "__main__":
    main()
