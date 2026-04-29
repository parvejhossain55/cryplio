-- ============================================
-- Migration 009: Notifications (DOWN)
-- ============================================

BEGIN;
-- Drop child table first
DROP TABLE IF EXISTS notification_preferences CASCADE;
-- Drop main notifications table
DROP TABLE IF EXISTS notifications CASCADE;
COMMIT;
