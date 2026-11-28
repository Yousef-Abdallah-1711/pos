# F-012 — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
At least six revenue/profit definitions disagree and are unlabeled.

## 2. Current behavior (verified)
Sales 121000; Profitability 110000; KPI 22000; GL revenue 20000; KPI net profit 6000 vs P&L 4000; KPI cash position 0 vs GL cash

## 3. Required behavior
Reports state basis or reconcile

## 4. Reproduction conditions
Open Sales Profitability KPI P&L after mixed paid/unpaid orders

## 5. Verified evidence
- 43-report-metric-runtime-results.csv

## 6. Root cause and confidence
CONFIRMED — three independent queries define revenue differently: Sales = SUM(o.total_amount) tax-inclusive (report_handler.go:2177), Profitability = SUM(oi.total_price) net of tax for non-cancelled orders (report_handler.go:2367), P&L = GL revenue-type accounts (financial_report_usecase.go:208); no shared definition or label; numeric bridge in evidence/network/RUN-003/PH-12/AUD-REV-TRACE.txt

## 7. Affected paths
- Backend/database: `backend/internal/handler/http/report_handler.go`:2177, `backend/internal/handler/http/report_handler.go`:2367, `backend/internal/usecase/financial_report_usecase.go`:208
- Frontend: `frontend/apps/dashboard/src/app/[locale]/dashboard/reports`
- API contracts: report responses gain additive "basis" fields; no field renamed
- Accounting logic / posting rules: none named
- Reports affected: Reports; Dashboard

## 8. Impact
- Financial: up to 101000 spread
- Operational: Reports; Dashboard workflows behave as in section 2.
- Security: No direct security impact identified.
- Data integrity: Yes — monetary/ledger data can become inconsistent (see financial-impact.md).

## 9. Proposed remediation design (Proposed)
Proposed: publish one definition per concept and label every report with its basis. (1) Adopt the decisions in POL-01 (revenue recognition basis) and POL-12 (single source of truth for dashboards/KPIs): "Sales" (operational, tax-inclusive or net as approved), "Profitability revenue" (net of tax, paid or all as approved) and "GL revenue" (accrual); (2) add a basis label and tooltip to each tile and CSV/PDF header (frontend reports pages and report_handler.go responses gain a "basis" field additively); (3) move the three queries (report_handler.go:2177, :2367 and financial_report_usecase.go:208) behind one shared revenue-definition package so a basis change happens in one place; (4) provide a reconciliation panel showing the bridge Sales -> Profitability -> GL (tax, unpaid orders, cancelled orders).

## 10. Alternative designs
A: shared definitions + labels + bridge (recommended). B: label only, keep three queries (cheap, but definitions can drift again). C: make every report read the GL (consistent, but operational reports lag the outbox and unpaid orders disappear).

## 11. Required policy decisions
POL-01 and POL-12 in 56; implementation blocked until approved.

## 12. Required data migrations
No.

## 13. Historical data implications
No stored data changes; historical exports remain as exported (their basis was undeclared).

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
none and shared-file findings F-002, F-003, F-011, F-046, F-067

## 21. Risks
- Regression in modules sharing the same source files (backend/internal/handler/http/report_handler.go, backend/internal/usecase/financial_report_usecase.go).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
Cross-report reconciliation suite (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
- Every revenue figure in the UI, CSV and PDF shows its basis label.
- For the RUN-003 dataset the bridge reproduces 517,000 -> 470,000 -> 290,000 with tax 47,000 and unpaid 180,000 explained line by line.
- A contract test fails if two reports that claim the same basis return different totals for the same dataset.
