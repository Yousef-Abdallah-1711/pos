-- Phase 7-9 extensions: advanced accounting, analytics/compliance support.

CREATE TABLE IF NOT EXISTS bank_reconciliation_sessions (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    bank_account_id INT NOT NULL REFERENCES bank_accounts_ledger(id) ON DELETE CASCADE,
    statement_date DATE NOT NULL,
    statement_balance NUMERIC(15,2) NOT NULL,
    book_balance NUMERIC(15,2) NOT NULL DEFAULT 0,
    difference NUMERIC(15,2) NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    notes TEXT,
    completed_at TIMESTAMP NULL,
    completed_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    updated_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_bank_reconciliation_status CHECK (status IN ('draft', 'in_progress', 'completed', 'cancelled'))
);

CREATE INDEX IF NOT EXISTS idx_bank_reconciliation_sessions_tenant ON bank_reconciliation_sessions(tenant_id, restaurant_id);
CREATE INDEX IF NOT EXISTS idx_bank_reconciliation_sessions_account_status ON bank_reconciliation_sessions(tenant_id, restaurant_id, bank_account_id, status);

CREATE TABLE IF NOT EXISTS bank_reconciliation_lines (
    id SERIAL PRIMARY KEY,
    reconciliation_id INT NOT NULL REFERENCES bank_reconciliation_sessions(id) ON DELETE CASCADE,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    journal_entry_id INT NULL REFERENCES journal_entries(id) ON DELETE SET NULL,
    line_description TEXT,
    statement_amount NUMERIC(15,2) NOT NULL,
    book_amount NUMERIC(15,2) NOT NULL,
    match_status VARCHAR(20) NOT NULL DEFAULT 'matched',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_bank_reconciliation_line_status CHECK (match_status IN ('matched', 'unmatched', 'manual'))
);

CREATE INDEX IF NOT EXISTS idx_bank_reconciliation_lines_reconciliation ON bank_reconciliation_lines(reconciliation_id);

CREATE TABLE IF NOT EXISTS recurring_journal_templates (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    template_code VARCHAR(60) NOT NULL,
    template_name VARCHAR(255) NOT NULL,
    description TEXT,
    frequency VARCHAR(20) NOT NULL,
    start_date DATE NOT NULL,
    next_run_date DATE NOT NULL,
    last_run_date DATE NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    auto_post BOOLEAN NOT NULL DEFAULT FALSE,
    lines_json JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    updated_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE (tenant_id, restaurant_id, template_code),
    CONSTRAINT chk_recurring_journal_frequency CHECK (frequency IN ('weekly', 'monthly', 'quarterly', 'yearly'))
);

CREATE INDEX IF NOT EXISTS idx_recurring_journal_templates_next_run ON recurring_journal_templates(tenant_id, restaurant_id, next_run_date) WHERE is_active = TRUE;

CREATE TABLE IF NOT EXISTS account_budgets (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    account_id INT NOT NULL REFERENCES chart_of_accounts(id) ON DELETE CASCADE,
    fiscal_year INT NOT NULL,
    period_month INT NOT NULL,
    budget_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    actual_amount NUMERIC(15,2) NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'EGP',
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, restaurant_id, account_id, fiscal_year, period_month),
    CONSTRAINT chk_account_budgets_month CHECK (period_month BETWEEN 1 AND 12)
);

CREATE INDEX IF NOT EXISTS idx_account_budgets_year_month ON account_budgets(tenant_id, restaurant_id, fiscal_year, period_month);

CREATE TABLE IF NOT EXISTS exchange_rates (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    base_currency VARCHAR(3) NOT NULL,
    quote_currency VARCHAR(3) NOT NULL,
    rate NUMERIC(18,8) NOT NULL,
    effective_date DATE NOT NULL,
    source VARCHAR(80),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, restaurant_id, base_currency, quote_currency, effective_date),
    CONSTRAINT chk_exchange_rates_rate_positive CHECK (rate > 0)
);

CREATE INDEX IF NOT EXISTS idx_exchange_rates_pair_date ON exchange_rates(tenant_id, restaurant_id, base_currency, quote_currency, effective_date DESC);

CREATE TABLE IF NOT EXISTS supplier_payment_schedules (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    supplier_invoice_id INT NULL REFERENCES invoices_supplier(id) ON DELETE SET NULL,
    supplier_id INT NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    scheduled_date DATE NOT NULL,
    amount NUMERIC(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'EGP',
    reminder_days INT NOT NULL DEFAULT 7,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    reminder_sent_at TIMESTAMP NULL,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    updated_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_supplier_payment_schedule_status CHECK (status IN ('scheduled', 'reminded', 'paid', 'cancelled')),
    CONSTRAINT chk_supplier_payment_schedule_amount CHECK (amount > 0)
);

CREATE INDEX IF NOT EXISTS idx_supplier_payment_schedule_due ON supplier_payment_schedules(tenant_id, restaurant_id, scheduled_date, status);

CREATE TABLE IF NOT EXISTS invoice_delivery_logs (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    invoice_id INT NOT NULL REFERENCES invoices_customer(id) ON DELETE CASCADE,
    channel VARCHAR(20) NOT NULL DEFAULT 'email',
    recipient_email VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'sent',
    sent_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_invoice_delivery_status CHECK (status IN ('sent', 'queued', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_invoice_delivery_logs_invoice ON invoice_delivery_logs(tenant_id, restaurant_id, invoice_id, sent_at DESC);

CREATE TABLE IF NOT EXISTS e_invoice_submissions (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    invoice_type VARCHAR(20) NOT NULL,
    invoice_id INT NOT NULL,
    eta_uuid VARCHAR(120),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    submitted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    response_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    error_message TEXT,
    created_by INT NULL REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT chk_e_invoice_submissions_type CHECK (invoice_type IN ('customer', 'supplier')),
    CONSTRAINT chk_e_invoice_submissions_status CHECK (status IN ('pending', 'accepted', 'rejected', 'error'))
);

CREATE INDEX IF NOT EXISTS idx_e_invoice_submissions_invoice ON e_invoice_submissions(tenant_id, restaurant_id, invoice_type, invoice_id);

CREATE TABLE IF NOT EXISTS audit_export_logs (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    restaurant_id INT NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    export_type VARCHAR(40) NOT NULL,
    export_format VARCHAR(20) NOT NULL,
    params_json JSONB NOT NULL DEFAULT '{}'::jsonb,
    row_count INT NOT NULL DEFAULT 0,
    generated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    generated_by INT NULL REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_export_logs_tenant ON audit_export_logs(tenant_id, restaurant_id, generated_at DESC);
