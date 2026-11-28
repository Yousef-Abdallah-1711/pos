# F-074 — Root cause

**Confidence:** CONFIRMED

## First break (from the register)
CONFIRMED (RUN-005 source + runtime): the parent customers page shows a toast and then rethrows the error (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx:95-118, throw error at line 112), and the child form awaits onSubmit without a try/catch (components/accounting/CustomerForm.tsx:69-79; same pattern in SupplierForm.tsx:75-85 and AccountForm.tsx:81-117), so every server rejection becomes an unhandled promise rejection; the contracts page uses try/finally with no catch (app/[locale]/dashboard/hr/contracts/page.tsx:18-32) and ContractForm awaits without catch (components/hr-advanced/ContractForm.tsx:27-29), so no message is shown at all. Runtime (Next dev build): all four forms raised window unhandledrejection "Request failed with status code 400"; customers, suppliers and accounts showed a toast with the raw database text "pq: value too long for type character varying(100)/(20)"; contracts showed no toast


## Origin and propagation
- Origin (first verified source location): `frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx`:129.
- Propagation: affected modules Customers, Suppliers, Chart of accounts, Contracts (dashboard forms); related findings F-014, F-025, F-065; observed consequence: POST returns 400 (for example pq: value too long for type character varying(100)); in every case the page raises an unhandled promise rejection "Request failed with status code 400"; the customer, supplier and account forms show a toast con
- Cause versus symptom: the symptom is what users or reports show (POST returns 400 (for example pq: value too long for type character varying(100)); in every case the page raises an unhandled promise reject); the cause is the first break above.

## Investigation performed
- Executed audit tests: EX-RUN004-FORMMATRIX, EX-RUN005-F074-UNHANDLED; scripts: tests/runtime/ph06_r4_form_matrix.mjs, tests/runtime/ph06_r5_f074.mjs.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
| reference | lines | resolved path | resolution | range valid |
|---|---|---|---|---|
| frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | 129 | frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx | 95 | frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx | exact | yes |
| frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | 69 | frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx | 18 | frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx | exact | yes |
| frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx | 27 | frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx | exact | yes |
| frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx | 95-118 | frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx | exact | yes |
| components/accounting/CustomerForm.tsx | 69-79 | frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx | prefix | yes |
| SupplierForm.tsx | 75-85 | frontend/apps/dashboard/src/components/accounting/SupplierForm.tsx | basename | yes |
| AccountForm.tsx | 81-117 | frontend/apps/dashboard/src/components/accounting/AccountForm.tsx | basename | yes |
| app/[locale]/dashboard/hr/contracts/page.tsx | 18-32 | frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx | prefix | yes |
| components/hr-advanced/ContractForm.tsx | 27-29 | frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx | prefix | yes |

