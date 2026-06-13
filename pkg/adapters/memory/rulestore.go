package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type RuleStore struct {
	mu       sync.RWMutex
	latest   *ports.RuleBundle
	byVersion map[string]*ports.RuleBundle
}

func NewRuleStore() *RuleStore { return &RuleStore{byVersion: map[string]*ports.RuleBundle{}} }

func (s *RuleStore) seed() *ports.RuleBundle {
	payload := []byte(`{"version":"1.0.0","locale":"zh-Hans","checksum":"seed","keywords_block":["免费领取","中奖","贷款"],"keywords_allow":["验证码","verification code"]}`)
	return &ports.RuleBundle{Version: "1.0.0", Checksum: "seed", Payload: payload}
}

func (s *RuleStore) GetLatest(ctx context.Context) (*ports.RuleBundle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.latest == nil {
		return s.seed(), nil
	}
	return s.latest, nil
}

func (s *RuleStore) GetByVersion(ctx context.Context, version string) (*ports.RuleBundle, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if b, ok := s.byVersion[version]; ok {
		return b, nil
	}
	if s.latest != nil && s.latest.Version == version {
		return s.latest, nil
	}
	if version == "1.0.0" && s.latest == nil {
		return s.seed(), nil
	}
	return nil, fmt.Errorf("version not found")
}

func (s *RuleStore) Save(ctx context.Context, bundle ports.RuleBundle) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b := bundle
	s.latest = &b
	s.byVersion[b.Version] = &b
	return nil
}
