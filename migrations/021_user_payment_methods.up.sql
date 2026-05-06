BEGIN;

CREATE TABLE IF NOT EXISTS user_payment_methods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    payment_method_code VARCHAR(30) NOT NULL REFERENCES payment_methods(code),
    display_name VARCHAR(100) NOT NULL,
    account_name VARCHAR(100),
    account_number VARCHAR(100),
    bank_name VARCHAR(100),
    metadata JSONB,
    is_active BOOLEAN NOT NULL DEFAULT true,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_payment_methods_user ON user_payment_methods(user_id);
CREATE INDEX idx_user_payment_methods_active ON user_payment_methods(is_active);

COMMIT;
