# Phase 6 Production Deployment Checklist

**Module:** HR + Accounting System  
**Date:** March 3, 2026

## 1) Pre-deploy Gate

- [x] `go build ./cmd/api` passes
- [x] `go test ./internal/usecase -run Phase6 -count=1` passes
- [x] `go test ./internal/repository -run Phase6 -count=1` passes
- [x] `go test ./cmd/api -run Phase6Routes -count=1` passes
- [x] `pnpm --filter @pos-saas/dashboard lint` passes
- [x] `pnpm --filter @pos-saas/dashboard build` passes
- [x] DB migration dry-run on fresh database passes (`001..157`)

## 2) Data & Migration Safety

- [x] Database backup taken and restore tested
- [x] Migration order validated in staging
- [x] `schema_migrations` confirms accounting/HR range applied:
  - `136_create_chart_of_accounts.sql`
  - `...`
  - `157_optimize_financial_reporting_indexes.sql`
- [x] Rollback script available: `backend/scripts/rollback_phase6_accounting_hr.sql`

## 3) Security & Access

- [x] RBAC modules 15..20 exist in `modules`
- [x] Role mappings seeded for Admin/Manager/Employee
- [x] Key accounting and HR endpoints enforce `wrapWithPermission`
- [x] SQL-injection regression tests pass
- [x] Tenant isolation regression tests pass

## 4) Functional Smoke (Post-deploy)

- [x] Create account in chart of accounts
- [x] Create/post/void journal entry
- [x] Mark order paid -> sales + COGS auto-posted
- [x] Mark salary paid -> payroll auto-posted
- [x] Run financial reports: P&L, Balance Sheet, Cash Flow, VAT
- [x] Run HR reports dashboard
- [x] Verify ledger cursor pagination and period summaries

## 5) Performance Smoke

- [x] P&L query response under target on staging data
- [x] Balance Sheet query response under target
- [x] Ledger query response under target
- [x] Report caching active (5-minute TTL)

## 6) i18n/RTL

- [x] Arabic locale navigation for accounting/reports pages works
- [x] RTL layout verified on financial + HR reports pages
- [x] No missing critical labels in accounting ledger/report flows

## 7) Release Sign-off

- [x] Product owner sign-off
- [x] QA sign-off
- [x] Engineering sign-off
- [x] Release tag created


