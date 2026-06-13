# Product Flywheel

## 对应原始需求
埋点 → 指标 → 洞察 → 迭代规划

## 涉及文件
- `docs/product/METRICS.md`, `ITERATION_SOP.md`, `ROADMAP.md`
- `ml/product/aggregate_metrics.py`
- `deploy/grafana/dashboards/product.json`
- `GET /api/v1/admin/metrics/summary`

## 动手验收
```bash
cd ml && python product/aggregate_metrics.py
cat product/reports/weekly.json
```
**期望输出：** JSON with admin_summary, benchmark_overall, recommendations

## Debug 指南
- aggregate 失败 → 先启动 gateway，POST /api/v1/auth/token

## 扩展阅读
- `docs/product/ITERATION_SOP.md`
