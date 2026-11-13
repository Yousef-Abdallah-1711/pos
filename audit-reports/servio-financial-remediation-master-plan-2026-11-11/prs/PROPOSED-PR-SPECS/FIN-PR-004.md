# FIN-PR-004 — Reject Unsafe Automatic Balancing

**Status: PROPOSED PR specification; not created. Risk: High.**

## Objective and task IDs
FIN-T02-02

## Scope and affected files
backend/internal/usecase/auto_posting_engine.go

## Dependencies and schema
Requires dependencies in remediation-tasks.json and approved policies. Schema impact must be confirmed against deployed baseline; no migration is presumed.

## Tests and acceptance
Caller inventory, independently calculated fixtures, all auto-post event regression
Acceptance: Dr100/Cr90 rejected; bounded approved rounding only

## Review and rollback
Reviewers: domain owner, accounting reviewer when financial policy/effects change, security reviewer for access changes, DB reviewer for migrations. Rollback: Stop rollout if known callers depend on normalization; explicit remediation before re-enable.
Merge gate: all acceptance evidence reviewed, no cross-tenant leakage, no unresolved financial discrepancy.
