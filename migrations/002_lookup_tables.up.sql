-- ============================================
-- Migration 002: Lookup Tables
-- Reference data: cryptocurrencies, fiat currencies, payment methods, dispute reasons, fee tiers
-- ============================================

BEGIN;

-- Supported cryptocurrencies (per SRS 5.2)
CREATE TABLE IF NOT EXISTS crypto_assets (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(10) NOT NULL,
    name VARCHAR(50) NOT NULL,
    blockchain VARCHAR(20) NOT NULL,
    contract_address VARCHAR(42),
    decimals INT NOT NULL DEFAULT 18,
    min_confirmation INT NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(symbol, blockchain)
);

CREATE INDEX idx_crypto_assets_symbol ON crypto_assets(symbol);
CREATE INDEX idx_crypto_assets_active ON crypto_assets(is_active);

-- Supported fiat currencies
CREATE TABLE IF NOT EXISTS fiat_currencies (
    id SERIAL PRIMARY KEY,
    code CHAR(3) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    symbol VARCHAR(5) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fiat_currencies_code ON fiat_currencies(code);

-- Payment methods (per SRS 5.3)
CREATE TABLE IF NOT EXISTS payment_methods (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(30) NOT NULL UNIQUE,
    name VARCHAR(50) NOT NULL,
    category payment_category NOT NULL,
    icon_url VARCHAR(255),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_payment_methods_code ON payment_methods(code);
CREATE INDEX idx_payment_methods_category ON payment_methods(category);

-- Dispute reasons (per SRS 3.4)
CREATE TABLE IF NOT EXISTS dispute_reasons (
    id SMALLSERIAL PRIMARY KEY,
    code VARCHAR(30) NOT NULL UNIQUE,
    label VARCHAR(100) NOT NULL,
    description TEXT,
    category VARCHAR(30) NOT NULL
);

CREATE INDEX idx_dispute_reasons_code ON dispute_reasons(code);

-- Platform fee tiers (per SRS 3.8 & 3.7)
CREATE TABLE IF NOT EXISTS fee_tiers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    min_volume_usd DECIMAL(15,2) NOT NULL,
    max_volume_usd DECIMAL(15,2),
    fee_percentage DECIMAL(5,3) NOT NULL,
    fee_minimum DECIMAL(10,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_fee_tiers_volume_range ON fee_tiers(min_volume_usd, max_volume_usd);

-- ============================================
-- SEED DATA (from SRS 4.2, 5.2, 5.3)
-- ============================================

-- Cryptocurrencies (SRS 5.2 Supported Cryptocurrencies)
INSERT INTO crypto_assets (symbol, name, blockchain, decimals, min_confirmation) VALUES
('USDT', 'Tether', 'ERC20', 6, 12),
('USDT', 'Tether', 'TRC20', 6, 1),
('USDC', 'USD Coin', 'ERC20', 6, 12),
('ETH', 'Ethereum', 'ETH', 18, 1),
('BTC', 'Bitcoin', 'BTC', 8, 1),
('BNB', 'BNB', 'BSC', 18, 1)
ON CONFLICT DO NOTHING;

-- Fiat currencies (global support)
INSERT INTO fiat_currencies (code, name, symbol) VALUES
('USD', 'US Dollar', '$'),
('BDT', 'Bangladeshi Taka', '৳'),
('EUR', 'Euro', '€'),
('INR', 'Indian Rupee', '₹'),
('NGN', 'Nigerian Naira', '₦'),
('KES', 'Kenyan Shilling', 'KSh'),
('PHP', 'Philippine Peso', '₱'),
('GBP', 'British Pound', '£')
ON CONFLICT (code) DO NOTHING;

-- Payment methods (SRS 5.3 Supported Payment Methods)
INSERT INTO payment_methods (code, name, category, description, sort_order) VALUES
('bkash', 'bKash', 'mobile_money', 'Bangladesh mobile money', 1),
('nagad', 'Nagad', 'mobile_money', 'Bangladesh mobile money', 2),
('bank_transfer', 'Bank Transfer', 'bank_transfer', 'Direct bank transfer', 3),
('wise', 'Wise', 'online_wallet', 'International money transfer', 4),
('paypal', 'PayPal', 'online_wallet', 'PayPal balance', 5),
('sepa', 'SEPA Transfer', 'bank_transfer', 'EU bank transfer (SEPA)', 6),
('upi', 'UPI', 'mobile_money', 'India Unified Payments Interface', 7),
('gcash', 'GCash', 'mobile_money', 'Philippines mobile wallet', 8),
('mpesa', 'M-Pesa', 'mobile_money', 'East Africa mobile money', 9)
ON CONFLICT (code) DO NOTHING;

-- Dispute reasons (SRS FR-401)
INSERT INTO dispute_reasons (code, label, description, category) VALUES
('payment_not_received', 'Payment not received', 'Buyer did not send payment', 'buyer'),
('wrong_amount', 'Wrong amount sent', 'Buyer sent incorrect fiat amount', 'buyer'),
('payment_did_not_clear', 'Payment did not clear', 'Payment failed or did not clear bank', 'both'),
('fake_proof', 'Fake payment proof', 'Buyer provided fabricated payment proof', 'seller'),
('account_problem', 'Account problem', 'Bank/mobile account issue prevented payment', 'both'),
('escrow_not_released', 'Escrow not released', 'Seller not releasing after payment confirmation', 'seller'),
('fraud', 'Suspected fraud', 'Possible scam or fraudulent activity', 'both'),
('other', 'Other', 'Other reason not listed', 'both')
ON CONFLICT (code) DO NOTHING;

-- Fee tiers (SRS 3.7 & 3.8 - Referral & Merchant)
INSERT INTO fee_tiers (name, min_volume_usd, max_volume_usd, fee_percentage, fee_minimum) VALUES
('Starter', 0, 10000, 0.010, 0),      -- 1%
('Volume', 10000, 100000, 0.008, 0),  -- 0.8%
('High Volume', 100000, NULL, 0.005, 0)  -- 0.5%
ON CONFLICT DO NOTHING;

COMMIT;
