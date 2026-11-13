# PHASE 06 — Financial reporting

**Status: PROPOSED. No implementation performed.**

## Objective and scope
Repair FIN-001 under approved close policy; then reconcile statements and comparative periods at shared watermark. FIN-001 minimum repair is urgent Phase 00 prerequisite.

## Tasks and dependencies
- **FIN-T06-01 — Repair earnings and Balance Sheet equation** (FIN-001). Depends on: FIN-T00-01, FIN-T01-01, FIN-T01-02. Acceptance: Assets=liabilities+equity at all approved cutoffs, including current-period earnings.
- **FIN-T06-02 — Reconcile all financial reports to ledger watermark** (REPORTS). Depends on: FIN-T06-01, FIN-T04-01, FIN-T05-01. Acceptance: All statements agree to posted ledger and documented policy.

## Files and migrations
- backend/internal/usecase/financial_report_usecase.go
- DISCOVERY REQUIRED: report usecases/queries
Migrations require backup/restore proof, precondition inventory, forward validation, locking analysis, and forward correction. Exact schema changes are DISCOVERY REQUIRED unless explicitly identified.

## Tests, rollback and Definition of Done
- four historical dated fixtures
- earnings close/cutoff and TB bridge
- P&L, cash flow, aging, VAT, inventory and consolidation independent fixtures
Rollback cannot erase posted history or provider side effects; use idempotent retry or approved compensating entries. Done means independent evidence is reviewed, alerts and runbooks exist, all acceptance criteria pass, and no NO-GO condition remains. Block on unclear accounting policy, unknown live schema, duplicates, irreconcilable history, or cross-tenant leakage. Handoff only the approved contracts and unresolved exceptions to the next phase.
