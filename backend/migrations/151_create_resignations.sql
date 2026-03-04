CREATE TABLE IF NOT EXISTS resignations (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    restaurant_id INT NOT NULL REFERENCES restaurants(id),
    employee_id INT NOT NULL REFERENCES employees(id),
    resignation_date DATE NOT NULL,
    last_working_date DATE NOT NULL,
    reason TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'submitted',
    notice_period_days INT NOT NULL DEFAULT 0,
    pending_salary NUMERIC(12,2) NOT NULL DEFAULT 0,
    unused_leave_balance NUMERIC(12,2) NOT NULL DEFAULT 0,
    settlement_amount NUMERIC(12,2) NOT NULL DEFAULT 0,
    acknowledged_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL REFERENCES users(id),
    updated_by INT NULL REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_resignations_tenant ON resignations(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_resignations_employee ON resignations(tenant_id, restaurant_id, employee_id);
CREATE INDEX IF NOT EXISTS idx_resignations_status ON resignations(tenant_id, restaurant_id, status);
