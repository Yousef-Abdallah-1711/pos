CREATE TABLE IF NOT EXISTS performance_reviews (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    restaurant_id INT NOT NULL REFERENCES restaurants(id),
    employee_id INT NOT NULL REFERENCES employees(id),
    reviewer_id INT NOT NULL REFERENCES users(id),
    review_period_start DATE NOT NULL,
    review_period_end DATE NOT NULL,
    review_date DATE NOT NULL,
    goals_achievement NUMERIC(5,2) NOT NULL DEFAULT 0,
    technical_skills NUMERIC(5,2) NOT NULL DEFAULT 0,
    communication NUMERIC(5,2) NOT NULL DEFAULT 0,
    leadership NUMERIC(5,2) NOT NULL DEFAULT 0,
    overall_rating NUMERIC(5,2) NOT NULL DEFAULT 0,
    strengths TEXT,
    areas_to_improve TEXT,
    action_plan TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    acknowledged_at TIMESTAMP NULL,
    acknowledged_by INT NULL REFERENCES users(id),
    completed_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL REFERENCES users(id),
    updated_by INT NULL REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_performance_reviews_tenant ON performance_reviews(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_employee ON performance_reviews(tenant_id, restaurant_id, employee_id);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_status ON performance_reviews(tenant_id, restaurant_id, status);
