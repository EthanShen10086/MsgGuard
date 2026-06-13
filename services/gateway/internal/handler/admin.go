package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/EthanShen10086/voxera-kit/aiquota"
	"github.com/EthanShen10086/voxera-kit/auth"
	"github.com/EthanShen10086/voxera-kit/featureflag"
	"github.com/EthanShen10086/voxera-kit/observability/logger"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type AdminHandler struct {
	log           logger.Logger
	analytics     ports.AnalyticsStore
	feedback      ports.FeedbackStore
	authenticator auth.Authenticator
	authorizer    auth.Authorizer
	quota         aiquota.Manager
	flags         featureflag.Store
}

func NewAdminHandler(
	log logger.Logger,
	analytics ports.AnalyticsStore,
	feedback ports.FeedbackStore,
	authenticator auth.Authenticator,
	authorizer auth.Authorizer,
	quota aiquota.Manager,
	flags featureflag.Store,
) *AdminHandler {
	return &AdminHandler{
		log: log, analytics: analytics, feedback: feedback,
		authenticator: authenticator, authorizer: authorizer,
		quota: quota, flags: flags,
	}
}

func (h *AdminHandler) IssueToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		UserID string   `json:"user_id"`
		Roles  []string `json:"roles"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if len(req.Roles) == 0 {
		req.Roles = []string{"admin"}
	}
	pair, err := h.authenticator.GenerateToken(r.Context(), &auth.Claims{
		UserID: req.UserID, Roles: req.Roles,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{
		"access_token": pair.AccessToken,
		"token_type":   pair.TokenType,
		"expires_at":   pair.ExpiresAt,
	})
}

func (h *AdminHandler) MetricsSummary(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	ok, _ := h.authorizer.Authorize(r.Context(), claims, "admin", "read")
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	since := time.Now().Add(-7 * 24 * time.Hour)
	counts, _ := h.analytics.CountByName(r.Context(), since)
	feedback, _ := h.feedback.List(r.Context(), 1000)
	writeJSON(w, map[string]any{
		"period_days":    7,
		"event_counts":   counts,
		"feedback_total": len(feedback),
		"generated_at":   time.Now().UTC(),
	})
}

func (h *AdminHandler) QuotaWhitelist(w http.ResponseWriter, r *http.Request) {
	if err := h.requireAdmin(r); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	switch r.Method {
	case http.MethodGet:
		entries, _ := h.quota.ListWhitelist(r.Context())
		writeJSON(w, entries)
	case http.MethodPost:
		var req struct {
			UserID string `json:"user_id"`
			Reason string `json:"reason"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = h.quota.AddToWhitelist(r.Context(), aiquota.WhitelistEntry{UserID: req.UserID, Reason: req.Reason})
		w.WriteHeader(http.StatusCreated)
		writeJSON(w, map[string]string{"status": "whitelisted"})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AdminHandler) FeatureFlags(w http.ResponseWriter, r *http.Request) {
	if err := h.requireAdmin(r); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	switch r.Method {
	case http.MethodGet:
		flags, _ := h.flags.GetFlags(r.Context())
		writeJSON(w, flags)
	case http.MethodPost, http.MethodPut:
		var flag featureflag.Flag
		if err := json.NewDecoder(r.Body).Decode(&flag); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = h.flags.SetFlag(r.Context(), flag)
		writeJSON(w, map[string]string{"status": "updated"})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AdminHandler) requireAdmin(r *http.Request) error {
	claims, err := h.authenticate(r)
	if err != nil {
		return err
	}
	ok, _ := h.authorizer.Authorize(r.Context(), claims, "admin", "read")
	if !ok {
		return errors.New("forbidden")
	}
	return nil
}

func (h *AdminHandler) authenticate(r *http.Request) (*auth.Claims, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("missing authorization")
	}
	return h.authenticator.Authenticate(r.Context(), token)
}
