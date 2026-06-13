#!/usr/bin/env python3
"""CLI Bayes inference."""
import json
import sys
from pathlib import Path

import joblib

ROOT = Path(__file__).parent.parent
PIPE = ROOT / "output" / "bayes_pipeline.joblib"


def main():
    text = sys.argv[1] if len(sys.argv) > 1 else "免费贷款无抵押"
    pipe = joblib.load(PIPE)
    pred = pipe.predict([text])[0]
    proba = pipe.predict_proba([text])[0].max() if hasattr(pipe, "predict_proba") else 0.0
    print(json.dumps({"text": text, "label": str(pred), "confidence": round(float(proba), 4)}))


if __name__ == "__main__":
    main()
