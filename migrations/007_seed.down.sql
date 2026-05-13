-- ============================================
-- Migration 007: Seed Data (Down)
-- ============================================

BEGIN;

DELETE FROM fee_tiers;
DELETE FROM dispute_reasons;
DELETE FROM payment_methods;
DELETE FROM fiat_currencies;
DELETE FROM crypto_assets;

COMMIT;
