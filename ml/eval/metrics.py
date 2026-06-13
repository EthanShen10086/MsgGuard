"""Evaluation harness for spam classifiers."""

import json
from pathlib import Path


def load_metrics(name: str) -> dict:
    path = Path(__file__).parent.parent / "output" / f"{name}_metrics.json"
    if not path.exists():
        return {"error": f"{path} not found"}
    return json.loads(path.read_text())


if __name__ == "__main__":
    for model in ("bayes", "coreml"):
        print(f"=== {model} ===")
        print(json.dumps(load_metrics(model), indent=2))
