module github.com/EthanShen10086/msgguard/pkg/app

go 1.25.0

require (
	github.com/EthanShen10086/msgguard/pkg/adapters v0.0.0
	github.com/EthanShen10086/msgguard/pkg/config v0.0.0
	github.com/EthanShen10086/msgguard/pkg/ports v0.0.0
	github.com/EthanShen10086/voxera-kit/aiquota v0.0.0
	github.com/EthanShen10086/voxera-kit/audit v0.0.0
	github.com/EthanShen10086/voxera-kit/auth v0.0.0
	github.com/EthanShen10086/voxera-kit/circuitbreaker v0.0.0
	github.com/EthanShen10086/voxera-kit/featureflag v0.0.0
	github.com/EthanShen10086/voxera-kit/llm v0.0.0
	github.com/EthanShen10086/voxera-kit/observability v0.0.0
	github.com/EthanShen10086/voxera-kit/ratelimiter v0.0.0
)

require (
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.28.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/nats-io/nats.go v1.37.0 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/redis/go-redis/v9 v9.7.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver v1.17.1 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.43.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.43.0 // indirect
	go.opentelemetry.io/otel/metric v1.43.0 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.opentelemetry.io/otel/trace v1.43.0 // indirect
	go.opentelemetry.io/proto/otlp v1.10.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260401024825-9d38bb4040a9 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260401024825-9d38bb4040a9 // indirect
	google.golang.org/grpc v1.80.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/EthanShen10086/msgguard/pkg/adapters => ../adapters
	github.com/EthanShen10086/msgguard/pkg/config => ../config
	github.com/EthanShen10086/msgguard/pkg/ports => ../ports
	github.com/EthanShen10086/voxera-kit/aiquota => ../../../voxera-kit/backend/aiquota
	github.com/EthanShen10086/voxera-kit/audit => ../../../voxera-kit/backend/audit
	github.com/EthanShen10086/voxera-kit/auth => ../../../voxera-kit/backend/auth
	github.com/EthanShen10086/voxera-kit/circuitbreaker => ../../../voxera-kit/backend/circuitbreaker
	github.com/EthanShen10086/voxera-kit/featureflag => ../../../voxera-kit/backend/featureflag
	github.com/EthanShen10086/voxera-kit/llm => ../../../voxera-kit/backend/llm
	github.com/EthanShen10086/voxera-kit/loadshed => ../../../voxera-kit/backend/loadshed
	github.com/EthanShen10086/voxera-kit/middleware => ../../../voxera-kit/backend/middleware
	github.com/EthanShen10086/voxera-kit/observability => ../../../voxera-kit/backend/observability
	github.com/EthanShen10086/voxera-kit/pii => ../../../voxera-kit/backend/pii
	github.com/EthanShen10086/voxera-kit/ratelimiter => ../../../voxera-kit/backend/ratelimiter
	github.com/EthanShen10086/voxera-kit/security => ../../../voxera-kit/backend/security
)
