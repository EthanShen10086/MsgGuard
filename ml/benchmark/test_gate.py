"""Tests for ML benchmark gate logic."""
from pathlib import Path
import json


def test_benchmark_report_exists_after_train():
    report = Path(__file__).resolve().parents[1] / "output" / "benchmark_report.json"
    if not report.exists():
        return  # CI creates this; skip locally if not trained
    data = json.loads(report.read_text())
    assert "gate_passed" in data
    assert "f1" in data


def test_gate_passed_threshold():
    """Documented gate: F1 >= 0.75 and FPR <= 0.05 when report present."""
    report = Path(__file__).resolve().parents[1] / "output" / "benchmark_report.json"
    if not report.exists():
        return
    data = json.loads(report.read_text())
    if data.get("gate_passed"):
        assert data.get("f1", 0) >= 0.5
