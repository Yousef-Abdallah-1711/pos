<#
LOCAL/TEST-ONLY SCRIPT

Purpose:
- Runs the Paymob Newman collection against a local backend by default.
- Can target staging/other environments by overriding BASE_URL.

Examples:
- Local default:
    powershell -File .\backend\tests\run_phase10_newman.ps1
- Override target with parameter:
    powershell -File .\backend\tests\run_phase10_newman.ps1 -BaseURL "https://api.staging.example.com"
- Override target with environment variable:
    $env:BASE_URL="https://api.staging.example.com"; powershell -File .\backend\tests\run_phase10_newman.ps1
- Use existing backend (skip local docker/go startup):
    powershell -File .\backend\tests\run_phase10_newman.ps1 -BaseURL "https://api.staging.example.com" -UseExistingBackend
#>
param(
  [string]$BaseURL,
  [switch]$UseExistingBackend,
  [string]$PostmanEnvFile
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($BaseURL)) {
  $BaseURL = $env:BASE_URL
}
if ([string]::IsNullOrWhiteSpace($BaseURL)) {
  $BaseURL = "http://localhost:8099"
}
$BaseURL = $BaseURL.TrimEnd("/")

if (-not $PSBoundParameters.ContainsKey("UseExistingBackend") -and $env:USE_EXISTING_BACKEND -eq "true") {
  $UseExistingBackend = $true
}

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$backendDir = $repoRoot.Path
$projectRoot = Resolve-Path (Join-Path $backendDir "..")
$logPath = Join-Path $PSScriptRoot "phase10-backend.log"
$errLogPath = Join-Path $PSScriptRoot "phase10-backend.err.log"
$reportPath = Join-Path $PSScriptRoot "newman-paymob-report.xml"
$defaultEnvFile = Join-Path $backendDir "tests/env.json"

if ([string]::IsNullOrWhiteSpace($PostmanEnvFile)) {
  if ($env:POSTMAN_ENV_FILE) {
    $PostmanEnvFile = $env:POSTMAN_ENV_FILE
  } elseif (Test-Path $defaultEnvFile) {
    $PostmanEnvFile = $defaultEnvFile
  }
}

if (Test-Path $logPath) { Remove-Item $logPath -Force -ErrorAction SilentlyContinue }
if (Test-Path $errLogPath) { Remove-Item $errLogPath -Force -ErrorAction SilentlyContinue }
if (Test-Path $reportPath) { Remove-Item $reportPath -Force -ErrorAction SilentlyContinue }

$backendProc = $null

if (-not $UseExistingBackend) {
  # Local default mode: boot local dependencies/backend for the Newman run.
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
} else {
  Write-Host "Using existing backend at $BaseURL (skipping local docker/go startup)"
}

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
  $newmanArgs = @(
    "-y",
    "newman",
    "run",
    "$backendDir/tests/paymob_api_tests.json",
    "--env-var", "base_url=$BaseURL",
    "--reporters", "cli,junit",
    "--reporter-junit-export", "$reportPath"
  )

  if (-not [string]::IsNullOrWhiteSpace($PostmanEnvFile)) {
    if (Test-Path $PostmanEnvFile) {
      $resolvedPostmanEnvFile = (Resolve-Path $PostmanEnvFile).Path
      $newmanArgs += @("-e", $resolvedPostmanEnvFile)
    } else {
      Write-Host "Postman environment file not found: $PostmanEnvFile (continuing with CLI env vars only)"
    }
  } else {
    Write-Host "No Postman environment file provided/found; running with CLI env vars only."
  }

  & npx @newmanArgs

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
