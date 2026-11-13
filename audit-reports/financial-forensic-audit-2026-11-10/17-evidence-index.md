# Evidence index

Baseline e519f52eb0badfb80c283f79a52c51acc1c663c6. Runtime: isolated Compose project audit-pos-financial-20261110, DB pos_phase2a_e2e, API 18082, Dashboard 13112, Landing 13111. Synthetic password/token omitted.

## Browser evidence

- evidence/registration-result.yml and register-page-snapshot.yml: registration.
- evidence/snapshots/: route snapshots; most reflect empty tenant.
- evidence/screenshots/dashboard-en-desktop.png: English dashboard.
- evidence/screenshots/journal-entry-ar-desktop.png and journal-entry-ar-mobile.png: Arabic detail at desktop and 390x844.
- evidence/screenshots/profit-loss-en-desktop.png: populated P&L.
- evidence/snapshots/profit-loss-2026-10-12-to-2026-11-10.pdf: downloaded report PDF.
- evidence/snapshots/page-*.yml and console-*.log: recent interaction snapshots and console capture.

## Controlled data evidence

- Browser request 736 showed journal UI POST without attachments and HTTP 400.
- Response body: pq invalid input syntax for type json.
- Authenticated browser fetch with attachments: [] returned 201 and draft id 1; UI post changed status to posted.
- Read-only psql output summarized in report 05. Journal lines and ledger rows match.
- API report responses summarized in report 08.
- Unbalanced 100/90 returned HTTP 400 journal entry is not balanced; one-line test returned HTTP 400 requiring at least two lines.

## Static references

- Journal schema: backend/migrations/137_create_journal_entries.sql.
- Ledger schema: backend/migrations/154_create_general_ledger.sql.
- Indexes: backend/migrations/157_optimize_financial_reporting_indexes.sql.
- Routes: backend/cmd/api/main.go:1051-1058, 1405-1416.
- Request/repository: backend/internal/domain/accounting.go; backend/internal/handler/http/journal_entry_handler.go; backend/internal/repository/journal_entry_repo.go.
- Posting: backend/internal/usecase/auto_posting_engine.go:522 onward and normalizeBalancedLines at 633 onward.
- Tax: backend/internal/usecase/order_usecase.go:173.
- Export: backend/internal/handler/http/accounting_advanced_handler.go:985 onward; Dashboard PDF component in frontend/apps/dashboard/src/components/financial-reports/ReportExportButton.tsx.

## External references

- IFRS IAS 1: https://www.ifrs.org/issued-standards/list-of-standards/ias-1-presentation-of-financial-statements/
- IFRS IAS 2: https://www.ifrs.org/issued-standards/list-of-standards/ias-2-inventories/
- IFRS 15: https://www.ifrs.org/issued-standards/list-of-standards/ifrs-15-revenue-from-contracts-with-customers/
- PostgreSQL 15 transaction isolation: https://www.postgresql.org/docs/15/transaction-iso.html
- PostgreSQL 15 constraints: https://www.postgresql.org/docs/15/ddl.html
- OWASP API BOLA test: https://wstg.owasp.org/latest/4-Web_Application_Security_Testing/12-API_Testing/02-API_Broken_Object_Level_Authorization/

Repository/runtime evidence supports observed behavior; external standards do not certify this product.

## Continuation evidence index — 2026-11-11

- **Source root cause:** `backend/internal/usecase/financial_report_usecase.go:261-336` Balance Sheet query and equity totals; excludes income/expense and calculates no current earnings/retained earnings.
- **Fixture summary:** four sequential posted accounting fixture states and P&L comparison are recorded in `08-financial-report-validation.md` and `07-accounting-correctness.md`.
- **FIN-003 and automatic posting:** `evidence/scripts/audit_closure_probe_test.go`; overlay map `evidence/scripts/overlay.json`. Tests actual supplier method failure from nil JSONB attachments, normalizer EGP 10 adjustment, conditional repository persistence, reverse-direction validation, and `entry_type` width failure.
- **Supplier type schema:** `backend/internal/domain/auto_posting.go` constant `auto_supplier_invoice`; `backend/migrations/137_create_journal_entries.sql` defines `entry_type VARCHAR(20)`.
- **Order error handling:** `backend/internal/usecase/order_usecase.go:492-520`; `backend/internal/handler/http/admin_order_handler.go:500-507`; `backend/internal/handler/http/public_order_handler.go:82-84`. Error logging/suppression is source evidence; paid-order runtime outcome NOT TESTED.
- **Git:** audited worktree HEAD and freshly fetched `origin/main` are both `e519f52eb0badfb80c283f79a52c51acc1c663c6`. `git status --short` shows only the untracked audit report directory; no tracked source changes.
- **Coverage closure:** `19-coverage-matrix.md` is the consolidated matrix; status details and non-tested scope are explicit.

Disposable runtime remains `audit-pos-financial-20261110` / database `pos_phase2a_e2e` during this evidence closure. No production/shared DB was modified. Prior primary-checkout Playwright artifacts could not be removed due to execution-policy rejection and remain outside the audit worktree; they are not used as new evidence here.

## Execution environment limitation

The previously executed overlay probe is represented by `evidence/scripts/audit_closure_probe_test.go` and `overlay.json`. A repeat run was attempted during this closure, but the host has no `go` executable (`go: The term 'go' is not recognized`) and the isolated API container does not expose a Go toolchain/source mount. No replacement run was claimed and no test result log was fabricated. The earlier runtime results in this report are retained as historical evidence from the preceding audit turn; exact raw console output was unavailable in this continuation.
