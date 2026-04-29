-- ============================================
-- Migration 001: Core User & Auth Tables (DOWN)
-- Drops all user-related tables
-- ============================================

BEGIN;

-- Drop child tables first (respect FK dependencies)
DROP TABLE IF EXISTS user_blocks CASCADE;
DROP TABLE IF EXISTS email_verification_tokens CASCADE;
DROP TABLE IF EXISTS password_reset_tokens CASCADE;
DROP TABLE IF EXISTS user_sessions CASCADE;
DROP TABLE IF EXISTS user_stats CASCADE;

-- Drop main users table (parent)
DROP TABLE IF EXISTS users CASCADE;

COMMIT;
