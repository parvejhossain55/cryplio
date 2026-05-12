-- ============================================
-- Migration 029: Add Encrypted Private Key to Wallets
-- Store encrypted private keys for custodial wallets
-- ============================================

BEGIN;

ALTER TABLE wallets ADD COLUMN encrypted_private_key TEXT;

COMMENT ON COLUMN wallets.encrypted_private_key IS 'AES-encrypted private key for the wallet';

COMMIT;
