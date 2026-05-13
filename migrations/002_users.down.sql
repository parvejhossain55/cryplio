-- ============================================
-- Migration 002: Identity & Access (Down)
-- ============================================

BEGIN;

DROP TABLE IF EXISTS two_factor_pending;
DROP TABLE IF EXISTS user_oauth;
DROP TABLE IF EXISTS email_verification_tokens;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS user_sessions;
DROP TABLE IF EXISTS user_stats;
DROP TABLE IF EXISTS users;

COMMIT;
