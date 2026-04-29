-- ============================================
-- Migration 015: Triggers for Timestamp Updates
-- Auto-update updated_at on row updates
-- ============================================

BEGIN;

-- Create helper function if not exists
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply triggers to all tables with updated_at column
DO $$
DECLARE
    table_record RECORD;
BEGIN
    FOR table_record IN
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public'
          AND tablename IN (
              'users', 'user_stats', 'kyc_records', 'trade_ads', 'trades',
              'wallets', 'wallet_transactions', 'disputes', 'notifications',
              'referrals', 'merchant_applications', 'announcements', 'platform_config',
              'admin_actions', 'blockchain_links', 'merchant_analytics', 'user_blocks'
          )
    LOOP
        EXECUTE format(
            'DROP TRIGGER IF EXISTS update_%s_updated_at ON %I;
             CREATE TRIGGER update_%s_updated_at
             BEFORE UPDATE ON %I
             FOR EACH ROW
             EXECUTE FUNCTION update_updated_at_column();',
            table_record.tablename, table_record.tablename,
            table_record.tablename, table_record.tablename
        );
    END LOOP;
END $$;

COMMIT;
