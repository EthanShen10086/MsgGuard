package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/EthanShen10086/voxera-kit/observability/logger"
)

type ShadowHandler struct {
	log      logger.Logger
	classify *ClassifyHandler
	mu       sync.Mutex
	metrics  []shadowRecord
}

type shadowRecord struct {
	Body         string `json:"body"`
	LocalAction  string `json:"local_action"`
	CloudAction  string `json:"cloud_action"`
	Agree        bool   `json:"agree"`
}

func NewShadowHandler(log logger.Logger, classify *ClassifyHandler) *ShadowHandler {
	return &ShadowHandler{log: log, classify: classify}
}

func (h *ShadowHandler) Compare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var req classifyRequest
	_ = json.Unmarshal(body, &req)
	local := heuristicClassify(req.Body)
	cloud, _ := h.classify.runInternal(r, req)
	rec := shadowRecord{
		Body: req.Body, LocalAction: local.Action,
		CloudAction: cloud.Action, Agree: local.Action == cloud.Action,
	}
	h.mu.Lock()
	h.metrics = append(h.metrics, rec)
	h.mu.Unlock()
	writeJSON(w, rec)
}

func (h *ShadowHandler) Stats(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()
	agree := 0
	for _, m := range h.metrics {
		if m.Agree {
			agree++
		}
	}
	writeJSON(w, map[string]any{"total": len(h.metrics), "agree": agree})
}
