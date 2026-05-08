-- Down migration for trade_status_log and email service tables

-- Drop triggers first
DROP TRIGGER IF EXISTS update_email_queue_updated_at ON email_queue;
DROP TRIGGER IF EXISTS update_email_templates_updated_at ON email_templates;

-- Drop indexes
DROP INDEX IF EXISTS idx_email_queue_next_attempt;
DROP INDEX IF EXISTS idx_email_queue_created_at;
DROP INDEX IF EXISTS idx_email_queue_priority;
DROP INDEX IF EXISTS idx_email_queue_status;
DROP INDEX IF EXISTS idx_trade_status_log_created_at;
DROP INDEX IF EXISTS idx_trade_status_log_trade_id;

-- Drop tables (reverse order of creation to handle FK constraints)
DROP TABLE IF EXISTS email_queue;
DROP TABLE IF EXISTS email_templates;
DROP TABLE IF EXISTS trade_status_log;
