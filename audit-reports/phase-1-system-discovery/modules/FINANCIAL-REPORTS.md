# Reports and financial reports (dashboard) - static audit

Status: **Verified implemented** (pages + API wiring for all main statements); calculations are server-side except display aggregation and PDF generation; financial report pages are ungated client-side.

## Pages (15)
Operational: `reports` (index), `reports/overview`, `sales`, `profitability`, `inventory`, `hr` (gates REPORTS 3-4 refs each; custom i18n; RTL refs 5-7). Financial: `reports/financial` (index + KPIs), `profit-loss`, `balance-sheet`, `cash-flow`, `receivables-aging`, `payables-aging`, `inventory-valuation`, `vat`, `advanced` - **no MODULE_IDS gate in any** (next-intl + custom). Hooks `useReports.ts` (cache + polling `setInterval`, :81), `useFinancialReports.ts`; API `reportsApi.ts`, `financialReportsApi.ts`; components `components/reports`, `components/financial-reports` (VATReport, ReportExportButton -> PDF).

## Contract (all match main.go:1328-1351; RBAC module 8 operational, 17 financial)
/reports/{overview,sales,profitability,inventory}, POST /reports/export, /reports/financial/{profit-loss,balance-sheet,cash-flow,receivables-aging,payables-aging,inventory-valuation,trial-balance,kpis,menu-profitability,break-even,predictive-cash-flow}, /accounting/tax/{periods,vat-report}. Unused: `/reports/hr/social-insurance`, `/audit/export`.

## Where calculations happen
- Server: all statement totals (P&L, balance sheet, cash flow, aging, VAT) come from the API; pages only render (`vat/page.tsx:34-36` maps output/input/net VAT).
- Client: profitability chart rounding (`profitability/page.tsx:137-138`), top-N slices, recharts formatting, KPI cards; dashboard overview revenue tile = sum of one recent-orders page (`dashboard/page.tsx:366-368`) - NOT a report aggregate.
- Export: `components/financial-reports/ReportExportButton.tsx` + `lib/pdf-export.ts` render DOM with html2canvas into jsPDF (image-based, not authoritative; dynamic import). `POST /reports/export` (server export) is wired in `reportsApi.ts`.
- Currency hard-coded EGP in financial advanced/index (`reports/financial/advanced/page.tsx:47`, `financial/page.tsx:33`).
- Recharts v3 `Formatter` typing is a known TS-error source (FE-DASH-009; not re-verified, tsc unavailable).

## Findings
FE-DASH-014 (financial pages ungated; backend does enforce module 17), FE-DASH-018, FE-DASH-023 (6s timeout for exports), FE-DASH-009.

## Unverified
Numeric correctness of server statements (backend/quality auditors), tz/date-range boundaries (`date_from/date_to` strings).
