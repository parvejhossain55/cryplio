-- ============================================
-- Migration 017: OAuth Accounts (DOWN)
-- ============================================

BEGIN;

DROP TRIGGER IF EXISTS update_user_oauth_updated_at ON user_oauth;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS user_oauth;

COMMIT;
