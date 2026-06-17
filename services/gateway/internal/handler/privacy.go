package handler

import (
	"net/http"

	"github.com/EthanShen10086/voxera-kit/observability/logger"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

// PrivacyHandler handles GDPR-style data subject requests.
type PrivacyHandler struct {
	log       logger.Logger
	analytics ports.AnalyticsStore
}

func NewPrivacyHandler(log logger.Logger, analytics ports.AnalyticsStore) *PrivacyHandler {
	return &PrivacyHandler{log: log, analytics: analytics}
}

// DeleteMe handles DELETE /api/v1/privacy/me?device_id=<uuid>
func (h *PrivacyHandler) DeleteMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	deviceID := r.URL.Query().Get("device_id")
	if deviceID == "" {
		http.Error(w, "device_id required", http.StatusBadRequest)
		return
	}
	deleted, err := h.analytics.DeleteByDeviceID(r.Context(), deviceID)
	if err != nil {
		h.log.Error("privacy delete failed", logger.Field{Key: "error", Value: err.Error()})
		http.Error(w, "deletion failed", http.StatusInternalServerError)
		return
	}
	if deleted == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	h.log.Info("privacy delete", logger.Field{Key: "device_id", Value: deviceID}, logger.Field{Key: "deleted", Value: deleted})
	writeJSON(w, map[string]any{
		"status":         "deleted",
		"device_id":      deviceID,
		"analytics_rows": deleted,
	})
}
