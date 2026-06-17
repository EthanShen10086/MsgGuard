package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/EthanShen10086/voxera-kit/auth"
	"github.com/EthanShen10086/voxera-kit/observability/logger"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

// ModelAdminHandler promotes or rolls back on-device model versions in the registry.
type ModelAdminHandler struct {
	log      logger.Logger
	registry ports.ModelRegistry
	auth     auth.Authenticator
	authz    auth.Authorizer
}

func NewModelAdminHandler(
	log logger.Logger,
	registry ports.ModelRegistry,
	authenticator auth.Authenticator,
	authorizer auth.Authorizer,
) *ModelAdminHandler {
	return &ModelAdminHandler{
		log: log, registry: registry, auth: authenticator, authz: authorizer,
	}
}

type promoteRequest struct {
	Version   string                `json:"version"`
	Locale    string                `json:"locale"`
	Checksum  string                `json:"checksum"`
	Artifacts []ports.ModelArtifact `json:"artifacts"`
}

type rollbackRequest struct {
	Locale  string `json:"locale"`
	Version string `json:"version"`
}

func (h *ModelAdminHandler) Promote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.requireAdmin(r); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req promoteRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Locale == "" {
		req.Locale = "zh-Hans"
	}
	if req.Version == "" {
		http.Error(w, "version required", http.StatusBadRequest)
		return
	}
	meta := ports.ModelMeta{
		Version: req.Version, Locale: req.Locale,
		Checksum: req.Checksum, Artifacts: req.Artifacts,
	}
	if meta.Checksum == "" {
		meta.Checksum = "sha256:promoted-" + req.Version
	}
	if err := h.registry.Register(r.Context(), meta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.log.Info("model promoted", logger.Field{Key: "locale", Value: req.Locale}, logger.Field{Key: "version", Value: req.Version})
	writeJSON(w, map[string]any{"status": "promoted", "locale": req.Locale, "version": req.Version})
}

func (h *ModelAdminHandler) Rollback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.requireAdmin(r); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var req rollbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Locale == "" {
		req.Locale = "zh-Hans"
	}
	rollbacker, ok := h.registry.(ports.ModelRollback)
	if !ok {
		http.Error(w, "rollback not supported by registry", http.StatusNotImplemented)
		return
	}
	meta, err := rollbacker.Rollback(r.Context(), req.Locale, req.Version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.log.Info("model rolled back", logger.Field{Key: "locale", Value: req.Locale}, logger.Field{Key: "version", Value: meta.Version})
	writeJSON(w, map[string]any{"status": "rolled_back", "locale": meta.Locale, "version": meta.Version})
}

func (h *ModelAdminHandler) requireAdmin(r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token == "" {
		return errors.New("unauthorized")
	}
	claims, err := h.auth.Authenticate(r.Context(), token)
	if err != nil {
		return err
	}
	ok, _ := h.authz.Authorize(r.Context(), claims, "models", "write")
	if !ok {
		return errors.New("forbidden")
	}
	return nil
}
