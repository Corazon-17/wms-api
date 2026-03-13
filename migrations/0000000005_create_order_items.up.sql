CREATE TABLE order_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    sku TEXT NOT NULL,
    quantity INT NOT NULL,
    price NUMERIC(12,2) NOT NULL
);

CREATE INDEX idx_order_items_order_id ON order_items (order_id);