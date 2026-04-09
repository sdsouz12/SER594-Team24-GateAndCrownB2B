-- Drop index
DROP INDEX IF EXISTS idx_sys_user_phone_number;

-- Drop column
ALTER TABLE sys_user DROP COLUMN IF EXISTS phone_number;
