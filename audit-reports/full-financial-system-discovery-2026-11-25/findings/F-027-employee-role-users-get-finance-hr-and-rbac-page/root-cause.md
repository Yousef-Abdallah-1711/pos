# F-027 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-004 source, partial scope): only 6 dashboard page files reference a permission hook; e.g. settings/roles/page.tsx has no guard reference, so the shell renders for any signed-in role while the API returns 403; other pages guard inconsistently


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx`.
- Propagation: affected modules Frontend routing; RBAC; related findings none recorded; observed consequence: Shells render (P&L, BS, CF, VAT, aging, payroll, RBAC roles) with blank data; only 24 pages gated
- Cause versus symptom: the symptom is what users or reports show (Shells render (P&L, BS, CF, VAT, aging, payroll, RBAC roles) with blank data; only 24 pages gated); the cause is the first break above.

## Investigation performed
- Executed audit tests: none recorded in 34b; scripts: reproduced through UI/API/SQL steps in reproduction.md.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx | - | frontend/apps/dashboard/src/app/[locale]/dashboard/settings/roles/page.tsx | exact | yes |
| settings/roles/page.tsx | - | frontend/apps/dashboard/src/app/[locale]/auth/callback/page.tsx | basename-ambiguous(117) | yes |

