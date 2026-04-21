-- Enable pgvector extension for semantic search
CREATE EXTENSION IF NOT EXISTS vector;

-- Catalog products table with embedding support
CREATE TABLE IF NOT EXISTS catalog_product (
    product_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    price_range VARCHAR(100),
    material VARCHAR(100),
    dimensions VARCHAR(200),
    status VARCHAR(20) DEFAULT 'active',
    embedding vector(384),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_catalog_product_category ON catalog_product(category);
CREATE INDEX IF NOT EXISTS idx_catalog_product_status ON catalog_product(status);
CREATE INDEX IF NOT EXISTS idx_catalog_product_embedding ON catalog_product
    USING ivfflat (embedding vector_cosine_ops) WITH (lists = 10);

-- Seed catalog data
INSERT INTO catalog_product (name, category, description, price_range, material, dimensions) VALUES
    ('Classic Frame Gate', 'frame', 'Traditional frame gate with clean lines. Ideal for residential driveways and commercial entrances.', '$800 - $1,500', 'Steel', '3m x 2m'),
    ('Heavy-Duty Industrial Gate', 'frame', 'Reinforced steel frame gate for industrial sites and warehouses. Includes motorized option.', '$2,000 - $4,500', 'Heavy Steel', '5m x 3m'),
    ('Sliding Frame Gate', 'frame', 'Space-saving sliding mechanism. Perfect for narrow driveways and tight spaces.', '$1,200 - $2,800', 'Galvanized Steel', '4m x 2m'),
    ('Crown Panel Standard', 'crown', 'Premium crown panel with decorative top. Adds elegance to any entrance gate.', '$300 - $600', 'Wrought Iron', '1m x 0.5m'),
    ('Crown Panel Ornate', 'crown', 'Ornate crown with floral patterns. Handcrafted decorative panel for luxury gates.', '$500 - $1,000', 'Cast Iron', '1.2m x 0.6m'),
    ('Crown Panel Modern', 'crown', 'Minimalist modern crown design. Clean geometric shapes for contemporary architecture.', '$250 - $500', 'Aluminum', '0.8m x 0.4m'),
    ('MDF Interior Door', 'mdf', 'Smooth MDF door panel for interior applications. Available in various finishes.', '$150 - $350', 'MDF', '2.1m x 0.9m'),
    ('MDF Exterior Panel', 'mdf', 'Weather-resistant MDF panel with protective coating. Suitable for outdoor use.', '$200 - $450', 'Waterproof MDF', '2m x 1m'),
    ('Tempered Glass Panel', 'glass', 'Safety tempered glass panel for gates. Provides visibility while maintaining security.', '$400 - $900', 'Tempered Glass', '0.5m x 1.5m'),
    ('Frosted Glass Insert', 'glass', 'Privacy frosted glass insert. Diffuses light while blocking direct view.', '$350 - $750', 'Frosted Glass', '0.4m x 1.2m'),
    ('Standard Hinge Set', 'fittings', 'Heavy-duty hinge set for gates up to 200kg. Includes all mounting hardware.', '$50 - $120', 'Stainless Steel', 'N/A'),
    ('Gate Lock System', 'fittings', 'Electronic lock system with remote control and keypad. Includes installation kit.', '$200 - $500', 'Steel/Electronics', 'N/A'),
    ('Gate Handle Set', 'fittings', 'Ergonomic handle set with anti-corrosion coating. Available in chrome or matte black.', '$30 - $80', 'Zinc Alloy', 'N/A'),
    ('Custom Cutting Service', 'cutting', 'Precision CNC cutting for custom gate panels and decorative elements.', '$100 - $800', 'Varies', 'Custom'),
    ('Laser Engraving', 'cutting', 'Laser engraving for logos, patterns, and custom designs on metal gates.', '$150 - $600', 'Metal', 'Custom'),
    ('Powder Coating', 'cutting', 'Professional powder coating service. Full color range, weather and UV resistant.', '$80 - $300', 'Coating', 'N/A'),
    ('Automatic Gate Motor', 'frame', 'Swing gate motor with 24V DC system. Supports gates up to 400kg per leaf.', '$600 - $1,200', 'Steel/Motor', 'N/A'),
    ('Access Control Panel', 'fittings', 'Smart access control with RFID cards, PIN pad, and mobile app integration.', '$300 - $800', 'Electronics', 'N/A'),
    ('Decorative Scroll Pack', 'crown', 'Pack of decorative iron scrolls for custom gate ornamentation.', '$40 - $150', 'Wrought Iron', '20cm x 30cm'),
    ('Safety Spike Set', 'crown', 'Anti-climb spike set for security gates. Available in pointed or rounded tips.', '$60 - $200', 'Hardened Steel', '15cm per spike')
ON CONFLICT DO NOTHING;
