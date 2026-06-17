package filesystem

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type ModelRegistry struct {
	mu      sync.RWMutex
	dir     string
	meta    map[string]ports.ModelMeta
	history map[string][]ports.ModelMeta
}

func NewModelRegistry(dir string) (*ModelRegistry, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	r := &ModelRegistry{dir: dir, meta: map[string]ports.ModelMeta{}, history: map[string][]ports.ModelMeta{}}
	if data, err := os.ReadFile(filepath.Join(dir, "registry.json")); err == nil {
		_ = json.Unmarshal(data, &r.meta)
	}
	if data, err := os.ReadFile(filepath.Join(dir, "history.json")); err == nil {
		_ = json.Unmarshal(data, &r.history)
	}
	return r, nil
}

func (r *ModelRegistry) Register(ctx context.Context, meta ports.ModelMeta) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if prev, ok := r.meta[meta.Locale]; ok && prev.Version != meta.Version {
		r.history[meta.Locale] = append(r.history[meta.Locale], prev)
	}
	r.meta[meta.Locale] = meta
	if err := r.persistLocked(); err != nil {
		return err
	}
	_ = ctx
	return nil
}

func (r *ModelRegistry) Rollback(ctx context.Context, locale, targetVersion string) (*ports.ModelMeta, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	stack := r.history[locale]
	if len(stack) == 0 {
		return nil, fmt.Errorf("no history for locale %s", locale)
	}
	var meta ports.ModelMeta
	if targetVersion != "" {
		idx := -1
		for i := len(stack) - 1; i >= 0; i-- {
			if stack[i].Version == targetVersion {
				idx = i
				break
			}
		}
		if idx < 0 {
			return nil, fmt.Errorf("version %s not in history", targetVersion)
		}
		meta = stack[idx]
		r.history[locale] = append(stack[:idx], stack[idx+1:]...)
	} else {
		meta = stack[len(stack)-1]
		r.history[locale] = stack[:len(stack)-1]
	}
	if cur, ok := r.meta[locale]; ok {
		r.history[locale] = append(r.history[locale], cur)
	}
	r.meta[locale] = meta
	if err := r.persistLocked(); err != nil {
		return nil, err
	}
	cp := meta
	_ = ctx
	return &cp, nil
}

func (r *ModelRegistry) persistLocked() error {
	data, _ := json.MarshalIndent(r.meta, "", "  ")
	if err := os.WriteFile(filepath.Join(r.dir, "registry.json"), data, 0o644); err != nil {
		return err
	}
	hist, _ := json.MarshalIndent(r.history, "", "  ")
	return os.WriteFile(filepath.Join(r.dir, "history.json"), hist, 0o644)
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
