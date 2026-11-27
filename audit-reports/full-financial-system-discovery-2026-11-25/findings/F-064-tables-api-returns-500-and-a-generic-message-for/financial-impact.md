# F-064 — Financial impact

- **Quantified impact (register):** 0 (dine-in orders can attach to an ambiguous table)
- **Accounts / statements named in the evidence:** none named
- **Impact classes:** concurrency, validation

## What is wrong in the books or reports
RUN-004 correction: duplicate table numbers are rejected with 409 "table number already exists for this restaurant" (constraint uq_restaurant_tables_number, table_repo.go:81-83). The RUN-003 result (TBL-01 201) was a misread: table id 1 had been renumbered to 4 at 20:53, so number 1 did not exist when id 2 was created at 21:10. Still reproduced: a 1,800-char name -> 500 "Failed to create table" (column name varchar(100); the DB error is not mapped to 400). | RUN-003: duplicate table number accepted (201); 1,800-char name -> 500 (REPRODUCED). 

## Which reports and decisions are affected
No financial statement is directly changed; impact is operational/security (see remediation-spec.md).

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.

