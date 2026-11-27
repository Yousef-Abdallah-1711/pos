# F-033 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source + RUN-003): same seeding/guard pattern as F-020 (backend/internal/repository/auth_repo.go:380 default Employee permissions; backend/cmd/api/main.go route wrappers) lets the Employee role read HR and export orders


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/auth_repo.go | 380 | backend/internal/repository/auth_repo.go | exact | yes |
| backend/cmd/api/main.go | 877 | backend/cmd/api/main.go | exact | yes |

