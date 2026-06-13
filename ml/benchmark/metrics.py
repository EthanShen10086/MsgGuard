"""Benchmark metrics helpers."""
from __future__ import annotations

from dataclasses import dataclass
from typing import Iterable


@dataclass
class BenchmarkResult:
    f1: float
    precision: float
    recall: float
    fpr: float
    fnr: float
    total: int
    tp: int
    fp: int
    tn: int
    fn: int


def compute_metrics(y_true: Iterable[str], y_pred: Iterable[str]) -> BenchmarkResult:
    tp = fp = tn = fn = 0
    for t, p in zip(y_true, y_pred):
        t_spam = t in ("spam", "phishing", "promotion")
        p_spam = p in ("spam", "phishing", "promotion")
        if t_spam and p_spam:
            tp += 1
        elif not t_spam and p_spam:
            fp += 1
        elif not t_spam and not p_spam:
            tn += 1
        else:
            fn += 1
    total = tp + fp + tn + fn
    precision = tp / (tp + fp) if (tp + fp) else 0.0
    recall = tp / (tp + fn) if (tp + fn) else 0.0
    f1 = 2 * precision * recall / (precision + recall) if (precision + recall) else 0.0
    fpr = fp / (fp + tn) if (fp + tn) else 0.0
    fnr = fn / (fn + tp) if (fn + tp) else 0.0
    return BenchmarkResult(f1, precision, recall, fpr, fnr, total, tp, fp, tn, fn)
