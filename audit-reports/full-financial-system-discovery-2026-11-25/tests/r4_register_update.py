"""RUN-004 register update (idempotent): adds F-066..F-070, 43 rows for the six remaining metrics, 30c statuses, 34b rows. Audit tooling only."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
NET = 'evidence/network/RUN-004/PH-12/'
DB = 'evidence/database/RUN-004/PH-12/'
SH = 'evidence/screenshots/RUN-004/PH-12/'
CAP = NET + 'AUD-R4-MET-CAPTURE.json; ' + NET + 'AUD-R4-MET2-CAPTURE.json'
F = {
    'F-066': dict(
        title='Dashboard home "Products" tile shows the page length (1) instead of the product total (12)',
        category='Software defect (UI aggregation)', severity='Low', evidence_status='REPRODUCED',
        affected_modules='Dashboard home',
        reproduction='Headed Chrome, tenant 4 owner: open /en/dashboard; compare the Products tile with GET /products?page=1&page_size=1 (pagination.total=12) and SQL (13 rows, 1 with status=deleted)',
        expected='Products tile = 12 (non-deleted menu items)',
        actual='Tile shows 1. The dashboard requests page_size=1 and reads total at the payload root, but the API returns it under pagination.total, so the code falls back to items.length = 1',
        monetary_impact_EGP='0 (wrong operational count; misleads the owner about catalogue size)',
        root_cause_first_break='CONFIRMED (source + captured response): frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:173-178 extractTotalFromCollection reads direct?.total (undefined) not pagination.total; request with page_size 1 at :314 and call at :329',
        source_refs='frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:173; frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:329; frontend/apps/dashboard/src/app/[locale]/dashboard/page.tsx:389',
        related_static_ids='MET-DASH-004', evidence_paths=CAP + '; ' + SH + 'AUD-R4-MET-dash.png',
        remediation_direction='Read pagination.total (or add a dedicated count endpoint); add a contract test for the tile',
        regression_tests_needed='Dashboard tile equals the API total when there is more than one product'),
    'F-067': dict(
        title='Profitability "Lowest Margin Items" is the profit-amount ranking read backwards, so a 100%-margin item ranks second-lowest',
        category='Software defect (report logic)', severity='Medium', evidence_status='REPRODUCED',
        affected_modules='Reports > Profitability (API, UI, CSV/PDF)',
        reproduction='Headed Chrome /en/dashboard/reports/profitability, tenant 4, 2026-10-29..2026-11-27; compare lowest_margin_items with margin percentages from SQL',
        expected='Items ordered by gross_margin_pct ascending (0%, 60%, 60%, 60%, 100%)',
        actual='Order is MD probe 0%, MD probe 100%, Probe A 60%, Probe B 60%, Burger 60%: the list is the gross-PROFIT descending list iterated backwards',
        monetary_impact_EGP='0 direct (decision quality: menu pricing review points at the wrong items)',
        root_cause_first_break='CONFIRMED: report_handler.go:2432 ORDER BY gross_profit DESC, then :2470-2476 builds the lowest list by iterating the same slice backwards; no margin sort exists',
        source_refs='backend/internal/handler/http/report_handler.go:2432; backend/internal/handler/http/report_handler.go:2470; backend/internal/handler/http/report_handler.go:212',
        related_static_ids='MET-PROF-007; F-040',
        evidence_paths=CAP + '; ' + DB + 'AUD-R4-ORACLE-PROF-CF.json; ' + SH + 'AUD-R4-MET2-prof.png',
        remediation_direction='Sort the lowest list by gross_margin_pct ASC (tie-break revenue) in SQL; rename the section if profit ranking is intended',
        regression_tests_needed='Fixture with items of different margin and profit; assert both lists'),
    'F-068': dict(
        title='Social-insurance CSV export lists an employee once per active insurance policy, repeating base salary and deductions',
        category='Software defect (export)', severity='Medium', evidence_status='REPRODUCED',
        affected_modules='HR reports > Social insurance export',
        reproduction='Tenant 4, month 2026-11: export with one active policy, add a second active policy for the same employee (API-DIAGNOSTIC), export again',
        expected='One row per employee salary record',
        actual='1 salary row (SQL) produced 2 data rows after the second active policy; both repeat base_salary 5,000.00 and the deductions, so any total over the file double-counts',
        monetary_impact_EGP='Proportional: each extra active policy repeats the employee base salary and deductions (here 5,000.00 base salary listed twice)',
        root_cause_first_break='CONFIRMED: accounting_advanced_repo.go:1421-1446 LEFT JOIN insurance_tracking it ON ... AND it.status active has no DISTINCT or per-employee limit',
        source_refs='backend/internal/repository/accounting_advanced_repo.go:1437; backend/internal/handler/http/accounting_advanced_handler.go:912',
        related_static_ids='MET-HRR-007',
        evidence_paths=NET + 'AUD-R4-HRR007-before.csv; ' + NET + 'AUD-R4-HRR007-after-second-policy.csv; ' + DB + 'AUD-R4-HR.json',
        remediation_direction='Aggregate policies per employee (string_agg or LATERAL LIMIT 1) so each salary row appears once',
        regression_tests_needed='Two active policies -> one row'),
    'F-069': dict(
        title='Completing a resignation never terminates the employee: turnover stays 0 and the person stays active and payable until the employee is deleted',
        category='Software defect (HR lifecycle)', severity='Medium', evidence_status='REPRODUCED',
        affected_modules='HR resignations, HR turnover report, employees, payroll eligibility',
        reproduction='Tenant 4: create employee (backend-native keys), resignation 2026-11-20 / last day 2026-11-25, acknowledge, complete (all 2xx); read the employee row and the turnover report',
        expected='Employee becomes terminated with termination_date = last working date; turnover departures = 1',
        actual='Employee remains employment_status=active, is_active=true, termination_date NULL; turnover departures 0, rate 0. Only DELETE /hr/employees/{id} sets termination_date (to CURRENT_DATE, not the last working date); then turnover = 1/1 = 100%',
        monetary_impact_EGP='Payroll and attendance continue for a resigned employee; settlement_amount is stored on the resignation with no payroll or GL effect',
        root_cause_first_break='CONFIRMED: resignation_repo.go:201-215 Complete updates only resignations; employee_repo.go:370-380 is the only termination_date writer; hr_report_handler.go:276-285 counts departures only from employees.termination_date',
        source_refs='backend/internal/repository/resignation_repo.go:201; backend/internal/repository/employee_repo.go:370; backend/internal/handler/http/hr_report_handler.go:276',
        related_static_ids='F-014; F-056; MET-HRR-002',
        evidence_paths=NET + 'AUD-R4-MET2-CAPTURE.json; ' + DB + 'AUD-R4-HR.json',
        remediation_direction='On resignation completion set employee status and termination_date from last_working_date in the same transaction; define notice-period policy',
        regression_tests_needed='Complete resignation -> employee terminated, departure counted, excluded from active payroll'),
    'F-070': dict(
        title='KDS busiest-hours and "today" boundaries are computed in UTC, not the restaurant local time',
        category='Software defect (time zone)', severity='Low', evidence_status='REPRODUCED (KDS only; other reports unverified)',
        affected_modules='KDS analytics and stats',
        reproduction='Tenant 4: compare /kds/analytics busiest_hours with SQL bucketed by UTC and by Africa/Cairo for the same kds_queue rows',
        expected='Hours and the "today" cut-off follow the restaurant time zone (Egypt UTC+2/+3)',
        actual='API hours equal the UTC buckets exactly (1:7, 17:4, 19:4, 20:18, 22:11); the Cairo buckets differ (0:11, 3:7, 19:4, 21:4, 22:18); kds_queue.received_at is timestamp without time zone; no tenant time-zone setting was found',
        monetary_impact_EGP='0 (scheduling signal shifted 2-3 hours); other date-bounded reports may share the cut-off issue (CANDIDATE, not tested)',
        root_cause_first_break='CONFIRMED for KDS: kds_repo.go:418-430 EXTRACT(HOUR FROM received_at) and CURRENT_DATE at :314-317; DB timezone is UTC',
        source_refs='backend/internal/repository/kds_repo.go:418; backend/internal/repository/kds_repo.go:314',
        related_static_ids='MET-KDS-003; F-044',
        evidence_paths=NET + 'AUD-R4-KDS-after-read.json; ' + NET + 'AUD-R4-KDS.json; ' + SH + 'AUD-R4-KDS-after-bump.png',
        remediation_direction='Introduce a restaurant time-zone setting and convert at query time; test cut-offs around local midnight',
        regression_tests_needed='Orders at 23:30 and 00:30 local land in the right day and hour'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
have = {r[0] for r in rows}
for fid, v in F.items():
    if fid in have:
        continue
    v = dict(v, finding_id=fid)
    rows.append([v.get(k, '') for k in hd])
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)

m43 = [
    ('MET-CF-002', 'Cash flow', 'investing_cash_flow (tenant 4, 2026-10-29..11-27)', 20000, -20000, 40000, 'FAIL',
     'RUN-004 headed Chrome + SQL: JE-2026-00082 Dr 1250 Delivery Van 20,000 / Cr 1110 Cash 20,000 (cash fell 20,000). Report shows +20,000 (debit-credit on 12%); correct cash-flow sign is -20,000. Same family as F-003 (mis-signed investing, catalog 12 line 47). Evidence ' + NET + 'AUD-R4-MET2-CAPTURE.json; ' + DB + 'AUD-R4-ORACLE-PROF-CF.json'),
    ('MET-PROF-007', 'Profitability', 'top_items / lowest_margin_items', 'top: Burger 288,000; Probe B 1,320; Probe A 240; MD 30.36; MD 0', 'identical (SQL)', 0, 'PASS (top) / FAIL (lowest) - F-067',
     'RUN-004: summary revenue 482,630.36 COGS 193,040.00 GP 289,590.36 equals SQL; top list equals SQL; lowest_margin_items is the profit ranking reversed. Lines of unpaid orders (191,400) are included (policy). Evidence ' + CAP + '; ' + DB + 'AUD-R4-ORACLE-PROF-CF.json'),
    ('MET-DASH-004', 'Dashboard', 'Recent Paid Revenue / Products tile', '33.40 / 1', '33.40 / 12', '0 / -11', 'PASS (revenue per its label) / FAIL (products) - F-066',
     'RUN-004: 33.40 = paid orders 74 (33.40) + 73 (0) among the latest 5 only. Orders 59, Pending 15, Staff 1, Drivers 2 match SQL; Categories 2 = active only; Products tile 1 vs 12. Evidence ' + CAP),
    ('MET-KDS-003', 'KDS', 'avg prep / total completed / trend / busiest hours (no late-orders metric exists)', '917.8365 min; 2', '917.8365; 2', 0, 'PASS (arithmetic) / FAIL (time zone) - F-070',
     'RUN-004 headed Chrome (orders bumped through the KDS UI) + SQL: analytics equal SQL. Busiest hours are UTC buckets. Lateness exists only as client-side colour thresholds (12/25 min); the inventory label Late orders was a guess. Evidence ' + NET + 'AUD-R4-KDS-after-read.json'),
    ('MET-HRR-002', 'HR reports', 'turnover_rate', '0 after completed resignation; 100 after DELETE', '100 after resignation completion', -100, 'FAIL - F-069',
     'RUN-004 API-DIAGNOSTIC + headed Chrome page capture: formula departures/avg headcount*100 is arithmetically right but departures come only from employees.termination_date, which resignation completion never sets. Evidence ' + DB + 'AUD-R4-HR.json'),
    ('MET-HRR-007', 'HR reports', 'social insurance CSV rows', 2, 1, 1, 'FAIL - F-068',
     'RUN-004 API-DIAGNOSTIC: one salary row gives two CSV rows with two active policies. Evidence ' + NET + 'AUD-R4-HRR007-after-second-policy.csv'),
]
ex = {r[0] for r in csv.reader(open('43-report-metric-runtime-results.csv', encoding='utf-8'))}
with open('43-report-metric-runtime-results.csv', 'a', newline='', encoding='utf-8') as h:
    w = csv.writer(h)
    for r in m43:
        if r[0] not in ex:
            w.writerow(r)

rows = list(csv.reader(open('30c-metric-coverage-reconciliation.csv', encoding='utf-8')))
st = {r[0]: r[6].split(' - ')[0] for r in m43}
for r in rows[1:]:
    if r[0] in st and r[4] == 'NOT EXECUTED':
        r[3] = 'YES'
        r[4] = 'EXECUTED (RUN-004): ' + st[r[0]]
        r[5] = 'MEASURED vs independent value (RUN-004, headed Chrome + SQL)'
csv.writer(open('30c-metric-coverage-reconciliation.csv', 'w', newline='', encoding='utf-8')).writerows(rows)

ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [
    ('EX-RUN004-MET-CF002', 'UI-HEADED-CHROME+SQL', 'MET-CF-002 investing cash flow vs ledger', 'owner', 'FAIL', 'report +20,000 vs ledger cash movement -20,000 (F-003 family)'),
    ('EX-RUN004-MET-PROF007', 'UI-HEADED-CHROME+SQL', 'MET-PROF-007 profitability top/lowest items', 'owner', 'FAIL', 'top list = SQL; lowest_margin list ranks by profit (F-067)'),
    ('EX-RUN004-MET-DASH004', 'UI-HEADED-CHROME+SQL', 'MET-DASH-004 dashboard tiles vs API/SQL', 'owner', 'FAIL', 'recent paid revenue 33.40 correct per label; Products tile 1 vs 12 (F-066)'),
    ('EX-RUN004-MET-KDS003', 'UI-HEADED-CHROME+SQL', 'MET-KDS-003 KDS analytics after UI bump to completed', 'owner', 'PARTIAL', 'averages/totals = SQL; hours are UTC buckets (F-070); no late-orders metric exists'),
    ('EX-RUN004-MET-HRR002', 'API-DIAGNOSTIC+HEADED-CHROME', 'MET-HRR-002 turnover through resignation lifecycle', 'owner', 'FAIL', 'completed resignation not counted as departure; employee stays active (F-069)'),
    ('EX-RUN004-MET-HRR007', 'API-DIAGNOSTIC', 'MET-HRR-007 social insurance export duplicates', 'owner', 'FAIL', 'one salary row gives two CSV rows with two active policies (F-068)'),
]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    w = csv.writer(h)
    for r in new:
        if r[0] not in ex:
            w.writerow(r)
print('done; findings now', len(list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))) - 1)
