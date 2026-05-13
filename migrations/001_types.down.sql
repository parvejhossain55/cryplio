-- ============================================
-- Migration 001: Types & Utility Functions (Down)
-- ============================================

BEGIN;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TYPE IF EXISTS admin_action_type;
DROP TYPE IF EXISTS payment_category;
DROP TYPE IF EXISTS feedback_rating;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS dispute_resolution;
DROP TYPE IF EXISTS dispute_status;
DROP TYPE IF EXISTS trade_status;
DROP TYPE IF EXISTS price_type;
DROP TYPE IF EXISTS ad_type;
DROP TYPE IF EXISTS user_role;
DROP TYPE IF EXISTS user_status;

COMMIT;
