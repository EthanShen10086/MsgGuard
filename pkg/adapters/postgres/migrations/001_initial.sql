-- MsgGuard initial schema (001)

CREATE TABLE IF NOT EXISTS feedback (
    id TEXT PRIMARY KEY,
    body TEXT,
    label TEXT,
    locale TEXT,
    trace_id TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS rules (
    version TEXT PRIMARY KEY,
    checksum TEXT,
    payload JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS analytics_events (
    id TEXT PRIMARY KEY,
    name TEXT,
    props JSONB,
    device_id TEXT,
    trace_id TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subscriptions (
    device_id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL DEFAULT '',
    signed_transaction TEXT,
    is_pro BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS quota_whitelist (
    user_id TEXT PRIMARY KEY,
    reason TEXT,
    granted_by TEXT,
    granted_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS feature_flags (
    key TEXT PRIMARY KEY,
    enabled BOOLEAN NOT NULL DEFAULT false,
    percentage DOUBLE PRECISION NOT NULL DEFAULT 0,
    allow_list JSONB NOT NULL DEFAULT '[]',
    deny_list JSONB NOT NULL DEFAULT '[]',
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS audit_log (
    id TEXT PRIMARY KEY,
    actor_id TEXT,
    action TEXT,
    resource_type TEXT,
    resource_id TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analytics_events_created_at ON analytics_events (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_feedback_created_at ON feedback (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON audit_log (created_at DESC);
