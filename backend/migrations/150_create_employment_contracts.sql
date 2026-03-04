CREATE TABLE IF NOT EXISTS employment_contracts (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    restaurant_id INT NOT NULL REFERENCES restaurants(id),
    employee_id INT NOT NULL REFERENCES employees(id),
    contract_type VARCHAR(50) NOT NULL,
    contract_number VARCHAR(120) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    salary_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency VARCHAR(10) NOT NULL DEFAULT 'EGP',
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    renewed_from_id INT NULL REFERENCES employment_contracts(id),
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL REFERENCES users(id),
    updated_by INT NULL REFERENCES users(id),
    UNIQUE (tenant_id, restaurant_id, contract_number)
);

CREATE INDEX IF NOT EXISTS idx_employment_contracts_tenant ON employment_contracts(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_employment_contracts_employee ON employment_contracts(tenant_id, restaurant_id, employee_id);
CREATE INDEX IF NOT EXISTS idx_employment_contracts_end_date ON employment_contracts(tenant_id, restaurant_id, end_date);
