package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/EthanShen10086/voxera-kit/observability/logger"

	"github.com/EthanShen10086/msgguard/pkg/adapters/appstore"
	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type EntitlementsHandler struct {
	log      logger.Logger
	store    ports.SubscriptionStore
	verifier *appstore.Verifier
}

func NewEntitlementsHandler(log logger.Logger, store ports.SubscriptionStore, verifier *appstore.Verifier) *EntitlementsHandler {
	return &EntitlementsHandler{log: log, store: store, verifier: verifier}
}

type verifyEntitlementRequest struct {
	DeviceID          string `json:"device_id"`
	SignedTransaction string `json:"signed_transaction"`
	ProductID         string `json:"product_id"`
}

type entitlementStatusResponse struct {
	DeviceID  string `json:"device_id"`
	ProductID string `json:"product_id,omitempty"`
	IsPro     bool   `json:"is_pro"`
}

func (h *EntitlementsHandler) Verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req verifyEntitlementRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.DeviceID = strings.TrimSpace(req.DeviceID)
	if req.DeviceID == "" {
		http.Error(w, "device_id required", http.StatusBadRequest)
		return
	}

	result, err := h.verifier.VerifySignedTransaction(req.SignedTransaction, req.ProductID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub := ports.Subscription{
		DeviceID:          req.DeviceID,
		ProductID:         result.ProductID,
		SignedTransaction: req.SignedTransaction,
		IsPro:             result.Valid,
		ExpiresAt:         result.ExpiresAt,
		UpdatedAt:         time.Now().UTC(),
	}
	if err := h.store.Upsert(r.Context(), sub); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.log.Info("entitlement verified", logger.Field{Key: "device_id", Value: req.DeviceID})
	writeJSON(w, entitlementStatusResponse{
		DeviceID: req.DeviceID, ProductID: sub.ProductID, IsPro: sub.IsPro,
	})
}

func (h *EntitlementsHandler) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	deviceID := strings.TrimSpace(r.URL.Query().Get("device_id"))
	if deviceID == "" {
		http.Error(w, "device_id required", http.StatusBadRequest)
		return
	}
	sub, err := h.store.GetByDeviceID(r.Context(), deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, entitlementStatusResponse{
		DeviceID: sub.DeviceID, ProductID: sub.ProductID, IsPro: sub.IsPro,
	})
}
