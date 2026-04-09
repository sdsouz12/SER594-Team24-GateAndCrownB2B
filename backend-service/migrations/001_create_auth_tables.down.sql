-- Drop indexes
DROP INDEX IF EXISTS idx_sys_user_role_rel_role_id;
DROP INDEX IF EXISTS idx_sys_user_role_rel_user_id;
DROP INDEX IF EXISTS idx_sys_user_organization_id;
DROP INDEX IF EXISTS idx_sys_user_provider_id;
DROP INDEX IF EXISTS idx_sys_user_email;
DROP INDEX IF EXISTS idx_sys_user_username;

-- Drop tables
DROP TABLE IF EXISTS sys_user_role_rel;
DROP TABLE IF EXISTS sys_user;
DROP TABLE IF EXISTS sys_role;
