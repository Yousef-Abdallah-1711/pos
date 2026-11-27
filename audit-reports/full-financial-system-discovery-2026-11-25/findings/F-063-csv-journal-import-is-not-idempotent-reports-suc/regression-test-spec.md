# F-063 — Regression test specification

| test id | type | description | assertion | automation |
|---|---|---|---|---|
| REG-F-063-01 | regression (API/DB) | Re-upload yields the same batch result; all-failed import returns an error status; 6,000-row import returns a job id within 5 s | Expected: Duplicate upload rejected or deduplicated; failed rows reported as an error status; large imports asynchronous or chunked with a job id | automated (Go integration / API) |
| REG-F-063-02 | real Chrome UI | Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved) | UI value equals API value equals stored value; screenshot captured | Playwright headed Chrome |
| REG-F-063-03 | financial reconciliation | Independent oracle comparison (Decimal / integer minor units) for the dataset | Residual 0.00 EGP for this finding; unrelated residuals unchanged | SQL + oracle script |
| REG-F-063-04 | concurrency | Fire N parallel requests (N>=5) for the affected operation | Exactly the arithmetically allowed number succeed; totals equal oracle; no duplicate rows | node/k6 parallel script |

## Preconditions per test
- REG-F-063-01: isolated stack + dataset from reproduction.md (original audit script: tests/runtime/ph11_import.mjs)
- REG-F-063-02: diagnostic tenant, logged in via real login UI (original audit script: tests/runtime/ph11_import.mjs)
- REG-F-063-03: GL, TB and affected report captured before and after (original audit script: tests/runtime/ph11_import.mjs)
- REG-F-063-04: fresh invoice/order with known balance (original audit script: tests/runtime/ph11_import.mjs)

## Pass criteria
All tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.
