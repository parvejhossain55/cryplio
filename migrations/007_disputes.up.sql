-- ============================================
-- Migration 007: Disputes
-- Trade dispute resolution system (SRS 3.4)
-- ============================================

BEGIN;

CREATE TABLE IF NOT EXISTS disputes (
    dispute_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    trade_id UUID NOT NULL UNIQUE REFERENCES trades(trade_id),
    raised_by UUID NOT NULL REFERENCES users(user_id),
    reason_code VARCHAR(30) NOT NULL,
    reason_text TEXT,
    evidence_links JSONB DEFAULT '[]',
    status dispute_status NOT NULL DEFAULT 'pending',
    assigned_admin UUID REFERENCES users(user_id),
    assigned_at TIMESTAMP,
    escalated_at TIMESTAMP,
    resolution_type dispute_resolution,
    resolution_note TEXT,
    released_amount DECIMAL(20,8),
    returned_amount DECIMAL(20,8),
    resolved_at TIMESTAMP,
    closed_at TIMESTAMP,
    appeal_status BOOLEAN DEFAULT false,
    appeal_note TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_disputes_trade ON disputes(trade_id);
CREATE INDEX idx_disputes_raised_by ON disputes(raised_by);
CREATE INDEX idx_disputes_admin ON disputes(assigned_admin);
CREATE INDEX idx_disputes_status ON disputes(status);
CREATE INDEX idx_disputes_created ON disputes(created_at DESC);

-- Dispute chat messages (3-way)
CREATE TABLE IF NOT EXISTS dispute_messages (
    message_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dispute_id UUID NOT NULL REFERENCES disputes(dispute_id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(user_id),
    message_type VARCHAR(20) NOT NULL DEFAULT 'text',
    content TEXT NOT NULL,
    file_url VARCHAR(500),
    is_internal BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dispute_messages_dispute ON dispute_messages(dispute_id, created_at DESC);

COMMIT;
