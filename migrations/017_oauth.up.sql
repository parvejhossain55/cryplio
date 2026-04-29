-- ============================================
-- Migration 017: OAuth Accounts
-- Link external OAuth providers (Google, etc.) to user accounts
-- ============================================

BEGIN;

-- OAuth accounts table
CREATE TABLE IF NOT EXISTS user_oauth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- 'google', 'github', etc.
    provider_user_id VARCHAR(255) NOT NULL, -- OAuth provider's user ID
    provider_email VARCHAR(255),
    provider_username VARCHAR(255),
    access_token TEXT,
    refresh_token TEXT,
    token_expiry TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(provider, provider_user_id),
    UNIQUE(user_id, provider)
);

-- Indexes for OAuth lookups
CREATE INDEX idx_user_oauth_user_id ON user_oauth(user_id);
CREATE INDEX idx_user_oauth_provider ON user_oauth(provider, provider_user_id);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger for user_oauth
CREATE TRIGGER update_user_oauth_updated_at
    BEFORE UPDATE ON user_oauth
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

COMMIT;
