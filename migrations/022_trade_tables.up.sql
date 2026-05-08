-- ============================================
-- Migration 002: Trade & Marketplace Tables
-- Trade ads, trades, wallets, transactions, disputes, messages, feedback, notifications
-- ============================================

BEGIN;

-- Trade advertisements (marketplace listings)
CREATE TABLE IF NOT EXISTS trade_ads (
    ad_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type VARCHAR(10) NOT NULL CHECK (type IN ('buy', 'sell')),
    crypto_symbol VARCHAR(10) NOT NULL DEFAULT 'USDT',
    fiat_symbol VARCHAR(10) NOT NULL,
    price_type VARCHAR(10) NOT NULL CHECK (price_type IN ('fixed', 'floating')),
    price DECIMAL(20,8) NOT NULL,
    min_amount DECIMAL(20,8) NOT NULL,
    max_amount DECIMAL(20,8) NOT NULL,
    payment_methods TEXT[] NOT NULL,
    payment_window_minutes INT NOT NULL DEFAULT 15,
    instructions TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused', 'closed')),
    trade_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for trade_ads
CREATE INDEX IF NOT EXISTS idx_trade_ads_user_id ON trade_ads(user_id);
CREATE INDEX IF NOT EXISTS idx_trade_ads_type ON trade_ads(type);
CREATE INDEX IF NOT EXISTS idx_trade_ads_crypto_fiat ON trade_ads(crypto_symbol, fiat_symbol);
CREATE INDEX IF NOT EXISTS idx_trade_ads_status ON trade_ads(status);
CREATE INDEX IF NOT EXISTS idx_trade_ads_price ON trade_ads(price);
CREATE INDEX IF NOT EXISTS idx_trade_ads_created_at ON trade_ads(created_at);

-- Trades (executed trades between users)
CREATE TABLE IF NOT EXISTS trades (
    trade_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ad_id UUID NOT NULL REFERENCES trade_ads(ad_id) ON DELETE CASCADE,
    buyer_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    seller_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    crypto_amount DECIMAL(20,8) NOT NULL,
    fiat_amount DECIMAL(20,8) NOT NULL,
    crypto_symbol VARCHAR(10) NOT NULL DEFAULT 'USDT',
    fiat_symbol VARCHAR(10) NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'paid', 'released', 'completed', 'cancelled', 'disputed')),
    type VARCHAR(10) NOT NULL CHECK (type IN ('buy', 'sell')),
    payment_method VARCHAR(50),
    payment_details JSONB,
    escrow_id VARCHAR(100),
    blockchain_tx_hash VARCHAR(100),
    payment_window_minutes INT NOT NULL DEFAULT 15,
    timer_expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for trades
CREATE INDEX IF NOT EXISTS idx_trades_ad_id ON trades(ad_id);
CREATE INDEX IF NOT EXISTS idx_trades_buyer_id ON trades(buyer_id);
CREATE INDEX IF NOT EXISTS idx_trades_seller_id ON trades(seller_id);
CREATE INDEX IF NOT EXISTS idx_trades_status ON trades(status);
CREATE INDEX IF NOT EXISTS idx_trades_created_at ON trades(created_at);
CREATE INDEX IF NOT EXISTS idx_trades_timer_expires ON trades(timer_expires_at);

-- User cryptocurrency wallets
CREATE TABLE IF NOT EXISTS wallets (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    crypto_symbol VARCHAR(10) NOT NULL,
    blockchain_address VARCHAR(100) NOT NULL UNIQUE,
    private_key_encrypted TEXT, -- Encrypted private key
    balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    locked_balance DECIMAL(20,8) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, crypto_symbol)
);

-- Create indexes for wallets
CREATE INDEX IF NOT EXISTS idx_wallets_user_id ON wallets(user_id);
CREATE INDEX IF NOT EXISTS idx_wallets_crypto_symbol ON wallets(crypto_symbol);
CREATE INDEX IF NOT EXISTS idx_wallets_address ON wallets(blockchain_address);

-- Wallet transactions (deposits, withdrawals, escrow locks/unlocks)
CREATE TABLE IF NOT EXISTS wallet_transactions (
    tx_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID NOT NULL REFERENCES wallets(wallet_id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('deposit', 'withdrawal', 'escrow_lock', 'escrow_release', 'escrow_refund')),
    amount DECIMAL(20,8) NOT NULL,
    balance_before DECIMAL(20,8) NOT NULL,
    balance_after DECIMAL(20,8) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'failed', 'cancelled')),
    blockchain_tx_hash VARCHAR(100),
    to_address VARCHAR(100),
    from_address VARCHAR(100),
    network_fee DECIMAL(20,8),
    trade_id UUID REFERENCES trades(trade_id) ON DELETE SET NULL,
    withdrawal_request_id UUID,
    two_fa_verified BOOLEAN NOT NULL DEFAULT false,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for wallet_transactions
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_wallet_id ON wallet_transactions(wallet_id);
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_type ON wallet_transactions(type);
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_status ON wallet_transactions(status);
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_trade_id ON wallet_transactions(trade_id);
CREATE INDEX IF NOT EXISTS idx_wallet_transactions_created_at ON wallet_transactions(created_at);

-- Disputes
CREATE TABLE IF NOT EXISTS disputes (
    dispute_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    raiser_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    respondent_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    reason VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'investigating', 'resolved', 'closed')),
    resolution VARCHAR(50),
    resolution_notes TEXT,
    resolved_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
    resolved_at TIMESTAMP,
    evidence_files TEXT[], -- Array of file URLs
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for disputes
CREATE INDEX IF NOT EXISTS idx_disputes_trade_id ON disputes(trade_id);
CREATE INDEX IF NOT EXISTS idx_disputes_raiser_id ON disputes(raiser_id);
CREATE INDEX IF NOT EXISTS idx_disputes_status ON disputes(status);
CREATE INDEX IF NOT EXISTS idx_disputes_created_at ON disputes(created_at);

-- Trade chat messages
CREATE TABLE IF NOT EXISTS trade_messages (
    message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    content TEXT,
    file_url VARCHAR(500),
    file_name VARCHAR(255),
    file_size BIGINT,
    mime_type VARCHAR(100),
    message_type VARCHAR(20) NOT NULL DEFAULT 'text' CHECK (message_type IN ('text', 'file', 'system')),
    is_edited BOOLEAN NOT NULL DEFAULT false,
    edited_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for trade_messages
CREATE INDEX IF NOT EXISTS idx_trade_messages_trade_id ON trade_messages(trade_id);
CREATE INDEX IF NOT EXISTS idx_trade_messages_sender_id ON trade_messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_trade_messages_created_at ON trade_messages(created_at);

-- User feedback/ratings after trades
CREATE TABLE IF NOT EXISTS trade_feedback (
    feedback_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    reviewer_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    reviewee_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    feedback_type VARCHAR(20) NOT NULL CHECK (feedback_type IN ('buyer_to_seller', 'seller_to_buyer')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(trade_id, reviewer_id, feedback_type)
);

-- Create indexes for trade_feedback
CREATE INDEX IF NOT EXISTS idx_trade_feedback_trade_id ON trade_feedback(trade_id);
CREATE INDEX IF NOT EXISTS idx_trade_feedback_reviewer_id ON trade_feedback(reviewer_id);
CREATE INDEX IF NOT EXISTS idx_trade_feedback_reviewee_id ON trade_feedback(reviewee_id);
CREATE INDEX IF NOT EXISTS idx_trade_feedback_rating ON trade_feedback(rating);

-- User notifications
CREATE TABLE IF NOT EXISTS notifications (
    notification_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB, -- Additional data related to the notification
    is_read BOOLEAN NOT NULL DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for notifications
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- Withdrawal requests (for admin approval)
CREATE TABLE IF NOT EXISTS withdrawal_requests (
    request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    wallet_id UUID NOT NULL REFERENCES wallets(wallet_id) ON DELETE CASCADE,
    amount DECIMAL(20,8) NOT NULL,
    to_address VARCHAR(100) NOT NULL,
    blockchain_tx_hash VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'processing', 'completed', 'failed')),
    approved_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
    approved_at TIMESTAMP,
    rejected_by UUID REFERENCES users(user_id) ON DELETE SET NULL,
    rejected_at TIMESTAMP,
    rejection_reason TEXT,
    two_fa_code VARCHAR(10),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes for withdrawal_requests
CREATE INDEX IF NOT EXISTS idx_withdrawal_requests_user_id ON withdrawal_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_withdrawal_requests_status ON withdrawal_requests(status);
CREATE INDEX IF NOT EXISTS idx_withdrawal_requests_created_at ON withdrawal_requests(created_at);

-- Update the updated_at trigger function if it doesn't exist
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_trade_ads_updated_at BEFORE UPDATE ON trade_ads FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_trades_updated_at BEFORE UPDATE ON trades FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_wallet_transactions_updated_at BEFORE UPDATE ON wallet_transactions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_disputes_updated_at BEFORE UPDATE ON disputes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_withdrawal_requests_updated_at BEFORE UPDATE ON withdrawal_requests FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMIT;
