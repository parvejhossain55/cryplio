-- ============================================
-- Migration 004: Wallets & Trading (Aligned with live DB)
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
    price_type price_type NOT NULL DEFAULT 'fixed',
    price DECIMAL(20,8) NOT NULL,
    floating_markup DECIMAL(5,2),
    min_amount DECIMAL(20,2) NOT NULL,
    max_amount DECIMAL(20,2) NOT NULL,
    payment_methods INTEGER[] NOT NULL DEFAULT '{}',
    trade_terms TEXT,
    payment_window_minutes INT NOT NULL DEFAULT 30,
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

CREATE TABLE trades (
    trade_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ad_id UUID NOT NULL REFERENCES trade_ads(ad_id),
    buyer_id UUID NOT NULL REFERENCES users(user_id),
    seller_id UUID NOT NULL REFERENCES users(user_id),
    crypto_id INT NOT NULL REFERENCES crypto_assets(id),
    fiat_id INT NOT NULL REFERENCES fiat_currencies(id),
    crypto_amount DECIMAL(20,8) NOT NULL,
    fiat_amount DECIMAL(20,2) NOT NULL,
    exchange_rate DECIMAL(20,8) NOT NULL,
    agreed_price DECIMAL(20,8) NOT NULL,
    payment_method INT NOT NULL REFERENCES payment_methods(id),
    price_type price_type NOT NULL DEFAULT 'fixed',
    status trade_status NOT NULL DEFAULT 'pending',
    dispute_id INT,
    chat_room_id VARCHAR(255) UNIQUE,
    started_at TIMESTAMP,
    payment_marked_at TIMESTAMP,
    released_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    completed_at TIMESTAMP,
    expired_at TIMESTAMP,
    payment_window_minutes INT NOT NULL DEFAULT 30,
    time_remaining_seconds INT,
    is_auto_dispute_triggered BOOLEAN NOT NULL DEFAULT false,
    cancel_reason TEXT,
    escrow_txn_hash VARCHAR(66),
    escrow_contract_address VARCHAR(42),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
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

COMMIT;
