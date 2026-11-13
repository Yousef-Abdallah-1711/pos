# FIN-PR-001 — Financial Statement Safety and Reconciliation Controls

**Status: PROPOSED PR specification; not created. Risk: High.**

## Objective and task IDs
FIN-T00-01, FIN-T06-01

## Scope and affected files
backend/internal/usecase/financial_report_usecase.go; confirm report routes in current base

## Dependencies and schema
Requires dependencies in remediation-tasks.json and approved policies. Schema impact must be confirmed against deployed baseline; no migration is presumed.

## Tests and acceptance
Finance/accountant review required; independent expected Balance Sheet/P&L/TB; no schema unless discovery requires it
Acceptance: dated report equation and cutoff fixtures

## Review and rollback
Reviewers: domain owner, accounting reviewer when financial policy/effects change, security reviewer for access changes, DB reviewer for migrations. Rollback: Rollback report availability only to safe warning state; never revert to knowingly misleading statement behavior.
Merge gate: all acceptance evidence reviewed, no cross-tenant leakage, no unresolved financial discrepancy.
