package memory

import (
	"context"
	"sync"
	"time"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type AnalyticsStore struct {
	mu     sync.Mutex
	events []ports.AnalyticsEvent
}

func NewAnalyticsStore() *AnalyticsStore { return &AnalyticsStore{} }

func (s *AnalyticsStore) Insert(ctx context.Context, event ports.AnalyticsEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}
	s.events = append(s.events, event)
	return nil
}

func (s *AnalyticsStore) List(ctx context.Context, since time.Time, limit int) ([]ports.AnalyticsEvent, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []ports.AnalyticsEvent
	for i := len(s.events) - 1; i >= 0; i-- {
		if !since.IsZero() && s.events[i].Timestamp.Before(since) {
			continue
		}
		out = append(out, s.events[i])
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (s *AnalyticsStore) CountByName(ctx context.Context, since time.Time) (map[string]int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	counts := map[string]int{}
	for _, e := range s.events {
		if !since.IsZero() && e.Timestamp.Before(since) {
			continue
		}
		counts[e.Name]++
	}
	return counts, nil
}

func (s *AnalyticsStore) DeleteByDeviceID(ctx context.Context, deviceID string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if deviceID == "" {
		return 0, nil
	}
	var kept []ports.AnalyticsEvent
	deleted := 0
	for _, e := range s.events {
		if e.DeviceID == deviceID {
			deleted++
			continue
		}
		kept = append(kept, e)
	}
	s.events = kept
	return deleted, nil
}
