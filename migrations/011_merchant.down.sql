-- ============================================
-- Migration 011: Merchant (DOWN)
-- ============================================

BEGIN;
DROP TABLE IF EXISTS merchant_analytics CASCADE;
DROP TABLE IF EXISTS merchant_applications CASCADE;
COMMIT;
