package handler

import (
	"encoding/json"
	"net/http"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type RulesHandler struct {
	store ports.RuleStore
}

func NewRulesHandler(store ports.RuleStore) *RulesHandler {
	return &RulesHandler{store: store}
}

func (h *RulesHandler) Latest(w http.ResponseWriter, r *http.Request) {
	bundle, err := h.store.GetLatest(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ETag", bundle.Checksum)
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
