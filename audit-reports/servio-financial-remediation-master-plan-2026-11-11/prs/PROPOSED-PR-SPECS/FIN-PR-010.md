# FIN-PR-010 — Financial Export Contracts

**Status: PROPOSED PR specification; not created. Risk: Medium.**

## Objective and task IDs
FIN-T08-01

## Scope and affected files
DISCOVERY REQUIRED: export handlers and UI

## Dependencies and schema
Requires dependencies in remediation-tasks.json and approved policies. Schema impact must be confirmed against deployed baseline; no migration is presumed.

## Tests and acceptance
large file memory, permission, Excel-compatible reader, RTL
Acceptance: screen/export parity, genuine format, formula safety

## Review and rollback
Reviewers: domain owner, accounting reviewer when financial policy/effects change, security reviewer for access changes, DB reviewer for migrations. Rollback: Disable unsafe format path; preserve source ledger.
Merge gate: all acceptance evidence reviewed, no cross-tenant leakage, no unresolved financial discrepancy.
