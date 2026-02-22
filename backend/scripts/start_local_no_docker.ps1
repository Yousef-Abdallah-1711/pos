param(
  [string]$Port = "8080",
  [switch]$SkipMigrate,
  [switch]$SkipSeed
)

$ErrorActionPreference = "Stop"

$backendRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
Push-Location $backendRoot

$embeddedProcess = $null
$embeddedStdout = Join-Path $backendRoot "scripts/embedded_pg_daemon.out.log"
$embeddedStderr = Join-Path $backendRoot "scripts/embedded_pg_daemon.err.log"

try {
  Write-Host "[local-backend] Starting embedded PostgreSQL..." -ForegroundColor Cyan

  if (Test-Path $embeddedStdout) { Remove-Item $embeddedStdout -Force }
  if (Test-Path $embeddedStderr) { Remove-Item $embeddedStderr -Force }

  $embeddedProcess = Start-Process `
    -FilePath "go" `
    -ArgumentList "run", "./scripts/embedded_pg_daemon.go" `
    -PassThru `
    -RedirectStandardOutput $embeddedStdout `
    -RedirectStandardError $embeddedStderr

  Start-Sleep -Seconds 4

  if (-not (Get-Process -Id $embeddedProcess.Id -ErrorAction SilentlyContinue)) {
    Write-Error "Embedded PostgreSQL failed to start. Check $embeddedStderr"
  }

  $env:DB_HOST = "127.0.0.1"
  $env:DB_PORT = "5432"
  $env:DB_USER = "postgres"
  $env:DB_PASSWORD = "postgres"
  $env:DB_NAME = "pos_saas"
  $env:DB_SSLMODE = "disable"
  $env:PAYMOB_MODE = "mock"
  $env:PORT = $Port

  if (-not $SkipMigrate) {
    Write-Host "[local-backend] Running migrations..." -ForegroundColor Cyan
    go run ./cmd/migrate/main.go
  }

  if ((-not $SkipSeed) -and (Test-Path "./scripts/paymob_phase9_seed.sql")) {
    Write-Host "[local-backend] Seeding payment verification data..." -ForegroundColor Cyan
    go run ./scripts/sql_exec.go -file ./scripts/paymob_phase9_seed.sql
  }

  Write-Host "[local-backend] Starting API on http://localhost:$Port ..." -ForegroundColor Green
  go run ./cmd/api/main.go
}
finally {
  if ($embeddedProcess -and (Get-Process -Id $embeddedProcess.Id -ErrorAction SilentlyContinue)) {
    Write-Host "[local-backend] Stopping embedded PostgreSQL..." -ForegroundColor Yellow
    Stop-Process -Id $embeddedProcess.Id -Force -ErrorAction SilentlyContinue
  }

  Pop-Location
}
