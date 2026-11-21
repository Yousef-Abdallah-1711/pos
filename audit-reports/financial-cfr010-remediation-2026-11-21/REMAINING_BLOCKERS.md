# Remaining blockers

| ID | Impact | Owner | Required action / evidence | State |
|---|---|---|---|---|
| CFR010-SRC-01 | Original operating root cause cannot be established or closed. | Finance data owner | Provide provenance-preserving original dates/scope, signed opening TB, source bank balances, opening journals, journal/GL exports, reports, close state, and producing environment/commit; or provide formal authorized disposition. | Open; no production access used |
| CFR010-POLICY-01 | Opening, earnings, and close behavior cannot be changed safely. | Qualified accountant/controller; fiscal calendar owner for close | Record signed OPEN-01, EARN-01, CLOSE-01 decisions with scope/effective date, implications, and correction/rollback authority. | Open; all three pending |
| CFR010-PERF-01 | Reconciliation scan performance is unmeasured. | Engineering / database owner | Run reconciliation-specific high-volume disposable workload; capture dataset counts, PostgreSQL config, EXPLAIN ANALYZE BUFFERS, p50/p95/max, concurrency and timeout. | Open; performance is not claimed |
| CFR010-CROSS-01 | Full Trial Balance/GL/P&L/Balance Sheet/Cash Flow source consistency diagnostics are incomplete. | Accounting engineering | Implement remaining policy-neutral checks and validate date boundaries, scope filters, classifications, missing/duplicate sources against independent SQL oracles. | Partial; current run verifies exercised report fixtures only |
| CFR010-OPEN-01 | Opening source matching and safeguards are not implemented. | Accounting engineering + qualified approver | Build read-only source-to-GL/cutover/duplicate/mapping/overlap diagnostics after the OPEN-01 contract is approved. | Open |
| CFR010-ACCEPT-01 | CFR-010 is not release-accepted. | Engineering, finance, independent reviewer | Final exact-SHA CI, review, PR/merge, performance and source/policy closure evidence. | Open; production remains NO-GO |

Completed on candidate SHA `e68d70d`: real disposable PostgreSQL reconciliation, authenticated API checks, EN/AR desktop/mobile browser acceptance, integration suite, and independent code review. Details are in `POSTGRESQL_ACCEPTANCE.md`, `SECURITY_ACCEPTANCE.md`, `BROWSER_ACCEPTANCE.md`, and `PR_AND_CI_EVIDENCE.md`.
