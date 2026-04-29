-- ============================================
-- Migration 012: Audit & Admin (DOWN)
-- Drops audit_logs, admin_actions, platform_config, announcements
-- ============================================

BEGIN;
DROP TABLE IF EXISTS announcements CASCADE;
DROP TABLE IF EXISTS platform_config CASCADE;
DROP TABLE IF EXISTS admin_actions CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;
COMMIT;
