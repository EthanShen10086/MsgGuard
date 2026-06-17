package memory

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/EthanShen10086/voxera-kit/auth"
)

var rolePermissions = map[string][]string{
	"admin":       {"models:write", "models:read", "feedback:read", "admin:read", "quota:write", "rules:write", "analytics:write"},
	"ml_engineer": {"models:write", "models:read", "feedback:read", "rules:write"},
	"device":      {"models:read", "analytics:write", "feedback:write"},
	"pro":         {"analytics:write", "models:read"},
	"user":        {"analytics:write", "models:read"},
}

type tokenPayload struct {
	UserID string   `json:"uid"`
	Roles  []string `json:"roles"`
	Exp    int64    `json:"exp"`
}

// Auth implements Authenticator + Authorizer with HMAC tokens (dev/production bootstrap).
type Auth struct {
	secret   string
	revoked  map[string]struct{}
	allowDev bool
	mu       sync.RWMutex
}

// NewAuth creates an auth service. allowDevSecret permits the default dev secret when true.
func NewAuth(secret string) *Auth {
	return NewAuthWithOptions(secret, true)
}

// NewAuthWithOptions creates auth with explicit dev-secret policy.
func NewAuthWithOptions(secret string, allowDevSecret bool) *Auth {
	if secret == "" {
		if allowDevSecret {
			secret = "msgguard-dev-secret"
		}
	}
	return &Auth{secret: secret, revoked: map[string]struct{}{}, allowDev: allowDevSecret}
}

// ValidateSecret returns an error if the secret is missing or insecure in production mode.
func ValidateSecret(secret string, env string, allowDev bool) error {
	if secret == "" {
		if allowDev && env != "production" {
			return nil
		}
		return errors.New("AUTH_SECRET is required")
	}
	if env == "production" && (secret == "msgguard-dev-secret" || len(secret) < 32) {
		return errors.New("AUTH_SECRET must be at least 32 chars in production")
	}
	return nil
}

func (a *Auth) Authenticate(ctx context.Context, token string) (*auth.Claims, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token")
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, err
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, []byte(a.secret))
	mac.Write(payloadBytes)
	if !hmac.Equal(sig, mac.Sum(nil)) {
		return nil, errors.New("invalid signature")
	}
	if _, revoked := a.isRevoked(token); revoked {
		return nil, errors.New("token revoked")
	}
	var p tokenPayload
	if err := json.Unmarshal(payloadBytes, &p); err != nil {
		return nil, err
	}
	if p.Exp > 0 && time.Now().Unix() > p.Exp {
		return nil, errors.New("token expired")
	}
	perms := []string{}
	for _, r := range p.Roles {
		perms = append(perms, rolePermissions[r]...)
	}
	return &auth.Claims{
		UserID: p.UserID, Roles: p.Roles, Permissions: perms,
		ExpiresAt: time.Unix(p.Exp, 0),
	}, nil
}

func (a *Auth) GenerateToken(ctx context.Context, claims *auth.Claims) (*auth.TokenPair, error) {
	exp := time.Now().Add(24 * time.Hour).Unix()
	if claims != nil && !claims.ExpiresAt.IsZero() {
		exp = claims.ExpiresAt.Unix()
	}
	uid := "admin"
	roles := []string{"admin"}
	if claims != nil {
		if claims.UserID != "" {
			uid = claims.UserID
		}
		if len(claims.Roles) > 0 {
			roles = claims.Roles
		}
	}
	token, err := a.sign(tokenPayload{UserID: uid, Roles: roles, Exp: exp})
	if err != nil {
		return nil, err
	}
	return &auth.TokenPair{AccessToken: token, TokenType: "Bearer", ExpiresAt: time.Unix(exp, 0)}, nil
}

func (a *Auth) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	claims, err := a.Authenticate(ctx, refreshToken)
	if err != nil {
		return nil, err
	}
	return a.GenerateToken(ctx, &auth.Claims{
		UserID: claims.UserID, Roles: claims.Roles,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
}

func (a *Auth) RevokeToken(ctx context.Context, token string) error {
	token = strings.TrimPrefix(token, "Bearer ")
	a.mu.Lock()
	a.revoked[token] = struct{}{}
	a.mu.Unlock()
	return nil
}

func (a *Auth) isRevoked(token string) (string, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, ok := a.revoked[token]
	return token, ok
}

func (a *Auth) Authorize(ctx context.Context, claims *auth.Claims, resource, action string) (bool, error) {
	return a.HasPermission(ctx, claims, auth.Permission{Resource: resource, Action: action})
}

func (a *Auth) HasPermission(ctx context.Context, claims *auth.Claims, perm auth.Permission) (bool, error) {
	if claims == nil {
		return false, nil
	}
	need := fmt.Sprintf("%s:%s", perm.Resource, perm.Action)
	for _, p := range claims.Permissions {
		if p == need {
			return true, nil
		}
	}
	for _, role := range claims.Roles {
		for _, p := range rolePermissions[role] {
			if p == need {
				return true, nil
			}
		}
	}
	return false, nil
}

func (a *Auth) HasRole(ctx context.Context, claims *auth.Claims, role string) bool {
	if claims == nil {
		return false
	}
	for _, r := range claims.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (a *Auth) sign(p tokenPayload) (string, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, []byte(a.secret))
	mac.Write(b)
	return base64.RawURLEncoding.EncodeToString(b) + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

var (
	_ auth.Authenticator = (*Auth)(nil)
	_ auth.Authorizer    = (*Auth)(nil)
)
