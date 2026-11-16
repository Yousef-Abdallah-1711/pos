# SERVIO Phase 4 Financial Engineering — CSV Safety and Contract

**Assessment date:** 2026-11-16 (Africa/Cairo)
**Starting `origin/main`:** `e4333270d3057782dea2e6214f18f9c63e0fe8c9`
**Historical Phase 3E evidence:** [closure report](../financial-acceptance-phase3b-2026-11-15/PHASE3E_CLOSURE.md)

## Scope and disposition

This increment addresses two verified CSV export gaps without changing accounting calculations or policy:

- Chart of Accounts and journal-entry exports/imports now have truthful `/api/v1/accounting/csv/...` routes. The dashboard uses these routes. Existing `/excel/...` paths remain permission-protected compatibility aliases and continue returning/accepting CSV; no XLSX support is implied.
- Spreadsheet formula-leading text is now neutralized in Chart of Accounts and journal-entry CSV text columns, using the same shared rule as report exports. Numeric debit and credit cells retain their numeric representation.
- An authenticated PostgreSQL acceptance assertion covers CSV media type, `.csv` disposition, parseable output, the new paths, and compatibility paths. It runs inside the existing isolated financial acceptance workflow.

## Verification

Focused local checks passed:

- `go test ./internal/pkg/csvutil ./internal/repository -count=1`
- `go test ./internal/handler/http -run '^(TestWriteCSVAttachmentAdvertisesCSVContract|TestCSVSafeTextCellPreventsSpreadsheetFormulaInjection|TestReportCSVExportsEscapeUntrustedNames)$' -count=1`
- `go test ./cmd/api -run '^TestRegisteredRouteManifestParity$' -count=1`
- `go test -tags integration ./cmd/api -run '^TestFinancialAcceptanceLiveAPIOrderToLedger$' -count=1` (compilation/package check only; without the acceptance API URL the integration scenario skips)
- `git diff --check`

The focused dashboard ESLint check passed for the changed page. The local Docker engine is unavailable (`npipe:////./pipe/dockerDesktopLinuxEngine` does not exist), so no local disposable PostgreSQL/API run or authenticated financial Playwright session was executed. Dashboard-wide typecheck failed on pre-existing errors in customer and supplier invoice pages plus a missing `@radix-ui/react-tabs` module. Dashboard Vitest did not execute tests: all 53 suites failed to load because `@testing-library/jest-dom/matchers` is missing. Both commands used local Node `24.16.0`, below the repository's declared `>=24.22.3 <25` range. The full ESLint run was stopped after the targeted changed-file lint passed. No public, staging, or production API was used.

The fetched current-main CI runs were green before this change: [financial acceptance 43777741651](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43777741651) and [verified CI 43777741716](https://github.com/Yousef-Abdallah-1711/pos/actions/runs/43777741716). They are baseline evidence, not validation of this branch.

## Finding tracker and coverage boundary

The 21 canonical CFR/FIN finding dispositions were reviewed against the current Phase 3E tracker. FIN-019 gains a local implementation and focused tests for CSV formula safety and explicit CSV route contracts; it is **implemented locally, not yet hosted-CI verified or merged**. The existing 21-row disposition table remains the source for the other finding statuses; no evidence here upgrades a subsystem to fully accepted.

Search of tracked audit reports and available repository history found the Phase 3E note that the older 30-subsystem mapping is historical, but no authoritative 30-row source mapping. The mapping remains **UNAVAILABLE**. No missing rows or completion percentage have been inferred.

All ten policy decisions remain **PENDING APPROVAL**: CFR-011, OPEN-01, EARN-01, CLOSE-01, TAX-01, ROUND-01, REFUND-01, INV-01, PAY-01, and CF-01. No accounting semantics were changed. The existing 5,000-cent Balance Sheet residual remains visibly `UNRECONCILED`; no balancing entry was added. Cash Flow classification remains policy-blocked.

## Acceptance still open

- Authenticated financial browser coverage across EN/AR desktop and mobile, with accountant, cashier, and second-tenant sessions, remains **BLOCKED / NOT RUN** in this local environment. No screenshot, trace, or visual review is claimed.
- Financial production readiness remains **NO-GO** pending the existing accounting-policy, subsystem, and browser acceptance gates, irrespective of this scoped export correction.
