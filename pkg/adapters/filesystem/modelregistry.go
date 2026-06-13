package filesystem

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type ModelRegistry struct {
	mu   sync.RWMutex
	dir  string
	meta map[string]ports.ModelMeta
}

func NewModelRegistry(dir string) (*ModelRegistry, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	r := &ModelRegistry{dir: dir, meta: map[string]ports.ModelMeta{}}
	if data, err := os.ReadFile(filepath.Join(dir, "registry.json")); err == nil {
		_ = json.Unmarshal(data, &r.meta)
	}
	return r, nil
}

func (r *ModelRegistry) Register(ctx context.Context, meta ports.ModelMeta) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.meta[meta.Locale] = meta
	data, _ := json.MarshalIndent(r.meta, "", "  ")
	return os.WriteFile(filepath.Join(r.dir, "registry.json"), data, 0o644)
}

func (r *ModelRegistry) GetLatest(ctx context.Context, locale string) (*ports.ModelMeta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if m, ok := r.meta[locale]; ok {
		cp := m
		return &cp, nil
	}
	return &ports.ModelMeta{Version: "1.0.0", Locale: locale, Checksum: "sha256:seed"}, nil
}

func (r *ModelRegistry) GetArtifact(ctx context.Context, version, name string) ([]byte, error) {
	for locale := range r.meta {
		path := filepath.Join(r.dir, locale, version, name)
		if data, err := os.ReadFile(path); err == nil {
			return data, nil
		}
	}
	path := filepath.Join(r.dir, version, name)
	return os.ReadFile(path)
}
