# F-066 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Headed Chrome, tenant 4 owner: open /en/dashboard; compare the Products tile with GET /products?page=1&page_size=1 (pagination.total=12) and SQL (13 rows, 1 with status=deleted)

## Expected result
Products tile = 12 (non-deleted menu items)

## Actual result (measured)
Tile shows 1. The dashboard requests page_size=1 and reads total at the payload root, but the API returns it under pagination.total, so the code falls back to items.length = 1

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
