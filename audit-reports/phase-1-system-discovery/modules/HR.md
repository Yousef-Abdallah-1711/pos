# HR (dashboard) - static audit

Status: **Verified implemented** for employees, attendance, leaves, payroll (salaries); **Partial** for advanced HR (contracts, insurance, performance, training, resignations, HR reports) - UI present, no page gating, thin pages.

## Pages (9 + report)
Core (custom i18n, gated, direct `api.*` calls): `hr/employees` (281 lines; 8 gate refs), `hr/attendance` (336), `hr/leaves` (442), `hr/payroll` (301). Advanced (next-intl, no gate in page, delegate to hooks): `hr/contracts` (75), `hr/insurance` (68), `hr/performance` (73), `hr/training` (68), `hr/resignations` (89). `reports/hr/page.tsx` (259, gate REPORTS). Components: `components/hr/*` (EmployeeForm/List, AttendanceForm, LeaveForm, SalaryForm), `components/hr-advanced/*`; hooks `useContracts`, `useInsurance`, `usePerformanceReviews`, `useTraining`, `useResignations`, `useHRReports`, `useHRAdvanced` (unreferenced).

## Contract (all used paths match main.go)
Core: `/hr/employees` CRUD, `/hr/attendance` (+`/clock-in`), `/hr/leaves` (+approve/reject), `/hr/salaries` (+`/{id}/mark-paid`). Advanced: `/hr/contracts` (+renew, expiring), `/hr/insurance` (+expiring), `/hr/performance` (+acknowledge, complete), `/hr/training` (+expiring), `/hr/resignations` (+acknowledge, complete, settlement). Reports `/reports/hr/{headcount,turnover,absenteeism,overtime,leave-balances,training-compliance,cost-analysis}`. RBAC modules 3,11,12,18,19,20 server-side (wrapWithPermission).
Unused backend: `/hr/roles` CRUD, GET `/hr/employees/{id}`, `/hr/attendance/{id}`, POST `/hr/attendance` (manual), `/hr/leaves/{id}`, `/hr/salaries/{id}`, `/hr/performance/{id}`, `/hr/performance/employee/{id}`, `/hr/resignations/{id}`, `/reports/hr/social-insurance`.

## Notes
- Payroll: page lists/creates salary rows and marks paid via API; the pay calculation (net, deductions, insurance) is not computed on the client in the page (SalaryForm values sent). Server-side posting to accounting via `autoPostingEngine` (main.go:336) - not verified here.
- `employees/page.tsx` defines `interface Employee { [key: string]: any }` and holds `debugInfo` state; type looseness is a source of implicit-any/TS errors (cf. FE-DASH-009).
- Advanced pages use `useTranslations('hrAdvanced.*')`; keys exist in en/ar (parity verified 1719 = 1719 keys).
- Sidebar gating: contracts/insurance/resignations use CONTRACTS (20) (`Sidebar.tsx:144-146`).

## Findings
FE-DASH-014 (5 advanced pages ungated), FE-DASH-009 (implicit any/typing), FE-DASH-016 (mixed i18n).

## Unverified
Statutory calculations (Egyptian social insurance) location, attendance geofencing, payroll period locking.
