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

# Subscription funnel steps (analytics v2) — populated when events exist
FUNNEL_STEPS = [
    "paywall_view",
    "product_selected",
    "purchase_attempt",
    "purchase_success",
    "purchase_cancel",
    "purchase_failed",
    "restore_attempt",
    "restore_success",
]


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


def build_subscription_funnel(event_counts: dict) -> dict:
    """Derive funnel counts from event_counts; placeholders until v2 events ship."""
    funnel_events = event_counts.get("subscription_funnel", 0)
    purchase_completed = event_counts.get("purchase_completed", 0)
    purchase_started = event_counts.get("purchase_started", 0)
    onboarding = event_counts.get("onboarding_completed", 0)

    steps = {step: 0 for step in FUNNEL_STEPS}
    # Map legacy events to funnel steps when dedicated funnel events absent
    if funnel_events == 0:
        steps["paywall_view"] = purchase_started
        steps["purchase_success"] = purchase_completed
    else:
        steps["paywall_view"] = funnel_events  # refined when props.step aggregated server-side

    conversion = None
    if steps["paywall_view"] > 0:
        conversion = round(steps["purchase_success"] / steps["paywall_view"], 4)

    return {
        "steps": steps,
        "onboarding_completed": onboarding,
        "conversion_paywall_to_purchase": conversion,
        "note": "Placeholder until subscription_funnel v2 events aggregated by step",
    }


def main():
    token = issue_token()
    summary = fetch_admin_summary(token)
    benchmark = load_benchmark()
    event_counts = summary.get("event_counts") or {}
    report = {
        "generated_at": datetime.now(timezone.utc).isoformat(),
        "period_days": summary.get("period_days", 7),
        "admin_summary": summary,
        "subscription_funnel": build_subscription_funnel(event_counts),
        "benchmark_overall": benchmark.get("overall", {}),
        "gate_passed": benchmark.get("gate_passed"),
        "recommendations": [],
    }
    overall = benchmark.get("overall", {})
    if overall.get("fpr", 0) > 0.02:
        report["recommendations"].append("FPR above 2%: expand adversarial dataset")
    if summary.get("feedback_total", 0) > 50:
        report["recommendations"].append("High feedback volume: review misclassification samples")
    funnel = report["subscription_funnel"]
    if funnel.get("conversion_paywall_to_purchase") is not None:
        if funnel["conversion_paywall_to_purchase"] < 0.05:
            report["recommendations"].append("Subscription conversion below 5%: review onboarding A/B")
    OUT.parent.mkdir(parents=True, exist_ok=True)
    OUT.write_text(json.dumps(report, indent=2))
    print(json.dumps(report, indent=2))


if __name__ == "__main__":
    main()
