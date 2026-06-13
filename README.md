# MsgGuard

Hybrid AI spam SMS filter — iOS app + Message Filter Extension + Go backend + Python ML pipeline.

[![CI](https://github.com/EthanShen10086/MsgGuard/actions/workflows/ci.yml/badge.svg)](https://github.com/EthanShen10086/MsgGuard/actions)

## Quick Start (验收)

```bash
# ML pipeline
cd ml && pip install -r requirements.txt && make data && make train && make benchmark

# Backend
cd services/gateway && go build ./cmd/server && CONFIG_PATH=../../deploy/config.yaml go run ./cmd/server &
curl localhost:8080/health   # expect: ok

# iOS
cd apps/ios && bash setup.sh
xcodebuild -scheme MsgGuard-iOS -destination 'platform=iOS Simulator,name=iPhone 16' build

# Docker Tier 1
./deploy/tiers/tier1-compose.sh
```

## Architecture

- [docs/architecture.md](docs/architecture.md)
- [docs/ACCEPTANCE.md](docs/ACCEPTANCE.md) — full checklist
- [docs/cookbook/](docs/cookbook/) — chapters 00–26
- [docs/api/openapi.yaml](docs/api/openapi.yaml)

## Tier Switching

| Tier | Command |
|------|---------|
| 0 Local | `./deploy/tiers/tier0-local.sh` |
| 1 Compose | `./deploy/tiers/tier1-compose.sh` |
| 2 GPU | `./deploy/tiers/tier2-compose-gpu.sh` |
| 3 K8s Staging | `./deploy/tiers/tier3-helm-staging.sh` |
| 4 K8s Prod | `./deploy/tiers/tier4-helm-prod.sh` |

## License

MIT
