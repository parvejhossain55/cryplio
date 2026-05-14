-- ============================================
-- Migration 008: Chat & Feedback
-- ============================================

BEGIN;

DROP TABLE IF EXISTS trade_feedback CASCADE;
DROP TABLE IF EXISTS trade_messages CASCADE;

CREATE TABLE trade_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    is_system BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE trade_feedback (
    feedback_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL REFERENCES trades(trade_id) ON DELETE CASCADE,
    from_user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    to_user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    rating VARCHAR(20) NOT NULL,
    comment TEXT,
    is_public BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(trade_id, from_user_id)
);

-- Index for faster chat lookup
CREATE INDEX idx_trade_messages_trade_id ON trade_messages(trade_id);

COMMIT;
