package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/EthanShen10086/voxera-kit/featureflag"
	"github.com/EthanShen10086/voxera-kit/observability/logger"
)

type ShadowHandler struct {
	log      logger.Logger
	classify *ClassifyHandler
	flags    featureflag.Store
	mu       sync.Mutex
	metrics  []shadowRecord
}

type shadowRecord struct {
	Body        string `json:"body"`
	LocalAction string `json:"local_action"`
	CloudAction string `json:"cloud_action"`
	Agree       bool   `json:"agree"`
}

func NewShadowHandler(log logger.Logger, classify *ClassifyHandler, flags featureflag.Store) *ShadowHandler {
	return &ShadowHandler{log: log, classify: classify, flags: flags}
}

func (h *ShadowHandler) Compare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if h.flags != nil {
		enabled, _ := h.flags.IsEnabled(r.Context(), "shadow_mode", featureflag.EvalContext{})
		if !enabled {
			http.Error(w, "shadow mode disabled", http.StatusServiceUnavailable)
			return
		}
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
	recordShadowCompare(rec.Agree)
	writeJSON(w, rec)
}

func (h *ShadowHandler) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, shadowStatsSnapshot())
}

func (h *ShadowHandler) Prometheus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeShadowPrometheus(w)
}
