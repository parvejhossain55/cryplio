-- ============================================
-- Migration 006: Trade Feedback (DOWN)
-- ============================================

BEGIN;
DROP TABLE IF EXISTS trade_feedback CASCADE;
COMMIT;
