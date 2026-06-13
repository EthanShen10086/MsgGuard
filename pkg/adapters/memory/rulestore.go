package memory

import (
	"context"
	"sync"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type RuleStore struct {
	mu     sync.RWMutex
	latest *ports.RuleBundle
}

func NewRuleStore() *RuleStore { return &RuleStore{} }

func (s *RuleStore) GetLatest(ctx context.Context) (*ports.RuleBundle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.latest == nil {
		return &ports.RuleBundle{Version: "1.0.0", Checksum: "seed", Payload: []byte(`{"keywords":["免费","中奖"]}`)}, nil
	}
	return s.latest, nil
}

func (s *RuleStore) Save(ctx context.Context, bundle ports.RuleBundle) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b := bundle
	s.latest = &b
	return nil
}
