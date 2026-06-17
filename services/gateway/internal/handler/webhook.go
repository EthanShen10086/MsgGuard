package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/EthanShen10086/voxera-kit/observability/logger"

	"github.com/EthanShen10086/msgguard/pkg/adapters/appstore"
	"github.com/EthanShen10086/msgguard/pkg/ports"
)

// WebhookHandler processes App Store Server Notifications V2.
type WebhookHandler struct {
	log      logger.Logger
	store    ports.SubscriptionStore
	verifier *appstore.Verifier
}

func NewWebhookHandler(log logger.Logger, store ports.SubscriptionStore, verifier *appstore.Verifier) *WebhookHandler {
	return &WebhookHandler{log: log, store: store, verifier: verifier}
}

type appStoreNotification struct {
	SignedPayload string `json:"signedPayload"`
}

// AppStore handles POST /api/v1/webhooks/appstore
func (h *WebhookHandler) AppStore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var note appStoreNotification
	if err := json.Unmarshal(body, &note); err != nil || note.SignedPayload == "" {
		http.Error(w, "signedPayload required", http.StatusBadRequest)
		return
	}
	result, err := h.verifier.VerifySignedTransaction(note.SignedPayload, "com.ethanshen.msgguard.pro.monthly")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Device mapping requires prior verify linking transaction to device_id — log for ops.
	h.log.Info("appstore webhook",
		logger.Field{Key: "transaction_id", Value: result.TransactionID},
		logger.Field{Key: "product_id", Value: result.ProductID},
	)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"accepted"}`))
}
