-- Migration: Create restaurant_tables table
-- Description: Stores physical table metadata and floor positioning for dine-in flows.

CREATE TABLE IF NOT EXISTS restaurant_tables (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    restaurant_id BIGINT NOT NULL,
    table_number INTEGER NOT NULL,
    name VARCHAR(100),
    capacity INTEGER DEFAULT 4,
    location VARCHAR(50) DEFAULT 'indoor',
    status VARCHAR(20) DEFAULT 'available',
    description TEXT,
    qr_code_config JSONB DEFAULT '{}'::jsonb,
    position_x DOUBLE PRECISION DEFAULT 0,
    position_y DOUBLE PRECISION DEFAULT 0,
    shape VARCHAR(20) DEFAULT 'square',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_restaurant_tables_tenant
        FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_restaurant_tables_restaurant
        FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE,
    CONSTRAINT uq_restaurant_tables_number
        UNIQUE (restaurant_id, table_number),
    CONSTRAINT chk_restaurant_tables_status
        CHECK (status IN ('available', 'occupied', 'reserved', 'maintenance')),
    CONSTRAINT chk_restaurant_tables_location
        CHECK (location IN ('indoor', 'outdoor', 'vip', 'rooftop', 'private')),
    CONSTRAINT chk_restaurant_tables_capacity
        CHECK (capacity > 0 AND capacity <= 50),
    CONSTRAINT chk_restaurant_tables_shape
        CHECK (shape IN ('square', 'circle', 'rectangle'))
);

CREATE INDEX IF NOT EXISTS idx_restaurant_tables_tenant_restaurant
    ON restaurant_tables(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_restaurant_tables_status
    ON restaurant_tables(restaurant_id, status);
CREATE INDEX IF NOT EXISTS idx_restaurant_tables_number
    ON restaurant_tables(restaurant_id, table_number);
