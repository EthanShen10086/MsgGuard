package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/EthanShen10086/voxera-kit/observability/logger"
)

type AnalyticsHandler struct {
	log logger.Logger
	mu  sync.Mutex
	events []analyticsEvent
}

type analyticsEvent struct {
	Name      string         `json:"name"`
	Props     map[string]any `json:"props"`
	Timestamp time.Time      `json:"timestamp"`
}

func NewAnalyticsHandler(log logger.Logger) *AnalyticsHandler {
	return &AnalyticsHandler{log: log}
}

func (h *AnalyticsHandler) Ingest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var req struct {
		Name  string         `json:"name"`
		Props map[string]any `json:"props"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.mu.Lock()
	h.events = append(h.events, analyticsEvent{Name: req.Name, Props: req.Props, Timestamp: time.Now().UTC()})
	h.mu.Unlock()
	h.log.Info("analytics", logger.Field{Key: "event", Value: req.Name})
	writeJSON(w, map[string]string{"status": "ok"})
}

func (h *AnalyticsHandler) Count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.events)
}
