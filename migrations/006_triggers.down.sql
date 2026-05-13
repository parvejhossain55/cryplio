-- ============================================
-- Migration 006: Triggers (Down)
-- ============================================

BEGIN;

DROP TRIGGER IF EXISTS update_email_templates_updated_at ON email_templates;
DROP TRIGGER IF EXISTS update_withdrawal_requests_updated_at ON withdrawal_requests;
DROP TRIGGER IF EXISTS update_disputes_updated_at ON disputes;
DROP TRIGGER IF EXISTS update_wallet_transactions_updated_at ON wallet_transactions;
DROP TRIGGER IF EXISTS update_trades_updated_at ON trades;
DROP TRIGGER IF EXISTS update_trade_ads_updated_at ON trade_ads;
DROP TRIGGER IF EXISTS update_wallets_updated_at ON wallets;
DROP TRIGGER IF EXISTS update_user_payment_methods_updated_at ON user_payment_methods;
DROP TRIGGER IF EXISTS update_crypto_assets_updated_at ON crypto_assets;
DROP TRIGGER IF EXISTS update_user_oauth_updated_at ON user_oauth;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

COMMIT;
