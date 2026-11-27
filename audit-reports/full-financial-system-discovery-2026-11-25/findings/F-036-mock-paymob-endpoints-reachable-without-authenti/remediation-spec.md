# F-036 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Mock Paymob endpoints reachable without authentication in non-production environments.

## 2. Current behavior (verified)
Return 400 (reachable) without token

## 3. Required behavior
Mock disabled outside test

## 4. Reproduction conditions
Unauthenticated GET /mock/paymob/iframe and POST /mock/paymob/process

## 5. Verified evidence
- evidence/api-sweep-results.csv

## 6. Root cause and confidence
CONFIRMED — backend/cmd/api/main.go:667-668 registers GET /mock/paymob/iframe and POST /mock/paymob/process with mux.HandleFunc and no auth wrapper (non-production only); RUN-003 PAY-01 reached the handler unauthenticated (400 payment_token required)

## 7. Affected paths
- Backend/database: `backend/cmd/api/main.go`:667-668, `backend/internal/handler/http/mock_paymob_handler.go`:46-85
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/products`
- API contracts: /mock/paymob/* removed or gated
- Accounting logic / posting rules: none named
- Reports affected: Payments (mock)

## 8. Impact
- Financial: n/a
- Operational: Payments (mock) workflows behave as in section 2.
- Security: Yes — see category: Security (environment)
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: gate the mock Paymob routes registered at backend/cmd/api/main.go:667-668 behind an explicit flag (PAYMOB_MODE=mock AND ENABLE_MOCK_PAYMENT_ROUTES=true) and require the same authentication or a signed token as the real flow; never register them when PAYMOB_MODE is live; with F-071 also refuse mock mode when ENV is staging.

## 10. Alternative designs
A: explicit flag + signature (recommended). B: remove the mock routes and use a test double in tests only.

## 11. Required policy decisions
None.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
Keep RBAC module/level checks unchanged; add negative tests for roles without permission.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-018, F-020, F-025, F-030, F-031, F-033, F-037, F-071

## 21. Risks
- Regression in modules sharing the same source files (backend/cmd/api/main.go, backend/internal/handler/http/mock_paymob_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Env guard tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- With PAYMOB_MODE=live the mock routes return 404.
- In mock mode they require the flag and a valid token.
- Staging refuses to start in mock mode.
