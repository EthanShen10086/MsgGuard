#!/usr/bin/env python3
"""Publish trained models for all supported locales."""
import hashlib
import json
import os
import urllib.request
from pathlib import Path

ROOT = Path(__file__).parent.parent
OUTPUT = ROOT / "output"
BASE = os.environ.get("MODEL_URL", "http://localhost:8083")
LOCALES = os.environ.get("MODEL_LOCALES", "zh-Hans,en-US").split(",")


def publish_locale(locale: str) -> None:
    artifacts = []
    for name in ["bayes_pipeline.joblib", "coreml_pipeline.joblib", "spam_classifier.mlmodel"]:
        path = OUTPUT / name
        if path.exists():
            h = hashlib.sha256(path.read_bytes()).hexdigest()
            artifacts.append({"name": name, "checksum": f"sha256:{h}", "size": path.stat().st_size})
    payload = {"version": "1.0.0", "locale": locale.strip(), "artifacts": artifacts}
    url = BASE + "/api/v1/models/register"
    req = urllib.request.Request(
        url, data=json.dumps(payload).encode(),
        headers={"Content-Type": "application/json"}, method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            print(f"{locale}: {resp.read().decode()}")
    except Exception as e:
        print(f"{locale}: saved locally ({e})")
        (OUTPUT / f"publish_manifest_{locale.strip()}.json").write_text(json.dumps(payload, indent=2))


def main():
    for locale in LOCALES:
        if locale.strip():
            publish_locale(locale)


if __name__ == "__main__":
    main()
