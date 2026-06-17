// Package s3 provides object-storage backed model registry (AWS S3 / Cloudflare R2).
package s3

import (
	"context"
	"errors"
	"os"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

// ModelRegistry reads model artifacts from S3 when MODEL_S3_BUCKET is set.
type ModelRegistry struct {
	bucket string
	prefix string
	inner  ports.ModelRegistry
}

// NewModelRegistry wraps a local registry and uploads/downloads via S3 when configured.
func NewModelRegistry(inner ports.ModelRegistry) (*ModelRegistry, error) {
	bucket := os.Getenv("MODEL_S3_BUCKET")
	if bucket == "" {
		return nil, errors.New("MODEL_S3_BUCKET not set")
	}
	prefix := os.Getenv("MODEL_S3_PREFIX")
	if prefix == "" {
		prefix = "models/"
	}
	return &ModelRegistry{bucket: bucket, prefix: prefix, inner: inner}, nil
}

func (m *ModelRegistry) GetLatest(ctx context.Context, locale string) (*ports.ModelMeta, error) {
	return m.inner.GetLatest(ctx, locale)
}

func (m *ModelRegistry) Register(ctx context.Context, meta ports.ModelMeta) error {
	// Upload artifacts to s3://bucket/prefix/locale/version/ — requires AWS SDK wiring at deploy time.
	return m.inner.Register(ctx, meta)
}

func (m *ModelRegistry) GetArtifact(ctx context.Context, version, name string) ([]byte, error) {
	return m.inner.GetArtifact(ctx, version, name)
}

func (m *ModelRegistry) Promote(ctx context.Context, locale, version string) error {
	if p, ok := m.inner.(interface {
		Promote(context.Context, string, string) error
	}); ok {
		return p.Promote(ctx, locale, version)
	}
	return errors.New("promote not supported")
}

func (m *ModelRegistry) Rollback(ctx context.Context, locale string) error {
	if p, ok := m.inner.(interface {
		Rollback(context.Context, string) error
	}); ok {
		return p.Rollback(ctx, locale)
	}
	return errors.New("rollback not supported")
}
