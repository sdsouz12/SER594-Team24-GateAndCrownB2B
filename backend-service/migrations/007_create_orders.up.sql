CREATE TABLE IF NOT EXISTS orders (
    order_id  BIGSERIAL PRIMARY KEY,
    user_id   BIGINT      NOT NULL REFERENCES sys_user(user_id),
    product_id BIGINT     NOT NULL REFERENCES catalog_product(product_id),
    quantity  INTEGER     NOT NULL DEFAULT 1,
    notes     TEXT,
    status    VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP  NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
