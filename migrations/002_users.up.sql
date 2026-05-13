-- ============================================
-- Migration 002: Identity & Access
-- ============================================

BEGIN;

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(30) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    status user_status NOT NULL DEFAULT 'active',
    phone_country_code VARCHAR(5),
    phone_number VARCHAR(20),
    phone_verified BOOLEAN NOT NULL DEFAULT false,
    email_verified BOOLEAN NOT NULL DEFAULT false,
    avatar_url VARCHAR(500),
    bio VARCHAR(200),
    timezone VARCHAR(50) DEFAULT 'UTC',
    locale VARCHAR(10) DEFAULT 'en',
    is_suspended BOOLEAN NOT NULL DEFAULT false,
    suspension_reason TEXT,
    suspended_at TIMESTAMP,
    suspended_until TIMESTAMP,
    last_login_at TIMESTAMP,
    last_seen_at TIMESTAMP,
    login_count INT NOT NULL DEFAULT 0,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    two_fa_secret VARCHAR(255),
    remember_2fa BOOLEAN NOT NULL DEFAULT false,
    remember_2fa_expiry TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE user_stats (
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

CREATE TABLE user_sessions (
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

CREATE TABLE password_reset_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    ip_address INET,
    expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '15 minutes'),
    used_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE email_verification_tokens (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL DEFAULT (NOW() + INTERVAL '24 hours'),
    verified_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_oauth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
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

CREATE TABLE two_factor_pending (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL UNIQUE REFERENCES users(user_id) ON DELETE CASCADE,
    secret VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL
);

COMMIT;
