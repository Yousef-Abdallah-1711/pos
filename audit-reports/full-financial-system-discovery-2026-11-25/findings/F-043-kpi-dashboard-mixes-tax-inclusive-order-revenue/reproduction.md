# F-043 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
GET /reports/financial/kpis (tenant 4, window 2026-10-28..2026-11-26)

## Expected result
Revenue/expenses/net profit tie to P&L; cash position ties to GL cash and bank accounts

## Actual result (measured)
KPI revenue 319,660 (paid orders incl tax) vs P&L revenue 305,300; net_profit 181,540 vs P&L 167,180 (+14,360); cash_position 50,000 (bank_accounts_ledger) vs GL cash+bank 276,270.01 | RUN-003: KPI net profit 181,540 vs P&L 167,180; cash_position 50,000 vs GL cash+bank 276,270.01 (CONFIRMED accounting_advanced_repo.go:838-880). 

## Re-runnable scripts / evidence
- `tests/runtime/ph12_batch_apis.mjs`

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
