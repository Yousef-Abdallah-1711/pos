-- Migration: Create fixed assets register for accounting module
-- Purpose: Track long-term assets and depreciation inputs for reporting/depreciation postings.

CREATE TABLE IF NOT EXISTS fixed_assets (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    restaurant_id INT NOT NULL REFERENCES restaurants(id),
    asset_code VARCHAR(30) NOT NULL,
    asset_name VARCHAR(150) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'general',
    purchase_date DATE NOT NULL,
    capitalization_date DATE NULL,
    purchase_cost NUMERIC(14,2) NOT NULL CHECK (purchase_cost >= 0),
    salvage_value NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (salvage_value >= 0),
    useful_life_months INT NOT NULL CHECK (useful_life_months > 0),
    depreciation_method VARCHAR(20) NOT NULL DEFAULT 'straight_line',
    depreciation_start_date DATE NULL,
    accumulated_depreciation NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (accumulated_depreciation >= 0),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (tenant_id, restaurant_id, asset_code)
);

CREATE INDEX IF NOT EXISTS idx_fixed_assets_tenant_restaurant
    ON fixed_assets(tenant_id, restaurant_id);

CREATE INDEX IF NOT EXISTS idx_fixed_assets_category
    ON fixed_assets(tenant_id, restaurant_id, category);

CREATE INDEX IF NOT EXISTS idx_fixed_assets_active
    ON fixed_assets(tenant_id, restaurant_id, is_active);

