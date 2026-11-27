# F-072 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Bank reconciliation is not tied to the ledger: book balance is the sum of client-supplied line amounts, so a session can be "completed" with invented numbers and stamps the bank account as reconciled.

## 2. Current behavior (verified)
Session status completed, difference 0.00, bank_accounts_ledger.last_reconciled_balance = 777,777.77 and last_reconciled_date set, while the ledger balance is 12,000.00 and the journal has 0 bank-account lines. A session with difference 123,456.78 or 500.00 is also stamped completed_at although its status stays in_progress, and the complete call returns HTTP 200

## 3. Required behavior
Book balance comes from the GL/bank ledger; lines must reference journals that touch that bank account; a session completes only when matched amounts agree with the statement and the ledger

## 4. Reproduction conditions
Tenant 4 bank account 2 (ledger 12,000.00, GL 0): create a session with statement balance 777,777.77, add one line referencing a real journal that has no bank line with book_amount 777,777.77, complete the session (API-DIAGNOSTIC)

## 5. Verified evidence
- evidence/database/RUN-004/PH-11/AUD-R4-BANK-RECON.json
- evidence/database/RUN-004/PH-11/AUD-R4-BANK-RECON2.json
- tests/runtime/ph11_r4_bank_recon.mjs
- tests/runtime/ph11_r4_bank_recon2.mjs

## 6. Root cause and confidence
CONFIRMED — refreshReconciliationSummaryTx sums bank_reconciliation_lines.book_amount supplied by the client (accounting_advanced_repo.go:2126-2149); AddBankReconciliationLine inserts without checking that the journal touches the bank GL account or that amounts equal the journal line (:113-166); CompleteBankReconciliationSession writes last_reconciled_* from the summed client amounts and sets completed_at regardless of status (:169-230)

## 7. Affected paths
- Backend/database: `backend/internal/repository/accounting_advanced_repo.go`:113, `backend/internal/repository/accounting_advanced_repo.go`:169, `backend/internal/repository/accounting_advanced_repo.go`:2126, `backend/internal/handler/http/accounting_advanced_handler.go`:58, `backend/internal/repository/accounting_advanced_repo.go`:2126-2149
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: POST /api/v1/accounting/bank-reconciliation/sessions: book_balance ignored (derived); POST .../sessions/{id}/lines: journal line id required; 400 for invalid lines; POST .../sessions/{id}/complete: 409 when difference is not zero
- Accounting logic / posting rules: none named
- Reports affected: Accounting > Bank reconciliation, bank accounts

## 8. Impact
- Financial: Control failure: a reconciled flag can be produced without any real matching (illustrated with 777,777.77 against a 12,000.00 ledger balance)
- Operational: Accounting > Bank reconciliation, bank accounts workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: rebuild reconciliation on ledger data. (1) Session creation reads the book balance from the GL for the bank account (general_ledger rows of bank_accounts_ledger.gl_account_id up to the statement date) and stores it; the client no longer supplies book_balance; (2) AddBankReconciliationLine (backend/internal/repository/accounting_advanced_repo.go:113-166) takes the journal line id, verifies it is a line on the bank GL account in the same tenant/restaurant, takes statement/book amounts from that line (the client supplies only the statement amount to match) and rejects a line already matched in another session; (3) refreshReconciliationSummaryTx (:2126-2149) computes difference = statement balance - (ledger balance adjusted by outstanding items) from those matched lines; (4) CompleteBankReconciliationSession (:169-230) completes only when the difference is zero (or an explicit adjustment journal is recorded) and otherwise returns 409 without setting completed_at; last_reconciled_* are written only on completion; (5) the dashboard shows the ledger balance next to the statement balance.

## 10. Alternative designs
A: ledger-derived book balance with line validation (recommended). B: keep manual entry but add a visible "unverified" flag (does not remove the fabricated-completion risk). C: automatic statement import and matching (larger, later phase).

## 11. Required policy decisions
Reconciliation tolerance and adjustment-journal rules are an accounting decision (proposed tolerance 0.01).

## 12. Required data migrations
Yes: unique matched journal-line per active session and a CHECK/foreign-key set tying lines to the session tenant; back-fill review of existing sessions.

## 13. Historical data implications
Existing completed sessions must be re-verified against the ledger; list sessions whose book balance differs from the GL at the statement date.

## 14. Transaction-integrity requirements
Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.

## 15. Idempotency requirements
Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.

## 16. Authorization requirements
Reconciliation completion should require a finance permission distinct from data entry.

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{"error":..., "success":false}`.

## 18. Recovery requirements
If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
F-029, F-043 and shared-file findings F-009, F-016, F-043, F-045, F-047, F-048, F-063, F-068, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/accounting_advanced_handler.go, backend/internal/repository/accounting_advanced_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Session with invented amounts cannot complete; matched real lines complete and update last_reconciled values; non-zero difference returns 409 (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- A session with a line whose journal has no bank-account line is rejected (400).
- Completing a session whose difference is not zero returns 409, leaves status in_progress and completed_at null.
- A session matched to real ledger lines completes and sets last_reconciled_balance equal to the ledger balance at the statement date.
- The book balance shown equals the GL balance of the bank account (see F-029).
