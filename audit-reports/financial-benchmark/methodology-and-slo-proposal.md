# Financial query benchmark methodology and SLO proposal

## Scope

This benchmark is engineering evidence from a disposable PostgreSQL database. It is not production traffic, a capacity forecast, or release acceptance. The job uses a pinned PostgreSQL image with `pg_stat_statements`, captures the runner and database settings, and uploads the JSON report as a workflow artifact.

## Dataset and workload

- The primary tenant receives the selected journal volume (default 10,000; dispatch choices 1,000, 10,000, or 50,000), posting dates distributed across 364 days of the year, five ledger lines per seeded journal, and a high-volume account history.
- A second tenant receives a smaller independently posted journal history. Its trial balance is measured separately.
- The primary tenant also receives non-empty AR invoices and receipts, AP invoices and supplier payments, inventory valuation data, and inventory movement history.
- Workloads cover GL account pagination, movement-derived running balances, Trial Balance, P&L, Balance Sheet, Cash Flow, AR/AP aging, inventory valuation, journal and account CSV exports, sales CSV/PDF exports, backdated postings, and a mixed read/write workload.
- Each single-worker measurement uses five warmups and 30 measured samples. The mixed workload uses four workers, four warmups, and 30 measured samples per worker.
- The GL index comparison uses the same fixture and query in sequence: first with the selected index removed, then with it restored. Both plans, index size, PostgreSQL query execution totals, latency percentiles, throughput, and sample error counts/rates are recorded.
- Query errors fail the benchmark job. A published successful measurement therefore records zero errors; failed jobs remain failures and must not be interpreted as performance acceptance.

## Interpretation

Compare before/after values only when the report shows the same commit family, dataset counts, PostgreSQL settings, runner specifications, workload parameters, and sample counts. The benchmark currently performs an index A/B comparison; it does not establish a universal baseline for other code or hardware changes. Repeat runs are needed to characterize variance before drawing conclusions.

## Proposed SLO process — pending owner approval

Adopt workload-specific latency objectives for GL pagination, financial reports, aging/valuation, exports, and concurrent posting. For each workload, the product and operations owners should approve a target percentile, a target value, the dataset/concurrency profile, the measurement window, and the allowable error rate. The report provides p50/p95/p99, throughput, and error metrics to support that decision, but this file sets no numeric target. Until approval is recorded, the benchmark is engineering evidence only and creates no release threshold.
