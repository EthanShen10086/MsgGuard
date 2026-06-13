#!/usr/bin/env python3
"""Dedupe and merge user feedback with processed data."""
import subprocess
import sys
from pathlib import Path

ROOT = Path(__file__).parent.parent


def main():
    merge = ROOT / "pipeline" / "merge.py"
    subprocess.check_call([sys.executable, str(merge)])
    split = ROOT / "pipeline" / "split.py"
    subprocess.check_call([sys.executable, str(split)])
    print("dedupe_and_merge: done")


if __name__ == "__main__":
    main()
