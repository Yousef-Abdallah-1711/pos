import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
NET = 'evidence/network/RUN-004/PH-06/'
F = {
    'F-079': dict(
        title='Leave requests cannot be created through the UI: the form sends ISO date-times and the API accepts only YYYY-MM-DD (400)',
        category='Software defect (UI/API contract)', severity='Medium', evidence_status='REPRODUCED',
        affected_modules='HR > Leave requests (dashboard form, API)',
        reproduction='Headed Chrome tenant 4 owner: HR > Leaves > Request Leave, choose valid start and end dates, submit',
        expected='201 and the request appears in the list', actual='POST /hr/leaves returns 400 {"error":"Invalid start_date format (expected YYYY-MM-DD)"}; the page shows "Request failed with status code 400" for every valid submission; the request body carries start_date "2026-11-20T00:00:00.000Z"',
        monetary_impact_EGP='0 direct; blocks the leave workflow, so entitlement and payroll deductions cannot be exercised from the UI (see F-054)',
        root_cause_first_break='CONFIRMED (source + request capture): frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114-115 send new Date(...).toISOString(); backend/internal/handler/http/hr_leave_handler.go:186 parses the field as YYYY-MM-DD only',
        source_refs='frontend/apps/dashboard/src/components/hr/LeaveForm.tsx:114; backend/internal/handler/http/hr_leave_handler.go:186',
        related_static_ids='F-014; F-054; F-074', evidence_paths=NET + 'AUD-R4-FORM-MATRIX-employees_leaves_roles.json; tests/runtime/ph06_r4_form_matrix.mjs',
        remediation_direction='Send the date input value (YYYY-MM-DD) unchanged, or make the API accept RFC3339 and normalise; add a contract test posting the form payload',
        regression_tests_needed='Submitting the leave form with valid dates returns 201 and lists the request'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
have = {r[0] for r in rows}
for fid, v in F.items():
    if fid not in have:
        rows.append([dict(v, finding_id=fid).get(k, '') for k in hd])
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [('EX-RUN004-FORMMATRIX2', 'UI-HEADED-CHROME', 'Real UI form matrix part 2: employees, leaves, roles x 5 modes (15 runs)', 'owner', 'FAIL',
        'employees valid -> 500 positive_working_hours (F-014 reconfirmed on the current tree); empty employee submit sends nothing and shows no message; leave requests always 400 date format (F-079); roles accept HTML/SQL characters in the name (201) and an over-long name returns 500 (F-062/F-031 class); evidence ' + NET + 'AUD-R4-FORM-MATRIX-employees_leaves_roles.json')]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    for r in new:
        if r[0] not in ex:
            csv.writer(h).writerow(r)
print(len(rows) - 1)
