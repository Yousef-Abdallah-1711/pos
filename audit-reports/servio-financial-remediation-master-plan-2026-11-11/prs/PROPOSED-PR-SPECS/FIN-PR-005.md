# FIN-PR-005 — Durable Accounting Event Processing

**Status: PROPOSED PR specification; not created. Risk: High.**

## Objective and task IDs
FIN-T03-01, FIN-T03-02

## Scope and affected files
DISCOVERY REQUIRED: source write paths, outbox/inbox, worker and migrations

## Dependencies and schema
Requires dependencies in remediation-tasks.json and approved policies. Schema impact must be confirmed against deployed baseline; no migration is presumed.

## Tests and acceptance
Failure injection, PostgreSQL transaction proof, queue alert/runbook
Acceptance: crash, replay, duplicate, concurrent event produces one effect

## Review and rollback
Reviewers: domain owner, accounting reviewer when financial policy/effects change, security reviewer for access changes, DB reviewer for migrations. Rollback: Do not enable async path until reconciliation and dual-read parity proven.
Merge gate: all acceptance evidence reviewed, no cross-tenant leakage, no unresolved financial discrepancy.
