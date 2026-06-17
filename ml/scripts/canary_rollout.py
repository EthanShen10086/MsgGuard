#!/usr/bin/env python3
"""Read model registry and set canary percentage via gateway admin flags API."""
import argparse
import json
import os
import sys
import urllib.error
import urllib.request
from pathlib import Path

DEFAULT_REGISTRY = Path(__file__).parent.parent.parent / "deploy" / "models" / "registry.json"
FLAG_KEY = "model_canary"


def load_registry(path: Path) -> dict:
    if not path.is_file():
        print(f"registry not found: {path}", file=sys.stderr)
        return {}
    return json.loads(path.read_text(encoding="utf-8"))


def api_request(base: str, token: str, method: str, path: str, body: dict | None = None) -> dict:
    url = base.rstrip("/") + path
    data = json.dumps(body).encode() if body is not None else None
    req = urllib.request.Request(url, data=data, method=method)
    req.add_header("Content-Type", "application/json")
    req.add_header("Authorization", f"Bearer {token}")
    with urllib.request.urlopen(req, timeout=30) as resp:
        return json.loads(resp.read().decode())


def main() -> None:
    parser = argparse.ArgumentParser(description="Set model canary rollout percentage")
    parser.add_argument("--registry", type=Path, default=DEFAULT_REGISTRY)
    parser.add_argument("--gateway", default=os.environ.get("GATEWAY_URL", "http://localhost:8080"))
    parser.add_argument("--token", default=os.environ.get("ADMIN_TOKEN", ""))
    parser.add_argument("--percentage", type=int, required=True, help="0-100 canary traffic")
    parser.add_argument("--locale", default="", help="Optional locale to log from registry")
    args = parser.parse_args()

    if not args.token:
        print("ADMIN_TOKEN or --token required", file=sys.stderr)
        sys.exit(1)
    if not 0 <= args.percentage <= 100:
        print("percentage must be 0-100", file=sys.stderr)
        sys.exit(1)

    registry = load_registry(args.registry)
    locales = list(registry.keys()) if registry else []
    if args.locale and args.locale in registry:
        meta = registry[args.locale]
        print(f"canary for {args.locale} version={meta.get('version')} -> {args.percentage}%")
    elif locales:
        print(f"registry locales: {locales}; setting flag {FLAG_KEY}={args.percentage}%")
    else:
        print(f"no registry entries; setting flag {FLAG_KEY}={args.percentage}%")

    payload = {
        "key": FLAG_KEY,
        "enabled": args.percentage > 0,
        "percentage": args.percentage,
    }
    try:
        result = api_request(args.gateway, args.token, "POST", "/api/v1/admin/flags", payload)
    except urllib.error.HTTPError as e:
        print(f"admin API error: {e.code} {e.read().decode()}", file=sys.stderr)
        sys.exit(1)
    print(json.dumps(result, indent=2))


if __name__ == "__main__":
    main()
