#!/usr/bin/env python3
"""Collect from HuggingFace (fallback to seed)."""
import sys
from pathlib import Path

ROOT = Path(__file__).parent.parent
OUT = ROOT / "data" / "raw" / "hf"
OUT.mkdir(parents=True, exist_ok=True)

try:
    from datasets import load_dataset
    ds = load_dataset("dbarbedillo/SMS_Spam_Multilingual_Collection_Dataset", split="train")
    path = OUT / "multilingual.csv"
    ds.to_csv(str(path))
    print(f"collect_hf: wrote {path}")
except Exception as e:
    print(f"collect_hf: fallback to seed ({e})", file=sys.stderr)
    seed = ROOT / "datasets" / "seed" / "en_spam_ham.csv"
    (OUT / "en_spam_ham.csv").write_bytes(seed.read_bytes())
    print(f"collect_hf: copied seed")
