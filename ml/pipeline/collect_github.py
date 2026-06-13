#!/usr/bin/env python3
"""Collect datasets from GitHub (fallback to seed)."""
import csv
import sys
from pathlib import Path

ROOT = Path(__file__).parent.parent
OUT = ROOT / "data" / "raw" / "github"
OUT.mkdir(parents=True, exist_ok=True)

try:
    import urllib.request
    url = "https://raw.githubusercontent.com/Cypher-Z/FBS_SMS_Dataset/main/data.csv"
    dest = OUT / "fbs_sms.csv"
    urllib.request.urlretrieve(url, dest)
    print(f"collect_github: downloaded {dest}")
except Exception as e:
    print(f"collect_github: fallback to seed ({e})", file=sys.stderr)
    seed = ROOT / "datasets" / "seed"
    for f in seed.glob("*.csv"):
        (OUT / f.name).write_bytes(f.read_bytes())
    print(f"collect_github: copied seed to {OUT}")
