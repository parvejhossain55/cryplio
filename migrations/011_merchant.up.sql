-- ============================================
-- Migration 011: Merchant System
-- Merchant applications, subscriptions, analytics (SRS 3.6)
-- ============================================

BEGIN;

CREATE TABLE IF NOT EXISTS merchant_applications (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(user_id),
    business_name VARCHAR(100),
    registration_number VARCHAR(100),
    tax_id VARCHAR(100),
    address TEXT,
    website VARCHAR(255),
    operational_volume DECIMAL(15,2),
    applied_at TIMESTAMP NOT NULL DEFAULT NOW(),
    status merchant_status NOT NULL DEFAULT 'pending',
    reviewed_by UUID REFERENCES users(user_id),
    reviewed_at TIMESTAMP,
    rejection_reason TEXT,
    subscription_plan VARCHAR(30),
    subscription_start_at TIMESTAMP,
    subscription_end_at TIMESTAMP,
    auto_renew BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_merchant_app_user ON merchant_applications(user_id);
CREATE INDEX idx_merchant_app_status ON merchant_applications(status);

-- Merchant daily analytics (rollup)
CREATE TABLE IF NOT EXISTS merchant_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    merchant_id UUID NOT NULL REFERENCES users(user_id),
    date DATE NOT NULL,
    total_trades INT NOT NULL DEFAULT 0,
    total_volume_usd DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_fees_earned DECIMAL(15,2) NOT NULL DEFAULT 0,
    new_customers INT NOT NULL DEFAULT 0,
    avg_response_time_seconds INT,
    completion_rate DECIMAL(5,2) DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(merchant_id, date)
);

CREATE INDEX idx_merchant_analytics_merchant ON merchant_analytics(merchant_id, date DESC);

COMMIT;
