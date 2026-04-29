-- ============================================
-- Migration 010: Referrals
-- User referral tracking and commission system (SRS 3.7)
-- ============================================

BEGIN;

CREATE TABLE IF NOT EXISTS referrals (
    referral_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referrer_id UUID NOT NULL REFERENCES users(user_id),
    referee_id UUID NOT NULL REFERENCES users(user_id),
    referral_code VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    reward_level VARCHAR(20) DEFAULT 'standard',
    reward_percentage DECIMAL(5,3) NOT NULL DEFAULT 0.20, -- 20% of platform fee
    total_commission_earned DECIMAL(15,2) NOT NULL DEFAULT 0,
    lifetime_volume_usd DECIMAL(15,2) NOT NULL DEFAULT 0,
    first_trade_at TIMESTAMP,
    last_commission_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_referrals_referrer ON referrals(referrer_id);
CREATE INDEX idx_referrals_referee ON referrals(referee_id);
CREATE UNIQUE INDEX idx_referrals_code ON referrals(referral_code);
CREATE INDEX idx_referrals_status ON referrals(status);

COMMIT;
