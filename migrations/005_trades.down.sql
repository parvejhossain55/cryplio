-- ============================================
-- Migration 005: Trades (DOWN)
-- Drops trade_attachments, trade_messages, and trades
-- ============================================

BEGIN;
-- Drop child tables first
DROP TABLE IF EXISTS trade_attachments CASCADE;
DROP TABLE IF EXISTS trade_messages CASCADE;
-- Drop main trades table
DROP TABLE IF EXISTS trades CASCADE;
COMMIT;
