# Opening balance reconciliation and safeguards

**Status: OPEN.** This candidate does not post, rewrite, or auto-balance opening entries.

## Evidence and current behavior

- The repository-level journal/GL diagnostic only compares journal lines to immutable ledger effects. It does not infer an opening source, bank statement, cutover, or offset account.
- Existing synthetic coverage includes `TestCFR010PostgresReproducesOpeningBalanceOutsideGL` in `backend/internal/usecase/cfr010_balance_sheet_integration_test.go`. It demonstrates that an opening cash amount held outside the GL produces a separate financial report discrepancy. It is not original operating evidence and does not establish a production root cause.
- No new opening-balance validator, approval metadata enforcement, or migration was included in this candidate.

## Required validation before OPEN-01 implementation

1. Preserve and identify the approved opening Trial Balance and source bank/cash balances by tenant, restaurant, and cutover date.
2. Compare each approved source line to opening journal and GL entries by account, amount, date, period, and source reference.
3. Detect duplicate opening movements, missing source links, historical journal overlap, account mapping mismatch, and cutover-date mismatch.
4. Require approval provenance where the existing workflow supports it. Do not invent a default equity/offset account.
5. Report findings read-only. Any journal creation, correction, migration, or historical restatement must remain separately approved and auditable.

## Approval boundary

OPEN-01 is **PENDING**. Accounting must decide source-of-truth precedence, cutover, account mapping and offset treatment, effective date, correction behavior, and migration/rollback policy. The engineering recommendation is to add deterministic read-only comparisons and block posting until required source and approval metadata are present; this recommendation is not approval.

No original statements or approved opening schedules were recovered in this task. Production opening reconciliation is therefore **NOT ESTABLISHED**.
