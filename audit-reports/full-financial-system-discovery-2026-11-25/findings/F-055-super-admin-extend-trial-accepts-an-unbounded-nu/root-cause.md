# F-055 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source): ExtendTenantTrial validates only days > 0 (subscription_handler.go:128-138); there is no upper bound


## Origin and propagation
- Origin (first verified source location): `backend/internal/handler/http/subscription_handler.go`:137.
- Propagation: affected modules Super Admin tenant management; SaaS access control; related findings none recorded; observed consequence: 200 days_added 100000, new_trial_end_date 2300-10-02 (negative days correctly rejected) | RUN-003: extend-trial 100000 days accepted -> 2300-10-02 (REPRODUCED; tenant 8 restored by fixture). 
- Cause versus symptom: the symptom is what users or reports show (200 days_added 100000, new_trial_end_date 2300-10-02 (negative days correctly rejected) | RUN-003: extend-trial 100000 days accepted -> 2300); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-AUD-SAAS-005; scripts: tests/runtime/ph13_saas.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| backend/internal/handler/http/subscription_handler.go | 137 | backend/internal/handler/http/subscription_handler.go | exact | yes |
| subscription_handler.go | 128-138 | backend/internal/handler/http/subscription_handler.go | prefix | yes |

