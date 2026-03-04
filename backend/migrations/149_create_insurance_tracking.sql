CREATE TABLE IF NOT EXISTS insurance_tracking (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    restaurant_id INT NOT NULL REFERENCES restaurants(id),
    employee_id INT NOT NULL REFERENCES employees(id),
    policy_type VARCHAR(50) NOT NULL,
    provider_name VARCHAR(255) NOT NULL,
    policy_number VARCHAR(120) NOT NULL,
    coverage_start DATE NOT NULL,
    coverage_end DATE NOT NULL,
    monthly_premium NUMERIC(12,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL REFERENCES users(id),
    updated_by INT NULL REFERENCES users(id),
    UNIQUE (tenant_id, restaurant_id, policy_number)
);

CREATE INDEX IF NOT EXISTS idx_insurance_tracking_tenant ON insurance_tracking(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_insurance_tracking_employee ON insurance_tracking(tenant_id, restaurant_id, employee_id);
CREATE INDEX IF NOT EXISTS idx_insurance_tracking_coverage_end ON insurance_tracking(tenant_id, restaurant_id, coverage_end);
