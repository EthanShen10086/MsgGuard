package ports

import "context"

type RuleBundle struct {
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
	Payload  []byte `json:"-"`
}

type RuleStore interface {
	GetLatest(ctx context.Context) (*RuleBundle, error)
	GetByVersion(ctx context.Context, version string) (*RuleBundle, error)
	Save(ctx context.Context, bundle RuleBundle) error
}
