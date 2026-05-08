-- Add the missing 'trade_completed' value to the notification_type enum.
-- ALTER TYPE ... ADD VALUE is transactional in PostgreSQL 12+.
ALTER TYPE notification_type ADD VALUE IF NOT EXISTS 'trade_completed';
