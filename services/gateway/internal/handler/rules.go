package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/EthanShen10086/voxera-kit/auth"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type RulesHandler struct {
	store         ports.RuleStore
	authenticator auth.Authenticator
	authorizer    auth.Authorizer
}

func NewRulesHandler(store ports.RuleStore, authenticator auth.Authenticator, authorizer auth.Authorizer) *RulesHandler {
	return &RulesHandler{store: store, authenticator: authenticator, authorizer: authorizer}
}

func (h *RulesHandler) Latest(w http.ResponseWriter, r *http.Request) {
	bundle, err := h.store.GetLatest(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if match := r.Header.Get("If-None-Match"); match != "" && strings.Trim(match, `"`) == bundle.Checksum {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	h.writeBundle(w, bundle)
}

func (h *RulesHandler) ByVersion(w http.ResponseWriter, r *http.Request) {
	version := strings.TrimPrefix(r.URL.Path, "/api/v1/rules/")
	if version == "" || version == "latest" {
		h.Latest(w, r)
		return
	}
	bundle, err := h.store.GetByVersion(r.Context(), version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if match := r.Header.Get("If-None-Match"); match != "" && strings.Trim(match, `"`) == bundle.Checksum {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	h.writeBundle(w, bundle)
}

func (h *RulesHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	claims, err := h.authenticator.Authenticate(r.Context(), r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	ok, _ := h.authorizer.Authorize(r.Context(), claims, "rules", "write")
	if !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	var req struct {
		Version  string          `json:"version"`
		Checksum string          `json:"checksum"`
		Locale   string          `json:"locale"`
		Payload  json.RawMessage `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Version == "" {
		http.Error(w, "version required", http.StatusBadRequest)
		return
	}
	body := req.Payload
	if len(body) == 0 {
		body, _ = json.Marshal(map[string]any{
			"version": req.Version, "checksum": req.Checksum, "locale": req.Locale,
		})
	}
	bundle := ports.RuleBundle{Version: req.Version, Checksum: req.Checksum, Payload: body}
	if bundle.Checksum == "" {
		bundle.Checksum = req.Version
	}
	if err := h.store.Save(r.Context(), bundle); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]string{"status": "registered", "version": req.Version})
}

func (h *RulesHandler) writeBundle(w http.ResponseWriter, bundle *ports.RuleBundle) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ETag", `"`+bundle.Checksum+`"`)
	if len(bundle.Payload) > 0 {
		w.Write(bundle.Payload)
		return
	}
	json.NewEncoder(w).Encode(map[string]any{
		"version": bundle.Version, "checksum": bundle.Checksum,
		"keywords_block": []string{"免费领取", "中奖", "贷款"},
		"keywords_allow": []string{"验证码", "verification code"},
	})
}
