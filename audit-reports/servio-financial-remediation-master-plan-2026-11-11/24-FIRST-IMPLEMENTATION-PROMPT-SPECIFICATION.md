# 24 — First implementation prompt specification

**Status: PROPOSED; do not execute from this planning package.** The future first implementation is a narrowly scoped FIN-001 reporting safety and correctness change. It is not permission to start coding.

## Objective and prerequisite approvals
After an authorized owner/accountant approves the fiscal period and retained/current earnings presentation rule, prevent the tested Balance Sheet equation failure. Before work, refresh the source baseline and actual deployed schema, confirm the FIN-001 function remains present, and review any in-scope changes. If the accounting policy or baseline differs, stop and re-scope.

## Allowed scope
- Isolated worktree/branch from the explicitly selected clean implementation base.
- Confirmed Balance Sheet calculation in `backend/internal/usecase/financial_report_usecase.go`.
- Focused unit/integration fixtures and only necessary report contract types.
- Any migration only if discovery proves it is required and it is reviewed independently.

Forbidden: source/order/payment logic, auto-posting, provider behavior, tax rate selection, broad refactor, production/shared DB, existing dirty worktrees, destructive backfill, unrelated formatting. No commit/push/merge/deploy until separately authorized.

## Required tests and evidence
- Reproduce four historical dated fixtures plus independently calculated opening equity, current earnings and prior-period close boundaries.
- Real PostgreSQL integration against isolated disposable database; assert assets = liabilities + equity and reconcile P&L/Trial Balance at matching cutoff.
- Include tenant/restaurant scope and report authorization regression.
- Capture before/after source SHA, exact query/function, fixtures, DB schema version, test output, report equation results, and reviewer/accountant approval.

## Branch, migration and recovery
Create a clean isolated worktree based on explicitly confirmed current main; do not use the dirty primary checkout. Make no migration unless required. If a migration is proposed: backup and restore rehearsal, forward migration on a copy, row/count/sum validation, lock analysis, rollback limits and forward-correction plan. Never update/delete posted history. If results disagree, stop and report; do not plug the difference.

## Acceptance and stop conditions
Acceptance: all fixtures reconcile; no current or prior period double-counting; report outputs agree with TB/P&L under approved policy; authorization tests pass; no unrelated files change. Stop on policy ambiguity, live-schema uncertainty, unexpected source divergence, any tenant leak, failing historical period, or need for data repair beyond the approved scope. Future PR must contain exact evidence and receive engineering plus accounting review.