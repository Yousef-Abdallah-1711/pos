# F-030 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Empty-body POST /users returns 500 but creates an active blank user; other empty-body POSTs return 500.

## 2. Current behavior (verified)
500 returned yet user id 6 created (blank email and name; role assigned; 60-char hash of empty password); PUT /users/{id} and POST /notifications also 500

## 3. Required behavior
400 and no row

## 4. Reproduction conditions
POST /users with {} as owner

## 5. Verified evidence
- evidence/runtime/43
- evidence/api-sweep-results.csv

## 6. Root cause and confidence
CONFIRMED — POST /users handler decodes name/email/phone/password with no validation, hashes the (possibly empty) password and calls CreateUserInTenant, then assigns the EMPLOYEE role with the error discarded (_ =) in separate non-atomic steps (backend/cmd/api/main.go:1261-1335)

## 7. Affected paths
- Backend/database: `backend/cmd/api/main.go`:1261, `backend/cmd/api/main.go`:1300, `backend/cmd/api/main.go`:1261-1335
- Frontend: none
- API contracts: POST /api/v1/users: 400 field errors (was 500)
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: n/a
- Operational: Users; Notifications workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: validate and make user creation atomic. Move the inline handler at backend/cmd/api/main.go:1261-1335 into the handler layer; validate name (required, length), email (format, unique per tenant) and password (minimum length/complexity) before hashing; create the user and assign the EMPLOYEE role in one transaction and fail the whole request if the role assignment fails (today the error is discarded with _ =); return 400 with field errors for empty bodies instead of 500.

## 10. Alternative designs
A: validation + transaction (recommended). B: validation only (role assignment can still fail silently).

## 11. Required policy decisions
Password policy parameters are an owner decision.

## 12. Required data migrations
No.

## 13. Historical data implications
Blank active users created by the defect are listed and deactivated after review.

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
none and shared-file findings F-018, F-020, F-025, F-031, F-033, F-036, F-037, F-071

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/main.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
User API validation tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- POST /users with {} returns 400 with field errors and creates no row.
- A failing role assignment rolls back the user.
- Concurrent creation with the same email creates exactly one user.
