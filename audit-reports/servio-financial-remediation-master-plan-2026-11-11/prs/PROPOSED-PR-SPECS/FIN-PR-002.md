# FIN-PR-002 — Money and Ledger Invariants

**Status: PROPOSED PR specification; not created. Risk: High.**

## Objective and task IDs
FIN-T01-01, FIN-T01-02

## Scope and affected files
DISCOVERY REQUIRED: money/domain/API/database contracts

## Dependencies and schema
Requires dependencies in remediation-tasks.json and approved policies. Schema impact must be confirmed against deployed baseline; no migration is presumed.

## Tests and acceptance
Approved policy, real PostgreSQL, property suite, tenant checks
Acceptance: precision boundaries, balanced journals, immutability, period constraints

## Review and rollback
Reviewers: domain owner, accounting reviewer when financial policy/effects change, security reviewer for access changes, DB reviewer for migrations. Rollback: Forward correction preferred after data writes; no silent rounding.
Merge gate: all acceptance evidence reviewed, no cross-tenant leakage, no unresolved financial discrepancy.
