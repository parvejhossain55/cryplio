-- ============================================
-- Migration 013: Monitoring (DOWN)
-- Drops rate_limit_counts, login_attempts, api_request_logs
-- ============================================

BEGIN;
DROP TABLE IF EXISTS api_request_logs CASCADE;
DROP TABLE IF EXISTS login_attempts CASCADE;
DROP TABLE IF EXISTS rate_limit_counts CASCADE;
COMMIT;
