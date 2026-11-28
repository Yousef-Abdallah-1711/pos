"""RUN-005: F-071/F-020/F-033/F-032 evidence expansion from the authorization matrices (idempotent)."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
N = 'evidence/network/RUN-005/PH-17/'
U = {
    'F-071': (N + 'AUD-R5-RBAC-MATRIX-ENFORCE.json; ' + N + 'AUD-R5-RBAC-MATRIX-STAGING.json; tests/runtime/ph17_r5_rbac_matrix.mjs',
              ' RUN-005 expansion (4 roles x 28 GET endpoints + 8 write probes, same database, same users): enforce-mode API: ADMIN 27/28 GET 200, MANAGER 15/28, EMPLOYEE 10/28, DRIVER 3/28; write probes authorised 8, 2, 0, 0 of 8. ENV=staging API with RBAC unset: EVERY role (including DRIVER) gets 27/28 GET 200 with salary/bank/national-id fields in 4 responses and passes authorization on 8/8 write probes (400 body validation instead of 403). Real data was returned (sensitive-field responses 4 per role), not only status codes.'),
    'F-020': (N + 'AUD-R5-RBAC-MATRIX-ENFORCE.json',
              ' RUN-005: in enforce mode the EMPLOYEE (cashier) role gets 200 from GET /hr/employees with base_salary (5000) and email in the payload, and 200 from /hr/attendance and /hr/leaves; MANAGER additionally reads /hr/salaries (4 rows, salary fields) and the user and role lists (policy question POL-15, not asserted as a defect); DRIVER sees 0 orders (scoped) and 3/28 readable endpoints.'),
    'F-033': (N + 'AUD-R5-RBAC-MATRIX-ENFORCE.json', ' RUN-005: the same enforce-mode matrix shows the EMPLOYEE role reading HR employee records with salary fields; billing/subscription route GET /payments/subscription returned 404 for all roles in this build, so the billing part of the original claim is not re-verified.'),
    'F-032': ('evidence/network/RUN-005/PH-17/AUD-R5-OBSERVED-REPRO.json', ' RUN-005: GET http://localhost:13112/en/test-url returned HTTP 200 (curl) on the current isolated dashboard.'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
ix = {k: i for i, k in enumerate(hd)}
for r in rows[1:]:
    u = U.get(r[0])
    if u:
        for e in u[0].split('; '):
            if e not in r[ix['evidence_paths']]:
                r[ix['evidence_paths']] += '; ' + e
        if u[1].strip() not in r[ix['actual']]:
            r[ix['actual']] += u[1]
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    if 'EX-RUN005-RBAC-MATRIX' not in ex:
        csv.writer(h).writerow(['EX-RUN005-RBAC-MATRIX', 'API-DIAGNOSTIC', 'Authorization matrix: roles ADMIN/MANAGER/EMPLOYEE/DRIVER x 28 GET + 8 write probes on the enforce-mode API and on an ENV=staging API (actual returned data inspected)', 'admin/manager/employee/driver', 'FAIL',
                                'enforce: ADMIN 27, MANAGER 15, EMPLOYEE 10, DRIVER 3 of 28 GETs; staging: every role 27/28 GET and 8/8 write probes authorised, salary fields exposed to all (F-071); employee reads base_salary (F-020); evidence ' + N + 'AUD-R5-RBAC-MATRIX-*.json'])
print('ok')
