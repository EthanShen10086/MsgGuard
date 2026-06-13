#!/usr/bin/env python3
"""Generate HTML report from benchmark JSON."""
import json
from pathlib import Path

ROOT = Path(__file__).parent.parent
REPORT = ROOT / "benchmark" / "reports" / "latest" / "report.json"


def main():
    if not REPORT.exists():
        print("Run: make benchmark first")
        return
    data = json.loads(REPORT.read_text())
    print(json.dumps(data, indent=2))


if __name__ == "__main__":
    main()
