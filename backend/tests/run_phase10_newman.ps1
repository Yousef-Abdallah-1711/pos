param(
  [string]$BaseURL = "http://localhost:8099"
)

$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$backendDir = $repoRoot.Path
$projectRoot = Resolve-Path (Join-Path $backendDir "..")
$logPath = Join-Path $PSScriptRoot "phase10-backend.log"
$errLogPath = Join-Path $PSScriptRoot "phase10-backend.err.log"
$reportPath = Join-Path $PSScriptRoot "newman-paymob-report.xml"

if (Test-Path $logPath) { Remove-Item $logPath -Force -ErrorAction SilentlyContinue }
if (Test-Path $errLogPath) { Remove-Item $errLogPath -Force -ErrorAction SilentlyContinue }
if (Test-Path $reportPath) { Remove-Item $reportPath -Force -ErrorAction SilentlyContinue }

# Ensure Postgres is up (used by payment repository paths).
Set-Location $projectRoot.Path
docker compose up -d postgres | Out-Null

$env:PAYMOB_MODE = "mock"
$env:PORT = "8099"
$env:DB_HOST = "127.0.0.1"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "postgres"
$env:DB_NAME = "pos_saas"
$env:DB_SSLMODE = "disable"

$backendProc = Start-Process -FilePath "go" `
  -ArgumentList "run cmd/api/main.go" `
  -WorkingDirectory $backendDir `
  -PassThru `
  -RedirectStandardOutput $logPath `
  -RedirectStandardError $errLogPath

try {
  $healthy = $false
  for ($i = 0; $i -lt 120; $i++) {
    Start-Sleep -Seconds 1
    try {
      $healthResp = Invoke-WebRequest -UseBasicParsing -Uri "$BaseURL/api/v1/health" -TimeoutSec 2
      if ([int]$healthResp.StatusCode -eq 200) {
        $healthy = $true
        break
      }
    } catch {}
  }

  if (-not $healthy) {
    Write-Host "Backend failed health check at $BaseURL"
    if (Test-Path $logPath) { Get-Content $logPath -Tail 200 }
    if (Test-Path $errLogPath) { Get-Content $errLogPath -Tail 200 }
    exit 1
  }

  Set-Location $projectRoot.Path
  npx -y newman run "$backendDir/tests/paymob_api_tests.json" `
    -e "$backendDir/tests/env.json" `
    --env-var "base_url=$BaseURL" `
    --reporters "cli,junit" `
    --reporter-junit-export "$reportPath"

  $newmanExit = $LASTEXITCODE
  if ($newmanExit -ne 0) {
    Write-Host "Newman run failed with exit code $newmanExit"
    if (Test-Path $logPath) { Get-Content $logPath -Tail 200 }
    if (Test-Path $errLogPath) { Get-Content $errLogPath -Tail 200 }
    exit $newmanExit
  }

  Write-Host "Newman run passed. Report: $reportPath"
}
finally {
  if ($backendProc -and -not $backendProc.HasExited) {
    Stop-Process -Id $backendProc.Id -Force -ErrorAction SilentlyContinue
  }
}
