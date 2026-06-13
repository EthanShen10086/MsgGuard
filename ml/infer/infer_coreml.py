#!/usr/bin/env python3
"""CLI Core ML pipeline inference (via joblib)."""
import json
import sys
from pathlib import Path

import joblib

ROOT = Path(__file__).parent.parent
PIPE = ROOT / "output" / "coreml_pipeline.joblib"


def main():
    text = sys.argv[1] if len(sys.argv) > 1 else "Free gift click here"
    pipe = joblib.load(PIPE)
    pred = int(pipe.predict([text])[0])
    print(json.dumps({"text": text, "label": "spam" if pred else "ham"}))


if __name__ == "__main__":
    main()
