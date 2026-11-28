"""RUN-004 addenda for 45, 48, 62, README and progress (idempotent: guarded by a marker). Counts are computed from the registers. Audit tooling only."""
import csv, os, re, collections
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
find = list(csv.DictReader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
sev = collections.Counter(x['severity'].split(' ')[0] for x in find)
conf = collections.Counter((re.match(r'\s*(CONFIRMED|PROBABLE|CANDIDATE|UNKNOWN)', x['root_cause_first_break']) or [None, 'UNSTATED'])[1] for x in find)
tasks = list(csv.DictReader(open('57-future-remediation-task-register.csv', encoding='utf-8')))
regs = list(csv.DictReader(open('59-post-remediation-regression-test-register.csv', encoding='utf-8')))
t51 = collections.Counter(r['status'] for r in csv.DictReader(open('51-phased-audit-task-register.csv', encoding='utf-8')))
M = '## RUN-004 addendum (2026-11-27)'
status_line = f"Findings {len(find)} (Critical {sev['Critical']}, High {sev['High']}, Medium {sev['Medium']}, Low {sev['Low']}); root cause {conf.get('CONFIRMED', 0)} CONFIRMED / {conf.get('PROBABLE', 0)} PROBABLE; future tasks {len(tasks)} (all NOT STARTED); future regression tests {len(regs)}; register 51 statuses {dict(t51)}."
A45 = f"""
{M}
{status_line}

Executed in RUN-004 (details in 34b rows EX-RUN004-*): six remaining report metrics (all 119 metric rows now measured; CF-002, PROF-007, DASH-004, HRR-002, HRR-007 fail, KDS-003 partial); SEC-013 staging comparison (F-071); write load 10-100 virtual users, a 600-payment burst and a 54,000-line report benchmark on a synthetic tenant; real-UI form matrix (17 forms x 5 modes); per-control interaction matrix (644 controls on 59 dashboard routes and 4 reached Super Admin routes, one safe action each); bank reconciliation (F-072, F-073); Theme Builder UI (F-077, F-078); landing smoke; 54 and 07 denominator correction (54b).

Still NOT executed (feasible but not done): single-field isolation tests (min/max, precision, date pairs, file types) for any of the 531 logical fields; 437 of the 531 logical fields belong to forms that were never opened (54b); Arabic and mobile INTERACTION passes (only the RUN-003 load matrix exists); invoice, payment, bank-account, POS, settings sub-form and storefront checkout form matrices; 116 destructive or money-moving controls (SKIPPED-DESTRUCTIVE in 53b) that need purpose-built fixtures; 66 controls with BLOCKED click timeouts (tooling); Super Admin audit-logs and settings pages; theme import, deletion of the active theme; browser rendering latency on the 54,000-line tenant; pool-exhaustion and database-restart reliability variants; per-tenant timezone verification of other reports (F-070 candidate).

Reliability reclassification: the Redis outage (F-018) is recorded as FAIL, not PASS: authenticated requests and login hang beyond 30 s while /health stays 200; after Redis returns, service resumes only after re-login because every session was lost (recovery PARTIAL). API restart with a pending outbox event (CFR-004) PASS; concurrent redrive PASS.

Blocked by product capability (unchanged, re-verified on the current baseline): F-001 new-tenant accounting bootstrap (tenants 5, 6, 7 remain unseeded), F-013 period close, F-014 employee create (reconfirmed through the UI), F-049 refunds, paid subscription checkout (offer.json checkout_available=false, enforced by middleware and a unit test), depreciation (fixed_assets table with no route or caller).
"""
A48 = f"""
{M}
Blocked versus unexecuted after RUN-004: see 45 addendum. New blockers: Super Admin billing lifecycle (product constant checkout_available=false); depreciation (no capability); SKIPPED-DESTRUCTIVE controls (fixture needed, not a product blocker); click-timeout BLOCKED rows are tooling, not product verdicts.
Prerequisite fixes that would unblock further audit work: F-014 (employee creation for payroll), F-079 (UI leave requests), F-013 (period close), F-049 (refund path), F-001 (tenant bootstrap).
"""
for f, text in (('45-final-coverage-and-unexecuted-tests.md', A45), ('48-blocked-vs-unexecuted-and-prerequisites.md', A48)):
    s = open(f, encoding='utf-8').read()
    if M not in s:
        open(f, 'a', encoding='utf-8').write(text)
print(status_line)
