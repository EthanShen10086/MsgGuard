module github.com/EthanShen10086/msgguard/services/model

go 1.25.0

require (
	github.com/EthanShen10086/msgguard/pkg/adapters v0.0.0
	github.com/EthanShen10086/msgguard/pkg/ports v0.0.0
)

replace (
	github.com/EthanShen10086/msgguard/pkg/adapters => ../../pkg/adapters
	github.com/EthanShen10086/msgguard/pkg/ports => ../../pkg/ports
)
