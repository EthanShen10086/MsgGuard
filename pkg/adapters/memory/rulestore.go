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
		payload := []byte(`{"version":"1.0.0","locale":"zh-Hans","checksum":"seed","keywords_block":["免费领取","中奖","贷款"],"keywords_allow":["验证码","verification code"]}`)
		return &ports.RuleBundle{Version: "1.0.0", Checksum: "seed", Payload: payload}, nil
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
