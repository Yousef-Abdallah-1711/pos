# F-020 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + DB): default Employee role is seeded with HR=READ, SALARIES=NONE (auth_repo.go seedDefaultRolesAndPermissions :380), but GET /hr/employees is guarded only by module 3 READ (cmd/api/main.go:877) and serializes base_salary (domain/hr.go:44), so a cashier reads salaries (RUN-003 SEC-009)


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/auth_repo.go`:380.
- Propagation: affected modules RBAC; HR; POS; related findings none recorded; observed consequence: GET /hr/employees returns base_salary; POS page Access denied; cannot mark orders paid | RUN-003 (AUD-SEC-BATTERY SEC-009): cashier (default Employee role: HR READ, SALARIES NONE) GET /hr/employees -> 200 with base_salary 5000 and contact d
- Cause versus symptom: the symptom is what users or reports show (GET /hr/employees returns base_salary; POS page Access denied; cannot mark orders paid | RUN-003 (AUD-SEC-BATTERY SEC-009): cashier (default); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-SEC-009, EX-RUN005-RBAC-MATRIX; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/auth_repo.go | 380 | backend/internal/repository/auth_repo.go | exact | yes |
| backend/cmd/api/main.go | 877 | backend/cmd/api/main.go | exact | yes |
| backend/internal/domain/hr.go | 44 | backend/internal/domain/hr.go | exact | yes |
| auth_repo.go | - | backend/internal/repository/auth_repo.go | prefix | yes |

