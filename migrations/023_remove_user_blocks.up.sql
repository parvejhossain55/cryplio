BEGIN;

DROP TRIGGER IF EXISTS update_user_blocks_updated_at ON user_blocks;
DROP TABLE IF EXISTS user_blocks CASCADE;

COMMIT;
