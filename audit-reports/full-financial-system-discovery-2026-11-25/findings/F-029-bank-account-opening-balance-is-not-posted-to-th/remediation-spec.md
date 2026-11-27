# F-029 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Bank account opening balance is not posted to the ledger; bank module and KPI cash differ from GL.

## 2. Current behavior (verified)
No journal; account shows Current balance 50000 while its transaction list runs -10000; KPI cash_position 50000 vs GL 1111 -10000 | RUN-003 (AUD-ACC-FLOW BNK-01): new bank account opening_balance 12,000 linked to GL 1112 creates no journal (GL 1112 stays 0) (REPRODUCED). 

## 3. Required behavior
Opening journal or clear disclosure; balances agree

## 4. Reproduction conditions
Create bank account via UI with opening balance 50000

## 5. Verified evidence
- evidence/runtime/38 | evidence/database/RUN-003/PH-11/AUD-ACC-FLOW.json | tests/runtime/ph11_accounting_flow.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/bank_account_repo.go:101-122 inserts opening_balance/current_balance into the bank ledger with no journal entry (the file only reads journal ids in queries at :224-270); the bank handler has no posting code; fresh account in RUN-003 left GL 1112 = 0

## 7. Affected paths
- Backend/database: `backend/internal/repository/bank_account_repo.go`:101-122, `backend/internal/handler/http/bank_account_handler.go`:97
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`, `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: POST /api/v1/accounting/bank-accounts: response includes opening_journal_entry_id (additive)
- Accounting logic / posting rules: accounts 1111, 1112
- Reports affected: Bank accounts; Reports KPI; Accounting

## 8. Impact
- Financial: 60000 gap
- Operational: Bank accounts; Reports KPI; Accounting workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: post the bank opening balance. In backend/internal/repository/bank_account_repo.go:101-122 create the account and a balanced opening journal in one transaction (Dr bank account GL / Cr opening balance equity, or Cr the account chosen by the accountant) dated at the opening date; make the KPI cash position read the GL (see F-043) so the bank module and the ledger agree.

## 10. Alternative designs
A: opening journal at creation (recommended). B: require the accountant to enter the opening journal manually and block the account until it exists. C: ignore the opening balance in the bank module (rejected: modules disagree).

## 11. Required policy decisions
Proposed POL-17 (opening balance equity account and date) must be added to 56 before implementation.

## 12. Required data migrations
No schema change; data-correction journals for existing accounts behind accountant approval.

## 13. Historical data implications
Existing bank accounts with unposted opening balances are listed with the difference to the GL.

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
none

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/bank_account_handler.go, backend/internal/repository/bank_account_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Bank vs GL reconciliation (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Creating a bank account with opening balance 50,000 produces one journal Dr bank / Cr equity of 50,000 and GL cash-and-bank equals the bank module balance.
- Re-saving the account does not post again.
- The balance sheet stays balanced.
