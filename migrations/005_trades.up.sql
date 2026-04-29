-- ============================================
-- Migration 005: Trades
-- Individual trade executions between buyer and seller
-- ============================================

BEGIN;

CREATE TABLE IF NOT EXISTS trades (
    trade_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ad_id UUID NOT NULL REFERENCES trade_ads(ad_id),
    buyer_id UUID NOT NULL REFERENCES users(user_id),
    seller_id UUID NOT NULL REFERENCES users(user_id),
    crypto_id INT NOT NULL REFERENCES crypto_assets(id),
    fiat_id INT NOT NULL REFERENCES fiat_currencies(id),
    crypto_amount DECIMAL(20,8) NOT NULL,
    fiat_amount DECIMAL(20,2) NOT NULL,
    exchange_rate DECIMAL(20,8) NOT NULL,
    payment_method INT NOT NULL REFERENCES payment_methods(id),
    price_type price_type NOT NULL DEFAULT 'fixed',
    agreed_price DECIMAL(20,8) NOT NULL,
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

CREATE INDEX idx_trades_ad_id ON trades(ad_id);
CREATE INDEX idx_trades_buyer ON trades(buyer_id);
CREATE INDEX idx_trades_seller ON trades(seller_id);
CREATE INDEX idx_trades_status ON trades(status);
CREATE INDEX idx_trades_created ON trades(created_at DESC);
CREATE INDEX idx_trades_dispute ON trades(dispute_id) WHERE dispute_id IS NOT NULL;
CREATE INDEX idx_trades_active ON trades(status, created_at)
WHERE status IN ('pending', 'active', 'paid');
CREATE INDEX idx_trades_user_pairs ON trades(buyer_id, seller_id);

-- Trade chat messages
CREATE TABLE IF NOT EXISTS trade_messages (
    message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(user_id),
    message_type VARCHAR(20) NOT NULL DEFAULT 'text',
    content TEXT,
    file_url VARCHAR(500),
    file_mime_type VARCHAR(100),
    file_size INT,
    is_read BOOLEAN NOT NULL DEFAULT false,
    read_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_trade_messages_trade ON trade_messages(trade_id, created_at DESC);
CREATE INDEX idx_trade_messages_sender ON trade_messages(sender_id);

-- Trade attachments
CREATE TABLE IF NOT EXISTS trade_attachments (
    id SERIAL PRIMARY KEY,
    message_id UUID NOT NULL REFERENCES trade_messages(message_id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size_bytes BIGINT NOT NULL,
    uploader_id UUID NOT NULL REFERENCES users(user_id),
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_trade_attachments_message ON trade_attachments(message_id);

COMMIT;
