# Scenario matrix

| Scenario | Result | Evidence / notes |
| --- | --- | --- |
| Authenticated order → payment event → journal → GL → exports/reports | PASS, synthetic | API log records 5,500-cent order, one paid event, balanced 5,500-cent trial balance, CSV/PDF checks; see transaction trace. |
| Duplicate paid callback | PASS | One event/journal effect; repeated callback returned HTTP 200. |
| Accountant journal and cashier denial | PASS | Balanced posted journal; cashier post denied; duplicate post rejected. |
| AR invoice and partial receipt | PASS | 10,000-cent invoice, 4,000-cent receipt, 6,000-cent open balance. |
| Payroll lifecycle and rollback/retry | PASS | Synthetic 150,000-cent gross/net; failed post rolled back; retry posted. |
| Supplier invoice and payment | PASS | 22,000-cent invoice and payment journal; invoice paid. |
| Concurrent supplier payments / overpayment | PASS | Four successes, 20,000 cents paid, 2,000 cents payable; repeat/overpayment rejected. |
| Browser UI, Arabic/English, mobile/desktop, RBAC, tenant isolation | PASS, 11/11 | `financial-browser-tests.txt` and JUnit report. |
| Integration suite | PASS, 1,011 passed; 12 skipped | Opt-in live/browser/benchmark acceptance was run separately where applicable. |
| Balance Sheet accounting equation | FAIL / UNRECONCILED | Assets 5,500, liabilities 500, equity 0; 5,000-cent residual. No balancing posting made. |
| Cash Flow classification | POLICY_PENDING | API/SQL net 5,000 cents; classification requires approved accounting policy. |
| OPEN-01 opening balance, EARN-01 current earnings, CLOSE-01 close/lock | BLOCKED / PENDING | No approved accounting behavior executed. |
| Cash-shift close, bank reconciliation, inventory/COGS, full refunds, split payments, worker crash/retry, fiscal year boundary | NOT RUN / NOT ESTABLISHED | No pass inferred. |
| Original operating-data reconciliation and financial release acceptance | BLOCKED | Source evidence and accounting approvals unavailable. |

Raw results: `financial-live-api-tests.jsonl`, `financial-browser-tests.txt`, `financial-browser-junit.xml`, `financial-acceptance-tests.jsonl`, and `financial-benchmark-results-1000.json`.
