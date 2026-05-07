-- ============================================
-- Migration 024: Withdrawal Approval Fields
-- Add admin approval columns to wallet_transactions
-- ============================================

BEGIN;

ALTER TABLE wallet_transactions
    ADD COLUMN IF NOT EXISTS approved_by UUID REFERENCES users(user_id),
    ADD COLUMN IF NOT EXISTS approved_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS requires_approval BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS destination_address VARCHAR(255);

COMMIT;
