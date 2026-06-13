#!/usr/bin/env python3
"""Generate Cookbook chapters 00-26 with acceptance commands."""
from pathlib import Path

ROOT = Path(__file__).parent.parent
COOKBOOK = ROOT / "docs" / "cookbook"

CHAPTERS = [
    ("00-setup", "环境搭建", "apps/ios/setup.sh, ml/requirements.txt", "cd apps/ios && bash setup.sh\ncd ml && pip install -r requirements.txt", "setup completes without error"),
    ("01-extension", "Message Filter Extension", "apps/ios/MessageFilterExtension/MessageFilterExtension.swift", "cd apps/ios && xcodebuild -scheme MsgGuard-iOS -destination 'platform=iOS Simulator,name=iPhone 16' build", "BUILD SUCCEEDED"),
    ("02-bayes", "贝叶斯 L1", "ml/train/train_bayes.py", "cd ml && make train", "f1 > 0.5 in output"),
    ("03-coreml", "Core ML L2", "ml/export/export_coreml.py", "cd ml && make export", "coreml_export.json created"),
    ("04-backend", "Go 后端 Gateway", "services/gateway/cmd/server/main.go", "cd services/gateway && go build ./cmd/server", "build succeeds"),
    ("05-llm", "云端 LLM L3", "services/gateway/internal/handler/classify.go", "curl -X POST localhost:8080/api/v1/classify -d '{\"body\":\"free gift\"}'", "JSON action field"),
    ("06-decision-tree", "混合决策树", "apps/ios/Packages/FilterEngine/Sources/FilterEngine/HybridFilterEngine.swift", "cd apps/ios && swift test --package-path Packages/FilterEngine", "tests pass"),
    ("07-appstore", "App Store 上架", "docs/app-store/metadata.md", "open docs/app-store/metadata.md", "metadata checklist present"),
    ("08-maintenance", "日常维护", "ml/flywheel/schedule_retrain.sh", "bash ml/flywheel/schedule_retrain.sh", "retrain pipeline runs"),
    ("09-traceid", "TraceID 追溯", "services/gateway/cmd/server/main.go", "curl -v localhost:8080/health 2>&1 | grep -i x-request", "X-Request-ID header"),
    ("10-config", "配置切换", "deploy/config.yaml", "cat deploy/config.yaml", "database/cache sections present"),
    ("11-data-pipeline", "数据流水线总览", "ml/Makefile", "cd ml && make help", "targets listed"),
    ("12-deploy-compose", "单机 Docker 部署", "deploy/docker-compose.yml", "./deploy/tiers/tier1-compose.sh", "curl localhost:8080/health -> ok"),
    ("13-deploy-k8s", "K8s Helm 部署", "deploy/helm/msgguard/", "helm template msgguard deploy/helm/msgguard", "Deployment manifest rendered"),
    ("14-gpu-train", "GPU 训练", "deploy/k8s/gpu-training-job.yaml", "kubectl apply -f deploy/k8s/gpu-training-job.yaml --dry-run=client", "job valid"),
    ("15-observability", "可观测 Debug", "deploy/prometheus/prometheus.yml", "open http://localhost:16686", "Jaeger UI"),
    ("16-pluggable", "可插拔切换", "pkg/ports/", "grep driver deploy/config.yaml", "postgres/memory drivers"),
    ("17-fallback", "兜底降级", "services/gateway/internal/handler/classify.go", "unset QWEN_API_KEY && curl classify", "heuristic fallback"),
    ("18-data-collection", "数据采集", "ml/pipeline/collect_uci.py", "cd ml && python3 pipeline/collect_uci.py", "data/raw/ populated"),
    ("19-clean-label", "清洗标注", "ml/pipeline/clean.py", "cd ml && python3 pipeline/clean.py", "data/processed/all.csv"),
    ("20-flywheel", "数据飞轮", "ml/flywheel/ingest_feedback.py", "cd ml && make flywheel", "merge completes"),
    ("21-infer", "推理脚本", "ml/infer/infer_bayes.py", "cd ml && make infer TEXT='免费贷款'", "JSON label spam"),
    ("22-benchmark", "Benchmark 成本", "ml/benchmark/run_benchmark.py", "cd ml && make benchmark", "gate_passed=True"),
    ("23-ai-sync", "AI 后端客户端协同", "docs/api/openapi.yaml", "curl localhost:8083/api/v1/models/latest", "version JSON"),
    ("24-eval", "效果评测体系", "ml/benchmark/baselines.yaml", "cat ml/benchmark/reports/latest/report.json", "f1 and fpr fields"),
    ("25-tier-switch", "分级部署切换", "deploy/tiers/", "ls deploy/tiers/", "tier0-4 scripts"),
    ("26-runbook", "运维 Runbook", "docs/runbooks/", "cat docs/runbooks/llm-outage.md", "SOP steps"),
]

TEMPLATE = """# {title}

## 对应原始需求
{req}

## 涉及文件
{files}

## 架构图
见 [architecture.md](../architecture.md)

## 动手验收
```bash
{cmd}
```
**期望输出：** {expect}

## Debug 指南
- 查 TraceID：`curl -v` 响应头 `X-Request-ID`
- Gateway 日志：docker compose logs gateway
- iOS 日志：Console.app 过滤 `com.msgguard`

## 扩展阅读
- [ACCEPTANCE.md](../ACCEPTANCE.md)
- [SOFTWARE_STACK.md](../SOFTWARE_STACK.md)
"""

for slug, title, files, cmd, expect in CHAPTERS:
    path = COOKBOOK / f"{slug}.md"
    # Keep existing 00-10 names where different
    name_map = {
        "11-data-pipeline": "11-data-pipeline.md",
        "12-deploy-compose": "12-deploy-compose.md",
    }
    content = TEMPLATE.format(title=title, req=title, files=files, cmd=cmd, expect=expect)
    path.write_text(content, encoding="utf-8")
    print(f"wrote {path.name}")

# Rename old 11 if needed - we have new chapters with new names
print("done")
