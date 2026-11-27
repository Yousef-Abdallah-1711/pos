import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
NET = 'evidence/network/RUN-004/PH-13/'
F = {'F-076': dict(
    title='Super Admin Users list links every row to /dashboard/users/{id}, a route that does not exist (404)',
    category='Software defect (UI navigation)', severity='Low', evidence_status='REPRODUCED',
    affected_modules='Super Admin > Users', reproduction='Headed Chrome, Super Admin app (port 13114), Users page: click View on a user row',
    expected='A user detail page (or no link)', actual='Navigation to /dashboard/users/{id} returns the 404 page (7 consecutive rows in the control matrix logged a 404 resource error); only dashboard/tenants/[id] exists as a detail route',
    monetary_impact_EGP='0', root_cause_first_break='CONFIRMED (source): super-admin/src/app/dashboard/users/page.tsx:132 renders href /dashboard/users/${user.id}; the directory app/dashboard/users contains only page.tsx (no [id] route)',
    source_refs='frontend/apps/super-admin/src/app/dashboard/users/page.tsx:132', related_static_ids='', evidence_paths=NET + 'AUD-R4-SUPERADMIN-CONTROL-MATRIX.json; tests/runtime/ph13_r4_superadmin_ui.mjs',
    remediation_direction='Implement the user detail page or remove the link', regression_tests_needed='Every link rendered on the Users page resolves with 200')}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8'))); hd = rows[0]; have = {r[0] for r in rows}
for fid, v in F.items():
    if fid not in have: rows.append([dict(v, finding_id=fid).get(k, '') for k in hd])
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    if 'EX-RUN004-SUPERADMIN-UI' not in ex:
        csv.writer(h).writerow(['EX-RUN004-SUPERADMIN-UI', 'UI-HEADED-CHROME', 'Super Admin app control matrix: dashboard, tenants, payments, users (audit-logs and settings not reached: run cut at the time limit)', 'super-admin', 'PARTIAL', 'dashboard 5/5 controls PASS; tenants/payments partly; users View links 404 (F-076); several controls BLOCKED by click timeouts and a Chrome tab crash on the users page (tooling, not attributed to the product); paid checkout is disabled by the embedded offer (checkout_available=false) so billing lifecycle cannot run'])
print(len(rows) - 1)
