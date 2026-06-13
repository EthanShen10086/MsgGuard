package httpauth

import (
	"net/http"

	"github.com/EthanShen10086/voxera-kit/auth"
)

// RequirePermission wraps a handler with Bearer token auth and RBAC check.
func RequirePermission(
	authenticator auth.Authenticator,
	authorizer auth.Authorizer,
	resource, action string,
	next http.HandlerFunc,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		claims, err := authenticator.Authenticate(r.Context(), token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ok, _ := authorizer.Authorize(r.Context(), claims, resource, action)
		if !ok {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
