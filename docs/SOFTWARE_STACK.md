# MsgGuard Software Stack

## Development (Mac)
| Tool | Version | Purpose |
|------|---------|---------|
| Xcode | 16+ | iOS build |
| XcodeGen | latest | project.yml → xcodeproj |
| SwiftLint | optional | lint |
| Go | 1.25+ | backend services |
| Python | 3.11+ | ML pipeline |
| Docker Desktop | latest | Tier 1 compose |

## ML Training (CPU — no GPU required)
| Package | Purpose |
|---------|---------|
| scikit-learn | Bayes + LogisticRegression |
| pandas | data processing |
| joblib | model serialization |
| coremltools | .mlmodel export (optional) |

## Backend Services
| Component | Default | Alternatives |
|-----------|---------|--------------|
| PostgreSQL 16 | feedback, rules | MongoDB (port adapter) |
| Redis 7 | LLM cache | memory fallback |
| NATS 2 | queue (optional) | noop |
| Jaeger | tracing | — |
| Prometheus + Grafana | metrics | — |
| Caddy | TLS/reverse proxy | Nginx |

## Cloud (Tier 3+)
| Tool | Purpose |
|------|---------|
| kubectl | K8s admin |
| Helm 3 | deploy/helm/msgguard |
| Cloud CLI | S3/R2 for model CDN |

## GPU (Tier 2/4 only)
- NVIDIA runtime + CUDA for large-scale fine-tune
- **Not required** for L0/L1/Bayes MVP
