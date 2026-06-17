package httpauth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/EthanShen10086/voxera-kit/auth"
)

// DeviceTokenIssuer mints short-lived device-scoped tokens.
type DeviceTokenIssuer interface {
	IssueDeviceToken(ctx context.Context, deviceID string) (*auth.TokenPair, error)
}

// RequireDeviceOrBearer accepts a valid Bearer token with device/user/pro/admin role.
func RequireDeviceOrBearer(authenticator auth.Authenticator, next http.HandlerFunc) http.HandlerFunc {
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
		for _, role := range claims.Roles {
			if role == "device" || role == "admin" || role == "user" || role == "pro" {
				next(w, r)
				return
			}
		}
		http.Error(w, "forbidden", http.StatusForbidden)
	}
}

// RequireModelRead allows model download for roles with models:read permission.
func RequireModelRead(authenticator auth.Authenticator, authorizer auth.Authorizer, next http.HandlerFunc) http.HandlerFunc {
	return RequirePermission(authenticator, authorizer, "models", "read", next)
}

// IssueDeviceTokenHandler exposes POST /api/v1/auth/device for app bootstrap.
func IssueDeviceTokenHandler(issuer DeviceTokenIssuer, bootstrapEnabled bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if !bootstrapEnabled {
			http.Error(w, "device token issuance disabled", http.StatusForbidden)
			return
		}
		var req struct {
			DeviceID string `json:"device_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.DeviceID) == "" {
			http.Error(w, "device_id required", http.StatusBadRequest)
			return
		}
		pair, err := issuer.IssueDeviceToken(r.Context(), req.DeviceID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": pair.AccessToken,
			"token_type":   pair.TokenType,
			"expires_at":   pair.ExpiresAt,
		})
	}
}

// MemoryDeviceIssuer issues device-role tokens via Authenticator.
type MemoryDeviceIssuer struct {
	Auth auth.Authenticator
	TTL  time.Duration
}

func (m *MemoryDeviceIssuer) IssueDeviceToken(ctx context.Context, deviceID string) (*auth.TokenPair, error) {
	ttl := m.TTL
	if ttl == 0 {
		ttl = 30 * 24 * time.Hour
	}
	return m.Auth.GenerateToken(ctx, &auth.Claims{
		UserID: deviceID,
		Roles:  []string{"device"},
		ExpiresAt: time.Now().Add(ttl),
	})
}
