-- ============================================
-- Migration 027 Down: Revert to crypto-specific wallets
-- ============================================

BEGIN;

-- Drop the single-wallet-per-user constraint
ALTER TABLE wallets DROP CONSTRAINT IF EXISTS wallets_user_id_key;

-- Re-add the original unique constraint
ALTER TABLE wallets ADD CONSTRAINT wallets_user_id_crypto_id_key UNIQUE (user_id, crypto_id);

-- Make crypto_id required again
ALTER TABLE wallets ALTER COLUMN crypto_id SET NOT NULL;

COMMENT ON TABLE wallets IS 'User wallets - one per crypto asset';

COMMIT;
