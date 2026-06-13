module github.com/EthanShen10086/msgguard/pkg/adapters

go 1.25.0

require (
	github.com/EthanShen10086/msgguard/pkg/ports v0.0.0
	github.com/EthanShen10086/voxera-kit/auth v0.0.0
	github.com/EthanShen10086/voxera-kit/llm v0.0.0
	github.com/lib/pq v1.10.9
	github.com/nats-io/nats.go v1.37.0
	github.com/redis/go-redis/v9 v9.7.0
)

replace (
	github.com/EthanShen10086/msgguard/pkg/ports => ../ports
	github.com/EthanShen10086/voxera-kit/auth => ../../../voxera-kit/backend/auth
	github.com/EthanShen10086/voxera-kit/llm => ../../../voxera-kit/backend/llm
)
