package httpauth

import (
	"net/http"
	"os"
	"strings"
)

// OIDCConfigured reports whether OIDC middleware is fully wired (stub: never true).
func OIDCConfigured() bool {
	return os.Getenv("OIDC_CLIENT_ID") != "" && os.Getenv("OIDC_CLIENT_SECRET") != ""
}

// OIDCMiddleware is a placeholder for enterprise SSO.
// When OIDC_ISSUER is set but client credentials are missing, protected routes return 501.
func OIDCMiddleware(prefixes ...string) func(http.Handler) http.Handler {
	issuer := strings.TrimSpace(os.Getenv("OIDC_ISSUER"))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if issuer == "" {
				next.ServeHTTP(w, r)
				return
			}
			if !matchesPrefix(r.URL.Path, prefixes) {
				next.ServeHTTP(w, r)
				return
			}
			if !OIDCConfigured() {
				http.Error(w, "OIDC issuer configured but OIDC middleware not implemented", http.StatusNotImplemented)
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
