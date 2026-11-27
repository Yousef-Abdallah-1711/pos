"""Applies RUN-004 results to 51-phased-audit-task-register.csv (idempotent; run AFTER task_status_overlay.py). Status vocabulary unchanged.
Each row separates: deliverable created / verified / execution completed / remaining. Audit tooling only."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
TAG = ' | RUN-004: '
U = {
    'T-P12-002': ('EXECUTED-PARTIAL', 'all 119 metric rows now measured (6 added in RUN-004: CF-002, PROF-007, DASH-004, KDS-003, HRR-002, HRR-007); remaining: export/print comparison and multi-dataset variants', 'F-003,F-066,F-067'),
    'T-P12-003': ('EXECUTED-PARTIAL', 'KDS analytics measured after real UI bump (F-070), HR turnover/export (F-068, F-069), dashboard tiles (F-066); remaining: empty/single/multi dataset variants and date boundaries', 'F-066,F-068,F-069,F-070'),
    'T-P17-001': ('EXECUTED-PARTIAL', 'SEC-013 EXECUTED (F-071): staging API with RBAC unset lets a cashier reach finance/HR/user endpoints; production probe refuses mock payments; remaining: SEC-011/012/014 re-run on the current tree', 'F-071'),
    'T-P11-003': ('EXECUTED-PARTIAL', 'bank reconciliation executed (F-072, F-073); depreciation is BLOCKED-CAPABILITY (fixed_assets table, no route/UI); import/export covered earlier; budgets/FX/recurring covered in RUN-003', 'F-072,F-073'),
    'T-P18-002': ('EXECUTED-PASS', 'write load 10-100 VU (5,719 orders, 0 server errors), 600-payment burst with outbox drain, and 54,000-line report benchmark executed and reconciled; browser rendering latency for the 54k tenant not measured', ''),
    'T-P04-001': ('EXECUTED-PARTIAL', 'per-control interaction matrix: 644 unique controls, 59 dashboard routes + 4 Super Admin routes, one safe action each (53b); 88 destructive/money-moving controls skipped, 58 tooling-blocked, cap of 22 controls per route', 'F-009,F-074,F-076,F-077'),
    'T-P04-002': ('EXECUTED-PARTIAL', 'form-level matrices (empty/invalid/Arabic/long/valid) for 17 forms; field-level isolation tests NOT done; denominator corrected to 531 logical fields (54b)', 'F-074,F-075,F-079'),
    'T-P05-003': ('EXECUTED-PARTIAL', 'user and role creation through the real UI in the form matrix (valid and Arabic persist; HTML/SQL text accepted; long name 500); user-role assignment screens only loaded', 'F-062,F-031'),
    'T-P08-002': ('EXECUTED-PARTIAL', 'Theme Builder real-UI run: export, duplicate, editor edit, activation, storefront/public effect, Appearance coupling, restore; import and deletion of the active theme not exercised', 'F-077,F-078,F-059,F-060'),
    'T-P08-004': ('EXECUTED-PARTIAL', 'driver create form exercised (valid/Arabic persist, long name 500); assignment, completion and COD settlement not executed (F-057 capability gap)', 'F-057,F-058'),
    'T-P20-001': ('EXECUTED-PARTIAL', 'deliverable created: findings/ folders exist for 64 findings (F-066..F-070 pending generation); verified: evidence 69/69 resolved, root cause 66 CONFIRMED / 3 PROBABLE; NOT verified: folder contents are largely generated from register fields (thin design/paths) and must be enriched per finding', ''),
    'T-P21-001': ('EXECUTED-PARTIAL', 'RUN-004 correction (was EXECUTED-PASS): specs/tasks exist for every finding but most are template-level (e.g. F-031: no resolved backend path, one-line design); acceptance criteria "implementation-ready" not met for those', ''),
    'T-P22-001': ('EXECUTED-PARTIAL', 'deliverable created: 61 closure report and 60 manifest exist; execution NOT completed (remaining audit work and restoration still open); report 61 must be regenerated only after remaining work', ''),
    'T-P22-002': ('PLANNED', 'restoration NOT done; port 3001 owner (pid 9064, landing next start) stopped in RUN-004; dev containers still stopped until final closure', ''),
}
NEW = [
    ('T-R04-BASE-001', '00', 'audit', 'Verify new git baseline, record in 63; separate historical/authorized/current/pre-existing/audit files', 'EXECUTED-PASS', 'HEAD unchanged 133260d; frontend work is uncommitted; backend diff none; 63 written; evidence RUN-004/01-git-baseline.txt', ''),
    ('T-R04-BASE-002', '00', 'dashboard', 'Review dashboard api.ts change (401 redirect) and smoke-test login + one mutation on current frontend', 'EXECUTED-PARTIAL', 'diff reviewed (401 redirect only); headed Chrome login and mutations (KDS bump, settings) worked on the current working tree; dedicated smoke script not yet written', ''),
    ('T-R04-BASE-003', '00', 'landing', 'Verify landing registration/login/language switching on current tree (read-only)', 'EXECUTED-PASS', 'smoke of 12 routes desktop+mobile EN/AR passed; language switch works; registration and login pages render (submission itself exercised in RUN-003); hydration warning observed (OBS-LAND-001)', ''),
    ('T-R04-ENV-001', '02', 'env', 'Resume isolated stack, port 3001 conflict, CORS wildcard', 'EXECUTED-PASS', 'stack up on persistent volume; pid 9064 stopped; API recreated with storefront CORS; frontends recreated', ''),
    ('T-R04-RC-001', '20', 'findings', 'Reconcile root-cause classification for all findings', 'EXECUTED-PASS', '69 findings: 66 CONFIRMED, 3 PROBABLE (F-018, F-021, F-038 with exact gaps), 0 unstated; F-064 duplicate claim retracted', 'F-064'),
]
rows = list(csv.reader(open('51-phased-audit-task-register.csv', encoding='utf-8')))
hd = rows[0]
ix = {k: i for i, k in enumerate(hd)}
for r in rows[1:]:
    t = r[ix['task_id']]
    if t in U:
        st, res, fin = U[t]
        r[ix['status']] = st
        if TAG not in r[ix['actual_result']]:
            r[ix['actual_result']] = r[ix['actual_result']] + TAG + res
        if fin:
            ids = set(filter(None, r[ix['finding_ids']].split(',')))
            r[ix['finding_ids']] = ','.join(sorted(ids | set(fin.split(','))))
have = {r[ix['task_id']] for r in rows[1:]}
for t, ph, mod, obj, st, res, fin in NEW:
    if t in have:
        continue
    row = [''] * len(hd)
    row[ix['task_id']], row[ix['phase']], row[ix['module']], row[ix['objective']] = t, ph, mod, obj
    row[ix['status']], row[ix['actual_result']], row[ix['finding_ids']] = st, 'RUN-004: ' + res, fin
    rows.append(row)
csv.writer(open('51-phased-audit-task-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
print(len(rows) - 1, 'tasks')
