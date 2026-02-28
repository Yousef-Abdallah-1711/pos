param(
  [int]$Port = 8080,
  [switch]$Rebuild
)

$ErrorActionPreference = 'Stop'

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot '..')
Set-Location $repoRoot

Write-Host "[backend] Using repo: $repoRoot"

$listeners = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
if ($listeners) {
  $ownerPids = $listeners | Select-Object -ExpandProperty OwningProcess -Unique
  foreach ($ownerPid in $ownerPids) {
    try {
      $proc = Get-Process -Id $ownerPid -ErrorAction Stop
      Write-Host "[backend] Stopping process on port ${Port}: PID=${ownerPid} Name=$($proc.ProcessName)"
      Stop-Process -Id $ownerPid -Force
    } catch {
      Write-Warning "[backend] Failed to stop PID ${ownerPid}: $($_.Exception.Message)"
    }
  }
}

$exePath = Join-Path $repoRoot "api-local.exe"
if ($Rebuild -or -not (Test-Path $exePath)) {
  Write-Host "[backend] Building api-local.exe from this repo..."
  go build -o $exePath ./cmd/api
}

$outLog = Join-Path $repoRoot "api-local.out.log"
$errLog = Join-Path $repoRoot "api-local.err.log"

Write-Host "[backend] Starting api-local.exe on port $Port"
$proc = Start-Process -FilePath $exePath `
  -WorkingDirectory $repoRoot `
  -RedirectStandardOutput $outLog `
  -RedirectStandardError $errLog `
  -PassThru

Write-Host "[backend] Started PID=$($proc.Id). Waiting for health..."

$ready = $false
for ($attempt = 1; $attempt -le 12; $attempt++) {
  Start-Sleep -Seconds 1
  try {
    $response = Invoke-RestMethod -Uri "http://localhost:$Port/api/v1/health" -TimeoutSec 2
    Write-Host "[backend] Health check passed on attempt ${attempt}: status=$($response.status) db_healthy=$($response.db_healthy)"
    $ready = $true
    break
  } catch {
    Write-Host "[backend] Waiting... attempt $attempt failed"
  }
}

if (-not $ready) {
  Write-Warning "[backend] Health endpoint did not pass. Check logs:"
  Write-Warning "  $outLog"
  Write-Warning "  $errLog"
  exit 1
}

Write-Host "[backend] Backend is running from this repo."
