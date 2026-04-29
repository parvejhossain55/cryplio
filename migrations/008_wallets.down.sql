-- ============================================
-- Migration 008: Wallets & Transactions (DOWN)
-- ============================================

BEGIN;
-- Drop child table first
DROP TABLE IF EXISTS wallet_transactions CASCADE;
-- Drop main wallets table
DROP TABLE IF EXISTS wallets CASCADE;
COMMIT;
