# F-079 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Leave requests cannot be created through the UI: the form sends ISO date-times and the API accepts only YYYY-MM-DD (400).

## 2. Current behavior (verified)
POST /hr/leaves returns 400 {"error":"Invalid start_date format (expected YYYY-MM-DD)"}; the page shows "Request failed with status code 400" for every valid submission; the request body carries start_date "2026-11-20T00:00:00.000Z"

## 3. Required behavior
201 and the request appears in the list

## 4. Reproduction conditions
Headed Chrome tenant 4 owner: HR > Leaves > Request Leave, choose valid start and end dates, submit

## 5. Verified evidence
- evidence/network/RUN-004/PH-06/AUD-R4-FORM-MATRIX-employees_leaves_roles.json
- tests/runtime/ph06_r4_form_matrix.mjs

## 6. Root cause and confidence
CONFIRMED — frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114-115 send new Date(...).toISOString(); backend/internal/handler/http/hr_leave_handler.go:186 parses the field as YYYY-MM-DD only

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/components/hr/LeaveForm.tsx`:114, `backend/internal/handler/http/hr_leave_handler.go`:186, `frontend/apps/dashboard/src/components/hr/LeaveForm.tsx`:114-115
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: POST/PUT /api/v1/hr/leaves: accepts both date formats (additive)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0 direct; blocks the leave workflow, so entitlement and payroll deductions cannot be exercised from the UI (see F-054)
- Operational: HR > Leave requests (dashboard form, API) workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: in frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114-115 send data.start_date and data.end_date exactly as the date input returns them (YYYY-MM-DD) instead of new Date(...).toISOString(); additionally make backend/internal/handler/http/hr_leave_handler.go:186 accept both YYYY-MM-DD and RFC3339 (normalise to a date in the tenant time zone, see F-070) so other clients are not affected; add a contract test that posts the exact form payload.

## 10. Alternative designs
A: fix the client and widen the API parser (recommended). B: client only (other API consumers still fail on RFC3339).

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-014, F-054, F-074 and shared-file findings F-054

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/hr_leave_handler.go, frontend/apps/dashboard/src/components/hr/LeaveForm.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Submitting the leave form with valid dates returns 201 and lists the request (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Submitting the leave form with valid dates returns 201 and the request appears in the list with the same dates.
- The API accepts YYYY-MM-DD and RFC3339 and rejects other formats with a field error.
- A form-contract test fails if the payload format changes.
