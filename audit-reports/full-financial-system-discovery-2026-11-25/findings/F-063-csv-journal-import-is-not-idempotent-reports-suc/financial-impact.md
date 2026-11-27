# F-063 — Financial impact

- **Quantified impact (register):** Duplicated or unknown-outcome imports distort the ledger once posted (auto_post=true posts immediately)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** accounting, concurrency, report

## What is wrong in the books or reports
identical re-upload creates a second journal (IMP-03); unbalanced, unknown account, bad date and wrong header all return HTTP 200 success:true with failed_entries 1 and created_entries 0; a 6,000-journal import ran ~122 s, the client socket was closed (no HTTP response) while the server completed all 6,000 journals, so a client retry duplicates them (17,000 draft journals created in the audit tests); formula characters are stored verbatim (export escapes them) | RUN-003 import battery: re-upload duplicates, all-failed imports answer 200 success, 6,000-journal import completed server-side after the client socket closed (REPRODUCED). 

## Which reports and decisions are affected
Financial reports that read the affected data: Accounting CSV import; journal entries

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

