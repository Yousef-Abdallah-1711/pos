# Balance Sheet Scenario Results

Amounts are in currency units (1.00 = 100 cents). A `POLICY_PENDING` status means the production expected equity or disposition is not approved. `NOT_RUN` means this work has no execution evidence for that scenario.

| # | Scenario | Expected / actual evidence | Result |
|---:|---|---|---|
| 1 | Single cash sale | Source 55.00; assets 55.00, liabilities 5.00, equity 0.00; residual 50.00 | POLICY_PENDING (observed) |
| 2 | Multiple sales | Two synthetic posted sales: assets 77.00, liabilities 7.00, equity 0, residual/P&L 70.00, Trial Balance 77.00/77.00; independent GL oracle matches | POLICY_PENDING |
| 3 | Sale with tax | Sale above; tax payable 5.00 | POLICY_PENDING (observed) |
| 4 | Sale with discount | No isolated fixture result captured | NOT_RUN |
| 5 | Refund/reversal | Synthetic partial reversal journal: assets 66.00, liabilities 6.00, equity 0, residual/P&L 60.00, Trial Balance 66.00/66.00; independent GL oracle matches. The customer refund workflow itself was not run. | POLICY_PENDING (journal reversal); refund workflow NOT_RUN |
| 6 | Net profit | P&L net 50.00; Balance Sheet residual 50.00 | POLICY_PENDING |
| 7 | Net loss | Integration fixture P&L -20.00; open-period residual -20.00 | POLICY_PENDING |
| 8 | Zero activity | Clean report fixture returned zero totals and `RECONCILED` before posting | PASS for zero-activity behavior |
| 9 | Prior-period earnings | 2025 P&L net 200.00; pre-close residual 200.00 | POLICY_PENDING |
| 10 | Current-period earnings | Current-period income/loss not rolled into equity | POLICY_PENDING |
| 11 | Existing retained earnings | Explicit test close: assets/equity 1,150.00; retained 200.00 | PASS for explicit synthetic journal math; production policy pending |
| 12 | Closing journal present | Explicit close fixture balanced at 1,150.00 | PASS for fixture; production policy pending |
| 13 | Fiscal-year boundary | Date filters tested; approved fiscal-year boundary not configured/evidenced | NOT_RUN |
| 14 | Opening balance missing | Opening balance outside GL reproduced; report includes GL -100.00, not bank opening 1,000.00 | PASS for defect reproduction |
| 15 | Opening balance with authorized mapping | No approved mapping supplied | POLICY_PENDING |
| 16 | Incorrect account classification | No dedicated fixture result captured | NOT_RUN |
| 17 | Duplicate journal effects | Duplicate journal post rejected in prior browser/API acceptance; no duplicate GL effect shown | PASS for duplicate post rejection |
| 18 | Voided journal | Date handling tested by PR #108; dedicated Balance Sheet void scenario not run here | NOT_RUN |
| 19 | Multiple restaurants, one tenant | Integration test adds sibling-restaurant posting and confirms no report change | PASS |
| 20 | Separate tenants | Integration test adds another tenant posting and confirms no report change | PASS |
| 21 | Negative balances | Existing unit test for contra-balance signs passes | PASS for tested sign case |
| 22 | High-precision decimal boundary | No dedicated boundary fixture captured | NOT_RUN |
| 23 | Concurrent posting/report read | No dedicated concurrency fixture captured | NOT_RUN |

The request for full acceptance is not complete: scenarios marked `NOT_RUN` and approval-gated cases remain blockers. No scenario is upgraded to a production accounting pass by this cache correction.

The multiple-sales and reversal evidence is from `TestCFR010SalesAndReversalAccounting` on a newly migrated disposable PostgreSQL database. Its exact figures are in currency units and it uses a scoped, separate GL aggregate as its oracle.
