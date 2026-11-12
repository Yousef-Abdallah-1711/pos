param(
    [string]$ApiBase = 'http://127.0.0.1:18081',
    [string]$DatabaseContainer = 'phase00-closure-20261111-postgres-1',
    [string]$Origin = 'http://localhost:3002'
)

$ErrorActionPreference = 'Stop'
$runId = [Guid]::NewGuid().ToString('N').Substring(0, 8)
$date = '2026-11-11'
$results = [System.Collections.Generic.List[object]]::new()

function Invoke-ClosureRequest {
    param([string]$Method, [string]$Path, [object]$Body, [object]$Session, [hashtable]$ExtraHeaders = @{})
    $headers = @{ Origin = $Origin }
    foreach ($key in $ExtraHeaders.Keys) { $headers[$key] = $ExtraHeaders[$key] }
    $request = @{ Uri = "$ApiBase$Path"; Method = $Method; Headers = $headers; SkipHttpErrorCheck = $true }
    if ($null -ne $Session) { $request.WebSession = $Session }
    if ($null -ne $Body) { $request.ContentType = 'application/json'; $request.Body = $Body | ConvertTo-Json -Compress }
    $response = Invoke-WebRequest @request
    $parsed = $null
    try { $parsed = $response.Content | ConvertFrom-Json -AsHashtable } catch { }
    return [pscustomobject]@{ Status = [int]$response.StatusCode; Content = [string]$response.Content; Json = $parsed; Headers = $response.Headers }
}

function Add-Result([string]$Name, [int]$Status, [string]$Expected, [bool]$Passed, [hashtable]$Details = @{}) {
    $results.Add([pscustomobject]@{ name = $Name; status = $Status; expected = $Expected; result = $(if ($Passed) { 'PASS' } else { 'FAIL' }); details = $Details })
}

function New-RegisteredOwner([string]$Suffix, [string]$Name, [string]$Slug) {
    $email = "$Name-$Suffix@example.test"
    $restaurantName = "$Name-$Suffix"
    $password = [Guid]::NewGuid().ToString('N') + 'Aa!9'
    $session = New-Object Microsoft.PowerShell.Commands.WebRequestSession
    $registration = Invoke-ClosureRequest 'POST' '/api/v1/auth/register' @{ name = $restaurantName; email = $email; password = $password; restaurant_name = $restaurantName; restaurant_slug = $Slug; phone = '+201000000001' } $session
    if ($registration.Status -ne 201) { throw "Synthetic registration failed for $Name (HTTP $($registration.Status)); see local API logs." }
    return [pscustomobject]@{ Email = $email; Password = $password; Session = $session; User = $registration.Json.user }
}

function Invoke-RealLogin($Account) {
    $response = Invoke-ClosureRequest 'POST' '/api/v1/auth/login' @{ email = $Account.Email; password = $Account.Password } $Account.Session
    if ($response.Status -ne 200) { throw "Synthetic login failed (HTTP $($response.Status)); see local API logs." }
    $sessionResponse = Invoke-ClosureRequest 'GET' '/api/v1/auth/session' $null $Account.Session
    return [pscustomobject]@{ Account = $Account; LoginStatus = $response.Status; SessionStatus = $sessionResponse.Status; SessionUser = $sessionResponse.Json }
}

function Invoke-DisposableSql([string]$Sql) {
    $output = $Sql | docker exec -i $DatabaseContainer psql -U closure_test -d pos_final_integration_test -v ON_ERROR_STOP=1 -At
    if ($LASTEXITCODE -ne 0) { throw 'Disposable database fixture/assertion failed.' }
    return (($output | Out-String).Trim())
}

$health = Invoke-ClosureRequest 'GET' '/api/v1/health/ready' $null $null
if ($health.Status -ne 200) { throw "Isolated API health check failed (HTTP $($health.Status))." }

$ownerAAccount = New-RegisteredOwner $runId 'phase00-owner-a' "p00a-$runId"
$ownerBAccount = New-RegisteredOwner $runId 'phase00-owner-b' "p00b-$runId"
$cashierAccount = New-RegisteredOwner $runId 'phase00-cashier' "p00c-$runId"
$ownerA = Invoke-RealLogin $ownerAAccount
$ownerB = Invoke-RealLogin $ownerBAccount

$ids = @($ownerA.SessionUser, $ownerB.SessionUser, $cashierAccount.User)
foreach ($id in $ids) {
    if (-not $id.tenant_id -or -not $id.restaurant_id -or -not $id.id) { throw 'Registration response omitted a test identity scope.' }
}

$cashierId = [int]$cashierAccount.User.id
$cashierTenantId = [int]$cashierAccount.User.tenant_id
$roleUpdate = Invoke-DisposableSql "UPDATE users SET role='cashier' WHERE id=$cashierId AND tenant_id=$cashierTenantId; DELETE FROM user_roles WHERE user_id=$cashierId AND tenant_id=$cashierTenantId; SELECT role FROM users WHERE id=$cashierId AND tenant_id=$cashierTenantId;"
if (-not $roleUpdate.Trim().EndsWith('cashier')) { throw 'Could not establish a no-permissions cashier fixture in the disposable database.' }
$cashier = Invoke-RealLogin $cashierAccount

$tenantA = [int]$ownerA.SessionUser.tenant_id
$restaurantA = [int]$ownerA.SessionUser.restaurant_id
$tenantB = [int]$ownerB.SessionUser.tenant_id
$restaurantB = [int]$ownerB.SessionUser.restaurant_id
$cashierTenant = [int]$cashier.SessionUser.tenant_id
$cashierRestaurant = [int]$cashier.SessionUser.restaurant_id
$cash = "P00A-$runId-CASH"
$capital = "P00A-$runId-CAP"
$cashB = "P00B-$runId-CASH"
$capitalB = "P00B-$runId-CAP"
$entryA = "P00-A-$runId"
$entryB = "P00-B-$runId"
$orderA = "P00-A-$runId"
$orderB = "P00-B-$runId"

$fixtureSql = @"
BEGIN;
INSERT INTO chart_of_accounts (tenant_id,restaurant_id,account_code,account_name,account_type,normal_balance,level,is_header,current_balance,currency,is_active,is_system)
VALUES ($tenantA,$restaurantA,'$cash','Synthetic Tenant A cash','asset','debit',3,FALSE,0,'EGP',TRUE,FALSE),
       ($tenantA,$restaurantA,'$capital','Synthetic Tenant A capital','equity','credit',2,FALSE,0,'EGP',TRUE,FALSE),
       ($tenantB,$restaurantB,'$cashB','Synthetic Tenant B cash','asset','debit',3,FALSE,0,'EGP',TRUE,FALSE),
       ($tenantB,$restaurantB,'$capitalB','Synthetic Tenant B capital','equity','credit',2,FALSE,0,'EGP',TRUE,FALSE);
WITH fixtures(tenant_id,restaurant_id,entry_number,amount,asset_code,equity_code,marker) AS (
 VALUES ($tenantA,$restaurantA,'$entryA',110::numeric,'$cash','$capital','A'),
        ($tenantB,$restaurantB,'$entryB',220::numeric,'$cashB','$capitalB','B')
), entries AS (
 INSERT INTO journal_entries (tenant_id,restaurant_id,entry_number,entry_date,description,entry_type,status,posted_at)
 SELECT tenant_id,restaurant_id,entry_number,DATE '$date','Synthetic Phase 00 isolation fixture '||marker,'manual','posted',NOW() FROM fixtures
 RETURNING id,tenant_id,restaurant_id,entry_date,entry_number
), accounts AS (
 SELECT e.id AS journal_entry_id,e.tenant_id,e.restaurant_id,e.entry_date,
        coa.id AS account_id,coa.account_code,coa.account_type,f.amount
 FROM entries e JOIN fixtures f USING (tenant_id,restaurant_id,entry_number)
 JOIN chart_of_accounts coa ON coa.tenant_id=e.tenant_id AND coa.restaurant_id=e.restaurant_id
   AND coa.account_code IN (f.asset_code,f.equity_code)
), lines AS (
 INSERT INTO journal_entry_lines (journal_entry_id,tenant_id,account_id,account_code,debit_amount,credit_amount,description,line_order)
 SELECT journal_entry_id,tenant_id,account_id,account_code,
        CASE WHEN account_type='asset' THEN amount ELSE 0 END,
        CASE WHEN account_type='equity' THEN amount ELSE 0 END,
        'Synthetic tenant isolation fixture',ROW_NUMBER() OVER(PARTITION BY journal_entry_id ORDER BY account_code)
 FROM accounts RETURNING id,journal_entry_id,tenant_id,account_id,debit_amount,credit_amount
), ledger AS (
 INSERT INTO general_ledger (tenant_id,restaurant_id,account_id,journal_entry_id,journal_line_id,entry_date,fiscal_period,debit_amount,credit_amount,running_balance)
 SELECT l.tenant_id,a.restaurant_id,l.account_id,l.journal_entry_id,l.id,a.entry_date,'2026-11',l.debit_amount,l.credit_amount,l.debit_amount-l.credit_amount
 FROM lines l JOIN accounts a ON a.journal_entry_id=l.journal_entry_id AND a.account_id=l.account_id
 RETURNING id
)
SELECT 'ledger_fixture_rows='||COUNT(*) FROM ledger;
INSERT INTO orders (tenant_id,restaurant_id,order_number,customer_name,customer_phone,subtotal,tax_amount,discount_amount,delivery_fee,total_amount,payment_method,payment_status,status,order_source,created_at)
VALUES ($tenantA,$restaurantA,'$orderA','Synthetic A','+201000000001',123.45,0,0,0,123.45,'cash','paid','confirmed','in_store',TIMESTAMP '$date 12:00:00'),
       ($tenantB,$restaurantB,'$orderB','Synthetic B','+201000000002',987.65,0,0,0,987.65,'cash','paid','confirmed','in_store',TIMESTAMP '$date 12:01:00');
COMMIT;
"@
$fixtureOutput = Invoke-DisposableSql $fixtureSql

if ($ownerA.LoginStatus -eq 200 -and $ownerA.SessionStatus -eq 200 -and $ownerA.SessionUser.role -eq 'owner') { Add-Result 'Owner A actual login and session' 200 'owner session established' $true @{ tenant_id = $tenantA; restaurant_id = $restaurantA; role = 'owner' } } else { Add-Result 'Owner A actual login and session' $ownerA.LoginStatus 'owner session established' $false @{} }
if ($ownerB.LoginStatus -eq 200 -and $ownerB.SessionStatus -eq 200 -and $ownerB.SessionUser.role -eq 'owner') { Add-Result 'Owner B actual login and session' 200 'owner session established' $true @{ tenant_id = $tenantB; restaurant_id = $restaurantB; role = 'owner' } } else { Add-Result 'Owner B actual login and session' $ownerB.LoginStatus 'owner session established' $false @{} }

$path = "/api/v1/reports/financial/balance-sheet?as_of=$date"
$reportA = Invoke-ClosureRequest 'GET' $path $null $ownerA.Account.Session
$reportB = Invoke-ClosureRequest 'GET' $path $null $ownerB.Account.Session
$passA = $reportA.Status -eq 200 -and $reportA.Json.data.total_assets -eq 110 -and $reportA.Json.data.total_equity -eq 110
$passB = $reportB.Status -eq 200 -and $reportB.Json.data.total_assets -eq 220 -and $reportB.Json.data.total_equity -eq 220
Add-Result 'Owner A reads own report' $reportA.Status '200; assets/equity=110; no B account marker' ($passA -and -not $reportA.Content.Contains($cashB)) @{ total_assets = $reportA.Json.data.total_assets; total_equity = $reportA.Json.data.total_equity; tenant_id = $tenantA }
Add-Result 'Owner B reads own report' $reportB.Status '200; assets/equity=220; no A account marker' ($passB -and -not $reportB.Content.Contains($cash)) @{ total_assets = $reportB.Json.data.total_assets; total_equity = $reportB.Json.data.total_equity; tenant_id = $tenantB }

$forgedA = Invoke-ClosureRequest 'GET' "$path&tenant_id=$tenantB&restaurant_id=$restaurantB" $null $ownerA.Account.Session @{ 'X-Tenant-ID' = "$tenantB"; 'X-Restaurant-ID' = "$restaurantB" }
$forgedB = Invoke-ClosureRequest 'GET' "$path&tenant_id=$tenantA&restaurant_id=$restaurantA" $null $ownerB.Account.Session @{ 'X-Tenant-ID' = "$tenantA"; 'X-Restaurant-ID' = "$restaurantA" }
Add-Result 'Owner A cross-tenant/report/restaurant forgery' $forgedA.Status '200 scoped to A only; no B amounts or account marker' ($forgedA.Status -eq 200 -and $forgedA.Json.data.total_assets -eq 110 -and $forgedA.Json.data.total_equity -eq 110 -and -not $forgedA.Content.Contains($cashB) -and -not $forgedA.Content.Contains('220')) @{ attempted_tenant_id = $tenantB; attempted_restaurant_id = $restaurantB; returned_total_assets = $forgedA.Json.data.total_assets }
Add-Result 'Owner B cross-tenant/report/restaurant forgery' $forgedB.Status '200 scoped to B only; no A amounts or account marker' ($forgedB.Status -eq 200 -and $forgedB.Json.data.total_assets -eq 220 -and $forgedB.Json.data.total_equity -eq 220 -and -not $forgedB.Content.Contains($cash) -and -not $forgedB.Content.Contains('110')) @{ attempted_tenant_id = $tenantA; attempted_restaurant_id = $restaurantA; returned_total_assets = $forgedB.Json.data.total_assets }

$missing = Invoke-ClosureRequest 'GET' $path $null $null
Add-Result 'Missing authentication' $missing.Status '401' ($missing.Status -eq 401) @{}
$invalidSession = New-Object Microsoft.PowerShell.Commands.WebRequestSession
$invalidSession.Cookies.Add([Uri]$ApiBase,(New-Object Net.Cookie('poos_session','invalid-closure-session')))
$invalid = Invoke-ClosureRequest 'GET' $path $null $invalidSession
Add-Result 'Invalid authentication session' $invalid.Status '401' ($invalid.Status -eq 401) @{}

$cashierReport = Invoke-ClosureRequest 'GET' $path $null $cashier.Account.Session
Add-Result 'Lower-privilege cashier financial report' $cashierReport.Status '403' ($cashier.SessionStatus -eq 200 -and $cashier.SessionUser.role -eq 'cashier' -and $cashierReport.Status -eq 403) @{ role = $cashier.SessionUser.role; tenant_id = $cashierTenant }

$exportBody = @{ report_type = 'sales'; format = 'csv'; date_from = $date; date_to = $date; granularity = 'day'; top_limit = 10; locale = 'en' }
$exportA = Invoke-ClosureRequest 'POST' '/api/v1/reports/export' $exportBody $ownerA.Account.Session @{ 'X-Tenant-ID' = "$tenantB"; 'X-Restaurant-ID' = "$restaurantB" }
$exportSafe = $exportA.Status -eq 200 -and $exportA.Content.Contains('123.45') -and -not $exportA.Content.Contains('987.65')
Add-Result 'Owner A export with Tenant B headers' $exportA.Status '200 CSV scoped to A; contains 123.45 and not 987.65' $exportSafe @{ contains_A_total = $exportA.Content.Contains('123.45'); contains_B_total = $exportA.Content.Contains('987.65'); tenant_id = $tenantA; attempted_tenant_id = $tenantB }
$cashierExport = Invoke-ClosureRequest 'POST' '/api/v1/reports/export' $exportBody $cashier.Account.Session
Add-Result 'Lower-privilege cashier export' $cashierExport.Status '403' ($cashierExport.Status -eq 403) @{ role = $cashier.SessionUser.role }

$dbAssertionSql = "SELECT 'journal_ledger_balance='||COALESCE((SELECT (total_debit=total_credit)::text FROM journal_entries WHERE tenant_id=$tenantA AND restaurant_id=$restaurantA AND entry_number='$entryA'),'false'); SELECT 'tenant_a_ledger_rows='||COUNT(*)||';tenant_a_cash='||COALESCE(SUM(gl.running_balance) FILTER (WHERE coa.account_code='$cash'),0) FROM general_ledger gl JOIN chart_of_accounts coa ON coa.id=gl.account_id WHERE gl.tenant_id=$tenantA AND gl.restaurant_id=$restaurantA AND coa.account_code IN ('$cash','$capital'); SELECT 'tenant_b_ledger_rows='||COUNT(*)||';tenant_b_cash='||COALESCE(SUM(gl.running_balance) FILTER (WHERE coa.account_code='$cashB'),0) FROM general_ledger gl JOIN chart_of_accounts coa ON coa.id=gl.account_id WHERE gl.tenant_id=$tenantB AND gl.restaurant_id=$restaurantB AND coa.account_code IN ('$cashB','$capitalB'); SELECT 'export_audit='||tenant_id||':'||restaurant_id||':'||user_id||':'||report_type FROM report_export_audit_logs WHERE request_path='/api/v1/reports/export' ORDER BY id DESC LIMIT 1; SELECT 'cashier_report_fixture_role='||role FROM users WHERE id=$([int]$cashier.SessionUser.id) AND tenant_id=$cashierTenant;"
$dbAssertions = Invoke-DisposableSql $dbAssertionSql
$expectedAudit = "$tenantA`:$restaurantA`:$([int]$ownerA.SessionUser.id):sales"
$auditLine = ($dbAssertions -split "`r?`n" | Where-Object { $_ -like 'export_audit=*' } | Select-Object -First 1)
Add-Result 'Database fixture and export-audit assertions' 200 'balanced A/B journals; two scoped ledger rows each; export audited under A; cashier role persisted' ($dbAssertions.Contains('journal_ledger_balance=true') -and $dbAssertions.Contains('tenant_a_ledger_rows=2;tenant_a_cash=110.00') -and $dbAssertions.Contains('tenant_b_ledger_rows=2;tenant_b_cash=220.00') -and $dbAssertions.Contains("export_audit=$expectedAudit") -and $dbAssertions.Contains('cashier_report_fixture_role=cashier')) @{ assertions = $dbAssertions }

$failures = @($results | Where-Object result -ne 'PASS')
$evidence = [pscustomobject]@{ run_id = $runId; executed_at_utc = [DateTime]::UtcNow.ToString('o'); api_base = $ApiBase; database = 'pos_final_integration_test (disposable tmpfs PostgreSQL 15; no credentials or tokens recorded)'; api_mode = 'real HTTP; actual cookie-session auth; RBAC_ENFORCEMENT=enforce'; fixture_setup = $fixtureOutput; results = $results; overall = $(if ($failures.Count -eq 0) { 'PASS' } else { 'FAIL' }) }
$outputPath = Join-Path $PSScriptRoot 'authenticated-tenant-results.json'
$evidence | ConvertTo-Json -Depth 8 | Set-Content -LiteralPath $outputPath -Encoding utf8
$results | Select-Object name,status,expected,result | Format-Table -AutoSize
"Overall: $($evidence.overall); evidence: $outputPath"
if ($failures.Count -gt 0) { exit 1 }

