-- ============================================
-- Migration 003: Platform Settings (Down)
-- ============================================

BEGIN;

DROP TABLE IF EXISTS fee_tiers;
DROP TABLE IF EXISTS dispute_reasons;
DROP TABLE IF EXISTS user_payment_methods;
DROP TABLE IF EXISTS payment_methods;
DROP TABLE IF EXISTS fiat_currencies;
DROP TABLE IF EXISTS crypto_assets;

COMMIT;
