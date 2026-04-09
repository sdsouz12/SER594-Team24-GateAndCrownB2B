-- Add test users for organizations
-- Password for all test users: "admin123" (bcrypt hash: $2a$10$gd4ESQwt2XF0ILSIeY1v/eNlwVRVbuRUaTIYpZX/3Nxs0v6XgR8lK)

INSERT INTO sys_user (username, password_hash, full_name, email, phone_number, primary_role_id, organization_id, status) VALUES
    -- Acme Corporation users (organization_id = 1)
    ('john.doe', '$2a$10$gd4ESQwt2XF0ILSIeY1v/eNlwVRVbuRUaTIYpZX/3Nxs0v6XgR8lK', 'John Doe', 'john.doe@acme.com', '+1-555-1001', 'ADMIN', 1, 'active'),
    ('jane.smith', '$2a$10$gd4ESQwt2XF0ILSIeY1v/eNlwVRVbuRUaTIYpZX/3Nxs0v6XgR8lK', 'Jane Smith', 'jane.smith@acme.com', '+1-555-1002', 'ADMIN', 1, 'active'),
    
    -- Beta Solutions users (organization_id = 2)
    ('sarah.jones', '$2a$10$gd4ESQwt2XF0ILSIeY1v/eNlwVRVbuRUaTIYpZX/3Nxs0v6XgR8lK', 'Sarah Jones', 'sarah.jones@beta.com', '+1-555-2001', 'ADMIN', 2, 'active'),
    ('mike.brown', '$2a$10$gd4ESQwt2XF0ILSIeY1v/eNlwVRVbuRUaTIYpZX/3Nxs0v6XgR8lK', 'Mike Brown', 'mike.brown@beta.com', '+1-555-2002', 'ADMIN', 2, 'active')
ON CONFLICT (username) DO NOTHING;

-- Also add user_role relationships
INSERT INTO sys_user_role_rel (user_id, role_id) 
SELECT user_id, primary_role_id FROM sys_user WHERE username IN ('john.doe', 'jane.smith', 'sarah.jones', 'mike.brown')
ON CONFLICT DO NOTHING;
