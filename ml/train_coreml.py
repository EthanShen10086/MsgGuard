#!/usr/bin/env python3
"""Train TF-IDF + Logistic Regression and export Core ML model."""

import json
from pathlib import Path

import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.model_selection import train_test_split
from sklearn.pipeline import Pipeline
from sklearn.metrics import classification_report
import joblib

OUTPUT = Path(__file__).parent / "output"
OUTPUT.mkdir(exist_ok=True)


def load_data():
    rows = [
        ("免费贷款无抵押当天放款", "spam"), ("恭喜中奖点击链接", "spam"),
        ("推广优惠活动限时", "promotion"), ("您的验证码是123456", "ham"),
        ("快递取件码8888", "ham"), ("Free gift winner click", "spam"),
        ("Your order has shipped", "ham"),
    ]
    return pd.DataFrame(rows, columns=["text", "label"])


def main():
    df = load_data()
    X_train, X_test, y_train, y_test = train_test_split(df["text"], df["label"], test_size=0.2, random_state=42)
    pipe = Pipeline([
        ("tfidf", TfidfVectorizer(max_features=3000)),
        ("clf", LogisticRegression(max_iter=1000)),
    ])
    pipe.fit(X_train, y_train)
    preds = pipe.predict(X_test)
    report = classification_report(y_test, preds, output_dict=True)
    joblib.dump(pipe, OUTPUT / "coreml_pipeline.joblib")
    (OUTPUT / "coreml_metrics.json").write_text(json.dumps(report, indent=2))
    print("Training complete. Run export_coreml.py to convert.")


if __name__ == "__main__":
    main()
