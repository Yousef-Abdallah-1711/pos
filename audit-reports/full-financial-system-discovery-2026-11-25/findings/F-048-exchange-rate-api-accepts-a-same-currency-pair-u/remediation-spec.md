# F-048 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Exchange-rate API accepts a same-currency pair (USD/USD).

## 2. Current behavior (verified)
200 {id:2}: pair stored | RUN-003: USD/USD pair accepted (id 2). 

## 3. Required behavior
Rejected (base must differ from quote)

## 4. Reproduction conditions
POST /accounting/exchange-rates base_currency=USD quote_currency=USD rate=1

## 5. Verified evidence
- evidence/database/RUN-003/PH-11/AUD-ACC-FLOW2.json
- tests/runtime/ph11_accounting_flow2.mjs

## 6. Root cause and confidence
CONFIRMED — the handler only requires non-empty codes and a positive rate (accounting_advanced_handler.go:475-482); there is no base != quote check

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/accounting_advanced_handler.go`, `backend/internal/handler/http/accounting_advanced_handler.go`:481, `backend/internal/handler/http/accounting_advanced_handler.go`:475-482
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: POST exchange-rates: 400/409 for the listed cases
- Accounting logic / posting rules: none named
- Reports affected: Accounting - exchange rates

## 8. Impact
- Financial: 0 (data quality)
- Operational: Accounting - exchange rates workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: validate currency pairs. In backend/internal/handler/http/accounting_advanced_handler.go (:475-482) require ISO-4217 uppercase codes from an allow-list, base != quote, and a rate within a sane range (> 0 and < configurable maximum); enforce uniqueness of (tenant, base, quote, effective_date) in the repository (accounting_advanced_repo.go:596).

## 10. Alternative designs
A: handler validation + unique index (recommended). B: handler only.

## 11. Required policy decisions
POL-11 (currency policy) in 56.

## 12. Required data migrations
Optional unique index; list existing duplicates first.

## 13. Historical data implications
Existing same-currency pairs are listed and deactivated.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
No authorization change; keep existing wrapWithPermission levels.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
none and shared-file findings F-009, F-063, F-068, F-072, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/accounting_advanced_handler.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Same-currency and malformed code rejected (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- USD/USD returns 400.
- An unknown code returns 400.
- Two rates for the same pair and date return 409.
