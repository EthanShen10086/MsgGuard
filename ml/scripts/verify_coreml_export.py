#!/usr/bin/env python3
"""Verify Core ML artifacts exist and are non-trivial."""
import json
import sys
from pathlib import Path

ROOT = Path(__file__).parent.parent
sys.path.insert(0, str(ROOT))
from locale_utils import configured_locales  # noqa: E402


def main() -> int:
    failed = False
    for locale in configured_locales():
        out = ROOT / "output" / locale
        model = out / "spam_classifier.mlmodel"
        featurizer = out / "coreml_featurizer.json"
        export_meta = out / "coreml_export.json"
        for path in (model, featurizer, export_meta):
            if not path.exists():
                print(f"MISSING {path}")
                failed = True
                continue
        if model.exists() and model.stat().st_size < 512:
            print(f"TOO_SMALL {model} ({model.stat().st_size} bytes)")
            failed = True
        if featurizer.exists():
            data = json.loads(featurizer.read_text(encoding="utf-8"))
            if data.get("feature_count", 0) < 10:
                print(f"INVALID featurizer feature_count for {locale}")
                failed = True
        if not failed:
            print(f"OK {locale} model={model.stat().st_size}B features={json.loads(featurizer.read_text())['feature_count']}")
    return 1 if failed else 0


if __name__ == "__main__":
    raise SystemExit(main())
