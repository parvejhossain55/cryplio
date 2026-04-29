-- ============================================
-- Migration 010: Referrals (DOWN)
-- ============================================

BEGIN;
DROP TABLE IF EXISTS referrals CASCADE;
COMMIT;
