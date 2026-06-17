-- Optional multi-tenant column for enterprise deployments.
ALTER TABLE feedback ADD COLUMN IF NOT EXISTS tenant_id TEXT;
ALTER TABLE analytics_events ADD COLUMN IF NOT EXISTS tenant_id TEXT;
CREATE INDEX IF NOT EXISTS idx_feedback_tenant ON feedback (tenant_id) WHERE tenant_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_analytics_tenant ON analytics_events (tenant_id) WHERE tenant_id IS NOT NULL;
