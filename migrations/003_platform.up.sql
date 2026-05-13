-- ============================================
-- Migration 003: Platform Settings
-- ============================================

BEGIN;

CREATE TABLE crypto_assets (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    name VARCHAR(50) NOT NULL,
    blockchain VARCHAR(20) NOT NULL,
    contract_address VARCHAR(42),
    decimals INT NOT NULL DEFAULT 18,
    min_confirmation INT NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(symbol, blockchain)
);

CREATE TABLE fiat_currencies (
    id SERIAL PRIMARY KEY,
    code CHAR(3) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    symbol VARCHAR(5) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE payment_methods (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(30) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    category payment_category NOT NULL,
    icon_url VARCHAR(255),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE user_payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    payment_method_code VARCHAR(30) NOT NULL REFERENCES payment_methods(code),
    display_name VARCHAR(100) NOT NULL,
    account_name VARCHAR(100),
    account_number VARCHAR(100),
    bank_name VARCHAR(100),
    metadata JSONB,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE dispute_reasons (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(30) NOT NULL UNIQUE,
    label VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(30) NOT NULL
);

CREATE TABLE fee_tiers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    min_volume_usd DECIMAL(15,2) NOT NULL,
    max_volume_usd DECIMAL(15,2),
    fee_percentage DECIMAL(5,3) NOT NULL,
    fee_minimum DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

COMMIT;
