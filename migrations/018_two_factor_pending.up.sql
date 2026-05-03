-- ============================================
-- Migration 018: Two-Factor Authentication Pending
-- Temporary storage for 2FA setup before confirmation
-- ============================================

BEGIN;

CREATE TABLE IF NOT EXISTS two_factor_pending (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    secret VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_two_factor_pending_user ON two_factor_pending(user_id);
CREATE INDEX idx_two_factor_pending_expiry ON two_factor_pending(expires_at);

COMMIT;
