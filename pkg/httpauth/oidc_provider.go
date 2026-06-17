package httpauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/EthanShen10086/voxera-kit/auth"
	"golang.org/x/oauth2"
)

const oidcStateCookie = "msgguard_oidc_state"

// OIDCProvider implements authorization-code SSO for the admin console.
type OIDCProvider struct {
	enabled   bool
	oauth     *oauth2.Config
	verifier  *oidc.IDTokenVerifier
	adminUI   string
	admins    map[string]struct{}
	adminDom  string
	auth      auth.Authenticator
	enforce   bool
}

// NewOIDCProviderFromEnv wires Google/Okta/Azure AD style OIDC when credentials are set.
func NewOIDCProviderFromEnv(authenticator auth.Authenticator) (*OIDCProvider, error) {
	issuer := strings.TrimSpace(os.Getenv("OIDC_ISSUER"))
	clientID := strings.TrimSpace(os.Getenv("OIDC_CLIENT_ID"))
	secret := os.Getenv("OIDC_CLIENT_SECRET")
	if issuer == "" || clientID == "" || secret == "" {
		return &OIDCProvider{enabled: false, auth: authenticator}, nil
	}
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("oidc provider: %w", err)
	}
	redirect := strings.TrimSpace(os.Getenv("OIDC_REDIRECT_URL"))
	if redirect == "" {
		redirect = "http://localhost:8080/api/v1/auth/oidc/callback"
	}
	adminUI := strings.TrimSpace(os.Getenv("ADMIN_UI_URL"))
	if adminUI == "" {
		adminUI = "http://localhost:5173"
	}
	p := &OIDCProvider{
		enabled: true,
		oauth: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: secret,
			RedirectURL:  redirect,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		},
		verifier: provider.Verifier(&oidc.Config{ClientID: clientID}),
		adminUI:  strings.TrimRight(adminUI, "/"),
		admins:   parseAdminEmails(os.Getenv("OIDC_ADMIN_EMAILS")),
		adminDom: strings.TrimSpace(os.Getenv("OIDC_ADMIN_DOMAIN")),
		auth:     authenticator,
		enforce:  os.Getenv("OIDC_ENFORCE_ADMIN") == "true",
	}
	return p, nil
}

func (p *OIDCProvider) Enabled() bool { return p != nil && p.enabled }

func (p *OIDCProvider) EnforceAdmin() bool { return p.Enabled() && p.enforce }

// ConfigHandler GET /api/v1/auth/oidc/config
func (p *OIDCProvider) ConfigHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"enabled":        p.Enabled(),
			"login_url":      "/api/v1/auth/oidc/login",
			"enforce_admin":  p.EnforceAdmin(),
			"bootstrap_hint": !p.EnforceAdmin(),
		})
	}
}

// LoginHandler GET /api/v1/auth/oidc/login
func (p *OIDCProvider) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !p.Enabled() {
			http.Error(w, "oidc not configured", http.StatusNotFound)
			return
		}
		state, err := randomState()
		if err != nil {
			http.Error(w, "state error", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     oidcStateCookie,
			Value:    state,
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https"),
			SameSite: http.SameSiteLaxMode,
			MaxAge:   600,
		})
		http.Redirect(w, r, p.oauth.AuthCodeURL(state), http.StatusFound)
	}
}

// CallbackHandler GET /api/v1/auth/oidc/callback
func (p *OIDCProvider) CallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !p.Enabled() {
			http.Error(w, "oidc not configured", http.StatusNotFound)
			return
		}
		if errParam := r.URL.Query().Get("error"); errParam != "" {
			http.Error(w, errParam, http.StatusBadRequest)
			return
		}
		stateCookie, err := r.Cookie(oidcStateCookie)
		if err != nil || stateCookie.Value == "" || stateCookie.Value != r.URL.Query().Get("state") {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}
		http.SetCookie(w, &http.Cookie{Name: oidcStateCookie, Value: "", Path: "/", MaxAge: -1})

		tok, err := p.oauth.Exchange(r.Context(), r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "token exchange failed", http.StatusBadRequest)
			return
		}
		rawID, ok := tok.Extra("id_token").(string)
		if !ok || rawID == "" {
			http.Error(w, "id_token missing", http.StatusBadRequest)
			return
		}
		idToken, err := p.verifier.Verify(r.Context(), rawID)
		if err != nil {
			http.Error(w, "id_token invalid", http.StatusUnauthorized)
			return
		}
		var claims struct {
			Email         string `json:"email"`
			EmailVerified bool   `json:"email_verified"`
			Sub           string `json:"sub"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "claims parse failed", http.StatusBadRequest)
			return
		}
		email := strings.ToLower(strings.TrimSpace(claims.Email))
		if email == "" || !claims.EmailVerified {
			http.Error(w, "verified email required", http.StatusForbidden)
			return
		}
		if !p.isAdmin(email) {
			http.Error(w, "not an admin identity", http.StatusForbidden)
			return
		}
		userID := claims.Sub
		if userID == "" {
			userID = email
		}
		pair, err := p.auth.GenerateToken(r.Context(), &auth.Claims{
			UserID: userID,
			Roles:  []string{"admin"},
			Metadata: map[string]any{
				"email": email,
				"sso":   "oidc",
			},
			ExpiresAt: time.Now().Add(8 * time.Hour),
		})
		if err != nil {
			http.Error(w, "token issue failed", http.StatusInternalServerError)
			return
		}
		redirect := p.adminUI + "/?access_token=" + url.QueryEscape(pair.AccessToken)
		http.Redirect(w, r, redirect, http.StatusFound)
	}
}

func (p *OIDCProvider) isAdmin(email string) bool {
	if len(p.admins) > 0 {
		_, ok := p.admins[strings.ToLower(email)]
		return ok
	}
	if p.adminDom != "" {
		return strings.HasSuffix(email, "@"+strings.ToLower(p.adminDom))
	}
	// Dev-friendly fallback when OIDC is enabled but allowlist not set (non-prod only).
	if os.Getenv("MSGGUARD_ENV") != "production" && os.Getenv("MSGGUARD_ENV") != "prod" {
		return true
	}
	return false
}

func parseAdminEmails(raw string) map[string]struct{} {
	out := make(map[string]struct{})
	for _, part := range strings.Split(raw, ",") {
		e := strings.ToLower(strings.TrimSpace(part))
		if e != "" {
			out[e] = struct{}{}
		}
	}
	return out
}

func randomState() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
