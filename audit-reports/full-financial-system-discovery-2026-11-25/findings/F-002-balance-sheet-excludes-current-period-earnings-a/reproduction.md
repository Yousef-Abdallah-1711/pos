# F-002 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
S13 loss scenario then open Balance Sheet

## Expected result
A = L + E

## Actual result (measured)
equation_difference -2000.00 (net loss); equity shows 100000 only | RUN-003 clean-tenant month (AUD-MONTH): balance sheet equation_balanced=false, equation_difference -27,710.00 = net loss exactly (REPRODUCED). 

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
