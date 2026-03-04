CREATE TABLE IF NOT EXISTS tax_periods (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    restaurant_id INT NOT NULL REFERENCES restaurants(id),
    period_name VARCHAR(20) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    output_vat NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    input_vat NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    net_vat NUMERIC(15,2) NOT NULL DEFAULT 0.00,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    closed_at TIMESTAMP NULL,
    closed_by INT NULL REFERENCES users(id),
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, restaurant_id, period_name)
);

CREATE INDEX IF NOT EXISTS idx_tax_periods_tenant_restaurant
ON tax_periods(tenant_id, restaurant_id);

CREATE INDEX IF NOT EXISTS idx_tax_periods_status
ON tax_periods(tenant_id, restaurant_id, status);

CREATE INDEX IF NOT EXISTS idx_tax_periods_range
ON tax_periods(tenant_id, restaurant_id, period_start, period_end);

