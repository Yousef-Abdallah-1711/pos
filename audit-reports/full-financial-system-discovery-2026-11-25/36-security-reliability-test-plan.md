# 36 — Security, Concurrency & Recovery Test Plan (NOT EXECUTED)

All tests: **status NOT EXECUTED**. Environment for all: isolated stack `pos-p2a-ws05-e2e` ([37](37-isolated-environment-execution-runbook.md)); `pos_app` role, `RBAC_ENFORCEMENT=enforce`, accounting pool `pos_accounting`. Tests marked (PROD) additionally need an environment with `ENV=production` (owner approval, runbook §12.3). Never run on `pos-dev` (superuser, BYPASSRLS, audit-only RBAC — results invalid). Linked IDs: VER-* = [31](31-independent-finding-verification.csv), KG-* = [29](29-knowledge-gap-register.md). Routes: [08](08-api-inventory.csv).

Common evidence: request/response JSON, SQL output, `pg_stat_activity` snapshot where relevant, container log excerpt, before/after journal + TB export. Common recovery: restore golden dump (runbook §7.9).

## A. Tenant isolation, roles, RLS
| ID | Action | Expected | Expected DB/financial state | Links |
|---|---|---|---|---|
| SEC-001 | Verify runtime DB identity: `select current_user, rolsuper, rolbypassrls` through the API's pool and `pg_stat_activity.usename` | `pos_app` (no super/BYPASSRLS); accounting writes under `pos_accounting` | — | VER-019, KG-ORC-004 |
| SEC-002 | Tenant A token calls GET on every tenant-scoped financial route (journals, orders, invoices, reports) using Tenant B IDs | 403/404, never B data | no change | RLS covers 20 tables only (57/70 tenant tables app-level) |
| SEC-003 | Tenant A token with `tenant_id`/`restaurant_id` overridden in body/query/header on POST routes | tenant taken from JWT; override ignored/rejected | no B rows touched | |
| SEC-004 | As `pos_app` via psql with `SET app.tenant_id` unset: `select count(*)` on each of the 20 RLS tables and on 5 core financial tables without RLS (journal_entries, general_ledger, orders, payments, customer_invoices) | RLS tables → 0 rows; non-RLS tables → **all rows visible (documents the gap)** | | VER-015, KG-SEC |
| SEC-005 | As `pos_app`: `UPDATE`/`DELETE` on `general_ledger`, `journal_entries`, `journal_entry_lines`; then `TRUNCATE` | UPDATE/DELETE rejected by cfr002 triggers; TRUNCATE result recorded (no trigger found) | rows unchanged | VER-015 |
| SEC-006 | As `pos_app`: `UPDATE`/`DELETE` on `orders`, `payments`, `salaries`, `customer_invoices` (paid rows) | Record outcome (expected: allowed = gap) | restore dump after | KG-SEC |
| SEC-007 | Attempt DDL (`CREATE TABLE`, `ALTER`) as `pos_app` and `pos_accounting` | denied | | |
| SEC-008 | `pos_accounting`: attempt to read/write non-accounting tables | only granted tables accessible | | |
| SEC-009 | User with no accounting module permission calls accounting/report/outbox-redrive routes (15 WRITE) | 403 | no journals | RBAC modules 1–20 |
| SEC-010 | Per module (20) × level (READ/WRITE/…): matrix of one READ-only user calling each write route | 403 for every write route | | route manifest |
| SEC-011 | Expired / tampered / wrong-signature / other-tenant JWT on financial routes | 401 | | |
| SEC-012 | Login rate limit (20/min per IP) and Origin check (login without `Origin`) | 429 / "invalid request origin" | | |
| SEC-013 | (PROD) `ENV=staging` vs `production`: compare guards (accounting pool fallback, stub login, RBAC mode) | staging weaker = confirmed gap | | VER-020, VER-021 |
| SEC-014 | Stub/mock login: attempt login with `password123` for stub-mode users in isolated stack | outcome recorded; any success = critical | | VER-022 |
| SEC-015 | Storefront public routes: unauthenticated access to admin/accounting routes | 401 | | |

## B. Duplicate and concurrent operations
| ID | Action | Expected | Expected state | Links |
|---|---|---|---|---|
| REL-001 | POST the same order creation twice quickly (same payload, same client) | Record outcome: expected two orders (no idempotency key) = gap | count orders | KG, VER-009 |
| REL-002 | Submit `PATCH payment-status=paid` 20× in parallel for one order | exactly 1 outbox row (UNIQUE) and exactly 1 journal | journal count = 1 | VER-009, VER-014 |
| REL-003 | Mark paid → paid again after posting | no second journal | | |
| REL-004 | Redrive same failed event concurrently (2×) | one requeue / one journal | | VER-017 |
| REL-005 | Post same manual journal twice (same idempotency key if any) | record behaviour | | VER-014 |
| REL-006 | Settle one customer invoice with 5 parallel receipts totalling > balance | total paid ≤ invoice total; no negative balance | | |
| REL-007 | Concurrent supplier-invoice approval + payment | consistent balance_due | | VER-023 |
| REL-008 | 50 parallel orders consuming the same limited stock item | stock never negative / oversell policy recorded | stock = initial − sold | reservations migration 180 |
| REL-009 | Concurrent stock adjustment + order consumption | final stock = arithmetic result | | |
| REL-010 | Report call (TB/BS) every 200 ms while 100 postings happen | each report snapshot balanced (debits=credits) | | |
| REL-011 | Run the k6 100-VU profile used earlier against isolated API with RBAC enforce | error rate, p95 recorded | | earlier load*.js |

## C. Failure and recovery
| ID | Fault | Expected | Recovery check |
|---|---|---|---|
| REL-020 | Make journal posting fail (e.g. COA row missing for rule account) then mark order paid | order paid, outbox event retried then terminal "failed"; no partial journal | add account, redrive → exactly one journal (VER-011/017) |
| REL-021 | Stop API during outbox processing (`docker stop api`) mid-batch; restart | claimed events recovered; no duplicate journals | outbox states consistent |
| REL-022 | Stop Postgres for 30 s during order+payment; restart | clean 5xx, no half-written order/payment; client retry works | |
| REL-023 | Break only accounting pool credentials (wrong password via a compose copy) | payroll/outbox posting fails visibly; operational pool unaffected; no fallback in non-dev | VER-019 |
| REL-024 | Force a transaction rollback (violate a constraint after first insert in a multi-step flow: supplier-invoice approval, payroll mark-paid) | all-or-nothing | counts unchanged |
| REL-025 | Client timeout then retry of a payment POST (network delay via `tc`/proxy or abort in Playwright) | record whether duplicate money record appears | |
| REL-026 | Redis down | rate limiting/cache behaviour recorded; money flows unaffected | |
| REL-027 | Worker restart while event `processing` | lease/timeout recovery | |

## D. Financial-integrity boundary tests
| ID | Action | Expected | Links |
|---|---|---|---|
| FIN-001 | Fiscal period: with 0 rows post a journal dated in a past month; then (if a route/SQL exists) create a closed period and post into it; void a posted entry in a closed period | document fail-open and void bypass | VER-013 |
| FIN-002 | Invalid amounts on every money field: negative, 0, 1e308, NaN, "abc", 0.001, 3 decimals, huge precision, locale commas | rejected or rounded per policy; never stored as NaN/Inf | ROUND-01 |
| FIN-003 | Negative rounding: refund/loss of −1.005, −1.234 through P&L/BS | compare with oracle; expect −1.23 vs implementation −1.22 | VER-004 |
| FIN-004 | Tax: same cart in POS (14%) and stored order (10%) | compare display vs DB vs journal | VER-003 |
| FIN-005 | Unbalanced manual journal (debit≠credit by 0.01) | rejected | |
| FIN-006 | Journal to header (non-postable) account / inactive account / other-tenant account | rejected | |
| FIN-007 | Backdated and future-dated entries | policy outcome recorded | |
| FIN-008 | Void a posted journal; check TB/BS/CF and GL | reversal consistent; CF treatment per CFR-010 | VER-002 |
| FIN-009 | Fresh tenant first paid order (no COA) | observe failure mode, no silent loss | VER-011 |
| FIN-010 | Card/storefront order never marked paid → confirm absent from revenue/ledger; reconcile operational vs GL | divergence quantified | VER-008, VER-010 |

## E. Execution notes
- Order: SEC-001 first (gate: if roles are wrong, all SEC/REL results are void), then A, B, D, C (fault injection last; restore dump between).
- Severity if a test fails: cross-tenant data exposure or ledger corruption = Critical; duplicate journal/money = High; message/UX = Medium.
- Counts: SEC 15 · REL 11 + 8 · FIN 10 = **44 tests prepared**, 0 executed.
