"""RUN-004 register update 4 (idempotent): F-074, F-075 from the real-UI form matrix + 34b rows. Audit tooling only."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
NET = 'evidence/network/RUN-004/PH-06/'; SH = 'evidence/screenshots/RUN-004/PH-06/'
F = {
    'F-074': dict(
        title='Several dashboard forms leave a rejected submit unhandled: a server 400 surfaces as "Runtime AxiosError / Request failed with status code 400" and the server message is never shown',
        category='Software defect (UI error handling)', severity='Low', evidence_status='REPRODUCED',
        affected_modules='Customers, Suppliers, Chart of accounts, Contracts (dashboard forms)',
        reproduction='Headed Chrome tenant 4 owner: submit each form with a 300-character text value (or the contract form with an over-long enum value); read the on-screen text',
        expected='A field-level or toast message built from the server error, and the form stays usable',
        actual='POST returns 400 (for example pq: value too long for type character varying(20)); the page shows "Runtime AxiosError" and "Request failed with status code 400" (the unhandled-rejection overlay in this build) with no field message; empty submits do show client-side required messages, so only server-side rejections are unhandled',
        monetary_impact_EGP='0',
        root_cause_first_break='PROBABLE (observed behaviour consistent across four forms; the submit handlers were not individually traced): the form submit promise is not caught, so axios rejections propagate to the framework overlay; related to F-025 (server returns raw DB text) and F-014 (generic toast)',
        source_refs='frontend/apps/dashboard/src/components/categories/CategoryFormPage.tsx:129',
        related_static_ids='F-014; F-025; F-065', evidence_paths=NET + 'AUD-R4-FORM-MATRIX.json; ' + SH + 'AUD-R4-FORM-customers-long-after.png; tests/runtime/ph06_r4_form_matrix.mjs',
        remediation_direction='Wrap every submit in try/catch, map server field errors to inline messages, show a generic toast for the rest; add a shared useFormSubmit hook',
        regression_tests_needed='For each listed form, a rejected submit shows a readable message and no framework error overlay'),
    'F-075': dict(
        title='Chart of accounts accepts any text as an account code or name (HTML and SQL characters, non-numeric codes) with no format or length rule in the UI or API',
        category='Software defect (validation)', severity='Low', evidence_status='REPRODUCED',
        affected_modules='Accounting > Chart of accounts',
        reproduction='Headed Chrome tenant 4: open the account form, enter account code and names "<b>x</b>\'; --", submit',
        expected='400 with a field error: codes follow the numeric/hierarchical format of the chart; names are plain text within the column limit',
        actual='POST /accounting/accounts returns 201 and stores account id 389 with code <b>x</b>\'; -- (the response HTML-escapes it as unicode); a 300-character value returns 400 with the raw database length error; no maxlength on the inputs',
        monetary_impact_EGP='0 direct; report columns, exports and code-prefix classification (for example the 12% and 22% cash-flow rules) depend on well-formed codes',
        root_cause_first_break='CONFIRMED (source + runtime): CreateAccount only trims the strings and passes them on (backend/internal/handler/http/account_handler.go:124-170: no regexp, length or format check); repository errors are returned to the client as 400 text (F-025)',
        source_refs='backend/internal/handler/http/account_handler.go:124',
        related_static_ids='F-031; F-062; F-003', evidence_paths=NET + 'AUD-R4-FORM-MATRIX.json; ' + SH + 'AUD-R4-FORM-accounts-invalid-after.png',
        remediation_direction='Validate code against the chart pattern (digits, length, parent prefix), trim and limit names, set input maxlength, escape in exports',
        regression_tests_needed='Invalid code or name returns 400 with a field message; valid code accepted'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
have = {r[0] for r in rows}
for fid, v in F.items():
    if fid not in have:
        v = dict(v, finding_id=fid)
        rows.append([v.get(k, '') for k in hd])
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [
    ('EX-RUN004-FORMMATRIX', 'UI-HEADED-CHROME', 'Real UI form matrix: 14 targets x 5 modes (empty, invalid, arabic, long, valid) = 70 runs; customers, suppliers, accounts, journal, attendance, contracts, users, drivers, tables, categories, products exercised; employees (route renders no form), leaves (no add control found) and roles (no add control found) NOT exercised', 'owner', 'PARTIAL',
     'valid and Arabic input persisted for customers, suppliers, accounts, drivers, products; empty submits show required messages; long values: 400 raw pq text (customers/suppliers/accounts/contracts, F-074) and 500 (drivers, tables, products, F-031/F-062/F-064); account code accepts HTML/SQL characters (F-075); categories with display order 400 (F-065 reconfirmed); journal counts in this matrix are distorted by the concurrent benchmark; evidence ' + NET + 'AUD-R4-FORM-MATRIX.json'),
]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    w = csv.writer(h)
    for r in new:
        if r[0] not in ex:
            w.writerow(r)
print(len(rows) - 1, 'findings')
