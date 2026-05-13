-- ============================================
-- Migration 004: Wallets & Trading (Down)
-- ============================================

BEGIN;

DROP TABLE IF EXISTS trade_feedback;
DROP TABLE IF EXISTS trade_messages;
DROP TABLE IF EXISTS disputes;
DROP TABLE IF EXISTS wallet_transactions;
DROP TABLE IF EXISTS trades;
DROP TABLE IF EXISTS trade_ads;
DROP TABLE IF EXISTS wallets;

COMMIT;
