# F-038 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-005 source + RUN-003 runtime): in the current landing code (frontend/apps/landing/src/app/[locale]/auth/register/page.tsx) onSubmit sets the slug error with setError when availability is not available (lines 59-64) and nothing clears it later: the only useEffect (lines 20-24) handles billing and analytics, availability is read from useDebouncedAvailability (line 39) without an effect that calls clearErrors, and the schema resolver does not know availability, so the manual error persists until the field is edited or the form is submitted again. RUN-003 reproduced the persistence in headed Chrome (2 of 3 early-submit scenarios, error visible >= 8 s while the field showed "available"). The landing source moved under the committed landing workstream; line numbers are for HEAD 35100f7


## Origin and propagation
- Origin (first verified source location): `frontend/apps/landing/src/app/[locale]/auth/register/page.tsx`:34-63.
- Propagation: affected modules Landing / registration; related findings none recorded; observed consequence: No POST sent; red error persists >=8 s while the same field shows "Subdomain is available"; second click succeeds (clearErrors)
- Cause versus symptom: the symptom is what users or reports show (No POST sent; red error persists >=8 s while the same field shows "Subdomain is available"; second click succeeds (clearErrors)); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN003-OBS-ONB-001-S1, EX-RUN003-OBS-ONB-001-S4, EX-RUN003-OBS-ONB-001-S5; scripts: tests/runtime/ph05_subdomain_probe.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | 34-63 | frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | exact | yes |
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | 59 | frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | exact | yes |
| frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | 39 | frontend/apps/landing/src/app/[locale]/auth/register/page.tsx | exact | yes |

