-- Create bay_organization table
CREATE TABLE IF NOT EXISTS bay_organization (
    organization_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    address TEXT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_bay_organization_name ON bay_organization(name);
CREATE INDEX IF NOT EXISTS idx_bay_organization_status ON bay_organization(status);

-- Insert test organizations
INSERT INTO bay_organization (name, description, contact_email, status) VALUES
    ('Acme Corporation', 'Main testing organization', 'contact@acme.com', 'active'),
    ('Beta Solutions', 'Secondary organization for testing', 'info@beta.com', 'active')
ON CONFLICT DO NOTHING;
