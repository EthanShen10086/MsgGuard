package ports

import "context"

type ModelArtifact struct {
	Name     string `json:"name"`
	Checksum string `json:"checksum"`
	Size     int64  `json:"size"`
	Path     string `json:"path"`
}

type ModelMeta struct {
	Version   string          `json:"version"`
	Locale    string          `json:"locale"`
	Checksum  string          `json:"checksum"`
	Artifacts []ModelArtifact `json:"artifacts"`
}

type ModelRegistry interface {
	Register(ctx context.Context, meta ModelMeta) error
	GetLatest(ctx context.Context, locale string) (*ModelMeta, error)
	GetArtifact(ctx context.Context, version, name string) ([]byte, error)
}

// ModelRollback is implemented by registries that keep version history per locale.
type ModelRollback interface {
	Rollback(ctx context.Context, locale, targetVersion string) (*ModelMeta, error)
}
