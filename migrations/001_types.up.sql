-- ============================================
-- Migration 001: Types & Utility Functions
-- ============================================

BEGIN;

CREATE TYPE user_status AS ENUM ('pending', 'active', 'suspended', 'banned', 'deleted');
CREATE TYPE user_role AS ENUM ('user', 'admin');
CREATE TYPE ad_type AS ENUM ('buy', 'sell');
CREATE TYPE price_type AS ENUM ('fixed', 'floating');
CREATE TYPE trade_status AS ENUM ('pending', 'active', 'paid', 'released', 'cancelled', 'disputed', 'completed', 'expired');
CREATE TYPE dispute_status AS ENUM ('pending', 'assigned', 'under_review', 'resolved', 'appealed', 'closed');
CREATE TYPE dispute_resolution AS ENUM ('release_to_buyer', 'return_to_seller', 'partial_split', 'cancel');
CREATE TYPE transaction_type AS ENUM ('deposit', 'withdrawal', 'trade_sale', 'trade_purchase', 'fee', 'refund', 'escrow_lock', 'escrow_release', 'dispute_hold', 'dispute_refund');
CREATE TYPE transaction_status AS ENUM ('pending', 'confirmed', 'completed', 'failed', 'cancelled');
CREATE TYPE feedback_rating AS ENUM ('positive', 'neutral', 'negative');
CREATE TYPE payment_category AS ENUM ('mobile_money', 'bank_transfer', 'online_wallet', 'crypto', 'cash');
CREATE TYPE admin_action_type AS ENUM ('user_suspend', 'user_ban', 'user_unban', 'dispute_resolve', 'withdrawal_approve', 'withdrawal_reject', 'announcement_post', 'fee_update', 'config_change', 'bulk_message', 'report_generate');

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

COMMIT;
