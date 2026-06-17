package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type FeedbackStore struct {
	db *sql.DB
}

func NewFeedbackStore(dsn string) (*FeedbackStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	s := &FeedbackStore{db: db}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS feedback (
		id TEXT PRIMARY KEY, body TEXT, label TEXT, locale TEXT,
		tenant_id TEXT, trace_id TEXT, created_at TIMESTAMPTZ DEFAULT NOW()
	)`)
	if err == nil {
		_, _ = db.Exec(`ALTER TABLE feedback ADD COLUMN IF NOT EXISTS tenant_id TEXT`)
	}
	return s, err
}

func (s *FeedbackStore) Create(ctx context.Context, item ports.FeedbackItem) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO feedback (id, body, label, locale, tenant_id, trace_id, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		item.ID, item.Body, item.Label, item.Locale, nullIfEmpty(item.TenantID), item.TraceID, item.CreatedAt,
	)
	return err
}

func (s *FeedbackStore) List(ctx context.Context, limit int) ([]ports.FeedbackItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, body, label, locale, tenant_id, trace_id, created_at FROM feedback ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ports.FeedbackItem
	for rows.Next() {
		var item ports.FeedbackItem
		var tenantID sql.NullString
		if err := rows.Scan(&item.ID, &item.Body, &item.Label, &item.Locale, &tenantID, &item.TraceID, &item.CreatedAt); err != nil {
			return nil, err
		}
		if tenantID.Valid {
			item.TenantID = tenantID.String
		}
		out = append(out, item)
	}
	return out, nil
}

type RuleStore struct {
	db *sql.DB
}

func NewRuleStore(dsn string) (*RuleStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	s := &RuleStore{db: db}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS rules (
		version TEXT PRIMARY KEY, checksum TEXT, payload JSONB, created_at TIMESTAMPTZ DEFAULT NOW()
	)`)
	return s, err
}

func (s *RuleStore) GetLatest(ctx context.Context) (*ports.RuleBundle, error) {
	row := s.db.QueryRowContext(ctx, `SELECT version, checksum, payload FROM rules ORDER BY created_at DESC LIMIT 1`)
	var version, checksum string
	var payload []byte
	if err := row.Scan(&version, &checksum, &payload); err != nil {
		if err == sql.ErrNoRows {
			return &ports.RuleBundle{Version: "1.0.0", Checksum: "seed", Payload: []byte(`{"keywords":["免费"]}`)}, nil
		}
		return nil, err
	}
	return &ports.RuleBundle{Version: version, Checksum: checksum, Payload: payload}, nil
}

func (s *RuleStore) GetByVersion(ctx context.Context, version string) (*ports.RuleBundle, error) {
	row := s.db.QueryRowContext(ctx, `SELECT version, checksum, payload FROM rules WHERE version=$1`, version)
	var v, checksum string
	var payload []byte
	if err := row.Scan(&v, &checksum, &payload); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("version not found")
		}
		return nil, err
	}
	return &ports.RuleBundle{Version: v, Checksum: checksum, Payload: payload}, nil
}

func (s *RuleStore) Save(ctx context.Context, bundle ports.RuleBundle) error {
	if !json.Valid(bundle.Payload) {
		bundle.Payload = []byte(`{}`)
	}
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO rules (version, checksum, payload) VALUES ($1,$2,$3) ON CONFLICT (version) DO UPDATE SET checksum=$2, payload=$3`,
		bundle.Version, bundle.Checksum, bundle.Payload,
	)
	return err
}

func (s *FeedbackStore) Close() error { return s.db.Close() }

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func NowUTC() time.Time { return time.Now().UTC() }
