# F-061 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
A paid order with total 0 (free item) creates an accounting event that can never post ("at least two journal lines are required").

## 2. Current behavior (verified)
order 73 total 0.00 paid: outbox status retryable attempts 1/8 last_error "at least two journal lines are required"; it will exhaust retries and appear as a failed event needing manual redrive; no journal exists | RUN-003: zero-total paid order 73 -> outbox retryable, error at least two journal lines are required (CONFIRMED auto_posting_engine.go:535). 

## 3. Required behavior
Zero-value sale completes without a journal (or posts a zero-line memo) and the event is completed

## 4. Reproduction conditions
Sell a product with price 0 (product creation allows it), mark the order paid, read accounting_event_outbox

## 5. Verified evidence
- evidence/database/RUN-003/PH-16/AUD-MASTERDATA.json
- evidence/database/RUN-003/PH-16/AUD-MD-NEGPRICE.json
- tests/runtime/ph16_masterdata.mjs
- tests/runtime/ph16_negprice.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/usecase/auto_posting_engine.go:535-537 rejects any posting with fewer than two lines; PostSalesEntry builds no lines for a zero total instead of returning nil

## 7. Affected paths
- Backend/database: `backend/internal/usecase/auto_posting_engine.go`:50-142, `backend/internal/usecase/auto_posting_engine.go`:535-537
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: unchanged
- Accounting logic / posting rules: none named
- Reports affected: POS/Orders; Accounting event outbox

## 8. Impact
- Financial: 0 value, but permanent failed-event noise hides real failures and needs manual redrive
- Operational: POS/Orders; Accounting event outbox workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: complete zero-value events. In backend/internal/usecase/auto_posting_engine.go (the two-line check at :535-537 and PostSalesEntry/PostCOGSEntry) return a distinct "nothing to post" result when the total is zero and mark the outbox event completed with a note, instead of building an empty journal that can never post and retries forever; keep rejecting genuinely unbalanced journals.

## 10. Alternative designs
A: short-circuit zero-total events (recommended). B: forbid zero-total paid orders (policy decision on free items, see F-062 price validation).

## 11. Required policy decisions
Whether free items are allowed at all belongs with the product validation decision (F-062).

## 12. Required data migrations
No; stuck events are completed by a one-off script after review.

## 13. Historical data implications
Stuck zero-total events (retryable:1 in tenant 4) are listed and closed.

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
F-062 and shared-file findings F-015, F-026, F-040

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/usecase/auto_posting_engine.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Zero-total paid order leaves the outbox completed with no journal (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Paying a zero-total order completes its outbox event with status completed and no journal.
- The outbox shows no retryable entry for it.
- A non-zero order still posts exactly one journal.
