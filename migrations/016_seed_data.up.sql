-- ============================================
-- Migration 016: Seed Data
-- Initial data for lookup tables and platform config
-- ============================================

BEGIN;

-- ============================================
-- CRYPTO ASSETS (SRS Table 5.2)
-- ============================================
INSERT INTO crypto_assets (symbol, name, blockchain, contract_address, decimals, min_confirmation) VALUES
('USDT', 'Tether', 'ERC20', '0xdAC17F958D2ee523a2206206994597C13D831ec7', 6, 12),
('USDT', 'Tether', 'TRC20', 'TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t', 6, 1),
('USDC', 'USD Coin', 'ERC20', '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48', 6, 12),
('ETH', 'Ethereum', 'ETH', NULL, 18, 1),
('BTC', 'Bitcoin', 'BTC', NULL, 8, 1),
('BNB', 'BNB', 'BSC', NULL, 18, 1)
ON CONFLICT (symbol, blockchain) DO NOTHING;

-- ============================================
-- FIAT CURRENCIES (SRS Table 5.3)
-- ============================================
INSERT INTO fiat_currencies (code, name, symbol) VALUES
('USD', 'US Dollar', '$'),
('BDT', 'Bangladeshi Taka', '৳'),
('EUR', 'Euro', '€'),
('GBP', 'British Pound', '£'),
('INR', 'Indian Rupee', '₹'),
('NGN', 'Nigerian Naira', '₦'),
('PHP', 'Philippine Peso', '₱'),
('KES', 'Kenyan Shilling', 'KSh'),
('CAD', 'Canadian Dollar', 'C$'),
('AUD', 'Australian Dollar', 'A$')
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- PAYMENT METHODS (SRS Table 5.3)
-- ============================================
INSERT INTO payment_methods (code, name, category, description, sort_order) VALUES
('bkash', 'bKash', 'mobile_money', 'Bangladesh mobile money provider', 1),
('nagad', 'Nagad', 'mobile_money', 'Bangladesh mobile money provider', 2),
('bank_transfer', 'Bank Transfer', 'bank_transfer', 'Direct bank wire transfer', 3),
('wise', 'Wise', 'online_wallet', 'International money transfer service', 4),
('paypal', 'PayPal', 'online_wallet', 'PayPal balance transfer', 5),
('sepa', 'SEPA Transfer', 'bank_transfer', 'EU Single Euro Payments Area transfer', 6),
('upi', 'UPI', 'mobile_money', 'India Unified Payments Interface', 7),
('gcash', 'GCash', 'mobile_money', 'Philippines mobile wallet', 8),
('mpesa', 'M-Pesa', 'mobile_money', 'East Africa mobile money', 9),
('cash_deposit', 'Cash Deposit', 'cash', 'Physical cash deposit at partner locations', 10)
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- DISPUTE REASONS (SRS FR-401 to FR-403)
-- ============================================
INSERT INTO dispute_reasons (code, label, description, category) VALUES
('payment_not_received', 'Payment not received', 'Buyer did not send payment or did not mark as paid', 'buyer'),
('wrong_amount', 'Wrong amount sent', 'Buyer sent incorrect fiat amount', 'buyer'),
('payment_did_not_clear', 'Payment did not clear', 'Bank transfer did not clear/failed', 'both'),
('fake_proof', 'Fake payment proof', 'Buyer provided fabricated/edited payment proof', 'seller'),
('account_problem', 'Account problem', 'Bank/mobile account issue prevented payment', 'both'),
('escrow_not_released', 'Escrow not released', 'Seller not releasing crypto after payment confirmation', 'seller'),
('fraud', 'Suspected fraud', 'Possible scam or fraudulent activity detected', 'both'),
('other', 'Other', 'Other reason not listed above', 'both')
ON CONFLICT (code) DO NOTHING;

-- ============================================
-- FEE TIERS (SRS 3.7 & 3.8 - Referral Commission Structure)
-- ============================================
INSERT INTO fee_tiers (name, min_volume_usd, max_volume_usd, fee_percentage, fee_minimum) VALUES
('Starter', 0, 10000, 0.010, 0),       -- 1% platform fee
('Volume', 10000, 100000, 0.008, 0),   -- 0.8% platform fee
('High Volume', 100000, NULL, 0.005, 0) -- 0.5% platform fee
ON CONFLICT DO NOTHING;

-- ============================================
-- PLATFORM CONFIG (SRS FR-804 - Admin fee configuration)
-- ============================================
INSERT INTO platform_config (key, value, value_type, description, is_public) VALUES
('platform_fee_percentage', '0.01', 'float', 'Default platform fee percentage (1%)', true),
('min_trade_amount', '10', 'float', 'Minimum trade amount in USD equivalent', true),
('max_trade_amount', '50000', 'float', 'Maximum trade amount per trade in USD', true),
('require_kyc_level_for_trade', 'level_1', 'string', 'Minimum KYC level required to trade', true),
('email_verification_required', 'true', 'bool', 'Require email verification before trading', true),
('max_failed_login_attempts', '5', 'int', 'Maximum failed login attempts before lockout', false),
('auto_dispute_time_minutes', '60', 'int', 'Auto-dispute if seller does not release after payment (minutes)', false),
('maintenance_mode', 'false', 'bool', 'Site-wide maintenance mode', true),
('allowed_countries', '["BD", "NG", "KE", "IN", "PH", "PK", "ID", "NG", "GH"]', 'json', 'List of allowed country codes', true),
('referral_commission_percentage', '0.20', 'float', 'Referral commission percentage of platform fee (20%)', true),
('merchant_monthly_fee', '50', 'float', 'Monthly merchant subscription fee in USD', false),
('max_concurrent_trades_per_user', '5', 'int', 'Maximum active trades per user', false),
('withdrawal_2fa_required', 'true', 'bool', 'Require 2FA for withdrawals', false)
ON CONFLICT (key) DO NOTHING;

COMMIT;
