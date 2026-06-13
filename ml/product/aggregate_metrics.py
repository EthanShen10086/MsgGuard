#!/usr/bin/env python3
"""Aggregate product metrics from analytics, feedback, and benchmark."""
import json
import os
import urllib.request
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).parent.parent
OUT = ROOT / "product" / "reports" / "weekly.json"
GATEWAY = os.environ.get("GATEWAY_URL", "http://localhost:8080")


def fetch_admin_summary(token: str) -> dict:
    req = urllib.request.Request(
        f"{GATEWAY}/api/v1/admin/metrics/summary",
        headers={"Authorization": f"Bearer {token}"},
    )
    try:
        with urllib.request.urlopen(req, timeout=5) as resp:
            return json.loads(resp.read())
    except Exception as e:
        return {"error": str(e)}


def load_benchmark() -> dict:
    report = ROOT / "benchmark" / "reports" / "latest" / "report.json"
    if report.exists():
        return json.loads(report.read_text())
    return {}


def issue_token() -> str:
    req = urllib.request.Request(
        f"{GATEWAY}/api/v1/auth/token",
        data=json.dumps({"user_id": "metrics-bot", "roles": ["admin"]}).encode(),
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(req, timeout=5) as resp:
        data = json.loads(resp.read())
        return data.get("AccessToken") or data.get("access_token", "")


def main():
    token = issue_token()
    summary = fetch_admin_summary(token)
    benchmark = load_benchmark()
    report = {
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "period_days": 7,
        "admin_summary": summary,
        "benchmark_overall": benchmark.get("overall", {}),
        "gate_passed": benchmark.get("gate_passed"),
        "recommendations": [],
    }
    overall = benchmark.get("overall", {})
    if overall.get("fpr", 0) > 0.02:
        report["recommendations"].append("FPR above 2%: expand adversarial dataset")
    if summary.get("feedback_total", 0) > 50:
        report["recommendations"].append("High feedback volume: review misclassification samples")
    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text(json.dumps(report, indent=2))
    print(json.dumps(report, indent=2))


if __name__ == "__main__":
    main()
