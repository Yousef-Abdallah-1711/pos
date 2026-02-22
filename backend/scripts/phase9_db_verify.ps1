param(
  [switch]$UseDocker = $true,
  [switch]$StartBackend = $true,
  [int]$TenantID = 1,
  [int]$RestaurantID = 1,
  [string]$MockBaseURL = 'http://localhost:8080',
  [string]$LiveBaseURL = 'http://localhost:8081',
  [string]$OutputFile = 'phase9_db_results.json'
)

$ErrorActionPreference = 'Stop'

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Resolve-Path (Join-Path $scriptDir '..\..')
$backendDir = Join-Path $repoRoot 'backend'
$seedSQLPath = Join-Path $scriptDir 'paymob_phase9_seed.sql'
$outputPath = Join-Path $scriptDir $OutputFile
$mockLogPath = Join-Path $scriptDir 'phase9_mock_backend.log'
$liveLogPath = Join-Path $scriptDir 'phase9_live_backend.log'

function Stop-PortListeners {
  param([int[]]$Ports)
  foreach ($port in $Ports) {
    $listeners = Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue
    if (-not $listeners) { continue }
    $listeners | Select-Object -ExpandProperty OwningProcess -Unique | ForEach-Object {
      try { Stop-Process -Id $_ -Force -ErrorAction SilentlyContinue } catch {}
    }
  }
}

function Wait-Endpoint200 {
  param([string]$URL, [int]$TimeoutSec = 90)
  for ($i = 0; $i -lt $TimeoutSec; $i++) {
    Start-Sleep -Seconds 1
    try {
      $resp = Invoke-WebRequest -UseBasicParsing -Uri $URL -TimeoutSec 2
      if ([int]$resp.StatusCode -eq 200) { return $true }
    } catch {}
  }
  return $false
}

function Wait-TcpPort {
  param(
    [string]$Hostname = '127.0.0.1',
    [int]$Port,
    [int]$TimeoutSec = 60
  )
  for ($i = 0; $i -lt $TimeoutSec; $i++) {
    try {
      $client = New-Object System.Net.Sockets.TcpClient
      $iar = $client.BeginConnect($Hostname, $Port, $null, $null)
      $ok = $iar.AsyncWaitHandle.WaitOne(1000, $false)
      if ($ok -and $client.Connected) {
        $client.EndConnect($iar)
        $client.Close()
        return $true
      }
      $client.Close()
    } catch {}
    Start-Sleep -Milliseconds 500
  }
  return $false
}

function Invoke-API {
  param(
    [string]$Method,
    [string]$URL,
    [hashtable]$Headers = @{},
    [string]$Body = ''
  )

  try {
    if ([string]::IsNullOrWhiteSpace($Body)) {
      $resp = Invoke-WebRequest -UseBasicParsing -Method $Method -Uri $URL -Headers $Headers -TimeoutSec 20
    } else {
      $resp = Invoke-WebRequest -UseBasicParsing -Method $Method -Uri $URL -Headers $Headers -ContentType 'application/json' -Body $Body -TimeoutSec 20
    }
    return @{
      Status = [int]$resp.StatusCode
      Body   = $resp.Content
    }
  } catch {
    if ($_.Exception.Response) {
      $response = $_.Exception.Response
      $body = ''
      try {
        $reader = New-Object System.IO.StreamReader($response.GetResponseStream())
        $body = $reader.ReadToEnd()
        $reader.Close()
      } catch {}
      return @{
        Status = [int]$response.StatusCode
        Body   = $body
      }
    }
    return @{
      Status = 0
      Body   = $_.Exception.Message
    }
  }
}

function Convert-BodyToObject {
  param([string]$Body)
  if ([string]::IsNullOrWhiteSpace($Body)) { return $null }
  try { return ($Body | ConvertFrom-Json) } catch { return $null }
}

function New-MockWebhookPayload {
  param(
    [int64]$OrderID,
    [int]$AmountCents,
    [string]$Currency = 'EGP',
    [bool]$Success = $true
  )

  return @{
    type = 'TRANSACTION'
    hmac = 'mock_hmac'
    obj  = @{
      id                    = [int64](Get-Date -UFormat %s) * 1000000 + (Get-Random -Minimum 1000 -Maximum 9999)
      amount_cents          = $AmountCents
      success               = $Success
      currency              = $Currency
      is_auth               = $Success
      is_capture            = $Success
      is_refunded           = $false
      is_voided             = $false
      is_3d_secure          = $false
      is_standalone_payment = $true
      has_parent_transaction = $false
      error_occured         = (-not $Success)
      pending               = $false
      integration_id        = 0
      owner                 = 0
      created_at            = (Get-Date).ToString('o')
      order                 = @{
        id = $OrderID
        merchant_order_id = "tenant_${TenantID}_monthly_${OrderID}"
      }
      source_data = @{
        pan      = '5123456789011111'
        sub_type = 'mastercard'
        type     = 'card'
      }
    }
  } | ConvertTo-Json -Depth 8
}

$mockBackendJob = $null
$liveBackendJob = $null
$summary = [ordered]@{
  timestamp = (Get-Date).ToString('s')
  config = [ordered]@{
    use_docker = [bool]$UseDocker
    start_backend = [bool]$StartBackend
    tenant_id = $TenantID
    restaurant_id = $RestaurantID
    mock_base_url = $MockBaseURL
    live_base_url = $LiveBaseURL
  }
  checks = [ordered]@{}
  details = [ordered]@{}
  notes = @()
}

try {
  if ($UseDocker) {
    docker compose up -d postgres | Out-Null

    $healthy = $false
    for ($i = 0; $i -lt 120; $i++) {
      Start-Sleep -Seconds 1
      try {
        $status = docker inspect -f '{{.State.Health.Status}}' pos-postgres 2>$null
        if ($status -eq 'healthy') {
          $healthy = $true
          break
        }
      } catch {}
    }
    if (-not $healthy) {
      throw 'postgres container did not become healthy'
    }

    if (-not (Wait-TcpPort -Hostname '127.0.0.1' -Port 5432 -TimeoutSec 90)) {
      throw 'postgres port 5432 is not reachable from host'
    }

    if (Test-Path $seedSQLPath) {
      Get-Content -Path $seedSQLPath -Raw | docker exec -i pos-postgres psql -U postgres -d pos_saas | Out-Null
      $summary.notes += 'seed sql applied via docker exec'
    }
  } else {
    $summary.notes += 'docker setup skipped; ensure PostgreSQL + seeded data are ready manually'
  }

  if ($StartBackend) {
    Stop-PortListeners -Ports @(8080, 8081)
    if (Test-Path $mockLogPath) { Remove-Item -Force $mockLogPath -ErrorAction SilentlyContinue }
    $mockBackendJob = Start-Job -ScriptBlock {
      Set-Location $using:backendDir
      $env:PAYMOB_MODE = 'mock'
      $env:PORT = '8080'
      $env:DB_HOST = '127.0.0.1'
      $env:DB_PORT = '5432'
      $env:DB_USER = 'postgres'
      $env:DB_PASSWORD = 'postgres'
      $env:DB_NAME = 'pos_saas'
      $env:DB_SSLMODE = 'disable'
      go run cmd/api/main.go 2>&1 | Tee-Object -FilePath $using:mockLogPath
    }
  }

  if (-not (Wait-Endpoint200 -URL "$MockBaseURL/api/v1/health" -TimeoutSec 120)) {
    throw 'mock backend health endpoint is not reachable'
  }

  $authHeaders = @{
    'X-Impersonation-Token' = 'phase9-verify'
    'X-Tenant-ID' = "$TenantID"
    'X-Restaurant-ID' = "$RestaurantID"
  }

  $plansResp = Invoke-API -Method 'GET' -URL "$MockBaseURL/api/v1/subscription/plans"
  $summary.checks['9.2.plans'] = $plansResp.Status
  if ($plansResp.Status -ne 200) { $summary.details['9.2.plans.body'] = $plansResp.Body }

  $init1 = Invoke-API -Method 'POST' -URL "$MockBaseURL/api/v1/subscription/initiate" -Headers $authHeaders -Body '{"plan_id":"monthly"}'
  $summary.checks['9.2.initiate_auth'] = $init1.Status
  if ($init1.Status -ne 200) { $summary.details['9.2.initiate_auth.body'] = $init1.Body }
  $init1Obj = Convert-BodyToObject -Body $init1.Body

  $paymentsBefore = Invoke-API -Method 'GET' -URL "$MockBaseURL/api/v1/subscription/payments?page=1&page_size=20" -Headers $authHeaders
  $summary.checks['9.2.payments_auth_before'] = $paymentsBefore.Status
  if ($paymentsBefore.Status -ne 200) { $summary.details['9.2.payments_auth_before.body'] = $paymentsBefore.Body }

  $current = Invoke-API -Method 'GET' -URL "$MockBaseURL/api/v1/subscription/current" -Headers $authHeaders
  $summary.checks['9.2.current_auth'] = $current.Status
  if ($current.Status -ne 200) { $summary.details['9.2.current_auth.body'] = $current.Body }

  $missingJWT = Invoke-API -Method 'GET' -URL "$MockBaseURL/api/v1/subscription/current"
  $summary.checks['9.3.missing_auth'] = $missingJWT.Status
  if ($missingJWT.Status -ne 401) { $summary.details['9.3.missing_auth.body'] = $missingJWT.Body }

  $order1 = 0
  if ($init1Obj -and $init1Obj.data -and $init1Obj.data.paymob_order_id) {
    $order1 = [int64]$init1Obj.data.paymob_order_id
  }
  if ($order1 -gt 0) {
    $webhookSuccessPayload = New-MockWebhookPayload -OrderID $order1 -AmountCents 50000 -Currency 'EGP' -Success $true
    $webhook1 = Invoke-API -Method 'POST' -URL "$MockBaseURL/api/v1/webhook/paymob?hmac=mock_hmac" -Body $webhookSuccessPayload
    $webhook2 = Invoke-API -Method 'POST' -URL "$MockBaseURL/api/v1/webhook/paymob?hmac=mock_hmac" -Body $webhookSuccessPayload
    $summary.checks['9.2.webhook_success_first'] = $webhook1.Status
    $summary.checks['9.3.duplicate_webhook_second'] = $webhook2.Status
    if ($webhook1.Status -ne 200) { $summary.details['9.2.webhook_success_first.body'] = $webhook1.Body }
    if ($webhook2.Status -ne 200) { $summary.details['9.3.duplicate_webhook_second.body'] = $webhook2.Body }
  } else {
    $summary.notes += 'could not extract paymob_order_id from initiate response'
  }

  $paymentsAfter = Invoke-API -Method 'GET' -URL "$MockBaseURL/api/v1/subscription/payments?page=1&page_size=20" -Headers $authHeaders
  $summary.checks['9.2.payments_auth_after'] = $paymentsAfter.Status
  if ($paymentsAfter.Status -ne 200) { $summary.details['9.2.payments_auth_after.body'] = $paymentsAfter.Body }

  $init2 = Invoke-API -Method 'POST' -URL "$MockBaseURL/api/v1/subscription/initiate" -Headers $authHeaders -Body '{"plan_id":"monthly"}'
  $summary.checks['9.3.initiate_for_amount_test'] = $init2.Status
  if ($init2.Status -ne 200) { $summary.details['9.3.initiate_for_amount_test.body'] = $init2.Body }
  $init2Obj = Convert-BodyToObject -Body $init2.Body
  $order2 = 0
  if ($init2Obj -and $init2Obj.data -and $init2Obj.data.paymob_order_id) {
    $order2 = [int64]$init2Obj.data.paymob_order_id
  }
  if ($order2 -gt 0) {
    $amountMismatchPayload = New-MockWebhookPayload -OrderID $order2 -AmountCents 100 -Currency 'EGP' -Success $true
    $amountMismatch = Invoke-API -Method 'POST' -URL "$MockBaseURL/api/v1/webhook/paymob?hmac=mock_hmac" -Body $amountMismatchPayload
    $summary.checks['9.3.amount_mismatch'] = $amountMismatch.Status
    if ($amountMismatch.Status -eq 200) { $summary.details['9.3.amount_mismatch.body'] = $amountMismatch.Body }
  } else {
    $summary.notes += 'could not run amount mismatch check because second initiate failed'
  }

  $sqlInjection = Invoke-API -Method 'POST' -URL "$MockBaseURL/api/v1/subscription/initiate" -Headers $authHeaders -Body '{"plan_id":"monthly''; DROP TABLE payments; --"}'
  $summary.checks['9.3.sql_injection_input'] = $sqlInjection.Status
  if ($sqlInjection.Status -eq 200) { $summary.details['9.3.sql_injection_input.body'] = $sqlInjection.Body }

  if ($UseDocker -and $order1 -gt 0) {
    $countQuery = "SELECT COUNT(*) FROM payments WHERE paymob_order_id = $order1;"
    try {
      $countValue = (docker exec -i pos-postgres psql -U postgres -d pos_saas -t -A -c $countQuery).Trim()
      $summary.checks['9.3.duplicate_db_count'] = [int]$countValue
    } catch {
      $summary.notes += 'failed to query duplicate payment count from DB'
    }
  }

  # HMAC tamper check in live mode
  if ($StartBackend) {
    Stop-PortListeners -Ports @(8081)
    if (Test-Path $liveLogPath) { Remove-Item -Force $liveLogPath -ErrorAction SilentlyContinue }
    $liveBackendJob = Start-Job -ScriptBlock {
      Set-Location $using:backendDir
      $env:PAYMOB_MODE = 'live'
      $env:PAYMOB_API_KEY = 'dummy'
      $env:PAYMOB_IFRAME_ID = '1'
      $env:PAYMOB_HMAC_SECRET = 'test_secret'
      $env:PAYMOB_BASE_URL = 'https://accept.paymob.com/api'
      $env:PORT = '8081'
      $env:DB_HOST = '127.0.0.1'
      $env:DB_PORT = '5432'
      $env:DB_USER = 'postgres'
      $env:DB_PASSWORD = 'postgres'
      $env:DB_NAME = 'pos_saas'
      $env:DB_SSLMODE = 'disable'
      go run cmd/api/main.go 2>&1 | Tee-Object -FilePath $using:liveLogPath
    }
  }

  if (Wait-Endpoint200 -URL "$LiveBaseURL/api/v1/health" -TimeoutSec 120) {
    $tampered = Invoke-API -Method 'POST' -URL "$LiveBaseURL/api/v1/webhook/paymob?hmac=deadbeef" -Body (New-MockWebhookPayload -OrderID 999999 -AmountCents 50000 -Currency 'EGP' -Success $true)
    $summary.checks['9.3.hmac_tampered_live'] = $tampered.Status
    if ($tampered.Status -ne 401) { $summary.details['9.3.hmac_tampered_live.body'] = $tampered.Body }
  } else {
    $summary.notes += 'live backend not started; skipped hmac tamper runtime check'
  }
}
finally {
  if ($liveBackendJob) {
    Stop-Job $liveBackendJob -ErrorAction SilentlyContinue
    Remove-Job $liveBackendJob -ErrorAction SilentlyContinue
  }
  if ($mockBackendJob) {
    Stop-Job $mockBackendJob -ErrorAction SilentlyContinue
    Remove-Job $mockBackendJob -ErrorAction SilentlyContinue
  }
  Stop-PortListeners -Ports @(8080, 8081)
}

$summary | ConvertTo-Json -Depth 8 | Set-Content -Path $outputPath
Write-Output ("Saved: " + $outputPath)
Get-Content -Path $outputPath
