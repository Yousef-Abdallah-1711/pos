# F-015 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Create salary then Mark Paid

## Expected result
Approval gate; accrual per PAY-01

## Actual result (measured)
Paid with is_approved=false; journal only at payment as direct expense | RUN-003 SQL re-verification (AUD-EVD-F015): salaries id 1 and 2 are status=paid with is_approved=false; journals JE-2026-00007/00049 are Dr 5210 5,000 / Cr 1111 5,000 posted at payment only; no accrual journal exists (STILL REPRODUCIBLE; UI re-run pending T-P10-003). 

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
