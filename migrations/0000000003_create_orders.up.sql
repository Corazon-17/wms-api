CREATE TABLE orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_sn TEXT NOT NULL UNIQUE,
    shop_id TEXT NOT NULL,
    marketplace_status TEXT,
    shipping_status TEXT,
    wms_status_id TEXT NOT NULL REFERENCES wms_statuses(id) ON DELETE SET NULL,
    tracking_number TEXT,
    total_amount NUMERIC(12,2),
    raw_marketplace_payload JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_wms_status ON orders (wms_status_id);
CREATE INDEX idx_orders_tracking_number ON orders (tracking_number);
CREATE INDEX idx_orders_updated_at ON orders (updated_at DESC);