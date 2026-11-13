# FIN-PR-006 — Order and Payment Financial Integration

**Status: PROPOSED PR specification; not created. Risk: High.**

## Objective and task IDs
FIN-T04-01

## Scope and affected files
DISCOVERY REQUIRED: active POS/web/admin/payment paths

## Dependencies and schema
Requires dependencies in remediation-tasks.json and approved policies. Schema impact must be confirmed against deployed baseline; no migration is presumed.

## Tests and acceptance
Provider sandbox/mock, signed replay, tax/accountant approval
Acceptance: cash/card/split/fail/cancel/refund/settlement independently reconcile

## Review and rollback
Reviewers: domain owner, accounting reviewer when financial policy/effects change, security reviewer for access changes, DB reviewer for migrations. Rollback: Provider effects require reconciliation/compensation; no blind reversal.
Merge gate: all acceptance evidence reviewed, no cross-tenant leakage, no unresolved financial discrepancy.
