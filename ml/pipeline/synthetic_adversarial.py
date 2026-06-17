#!/usr/bin/env python3
"""Generate adversarial SMS samples that mimic ham (OTP, delivery) with spam intent."""
import argparse
import csv
import json
import random
from pathlib import Path

ROOT = Path(__file__).parent.parent
OUT_DIR = ROOT / "datasets" / "adversarial"

TEMPLATES_ZH = [
    ("您的验证码是{code}，5分钟内有效", "ham"),
    ("【顺丰快递】包裹已到达，取件码{code}", "ham"),
    ("【银行】您尾号{tail}消费{amt}元", "ham"),
    ("免费验证码领取点击{url}", "spam"),
    ("您的验证码中奖请点击{url}", "spam"),
    ("快递丢失请补填信息{url}", "spam"),
    ("恭喜验证码中奖加微信{wx}", "spam"),
    ("登录验证码{code}请勿转发给他人", "ham"),
]

TEMPLATES_EN = [
    ("Your verification code is {code}", "ham"),
    ("Amazon: Your order has shipped", "ham"),
    ("Verify account click here {url}", "spam"),
    ("Your package held pay fee {url}", "spam"),
    ("OTP {code} valid 5 min do not share", "ham"),
    ("Free gift verify now {url}", "spam"),
    ("WINNER claim prize {url}", "spam"),
]

URLS = ["http://bit.ly/x", "https://t.cn/abc", "http://fake-bank.com/verify"]
CODES = ["123456", "847291", "8888", "999888"]


def render(template: str) -> str:
    return template.format(
        code=random.choice(CODES),
        tail=random.randint(1000, 9999),
        amt=random.randint(10, 9999),
        url=random.choice(URLS),
        wx="wxid_" + str(random.randint(10000, 99999)),
    )


def generate(locale: str, count: int) -> list[dict]:
    templates = TEMPLATES_ZH if locale.startswith("zh") else TEMPLATES_EN
    rows = []
    for _ in range(count):
        text_tpl, label = random.choice(templates)
        rows.append({"text": render(text_tpl), "label": label, "locale": locale, "adversarial": True})
    return rows


def write_outputs(rows: list[dict], out: Path) -> None:
    out.parent.mkdir(parents=True, exist_ok=True)
    with out.with_suffix(".csv").open("w", newline="", encoding="utf-8") as f:
        w = csv.DictWriter(f, fieldnames=["text", "label", "locale", "adversarial"])
        w.writeheader()
        w.writerows(rows)
    with out.with_suffix(".jsonl").open("w", encoding="utf-8") as f:
        for row in rows:
            f.write(json.dumps(row, ensure_ascii=False) + "\n")


def main() -> None:
    parser = argparse.ArgumentParser(description="Generate adversarial SMS benchmark samples")
    parser.add_argument("--locale", default="zh-Hans", help="Locale tag (zh-Hans, en-US)")
    parser.add_argument("--count", type=int, default=200, help="Number of samples")
    parser.add_argument("--output", type=Path, default=OUT_DIR / "synthetic")
    args = parser.parse_args()

    rows = generate(args.locale, args.count)
    write_outputs(rows, args.output)
    spam = sum(1 for r in rows if r["label"] == "spam")
    print(f"Wrote {len(rows)} adversarial samples ({spam} spam) -> {args.output}.csv/jsonl")


if __name__ == "__main__":
    main()
