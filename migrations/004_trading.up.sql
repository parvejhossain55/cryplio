-- ============================================
-- Migration 004: Wallets & Trading
-- ============================================

BEGIN;

CREATE TABLE wallets (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    crypto_id INT REFERENCES crypto_assets(id),
    address VARCHAR(255) NOT NULL,
    address_label VARCHAR(100),
    balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    locked_balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, crypto_id),
    UNIQUE(address)
);

CREATE TABLE trade_ads (
    ad_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type ad_type NOT NULL,
    crypto_id INT NOT NULL REFERENCES crypto_assets(id),
    fiat_id INT NOT NULL REFERENCES fiat_currencies(id),
    price_type price_type NOT NULL,
    price DECIMAL(20,8) NOT NULL,
    min_amount DECIMAL(20,8) NOT NULL,
    max_amount DECIMAL(20,8) NOT NULL,
    payment_method_code VARCHAR(30) NOT NULL REFERENCES payment_methods(code),
    payment_window_minutes INT NOT NULL DEFAULT 15,
    terms TEXT,
    instructions TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE trades (
    trade_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ad_id UUID NOT NULL REFERENCES trade_ads(ad_id),
    buyer_id UUID NOT NULL REFERENCES users(user_id),
    seller_id UUID NOT NULL REFERENCES users(user_id),
    crypto_id INT NOT NULL REFERENCES crypto_assets(id),
    fiat_id INT NOT NULL REFERENCES fiat_currencies(id),
    crypto_amount DECIMAL(20,8) NOT NULL,
    fiat_amount DECIMAL(20,8) NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    status trade_status NOT NULL DEFAULT 'pending',
    payment_method_code VARCHAR(30) NOT NULL REFERENCES payment_methods(code),
    escrow_wallet_id UUID REFERENCES wallets(wallet_id),
    tx_hash VARCHAR(66),
    payment_window_minutes INT NOT NULL DEFAULT 15,
    expires_at TIMESTAMP NOT NULL,
    paid_at TIMESTAMP,
    released_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    disputed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(wallet_id),
    user_id UUID NOT NULL REFERENCES users(user_id),
    type transaction_type NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',
    amount DECIMAL(20,8) NOT NULL,
    fee DECIMAL(20,8) NOT NULL DEFAULT 0,
    tx_hash VARCHAR(66),
    reference_id UUID,
    from_address VARCHAR(255),
    to_address VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE disputes (
    dispute_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    raiser_id UUID NOT NULL REFERENCES users(user_id),
    reason_code VARCHAR(30) NOT NULL REFERENCES dispute_reasons(code),
    description TEXT NOT NULL,
    status dispute_status NOT NULL DEFAULT 'pending',
    resolution dispute_resolution,
    resolved_by UUID REFERENCES users(user_id),
    resolved_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE trade_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(user_id),
    message TEXT,
    is_system BOOLEAN NOT NULL DEFAULT false,
    attachment_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE trade_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(user_id),
    reviewee_id UUID NOT NULL REFERENCES users(user_id),
    rating feedback_rating NOT NULL,
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(trade_id, reviewer_id)
);

COMMIT;
