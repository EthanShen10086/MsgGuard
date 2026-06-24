module github.com/EthanShen10086/msgguard/pkg/adapters

go 1.25.0

require (
	github.com/EthanShen10086/msgguard/pkg/ports v0.0.0
	github.com/EthanShen10086/voxera-kit/aiquota v0.0.0
	github.com/EthanShen10086/voxera-kit/auth v0.0.0
	github.com/EthanShen10086/voxera-kit/featureflag v0.0.0
	github.com/EthanShen10086/voxera-kit/llm v0.0.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/lib/pq v1.10.9
	github.com/nats-io/nats.go v1.37.0
	github.com/redis/go-redis/v9 v9.7.0
	go.mongodb.org/mongo-driver v1.17.9
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/nats-io/nkeys v0.4.7 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
)

replace (
	github.com/EthanShen10086/msgguard/pkg/ports => ../ports
	github.com/EthanShen10086/voxera-kit/aiquota => ../../../voxera-kit/backend/aiquota
	github.com/EthanShen10086/voxera-kit/auth => ../../../voxera-kit/backend/auth
	github.com/EthanShen10086/voxera-kit/featureflag => ../../../voxera-kit/backend/featureflag
	github.com/EthanShen10086/voxera-kit/llm => ../../../voxera-kit/backend/llm
)
