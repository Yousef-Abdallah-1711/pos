# F-020 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + DB): default Employee role is seeded with HR=READ, SALARIES=NONE (auth_repo.go seedDefaultRolesAndPermissions :380), but GET /hr/employees is guarded only by module 3 READ (cmd/api/main.go:877) and serializes base_salary (domain/hr.go:44), so a cashier reads salaries (RUN-003 SEC-009)


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/auth_repo.go | 380 | backend/internal/repository/auth_repo.go | exact | yes |
| backend/cmd/api/main.go | 877 | backend/cmd/api/main.go | exact | yes |
| backend/internal/domain/hr.go | 44 | backend/internal/domain/hr.go | exact | yes |
| auth_repo.go | - | backend/internal/repository/auth_repo.go | prefix | yes |

