# F-033 — Reproduction

## Environment
Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.

## Steps (as recorded in the register)
API sweep as Employee user

## Expected result
Least privilege

## Actual result (measured)
GET /admin/orders/export 200; GET /subscription/current and /subscription/payments 200; GET /hr/employees (with base_salary), /hr/leaves, /hr/attendance, /hr/roles 200 RUN-005: the same enforce-mode matrix shows the EMPLOYEE role reading HR employee records with salary fields; billing/subscription route GET /payments/subscription returned 404 for all roles in this build, so the billing part of the original claim is not re-verified.

## Re-runnable scripts / evidence
- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
