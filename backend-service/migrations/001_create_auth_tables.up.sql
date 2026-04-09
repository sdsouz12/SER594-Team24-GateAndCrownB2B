-- Create sys_role table
CREATE TABLE IF NOT EXISTS sys_role (
    role_id VARCHAR(50) PRIMARY KEY,
    slug VARCHAR(100) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create sys_user table
CREATE TABLE IF NOT EXISTS sys_user (
    user_id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(200),
    primary_role_id VARCHAR(50) REFERENCES sys_role(role_id),
    organization_id BIGINT,
    provider_id BIGINT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create sys_user_role_rel table
CREATE TABLE IF NOT EXISTS sys_user_role_rel (
    rel_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES sys_user(user_id) ON DELETE CASCADE,
    role_id VARCHAR(50) NOT NULL REFERENCES sys_role(role_id) ON DELETE CASCADE,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_sys_user_username ON sys_user(username);
CREATE INDEX IF NOT EXISTS idx_sys_user_email ON sys_user(email);
CREATE INDEX IF NOT EXISTS idx_sys_user_provider_id ON sys_user(provider_id);
CREATE INDEX IF NOT EXISTS idx_sys_user_organization_id ON sys_user(organization_id);
CREATE INDEX IF NOT EXISTS idx_sys_user_role_rel_user_id ON sys_user_role_rel(user_id);
CREATE INDEX IF NOT EXISTS idx_sys_user_role_rel_role_id ON sys_user_role_rel(role_id);

-- Insert default roles
INSERT INTO sys_role (role_id, slug, name, description) VALUES
    ('CLIENT', 'CLIENT', 'Client', 'Client user who creates gate orders'),
    ('ADMIN', 'ADMIN', 'Admin', 'Admin for specific company/provider'),
    ('SUPERADMIN', 'SUPERADMIN', 'Super Admin', 'Full system access, sees all gate orders')
ON CONFLICT (role_id) DO NOTHING;

-- Insert a test super admin user (password: "admin123")
-- Password hash for "admin123" using bcrypt
INSERT INTO sys_user (username, email, password_hash, full_name, primary_role_id, status) VALUES
    ('admin', 'admin@gate-backoffice.com', '$2a$10$gd4ESQwt2XF0ILSIeY1v/eNlwVRVbuRUaTIYpZX/3Nxs0v6XgR8lK', 'System Administrator', 'SUPERADMIN', 'active')
ON CONFLICT (username) DO NOTHING;

-- Assign ADMIN and SUPERADMIN roles to test user (user_id = 1)
INSERT INTO sys_user_role_rel (user_id, role_id) VALUES
    (1, 'ADMIN'),
    (1, 'SUPERADMIN')
ON CONFLICT (user_id, role_id) DO NOTHING;
