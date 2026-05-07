-- ============================================
-- Migration 024: Withdrawal Approval Fields (Down)
-- ============================================

BEGIN;

ALTER TABLE wallet_transactions
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS requires_approval,
    DROP COLUMN IF EXISTS destination_address;

COMMIT;
