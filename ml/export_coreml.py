#!/usr/bin/env python3
"""Export sklearn pipeline to Core ML (when coremltools supports pipeline)."""

import json
from pathlib import Path

OUTPUT = Path(__file__).parent / "output"


def main():
    pipeline_path = OUTPUT / "coreml_pipeline.joblib"
    if not pipeline_path.exists():
        print("Run train_coreml.py first")
        return
    try:
        import coremltools as ct
        import joblib
        pipe = joblib.load(pipeline_path)
        # Placeholder: full sklearn→CoreML conversion requires custom wrapper for TF-IDF+LR
        meta = {
            "status": "placeholder",
            "note": "Integrate SKLearn converter or train natively with Create ML for production",
            "pipeline_path": str(pipeline_path),
        }
        (OUTPUT / "export_status.json").write_text(json.dumps(meta, indent=2))
        print(json.dumps(meta, indent=2))
    except ImportError:
        print("coremltools not installed")


if __name__ == "__main__":
    main()
