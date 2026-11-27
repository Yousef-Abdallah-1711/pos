# F-033 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + RUN-003): same seeding/guard pattern as F-020 (backend/internal/repository/auth_repo.go:380 default Employee permissions; backend/cmd/api/main.go route wrappers) lets the Employee role read HR and export orders


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/auth_repo.go`:380.
- Propagation: affected modules Orders; Subscription; HR; related findings F-020; observed consequence: GET /admin/orders/export 200; GET /subscription/current and /subscription/payments 200; GET /hr/employees (with base_salary), /hr/leaves, /hr/attendance, /hr/roles 200 RUN-005: the same enforce-mode matrix shows the EMPLOYEE role reading HR
- Cause versus symptom: the symptom is what users or reports show (GET /admin/orders/export 200; GET /subscription/current and /subscription/payments 200; GET /hr/employees (with base_salary), /hr/leaves, /h); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-SEC-009; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/auth_repo.go | 380 | backend/internal/repository/auth_repo.go | exact | yes |
| backend/cmd/api/main.go | 877 | backend/cmd/api/main.go | exact | yes |

