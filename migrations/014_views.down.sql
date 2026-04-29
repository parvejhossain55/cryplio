-- ============================================
-- Migration 014: Views (DOWN)
-- Drops all database views
-- ============================================

BEGIN;
DROP VIEW IF EXISTS open_disputes;
DROP VIEW IF EXISTS public_user_profiles;
DROP VIEW IF EXISTS completed_trades;
DROP VIEW IF EXISTS active_trade_ads;
COMMIT;
