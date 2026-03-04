CREATE TABLE IF NOT EXISTS fiscal_periods (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    restaurant_id INT NOT NULL REFERENCES restaurants(id),
    period_name VARCHAR(20) NOT NULL,
    period_type VARCHAR(10) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'open',
    closed_at TIMESTAMP NULL,
    closed_by INT NULL REFERENCES users(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, restaurant_id, period_name)
);

CREATE INDEX IF NOT EXISTS idx_fiscal_periods_tenant_restaurant
ON fiscal_periods(tenant_id, restaurant_id);

CREATE INDEX IF NOT EXISTS idx_fiscal_periods_dates
ON fiscal_periods(tenant_id, restaurant_id, start_date, end_date);

CREATE INDEX IF NOT EXISTS idx_fiscal_periods_status
ON fiscal_periods(tenant_id, restaurant_id, status);

