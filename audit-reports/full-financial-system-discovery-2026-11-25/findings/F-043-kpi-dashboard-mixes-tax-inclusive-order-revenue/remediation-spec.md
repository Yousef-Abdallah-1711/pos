# F-043 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
KPI dashboard mixes tax-inclusive order revenue with GL expenses, and cash position ignores GL cash accounts.

## 2. Current behavior (verified)
KPI revenue 319,660 (paid orders incl tax) vs P&L revenue 305,300; net_profit 181,540 vs P&L 167,180 (+14,360); cash_position 50,000 (bank_accounts_ledger) vs GL cash+bank 276,270.01 | RUN-003: KPI net profit 181,540 vs P&L 167,180; cash_position 50,000 vs GL cash+bank 276,270.01 (CONFIRMED accounting_advanced_repo.go:838-880). 

## 3. Required behavior
Revenue/expenses/net profit tie to P&L; cash position ties to GL cash and bank accounts

## 4. Reproduction conditions
GET /reports/financial/kpis (tenant 4, window 2026-10-28..2026-11-26)

## 5. Verified evidence
- evidence/database/RUN-003/PH-12/AUD-MET-ORD-KPI.txt
- evidence/network/RUN-003/PH-12/AUD-MET-BATCH.json
- tests/runtime/ph12_batch_apis.mjs

## 6. Root cause and confidence
CONFIRMED — backend/internal/repository/accounting_advanced_repo.go:838-880 revenue=SUM(orders.total_amount paid, created_at), expenses=GL expense lines, cash_position=SUM(bank_accounts_ledger.current_balance)

## 7. Affected paths
- Backend/database: `backend/internal/repository/accounting_advanced_repo.go`:838-880
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`, `frontend/apps/dashboard/src/app/[locale]/dashboard/orders`
- API contracts: GET KPI endpoint: same fields; values change; adds reconciliation block (additive)
- Accounting logic / posting rules: none named
- Reports affected: Financial KPI dashboard

## 8. Impact
- Financial: 14,360 profit overstatement; 226,270.01 cash understatement in this dataset
- Operational: Financial KPI dashboard workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: derive the KPI dashboard from the GL. In backend/internal/repository/accounting_advanced_repo.go:838-880 replace revenue = SUM(orders.total_amount) with the same GL revenue query used by the P&L (financial_report_usecase.go:208), expenses from the GL expense accounts, net profit = revenue - expenses, and cash position = balances of the cash and bank GL accounts (1110/1120-family) instead of bank_accounts_ledger.current_balance; show a reconciliation line when the bank module differs from the GL (see F-029).

## 10. Alternative designs
A: GL-derived KPIs (recommended; one source with the P&L). B: keep operational KPIs and label them (does not fix net profit mixing tax-inclusive revenue with GL expenses).

## 11. Required policy decisions
POL-12 (single source of truth for dashboards/KPIs) in 56.

## 12. Required data migrations
No.

## 13. Historical data implications
None.

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
F-029 and shared-file findings F-009, F-016, F-045, F-047, F-068, F-072, F-073

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/repository/accounting_advanced_repo.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
KPI summary equals P&L and TB-derived cash on controlled data (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- KPI revenue, expenses and net profit equal the P&L for the same period.
- KPI cash position equals the GL cash and bank balance.
- A difference between the bank module and the GL is displayed, not hidden.
