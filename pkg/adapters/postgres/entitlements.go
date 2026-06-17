package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type SubscriptionStore struct {
	db *sql.DB
}

func NewSubscriptionStore(dsn string) (*SubscriptionStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &SubscriptionStore{db: db}, nil
}

func (s *SubscriptionStore) Upsert(ctx context.Context, sub ports.Subscription) error {
	ts := sub.UpdatedAt
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO subscriptions (device_id, product_id, signed_transaction, is_pro, expires_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (device_id) DO UPDATE SET
			product_id = EXCLUDED.product_id,
			signed_transaction = EXCLUDED.signed_transaction,
			is_pro = EXCLUDED.is_pro,
			expires_at = EXCLUDED.expires_at,
			updated_at = EXCLUDED.updated_at`,
		sub.DeviceID, sub.ProductID, sub.SignedTransaction, sub.IsPro, sub.ExpiresAt, ts,
	)
	return err
}

func (s *SubscriptionStore) GetByDeviceID(ctx context.Context, deviceID string) (*ports.Subscription, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT device_id, product_id, signed_transaction, is_pro, expires_at, updated_at
		FROM subscriptions WHERE device_id = $1`, deviceID)
	var sub ports.Subscription
	var expires sql.NullTime
	if err := row.Scan(&sub.DeviceID, &sub.ProductID, &sub.SignedTransaction, &sub.IsPro, &expires, &sub.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return &ports.Subscription{DeviceID: deviceID, IsPro: false}, nil
		}
		return nil, err
	}
	if expires.Valid {
		t := expires.Time
		sub.ExpiresAt = &t
	}
	return &sub, nil
}
