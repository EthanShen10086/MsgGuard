package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/EthanShen10086/voxera-kit/auth"
	"github.com/EthanShen10086/voxera-kit/observability/logger"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type AdminHandler struct {
	log           logger.Logger
	analytics     ports.AnalyticsStore
	feedback      ports.FeedbackStore
	authenticator auth.Authenticator
	authorizer    auth.Authorizer
}

func NewAdminHandler(
	log logger.Logger,
	analytics ports.AnalyticsStore,
	feedback ports.FeedbackStore,
	authenticator auth.Authenticator,
	authorizer auth.Authorizer,
) *AdminHandler {
	return &AdminHandler{
		log: log, analytics: analytics, feedback: feedback,
		authenticator: authenticator, authorizer: authorizer,
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
		"period_days": 7,
		"event_counts": counts,
		"feedback_total": len(feedback),
		"generated_at": time.Now().UTC(),
	})
}

func (h *AdminHandler) authenticate(r *http.Request) (*auth.Claims, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return nil, errors.New("missing authorization")
	}
	return h.authenticator.Authenticate(r.Context(), token)
}
