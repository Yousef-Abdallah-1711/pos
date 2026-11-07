param(
    [string]$ProjectName = 'pos-p2a-ws05-e2e-20261103'
)

$ErrorActionPreference = 'Stop'

$project = $ProjectName
$api = 'http://localhost:18081/api/v1'
$database = 'pos_phase2a_e2e'
$container = "$project-postgres-1"
$password = 'WS05-Disposable-Only-2026!'
$superAdminEmail = 'phase2a-superadmin@example.test'
$superAdminPassword = 'WS05-SuperAdmin-Only-2026!'
$superAdminPasswordHash = '$2a$10$cZ6RYwpViOspI8u.NYyXkeCD1IOXHF0xussswOrgijec47mn2OVEe'

$labels = docker inspect --format '{{json .Config.Labels}}' $container | ConvertFrom-Json
if ($LASTEXITCODE -ne 0 -or $labels.'com.docker.compose.project' -ne $project -or $labels.'com.docker.compose.service' -ne 'postgres' -or $labels.'com.pos.disposable' -ne 'true') {
    throw "Refusing to seed: $container is not the labeled WS05 disposable PostgreSQL service."
}
$databaseName = docker exec $container psql -U postgres -d $database -Atc 'SELECT current_database()'
if ($LASTEXITCODE -ne 0 -or $databaseName.Trim() -ne $database) {
    throw "Refusing to seed: expected disposable database '$database'."
}
$health = Invoke-RestMethod -Uri 'http://localhost:18081/api/v1/health' -TimeoutSec 10
if (-not $health) { throw 'The isolated API health check failed.' }

function Invoke-DisposableSql([string]$sql) {
    $result = $sql | docker exec -i $container psql -U postgres -d $database -v ON_ERROR_STOP=1 -At
    if ($LASTEXITCODE -ne 0) { throw 'Disposable WS05 fixture SQL failed.' }
    return ($result -join "`n").Trim()
}

function Ensure-Tenant([string]$suffix, [string]$name) {
    $restaurantSlug = "phase2a-restaurant-$suffix"
    $tenantSlug = $restaurantSlug
    $email = "phase2a-owner-$suffix@example.test"
    $tenantID = Invoke-DisposableSql "SELECT id FROM tenants WHERE slug = '$tenantSlug' LIMIT 1"
    if (-not $tenantID) {
        $registration = @{
            name = "WS05 $name Owner"
            email = $email
            password = $password
            restaurant_name = "WS05 $name Restaurant"
            restaurant_slug = $restaurantSlug
            phone = '5550100200'
        } | ConvertTo-Json
        try {
            $null = Invoke-RestMethod -Method Post -Uri "$api/auth/register" -ContentType 'application/json' -Headers @{ Origin = 'http://localhost:13102' } -Body $registration -TimeoutSec 20
        } catch {
            throw "Could not create disposable WS05 $name tenant account: $($_.Exception.Message)"
        }
        $tenantID = Invoke-DisposableSql "SELECT id FROM tenants WHERE slug = '$tenantSlug' LIMIT 1"
    }
    $restaurantID = Invoke-DisposableSql "SELECT id FROM restaurants WHERE tenant_id = $tenantID AND slug = '$restaurantSlug' LIMIT 1"
    $user = Invoke-DisposableSql "SELECT id FROM users WHERE tenant_id = $tenantID AND email = '$email' AND status = 'active' LIMIT 1"
    if (-not $restaurantID -or -not $user) { throw "The disposable WS05 $name account fixture is incomplete." }

    $login = @{ email = $email; password = $password } | ConvertTo-Json
    $auth = Invoke-WebRequest -Method Post -Uri "$api/auth/login" -ContentType 'application/json' -Headers @{ Origin = 'http://localhost:13102' } -Body $login -TimeoutSec 20
    if ($auth.StatusCode -ne 200 -or -not ($auth.Headers['Set-Cookie'] -match 'poos_session=')) {
        throw "Could not establish the disposable WS05 $name browser session."
    }
    return [pscustomobject]@{ TenantID = [int]$tenantID; RestaurantID = [int]$restaurantID; UserID = [int]$user; Email = $email }
}

$tenantA = Ensure-Tenant 'a' 'Tenant A'
$tenantB = Ensure-Tenant 'b' 'Tenant B'

$fixtureSql = @"
BEGIN;
INSERT INTO categories (tenant_id, restaurant_id, name, display_order, is_active, created_at, updated_at)
SELECT $($tenantA.TenantID), $($tenantA.RestaurantID), 'WS05 E2E Menu', 1, TRUE, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE tenant_id = $($tenantA.TenantID) AND restaurant_id = $($tenantA.RestaurantID) AND name = 'WS05 E2E Menu');
INSERT INTO categories (tenant_id, restaurant_id, name, display_order, is_active, created_at, updated_at)
SELECT $($tenantB.TenantID), $($tenantB.RestaurantID), 'WS05 E2E Menu', 1, TRUE, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM categories WHERE tenant_id = $($tenantB.TenantID) AND restaurant_id = $($tenantB.RestaurantID) AND name = 'WS05 E2E Menu');
INSERT INTO products (tenant_id, restaurant_id, category_id, name, name_en, name_ar, description, price, available, is_available, quantity_in_stock, track_inventory, low_stock_threshold, status, created_at, updated_at)
SELECT $($tenantA.TenantID), $($tenantA.RestaurantID), c.id, 'WS05 E2E Burger A', 'WS05 E2E Burger A', 'وجبة اختبار أ', 'Disposable WS05 fixture', 12.50, TRUE, TRUE, 20, TRUE, 3, 'active', NOW(), NOW()
FROM categories c WHERE c.tenant_id = $($tenantA.TenantID) AND c.restaurant_id = $($tenantA.RestaurantID) AND c.name = 'WS05 E2E Menu'
AND NOT EXISTS (SELECT 1 FROM products p WHERE p.tenant_id = $($tenantA.TenantID) AND p.restaurant_id = $($tenantA.RestaurantID) AND p.name = 'WS05 E2E Burger A');
INSERT INTO products (tenant_id, restaurant_id, category_id, name, name_en, name_ar, description, price, available, is_available, quantity_in_stock, track_inventory, low_stock_threshold, status, created_at, updated_at)
SELECT $($tenantB.TenantID), $($tenantB.RestaurantID), c.id, 'WS05 E2E Burger B', 'WS05 E2E Burger B', 'وجبة اختبار ب', 'Disposable WS05 fixture', 18.75, TRUE, TRUE, 20, TRUE, 3, 'active', NOW(), NOW()
FROM categories c WHERE c.tenant_id = $($tenantB.TenantID) AND c.restaurant_id = $($tenantB.RestaurantID) AND c.name = 'WS05 E2E Menu'
AND NOT EXISTS (SELECT 1 FROM products p WHERE p.tenant_id = $($tenantB.TenantID) AND p.restaurant_id = $($tenantB.RestaurantID) AND p.name = 'WS05 E2E Burger B');
UPDATE products SET name_en = COALESCE(name_en, name), name_ar = COALESCE(name_ar, name), quantity_in_stock = 20, track_inventory = TRUE, low_stock_threshold = 3 WHERE tenant_id IN ($($tenantA.TenantID), $($tenantB.TenantID)) AND name IN ('WS05 E2E Burger A', 'WS05 E2E Burger B');
WITH fixture_super_admin AS (
    INSERT INTO users (tenant_id, restaurant_id, email, password_hash, name, phone, role, status, is_super_admin, created_at, updated_at)
    VALUES ($($tenantA.TenantID), NULL, '$superAdminEmail', '$superAdminPasswordHash', 'WS05 Disposable Super Admin', '5550100299', 'super_admin', 'active', TRUE, NOW(), NOW())
    ON CONFLICT (tenant_id, email) DO UPDATE SET
        restaurant_id = NULL,
        password_hash = EXCLUDED.password_hash,
        name = EXCLUDED.name,
        role = EXCLUDED.role,
        status = 'active',
        is_super_admin = TRUE,
        updated_at = NOW()
    RETURNING id
)
INSERT INTO super_admins (user_id, ip_whitelist, is_active, failed_login_attempts, locked_until, created_at, updated_at)
SELECT id, NULL, TRUE, 0, NULL, NOW(), NOW() FROM fixture_super_admin
ON CONFLICT (user_id) DO UPDATE SET
    ip_whitelist = NULL,
    is_active = TRUE,
    failed_login_attempts = 0,
    locked_until = NULL,
    updated_at = NOW();
COMMIT;
"@
$null = Invoke-DisposableSql $fixtureSql

Write-Output "Seeded isolated WS05 fixtures: Tenant A $($tenantA.TenantID), Tenant B $($tenantB.TenantID)."
Write-Output "Test-only account for both tenants: phase2a-owner-a@example.test / phase2a-owner-b@example.test"
Write-Output "Test-only password: $password"
Write-Output "Test-only super-admin account: $superAdminEmail"
Write-Output "Test-only super-admin password: $superAdminPassword"
