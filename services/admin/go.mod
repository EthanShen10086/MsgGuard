module github.com/EthanShen10086/msgguard/services/admin

go 1.25.0

require (
	github.com/EthanShen10086/msgguard/pkg/adapters v0.0.0
	github.com/EthanShen10086/voxera-kit/aiquota v0.0.0
	github.com/EthanShen10086/voxera-kit/featureflag v0.0.0
)

replace (
	github.com/EthanShen10086/msgguard/pkg/adapters => ../../pkg/adapters
	github.com/EthanShen10086/voxera-kit/aiquota => ../../../voxera-kit/backend/aiquota
	github.com/EthanShen10086/voxera-kit/featureflag => ../../../voxera-kit/backend/featureflag
)
