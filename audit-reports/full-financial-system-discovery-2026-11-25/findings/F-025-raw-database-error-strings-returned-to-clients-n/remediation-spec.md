# F-025 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Raw database error strings returned to clients (numeric overflow; constraint names).

## 2. Current behavior (verified)
pq: numeric field overflow; constraint names in messages | RUN-003 real-UI forms: customer, supplier, chart-of-accounts, contract and bank-account forms have no client-side length limits; over-long input reaches the API and the UI user receives raw text such as pq: value too long for type character varying(3) (REPRODUCED; the filler used the same long synthetic value in every text field, so the inputs are artificial but the exposure is real). 

## 3. Required behavior
Sanitised message

## 4. Reproduction conditions
Submit 1e15 amount

## 5. Verified evidence
- runtime/13 | evidence/network/RUN-003/PH-06/AUD-UIFORMS2.json | evidence/network/RUN-003/PH-06/AUD-UIFORMS.json | tests/runtime/ph06_ui_forms2.mjs

## 6. Root cause and confidence
CONFIRMED — handlers pass database errors to clients, e.g. backend/cmd/api/main.go:1302 (fmt.Sprintf("Failed to create user: %v", err)) and respondError(..., err.Error()) at admin_theme_handler_v2.go:80,120,154,361; 128 response sites in internal/handler/http use err.Error()

## 7. Affected paths
- Backend/database: `backend/cmd/api/main.go`:1302, `backend/internal/handler/http/admin_theme_handler_v2.go`:80
- Frontend: none
- API contracts: error bodies gain {code, request_id} additively; message text becomes generic
- Accounting logic / posting rules: none named
- Reports affected: none

## 8. Impact
- Financial: n/a
- Operational: API workflows behave as in section 2.
- Security: Yes — see category: Security/UX
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: introduce one error-mapping layer. (1) A shared helper (proposed backend/internal/handler/http/errors.go) maps known sentinel errors and pq error codes (23505 unique, 23503 foreign key, 23514 check, 22003 numeric overflow, 22001 value too long) to 400/409/422 with a stable code and a generic message; unknown errors log the detail with a request id and return 500 "internal error" plus the request id; (2) replace direct err.Error() / fmt.Sprintf(... %v, err) responses (e.g. backend/cmd/api/main.go:1302, admin_theme_handler_v2.go:80,120,154,361 and the other ~128 sites found by grep) with the helper; (3) add a test that scans handlers for err.Error() in response writers.

## 10. Alternative designs
A: central mapper + CI scan (recommended). B: fix only the known sites (leaks return with every new handler).

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
Keep RBAC module/level checks unchanged; add negative tests for roles without permission.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
Failed requests must leave no partial state.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-018, F-020, F-030, F-031, F-033, F-036, F-037, F-059, F-071

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/main.go, backend/internal/handler/http/admin_theme_handler_v2.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Error-contract tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- No API response body contains pq:, SQLSTATE, constraint names or table names for any induced failure in the fixture battery (duplicate account code, over-long name, numeric overflow).
- Duplicate keys return 409 with a stable error code; check violations return 422.
- Every 500 response carries a request id that appears in the server log.
