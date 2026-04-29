-- ============================================
-- Migration 004: Trade Advertisements
-- Trade ads created by makers (buy/sell offers)
-- ============================================

BEGIN;

CREATE TABLE IF NOT EXISTS trade_ads (
    ad_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id),
    type ad_type NOT NULL,
    crypto_id INT NOT NULL REFERENCES crypto_assets(id),
    fiat_id INT NOT NULL REFERENCES fiat_currencies(id),
    price_type price_type NOT NULL DEFAULT 'fixed',
    price DECIMAL(20,8) NOT NULL,
    floating_markup DECIMAL(5,2),
    min_amount DECIMAL(20,2) NOT NULL,
    max_amount DECIMAL(20,2) NOT NULL,
    payment_methods INT[] NOT NULL DEFAULT '{}',
    trade_terms TEXT,
    payment_window_minutes INT NOT NULL DEFAULT 30,
    requires_kyc_level kyc_level DEFAULT 'level_1',
    is_public BOOLEAN NOT NULL DEFAULT true,
    is_paused BOOLEAN NOT NULL DEFAULT false,
    visibility_start_at TIMESTAMP,
    visibility_end_at TIMESTAMP,
    timezone VARCHAR(50) DEFAULT 'UTC',
    auto_repost BOOLEAN NOT NULL DEFAULT false,
    repost_count INT NOT NULL DEFAULT 0,
    views_count INT NOT NULL DEFAULT 0,
    response_count INT NOT NULL DEFAULT 0,
    locked_balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    first_published_at TIMESTAMP,
    published_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Indexes
CREATE INDEX idx_ads_user_id ON trade_ads(user_id);
CREATE INDEX idx_ads_type_crypto ON trade_ads(type, crypto_id);
CREATE INDEX idx_ads_fiat ON trade_ads(fiat_id);
CREATE INDEX idx_ads_status ON trade_ads(status);
CREATE INDEX idx_ads_price ON trade_ads(price);
CREATE INDEX idx_ads_created ON trade_ads(created_at DESC);
CREATE INDEX idx_ads_payment_methods ON trade_ads USING GIN(payment_methods);
CREATE INDEX idx_ads_active ON trade_ads(type, status, is_public, crypto_id, fiat_id)
WHERE status = 'active' AND is_public = true AND deleted_at IS NULL;

COMMIT;
