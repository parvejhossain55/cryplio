-- ============================================
-- Migration 018: Two-Factor Authentication Pending (down)
-- ============================================

BEGIN;

DROP TABLE IF EXISTS two_factor_pending CASCADE;

COMMIT;
