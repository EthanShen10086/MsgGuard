#!/usr/bin/env python3
"""Export sklearn pipeline to Core ML .mlmodel (optionally per locale)."""
import argparse
import json
import sys
from pathlib import Path

import joblib

ROOT = Path(__file__).parent.parent
sys.path.insert(0, str(ROOT))
from locale_utils import normalize_locale, output_dir  # noqa: E402


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--locale", default=None)
    args = parser.parse_args()
    tag = normalize_locale(args.locale) if args.locale else None
    out = output_dir(ROOT, tag) if tag else ROOT / "output"
    pipe_path = out / "coreml_pipeline.joblib"
    model_out = out / "spam_classifier.mlmodel"
    if not pipe_path.exists():
        raise SystemExit(f"Run train for locale {tag or 'all'}")
    pipe = joblib.load(pipe_path)
    try:
        from coremltools.converters.sklearn import convert as sk_convert

        model = sk_convert(pipe, input_features="text", output_feature_names="label")
        model.author = "MsgGuard"
        model.short_description = f"SMS spam classifier ({tag or 'all'})"
        model.save(str(model_out))
        meta = {"path": str(model_out), "locale": tag or "all", "size_bytes": model_out.stat().st_size}
        (out / "coreml_export.json").write_text(json.dumps(meta, indent=2))
        print(json.dumps(meta, indent=2))
    except ImportError:
        meta = {"path": str(model_out), "locale": tag or "all", "status": "coremltools not installed"}
        (out / "coreml_export.json").write_text(json.dumps(meta, indent=2))
        print(json.dumps(meta, indent=2))


if __name__ == "__main__":
    main()
