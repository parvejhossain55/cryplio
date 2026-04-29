-- ============================================
-- Migration 000: ENUM Type Definitions (DOWN)
-- Drops all custom enumerations
-- ============================================

BEGIN;

-- Drop types in any order (no dependencies)
DROP TYPE IF EXISTS merchant_status CASCADE;
DROP TYPE IF EXISTS admin_action_type CASCADE;
DROP TYPE IF EXISTS notification_type CASCADE;
DROP TYPE IF EXISTS feedback_rating CASCADE;
DROP TYPE IF EXISTS transaction_status CASCADE;
DROP TYPE IF EXISTS transaction_type CASCADE;
DROP TYPE IF EXISTS dispute_resolution CASCADE;
DROP TYPE IF EXISTS dispute_status CASCADE;
DROP TYPE IF EXISTS trade_status CASCADE;
DROP TYPE IF EXISTS price_type CASCADE;
DROP TYPE IF EXISTS ad_type CASCADE;
DROP TYPE IF EXISTS payment_category CASCADE;
DROP TYPE IF EXISTS kyc_level CASCADE;
DROP TYPE IF EXISTS user_status CASCADE;

COMMIT;
