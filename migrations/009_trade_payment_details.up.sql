-- Migration 009: Add Payment Details to Trades
BEGIN;

ALTER TABLE trades ADD COLUMN payment_details JSONB;

COMMIT;
