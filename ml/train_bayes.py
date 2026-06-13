#!/usr/bin/env python3
"""Train Naive Bayes spam classifier and export metrics."""

import json
from pathlib import Path

import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics import classification_report, f1_score
from sklearn.model_selection import train_test_split
from sklearn.naive_bayes import MultinomialNB
from sklearn.pipeline import Pipeline
import joblib

DATASET = Path(__file__).parent / "datasets" / "sms_spam.csv"
OUTPUT = Path(__file__).parent / "output"
OUTPUT.mkdir(exist_ok=True)


def load_data():
    if not DATASET.exists():
        # Built-in seed data
        rows = [
            ("免费贷款无抵押", "spam"), ("恭喜中奖领取", "spam"),
            ("Your verification code is 1234", "ham"), ("快递已到达取件码888", "ham"),
            ("Free gift click here now", "spam"), ("Order shipped tracking", "ham"),
        ]
        df = pd.DataFrame(rows, columns=["text", "label"])
    else:
        df = pd.read_csv(DATASET)
    return df


def main():
    df = load_data()
    X_train, X_test, y_train, y_test = train_test_split(df["text"], df["label"], test_size=0.2, random_state=42)
    pipe = Pipeline([
        ("tfidf", TfidfVectorizer(max_features=5000, ngram_range=(1, 2))),
        ("clf", MultinomialNB()),
    ])
    pipe.fit(X_train, y_train)
    preds = pipe.predict(X_test)
    report = classification_report(y_test, preds, output_dict=True)
    metrics = {"f1": f1_score(y_test, preds, pos_label="spam"), "report": report}
    joblib.dump(pipe, OUTPUT / "bayes_pipeline.joblib")
    (OUTPUT / "bayes_metrics.json").write_text(json.dumps(metrics, indent=2))
    print(json.dumps(metrics, indent=2))


if __name__ == "__main__":
    main()
