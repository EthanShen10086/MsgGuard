"""Locale helpers for per-locale ML training and iOS/backend alignment."""
from __future__ import annotations

import os
from pathlib import Path

# Data pipeline uses short codes; iOS / backend use BCP-47 tags.
LOCALE_MAP = {
    "zh": "zh-Hans",
    "zh-Hans": "zh-Hans",
    "en": "en-US",
    "en-US": "en-US",
}

DEFAULT_LOCALES = ["zh-Hans", "en-US"]

# Benchmark jsonl → locale tag
BENCH_LOCALE = {
    "test_zh.jsonl": "zh-Hans",
    "test_en.jsonl": "en-US",
}


def normalize_locale(locale: str) -> str:
    return LOCALE_MAP.get(locale.strip(), locale.strip())


def data_locales_for(tag: str) -> list[str]:
    """Return CSV locale column values that belong to a BCP-47 tag."""
    tag = normalize_locale(tag)
    if tag == "zh-Hans":
        return ["zh", "zh-Hans"]
    if tag == "en-US":
        return ["en", "en-US"]
    return [tag]


def output_dir(root: Path, locale: str) -> Path:
    tag = normalize_locale(locale)
    out = root / "output" / tag
    out.mkdir(parents=True, exist_ok=True)
    return out


def configured_locales() -> list[str]:
    raw = os.environ.get("MODEL_LOCALES", ",".join(DEFAULT_LOCALES))
    return [normalize_locale(x) for x in raw.split(",") if x.strip()]
