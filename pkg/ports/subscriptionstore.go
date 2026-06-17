package ports

import (
	"context"
	"time"
)

// Subscription records App Store entitlement state for a device.
type Subscription struct {
	DeviceID          string     `json:"device_id"`
	ProductID         string     `json:"product_id"`
	SignedTransaction string     `json:"signed_transaction,omitempty"`
	IsPro             bool       `json:"is_pro"`
	ExpiresAt         *time.Time `json:"expires_at,omitempty"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// SubscriptionStore persists subscription / pro status per device.
type SubscriptionStore interface {
	Upsert(ctx context.Context, sub Subscription) error
	GetByDeviceID(ctx context.Context, deviceID string) (*Subscription, error)
}
