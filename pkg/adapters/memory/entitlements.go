package memory

import (
	"context"
	"sync"
	"time"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type SubscriptionStore struct {
	mu   sync.RWMutex
	subs map[string]ports.Subscription
}

func NewSubscriptionStore() *SubscriptionStore {
	return &SubscriptionStore{subs: make(map[string]ports.Subscription)}
}

func (s *SubscriptionStore) Upsert(ctx context.Context, sub ports.Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sub.UpdatedAt.IsZero() {
		sub.UpdatedAt = time.Now().UTC()
	}
	s.subs[sub.DeviceID] = sub
	return nil
}

func (s *SubscriptionStore) GetByDeviceID(ctx context.Context, deviceID string) (*ports.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sub, ok := s.subs[deviceID]
	if !ok {
		return &ports.Subscription{DeviceID: deviceID, IsPro: false}, nil
	}
	cpy := sub
	return &cpy, nil
}
