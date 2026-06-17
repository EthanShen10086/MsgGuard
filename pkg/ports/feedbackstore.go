package ports

import (
	"context"
	"time"
)

type FeedbackItem struct {
	ID        string    `json:"id"`
	Body      string    `json:"body"`
	Label     string    `json:"label"`
	Locale    string    `json:"locale"`
	TenantID  string    `json:"tenant_id,omitempty"`
	TraceID   string    `json:"trace_id"`
	CreatedAt time.Time `json:"created_at"`
}

type FeedbackStore interface {
	Create(ctx context.Context, item FeedbackItem) error
	List(ctx context.Context, limit int) ([]FeedbackItem, error)
}
