#!/usr/bin/env python3
"""Compare Bayes vs heuristic layer predictions."""
import json
from pathlib import Path

import joblib

ROOT = Path(__file__).parent.parent
BENCH = ROOT / "datasets" / "benchmark" / "test_zh.jsonl"
PIPE = ROOT / "output" / "bayes_pipeline.joblib"
SPAM_WORDS = ["免费", "中奖", "贷款", "free", "winner", "click"]


def heuristic(text: str) -> str:
    hits = sum(1 for w in SPAM_WORDS if w in text.lower())
    return "spam" if hits >= 2 else "ham"


def main():
    pipe = joblib.load(PIPE)
    disagreements = []
    for line in BENCH.read_text().splitlines():
        row = json.loads(line)
        bayes = "spam" if pipe.predict([row["text"]])[0] in (1, "spam") else "ham"
        heur = heuristic(row["text"])
        if bayes != heur:
            disagreements.append({"text": row["text"], "bayes": bayes, "heuristic": heur, "label": row["label"]})
    out = ROOT / "benchmark" / "reports" / "latest" / "compare_layers.json"
    out.write_text(json.dumps(disagreements[:50], indent=2, ensure_ascii=False))
    print(f"disagreements: {len(disagreements)}")


if __name__ == "__main__":
    main()
