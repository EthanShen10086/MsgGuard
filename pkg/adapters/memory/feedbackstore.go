package memory

import (
	"context"
	"sync"
	"time"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type FeedbackStore struct {
	mu    sync.Mutex
	items []ports.FeedbackItem
}

func NewFeedbackStore() *FeedbackStore { return &FeedbackStore{} }

func (s *FeedbackStore) Create(ctx context.Context, item ports.FeedbackItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = time.Now().UTC()
	}
	s.items = append(s.items, item)
	return nil
}

func (s *FeedbackStore) List(ctx context.Context, limit int) ([]ports.FeedbackItem, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 || limit > len(s.items) {
		limit = len(s.items)
	}
	start := len(s.items) - limit
	if start < 0 {
		start = 0
	}
	out := make([]ports.FeedbackItem, len(s.items[start:]))
	copy(out, s.items[start:])
	return out, nil
}
