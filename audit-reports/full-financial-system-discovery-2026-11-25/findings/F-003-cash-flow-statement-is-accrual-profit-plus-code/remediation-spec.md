# F-003 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
Cash Flow statement is accrual profit plus code-prefix guesses and does not tie to cash.

## 2. Current behavior (verified)
Operating = net profit (-2000); net cash flow 98000 vs GL cash+bank 93000 | RUN-003 (AUD-MET-FIN): cash-flow net 82,710.01 vs GL cash+bank 115,010.01, residual -32,300.00 (STILL REPRODUCIBLE).  | RUN-003 earlier run: cash-flow 82,710.01 vs GL cash+bank 115,010.01 (residual 32,300.00). Re-measured after more activity: report 172,710.01 vs GL 1110+1111 = 280,010.01, residual 107,300.00. ROOT CAUSE CONFIRMED (AUD-CF-TRACE): usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment that never reads cash/bank accounts and omits every working-capital movement (AR -6,520; inventory +96,000; input VAT -2,800; card clearing -11,000; AP +800; VAT payable +30,820 = +107,300); investing = '12%' debit-credit has inverted sign for asset purchases; financing excludes 21xx/2120 liabilities. 

## 3. Required behavior
Net cash change equals change in cash accounts

## 4. Reproduction conditions
Open Cash Flow after S13

## 5. Verified evidence
- evidence/reports C-s13a-reports_financial_cash_flow.json
- runtime/07 | evidence/network/RUN-003/PH-12/AUD-MET-FIN.json | evidence/screenshots/RUN-003/PH-12/AUD-MET-FIN-cf.png | tests/runtime/ph12_fin_statements.mjs | evidence/network/RUN-003/PH-12/AUD-MET-FIN.earlier-run.json | evidence/database/RUN-003/PH-12/AUD-CF-TRACE.txt | evidence/database/RUN-003/PH-11/AUD-ACC-FLOW.json | evidence/database/RUN-003/PH-14/AUD-MONTH.json

## 6. Root cause and confidence
CONFIRMED — backend/internal/usecase/financial_report_usecase.go:366-399 GenerateCashFlow is an indirect-method fragment: operating = P&L (revenue/expense credit-debit), investing = 12xx debit-credit (sign inverted), financing = equity or 22xx credit-debit; never reads cash/bank accounts and omits working-capital movements (see AUD-CF-TRACE.txt)

## 7. Affected paths
- Backend/database: `backend/internal/usecase/financial_report_usecase.go`:366-399
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: GET /api/v1/reports/financial/cash-flow: add fields opening_cash, closing_cash, reconciliation_difference, working_capital_lines[] additively
- Accounting logic / posting rules: accounts 1110, 1111, 2120
- Reports affected: Financial reports

## 8. Impact
- Financial: 5000 mismatch (this dataset)
- Operational: Financial reports workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: rewrite GenerateCashFlow (backend/internal/usecase/financial_report_usecase.go:366-399) as an indirect-method statement derived from the GL: net income + non-cash items +/- change in working-capital accounts (1120, 1130, 1150, 2101, 2120, others by account range) = operating; asset account movements (12xx) with the correct sign = investing; equity and long-term liability movements (3xxx, 22xx) = financing; closing check: opening cash + net change = closing balance of cash accounts (1110, 1111, 1112, 1113 per policy POL-04). Return the reconciliation line so a mismatch is visible instead of silent.

## 10. Alternative designs
A: indirect method from GL account ranges (recommended). B: direct method from cash/bank account movements grouped by counter-account. C: keep formula and rename the report 'profit and capital summary' (does not solve the defect).

## 11. Required policy decisions
POL-04 (method and classification) in 56.

## 12. Required data migrations
No schema change expected; confirm during implementation.

## 13. Historical data implications
Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.

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
none and shared-file findings F-002, F-011, F-012

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/usecase/financial_report_usecase.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
CF net = delta cash for scenarios (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- For the clean-tenant month dataset (evidence PH-14) the report net change equals the GL cash change 74,480.00 exactly.
- Posting a 20,000 asset purchase (Dr 1250 / Cr 1110) reduces net cash flow by 20,000 and shows investing -20,000.
- The report includes closing-balance reconciliation with difference 0.00; a deliberately inconsistent dataset shows a non-zero difference.
