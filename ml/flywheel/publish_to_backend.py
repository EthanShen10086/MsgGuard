#!/usr/bin/env python3
"""Publish trained model to backend model service."""
import hashlib
import json
import os
import urllib.request
from pathlib import Path

ROOT = Path(__file__).parent.parent
OUTPUT = ROOT / "output"
URL = os.environ.get("MODEL_URL", "http://localhost:8083") + "/api/v1/models/register"


def main():
    artifacts = []
    for name in ["bayes_pipeline.joblib", "coreml_pipeline.joblib", "spam_classifier.mlmodel"]:
        path = OUTPUT / name
        if path.exists():
            h = hashlib.sha256(path.read_bytes()).hexdigest()
            artifacts.append({"name": name, "checksum": f"sha256:{h}", "size": path.stat().st_size})
    payload = {"version": "1.0.0", "locale": "zh-Hans", "artifacts": artifacts}
    req = urllib.request.Request(
        URL, data=json.dumps(payload).encode(),
        headers={"Content-Type": "application/json"}, method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            print(resp.read().decode())
    except Exception as e:
        print(f"publish_to_backend: saved locally ({e})")
        (OUTPUT / "publish_manifest.json").write_text(json.dumps(payload, indent=2))


if __name__ == "__main__":
    main()
