#!/usr/bin/env python3
"""Run offline benchmark on fixed test sets."""
import json
import sys
from pathlib import Path

import joblib
import yaml

ROOT = Path(__file__).parent.parent
BENCH = ROOT / "datasets" / "benchmark"
OUTPUT = ROOT / "output" / "bayes_pipeline.joblib"
REPORTS = ROOT / "benchmark" / "reports" / "latest"
REPORTS.mkdir(parents=True, exist_ok=True)

sys.path.insert(0, str(ROOT / "benchmark"))
from metrics import compute_metrics  # noqa: E402


def load_jsonl(path: Path):
    rows = []
    for line in path.read_text(encoding="utf-8").splitlines():
        if line.strip():
            rows.append(json.loads(line))
    return rows


def predict(pipe, text: str) -> str:
    pred = pipe.predict([text])[0]
    return "spam" if pred in (1, "spam", "1") else "ham"


def main():
    if not OUTPUT.exists():
        raise SystemExit("Run: make train")
    pipe = joblib.load(OUTPUT)
    baselines = yaml.safe_load((ROOT / "benchmark" / "baselines.yaml").read_text())
    results = {}
    all_true, all_pred = [], []
    for name in ["test_zh.jsonl", "test_en.jsonl", "adversarial.jsonl"]:
        path = BENCH / name
        if not path.exists():
            continue
        rows = load_jsonl(path)
        y_true = [r["label"] for r in rows]
        y_pred = [predict(pipe, r["text"]) for r in rows]
        m = compute_metrics(y_true, y_pred)
        results[name] = {
            "f1": round(m.f1, 4), "precision": round(m.precision, 4),
            "recall": round(m.recall, 4), "fpr": round(m.fpr, 4),
            "fnr": round(m.fnr, 4), "total": m.total,
        }
        all_true.extend(y_true)
        all_pred.extend(y_pred)
    overall = compute_metrics(all_true, all_pred)
    report = {
        "overall": {
            "f1": round(overall.f1, 4), "fpr": round(overall.fpr, 4),
            "fnr": round(overall.fnr, 4), "total": overall.total,
        },
        "by_set": results,
        "gate_passed": overall.f1 >= baselines["f1_min"] and overall.fpr <= baselines["fpr_max"],
        "baselines": baselines,
    }
    (REPORTS / "report.json").write_text(json.dumps(report, indent=2))
    html = f"""<!DOCTYPE html><html><head><title>MsgGuard Benchmark</title></head>
<body><h1>MsgGuard Benchmark Report</h1>
<pre>{json.dumps(report, indent=2)}</pre>
<p>Gate: {'PASS' if report['gate_passed'] else 'FAIL'}</p></body></html>"""
    (REPORTS / "report.html").write_text(html)
    print(json.dumps(report["overall"], indent=2))
    print(f"gate_passed={report['gate_passed']}")
    sys.exit(0 if report["gate_passed"] else 1)


if __name__ == "__main__":
    main()
