# Financial forensic audit — workflow coverage matrix

**Overall: PARTIAL. Commercial readiness: NO-GO pending financial control remediation.** Evidence is for baseline `e519f52eb0badfb80c283f79a52c51acc1c663c6` and disposable tenant/database only. Statuses describe demonstrated evidence, not assumptions. `BLOCKED` means a specific integration defect prevents that path; `NOT TESTED` means no scenario was executed. No application source changes were made.

## Financial workflows

| Workflow / assertion | Status | Evidence and limits |
|---|---|---|
| Manual journal default UI submission | FAILED | Earlier browser/API evidence: omitted attachments produces PostgreSQL invalid JSON; explicit `attachments: []` creates draft and post succeeds. |
| Manual journal balanced posting / trial balance | VERIFIED | Original Dr/Cr EGP 100 journal posted; two ledger rows; trial balance 100/100. One fixture only. |
| FIN-001 Balance Sheet equation | FAILED | Four fixture dates; missing earnings/equity roll-up confirmed in `financial_report_usecase.go:261-336`; gap equals net profit (EGP 100 first fixture, then EGP 60). |
| P&L for fixture periods | VERIFIED (limited) | Compared net earnings across four controlled fixture dates. No full chart/classification or cutoff matrix. |
| Cash flow | PARTIAL | Earlier one-case response matched cash debit directionally; no independent operating/investing/financing reconciliation. |
| POS cash order | NOT TESTED | No paid POS order source→payment→journal→ledger trace. |
| POS card order | NOT TESTED | No gateway/payment fixture. |
| Website order | NOT TESTED | Public order implementation inspected; no runtime fixture through report/export. |
| Split payment | NOT TESTED | No test. |
| Failed payment | NOT TESTED | No test. |
| Canceled order | NOT TESTED | No test. |
| Partial refund | NOT TESTED | No test. |
| Full refund / reversal | NOT TESTED | No test. |
| Discount + tax combinations | NOT TESTED | Runtime combinations absent. Fixed 10% tax is source-verified at `order_usecase.go:173`; policy applicability unknown. |
| Inventory consumption and COGS | NOT TESTED | Posting source exists; no stock/order fixture reconciliation. |
| Supplier purchase / payable via automatic entry | BLOCKED | Actual auto call failed before persistence on nil JSONB attachments; with valid JSON, `auto_supplier_invoice` exceeds `VARCHAR(20)`. Other supplier repository routes not fully traced. |
| FIN-003 material imbalance normalizer | VERIFIED behavior; auto-path persistence BLOCKED | Dr 100 / Cr 90 becomes final credit 100; conditional repository route posted normalized 100/100 with two ledger rows. Actual supplier automatic method failed earlier at attachments. Reverse Dr100/Cr110 rejected. |
| Customer invoice / payment | NOT TESTED | No runtime invoice/payment trace. |
| Receivables / payables aging | NOT TESTED | Existing empty/insufficient snapshots do not reconcile amounts. |
| VAT report | NOT TESTED | No meaningful tax fixture or independent reconciliation. |
| Inventory valuation | NOT TESTED | No valuation fixture. |
| Financial PDF export | PARTIAL | Existing PDF download evidence; row count, exact amounts, filters, tenant boundary not verified. |
| CSV export | NOT TESTED | No runtime export evidence. |
| XLSX export | NOT TESTED | No runtime export evidence. |
| Spreadsheet formula safety | NOT TESTED | No CSV/XLSX output inspected. |

## Reliability, tenancy and performance

| Area | Status | Evidence and limits |
|---|---|---|
| Automatic posting JSONB insert | FAILED | Actual `PostSupplierInvoiceEntry` probe failed with invalid JSON; no journal persisted. Source `createAndPostEntry` omits attachments. |
| Supplier `entry_type` schema compatibility | FAILED | 21-character `auto_supplier_invoice` rejected against `VARCHAR(20)` when valid JSON supplied. |
| Paid-order auto-post error propagation | VERIFIED at source; runtime NOT TESTED | Sales/COGS errors are logged and suppressed (`order_usecase.go:492-520`); no actual paid-order scenario. |
| Posting transaction: ledger + status atomicity | SOURCE VERIFIED, failure injection NOT TESTED | Inspected repository transaction; no injected crash/retry. Draft creation remains separate. |
| Duplicate requests / idempotency | NOT TESTED | No repeated order/payment/webhook/post requests. |
| Concurrent payments / refunds / posting | NOT TESTED | No race/concurrency scenario. |
| Retry / crash recovery / failure injection | NOT TESTED | No injected failure. |
| Two-tenant financial isolation | NOT TESTED | Only the synthetic audit tenant was used in this continuation. |
| Cross-tenant object/export access | NOT TESTED | No adversarial IDs or export attempts. |
| Low-privilege RBAC | NOT TESTED | Synthetic owner context only; no role-negative matrix. |
| Development performance | PARTIAL / HISTORICAL | Earlier one sample per route, about 2.7–10.6 s, duplicate requests. Not repeatable distributions and no root-cause attribution. |
| Production-mode performance | NOT TESTED | No production build measurements. |
| Browser traces, waterfalls, query plans, React profiling | NOT TESTED | None captured in this continuation. |

## Confirmed root causes and required action

- **FIN-001 (P1):** Balance Sheet filters to asset/liability/equity accounts and sums latest ledger balances; no current earnings/retained earnings computation. Remediation must follow approved close policy and reconcile independent multi-period fixtures.
- **FIN-003 (P1):** `normalizeBalancedLines` adjusts the final business line for any imbalance beyond 0.009. Reject material variance; only permit an approved, explicit, bounded rounding line.
- **FIN-008 (P2):** automatic journal construction leaves JSONB attachments nil. Set a valid empty JSON value and test every auto-post event type.
- **FIN-009 (P2):** automatic supplier entry type exceeds DB column width. Align allowed type/column through reviewed schema change and migration tests.
- **FIN-010 (P2):** paid-order accounting failures are logged but not surfaced/queued; order/payment can report success independently of accounting. Add durable retry/alert/reconciliation and exercise failing posting end-to-end.

## Git and environment

- Audited HEAD: `e519f52eb0badfb80c283f79a52c51acc1c663c6` (detached).
- Freshly fetched `origin/main`: same SHA; no upstream change relative to audited baseline.
- Worktree status: only `audit-reports/financial-forensic-audit-2026-11-10/` is untracked; no tracked application changes.
- Disposable Compose project: `audit-pos-financial-20261110`; DB `pos_phase2a_e2e`. No protected DB writes.
- Evidence: sibling `evidence/scripts/audit_closure_probe_test.go`, `evidence/scripts/overlay.json`, screenshots/snapshots, and reports 05–08, 12, 14, 17.
- Not complete: end-to-end payment/order/refund workflows, concurrency/idempotency, tenant/RBAC, production performance, and complete export/reconciliation matrix. This package does not claim complete coverage.

