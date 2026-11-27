# F-008 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Customer invoices post to the ledger at creation while status is draft.

## 2. Current behavior (verified)
Journal JE-14 posted for draft invoice (journal_entry_id returned) | RUN-003 SQL (AUD-EVD-F008) PARTIAL: invoices_customer id 2 (INV-2026-00001; now status paid) has journal JE-2026-00014 posted 2026-11-25 22:49:27 Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700, i.e. posted at creation; the draft state itself could not be re-observed because the invoice has since been paid - fresh draft re-run scheduled in T-P09-001.  | RUN-003 FULL REPRODUCTION (AUD-CINV-001, headed Chrome): new invoice id 3 status=draft; journal JE-2026-00050 posted at creation Dr 1120 8,320 / Cr 2120 1,120 / Cr 4100 7,200 (STILL REPRODUCIBLE). 

## 3. Required behavior
Draft has no accounting effect

## 4. Reproduction conditions
Create invoice with status draft

## 5. Verified evidence
- runtime TX-011 | evidence/database/RUN-003/PH-20/AUD-EVD-F008-F015-sql.txt | evidence/screenshots/RUN-003/PH-09/AUD-CINV-001-form-before-submit.png | evidence/screenshots/RUN-003/PH-09/AUD-CINV-001-after-submit.png | evidence/network/RUN-003/PH-09/AUD-CINV-001.json | evidence/database/RUN-003/PH-09/AUD-CINV-001-before-after.txt | tests/runtime/ph09_custinv_draft_headed.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/invoice_customer_repo.go:452-560 createInvoiceTx calls syncJournalEntryTx (:539) which builds (buildJournalEntry :561-636) and posts the AR/revenue/VAT journal for every invoice at creation regardless of status; draft invoices therefore hit AR, revenue and VAT

## 7. Affected paths
- Backend/database: `backend/internal/repository/invoice_customer_repo.go`:452-636, `backend/internal/repository/invoice_customer_repo.go`:452-560
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices`, `frontend/apps/dashboard/src/app/[locale]/dashboard/accounting`
- API contracts: POST /api/v1/accounting/customer-invoices: response journal_entry_id is null for drafts; POST /api/v1/accounting/customer-invoices/{id}/issue: posts the journal (confirm exact route in 08-api-inventory.csv)
- Accounting logic / posting rules: accounts 1120, 2120, 4100
- Reports affected: Customer invoices

## 8. Impact
- Financial: 5200 AR+revenue recognised at draft
- Operational: Customer invoices workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: post to the ledger only when an invoice leaves draft. In backend/internal/repository/invoice_customer_repo.go createInvoiceTx (:452-560) skip syncJournalEntryTx (:539) when status = draft; add the same call to the issue/send transition and to the approve transition; editing a draft must never touch the ledger. Existing draft invoices that already carry a journal are listed by an exception query and reversed by approved journals.

## 10. Alternative designs
A: post on issue (recommended; standard accrual practice). B: keep posting on create but label drafts as posted (rejected: misleading). C: separate "proforma" document type (larger change).

## 11. Required policy decisions
POL-09 (accounting effect of draft documents) in 56.

## 12. Required data migrations
No schema change; a data-correction script (reversal of journals for drafts) behind accountant approval.

## 13. Historical data implications
Draft invoices with posted journals (e.g. INV-2026-00002 / JE-2026-00050 in tenant 4) must be listed and reversed explicitly.

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
- Regression in modules sharing the same source files (backend/internal/repository/invoice_customer_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Draft invoice accounting test (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Creating a draft invoice creates no journal entry and leaves AR (1120) and revenue (4100) unchanged.
- Issuing the invoice creates exactly one journal; re-issuing does not create a second.
- Deleting or voiding a draft leaves the ledger unchanged.
