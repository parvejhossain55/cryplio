-- ============================================
-- Migration 005: System & Support (Down)
-- ============================================

BEGIN;

DROP TABLE IF EXISTS email_queue;
DROP TABLE IF EXISTS email_templates;
DROP TABLE IF EXISTS admin_audit_log;
DROP TABLE IF EXISTS withdrawal_requests;
DROP TABLE IF EXISTS notifications;

COMMIT;
