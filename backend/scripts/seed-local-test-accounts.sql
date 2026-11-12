-- Development-only accounts for the local POS Docker stack.
-- Run only through seed-local-test-accounts.ps1, which verifies the Compose project.
-- Passwords are supplied interactively by the PowerShell wrapper and only hashes persist.
\getenv seed_superadmin_password POS_LOCAL_SUPERADMIN_PASSWORD
\getenv seed_owner_password POS_LOCAL_OWNER_PASSWORD
\getenv seed_cashier_password POS_LOCAL_CASHIER_PASSWORD
SET pos.seed_superadmin_password TO :'seed_superadmin_password';
SET pos.seed_owner_password TO :'seed_owner_password';
SET pos.seed_cashier_password TO :'seed_cashier_password';
BEGIN;

DO $seed$
DECLARE
    demo_tenant_id INTEGER;
    demo_restaurant_id INTEGER;
    super_admin_user_id INTEGER;
    owner_user_id INTEGER;
    cashier_user_id INTEGER;
    admin_role_id INTEGER;
    employee_role_id INTEGER;
BEGIN
    IF current_database() <> 'pos_saas' THEN
        RAISE EXCEPTION 'Refusing to seed non-development database: %', current_database();
    END IF;

    SELECT id INTO demo_tenant_id
    FROM tenants
    WHERE slug = 'restaurant-1' AND name = 'Test Tenant' AND is_active = TRUE;

    IF demo_tenant_id IS NULL THEN
        RAISE EXCEPTION 'Expected local Test Tenant (slug restaurant-1) was not found';
    END IF;

    SELECT id INTO demo_restaurant_id
    FROM restaurants
    WHERE tenant_id = demo_tenant_id AND slug = 'restaurant-1' AND name = 'Test Restaurant';

    IF demo_restaurant_id IS NULL THEN
        RAISE EXCEPTION 'Expected local Test Restaurant was not found';
    END IF;

    SELECT id INTO admin_role_id FROM roles
    WHERE tenant_id = demo_tenant_id AND restaurant_id = demo_restaurant_id AND role_code = 'ADMIN' AND is_active = TRUE;

    SELECT id INTO employee_role_id FROM roles
    WHERE tenant_id = demo_tenant_id AND restaurant_id = demo_restaurant_id AND role_code = 'EMPLOYEE' AND is_active = TRUE;

    IF admin_role_id IS NULL OR employee_role_id IS NULL THEN
        RAISE EXCEPTION 'Expected active ADMIN and EMPLOYEE roles were not found; apply migrations first';
    END IF;

    INSERT INTO users (
        tenant_id, restaurant_id, email, password_hash, name, phone, role, status, is_super_admin, created_at, updated_at
    ) VALUES (
        demo_tenant_id, NULL, 'superadmin@pos-saas.com',
        crypt(current_setting('pos.seed_superadmin_password'), gen_salt('bf')),
        'Local POS Super Admin', '+15550001001', 'super_admin', 'active', TRUE, NOW(), NOW()
    )
    ON CONFLICT (tenant_id, email) DO UPDATE SET
        restaurant_id = NULL,
        password_hash = EXCLUDED.password_hash,
        name = EXCLUDED.name,
        phone = EXCLUDED.phone,
        role = EXCLUDED.role,
        status = 'active',
        is_super_admin = TRUE,
        updated_at = NOW()
    RETURNING id INTO super_admin_user_id;

    INSERT INTO super_admins (user_id, ip_whitelist, is_active, failed_login_attempts, locked_until, created_at, updated_at)
    VALUES (super_admin_user_id, NULL, TRUE, 0, NULL, NOW(), NOW())
    ON CONFLICT (user_id) DO UPDATE SET
        ip_whitelist = NULL,
        is_active = TRUE,
        failed_login_attempts = 0,
        locked_until = NULL,
        updated_at = NOW();

    INSERT INTO users (
        tenant_id, restaurant_id, email, password_hash, name, phone, role, status, is_super_admin, created_at, updated_at
    ) VALUES (
        demo_tenant_id, demo_restaurant_id, 'owner@pos-demo.local',
        crypt(current_setting('pos.seed_owner_password'), gen_salt('bf')),
        'Local Demo Owner', '+15550001002', 'owner', 'active', FALSE, NOW(), NOW()
    )
    ON CONFLICT (tenant_id, email) DO UPDATE SET
        restaurant_id = EXCLUDED.restaurant_id,
        password_hash = EXCLUDED.password_hash,
        name = EXCLUDED.name,
        phone = EXCLUDED.phone,
        role = EXCLUDED.role,
        status = 'active',
        is_super_admin = FALSE,
        updated_at = NOW()
    RETURNING id INTO owner_user_id;

    INSERT INTO users (
        tenant_id, restaurant_id, email, password_hash, name, phone, role, status, is_super_admin, created_at, updated_at
    ) VALUES (
        demo_tenant_id, demo_restaurant_id, 'cashier@pos-demo.local',
        crypt(current_setting('pos.seed_cashier_password'), gen_salt('bf')),
        'Local Demo Cashier', '+15550001003', 'cashier', 'active', FALSE, NOW(), NOW()
    )
    ON CONFLICT (tenant_id, email) DO UPDATE SET
        restaurant_id = EXCLUDED.restaurant_id,
        password_hash = EXCLUDED.password_hash,
        name = EXCLUDED.name,
        phone = EXCLUDED.phone,
        role = EXCLUDED.role,
        status = 'active',
        is_super_admin = FALSE,
        updated_at = NOW()
    RETURNING id INTO cashier_user_id;

    INSERT INTO user_roles (tenant_id, user_id, role_id, assigned_by, created_at, updated_at)
    VALUES (demo_tenant_id, owner_user_id, admin_role_id, owner_user_id, NOW(), NOW())
    ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;

    INSERT INTO user_roles (tenant_id, user_id, role_id, assigned_by, created_at, updated_at)
    VALUES (demo_tenant_id, cashier_user_id, employee_role_id, owner_user_id, NOW(), NOW())
    ON CONFLICT (tenant_id, user_id, role_id) DO NOTHING;
END
$seed$;

COMMIT;
