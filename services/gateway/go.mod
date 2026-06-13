module github.com/EthanShen10086/msgguard/services/gateway

go 1.25.0

require (
	github.com/EthanShen10086/msgguard/pkg/app v0.0.0
	github.com/EthanShen10086/msgguard/pkg/config v0.0.0
	github.com/EthanShen10086/msgguard/pkg/httpauth v0.0.0
	github.com/EthanShen10086/msgguard/pkg/ports v0.0.0
	github.com/EthanShen10086/voxera-kit/audit v0.0.0
	github.com/EthanShen10086/voxera-kit/auth v0.0.0
	github.com/EthanShen10086/voxera-kit/circuitbreaker v0.0.0
	github.com/EthanShen10086/voxera-kit/loadshed v0.0.0
	github.com/EthanShen10086/voxera-kit/middleware v0.0.0
	github.com/EthanShen10086/voxera-kit/observability v0.0.0
	github.com/EthanShen10086/voxera-kit/pii v0.0.0
	github.com/EthanShen10086/voxera-kit/security v0.0.0
	github.com/google/uuid v1.6.0
)

replace (
	github.com/EthanShen10086/msgguard/pkg/app => ../../pkg/app
	github.com/EthanShen10086/msgguard/pkg/config => ../../pkg/config
	github.com/EthanShen10086/msgguard/pkg/httpauth => ../../pkg/httpauth
	github.com/EthanShen10086/msgguard/pkg/ports => ../../pkg/ports
	github.com/EthanShen10086/voxera-kit/audit => ../../../voxera-kit/backend/audit
	github.com/EthanShen10086/voxera-kit/auth => ../../../voxera-kit/backend/auth
	github.com/EthanShen10086/voxera-kit/circuitbreaker => ../../../voxera-kit/backend/circuitbreaker
	github.com/EthanShen10086/voxera-kit/loadshed => ../../../voxera-kit/backend/loadshed
	github.com/EthanShen10086/voxera-kit/middleware => ../../../voxera-kit/backend/middleware
	github.com/EthanShen10086/voxera-kit/observability => ../../../voxera-kit/backend/observability
	github.com/EthanShen10086/voxera-kit/pii => ../../../voxera-kit/backend/pii
	github.com/EthanShen10086/voxera-kit/security => ../../../voxera-kit/backend/security
)
