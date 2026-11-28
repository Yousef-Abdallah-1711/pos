# F-032 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (source): frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx exists in the shipped dashboard app


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx`.
- Propagation: affected modules Dashboard app; related findings none recorded; observed consequence: Page renders TEST PAGE LOADED with no auth protection and a hydration error RUN-005: GET http://localhost:13112/en/test-url returned HTTP 200 (curl) on the current isolated dashboard.
- Cause versus symptom: the symptom is what users or reports show (Page renders TEST PAGE LOADED with no auth protection and a hydration error RUN-005: GET http://localhost:13112/en/test-url returned HTTP 20); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx | - | frontend/apps/dashboard/src/app/[locale]/test-url/page.tsx | exact | yes |

