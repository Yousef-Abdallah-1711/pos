# F-037 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Tenant contact-info endpoint stores unvalidated email/phone.

## 2. Current behavior (verified)
PUT /tenant/contact-info 200; tenants.email/phone persisted as entered; empty values also accepted; UI shows no validation for phone and only native browser validation for email (bypassed by API)

## 3. Required behavior
Rejected with validation error (400/422)

## 4. Reproduction conditions
Settings > Restaurant: enter email "not-an-email" and phone "abc-###", click Save Changes

## 5. Verified evidence
- evidence/network/RUN-003/PH-05/AUD-SET-001-011.json
- evidence/screenshots/RUN-003/PH-05/AUD-SET-001-invalid-email-after.png
- evidence/screenshots/RUN-003/PH-05/AUD-SET-002-invalid-phone-after.png
- tests/runtime/ph05_settings_forms.mjs

## 6. Root cause and confidence
CONFIRMED — inline handler cmd/api/main.go:845-873 decodes body and runs UPDATE tenants SET email=$1, phone=$2 with no validation or length limit

## 7. Affected paths
- Backend/database: `backend/cmd/api/main.go`:845-873
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/settings`
- API contracts: PUT contact-info: 400 field errors (was 200)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: none direct (tenant billing/notification contact may be wrong)
- Operational: Settings / tenant profile workflows behave as in section 2.
- Security: Yes — see category: Software defect (input validation)
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: move the inline handler at backend/cmd/api/main.go:845-873 to the handler layer, validate email (RFC format) and phone (E.164 using the registration validator), enforce length limits, and return 400 with field errors.

## 10. Alternative designs
A: reuse registration validators (recommended). B: client-side validation only (rejected).

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
Existing invalid values are listed.

## 14. Transaction-integrity requirements
Single-request change; keep existing transaction boundaries.

## 15. Idempotency requirements
Not applicable beyond keeping existing behaviour.

## 16. Authorization requirements
Keep RBAC module/level checks unchanged; add negative tests for roles without permission.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-018, F-020, F-025, F-030, F-031, F-033, F-036, F-071

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/main.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
API contract test: invalid email/phone => 400, no DB change (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Invalid email or phone returns 400 and changes nothing.
- Valid values persist and appear in the dashboard.
