module github.com/EthanShen10086/msgguard/pkg/httpauth

go 1.25.0

require (
	github.com/EthanShen10086/voxera-kit/auth v0.0.0
	github.com/coreos/go-oidc/v3 v3.18.0
	golang.org/x/oauth2 v0.36.0
)

require github.com/go-jose/go-jose/v4 v4.1.4 // indirect

replace github.com/EthanShen10086/voxera-kit/auth => ../../../voxera-kit/backend/auth
