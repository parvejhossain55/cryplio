-- ============================================
-- Migration 000: ENUM Type Definitions
-- Defines all custom enumerations used across tables
-- ============================================

BEGIN;

-- User-related enums
CREATE TYPE user_status AS ENUM ('pending', 'active', 'suspended', 'banned', 'deleted');
CREATE TYPE kyc_level AS ENUM ('level_0', 'level_1', 'level_2', 'level_3');

-- Trade-related enums
CREATE TYPE ad_type AS ENUM ('buy', 'sell');
CREATE TYPE price_type AS ENUM ('fixed', 'floating');
CREATE TYPE trade_status AS ENUM ('pending', 'active', 'paid', 'released', 'cancelled', 'disputed', 'completed', 'expired');

-- Escrow & dispute enums
CREATE TYPE dispute_status AS ENUM ('pending', 'assigned', 'under_review', 'resolved', 'appealed', 'closed');
CREATE TYPE dispute_resolution AS ENUM ('release_to_buyer', 'return_to_seller', 'partial_split', 'cancel');

-- Transaction enums
CREATE TYPE transaction_type AS ENUM ('deposit', 'withdrawal', 'trade_sale', 'trade_purchase', 'fee', 'refund', 'escrow_lock', 'escrow_release', 'dispute_hold', 'dispute_refund');
CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'completed', 'failed', 'cancelled');

-- Other enums
CREATE TYPE feedback_rating AS ENUM ('positive', 'neutral', 'negative');
CREATE TYPE payment_category AS ENUM ('mobile_money', 'bank_transfer', 'online_wallet', 'crypto', 'cash');
CREATE TYPE notification_type AS ENUM ('trade_started', 'trade_paid', 'trade_released', 'trade_cancelled', 'trade_disputed', 'dispute_resolved', 'new_message', 'deposit_received', 'withdrawal_completed', 'kyc_approved', 'kyc_rejected', 'merchant_approved', 'referral_earned', 'system_announcement');
CREATE TYPE admin_action_type AS ENUM ('user_suspend', 'user_ban', 'user_unban', 'dispute_resolve', 'withdrawal_approve', 'withdrawal_reject', 'kyc_review', 'announcement_post', 'fee_update', 'merchant_approve', 'merchant_reject', 'config_change', 'bulk_message', 'report_generate');
CREATE TYPE merchant_status AS ENUM ('none', 'pending', 'approved', 'rejected', 'suspended');

COMMIT;
