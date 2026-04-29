-- ============================================
-- Migration 016: Seed Data (DOWN)
-- NOTE: Seed data down migration is intentionally a no-op.
-- Lookup tables are referenced by other tables; deleting seed rows
-- would violate foreign key constraints. Full database reset requires
-- dropping all tables (migration 002 down or full drop).
-- ============================================

BEGIN;
-- No operation: seed data cannot be safely rolled back individually.
-- To reset all data, rollback all migrations or drop database.
COMMIT;
