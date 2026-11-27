# F-070 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED for KDS: kds_repo.go:418-430 EXTRACT(HOUR FROM received_at) and CURRENT_DATE at :314-317; DB timezone is UTC


## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/kds_repo.go | 418 | backend/internal/repository/kds_repo.go | exact | yes |
| backend/internal/repository/kds_repo.go | 314 | backend/internal/repository/kds_repo.go | exact | yes |
| kds_repo.go | 418-430 | backend/internal/repository/kds_repo.go | prefix | yes |

