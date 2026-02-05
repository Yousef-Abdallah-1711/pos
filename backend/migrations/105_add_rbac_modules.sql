-- Migration: Add missing RBAC modules
-- Description: Inserts Orders, Leaves, and Salaries modules if missing

INSERT INTO modules (name, display_name, description, path) VALUES
('ORDERS', 'Orders', 'Order management and fulfillment', '/dashboard/orders')
ON CONFLICT DO NOTHING;

INSERT INTO modules (name, display_name, description, path) VALUES
('LEAVES', 'Leaves', 'Employee leave management', '/dashboard/hr/leaves')
ON CONFLICT DO NOTHING;

INSERT INTO modules (name, display_name, description, path) VALUES
('SALARIES', 'Payroll', 'Payroll and salary management', '/dashboard/hr/payroll')
ON CONFLICT DO NOTHING;
