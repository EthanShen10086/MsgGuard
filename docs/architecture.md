# MsgGuard Architecture Overview

**Last updated:** 2026-06-17

MsgGuard is a hybrid on-device + optional cloud AI spam filter for SMS (iOS primary), with macOS Mail extension and a Go backend for rules/models, feedback, and admin.

## Target Architecture

```mermaid
flowchart TB
  subgraph clients [Clients]
    iOS[iOS App + Filter Extension]
    macOS[macOS Mail Extension]
    Android[Android Skeleton]
    AdminWeb[Admin Web SPA]
  end

  subgraph edge [Edge / CDN]
    Caddy[Caddy Static Site]
    RulesCDN[Rules CDN]
  end

  subgraph backend [Backend Platform]
    GW[Gateway :8080]
    Rules[Rules :8081]
    Model[Model :8083]
    Flywheel[Flywheel Worker]
  end

  subgraph data [Data]
    PG[(PostgreSQL / MongoDB)]
    Redis[(Redis Cache)]
    NATS[NATS Queue]
  end

  subgraph ml [ML Flywheel]
    Train[Train / Benchmark]
    Export[CoreML Export]
    Retrain[Cron Retrain]
  end

  subgraph obs [Observability]
    Prom[Prometheus]
    Graf[Grafana]
    Jaeger[Jaeger]
  end

  iOS -->|L0-L2 on-device| iOS
  iOS -->|opt-in defer| GW
  iOS -->|rules sync| RulesCDN
  macOS --> GW
  AdminWeb -->|admin API| GW
  GW --> Rules
  GW --> Model
  GW --> PG
  GW --> Redis
  GW --> NATS
  Flywheel --> NATS
  Flywheel --> Train
  Train --> Export
  Retrain --> Train
  GW --> Prom
  Prom --> Graf
  GW --> Jaeger
  Caddy -->|privacy/support/pricing/status| clients
```

## Sub-documents

| Area | Document | Status |
|------|----------|--------|
| iOS client | [client-ios.md](architecture/client-ios.md) | Implemented |
| macOS Mail | [client-macos.md](architecture/client-macos.md) | In Progress |
| Android | [client-android.md](architecture/client-android.md) | Planned |
| Backend platform | [backend-platform.md](architecture/backend-platform.md) | Implemented |
| ML flywheel | [ml-flywheel.md](architecture/ml-flywheel.md) | Implemented |
| Auth & security | [auth-security.md](architecture/auth-security.md) | Implemented (prod secrets at deploy) |
| Commercial / tiers | [commercial.md](architecture/commercial.md) | Implemented (Apple creds at deploy) |
| Admin web | [admin-web.md](architecture/admin-web.md) | Implemented |

## Related

- [Threat model](../security/THREAT_MODEL.md)
- [Tier matrix](../deploy/TIER_MATRIX.md)
- [OpenAPI](../api/openapi.yaml)
- [Commercial readiness](../COMMERCIAL_READINESS.md)
