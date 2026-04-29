-- ============================================
-- Migration 006: Trade Feedback
-- Post-trade ratings and reviews
-- ============================================

BEGIN;

CREATE TABLE IF NOT EXISTS trade_feedback (
    feedback_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL UNIQUE REFERENCES trades(trade_id),
    from_user_id UUID NOT NULL REFERENCES users(user_id),
    to_user_id UUID NOT NULL REFERENCES users(user_id),
    rating feedback_rating NOT NULL,
    comment TEXT,
    is_public BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_feedback_to_user ON trade_feedback(to_user_id);
CREATE INDEX idx_feedback_rating ON trade_feedback(rating);
CREATE INDEX idx_feedback_created ON trade_feedback(created_at DESC);

COMMIT;
