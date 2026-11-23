# Concurrency and failure tests

## Observed results

- Duplicate paid callback: one event/journal effect; retry returned success without duplicate financial posting.
- Duplicate manual journal post: rejected with HTTP 400 and no additional GL effect.
- Payroll forced post failure: rollback verified, retry succeeded, duplicate mark-paid did not create another journal.
- Concurrent supplier payment test: four successful requests settled 20,000 cents against a 22,000-cent invoice; remaining payable 2,000 cents. Repeated payment and overpayment requests were rejected.
- Initial 10,000-journal performance run seeded/reconciled but stopped before measurements because `pg_stat_statements` had not been preloaded. PostgreSQL configuration was corrected only in the disposable environment; subsequent 1,000-journal benchmark completed.
- Prior #108 retry-limit test passed in isolated suite. This does not prove the root cause of its earlier CI failure or establish that the original failure was flaky.

Raw application test output: `financial-live-api-tests.jsonl`, `financial-acceptance-tests.jsonl`; benchmark output: `financial-benchmark-run-1000.txt`.

## Not established

Worker crash/restart and production queue recovery were not validated here. The earlier intermittent worker test's mechanism remains unestablished. Do not label that failure a confirmed flake or infrastructure issue based only on a successful retry.
