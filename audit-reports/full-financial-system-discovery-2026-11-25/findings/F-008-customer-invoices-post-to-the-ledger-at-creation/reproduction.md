# F-008 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
Create invoice with status draft

## Expected result
Draft has no accounting effect

## Actual result (measured)
Journal JE-14 posted for draft invoice (journal_entry_id returned) | RUN-003 SQL (AUD-EVD-F008) PARTIAL: invoices_customer id 2 (INV-2026-00001; now status paid) has journal JE-2026-00014 posted 2026-11-25 22:49:27 Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700, i.e. posted at creation; the draft state itself could not be re-observed because the invoice has since been paid - fresh draft re-run scheduled in T-P09-001.  | RUN-003 FULL REPRODUCTION (AUD-CINV-001, headed Chrome): new invoice id 3 status=draft; journal JE-2026-00050 posted at creation Dr 1120 8,320 / Cr 2120 1,120 / Cr 4100 7,200 (STILL REPRODUCIBLE). 

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
