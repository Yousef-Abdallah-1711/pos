# CFR-010 continuation — 2026-11-22

Main `966a5004fc594d67f8ef1850909bd7f8e963ecf5`; candidate [#114](https://github.com/Yousef-Abdallah-1711/pos/pull/114) `02b145f9e7ee140238906847fdc37b16407b6a62`. Verify live state before acting. #108, #109, #110 and #112 also remain open. Independent reviews are absent.

Synthetic sale: assets 5,500 cents, liabilities 500, equity 0, P&L 5,000, residual 5,000. Balance Sheet excludes unclosed nominal earnings from equity. EARN-01, OPEN-01 and CLOSE-01 are pending. Original operating discrepancy is not established.

#114 contains cache bypass and historical P&L void correction plus Cash Flow movement-date correction `88037b2`. Existing Cash Flow formula gives +50/-50/0 while cash GL is +55/-55/0 for a taxed synthetic journal; the formula is not approved cash-flow methodology. Local targeted PostgreSQL regression and handler tests passed; actual sale/payment/refund, live network and manual Dashboard acceptance remain outstanding. The API launch was rejected by execution policy before startup; request an environment where launching a disposable local service is permitted, then use synthetic credentials and isolated PostgreSQL/Redis. Do not bypass restrictions.

Calendar boundary and bounded 20-post/300-read Balance Sheet tests passed. The matrix has 9 bounded PASS, 9 POLICY_PENDING, 3 PARTIAL, 2 NOT_RUN; actual refund remains unrun. Prior performance p50 174.044 vs 180.726 ms in separate hosted runs, no approved budget. Controlled comparison remains pending. Source evidence lives on #114 under `audit-reports/cfr010-balance-sheet-fix-2026-11-22/`; fetch that branch to inspect it.
