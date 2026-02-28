-- Migration: Create table_layouts and attach layout_id to restaurant_tables
-- Description: Stores floor/section canvas metadata and groups tables by layout.

CREATE TABLE IF NOT EXISTS table_layouts (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL,
    restaurant_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL DEFAULT 'Main Floor',
    floor_number INTEGER NOT NULL DEFAULT 1,
    canvas_width DOUBLE PRECISION DEFAULT 800,
    canvas_height DOUBLE PRECISION DEFAULT 600,
    background_image TEXT,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_table_layouts_tenant
        FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_table_layouts_restaurant
        FOREIGN KEY (restaurant_id) REFERENCES restaurants(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_table_layouts_tenant_restaurant
    ON table_layouts(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_table_layouts_restaurant_floor
    ON table_layouts(restaurant_id, floor_number);

ALTER TABLE restaurant_tables
ADD COLUMN IF NOT EXISTS layout_id BIGINT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_restaurant_tables_layout_id'
    ) THEN
        ALTER TABLE restaurant_tables
        ADD CONSTRAINT fk_restaurant_tables_layout_id
            FOREIGN KEY (layout_id) REFERENCES table_layouts(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_restaurant_tables_layout_id
    ON restaurant_tables(layout_id);
