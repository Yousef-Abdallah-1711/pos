# API contract inventory

**Status: PARTIAL.** Route registrations are in backend/cmd/api/main.go; request/response contracts are not fully OpenAPI-validated.

## Accounting, permission module 15

| Method | Route | Permission | Purpose |
|---|---|---|---|
| GET/POST | /api/v1/accounting/journal-entries | READ / WRITE | List/create |
| GET/PUT | /api/v1/accounting/journal-entries/{id} | READ / WRITE | Read/update |
| POST | /api/v1/accounting/journal-entries/{id}/post | WRITE | Post draft |
| POST | /api/v1/accounting/journal-entries/{id}/void | WRITE | Void |
| GET | /api/v1/accounting/trial-balance | READ | Trial balance; no Dashboard page found |
| GET | /api/v1/accounting/general-ledger | READ | Ledger |
| GET/POST/PUT | /api/v1/accounting/accounts... | route-specific | Chart of accounts |
| GET/POST/... | /api/v1/accounting/bank-accounts... | route-specific | Bank accounts |
| GET | /api/v1/accounting/excel/chart-of-accounts/export | EXPORT | Handler emits CSV |
| GET | /api/v1/accounting/excel/journal-entries/export | EXPORT | Handler emits CSV |
| GET | /api/v1/accounting/tax/vat-report | accounting read | VAT |

## Financial reports, module 17 READ

GET /api/v1/reports/financial/{profit-loss,balance-sheet,cash-flow,receivables-aging,payables-aging,inventory-valuation,trial-balance}; also predictive-cash-flow. Audit export is GET /api/v1/audit/export (EXPORT). POST /api/v1/reports/export uses module 8 EXPORT.

Financial report pages export PDF client-side. Route name excel is not proof of XLSX: chart/journal handlers call writeCSVAttachment. Negative permissions, pagination, complete parameter validation, and all response schemas remain unverified.
