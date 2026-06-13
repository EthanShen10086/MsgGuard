#!/usr/bin/env python3
"""Publish per-locale trained models to backend model service."""
import hashlib
import json
import os
import shutil
import urllib.request
from pathlib import Path

ROOT = Path(__file__).parent.parent
sys_path = ROOT
import sys

sys.path.insert(0, str(ROOT))
from locale_utils import configured_locales, normalize_locale, output_dir  # noqa: E402

BASE = os.environ.get("MODEL_URL", "http://localhost:8083")
STORAGE = Path(os.environ.get("MODEL_STORAGE_PATH", ROOT.parent / "deploy" / "models"))


def artifact_meta(path: Path):
    if not path.exists():
        return None
    h = hashlib.sha256(path.read_bytes()).hexdigest()
    return {"name": path.name, "checksum": f"sha256:{h}", "size": path.stat().st_size}


def publish_locale(locale: str) -> None:
    tag = normalize_locale(locale)
    out = output_dir(ROOT, tag)
    names = ["bayes_pipeline.joblib", "coreml_pipeline.joblib", "spam_classifier.mlmodel", "bayes_model.json", "coreml_featurizer.json"]
    artifacts = []
    version = "1.0.0"
    core_checksum = "seed"
    for name in names:
        meta = artifact_meta(out / name)
        if meta:
            artifacts.append(meta)
            if name == "spam_classifier.mlmodel":
                core_checksum = meta["checksum"]
    if not artifacts:
        print(f"{tag}: no artifacts in {out}")
        return

    dest = STORAGE / tag / version
    dest.mkdir(parents=True, exist_ok=True)
    for name in names:
        src = out / name
        if src.exists():
            shutil.copy2(src, dest / name)

    payload = {
        "version": version,
        "locale": tag,
        "checksum": core_checksum,
        "artifacts": artifacts,
    }
    url = BASE + "/api/v1/models/register"
    req = urllib.request.Request(
        url, data=json.dumps(payload).encode(),
        headers={"Content-Type": "application/json"}, method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            print(f"{tag}: {resp.read().decode()}")
    except Exception as e:
        print(f"{tag}: saved locally ({e})")
        (out / "publish_manifest.json").write_text(json.dumps(payload, indent=2))


def main():
    STORAGE.mkdir(parents=True, exist_ok=True)
    for locale in configured_locales():
        publish_locale(locale)


if __name__ == "__main__":
    main()
