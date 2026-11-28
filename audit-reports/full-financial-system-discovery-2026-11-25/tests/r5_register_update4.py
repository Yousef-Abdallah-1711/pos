"""RUN-005 register update 4: F-080, F-081 (new) and F-024/F-007/F-008 evidence from the real-UI invoice run (idempotent)."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
N9 = 'evidence/network/RUN-005/PH-09/AUD-R5-INVOICE-FORM-RESULT.json'
F = {
    'F-080': dict(
        title='Customer and supplier forms offer header and unrelated accounts as the receivable/payable account, the backend accepts them, and later invoicing fails with "account N is a header account"',
        category='Software defect (validation / data integrity)', severity='Medium', evidence_status='REPRODUCED',
        affected_modules='Accounting > Customers, Suppliers, Customer invoices, Payments',
        reproduction='Headed Chrome tenant 4: open the customer form; the Receivable account dropdown lists 1000 - Assets, 1100 - Current Assets, Petty Cash, Cash on Hand, Bank accounts and others; choose the first account, save; then create a customer invoice for that customer (the form matrix did this for customers 6-13 and suppliers 6, 7)',
        expected='Only postable receivable (or payable) accounts are offered, and the API rejects a header or non-receivable account',
        actual='Customer rows 6-13 were stored with receivable_account_id 63 (1000 Assets, is_header true) and suppliers 6 and 7 with payable_account_id 64 (2000 Liabilities, header); the create call returned 201; creating an invoice for such a customer returned 400 {"error":"account 63 is a header account and cannot be posted to"}, so the customer can never be invoiced until the account is edited',
        monetary_impact_EGP='0 direct; blocks invoicing and payments for affected customers and suppliers',
        root_cause_first_break='CONFIRMED (source + runtime): CustomerForm.tsx:61-62 filters the dropdown only by account_type === "asset" (the supplier form uses the same pattern), the customer handler passes ReceivableAccountID straight to the repository (backend/internal/handler/http/customer_acct_handler.go:141,244) without checking that it is a postable receivable account, and the first posting attempt fails at backend/internal/repository/journal_entry_repo.go:708',
        source_refs='frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx:61; backend/internal/handler/http/customer_acct_handler.go:141; backend/internal/repository/journal_entry_repo.go:708',
        related_static_ids='F-062; F-075', evidence_paths=N9 + '; evidence/network/RUN-005/PH-09/AUD-R5-INVOICE-FORM-DUMP.json; tests/runtime/ph09_r5_invoice_form.mjs',
        remediation_direction='Offer only active, non-header accounts of the receivable (payable) subtype, validate the account server-side on create and update, and show a field error',
        regression_tests_needed='Creating a customer or supplier with a header or non-receivable account returns 400; the dropdown never lists header accounts'),
    'F-081': dict(
        title='Invoice, payment and related forms load only the first 20 customers or suppliers, so customers beyond the first page cannot be selected',
        category='Software defect (UI/API contract)', severity='Medium', evidence_status='REPRODUCED',
        affected_modules='Accounting > Customer invoices, Supplier invoices (customer and supplier dropdowns)',
        reproduction='Headed Chrome tenant 4 (91 active customers): open New customer invoice and count the options in the customer dropdown',
        expected='Every active customer can be selected (search or paging)',
        actual='The dropdown held 21 options (placeholder + 20 customers: -5, 0, 0.001, ...) while 91 customers are active; a customer created for the test could only be chosen after renaming it so that it sorted into the first 20',
        monetary_impact_EGP='0 direct; customers beyond the first 20 (alphabetically) cannot be invoiced through the UI',
        root_cause_first_break='CONFIRMED (source + runtime): the customer list endpoint defaults to page_size 20 (backend/internal/handler/http/customer_acct_handler.go:29; suppliers: supplier_handler.go:29) and the invoice pages call getCustomerAccounts() / getSuppliers() without paging or search (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx:120; supplier-invoices/page.tsx:107)',
        source_refs='backend/internal/handler/http/customer_acct_handler.go:29; backend/internal/handler/http/supplier_handler.go:29; frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices/page.tsx:120; frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices/page.tsx:107',
        related_static_ids='F-080', evidence_paths=N9 + '; tests/runtime/ph09_r5_invoice_form.mjs',
        remediation_direction='Replace the plain select with a searchable async combobox (server-side search) or request all pages; keep the default page size for list screens only',
        regression_tests_needed='With more than 20 customers every customer can be found and selected in the invoice and receipt forms'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
ix = {k: i for i, k in enumerate(hd)}
have = {r[0] for r in rows}
for fid, v in F.items():
    if fid not in have:
        rows.append([dict(v, finding_id=fid).get(k, '') for k in hd])
U = {'F-024': ('REPRODUCED (RUN-005 headed Chrome)', ' RUN-005 (real invoice form, 1 x 5,000, 10% discount, 14% VAT): the POST payload carries "tax_amount":630.0000000000001 (raw float artefact) while the UI shows Tax EGP 630.00.'),
     'F-007': (None, ' RUN-005 real-UI reconfirmation: same invoice: UI preview tax 630.00 / total 5,130.00, stored tax 700.00 / total 5,200.00 (tax on the pre-discount 5,000), journal Dr 1120 5,200 / Cr 4100 4,500 / Cr 2120 700.'),
     'F-008': (None, ' RUN-005 real-UI reconfirmation: draft invoice INV-2026-00005 was posted to journal entry 44787 at creation.')}
for r in rows[1:]:
    u = U.get(r[0])
    if u:
        if u[0]:
            r[ix['evidence_status']] = u[0]
        if N9 not in r[ix['evidence_paths']]:
            r[ix['evidence_paths']] += '; ' + N9
        if u[1].strip() not in r[ix['actual']]:
            r[ix['actual']] += u[1]
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    if 'EX-RUN005-INVOICE-UI' not in ex:
        csv.writer(h).writerow(['EX-RUN005-INVOICE-UI', 'UI-HEADED-CHROME', 'Customer and invoice forms through the real UI: valid customer with receivable 1120; invoice 1 x 5,000, 10% discount, 14% VAT; dropdown counts; header-account failure', 'owner', 'FAIL',
                                'customer created (id with AR 74); invoice 201 with payload tax 630.0000000000001 (F-024), stored tax 700.00 vs UI 630.00 (F-007), draft journal posted (F-008); customer dropdown 20 of 91 (F-081); header receivable account accepted then invoice 400 (F-080); evidence ' + N9])
print(len(rows) - 1, 'findings')
