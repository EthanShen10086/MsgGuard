#!/usr/bin/env python3
"""Run offline benchmark with per-locale pipelines and gates."""
import json
import sys
from pathlib import Path

import joblib
import yaml

ROOT = Path(__file__).parent.parent
BENCH = ROOT / "datasets" / "benchmark"
REPORTS = ROOT / "benchmark" / "reports" / "latest"
REPORTS.mkdir(parents=True, exist_ok=True)

sys.path.insert(0, str(ROOT))
from locale_utils import BENCH_LOCALE, configured_locales, normalize_locale, output_dir  # noqa: E402

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


def pipeline_for_locale(locale: str):
    tag = normalize_locale(locale)
    path = output_dir(ROOT, tag) / "bayes_pipeline.joblib"
    if not path.exists():
        path = ROOT / "output" / "bayes_pipeline.joblib"
    if not path.exists():
        raise SystemExit(f"Missing pipeline for {tag}; run: make train-all-locales")
    return joblib.load(path), tag


def main():
    baselines = yaml.safe_load((ROOT / "benchmark" / "baselines.yaml").read_text())
    per_locale = baselines.get("locales", {})
    results = {}
    locale_gates = {}
    all_true, all_pred = [], []

    for bench_file, locale in BENCH_LOCALE.items():
        path = BENCH / bench_file
        if not path.exists():
            continue
        pipe, tag = pipeline_for_locale(locale)
        rows = load_jsonl(path)
        y_true = [r["label"] for r in rows]
        y_pred = [predict(pipe, r["text"]) for r in rows]
        m = compute_metrics(y_true, y_pred)
        results[bench_file] = {
            "locale": tag,
            "f1": round(m.f1, 4), "precision": round(m.precision, 4),
            "recall": round(m.recall, 4), "fpr": round(m.fpr, 4),
            "fnr": round(m.fnr, 4), "total": m.total,
        }
        gate = per_locale.get(tag, {})
        f1_min = gate.get("f1_min", baselines["f1_min"])
        fpr_max = gate.get("fpr_max", baselines["fpr_max"])
        locale_gates[tag] = m.f1 >= f1_min and m.fpr <= fpr_max
        all_true.extend(y_true)
        all_pred.extend(y_pred)

    adv_path = BENCH / "adversarial.jsonl"
    adv = {}
    adv_gate = True
    if adv_path.exists():
        pipe, _ = pipeline_for_locale("zh-Hans")
        rows = load_jsonl(adv_path)
        y_true = [r["label"] for r in rows]
        y_pred = [predict(pipe, r["text"]) for r in rows]
        m = compute_metrics(y_true, y_pred)
        adv = {"f1": round(m.f1, 4), "fpr": round(m.fpr, 4), "fnr": round(m.fnr, 4), "total": m.total}
        adv_gate = adv["fpr"] <= baselines.get("adversarial_fpr_max", baselines["fpr_max"])
        results["adversarial.jsonl"] = adv

    overall = compute_metrics(all_true, all_pred)
    gate_passed = (
        overall.f1 >= baselines["f1_min"]
        and overall.fpr <= baselines["fpr_max"]
        and adv_gate
        and all(locale_gates.values()) if locale_gates else True
    )
    report = {
        "overall": {
            "f1": round(overall.f1, 4), "fpr": round(overall.fpr, 4),
            "fnr": round(overall.fnr, 4), "total": overall.total,
        },
        "by_set": results,
        "locale_gates": locale_gates,
        "adversarial_fpr": adv.get("fpr"),
        "gate_passed": gate_passed,
        "baselines": baselines,
        "trained_locales": configured_locales(),
    }
    (REPORTS / "report.json").write_text(json.dumps(report, indent=2))
    html = f"""<!DOCTYPE html><html><head><title>MsgGuard Benchmark</title></head>
<body><h1>MsgGuard Benchmark Report</h1>
<pre>{json.dumps(report, indent=2)}</pre>
<p>Gate: {'PASS' if gate_passed else 'FAIL'}</p></body></html>"""
    (REPORTS / "report.html").write_text(html)
    print(json.dumps(report["overall"], indent=2))
    print(json.dumps({"locale_gates": locale_gates, "gate_passed": gate_passed}))
    sys.exit(0 if gate_passed else 1)


if __name__ == "__main__":
    main()
