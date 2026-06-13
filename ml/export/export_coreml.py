#!/usr/bin/env python3
"""Export sklearn pipeline to Core ML .mlmodel."""
import json
from pathlib import Path

import joblib

ROOT = Path(__file__).parent.parent
OUTPUT = ROOT / "output"
PIPE = OUTPUT / "coreml_pipeline.joblib"
MODEL_OUT = OUTPUT / "spam_classifier.mlmodel"


def main():
    if not PIPE.exists():
        raise SystemExit("Run: make train")
    pipe = joblib.load(PIPE)
    try:
        import coremltools as ct
        from coremltools.converters.sklearn import convert as sk_convert

        model = sk_convert(
            pipe,
            input_features="text",
            output_feature_names="label",
        )
        model.author = "MsgGuard"
        model.short_description = "SMS spam classifier"
        model.save(str(MODEL_OUT))
        meta = {"path": str(MODEL_OUT), "size_bytes": MODEL_OUT.stat().st_size}
        (OUTPUT / "coreml_export.json").write_text(json.dumps(meta, indent=2))
        print(json.dumps(meta, indent=2))
    except ImportError:
        # Fallback: write metadata only
        meta = {"path": str(MODEL_OUT), "status": "coremltools not installed, pipeline saved as joblib"}
        (OUTPUT / "coreml_export.json").write_text(json.dumps(meta, indent=2))
        print(json.dumps(meta, indent=2))


if __name__ == "__main__":
    main()
