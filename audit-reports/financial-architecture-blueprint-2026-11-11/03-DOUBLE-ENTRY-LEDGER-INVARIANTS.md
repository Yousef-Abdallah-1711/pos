# 03 — Double-entry ledger invariants (PROPOSED)

Evidence: historical audit FIN-001/003/008/009 and source paths in 21. This is target behavior, not current behavior.

A posted journal has at least two eligible lines; each line has exactly one positive debit or credit; all lines share tenant, restaurant, currency and functional basis. Debits equal credits exactly at approved currency precision. Material imbalance is rejected before persistence. Any allowed rounding is a separate line to an approved rounding account, bounded by an approved tolerance and retained alongside the original source amounts. Header totals equal recomputed line totals. Posting transaction writes journal, lines, GL and posted outcome atomically. A unique source-effect key prevents duplicate effects. Posted lines and attribution are immutable; corrections are linked reversals/replacements. Closed-period posts require policy-approved reopen or next-period correction.

GL is a one-to-one projection of posted journal lines; reconcile identity, amount, currency, account and dimensions. Running balances are derived/cache values. Trial Balance debits equal credits for declared tenant, restaurant, currency, cutoff and posting watermark.

Acceptance: failures before/after line insertion leave no partial posted state; repeated concurrent delivery creates one effect; independent events remain independent; correction preserves original. Failure injection and concurrency were NOT TESTED. FIN-003 changed a final business line to force balance; FIN-001 omitted current earnings in the Balance Sheet. Sources: `backend/internal/usecase/auto_posting_engine.go:633-660`, `backend/internal/repository/journal_entry_repo.go`, `general_ledger_repo.go`, `financial_report_usecase.go:261-336`.
