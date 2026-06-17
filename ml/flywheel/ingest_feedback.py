#!/usr/bin/env python3
"""Ingest feedback from gateway API into labeled/user/."""
import csv
import json
import os
import urllib.request
from pathlib import Path

ROOT = Path(__file__).parent.parent
OUT = ROOT / "data" / "labeled" / "user"
OUT.mkdir(parents=True, exist_ok=True)
URL = os.environ.get("GATEWAY_URL", "http://localhost:8080") + "/api/v1/feedback"


def _auth_header() -> dict:
    token = os.environ.get("FEEDBACK_INGEST_TOKEN", "")
    if not token:
        bootstrap = os.environ.get("GATEWAY_URL", "http://localhost:8080") + "/api/v1/auth/token"
        try:
            req = urllib.request.Request(
                bootstrap,
                data=json.dumps({"user_id": "flywheel", "roles": ["ml_engineer"]}).encode(),
                headers={"Content-Type": "application/json"},
                method="POST",
            )
            with urllib.request.urlopen(req, timeout=5) as resp:
                token = json.loads(resp.read()).get("access_token", "")
        except Exception:
            pass
    if token:
        return {"Authorization": f"Bearer {token}"}
    return {}


def main():
    try:
        req = urllib.request.Request(URL, method="GET", headers=_auth_header())
        with urllib.request.urlopen(req, timeout=5) as resp:
            items = json.loads(resp.read())
    except Exception as e:
        print(f"ingest_feedback: skip ({e})")
        items = []
    if not items:
        print("ingest_feedback: no items")
        return
    path = OUT / "feedback.csv"
    with path.open("w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=["text", "label", "source", "locale", "content_hash"])
        w.writeheader()
        for item in items:
            w.writerow({
                "text": item.get("body", "")[:200],
                "label": item.get("label", "ham"),
                "source": "user_feedback",
                "locale": item.get("locale", "zh"),
                "content_hash": item.get("id", ""),
            })
    print(f"ingest_feedback: {len(items)} rows -> {path}")


if __name__ == "__main__":
    main()
