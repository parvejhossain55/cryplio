-- ============================================
-- Migration 007: Seed Data
-- ============================================

BEGIN;

-- Cryptocurrencies
INSERT INTO crypto_assets (symbol, name, blockchain, decimals) VALUES
('USDT', 'Tether', 'ERC20', 6),
('USDT', 'Tether', 'TRC20', 6),
('USDC', 'USD Coin', 'ERC20', 6),
('ETH', 'Ethereum', 'ETH', 18),
('BTC', 'Bitcoin', 'BTC', 8),
('BNB', 'BNB', 'BSC', 18);

-- Fiat Currencies
INSERT INTO fiat_currencies (code, name, symbol) VALUES
('USD', 'US Dollar', '$'),
('BDT', 'Bangladeshi Taka', '৳'),
('EUR', 'Euro', '€'),
('INR', 'Indian Rupee', '₹'),
('NGN', 'Nigerian Naira', '₦'),
('KES', 'Kenyan Shilling', 'KSh'),
('PHP', 'Philippine Peso', '₱'),
('GBP', 'British Pound', '£');

-- Payment Methods
INSERT INTO payment_methods (code, name, category, sort_order) VALUES
('bkash', 'bKash', 'mobile_money', 1),
('nagad', 'Nagad', 'mobile_money', 2),
('bank_transfer', 'Bank Transfer', 'bank_transfer', 3),
('wise', 'Wise', 'online_wallet', 4),
('paypal', 'PayPal', 'online_wallet', 5);

-- Dispute Reasons
INSERT INTO dispute_reasons (code, label, category) VALUES
('payment_not_received', 'Payment not received', 'buyer'),
('wrong_amount', 'Wrong amount sent', 'buyer'),
('fake_proof', 'Fake payment proof', 'seller'),
('escrow_not_released', 'Escrow not released', 'seller'),
('other', 'Other', 'both');

-- Fee Tiers
INSERT INTO fee_tiers (name, min_volume_usd, max_volume_usd, fee_percentage) VALUES
('Starter', 0, 10000, 0.010),
('Volume', 10000, 100000, 0.008),
('High Volume', 100000, NULL, 0.005);

COMMIT;
