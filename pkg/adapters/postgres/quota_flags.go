package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/EthanShen10086/voxera-kit/aiquota"
	aiquotaMemory "github.com/EthanShen10086/voxera-kit/aiquota/memory"
	"github.com/EthanShen10086/voxera-kit/featureflag"
)

// QuotaManager persists whitelist entries to postgres while delegating metering to memory.
type QuotaManager struct {
	inner *aiquotaMemory.Store
	db    *sql.DB
}

func NewQuotaManager(dsn string) (*QuotaManager, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	m := &QuotaManager{inner: aiquotaMemory.NewStore(), db: db}
	if err := m.loadWhitelist(context.Background()); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *QuotaManager) loadWhitelist(ctx context.Context) error {
	rows, err := m.db.QueryContext(ctx, `SELECT user_id, reason, granted_by, granted_at FROM quota_whitelist`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var entry aiquota.WhitelistEntry
		if err := rows.Scan(&entry.UserID, &entry.Reason, &entry.GrantedBy, &entry.GrantedAt); err != nil {
			return err
		}
		_ = m.inner.AddToWhitelist(ctx, entry)
	}
	return rows.Err()
}

func (m *QuotaManager) CheckQuota(ctx context.Context, userID, model string, estimatedTokens int) error {
	return m.inner.CheckQuota(ctx, userID, model, estimatedTokens)
}

func (m *QuotaManager) RecordUsage(ctx context.Context, record aiquota.UsageRecord) error {
	return m.inner.RecordUsage(ctx, record)
}

func (m *QuotaManager) GetUsage(ctx context.Context, userID string) (*aiquota.Usage, error) {
	return m.inner.GetUsage(ctx, userID)
}

func (m *QuotaManager) GetQuota(ctx context.Context, userID string) (*aiquota.Quota, error) {
	return m.inner.GetQuota(ctx, userID)
}

func (m *QuotaManager) SetTier(ctx context.Context, userID string, tier aiquota.Tier) error {
	return m.inner.SetTier(ctx, userID, tier)
}

func (m *QuotaManager) IsWhitelisted(ctx context.Context, userID string) (bool, error) {
	return m.inner.IsWhitelisted(ctx, userID)
}

func (m *QuotaManager) AddToWhitelist(ctx context.Context, entry aiquota.WhitelistEntry) error {
	if entry.GrantedAt.IsZero() {
		entry.GrantedAt = time.Now().UTC()
	}
	_, err := m.db.ExecContext(ctx, `
		INSERT INTO quota_whitelist (user_id, reason, granted_by, granted_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET reason = EXCLUDED.reason, granted_by = EXCLUDED.granted_by, granted_at = EXCLUDED.granted_at`,
		entry.UserID, entry.Reason, entry.GrantedBy, entry.GrantedAt,
	)
	if err != nil {
		return err
	}
	return m.inner.AddToWhitelist(ctx, entry)
}

func (m *QuotaManager) RemoveFromWhitelist(ctx context.Context, userID string) error {
	if _, err := m.db.ExecContext(ctx, `DELETE FROM quota_whitelist WHERE user_id = $1`, userID); err != nil {
		return err
	}
	return m.inner.RemoveFromWhitelist(ctx, userID)
}

func (m *QuotaManager) ListWhitelist(ctx context.Context) ([]aiquota.WhitelistEntry, error) {
	return m.inner.ListWhitelist(ctx)
}

func (m *QuotaManager) GetCostReport(ctx context.Context, tenantID string, from, to time.Time) (*aiquota.CostReport, error) {
	return m.inner.GetCostReport(ctx, tenantID, from, to)
}

func (m *QuotaManager) AcquireConcurrency(ctx context.Context, userID string) (func(), error) {
	return m.inner.AcquireConcurrency(ctx, userID)
}

// FlagStore persists feature flags to postgres.
type FlagStore struct {
	mu    sync.RWMutex
	db    *sql.DB
	cache map[string]featureflag.Flag
}

func NewFlagStore(dsn string) (*FlagStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	s := &FlagStore{db: db, cache: make(map[string]featureflag.Flag)}
	if err := s.load(context.Background()); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *FlagStore) load(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `SELECT key, enabled, percentage, allow_list, deny_list FROM feature_flags`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var flag featureflag.Flag
		var allowRaw, denyRaw []byte
		if err := rows.Scan(&flag.Key, &flag.Enabled, &flag.Percentage, &allowRaw, &denyRaw); err != nil {
			return err
		}
		_ = json.Unmarshal(allowRaw, &flag.AllowList)
		_ = json.Unmarshal(denyRaw, &flag.DenyList)
		s.cache[flag.Key] = flag
	}
	return rows.Err()
}

func (s *FlagStore) IsEnabled(ctx context.Context, key string, evalCtx featureflag.EvalContext) (bool, error) {
	s.mu.RLock()
	flag, ok := s.cache[key]
	s.mu.RUnlock()
	if !ok || !flag.Enabled {
		return false, nil
	}
	for _, denied := range flag.DenyList {
		if denied == evalCtx.UserID {
			return false, nil
		}
	}
	for _, allowed := range flag.AllowList {
		if allowed == evalCtx.UserID {
			return true, nil
		}
	}
	if flag.Percentage <= 0 {
		return false, nil
	}
	if flag.Percentage >= 100 {
		return true, nil
	}
	bucket := hashBucket(evalCtx.UserID, key)
	return bucket < flag.Percentage, nil
}

func (s *FlagStore) GetFlags(ctx context.Context) ([]featureflag.Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	flags := make([]featureflag.Flag, 0, len(s.cache))
	for _, f := range s.cache {
		flags = append(flags, f)
	}
	return flags, nil
}

func (s *FlagStore) SetFlag(ctx context.Context, flag featureflag.Flag) error {
	allowRaw, _ := json.Marshal(flag.AllowList)
	denyRaw, _ := json.Marshal(flag.DenyList)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO feature_flags (key, enabled, percentage, allow_list, deny_list, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (key) DO UPDATE SET
			enabled = EXCLUDED.enabled,
			percentage = EXCLUDED.percentage,
			allow_list = EXCLUDED.allow_list,
			deny_list = EXCLUDED.deny_list,
			updated_at = NOW()`,
		flag.Key, flag.Enabled, flag.Percentage, allowRaw, denyRaw,
	)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.cache[flag.Key] = flag
	s.mu.Unlock()
	return nil
}

func hashBucket(userID, key string) float64 {
	h := sha256.New()
	h.Write([]byte(userID + ":" + key))
	sum := h.Sum(nil)
	val := binary.BigEndian.Uint32(sum[:4])
	return float64(val) / float64(^uint32(0)) * 100
}
