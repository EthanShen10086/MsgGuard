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

## CI/CD & Release

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `ci.yml` | push/PR `main`, `staging` | Build & test gate (Swift, Go, Helm, ML) |
| `deploy-staging.yml` | push `staging` | Helm validate, image build, optional K8s deploy |
| `release.yml` | tag `v*.*.*` | ML gate → GHCR → canary → prod |
| `ios-beta.yml` | manual | TestFlight (requires ASC secrets) |
| `security.yml` | weekly + PR | gosec, npm audit |

See [RELEASE.md](RELEASE.md) and [CONTRIBUTING.md](../CONTRIBUTING.md) for branch strategy.

## Cloud (Tier 3+)
| Tool | Purpose |
|------|---------|
| kubectl | K8s admin |
| Helm 3 | deploy/helm/msgguard |
| Cloud CLI | S3/R2 for model CDN |

## GPU (Tier 2/4 only)
- NVIDIA runtime + CUDA for large-scale fine-tune
- **Not required** for L0/L1/Bayes MVP
