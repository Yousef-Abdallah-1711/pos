# FIN-PR-008 — Financial Reports and Period Close

**Status: PROPOSED PR specification; not created. Risk: High.**

## Objective and task IDs
FIN-T06-02

## Scope and affected files
DISCOVERY REQUIRED: all report queries and close paths

## Dependencies and schema
Requires dependencies in remediation-tasks.json and approved policies. Schema impact must be confirmed against deployed baseline; no migration is presumed.

## Tests and acceptance
Accountant signoff, comparative/close/backdate fixtures
Acceptance: all statements reconcile at cutoff and watermark

## Review and rollback
Reviewers: domain owner, accounting reviewer when financial policy/effects change, security reviewer for access changes, DB reviewer for migrations. Rollback: Version/restate reports; never delete posted entries.
Merge gate: all acceptance evidence reviewed, no cross-tenant leakage, no unresolved financial discrepancy.
