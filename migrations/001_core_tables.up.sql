-- ============================================
-- Migration 001: Core User & Auth Tables
-- User accounts, sessions, password resets, email verification
-- ============================================

BEGIN;

-- Users table (core identity)
CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(30) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    phone_country_code VARCHAR(5),
    phone_number VARCHAR(20),
    phone_verified BOOLEAN NOT NULL DEFAULT false,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    kyc_level kyc_level NOT NULL DEFAULT 'level_0',
    kyc_last_updated TIMESTAMP,
    status user_status NOT NULL DEFAULT 'pending',
    avatar_url VARCHAR(500),
    bio VARCHAR(200),
    timezone VARCHAR(50) DEFAULT 'UTC',
    locale VARCHAR(10) DEFAULT 'en',
    is_merchant BOOLEAN NOT NULL DEFAULT false,
    is_suspended BOOLEAN NOT NULL DEFAULT false,
    suspension_reason TEXT,
    suspended_at TIMESTAMP,
    suspended_until TIMESTAMP,
    last_login_at TIMESTAMP,
    login_count INT NOT NULL DEFAULT 0,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    referral_code VARCHAR(15) UNIQUE,
    referred_by UUID REFERENCES users(user_id),
    two_fa_secret VARCHAR(255),
    remember_2fa BOOLEAN NOT NULL DEFAULT false,
    remember_2fa_expiry TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- User statistics (denormalized aggregates)
CREATE TABLE IF NOT EXISTS user_stats (
    user_id UUID PRIMARY KEY REFERENCES users(user_id) ON DELETE CASCADE,
    total_trades INT NOT NULL DEFAULT 0,
    successful_trades INT NOT NULL DEFAULT 0,
    dispute_rate DECIMAL(5,2) DEFAULT 0,
    avg_rating DECIMAL(3,2),
    positive_feedback_count INT NOT NULL DEFAULT 0,
    neutral_feedback_count INT NOT NULL DEFAULT 0,
    negative_feedback_count INT NOT NULL DEFAULT 0,
    total_volume_usd DECIMAL(15,2) NOT NULL DEFAULT 0,
    last_trade_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- User sessions (device management)
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token_id VARCHAR(255) UNIQUE NOT NULL,
    device_fingerprint VARCHAR(255),
    ip_address INET,
    user_agent TEXT,
    device_type VARCHAR(50),
    location VARCHAR(100),
    is_remembered BOOLEAN NOT NULL DEFAULT false,
    expires_at TIMESTAMP NOT NULL,
    last_used_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Password reset tokens
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    ip_address INET,
    expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '15 minutes'),
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Email verification tokens
CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '24 hours'),
    verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- User blocked users
CREATE TABLE IF NOT EXISTS user_blocks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    blocker_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    reason VARCHAR(200),
    is_permanent BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(blocker_id, blocked_id)
);

-- Indexes for user tables
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_kyc_level ON users(kyc_level);
CREATE INDEX idx_users_referral_code ON users(referral_code);
CREATE INDEX idx_users_referred_by ON users(referred_by);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token_id ON user_sessions(token_id);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

CREATE INDEX idx_password_reset_token ON password_reset_tokens(token_hash);
CREATE INDEX idx_password_reset_user ON password_reset_tokens(user_id);

CREATE INDEX idx_email_verification_token ON email_verification_tokens(token_hash);

CREATE INDEX idx_user_blocks_blocker ON user_blocks(blocker_id);
CREATE INDEX idx_user_blocks_blocked ON user_blocks(blocked_id);

COMMIT;
