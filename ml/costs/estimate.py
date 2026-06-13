#!/usr/bin/env python3
"""Estimate monthly costs for MsgGuard backend."""
import json

CONFIG = {
    "users": 10000,
    "llm_calls_per_user_day": 0.1,
    "llm_cost_per_1k_tokens": 0.002,
    "tokens_per_call": 500,
    "storage_gb": 50,
    "storage_cost_per_gb": 0.023,
    "redis_gb": 1,
    "redis_cost_per_gb": 0.05,
}


def main():
    llm_monthly = (
        CONFIG["users"] * CONFIG["llm_calls_per_user_day"] * 30
        * CONFIG["tokens_per_call"] / 1000 * CONFIG["llm_cost_per_1k_tokens"]
    )
    storage = CONFIG["storage_gb"] * CONFIG["storage_cost_per_gb"]
    redis = CONFIG["redis_gb"] * CONFIG["redis_cost_per_gb"]
    total = llm_monthly + storage + redis
    report = {
        "llm_usd": round(llm_monthly, 2),
        "storage_usd": round(storage, 2),
        "redis_usd": round(redis, 2),
        "total_usd": round(total, 2),
        "config": CONFIG,
    }
    print(json.dumps(report, indent=2))


if __name__ == "__main__":
    main()
