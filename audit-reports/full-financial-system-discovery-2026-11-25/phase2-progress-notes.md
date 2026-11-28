# Phase 2 progress notes (architecture, API and page inventory)

## Coverage (stated honestly)
- API inventory (`08-api-inventory.csv`): 374 routes, parsed programmatically (tests/gen_api_inventory.py + tests/build_api_csv.py) from `backend/cmd/api/*.go`; count equals `routes.manifest` (374). 4 commented-out routes excluded. Handler file:line resolved for all but inline closures (39, pointed at registration line). Module names for IDs 1-12 come from the frontend `MODULE_IDS` (lib/rbac.ts:21-33); backend migrations seed IDs 13-20 explicitly (126,132,155) but 1-12 are BIGSERIAL insertion order - mapping for 1-12 is UNVERIFIED against a live DB. Financial relevance was classified by rule + hand-enriched for ~35 routes; remaining "reads money"/"writes money" labels are rule-based and may be coarse.
- Page inventory (`06`): 119 pages/route handlers across 4 apps from a static import-closure scan; endpoints are literal-string matches only. Dynamic URL construction and server components fetching through helper wrappers may be missed. Restaurant-website `middleware.ts` not audited.
- Field inventory (`07`): 200 field rows; strongest coverage for accounting (journal, supplier/customer invoice, payments received/made, bank accounts, COA, suppliers, customers, advanced forms), payroll SalaryForm, products, inventory adjust, POS, storefront checkout. NOT covered or only noted: EmployeeForm (506 lines), contracts/insurance/resignation forms, leave/attendance forms, super-admin forms, auth/registration, RBAC role forms, billing PaymentModal fields. Validation column reflects frontend code plus a few backend checks cited; backend request-struct validation was not enumerated field by field. Everything is `NOT TESTED`.
- Architecture (`03`, `05`, `13`, diagrams): based on reading main.go, accounting_db.go, order/outbox/posting/journal/GL/invoice/payment/payroll code, migrations 141, 180-190. Not executed: no DB, no tests, no runtime verification.

## Items to hand to later phases
1. Verify (runtime) that `pos_app` can SELECT ledger tables for reports after migration 182.
2. Confirm supplier-invoice line-total trigger exists (customer one verified).
3. Confirm sources of `GetProductCost` and COGS cost drift.
4. Search for any scheduler that runs recurring journal templates.
5. Check deployment env actually sets `ENV=production` and `ACCOUNTING_DB_*` for the production compose path (docker-compose.prod.yml lacks the latter).
6. Verify what `ExportOrders`, `reports/*` do with pool/tenant scoping (raw SQL in handlers).

## Helper scripts (read-only, in tests/)
`gen_api_inventory.py`, `build_api_csv.py`, `analyze_pages.py`, `build_pages_md.py`, `build_fields_csv.py`, `fields_src.txt`, intermediate JSON `_routes_raw.json`, `_pages_raw.json`. Run from the audit folder with `python tests/<script>`.
