package httpauth

import (
	"net/http"
	"strings"
)

// RequireClientCert rejects requests to matching path prefixes when no TLS client
// certificate was presented (gateway terminates mTLS or receives TLS from upstream).
func RequireClientCert(pathPrefixes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !pathMatchesPrefix(r.URL.Path, pathPrefixes) {
				next.ServeHTTP(w, r)
				return
			}
			if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
				http.Error(w, "client certificate required", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireClientCertHeader validates mTLS at the edge when a proxy forwards the
// client subject DN (e.g. Caddy `{http.request.tls.client.subject}`).
func RequireClientCertHeader(header string, pathPrefixes []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !pathMatchesPrefix(r.URL.Path, pathPrefixes) {
				next.ServeHTTP(w, r)
				return
			}
			if strings.TrimSpace(r.Header.Get(header)) == "" {
				http.Error(w, "client certificate required", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func pathMatchesPrefix(path string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}
