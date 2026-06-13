package ports

import (
	"context"
	"time"
)

type AnalyticsEvent struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Props     map[string]any `json:"props"`
	DeviceID  string         `json:"device_id"`
	TraceID   string         `json:"trace_id"`
	Timestamp time.Time      `json:"timestamp"`
}

type AnalyticsStore interface {
	Insert(ctx context.Context, event AnalyticsEvent) error
	List(ctx context.Context, since time.Time, limit int) ([]AnalyticsEvent, error)
	CountByName(ctx context.Context, since time.Time) (map[string]int, error)
}
