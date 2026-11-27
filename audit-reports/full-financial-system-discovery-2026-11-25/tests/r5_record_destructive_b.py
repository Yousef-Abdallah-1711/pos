"""RUN-005: records destructive batch B (34b rows, 53b result updates, F-069 evidence). Idempotent."""
import csv, os, re
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
N = 'evidence/network/RUN-005/PH-04/'
FILES = 'AUD-R5-DESTRUCTIVE-B.run1.json .. run4-leave-resign.json; AUD-R5-DESTRUCTIVE-B.json (B4/B12/B13 rerun)'
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [
    ('EX-RUN005-DESTR-ACCOUNT', 'UI-HEADED-CHROME', 'Delete on a disposable chart-of-accounts account', 'owner', 'PASS', 'confirm dialog, DELETE /accounting/accounts/419 200, row removed; evidence ' + N + FILES),
    ('EX-RUN005-DESTR-CATEGORY', 'UI-HEADED-CHROME', 'Delete on a disposable category', 'owner', 'PASS', 'confirm dialog "Are you sure you want to delete this category?", DELETE /categories/8 200; soft delete (row kept, active=false); evidence ' + N + FILES),
    ('EX-RUN005-DESTR-BANK-ACCOUNT', 'UI-HEADED-CHROME', 'Delete on a disposable bank account (card)', 'owner', 'PASS', 'confirm "Delete bank account ...?", DELETE /accounting/bank-accounts/5 200, row removed; evidence ' + N + FILES),
    ('EX-RUN005-DESTR-HR-INSURANCE-TRAINING', 'UI-HEADED-CHROME', 'Delete on disposable HR insurance and training records', 'owner', 'PASS', 'DELETE /hr/insurance/3 200 and DELETE /hr/training/2 200, rows removed; evidence ' + N + FILES),
    ('EX-RUN005-DESTR-DRIVER-DEACTIVATE', 'UI-HEADED-CHROME', 'Deactivate on a disposable driver', 'owner', 'PASS', 'PUT /admin/drivers/46/status 200, status active -> inactive; evidence ' + N + FILES),
    ('EX-RUN005-DESTR-LEAVE', 'UI-HEADED-CHROME', 'Approve and Reject on disposable pending leave requests', 'owner', 'PASS', 'modal confirm; POST /hr/leaves/14/approve 200 -> approved; POST /hr/leaves/15/reject 200 -> rejected; evidence ' + N + FILES),
    ('EX-RUN005-DESTR-RESIGNATION', 'UI-HEADED-CHROME', 'Complete and Delete on disposable resignations', 'owner', 'FAIL', 'Complete: POST /hr/resignations/9/complete 200 with settlement_amount 0 and the employee stays active with no termination date (F-069 reproduced through the UI); Delete: DELETE /hr/resignations/8 200; evidence ' + N + FILES),
    ('EX-RUN005-DESTR-INVOICE-VOID', 'UI-HEADED-CHROME', 'Void on a disposable draft customer invoice', 'owner', 'PASS', 'prompt for reason; POST /accounting/customer-invoices/20/void 200; the draft invoice journal 44795 reversed (see journal checks); evidence ' + N + FILES),
    ('EX-RUN005-DESTR-CONTRACT-USER-ROLE', 'UI-HEADED-CHROME', 'Delete on a disposable contract, user and role', 'owner', 'PARTIAL', 'contract: fixture creation did not return an id in the runner (not executed); user and role: row and Delete control found but the click produced no confirmation dialog and no request (rows unchanged, 1 each); not classified as a product defect without a manual check; evidence ' + N + 'AUD-R5-DESTRUCTIVE-B.json'),
]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    for r in new:
        if r[0] not in ex:
            csv.writer(h).writerow(r)
EXEC = [('accounting/chart-of-accounts', r'Delete', 'EXECUTED-RUN005 (account deleted after confirm)'),
        ('categories', r'Delete', 'EXECUTED-RUN005 (soft delete after confirm)'),
        ('accounting/bank-accounts', r'Delete', 'EXECUTED-RUN005 (deleted after confirm)'),
        ('hr/insurance', r'Delete', 'EXECUTED-RUN005 (deleted)'), ('hr/training', r'Delete', 'EXECUTED-RUN005 (deleted)'),
        ('drivers', r'Deactivate', 'EXECUTED-RUN005 (status inactive)'),
        ('hr/leaves', r'Approve|Reject', 'EXECUTED-RUN005 (approved/rejected via modal)'),
        ('hr/resignations', r'Complete|Delete', 'EXECUTED-RUN005 (Complete leaves employee active, F-069; delete OK)'),
        ('customer-invoices', r'Void', 'EXECUTED-RUN005 (draft void reverses journal)')]
rows = list(csv.reader(open('53b-ui-control-interaction-results.csv', encoding='utf-8')))
ix = {k: i for i, k in enumerate(rows[0])}
n = 0
for r in rows[1:]:
    if r[ix['result']] != 'SKIPPED-DESTRUCTIVE':
        continue
    for route, pat, res in EXEC:
        if route in r[ix['route']] and re.search(pat, r[ix['label']], re.I):
            r[ix['result']] = res.split(' (')[0]; r[ix['actual']] = res; r[ix['blocker']] = ''; r[ix['evidence']] = N + 'AUD-R5-DESTRUCTIVE-B*.json'; n += 1; break
csv.writer(open('53b-ui-control-interaction-results.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
U = {'F-069': ' RUN-005 real-UI reproduction: Complete on resignation 9 returned 200 with settlement_amount 0; the employee stayed active and no termination date was set.'}
f = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
fx = {k: i for i, k in enumerate(f[0])}
for r in f[1:]:
    if r[0] in U and U[r[0]].strip() not in r[fx['actual']]:
        r[fx['actual']] += U[r[0]]
        r[fx['evidence_paths']] += '; ' + N + 'AUD-R5-DESTRUCTIVE-B.run4-leave-resign.json'
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(f)
print('53b rows updated', n)
