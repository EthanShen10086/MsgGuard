package httpauth

import (
	"net/http"
	"os"
	"strings"
)

// OIDCConfigured reports whether OIDC client credentials are present.
func OIDCConfigured() bool {
	return strings.TrimSpace(os.Getenv("OIDC_ISSUER")) != "" &&
		strings.TrimSpace(os.Getenv("OIDC_CLIENT_ID")) != "" &&
		os.Getenv("OIDC_CLIENT_SECRET") != ""
}

// OIDCEnforceAdmin blocks bootstrap token on admin routes when OIDC_ENFORCE_ADMIN=true.
func OIDCEnforceAdmin() bool {
	return os.Getenv("OIDC_ENFORCE_ADMIN") == "true"
}

// OIDCMiddleware optionally requires OIDC to be configured before serving admin prefixes.
func OIDCMiddleware(provider *OIDCProvider, prefixes ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if provider == nil || !provider.EnforceAdmin() {
				next.ServeHTTP(w, r)
				return
			}
			if !matchesPrefix(r.URL.Path, prefixes) {
				next.ServeHTTP(w, r)
				return
			}
			if !provider.Enabled() {
				http.Error(w, "OIDC required for admin but not configured", http.StatusServiceUnavailable)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func matchesPrefix(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
