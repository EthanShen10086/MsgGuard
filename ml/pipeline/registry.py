#!/usr/bin/env python3
"""Register dataset version manifest."""
import hashlib
import json
from pathlib import Path

ROOT = Path(__file__).parent.parent
REG = ROOT / "data" / "registry"
REG.mkdir(parents=True, exist_ok=True)


def main():
    all_csv = ROOT / "data" / "processed" / "all.csv"
    h = hashlib.sha256(all_csv.read_bytes()).hexdigest() if all_csv.exists() else "empty"
    manifest_path = REG / "manifest.json"
    manifest = json.loads(manifest_path.read_text()) if manifest_path.exists() else {}
    manifest["data_hash"] = h
    manifest_path.write_text(json.dumps(manifest, indent=2))
    print(f"registry: hash={h[:12]}...")


if __name__ == "__main__":
    main()
