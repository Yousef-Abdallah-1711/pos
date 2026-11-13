# 03 — Phase 00 changes

## Implemented

- `backend/internal/usecase/financial_report_usecase.go`: BalanceSheetReport now returns integrity_status, equation_difference and integrity_check_basis. GenerateBalanceSheet derives status from the pre-existing round2 equation check and preserves all totals.
- `backend/internal/usecase/order_usecase.go`: automatic paid-order sales/COGS/reload failures now log source_event_type, operation, tenant, restaurant, source event ID, failure category, recovery_state=not_durable and correlation_id=unavailable. Raw error text is not logged. Payment/order success semantics remain unchanged; no durable recovery is claimed.
- `backend/internal/usecase/financial_report_usecase_test.go`: balanced/unreconciled status and unchanged total fixtures.
- `backend/internal/usecase/accounting_posting_diagnostic_test.go`: structured field and sensitive error-text omission check.
- `backend/internal/usecase/financial_report_postgres_integration_test.go`: real isolated PostgreSQL report diagnostic fixture.
- `backend/internal/handler/http/financial_report_integrity_handler_test.go`: additive API fields and existing report totals contract.
- Dashboard Balance Sheet type/component/page/locales: accessible localized warning; original figures remain visible; PDF export requires explicit RECONCILED and equation_balanced true.
- Dashboard helper/tests: old/unknown/inconsistent status blocks PDF export.

## Exclusions

No current-period earnings accounting, ledger mutation, auto-normalizer change, attachment/schema fix, outbox, payment behavior change, monetary precision change, tax configuration, refund/inventory accounting, migration or backfill.