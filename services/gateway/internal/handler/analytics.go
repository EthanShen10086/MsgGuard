package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/EthanShen10086/voxera-kit/observability/logger"
	"github.com/google/uuid"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type AnalyticsHandler struct {
	log   logger.Logger
	store ports.AnalyticsStore
}

func NewAnalyticsHandler(log logger.Logger, store ports.AnalyticsStore) *AnalyticsHandler {
	return &AnalyticsHandler{log: log, store: store}
}

func (h *AnalyticsHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var req struct {
		Name     string         `json:"name"`
		Props    map[string]any `json:"props"`
		DeviceID string         `json:"device_id"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	traceID := r.Header.Get("X-Request-ID")
	event := ports.AnalyticsEvent{
		ID: uuid.NewString(), Name: req.Name, Props: req.Props,
		DeviceID: req.DeviceID, TraceID: traceID, Timestamp: time.Now().UTC(),
	}
	_ = h.store.Insert(r.Context(), event)
	h.log.Info("analytics", logger.Field{Key: "event", Value: req.Name})
	writeJSON(w, map[string]string{"status": "ok"})
}
