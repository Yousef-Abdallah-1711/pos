# FIN-PR-003 — Automatic Journal Persistence and Schema Compatibility

**Status: PROPOSED PR specification; not created. Risk: High.**

## Objective and task IDs
FIN-T02-01, FIN-T02-03, FIN-T02-04

## Scope and affected files
journal UI/API/repository; auto_posting domain/usecase; migration 137 contract

## Dependencies and schema
Requires dependencies in remediation-tasks.json and approved policies. Schema impact must be confirmed against deployed baseline; no migration is presumed.

## Tests and acceptance
Migration review, isolated forward/rollback attempt, preserved posted records, schema integration
Acceptance: attachments default and all entry types insert

## Review and rollback
Reviewers: domain owner, accounting reviewer when financial policy/effects change, security reviewer for access changes, DB reviewer for migrations. Rollback: Prefer forward correction after deployment; do not truncate type/history.
Merge gate: all acceptance evidence reviewed, no cross-tenant leakage, no unresolved financial discrepancy.
