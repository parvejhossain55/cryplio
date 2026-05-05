-- ============================================
-- Migration 015: Triggers (DOWN)
-- Drops all update_updated_at_column triggers
-- ============================================

BEGIN;
-- Drop triggers on all tables
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_user_stats_updated_at ON user_stats;
DROP TRIGGER IF EXISTS update_trade_ads_updated_at ON trade_ads;
DROP TRIGGER IF EXISTS update_trades_updated_at ON trades;
DROP TRIGGER IF EXISTS update_wallets_updated_at ON wallets;
DROP TRIGGER IF EXISTS update_wallet_transactions_updated_at ON wallet_transactions;
DROP TRIGGER IF EXISTS update_disputes_updated_at ON disputes;
DROP TRIGGER IF EXISTS update_notifications_updated_at ON notifications;
DROP TRIGGER IF EXISTS update_referrals_updated_at ON referrals;
DROP TRIGGER IF EXISTS update_merchant_applications_updated_at ON merchant_applications;
DROP TRIGGER IF EXISTS update_announcements_updated_at ON announcements;
DROP TRIGGER IF EXISTS update_platform_config_updated_at ON platform_config;
DROP TRIGGER IF EXISTS update_admin_actions_updated_at ON admin_actions;
DROP TRIGGER IF EXISTS update_blockchain_links_updated_at ON blockchain_links;
DROP TRIGGER IF EXISTS update_merchant_analytics_updated_at ON merchant_analytics;
DROP TRIGGER IF EXISTS update_user_blocks_updated_at ON user_blocks;

-- Drop helper function
DROP FUNCTION IF EXISTS update_updated_at_column();
COMMIT;
