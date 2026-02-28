-- Migration: Add table_id to orders
-- Description: Links orders to restaurant tables for dine-in tracking.

ALTER TABLE orders
ADD COLUMN IF NOT EXISTS table_id BIGINT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_orders_table'
    ) THEN
        ALTER TABLE orders
        ADD CONSTRAINT fk_orders_table
            FOREIGN KEY (table_id) REFERENCES restaurant_tables(id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_orders_table_id ON orders(table_id);
