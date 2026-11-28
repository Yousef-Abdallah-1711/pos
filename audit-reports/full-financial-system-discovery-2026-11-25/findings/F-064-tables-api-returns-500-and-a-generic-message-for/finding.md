# F-064 — Tables API returns 500 and a generic message for an over-long name (the duplicate-table-number claim was RETRACTED in RUN-004)

- **Severity:** Low  ·  **Category:** Validation gap  ·  **Evidence status:** REPRODUCED (name length); duplicate claim RETRACTED
- **Affected modules:** Tables; POS dine-in
- **Root-cause confidence:** CONFIRMED
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-TBL-001, EX-RUN004-FORMMATRIX
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 409 duplicate table number; 400 name length

Observed: RUN-004 correction: duplicate table numbers are rejected with 409 "table number already exists for this restaurant" (constraint uq_restaurant_tables_number, table_repo.go:81-83). The RUN-003 result (TBL-01 201) was a misread: table id 1 had been renumbered to 4 at 20:53, so number 1 did not exist when id 2 was created at 21:10. Still reproduced: a 1,800-char name -> 500 "Failed to create table" (column name varchar(100); the DB error is not mapped to 400). | RUN-003: duplicate table number accepted (201); 1,800-char name -> 500 (REPRODUCED). 

## Severity rationale
Rated **Low** (category: Validation gap; evidence: REPRODUCED (name length); duplicate claim RETRACTED). Basis recorded in the register: 0 (dine-in orders can attach to an ambiguous table). Root-cause confidence: CONFIRMED. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: Tables; POS dine-in. 0 (dine-in orders can attach to an ambiguous table)

## Technical impact
CONFIRMED (RUN-004): the name length is not validated in the handler and the varchar(100) overflow surfaces as a generic 500 (table_repo.go:78-86 wraps the error as "failed to create table"; handler table_handler.go:42 maps it to 500). Uniqueness IS enforced by uq_restaurant_tables_number.

Data integrity: no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access. Security relevance: none recorded. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
0 (dine-in orders can attach to an ambiguous table)

## Where it was reproduced
POST /tables {table_number:1,name:dup,capacity:4} when table 1 exists; POST /tables with a 1,800-character name

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
