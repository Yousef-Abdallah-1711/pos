# F-024 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Frontend sends float artifacts (630.0000000000001) in monetary payloads.

## 2. Current behavior (verified)
tax_amount 630.0000000000001 (server stores rounded)

## 3. Required behavior
Exact decimals

## 4. Reproduction conditions
Create invoice with 14% tax

## 5. Verified evidence
- network CI-C-001

## 6. Root cause and confidence
CONFIRMED — InvoiceForm computes discountAmount, taxAmount and lineTotal with raw floating-point multiplication and sends them unrounded (frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx:620-650)

## 7. Affected paths
- Backend/database: `frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx`:623, `frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx`:637, `frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx`:620-650
- Frontend: none
- API contracts: unchanged (values are rounded)
- Accounting logic / posting rules: none named
- Reports affected: Invoices

## 8. Impact
- Financial: n/a
- Operational: Invoices workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: compute money in integer minor units (cents) in frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx (calculateLineValues :620-635 and mapLineToPayload :637-650): convert unit price and quantity to cents once, apply discount and tax with an explicit rounding rule (half-up to 2 decimals, the same rule the backend trigger uses), and send values rounded to two decimals; better, send only quantity, unit price, discount rate and tax rate and let the server return the authoritative amounts for display.

## 10. Alternative designs
A: integer cents + one rounding rule shared with the server (recommended). B: toFixed(2) at the end (still accumulates float error across lines). C: server-only calculation with preview endpoint (best, larger change).

## 11. Required policy decisions
Rounding rule (half-up per line vs per document) under POL-02.

## 12. Required data migrations
No.

## 13. Historical data implications
Existing rows with float artefacts are rounded in the numeric(…,2) columns already; list any that differ from the rounded value.

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
- Regression in modules sharing the same source files (frontend/apps/dashboard/src/components/accounting/InvoiceForm.tsx).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Precision tests (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- The payload for 5,000 x 1 with 10% discount and 14% VAT contains tax_amount 630 and line_total 5130 exactly (no 630.0000000000001).
- For 1,000 random lines the UI total equals the server-stored total to the cent.
