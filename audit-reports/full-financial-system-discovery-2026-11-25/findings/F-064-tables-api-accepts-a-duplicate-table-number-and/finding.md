# F-064 — Tables API accepts a duplicate table number and returns 500 for an over-long name

- **Severity:** Low  ·  **Category:** Validation gap  ·  **Evidence status:** REPRODUCED
- **Affected modules:** Tables; POS dine-in
- **Root-cause confidence:** PROBABLE
- **Executed audit tests referencing this finding:** EX-RUN003-AUD-TBL-001
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
Expected: 409 duplicate table number; 400 name length

Observed: duplicate -> 201 (two tables numbered 1, ids 1 and 2); long name -> 500 "Failed to create table"; capacity (1-50), negative number and delete-with-active-orders (409) are validated correctly | RUN-003: duplicate table number accepted (201); 1,800-char name -> 500 (REPRODUCED). 

## Monetary / operational impact
0 (dine-in orders can attach to an ambiguous table)

## Where it was reproduced
POST /tables {table_number:1,name:dup,capacity:4} when table 1 exists; POST /tables with a 1,800-character name

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
