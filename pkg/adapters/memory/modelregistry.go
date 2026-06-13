package memory

import (
	"context"
	"sync"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type ModelRegistry struct {
	mu     sync.RWMutex
	latest map[string]ports.ModelMeta
}

func NewModelRegistry() *ModelRegistry {
	return &ModelRegistry{latest: map[string]ports.ModelMeta{}}
}

func (r *ModelRegistry) Register(ctx context.Context, meta ports.ModelMeta) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.latest[meta.Locale] = meta
	return nil
}

func (r *ModelRegistry) GetLatest(ctx context.Context, locale string) (*ports.ModelMeta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if m, ok := r.latest[locale]; ok {
		cp := m
		return &cp, nil
	}
	return &ports.ModelMeta{
		Version: "1.0.0", Locale: locale,
		Checksum: "sha256:seed",
	}, nil
}

func (r *ModelRegistry) GetArtifact(ctx context.Context, version, name string) ([]byte, error) {
	return nil, nil
}
