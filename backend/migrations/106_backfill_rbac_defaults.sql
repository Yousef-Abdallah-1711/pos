-- Migration: Backfill RBAC roles, permissions, and default user roles
-- Description: Ensures system roles, role permissions, and owner assignments exist for existing tenants

WITH tenant_restaurants AS (
    SELECT t.id AS tenant_id, MIN(r.id) AS restaurant_id
    FROM tenants t
    JOIN restaurants r ON r.tenant_id = t.id
    GROUP BY t.id
),
role_defaults AS (
    SELECT * FROM (VALUES
        ('ADMIN', 'Administrator', 'Full system access', 'admin', 1),
        ('MANAGER', 'Manager', 'Management access', 'manager', 2),
        ('EMPLOYEE', 'Employee', 'Staff access', 'basic', 3),
        ('DRIVER', 'Driver', 'Driver access to assigned orders', 'basic', 4)
    ) AS v(role_code, role_name, description, access_level, display_order)
),
missing_roles AS (
    SELECT tr.tenant_id,
           tr.restaurant_id,
           rd.role_code,
           rd.role_name,
           rd.description,
           rd.access_level,
           rd.display_order
    FROM tenant_restaurants tr
    CROSS JOIN role_defaults rd
    WHERE NOT EXISTS (
        SELECT 1
        FROM roles r
        WHERE r.tenant_id = tr.tenant_id
          AND r.role_code = rd.role_code
    )
)
INSERT INTO roles (
    tenant_id,
    restaurant_id,
    role_name,
    description,
    role_code,
    access_level,
    is_system_role,
    is_active,
    display_order,
    created_at,
    updated_at
)
SELECT tenant_id,
       restaurant_id,
       role_name,
       description,
       role_code,
       access_level,
       TRUE,
       TRUE,
       display_order,
       NOW(),
       NOW()
FROM missing_roles;

INSERT INTO role_permissions (
    tenant_id,
    role_id,
    module_id,
    permission_level,
    created_at,
    updated_at
)
SELECT r.tenant_id,
       r.id,
       m.id,
       CASE
           WHEN r.role_code = 'ADMIN' THEN 'ADMIN'
           WHEN r.role_code = 'MANAGER' THEN CASE
               WHEN m.name IN ('PRODUCTS', 'CATEGORIES', 'ORDERS', 'HR', 'LEAVES', 'NOTIFICATIONS') THEN 'WRITE'
               WHEN m.name IN ('SALARIES', 'SETTINGS', 'USERS', 'ROLES', 'REPORTS', 'THEMES') THEN 'READ'
               ELSE 'NONE'
           END
           WHEN r.role_code = 'EMPLOYEE' THEN CASE
               WHEN m.name IN ('PRODUCTS', 'CATEGORIES', 'ORDERS', 'HR', 'LEAVES', 'NOTIFICATIONS') THEN 'READ'
               ELSE 'NONE'
           END
           WHEN r.role_code = 'DRIVER' THEN CASE
               WHEN m.name IN ('ORDERS', 'NOTIFICATIONS') THEN 'READ'
               ELSE 'NONE'
           END
           ELSE 'NONE'
       END,
       NOW(),
       NOW()
FROM roles r
JOIN modules m ON m.is_active = TRUE
WHERE r.is_system_role = TRUE
  AND r.role_code IN ('ADMIN', 'MANAGER', 'EMPLOYEE', 'DRIVER')
ON CONFLICT (tenant_id, role_id, module_id)
DO UPDATE SET permission_level = EXCLUDED.permission_level,
              updated_at = EXCLUDED.updated_at;

INSERT INTO user_roles (
    tenant_id,
    user_id,
    role_id,
    assigned_by,
    created_at,
    updated_at
)
SELECT u.tenant_id,
       u.id,
       r.id,
       u.id,
       NOW(),
       NOW()
FROM users u
JOIN roles r ON r.tenant_id = u.tenant_id AND r.role_code = 'ADMIN'
WHERE u.role = 'owner'
ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;

INSERT INTO user_roles (
    tenant_id,
    user_id,
    role_id,
    assigned_by,
    created_at,
    updated_at
)
SELECT u.tenant_id,
       u.id,
       r.id,
       u.id,
       NOW(),
       NOW()
FROM users u
JOIN roles r ON r.tenant_id = u.tenant_id AND r.role_code = 'EMPLOYEE'
LEFT JOIN user_roles ur ON ur.user_id = u.id AND ur.tenant_id = u.tenant_id
WHERE ur.id IS NULL
ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;

INSERT INTO user_roles (
    tenant_id,
    user_id,
    role_id,
    assigned_by,
    created_at,
    updated_at
)
SELECT u.tenant_id,
       u.id,
       r.id,
       u.id,
       NOW(),
       NOW()
FROM users u
JOIN roles r ON r.tenant_id = u.tenant_id AND r.role_code = 'DRIVER'
WHERE u.role = 'driver'
ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;