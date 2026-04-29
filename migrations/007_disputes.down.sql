-- ============================================
-- Migration 007: Disputes (DOWN)
-- ============================================

BEGIN;
-- Drop child table first
DROP TABLE IF EXISTS dispute_messages CASCADE;
-- Drop main disputes table
DROP TABLE IF EXISTS disputes CASCADE;
COMMIT;
