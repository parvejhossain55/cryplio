-- ============================================
-- Migration 002: Lookup Tables (DOWN)
-- Drops lookup/reference tables
-- ============================================

BEGIN;

-- Drop in reverse dependency order
DROP TABLE IF EXISTS fee_tiers CASCADE;
DROP TABLE IF EXISTS dispute_reasons CASCADE;
DROP TABLE IF EXISTS payment_methods CASCADE;
DROP TABLE IF EXISTS fiat_currencies CASCADE;
DROP TABLE IF EXISTS crypto_assets CASCADE;

COMMIT;
