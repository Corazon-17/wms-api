CREATE TABLE wms_statuses (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

INSERT INTO wms_statuses (id, name) VALUES
    ('READY_TO_PICK', 'Ready To Pick'),
    ('PICKING', 'Picking'),
    ('PACKED', 'Packed'),
    ('SHIPPED', 'Shipped');