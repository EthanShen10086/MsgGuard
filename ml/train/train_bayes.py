#!/usr/bin/env python3
"""Train Naive Bayes from processed data."""
import json
from pathlib import Path

import joblib
import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics import classification_report, f1_score
from sklearn.model_selection import train_test_split
from sklearn.naive_bayes import MultinomialNB
from sklearn.pipeline import Pipeline

ROOT = Path(__file__).parent.parent
DATA = ROOT / "data" / "processed" / "train.csv"
FALLBACK = ROOT / "data" / "processed" / "all.csv"
OUTPUT = ROOT / "output"
OUTPUT.mkdir(exist_ok=True)


def load_data():
    path = DATA if DATA.exists() else FALLBACK
    if not path.exists():
        raise SystemExit("Run: make data")
    df = pd.read_csv(path)
    df["label"] = df["label"].apply(lambda x: "spam" if x in ("spam", "phishing", "promotion") else "ham")
    return df


def main():
    df = load_data()
    X_train, X_test, y_train, y_test = train_test_split(
        df["text"], df["label"], test_size=0.2, random_state=42, stratify=df["label"]
    )
    pipe = Pipeline([
        ("tfidf", TfidfVectorizer(max_features=8000, ngram_range=(1, 2))),
        ("clf", MultinomialNB()),
    ])
    pipe.fit(X_train, y_train)
    preds = pipe.predict(X_test)
    f1 = f1_score(y_test, preds, pos_label="spam")
    report = classification_report(y_test, preds, output_dict=True)
    metrics = {"f1": f1, "report": report}
    joblib.dump(pipe, OUTPUT / "bayes_pipeline.joblib")
    (OUTPUT / "bayes_metrics.json").write_text(json.dumps(metrics, indent=2))
    print(json.dumps({"f1": round(f1, 4)}, indent=2))


if __name__ == "__main__":
    main()
