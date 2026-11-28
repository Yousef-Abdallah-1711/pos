# F-027 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Crawl dashboard as Employee user: 50 of 79 routes render; 24 show Access denied; finance report pages render with 403 data calls

## Expected result
Page guard consistent with API permissions

## Actual result (measured)
Shells render (P&L, BS, CF, VAT, aging, payroll, RBAC roles) with blank data; only 24 pages gated

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
