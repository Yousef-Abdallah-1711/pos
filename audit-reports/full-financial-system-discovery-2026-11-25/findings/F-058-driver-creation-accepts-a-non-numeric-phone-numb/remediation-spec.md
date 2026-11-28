# F-058 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Driver creation accepts a non-numeric phone number.

## 2. Current behavior (verified)
201 driver id 2 stored with phone_number "abc" | RUN-003 (AUD-THEME-DRIVER / AUD-THEME2): driver phone abc accepted (201) (REPRODUCED). 

## 3. Required behavior
400 phone format error

## 4. Reproduction conditions
POST /admin/drivers with phone_number "abc"

## 5. Verified evidence
- evidence/database/RUN-003/PH-08/AUD-THEME-DRIVER.json
- evidence/database/RUN-003/PH-08/AUD-THEME2.json
- tests/runtime/ph08_theme_driver.mjs
- tests/runtime/ph08_theme2.mjs

## 6. Root cause and confidence
CONFIRMED — driver create requires non-empty fields only (admin_driver_handler.go:146-150); no phone format validation

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/admin_driver_handler.go`:129, `backend/internal/handler/http/admin_driver_handler.go`:146, `backend/internal/handler/http/admin_driver_handler.go`:146-150
- Frontend: none
- API contracts: POST /api/v1/admin/drivers: 400 for invalid phone
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: 0
- Operational: Drivers workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: validate driver phone numbers in backend/internal/handler/http/admin_driver_handler.go (:146-150) with the E.164 validator used by registration, plus email format and licence number length; return 400 with field errors.

## 10. Alternative designs
A: shared validator (recommended). B: client-side only.

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
Existing invalid phones are listed.

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
F-037 and shared-file findings F-057

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/admin_driver_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Malformed phone rejected for drivers, employees, customers (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Phone "abc" returns 400.
- A valid E.164 number is stored normalised.
