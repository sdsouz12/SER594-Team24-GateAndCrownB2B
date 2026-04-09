-- Remove test users
DELETE FROM sys_user_role_rel 
WHERE user_id IN (
    SELECT user_id FROM sys_user 
    WHERE username IN ('john.doe', 'jane.smith', 'bob.wilson', 'sarah.jones', 'mike.brown', 'lisa.davis')
);

DELETE FROM sys_user 
WHERE username IN ('john.doe', 'jane.smith', 'bob.wilson', 'sarah.jones', 'mike.brown', 'lisa.davis');
