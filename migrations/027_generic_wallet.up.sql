-- ============================================
-- Migration 027: Generic Wallet (one per user for all tokens)
-- Change wallet to support all tokens, not just one crypto
-- ============================================

BEGIN;

-- Make crypto_id nullable (wallet no longer tied to specific crypto)
ALTER TABLE wallets ALTER COLUMN crypto_id DROP NOT NULL;

-- Drop the old unique constraint (one wallet per user per crypto)
ALTER TABLE wallets DROP CONSTRAINT IF EXISTS wallets_user_id_crypto_id_key;

-- Add new unique constraint: only one wallet per user
ALTER TABLE wallets ADD CONSTRAINT wallets_user_id_key UNIQUE (user_id);

-- Update comment to reflect new design
COMMENT ON TABLE wallets IS 'User wallets - one per user, holds all token balances';

COMMIT;
