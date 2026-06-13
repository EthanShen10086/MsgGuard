#!/usr/bin/env python3
"""Build Core ML spam_classifier.mlmodel from sklearn TF-IDF + LogisticRegression."""
from __future__ import annotations

import json
import os
import sys
from pathlib import Path
from typing import Optional

import joblib
import numpy as np

ROOT = Path(__file__).parent.parent
sys.path.insert(0, str(ROOT))
from locale_utils import normalize_locale, output_dir  # noqa: E402


def _extract_pipeline(pipe):
    tfidf = pipe.named_steps["tfidf"]
    clf = pipe.named_steps["clf"]
    vocab = tfidf.vocabulary_
    idf = getattr(tfidf, "idf_", np.ones(len(vocab)))
    feature_names = [None] * len(vocab)
    for term, idx in vocab.items():
        feature_names[idx] = term
    coef = clf.coef_
    intercept = clf.intercept_
    if coef.shape[0] == 1:
        weights = coef[0]
        bias = float(intercept[0])
    else:
        weights = coef[1] - coef[0]
        bias = float(intercept[1] - intercept[0])
    return feature_names, idf, weights, bias


def export_locale(locale: Optional[str] = None) -> dict:
    from coremltools.models import MLModel, datatypes
    from coremltools.models.neural_network import NeuralNetworkBuilder

    tag = normalize_locale(locale) if locale else None
    out = output_dir(ROOT, tag) if tag else ROOT / "output"
    pipe_path = out / "coreml_pipeline.joblib"
    model_out = out / "spam_classifier.mlmodel"
    featurizer_out = out / "coreml_featurizer.json"
    if not pipe_path.exists():
        raise SystemExit(f"Missing pipeline at {pipe_path}; run make train-locale LOCALE={locale}")

    pipe = joblib.load(pipe_path)
    feature_names, idf, weights, bias = _extract_pipeline(pipe)
    feature_count = len(feature_names)

    featurizer = {
        "version": 1,
        "locale": tag or "all",
        "feature_count": feature_count,
        "vocabulary": {name: idx for idx, name in enumerate(feature_names) if name},
        "idf": idf.tolist(),
        "threshold": 0.5,
    }
    featurizer_out.write_text(json.dumps(featurizer, ensure_ascii=False), encoding="utf-8")

    builder = NeuralNetworkBuilder(
        input_features=[("features", datatypes.Array(feature_count))],
        output_features=[("spam_score", datatypes.Array(1))],
    )
    builder.add_inner_product(
        name="logit",
        W=weights.reshape(1, feature_count).astype(np.float32),
        b=np.array([bias], dtype=np.float32),
        input_channels=feature_count,
        output_channels=1,
        has_bias=True,
        input_name="features",
        output_name="logit",
    )
    builder.add_activation(
        name="spam_score",
        non_linearity="SIGMOID",
        input_name="logit",
        output_name="spam_score",
    )
    spec = builder.spec
    spec.description.metadata.shortDescription = f"MsgGuard logistic spam model ({tag or 'all'})"
    spec.description.metadata.author = "MsgGuard"
    model = MLModel(spec)
    model.save(str(model_out))

    size = model_out.stat().st_size
    meta = {
        "path": str(model_out),
        "locale": tag or "all",
        "featurizer": str(featurizer_out),
        "size_bytes": size,
        "feature_count": feature_count,
        "status": "exported",
    }
    (out / "coreml_export.json").write_text(json.dumps(meta, indent=2))
    if os.environ.get("CI") == "true" and size < 512:
        raise SystemExit(f"CoreML export too small: {size} bytes")
    if not featurizer_out.exists():
        raise SystemExit("Featurizer JSON missing after export")
    return meta


def main():
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("--locale", default=None)
    args = parser.parse_args()
    meta = export_locale(args.locale)
    print(json.dumps(meta, indent=2))


if __name__ == "__main__":
    main()
