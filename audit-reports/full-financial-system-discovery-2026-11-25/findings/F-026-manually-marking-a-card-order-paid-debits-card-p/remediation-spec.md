# F-026 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Manually marking a card order paid debits Card Payment Clearing with no clearing/settlement flow.

## 2. Current behavior (verified)
1113 balance 11000 persists

## 3. Required behavior
Clearing cleared on settlement

## 4. Reproduction conditions
Mark card order paid

## 5. Verified evidence
- runtime/24

## 6. Root cause and confidence
CONFIRMED — card payments resolve the debit account to 1113 Card Payment Clearing (auto_posting_engine.go:615-616; domain/auto_posting.go:37; seeded in migration 156:35); no code or migration references 1113 anywhere else, so no settlement/clearing flow exists

## 7. Affected paths
- Backend/database: `backend/internal/usecase/auto_posting_engine.go`:615, `backend/internal/domain/auto_posting.go`:37, `backend/migrations/156_seed_chart_of_accounts.sql`:35, `backend/internal/usecase/auto_posting_engine.go`:615-616
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: POST /api/v1/accounting/card-settlements (new); GET /api/v1/accounting/card-settlements/unsettled (new)
- Accounting logic / posting rules: accounts 1113
- Reports affected: Accounting

## 8. Impact
- Financial: 11000 stuck
- Operational: Accounting workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: add a settlement step for Card Payment Clearing (1113). (1) A "card settlement" operation (proposed POST /api/v1/accounting/card-settlements) records a bank deposit for a date range and amount, posting Dr bank account / Cr 1113 and, when fees are configured, Dr card-fee expense; (2) the clearing balance report lists unsettled card sales by order so the owner can match them to bank statements; (3) manual "mark paid" for card orders stays allowed only with a provider reference (see F-005).

## 10. Alternative designs
A: explicit settlement document (recommended). B: post card sales straight to the bank account (loses the timing difference and fees). C: leave the clearing balance to grow (current state).

## 11. Required policy decisions
POL-07 (non-cash tender settlement) in 56.

## 12. Required data migrations
Yes: settlement header/lines tables linking orders or payments to a settlement.

## 13. Historical data implications
Existing 1113 balances need an opening reconciliation against bank statements.

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
none and shared-file findings F-015, F-040, F-061

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/domain/auto_posting.go, backend/internal/usecase/auto_posting_engine.go, backend/migrations/156_seed_chart_of_accounts.sql).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Clearing reconciliation test (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Settling a batch of card orders moves exactly the settled amount from 1113 to the bank account and leaves the unsettled remainder visible.
- Fees reduce the deposit and post to the fee expense account.
- A second settlement of the same orders is rejected.
