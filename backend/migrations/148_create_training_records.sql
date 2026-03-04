CREATE TABLE IF NOT EXISTS training_records (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    restaurant_id INT NOT NULL REFERENCES restaurants(id),
    employee_id INT NOT NULL REFERENCES employees(id),
    training_title VARCHAR(255) NOT NULL,
    training_provider VARCHAR(255),
    training_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    certificate_number VARCHAR(100),
    certificate_expiry DATE,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL REFERENCES users(id),
    updated_by INT NULL REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_training_records_tenant ON training_records(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_training_records_employee ON training_records(tenant_id, restaurant_id, employee_id);
CREATE INDEX IF NOT EXISTS idx_training_records_certificate_expiry ON training_records(tenant_id, restaurant_id, certificate_expiry);
