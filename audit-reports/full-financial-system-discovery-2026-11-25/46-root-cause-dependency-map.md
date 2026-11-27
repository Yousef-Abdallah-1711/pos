# 46 — Root-cause and dependency map

Chains traced from UI to ledger; "first break" = first incorrect or missing step. Findings from [22](22-confirmed-findings-register.csv).

| Root cause (foundational) | First break | Symptoms (findings) |
|---|---|---|
| RC-1 No tenant accounting bootstrap | registration transaction (`auth_repo.go`) | F-001 (no posting), F-011 partly (no opening journals), F-022 (storefront setup also absent), downstream every report is empty until fixed |
| RC-2 No single money/tax contract between client and server | POS store / invoice form compute independently of server | F-004, F-007, F-024, F-012 (no common definitions) |
| RC-3 No tender/payment model for orders | POS marks paid only for cash; no gateway; clearing never cleared | F-005, F-026, F-017 (duplicate unpaid orders consume stock) |
| RC-4 Inventory not tied to accounting | stock via product form and order creation; no receiving/adjustment postings | F-011, F-021 (supplier balances), COGS timing |
| RC-5 Report layer queries bespoke sources | each report has its own SQL/filters | F-002, F-003, F-010, F-012, F-009 (missing column) |
| RC-6 Concurrency not enforced at DB level for customer settlement | balance check outside serialized unit | F-006 |
| RC-7 Period and approval lifecycle absent | void path skips check; no period API; no payroll approval gate | F-013, F-015 |
| RC-8 Contract mismatches / no API-UI contract tests | EmployeeForm payload keys | F-014, F-016 |
| RC-9 Dependency/health design | session store without timeouts; health ignores Redis | F-018 |
| RC-10 Weak tenant/privilege defense in depth | RLS on 20/70 tables; no immutability on paid orders/payments; coarse default roles | F-019, F-020, F-025 |

## Dependency graph (what must precede what)
RC-1 → (enables meaningful testing of everything financial) → RC-2/RC-3 (what gets posted) → RC-4 → RC-5 (reports can only be correct once sources are correct) → RC-7 (lifecycle) .
RC-6, RC-8, RC-9, RC-10 are independent of the accounting model and can be handled in parallel.
Do not patch reports (F-002/003/010/012) before RC-1..RC-4 and the policy decisions (EARN-01, TAX-01, REV, INV-01) are settled; report-level patches would be rework.

---
## RUN-003 additions (F-040 … F-055 and re-confirmed roots)
| Root cause | First break (file:line) | Symptoms |
|---|---|---|
| RC-10 Discount/tax trigger formula computes tax on gross | `migrations/141_create_invoices_customer.sql:79-83`, `142_create_invoices_supplier.sql:75-79` | F-007, F-041 |
| RC-11 No idempotency layer for money-creating POSTs | only `order_cancellation_handler.go:36` reads Idempotency-Key | F-017, F-052, F-047 (no period key) |
| RC-12 Check-then-write without row lock on customer settlement | `payment_acct_repo.go:445-449`, `:569-582` | F-006 |
| RC-13 Reports read operational tables instead of the GL | `order_repo.go:859-870`, `accounting_advanced_repo.go:838-880,1105-1175`, `tax_repo.go:260-291`, `financial_report_usecase.go:285-306,366-399,488-526` | F-002, F-003, F-010, F-011, F-042, F-043, F-045 |
| RC-14 Inventory accounting keyed on `track_inventory` and legacy table | `auto_posting_engine.go:173`, `report_handler.go:557-645` | F-040, F-046, F-034, F-011 |
| RC-15 Paid orders are terminal; no reversal workflow | `order_usecase.go:341`, `order_repo.go:634`, `domain/order.go:215-240` | F-049, F-050 |
| RC-16 Error mapping returns 500 + raw DB text | many handlers (`admin_order_handler.go:410`, HR/accounting handlers) | F-025, F-031, F-050, HR matrix errors |
| RC-17 Environment/config coupling | CORS allow-list lacks wildcard in isolated compose (`middleware/cors.go:66-92` supports it) | F-022 (UI symptom) |
Shared-file clusters (fix together): `payment_acct_repo.go` (F-006, F-052), `accounting_advanced_repo.go` (F-009, F-043, F-045, F-047), `financial_report_usecase.go` (F-002, F-003, F-011), `auto_posting_engine.go` (F-040, F-015), `order_usecase.go` (F-017, F-028, F-049).
