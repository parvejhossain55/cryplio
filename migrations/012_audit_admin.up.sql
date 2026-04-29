-- ============================================
-- Migration 012: Audit Logs & Admin Actions
-- Full audit trail for compliance and admin activity
-- ============================================

BEGIN;

-- Generic audit log (for all tables)
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name VARCHAR(50) NOT NULL,
    record_id VARCHAR(255) NOT NULL,
    action VARCHAR(20) NOT NULL,
    old_values JSONB,
    new_values JSONB,
    changed_fields JSONB,
    user_id UUID REFERENCES users(user_id),
    user_ip INET,
    user_agent TEXT,
    reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_logs_table ON audit_logs(table_name, record_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_created ON audit_logs(created_at DESC);

-- Admin action log (specific admin operations)
CREATE TABLE IF NOT EXISTS admin_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    admin_id UUID NOT NULL REFERENCES users(user_id),
    action_type admin_action_type NOT NULL,
    target_user_id UUID REFERENCES users(user_id),
    target_record_id VARCHAR(255),
    description TEXT NOT NULL,
    metadata JSONB DEFAULT '{}',
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_admin_actions_admin ON admin_actions(admin_id);
CREATE INDEX idx_admin_actions_target ON admin_actions(target_user_id);
CREATE INDEX idx_admin_actions_created ON admin_actions(created_at DESC);

-- Platform configuration (key-value store)
CREATE TABLE IF NOT EXISTS platform_config (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    value_type VARCHAR(20) NOT NULL DEFAULT 'string',
    description TEXT,
    is_public BOOLEAN NOT NULL DEFAULT false,
    updated_by UUID REFERENCES users(user_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_platform_config_key ON platform_config(key);

-- System announcements (SRS FR-806)
CREATE TABLE IF NOT EXISTS announcements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(200) NOT NULL,
    message TEXT NOT NULL,
    type VARCHAR(20) DEFAULT 'info',
    target_audience VARCHAR(20) DEFAULT 'all',
    is_active BOOLEAN NOT NULL DEFAULT true,
    published_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    created_by UUID NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_announcements_active ON announcements(is_active, published_at);

COMMIT;
