# F-070 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED for KDS: kds_repo.go:418-430 EXTRACT(HOUR FROM received_at) and CURRENT_DATE at :314-317; DB timezone is UTC


## Origin and propagation
- Origin (first verified source location): `backend/internal/repository/kds_repo.go`:418.
- Propagation: affected modules KDS analytics and stats; related findings F-044; observed consequence: API hours equal the UTC buckets exactly (1:7, 17:4, 19:4, 20:18, 22:11); the Cairo buckets differ (0:11, 3:7, 19:4, 21:4, 22:18); kds_queue.received_at is timestamp without time zone; no tenant time-zone setting was found
- Cause versus symptom: the symptom is what users or reports show (API hours equal the UTC buckets exactly (1:7, 17:4, 19:4, 20:18, 22:11); the Cairo buckets differ (0:11, 3:7, 19:4, 21:4, 22:18); kds_queue.); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-MET-KDS003; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/repository/kds_repo.go | 418 | backend/internal/repository/kds_repo.go | exact | yes |
| backend/internal/repository/kds_repo.go | 314 | backend/internal/repository/kds_repo.go | exact | yes |
| kds_repo.go | 418-430 | backend/internal/repository/kds_repo.go | prefix | yes |

