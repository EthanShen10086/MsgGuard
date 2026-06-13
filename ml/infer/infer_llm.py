#!/usr/bin/env python3
"""Call gateway classify API."""
import json
import os
import sys
import urllib.request

URL = os.environ.get("GATEWAY_URL", "http://localhost:8080") + "/api/v1/classify"


def main():
    text = sys.argv[1] if len(sys.argv) > 1 else "Free gift click here"
    req = urllib.request.Request(
        URL, data=json.dumps({"body": text, "sender": "+10000000000"}).encode(),
        headers={"Content-Type": "application/json"}, method="POST",
    )
    with urllib.request.urlopen(req, timeout=10) as resp:
        print(resp.read().decode())


if __name__ == "__main__":
    main()
