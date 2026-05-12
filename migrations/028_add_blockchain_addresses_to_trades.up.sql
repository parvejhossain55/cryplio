-- ============================================
-- Migration 028: Add Blockchain Addresses to Trades
-- Store buyer, seller, and token addresses for escrow interaction
-- ============================================

BEGIN;

ALTER TABLE trades ADD COLUMN buyer_address VARCHAR(42);
ALTER TABLE trades ADD COLUMN seller_address VARCHAR(42);
ALTER TABLE trades ADD COLUMN token_address VARCHAR(42);

COMMENT ON COLUMN trades.buyer_address IS 'On-chain address of the buyer';
COMMENT ON COLUMN trades.seller_address IS 'On-chain address of the seller';
COMMENT ON COLUMN trades.token_address IS 'On-chain address of the token being traded';

COMMIT;
