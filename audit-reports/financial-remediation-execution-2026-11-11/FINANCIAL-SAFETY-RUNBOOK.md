# Financial safety runbook — Balance Sheet integrity (Phase 00)

## Detect
Open Dashboard → Reports → Financial → Balance Sheet and select the relevant as-of date. The API response includes the legacy equation flag, explicit integrity_status, equation_difference and integrity_check_basis. A red accessible warning means UNRECONCILED or unrecognized/unverified status. Figures stay visible for investigation. The warning does not fix the accounting.

## Operational treatment
Do not treat an unreconciled or unverified Balance Sheet as final, approved or suitable for external filing. PDF export is withheld when status is not explicitly RECONCILED and the legacy equation flag is false. Do not alter or delete journal/GL rows to clear the warning. Other financial functions are unchanged; any external use remains subject to finance-owner controls.

## Investigate and preserve
Record tenant/restaurant, as-of date, request timestamp, response status and equation difference. Use read-only journal, Trial Balance and GL views for the same scope/date. Compare source journals and ledger rows; note current status only indicates equation arithmetic, not correct attribution, complete posting or accountant approval. Preserve screenshots and relevant log correlation IDs without copying credentials, tokens or sensitive payment data.

## Escalate
Route discrepancy and affected date range to finance owner and qualified accountant; include engineering for query/source diagnosis. If automatic posting failure is suspected, capture operation/source IDs from logs. Logs are not a durable recovery queue. Do not replay provider operations or post manual balancing entries without separately approved procedure.

## Recovery after an approved fix
Re-run independent dated fixtures and reports against isolated test DB; verify no source history was rewritten. Production correction/backfill requires separate authorization, approved accountant policy, a reviewed migration/repair plan and verified backup. This runbook grants no production-data write authority.