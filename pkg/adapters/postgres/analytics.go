package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type AnalyticsStore struct {
	db *sql.DB
}

func NewAnalyticsStore(dsn string) (*AnalyticsStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	s := &AnalyticsStore{db: db}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS analytics_events (
		id TEXT PRIMARY KEY, name TEXT, props JSONB, device_id TEXT,
		trace_id TEXT, created_at TIMESTAMPTZ DEFAULT NOW()
	)`)
	return s, err
}

func (s *AnalyticsStore) Insert(ctx context.Context, event ports.AnalyticsEvent) error {
	props, _ := json.Marshal(event.Props)
	ts := event.Timestamp
	if ts.IsZero() {
		ts = time.Now().UTC()
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO analytics_events (id, name, props, device_id, trace_id, created_at) VALUES ($1,$2,$3,$4,$5,$6)`,
		event.ID, event.Name, props, event.DeviceID, event.TraceID, ts,
	)
	return err
}

func (s *AnalyticsStore) List(ctx context.Context, since time.Time, limit int) ([]ports.AnalyticsEvent, error) {
	q := `SELECT id, name, props, device_id, trace_id, created_at FROM analytics_events WHERE created_at >= $1 ORDER BY created_at DESC LIMIT $2`
	if since.IsZero() {
		q = `SELECT id, name, props, device_id, trace_id, created_at FROM analytics_events ORDER BY created_at DESC LIMIT $1`
		rows, err := s.db.QueryContext(ctx, q, limit)
		return scanAnalytics(rows, err)
	}
	rows, err := s.db.QueryContext(ctx, q, since, limit)
	return scanAnalytics(rows, err)
}

func (s *AnalyticsStore) CountByName(ctx context.Context, since time.Time) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT name, COUNT(*) FROM analytics_events WHERE created_at >= $1 GROUP BY name`, since,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := map[string]int{}
	for rows.Next() {
		var name string
		var n int
		if err := rows.Scan(&name, &n); err != nil {
			return nil, err
		}
		counts[name] = n
	}
	return counts, nil
}

func scanAnalytics(rows *sql.Rows, err error) ([]ports.AnalyticsEvent, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ports.AnalyticsEvent
	for rows.Next() {
		var e ports.AnalyticsEvent
		var props []byte
		if err := rows.Scan(&e.ID, &e.Name, &props, &e.DeviceID, &e.TraceID, &e.Timestamp); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(props, &e.Props)
		out = append(out, e)
	}
	return out, nil
}
