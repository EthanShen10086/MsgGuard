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


def main():
    try:
        req = urllib.request.Request(URL, method="GET")
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
