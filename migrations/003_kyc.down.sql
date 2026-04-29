-- ============================================
-- Migration 003: KYC Tables (DOWN)
-- ============================================

BEGIN;
DROP TABLE IF EXISTS kyc_records CASCADE;
COMMIT;
