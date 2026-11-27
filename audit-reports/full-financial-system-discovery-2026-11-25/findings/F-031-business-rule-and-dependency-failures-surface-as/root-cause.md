# F-031 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): error mapping is done by passing err text to the client (see F-025 sites) and returning 500 for constraint violations (e.g. main.go:1302; hr handlers rely on DB check constraints)


## Origin and propagation
- Origin (first verified source location): `backend/cmd/api/main.go`:1302.
- Propagation: affected modules Orders; Receipts; Platform; related findings F-025; observed consequence: Paid-order cancel 500 (item cancel 409); receipt during outage 400 with raw role error pq role pos_accounting not permitted; DB down: /health 200 degraded and authenticated calls 403 TENANT_ACCESS_BLOCKED contact support | RUN-003 (AUD-ACC-
- Cause versus symptom: the symptom is what users or reports show (Paid-order cancel 500 (item cancel 409); receipt during outage 400 with raw role error pq role pos_accounting not permitted; DB down: /healt); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-ACC-001, EX-RUN003-AUD-ATT-001, EX-RUN003-AUD-LEV-001, EX-RUN003-AUD-CON-001, EX-RUN003-AUD-MD-001, EX-RUN003-AUD-UIFORM-001, EX-RUN003-AUD-UIFORM-003, EX-RUN004-FORMMATRIX, EX-RUN004-FORMMATRIX2; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/cmd/api/main.go | 1302 | backend/cmd/api/main.go | exact | yes |

