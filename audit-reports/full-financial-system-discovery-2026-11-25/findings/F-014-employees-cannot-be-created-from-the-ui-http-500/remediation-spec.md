# F-014 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Employees cannot be created from the UI (HTTP 500; UI shows generic toast, raw DB error exposed in API response body).

## 2. Current behavior (verified)
500 positive_working_hours violation; payload keys weekly_hours/monthly_salary vs backend working_hours_per_week/base_salary | RUN-003 (headed Chrome 150, tenant 4 owner): POST /hr/employees -> 500 reproduced; DB before=after (no orphan rows); UI shows only generic 'Request failed with status code 500' - raw DB error is in the API response body (see evidence). 

## 3. Required behavior
Employee created

## 4. Reproduction conditions
Fill EmployeeForm and Create

## 5. Verified evidence
- network (emp_create run) | evidence/screenshots/RUN-003/PH-10/AUD-HR-001-before-submit.png
- evidence/screenshots/RUN-003/PH-10/AUD-HR-001-after-submit.png
- evidence/network/RUN-003/PH-10/AUD-HR-001.json
- evidence/database/RUN-003/PH-10/AUD-HR-001-before-after.txt
- tests/runtime/ph10_emp_create_headed.mjs

## 6. Root cause and confidence
CONFIRMED — Frontend-backend contract mismatch

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx`:77, `backend/internal/handler/http/hr_employee_handler.go`:139-150
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/hr`
- API contracts: unchanged unless the design below says otherwise (see 08-api-inventory.csv)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: payroll blocked until API fixture
- Operational: HR workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: align the EmployeeForm payload with the backend contract (frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx:76-77 vs backend/internal/handler/http/hr_employee_handler.go:139-150) and make the handler return 400 with the offending field instead of a 500 with raw database text; add a contract test that posts the form's actual payload.

## 10. Alternative designs
Alternative A: fix at the source layer named in root-cause.md; Alternative B: add a compensating check in the calling layer (weaker, not recommended unless the source layer is vendor/owner-controlled).

## 11. Required policy decisions
No accounting policy decision is required; product decisions are limited to the behaviour described in Required behavior.

## 12. Required data migrations
No schema change expected; confirm during implementation.

## 13. Historical data implications
No financial history is affected; existing data needs no correction unless noted in the evidence.

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
none

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/hr_employee_handler.go, frontend/apps/dashboard/src/components/hr/EmployeeForm.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Employee create UI/API contract test (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Creating an employee through the real UI in headed Chrome returns 201 and the employee appears in the list.
- Invalid payloads return 400 naming the field; no raw DB error text in any response.
