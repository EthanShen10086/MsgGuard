#!/usr/bin/env python3
"""Collect UCI SMS dataset (fallback to seed)."""
import sys
import zipfile
import io
from pathlib import Path

ROOT = Path(__file__).parent.parent
OUT = ROOT / "data" / "raw" / "uci"
OUT.mkdir(parents=True, exist_ok=True)

try:
    import urllib.request
    url = "https://archive.ics.uci.edu/ml/machine-learning-databases/00228/smsspamcollection.zip"
    data = urllib.request.urlopen(url, timeout=30).read()
    with zipfile.ZipFile(io.BytesIO(data)) as zf:
        zf.extractall(OUT)
    print(f"collect_uci: extracted to {OUT}")
except Exception as e:
    print(f"collect_uci: fallback to seed ({e})", file=sys.stderr)
    seed = ROOT / "datasets" / "seed" / "en_spam_ham.csv"
    (OUT / "SMSSpamCollection").write_bytes(seed.read_bytes())
    print("collect_uci: copied seed")
