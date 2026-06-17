package httpauth

import (
	"net/http"
	"os"
	"strings"
)

// AppAttestMiddleware optionally requires X-App-Attest-Token on device-facing POST routes.
// Production: set APP_ATTEST_REQUIRED=true after wiring Apple App Attest verification.
func AppAttestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	required := strings.ToLower(os.Getenv("APP_ATTEST_REQUIRED")) == "true"
	return func(w http.ResponseWriter, r *http.Request) {
		if !required {
			next(w, r)
			return
		}
		if strings.TrimSpace(r.Header.Get("X-App-Attest-Token")) == "" {
			http.Error(w, "app attest required", http.StatusUnauthorized)
			return
		}
		// TODO: verify attestation with Apple when APP_ATTEST_TEAM_ID is configured
		next(w, r)
	}
}

// ChainDeviceAuth wraps device-token routes with optional App Attest.
func ChainDeviceAuth(next http.HandlerFunc) http.HandlerFunc {
	return AppAttestMiddleware(next)
}
