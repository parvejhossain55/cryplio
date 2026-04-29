-- ============================================
-- Migration 008: Wallets & Transactions
-- Crypto wallet management and transaction ledger
-- ============================================

BEGIN;

-- User wallets (one per crypto per user)
CREATE TABLE IF NOT EXISTS wallets (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    crypto_id INT NOT NULL REFERENCES crypto_assets(id),
    address VARCHAR(255) UNIQUE NOT NULL,
    address_label VARCHAR(100),
    balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    locked_balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    total_deposited DECIMAL(20,8) NOT NULL DEFAULT 0,
    total_withdrawn DECIMAL(20,8) NOT NULL DEFAULT 0,
    last_deposit_at TIMESTAMP,
    last_withdrawal_at TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, crypto_id)
);

CREATE INDEX idx_wallets_user ON wallets(user_id);
CREATE INDEX idx_wallets_address ON wallets(address);
CREATE INDEX idx_wallets_crypto ON wallets(crypto_id);

-- Wallet transaction history
CREATE TABLE IF NOT EXISTS wallet_transactions (
    txn_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(wallet_id),
    user_id UUID NOT NULL REFERENCES users(user_id),
    type transaction_type NOT NULL,
    subtype VARCHAR(30),
    amount DECIMAL(20,8) NOT NULL,
    fee DECIMAL(20,8) NOT NULL DEFAULT 0,
    fee_currency CHAR(3) DEFAULT 'USD',
    crypto_id INT REFERENCES crypto_assets(id),
    network VARCHAR(20),
    status transaction_status NOT NULL DEFAULT 'pending',
    reference_id UUID,
    tx_hash VARCHAR(66) UNIQUE,
    from_address VARCHAR(255),
    to_address VARCHAR(255),
    block_number BIGINT,
    confirmations INT DEFAULT 0,
    required_confirmations INT DEFAULT 1,
    explorer_url VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    metadata JSONB
);

CREATE INDEX idx_wallet_tx_wallet ON wallet_transactions(wallet_id, created_at DESC);
CREATE INDEX idx_wallet_tx_user ON wallet_transactions(user_id);
CREATE INDEX idx_wallet_tx_type ON wallet_transactions(type);
CREATE INDEX idx_wallet_tx_status ON wallet_transactions(status);
CREATE INDEX idx_wallet_tx_txhash ON wallet_transactions(tx_hash);
CREATE INDEX idx_wallet_tx_reference ON wallet_transactions(reference_id);

-- Withdrawals
CREATE TABLE IF NOT EXISTS withdrawals (
    withdrawal_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id),
    wallet_id UUID NOT NULL REFERENCES wallets(wallet_id),
    crypto_id INT NOT NULL REFERENCES crypto_assets(id),
    address VARCHAR(255) NOT NULL,
    amount DECIMAL(20,8) NOT NULL,
    fee DECIMAL(20,8) NOT NULL DEFAULT 0,
    net_amount DECIMAL(20,8) NOT NULL,
    network VARCHAR(20) NOT NULL DEFAULT 'TRC20',
    status transaction_status NOT NULL DEFAULT 'pending',
    tx_hash VARCHAR(66),
    admin_notes TEXT,
    requires_2fa BOOLEAN NOT NULL DEFAULT true,
    requires_email BOOLEAN NOT NULL DEFAULT true,
    confirmed_at TIMESTAMP,
    completed_at TIMESTAMP,
    cancelled_at TIMESTAMP,
    cancellation_reason TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_withdrawals_user ON withdrawals(user_id);
CREATE INDEX idx_withdrawals_status ON withdrawals(status);
CREATE INDEX idx_withdrawals_txhash ON withdrawals(tx_hash);

COMMIT;
