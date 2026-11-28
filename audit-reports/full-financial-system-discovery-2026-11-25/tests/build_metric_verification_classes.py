"""RUN-005 Phase I.1: adds an explicit verification_class column to 30c-metric-coverage-reconciliation.csv and writes 30d-metric-verification-classes.md. Idempotent. Audit tooling only."""
import csv, os, collections
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
m43 = {}
for r in csv.DictReader(open('43-report-metric-runtime-results.csv', encoding='utf-8')):
    m43.setdefault(r['metric_ref'], r)
rows = list(csv.reader(open('30c-metric-coverage-reconciliation.csv', encoding='utf-8')))
hd = rows[0]
if 'verification_class' not in hd:
    hd.append('verification_class')
vi = hd.index('verification_class')
def klass(metric, status, cov):
    s = (m43.get(metric, {}).get('status') or status or '').upper()
    c = cov.upper()
    if 'BLOCKED' in c or 'BLOCKED' in s:
        return '4 attempted but blocked'
    if 'POLICY' in s:
        return '5 policy dependent'
    if 'VACUOUS' in s or 'VACUOUS' in c or 'ZERO DATA' in s or 'NOT VERIFIED' in s or 'OBSERVED' in s:
        return '3 measured without sufficient independent verification'
    if 'BOUNDARIES UNTESTED' in c or 'ONE BUCKET' in s:
        return '3 measured without sufficient independent verification'
    if 'NOT IN 30 INVENTORY' in c:
        return '3 measured without sufficient independent verification' if not s.startswith(('PASS', 'FAIL')) else ('2 independently verified and incorrect' if s.startswith('FAIL') or 'INCONSISTENT' in s else '1 independently verified and correct')
    if s.startswith('FAIL') or 'INCONSISTENT' in s or 'MISLEADING' in s:
        return '2 independently verified and incorrect'
    if s.startswith('PASS') or 'MEASURED VS INDEPENDENT' in c:
        return '1 independently verified and correct'
    return '3 measured without sufficient independent verification'
cnt = collections.Counter()
for r in rows[1:]:
    while len(r) < len(hd):
        r.append('')
    r[vi] = klass(r[0], r[4], r[5]); cnt[r[vi]] += 1
csv.writer(open('30c-metric-coverage-reconciliation.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
cov = collections.Counter(r[5] for r in rows[1:])
open('30d-metric-verification-classes.md', 'w', encoding='utf-8').write(f"""# 30d - Metric verification classes (RUN-005)

Built by `tests/build_metric_verification_classes.py`; the class is stored per row in the new `verification_class` column of [30c](30c-metric-coverage-reconciliation.csv).

Total rows: {len(rows) - 1}. Class counts: {dict(sorted(cnt.items()))}.

Relationship to the RUN-004 numbers (coverage_class column): {dict(cov)}. In RUN-004 terms: 99 rows were MEASURED vs an independent value (56 + 26 + 11 + 6), 8 were measured without one (5 added after the inventory + 3 vacuous) and 12 were blocked by F-009 (menu profitability and break-even): 99 + 8 + 12 = 119. "Measured vs independent value" is NOT the same as "correct": the correct/incorrect split above comes from the runtime status in 43 (FAIL, INCONSISTENT and MISLEADING map to class 2).

Class definitions: 1 independently verified and correct; 2 independently verified and incorrect (a finding exists); 3 measured without sufficient independent verification (vacuous data, untested boundaries, one populated bucket, or not verified); 4 attempted but blocked by a defect; 5 policy dependent (the expected value depends on an unapproved decision in 56); 6 missing capability and 7 not applicable (none among the 119 rows).

Matching the application's own SQL is not accepted as proof of accounting correctness; rows whose only comparison is the application query are class 3.
""")
print(dict(cnt))
