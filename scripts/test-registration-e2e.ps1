$ErrorActionPreference = 'Stop'

$repoRoot = (Resolve-Path (Join-Path $PSScriptRoot '..')).Path
$projectName = "pos-registration-e2e-$([guid]::NewGuid().ToString('N').Substring(0, 8))"
$baseCompose = Join-Path $repoRoot 'docker-compose.phase2a-e2e.yml'
$overrideCompose = Join-Path $repoRoot 'docker-compose.registration-e2e.override.yml'
$fixturePath = Join-Path $env:TEMP "$projectName-fixture.json"
$outputPath = Join-Path $env:TEMP "$projectName-playwright"
$started = $false

function Get-FreeTcpPort {
    $listener = [System.Net.Sockets.TcpListener]::new([System.Net.IPAddress]::Loopback, 0)
    $listener.Start()
    try { return ([System.Net.IPEndPoint]$listener.LocalEndpoint).Port }
    finally { $listener.Stop() }
}

function Invoke-Compose([string[]]$Arguments) {
    & docker compose --project-name $projectName --file $baseCompose --file $overrideCompose @Arguments
    if ($LASTEXITCODE -ne 0) { throw "docker compose failed with exit code $LASTEXITCODE" }
}

try {
    $REG_E2E_PG_PORT = Get-FreeTcpPort
    $REG_E2E_API_PORT = Get-FreeTcpPort
    $REG_E2E_LANDING_PORT = Get-FreeTcpPort
    $REG_E2E_DASHBOARD_PORT = Get-FreeTcpPort
    $env:REG_E2E_PG_PORT = "$REG_E2E_PG_PORT"
    $env:REG_E2E_API_PORT = "$REG_E2E_API_PORT"
    $env:REG_E2E_LANDING_PORT = "$REG_E2E_LANDING_PORT"
    $env:REG_E2E_DASHBOARD_PORT = "$REG_E2E_DASHBOARD_PORT"
    $env:REG_E2E_API_URL = "http://localhost:$REG_E2E_API_PORT/api/v1"
    $env:REG_E2E_LANDING_URL = "http://localhost:$REG_E2E_LANDING_PORT"
    $env:REG_E2E_DASHBOARD_URL = "http://localhost:$REG_E2E_DASHBOARD_PORT"
    $env:REG_E2E_FIXTURE_FILE = $fixturePath
    $env:REG_E2E_OUTPUT_DIR = $outputPath

    $existingVolumes = @(docker volume ls --filter "label=com.docker.compose.project=$projectName" --format '{{.Name}}')
    if ($existingVolumes.Count -gt 0) { throw "Refusing to reuse existing project resources: $projectName" }

    Write-Output "Disposable project: $projectName"
    Invoke-Compose @('up', '-d', '--build', 'postgres', 'redis', 'init-roles', 'api', 'landing', 'dashboard')
    $started = $true

    $deadline = (Get-Date).AddMinutes(8)
    do {
        $apiReady = $false
        $landingReady = $false
        $dashboardReady = $false
        try { $apiReady = (Invoke-WebRequest "$env:REG_E2E_API_URL/health" -TimeoutSec 3).StatusCode -eq 200 } catch {}
        try { $landingReady = (Invoke-WebRequest "$env:REG_E2E_LANDING_URL/en/auth/register" -TimeoutSec 3).StatusCode -eq 200 } catch {}
        try { $dashboardReady = (Invoke-WebRequest "$env:REG_E2E_DASHBOARD_URL/en/auth/callback" -TimeoutSec 3).StatusCode -eq 200 } catch {}
        if (-not ($apiReady -and $landingReady -and $dashboardReady)) { Start-Sleep -Seconds 3 }
    } until (($apiReady -and $landingReady -and $dashboardReady) -or (Get-Date) -ge $deadline)
    if (-not ($apiReady -and $landingReady -and $dashboardReady)) { throw 'Disposable API and frontend services did not become healthy within eight minutes.' }

    Push-Location (Join-Path $repoRoot 'frontend')
    try {
        pnpm --filter @pos-saas/landing test:e2e:registration
        if ($LASTEXITCODE -ne 0) { throw "Playwright failed with exit code $LASTEXITCODE" }

        $WS05_TEST_PASSWORD = 'WS05-Disposable-Only-2026!'
        $authFixture = @{
            name = 'WS05 Disposable Auth Owner'
            email = 'phase2a-owner-a@example.test'
            password = $WS05_TEST_PASSWORD
            restaurant_name = 'WS05 Disposable Auth Restaurant'
            restaurant_slug = 'phase2a-restaurant-a'
        } | ConvertTo-Json
        $null = Invoke-RestMethod -Method Post -Uri "$env:REG_E2E_API_URL/auth/register" -ContentType 'application/json' -Headers @{ Origin = $env:REG_E2E_LANDING_URL } -Body $authFixture -TimeoutSec 20
        $env:WS05_LANDING_URL = $env:REG_E2E_LANDING_URL
        $env:WS05_LANDING_REDIRECT_URL = $env:REG_E2E_LANDING_URL
        $env:WS05_DASHBOARD_URL = $env:REG_E2E_DASHBOARD_URL
        $env:WS05_API_URL = $env:REG_E2E_API_URL
        $env:WS05_TEST_PASSWORD = $WS05_TEST_PASSWORD
        $env:WS05_TEST_RUNNER_ON_HOST = '1'
        pnpm --filter @pos-saas/dashboard exec playwright test --config playwright.phase2a.config.ts tests/phase2a-e2e/auth-lifecycle.spec.ts tests/phase2a-e2e/dashboard-auth-redirect.spec.ts
        if ($LASTEXITCODE -ne 0) { throw "Existing authentication Playwright tests failed with exit code $LASTEXITCODE" }
    }
    finally { Pop-Location }

    if (-not (Test-Path -LiteralPath $fixturePath)) { throw 'Registration browser test did not write its synthetic database verification fixture.' }
    $fixture = Get-Content -Raw -LiteralPath $fixturePath | ConvertFrom-Json
    $sql = @"
SELECT CASE WHEN EXISTS (
  SELECT 1
  FROM tenants t
  JOIN restaurants r ON r.tenant_id = t.id
  JOIN users u ON u.tenant_id = t.id AND u.restaurant_id = r.id
  JOIN user_roles ur ON ur.tenant_id = t.id AND ur.user_id = u.id
  JOIN roles role ON role.tenant_id = t.id AND role.restaurant_id = r.id AND role.id = ur.role_id
  JOIN themes_v2 theme ON theme.tenant_id = t.id AND theme.restaurant_id = r.id
  WHERE t.slug = :'slug'
    AND u.email = :'email'
    AND u.role = 'owner'
    AND role.role_code = 'ADMIN'
    AND jsonb_typeof(theme.config) = 'object'
) THEN 'PASS' ELSE 'FAIL' END;
"@
    $dbResult = $sql | docker compose --project-name $projectName --file $baseCompose --file $overrideCompose exec -T postgres psql -U postgres -d pos_phase2a_e2e -v ON_ERROR_STOP=1 -v "slug=$($fixture.slug)" -v "email=$($fixture.email)" -At
    if ($LASTEXITCODE -ne 0 -or $dbResult.Trim() -ne 'PASS') { throw 'Post-registration database assertion failed for tenant, restaurant, owner, role, or theme.' }
    Write-Output 'Post-registration database assertion: PASS'

    $rollbackSql = "DROP INDEX IF EXISTS idx_users_email_lower; DELETE FROM schema_migrations WHERE version = '179_index_users_email_lower.sql';"
    Invoke-Compose @('exec', '-T', 'postgres', 'psql', '-U', 'postgres', '-d', 'pos_phase2a_e2e', '-v', 'ON_ERROR_STOP=1', '-c', $rollbackSql)
    Invoke-Compose @('exec', '-T', '-e', 'DB_USER=pos_migrator', '-e', 'DB_PASSWORD=ws05-local-only-migrator', 'api', 'go', 'run', './cmd/migrate')
    $migrationCheck = docker compose --project-name $projectName --file $baseCompose --file $overrideCompose exec -T postgres psql -U postgres -d pos_phase2a_e2e -Atc "SELECT (to_regclass('public.idx_users_email_lower') IS NOT NULL)::text || ':' || EXISTS (SELECT 1 FROM schema_migrations WHERE version='179_index_users_email_lower.sql')::text"
    if ($LASTEXITCODE -ne 0 -or $migrationCheck.Trim() -ne 'true:true') { throw 'Migration 179 rollback/reapply verification failed.' }
    Write-Output 'Migration 179 rollback and reapply: PASS'
}
finally {
    if ($started) {
        Invoke-Compose @('down', '--remove-orphans')
        $projectVolumes = @(docker volume ls --filter "label=com.docker.compose.project=$projectName" --format '{{.Name}}')
        foreach ($volumeName in $projectVolumes) {
            if (-not $volumeName.StartsWith("${projectName}_", [System.StringComparison]::Ordinal)) {
                throw "Refusing to remove a volume outside the disposable project: $volumeName"
            }
            & docker volume rm $volumeName
            if ($LASTEXITCODE -ne 0) { throw "Could not remove disposable volume $volumeName" }
        }
    }
    Remove-Item -LiteralPath $fixturePath -Force -ErrorAction SilentlyContinue
}
