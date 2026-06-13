package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/EthanShen10086/voxera-kit/audit"
	"github.com/EthanShen10086/voxera-kit/observability/logger"
	"github.com/google/uuid"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type FeedbackHandler struct {
	log         logger.Logger
	store       ports.FeedbackStore
	auditWriter audit.Writer
	queue       ports.Queue
}

func NewFeedbackHandler(log logger.Logger, store ports.FeedbackStore, auditWriter audit.Writer, queue ports.Queue) *FeedbackHandler {
	return &FeedbackHandler{log: log, store: store, auditWriter: auditWriter, queue: queue}
}

type feedbackRequest struct {
	Description string `json:"description"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	Body        string `json:"body"`
	Locale      string `json:"locale"`
	TraceID     string `json:"traceID"`
}

type feedbackResponse struct {
	ID      string `json:"id"`
	TraceID string `json:"traceID"`
}

func (h *FeedbackHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req feedbackRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	traceID := r.Header.Get("X-Request-ID")
	if req.TraceID != "" {
		traceID = req.TraceID
	}
	id := traceID
	if id == "" {
		id = uuid.NewString()
	}
	label := req.Label
	if label == "" {
		label = req.Category
	}
	text := req.Body
	if text == "" {
		text = req.Description
	}
	item := ports.FeedbackItem{
		ID: id, Body: text, Label: label, Locale: req.Locale,
		TraceID: traceID, CreatedAt: time.Now().UTC(),
	}
	_ = h.store.Create(r.Context(), item)
	if h.auditWriter != nil {
		_ = h.auditWriter.Write(r.Context(), audit.Entry{
			ID: uuid.NewString(), ActorID: "user", Action: "feedback.create",
			ResourceType: "feedback", ResourceID: id, Timestamp: time.Now().UTC(),
		})
	}
	if h.queue != nil {
		payload, _ := json.Marshal(map[string]string{"id": id, "label": label})
		_ = h.queue.Publish(r.Context(), "msgguard.flywheel.trigger", payload)
	}
	h.log.Info("feedback received", logger.Field{Key: "trace_id", Value: traceID})
	writeJSON(w, feedbackResponse{ID: id, TraceID: traceID})
}

func (h *FeedbackHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	items, err := h.store.List(r.Context(), 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, items)
}
