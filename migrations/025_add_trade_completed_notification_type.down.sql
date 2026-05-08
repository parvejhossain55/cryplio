-- PostgreSQL does not support removing enum values.
-- To roll back, recreate the enum without 'trade_completed' (requires replacing all usages).
-- This migration is intentionally a no-op; manual rollback requires a table rebuild.
SELECT 1; -- no-op placeholder
