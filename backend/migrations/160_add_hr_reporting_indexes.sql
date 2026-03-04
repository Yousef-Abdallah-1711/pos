-- Migration: Add indexes for HR reporting filters and aggregations
-- Purpose:
--   1) Speed up HR reports by tenant/restaurant/date scoped scans
--   2) Reduce full table scans on attendance/salary/performance queries

CREATE INDEX IF NOT EXISTS idx_employees_scope_active_status_dept
  ON employees (tenant_id, restaurant_id, is_active, employment_status, department);

CREATE INDEX IF NOT EXISTS idx_employees_scope_hire_termination
  ON employees (tenant_id, restaurant_id, hire_date, termination_date);

CREATE INDEX IF NOT EXISTS idx_attendance_scope_date_status
  ON attendance (tenant_id, restaurant_id, attendance_date, status);

CREATE INDEX IF NOT EXISTS idx_attendance_scope_employee_date
  ON attendance (tenant_id, restaurant_id, employee_id, attendance_date);

CREATE INDEX IF NOT EXISTS idx_salaries_scope_period_status
  ON salaries (tenant_id, restaurant_id, pay_period_start, pay_period_end, status);

CREATE INDEX IF NOT EXISTS idx_performance_reviews_scope_date_status
  ON performance_reviews (tenant_id, restaurant_id, review_date, status);

CREATE INDEX IF NOT EXISTS idx_training_records_scope_status_expiry
  ON training_records (tenant_id, restaurant_id, status, certificate_expiry);

CREATE INDEX IF NOT EXISTS idx_leaves_scope_year_status
  ON leaves (tenant_id, restaurant_id, start_date, status);
