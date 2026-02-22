# Phase 9 DB-Backed Verification Path

This runbook closes the remaining runtime checks in Phase 9 that require a real database-backed backend.

## Purpose

Use this when `PAYMOB_MODE=mock` is running with PostgreSQL connected (not stub mode), so these checks can be validated end-to-end:

- Phase `9.2` protected API success path (`/subscription/initiate`, `/subscription/payments`, `/subscription/current`)
- Phase `9.3` amount mismatch rejection and duplicate webhook idempotency with persisted payment rows

## Added Verification Assets

- Seed SQL: `backend/scripts/paymob_phase9_seed.sql`
- DB-backed verification script: `backend/scripts/phase9_db_verify.ps1`
- Embedded DB daemon (no Docker): `backend/scripts/embedded_pg_daemon.go`
- SQL seed runner (no `psql` required): `backend/scripts/sql_exec.go`

## What The Script Verifies

`phase9_db_verify.ps1` performs:

1. Starts PostgreSQL in Docker (optional) and applies seed SQL.
2. Starts backend in `mock` mode with DB env vars.
3. Runs API checks:
   - `GET /api/v1/subscription/plans`
   - `POST /api/v1/subscription/initiate` (auth context)
   - `GET /api/v1/subscription/payments`
   - `GET /api/v1/subscription/current`
   - missing auth check on protected endpoint (`401`)
4. Runs webhook security checks:
   - success webhook
   - duplicate webhook (idempotent processing)
   - amount mismatch webhook (rejected)
   - SQL-injection-like plan input (rejected)
5. Starts backend in `live` mode on a second port and verifies tampered HMAC rejection (`401`).

Results are written to:

- `backend/scripts/phase9_db_results.json`

## Run It (Windows / PowerShell)

From repo root:

```powershell
powershell -ExecutionPolicy Bypass -File .\backend\scripts\phase9_db_verify.ps1
```

Optional flags:

```powershell
powershell -ExecutionPolicy Bypass -File .\backend\scripts\phase9_db_verify.ps1 `
  -UseDocker:$true `
  -StartBackend:$true `
  -TenantID 1 `
  -RestaurantID 1
```

## Docker-Free Fallback (Embedded PostgreSQL)

When Docker is unavailable, use embedded PostgreSQL and still run full DB-backed checks:

```powershell
cd backend
go run .\scripts\embedded_pg_daemon.go
```

In another terminal:

```powershell
cd backend
$env:DB_HOST='127.0.0.1'
$env:DB_PORT='5432'
$env:DB_USER='postgres'
$env:DB_PASSWORD='postgres'
$env:DB_NAME='pos_saas'
$env:DB_SSLMODE='disable'

go run .\cmd\migrate\main.go
go run .\scripts\sql_exec.go -file .\scripts\paymob_phase9_seed.sql
.\scripts\phase9_db_verify.ps1 -UseDocker:$false
```

## Preconditions

- Docker Desktop running (if `-UseDocker:$true`)
- Port `5432` available for PostgreSQL
- Ports `8080` and `8081` available for backend runs
- Migrations already applied (Docker path mounts `backend/migrations` by default)

## Expected Outcome

Mark these as complete when script output confirms:

- `9.2` protected endpoints return `200` in DB-backed mode.
- `9.3` amount mismatch is rejected and duplicate webhook does not create duplicate payment state.
- Tampered HMAC in live mode returns `401`.
