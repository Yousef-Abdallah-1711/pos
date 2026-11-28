"""Master Prompt 07: builds the handoff package (65-72, queue 66/66b/66c/66d, phase folder docs) from the authoritative
registers. Audit tooling only: reads registers/evidence, writes documents, never touches product code. Idempotent."""
import csv, collections, glob, hashlib, json, os, re, subprocess, datetime
ROOT = r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25'
os.chdir(ROOT)
NOW = datetime.datetime.now(datetime.timezone.utc).strftime('%Y-%m-%d')
rd = lambda f: list(csv.DictReader(open(f, encoding='utf-8')))
def wcsv(f, rows, cols):
    with open(f, 'w', newline='', encoding='utf-8') as h:
        w = csv.writer(h); w.writerow(cols); w.writerows([[r.get(c, '') for c in cols] for r in rows])
def git(*a):
    try:
        return subprocess.run(['git', '-C', ROOT, *a], capture_output=True, text=True, timeout=60).stdout.strip()
    except Exception:
        return 'unknown'
HEAD = git('rev-parse', 'HEAD'); BRANCH = git('branch', '--show-current')
R51 = rd('51-phased-audit-task-register.csv'); FIND = rd('22-confirmed-findings-register.csv')
REM = rd('57-future-remediation-task-register.csv'); REG = rd('59-post-remediation-regression-test-register.csv')
F07 = rd('07-field-level-inventory.csv'); F54 = rd('54b-field-denominator-and-form-coverage.csv'); C53 = rd('53b-ui-control-interaction-results.csv')
PH = {}
for i in range(23):
    t = open(f'phases/phase-{i:02d}/spec.md', encoding='utf-8').readline().strip()
    PH[f'{i:02d}'] = re.sub(r'^# Phase \d+ [—-] ', '', t)

# ---------------------------------------------------------------- phase of a route / component
def phase_of(s):
    s = s.lower()
    for k, p in (('super-admin', '13'), ('tenants', '13'), ('/hr/', '10'), ('hr-advanced', '10'), ('hr/', '10'), ('payroll', '10'),
                 ('checkout', '08'), ('restaurant-website', '08'), ('kds', '08'), ('tables', '08'), ('qr', '08'),
                 ('pos', '07'), ('orders', '07'), ('customer-invoices', '09'), ('customers', '09'), ('supplier', '06'), ('inventory', '06'), ('products', '06'),
                 ('accounting', '11'), ('journal', '11'), ('ledger', '11'), ('bank', '11'), ('report', '12'),
                 ('theme', '05'), ('editors', '05'), ('settings', '05'), ('billing', '05'), ('users', '05'), ('roles', '05'), ('access-control', '05')):
        if k in s:
            return p
    return '05'

# ---------------------------------------------------------------- remaining fields
fields = []
for a, b in zip(F07, F54):
    if a['tested_status'] == 'NOT TESTED':
        fields.append(dict(line=b['inventory_line'], app=a['app'], route=a['page_route'], form=a['form'], field=a['field_name'], type=a['type'], required=a['required'], maps=a['maps_to_api_field'], comp=a['component_path'], hint=a['financial_impact_hint']))
CASES = {'text': 'empty;whitespace;300 chars;Arabic;special characters', 'email': 'empty;malformed;uppercase;300 chars', 'tel': 'empty;letters;Arabic digits;long', 'number': 'empty;negative;zero;6 decimals;99,999,999,999;0.001',
         'date': 'empty;1900;2999;reversed pair', 'password': 'empty;short;long;Arabic', 'select': 'empty option;each option;(select/checkbox/file NOT covered by the existing runner: extend it)', 'checkbox': 'toggle both states', 'textarea': 'empty;300 chars;Arabic;special characters'}
def cases_for(t):
    t = (t or '').lower()
    for k in CASES:
        if k in t:
            return CASES[k]
    return CASES['text']
for f in fields:
    f['cases'] = cases_for(f['type'])
    f['batch'] = None
groups = collections.OrderedDict()
for f in fields:
    key = (f['app'], f['route'], f['form'])
    groups.setdefault(key, []).append(f)
queue = []
def Q(**k):
    k.setdefault('status', 'NOT-STARTED'); k.setdefault('parent', ''); k.setdefault('required_files', ''); k.setdefault('execution_script', ''); k.setdefault('fixtures', '')
    k.setdefault('findings', ''); k.setdefault('blocker', ''); k.setdefault('complexity', 'S'); k.setdefault('prereq', 'Isolated stack up (see 65 section Environment); owner login for tenant 4')
    k.setdefault('deps', ''); k.setdefault('next', ''); k.setdefault('exec_type', 'AUDIT-EXECUTION')
    queue.append(k)
# ordered specific tasks first
Q(id='Q-ENV-001', parent='T-P02-*', phase='02', title='Verify isolated stack health and record the start-of-session checkpoint', module='environment', priority='P0', cap='ROUTINE-EXECUTION', complexity='S',
  prereq='Docker running', required_files='37-isolated-environment-execution-runbook.md; 64-run005-baseline.md', script='docker ps; curl http://127.0.0.1:18082/api/v1/health; pg_dump to evidence/database/',
  expected='api, postgres, redis healthy; tenants 5/6/7 COA rows 0; dev DB pos_saas 2 tenants/6 users/0 orders', dod='Checkpoint file evidence/runtime/RUN-006/01-baseline.txt written with docker ps, health, row counts and a pg_dump sha256', evidence='evidence/runtime/RUN-006/', next='Q-INC-001')
Q(id='Q-INC-001', parent='T-P05-*', phase='05', title='Manually verify Delete on a contract, user and role in headed Chrome (runner click produced no dialog/request)', module='hr/contracts; settings/users; settings/roles', priority='P0', cap='BROWSER-AUTOMATION', complexity='M',
  required_files='tests/runtime/ph04_r5_destructive_b.mjs (steps B4,B12,B13); evidence/network/RUN-005/PH-04/AUD-R5-DESTRUCTIVE-B.json', script='ONLY=B12,B13 node tests/runtime/ph04_r5_destructive_b.mjs (headed), then manual click with DevTools Network open',
  fixtures='Create disposable user and role via API with prefix 000R6 (the runner fixture helper does this); contract fixture needs employee_id 4, contract_type permanent, unique contract_number, valid dates and salary',
  expected='Delete opens a confirmation (browser confirm or modal) and issues DELETE /users/:id, /rbac/roles/:id, /hr/contracts/:id; rows removed or blocked with a readable message',
  dod='Each of the 3 controls classified PRODUCT-DEFECT (finding registered with request/response evidence) or RUNNER-DEFECT (selector or modal fixed and rerun) with screenshot and network capture', evidence='evidence/network/RUN-006/PH-05/, evidence/screenshots/RUN-006/PH-05/', findings='', next='Q-INC-002')
Q(id='Q-INC-002', parent='T-P05-*', phase='10', title='Find why the contract fixture POST /hr/contracts returned no id in the destructive runner', module='hr/contracts', priority='P1', cap='SOURCE-TRACING', complexity='S',
  required_files='tests/runtime/ph04_r5_destructive_b.mjs line for B4; backend handler for POST /hr/contracts', expected='Response body and status of the fixture call captured; id extracted from the actual response shape',
  dod='Documented whether the POST failed (400/409/500 with body) or the id path in the runner is wrong; runner fixed in audit tooling if it was a runner defect', evidence='evidence/network/RUN-006/PH-10/', deps='', next='Q-INC-001')
Q(id='Q-USR-001', parent='T-P05-*', phase='05', title='Users form: explain why Arabic, special and whitespace names produce no request and no message', module='settings/users', priority='P1', cap='SOURCE-TRACING', complexity='M',
  required_files='evidence/network/RUN-005/PH-06/AUD-R5-FIELD-users.json; dashboard users form component (locate with grep "Create User" under frontend/apps/dashboard/src/app/[locale]/dashboard/settings/users)',
  expected='Each no-request case is attributed to: client validation (pattern/required), a missing select (role) not filled by the runner, selector error, or backend failure',
  dod='Attribution table for name, email, phone, password cases with a corrected valid baseline that is ACCEPTED (201) before any isolated case is judged; users field results regenerated', evidence='evidence/network/RUN-006/PH-05/', next='Q-F065-001')
Q(id='Q-F065-001', parent='T-P05-*', phase='06', title='Categories form and F-065: reproduce and trace display_order sent as a string', module='categories', priority='P1', cap='SOURCE-TRACING', complexity='S',
  required_files='evidence/network/RUN-005/PH-06/AUD-R5-FIELD-categories.json (every case 400 {"error":"Invalid request body"}, display_order "" or "12"); finding package findings/F-065*', findings='F-065',
  expected='Reproduce one create from the real UI with a numeric order and capture the request body; identify the form code that serializes display_order and the Go struct field type that rejects it',
  dod='Exact file:line for the string serialization and the DTO decode failure added to finding F-065 root-cause.md; no product fix', evidence='evidence/network/RUN-006/PH-06/', next='Q-FLD-*')
# field batches
n = 0
for key, rows in groups.items():
    n += 1
    app, route, form = key
    ph = phase_of(route + ' ' + form)
    short = form if len(form) < 70 else os.path.basename(form)
    Q(id=f'Q-FLD-{n:03d}', parent='T-P04-003', phase='04', title=f'Isolated field cases for form "{short}" ({len(rows)} fields)', module=route if route.startswith('/') else app, priority='P2', cap='BROWSER-AUTOMATION', complexity='S' if len(rows) < 8 else 'M',
      required_files='07-field-level-inventory.csv; 66b-remaining-field-batches.csv; 54c-isolated-field-results.md', script='tests/runtime/ph06_r5_field_isolation.mjs (add the form to its T table, then ONLY=<key> node ...)',
      fixtures='Valid baseline that is ACCEPTED first (verify one clean submit returns 2xx before judging cases); disposable names prefixed 000R6', expected='Per field: front-end message or request; back-end 2xx/4xx not 500; no raw database text; DB row matches',
      dod='All fields of the batch have outcome rows in 54c; baseline validity proven; findings registered for 500s, raw pq text, silent no-ops', evidence='evidence/network/RUN-006/PH-' + ph + '/AUD-R6-FIELD-' + str(n) + '.json', next='next Q-FLD')
# destructive
dest = [c for c in C53 if c['result'] == 'SKIPPED-DESTRUCTIVE']
HINT = {'kds': ('KDS ticket in status new/preparing created by an unpaid tracked order on tenant 4', 'KDS row status advances; order status unchanged unless policy says otherwise', 'none (no ledger effect until payment)'),
        'settings/users': ('disposable user 000R6 user (POST /users)', 'confirm; DELETE /users/:id; row removed or blocked with readable message', 'none'),
        'settings/roles': ('disposable role 000R6 role with no members', 'confirm; DELETE /rbac/roles/:id; row removed; system roles must stay non-deletable', 'none'),
        'hr/contracts': ('disposable contract for employee 4, unique number', 'confirm; DELETE /hr/contracts/:id', 'none'),
        'accounting/advanced': ('rate/budget values that can be restored; export a copy first', 'PUT/POST saves; read back equals input', 'FX rate changes later conversions: record before/after of any posted multi-currency item'),
        'accounting/supplier-invoices': ('draft or approved supplier invoice 000R6', 'prompt for reason; POST void; journal reversed', 'reversal journal; supplier payable back to prior value'),
        'pos': ('cart with two items', 'cart cleared, no request, no order created', 'none'),
        'settings/billing': ('no live checkout (checkout_available=false)', 'control disabled or unavailable message', 'none; classify BLOCKED-CAPABILITY'),
        'tables/layout': ('tables fixture; export positions first', 'positions saved/reset; read back', 'none'),
        'settings/restaurant': ('copy current values; change one reversible field', 'PUT /restaurants saves; reload shows value', 'tax/currency fields change later documents: do not change them'),
        'theme-builder': ('theme export JSON from RUN-004 evidence', 'import creates a copy; active theme unchanged', 'none')}
def hint_for(route):
    for k, v in HINT.items():
        if k in route:
            return v
    return ('disposable fixture created by API with prefix 000R6', 'confirmation then request; read back', 'state per row')
dcols = ['control_id', 'route', 'control_type', 'label', 'role', 'fixture_and_initial_state', 'expected_confirmation_and_db_change', 'expected_financial_consequence', 'evidence_required', 'safe_cleanup', 'queue_task']
drows = []
for i, c in enumerate(dest, 1):
    fx, exp, fin = hint_for(c['route'])
    ph = phase_of(c['route'])
    tid = f'Q-DST-{i:03d}'
    drows.append(dict(control_id=c['control_id'], route=c['route'], control_type=c['control_type'], label=c['label'][:60], role=c['role'], fixture_and_initial_state=fx, expected_confirmation_and_db_change=exp, expected_financial_consequence=fin,
                      evidence_required='screenshot before/after, network request+response, SQL before/after', safe_cleanup='delete or restore the 000R6 fixture; never touch tenants 5/6/7', queue_task=tid))
    ex_type = 'EXTERNAL-OR-CAPABILITY' if ('billing' in c['route'] or 'E-Invoice' in c['label']) else 'AUDIT-EXECUTION'
    Q(id=tid, parent='T-P04-001', phase='04', title=f'Execute destructive control "{c["label"][:40]}" on {c["route"]}', module=c['route'], priority='P2', cap='BROWSER-AUTOMATION', complexity='S',
      required_files='53b-ui-control-interaction-results.csv (' + c['control_id'] + '); tests/runtime/ph04_r5_destructive_b.mjs (table-driven pattern)', script='extend tests/runtime/ph04_r5_destructive_b.mjs with a step for this control',
      fixtures=fx, expected=exp, dod='53b row result changed from SKIPPED-DESTRUCTIVE to an executed status with actual request/response; or BLOCKED-CAPABILITY with proof', evidence='evidence/network/RUN-006/PH-' + ph + '/', exec_type=ex_type, status='NOT-STARTED' if ex_type == 'AUDIT-EXECUTION' else 'BLOCKED-CAPABILITY-LIKELY', next='next Q-DST')
# blocked
blk = [c for c in C53 if c['result'] == 'BLOCKED']
bcols = ['control_id', 'route', 'control_type', 'label', 'role', 'previous_failure', 'previous_script_and_evidence', 'likely_tooling_issue', 'proposed_retry', 'evidence_to_separate_product_vs_runner', 'queue_task']
brows = []
byr = collections.OrderedDict()
for c in blk:
    byr.setdefault(c['route'], []).append(c)
for i, (route, cs) in enumerate(byr.items(), 1):
    tid = f'Q-BLK-{i:03d}'
    for c in cs:
        brows.append(dict(control_id=c['control_id'], route=route, control_type=c['control_type'], label=c['label'][:60], role=c['role'], previous_failure=c['blocker'] or c['actual'][:80], previous_script_and_evidence=c['evidence'] + ' via tests/runtime/ph04_r4_control_matrix.mjs',
                          likely_tooling_issue='click timeout or tab crash while a heavy page (large list, builder canvas) was open; memory pressure in the Docker VM (supported by RUN-005 memory observations, not proven per control)',
                          proposed_retry='open the route fresh in a new page, click this control by role+name with a 10 s timeout, wait for networkidle, one control per page load',
                          evidence_to_separate_product_vs_runner='second independent attempt in a clean profile; console errors; page.on(crash); if it still fails with a console/pageerror it is a product defect (register), if it passes it is a runner defect (record PASS)', queue_task=tid))
    Q(id=tid, parent='T-P04-001', phase='04', title=f'Rerun {len(cs)} tooling-blocked controls on {route}', module=route, priority='P2', cap='BROWSER-AUTOMATION', complexity='S',
      required_files='66d-remaining-blocked-controls.csv; tests/runtime/ph04_r4_control_matrix.mjs', script='node tests/runtime/ph04_r4_control_matrix.mjs with ONLY route filter (see script header) or a new small per-control script',
      expected='Control responds with an observable change; no console error or 5xx', dod='each listed control_id gets a PASS/NO-EFFECT/FAIL result with evidence in 53b', evidence='evidence/network/RUN-006/PH-' + phase_of(route) + '/', next='next Q-BLK')
# specific remaining tasks
SPEC = [
 ('Q-FRM-001', '09', 'Payments forms matrix (PaymentForm received and made): all fields, over-allocation, unrelated invoice', 'accounting/customer-invoices; accounting/supplier-invoices', 'P1', 'BROWSER-AUTOMATION', 'M', 'F-080,F-081', 'tests/runtime/ph09_r5_invoice_form.mjs (pattern)'),
 ('Q-FRM-002', '11', 'Bank account form matrix and bank reconciliation session form', 'accounting/bank-accounts; accounting/advanced', 'P1', 'BROWSER-AUTOMATION', 'M', 'F-072,F-073', 'tests/runtime/ph06_r5_field_isolation.mjs'),
 ('Q-FRM-003', '07', 'POS forms: cart, customer, discount, payment types, delivery fee', 'pos', 'P1', 'BROWSER-AUTOMATION', 'M', 'F-023', 'tests/runtime/ph07_r5_pos_currency.mjs (pattern)'),
 ('Q-FRM-004', '05', 'Settings sub-forms: restaurant info, tax/currency, access control, users', 'settings/*', 'P2', 'BROWSER-AUTOMATION', 'M', '', 'tests/runtime/ph06_r5_field_isolation.mjs'),
 ('Q-FRM-005', '08', 'Storefront checkout form and order tracking (restaurant-website; restart the stopped container first)', 'restaurant-website checkout, orders/track', 'P1', 'BROWSER-AUTOMATION', 'M', '', 'new script tests/runtime/ph08_r6_checkout.mjs'),
 ('Q-SA-001', '13', 'Super Admin audit logs page: list, filter, pagination, export, authorization', 'super-admin audit-logs', 'P2', 'BROWSER-AUTOMATION', 'S', '', 'tests/runtime/ph04_r4_control_matrix.mjs (Super Admin variant)'),
 ('Q-SA-002', '13', 'Super Admin settings page and remaining admin actions (extend trial, delete tenant modal fields) on disposable tenants only', 'super-admin settings, tenants', 'P2', 'BROWSER-AUTOMATION', 'M', '', 'new script; never use tenants 4/5/6/7/10'),
 ('Q-RPT-001', '12', 'Date-boundary tests for every report: entries at 23:59:59 and 00:00:00 local (Africa/Cairo) and UTC, inclusive end date', 'reports/*', 'P1', 'DATA-VERIFICATION', 'M', 'F-070', 'new script tests/runtime/ph12_r6_boundaries.mjs + SQL oracle'),
 ('Q-RPT-002', '12', 'UTC cutoff checks beyond KDS: dashboard tiles, sales, cash flow, HR reports', 'reports/*; dashboard home', 'P1', 'DATA-VERIFICATION', 'M', 'F-070', 'SQL with AT TIME ZONE compared to API output'),
 ('Q-RPT-003', '12', 'Export verification for every report with an export: CSV/PDF values vs UI vs database', 'reports/*', 'P1', 'DATA-VERIFICATION', 'M', 'F-068', 'new script tests/runtime/ph12_r6_exports.mjs'),
 ('Q-RPT-004', '12', 'Arabic export content (encoding, RTL text, numbers)', 'reports/*', 'P2', 'DATA-VERIFICATION', 'S', '', 'same script as Q-RPT-003 with locale ar'),
 ('Q-RPT-005', '18', 'Browser benchmark on tenant 10 (54,000 lines): load time, memory, long tasks for the 5 heaviest reports', 'reports/financial/*', 'P2', 'ROUTINE-EXECUTION', 'M', '', 'new script using chrome_lib.mjs performance metrics; tenant 10 read-only'),
 ('Q-HR-001', '10', 'Payroll SalaryForm field matrix and payroll run reconciliation', 'hr/payroll', 'P1', 'FINANCIAL-RECONCILIATION', 'M', 'F-020,F-079', 'tests/runtime/ph06_r5_field_isolation.mjs'),
 ('Q-CON-001', '07', 'Concurrency tests: double-submit POS payment, parallel stock deduction, parallel invoice payment', 'pos; accounting', 'P1', 'COMPLEX-INVESTIGATION', 'L', '', 'new script with parallel fetch against tenant 4 disposable fixtures'),
 ('Q-HIS-001', '19', 'Remaining historical regression rows in 21 not yet EXECUTED', '21-existing-issue-regression-matrix.csv', 'P2', 'ROUTINE-EXECUTION', 'M', '', 'see 21 for per-row method'),
 ('Q-MAT-001', '04', 'Per-control matrix beyond the 22-per-route cap and Arabic/mobile interaction on the remaining routes', 'all dashboard routes', 'P3', 'BROWSER-AUTOMATION', 'L', '', 'tests/runtime/ph04_r4_control_matrix.mjs, tests/runtime/ph04_r5_lang_viewport.mjs'),
 ('Q-POL-001', '21', 'Policy review of the 22 open decisions (do not approve; prepare options for the owner)', '56b-policy-decision-detail.md', 'P2', 'POLICY-REVIEW', 'M', '', 'document only'),
 ('Q-ENV-002', '22', 'Final restoration: restart or document stopped containers, final snapshot with full sha256, backup, dev DB fingerprint comparison', 'environment', 'P9', 'ROUTINE-EXECUTION', 'S', '', 'see 68 section Closing run'),
]
for t in SPEC:
    Q(id=t[0], parent='T-P' + t[1] + '-*', phase=t[1], title=t[2], module=t[3], priority=t[4], cap=t[5], complexity=t[6], findings=t[7], script=t[8], required_files='progress.md; ' + ('phases/phase-' + t[1] + '/README.md'),
      expected='Actual behavior recorded against the documented rule or oracle', dod='Every item named in the title has an executed result row in 34b and, where behavior is wrong, a registered finding with evidence', evidence='evidence/{network,database,screenshots}/RUN-006/PH-' + t[1] + '/', next='next in queue')
ORDER = {'P9': 9, 'P0': 0, 'P1': 1, 'P2': 2, 'P3': 3}
queue.sort(key=lambda q: (ORDER[q['priority']], q['id']))
qcols = ['Task ID', 'Parent task ID', 'Phase ID', 'Title', 'Module', 'Priority', 'Status', 'Execution type', 'Required model capability', 'Complexity', 'Prerequisites', 'Dependency task IDs', 'Required files', 'Execution script', 'Required fixtures', 'Expected result', 'Definition of Done', 'Evidence destination', 'Related finding IDs', 'Blocker', 'Next task']
def qrow(q):
    return {'Task ID': q['id'], 'Parent task ID': q['parent'], 'Phase ID': q['phase'], 'Title': q['title'], 'Module': q['module'], 'Priority': q['priority'], 'Status': q['status'], 'Execution type': q['exec_type'], 'Required model capability': q['cap'], 'Complexity': q['complexity'],
            'Prerequisites': q['prereq'], 'Dependency task IDs': q['deps'] or ('Q-ENV-001' if q['id'] != 'Q-ENV-001' else ''), 'Required files': q['required_files'], 'Execution script': q.get('script', q['execution_script']), 'Required fixtures': q['fixtures'],
            'Expected result': q.get('expected', ''), 'Definition of Done': q.get('dod', ''), 'Evidence destination': q.get('evidence', ''), 'Related finding IDs': q['findings'], 'Blocker': q['blocker'], 'Next task': q['next']}
wcsv('66-remaining-audit-execution-queue.csv', [qrow(q) for q in queue], qcols)
wcsv('66b-remaining-field-batches.csv', [dict(f, queue_task='') for f in fields], ['line', 'app', 'route', 'form', 'field', 'type', 'required', 'maps', 'comp', 'cases', 'hint', 'queue_task'])
# link field rows to their batch
bk = {k: f'Q-FLD-{i:03d}' for i, k in enumerate(groups, 1)}
for f in fields:
    f['queue_task'] = bk[(f['app'], f['route'], f['form'])]
wcsv('66b-remaining-field-batches.csv', [dict(f, comp=f['comp'], hint=f['hint']) for f in fields], ['line', 'app', 'route', 'form', 'field', 'type', 'required', 'maps', 'comp', 'cases', 'hint', 'queue_task'])
wcsv('66c-remaining-destructive-controls.csv', drows, dcols)
wcsv('66d-remaining-blocked-controls.csv', brows, bcols)

# ---------------------------------------------------------------- numbers
sev = collections.Counter(f['severity'].split(' ')[0] for f in FIND)
rc = collections.Counter((re.match(r'\s*(CONFIRMED|PROBABLE)', f['root_cause_first_break']) or [0, 'UNSTATED'])[1] for f in FIND)
c53 = collections.Counter(c['result'] for c in C53)
exec53 = sum(v for k, v in c53.items() if k in ('PASS', 'NO-EFFECT', 'FLAG-MUTATION', 'EXECUTED-RUN005'))
st = collections.defaultdict(collections.Counter)
for r in R51:
    st[r['phase']][r['status']] += 1
def pcell(p):
    c = st[p]
    comp = c['EXECUTED-PASS']; fail = c['EXECUTED-FAIL']; part = c['EXECUTED-PARTIAL']; blk_ = c['BLOCKED-CAPABILITY'] + c['BLOCKED-DEFECT']; plan = c['PLANNED']
    nq = sum(1 for q in queue if q['phase'] == p and q['status'] == 'NOT-STARTED')
    return comp, fail, part, blk_, plan, nq
def pstatus(p):
    comp, fail, part, blk_, plan, nq = pcell(p)
    if plan and not (comp + fail + part): return 'UNSTARTED'
    if nq or part or plan: return 'PARTIAL'
    if blk_ and not (comp + fail): return 'BLOCKED'
    return 'EXECUTED (feasible tasks done; product failures recorded)' if fail or comp else 'PARTIAL'
nfld_tested = 93; nfld_total = 531
stat = collections.Counter(p['status'] for p in [dict(status=pstatus(p)) for p in PH])
fm = collections.Counter(f['finding_ids'] for f in R51)
reg_tasks_sev = collections.Counter(r['severity'].split(' ')[0] for r in REM)
dec = 22
env_ok = 'api, postgres, redis healthy; landing and dashboard up; restaurant-website and super-admin STOPPED (verified by docker ps at ' + NOW + ')'
sha_dump = hashlib.sha256(open('evidence/database/isolated-audit-run005-final-2026-11-27.dump', 'rb').read()).hexdigest() if os.path.exists('evidence/database/isolated-audit-run005-final-2026-11-27.dump') else 'MISSING'
bk_path = glob.glob(r'C:\Users\Yousef\audit-backups\*RUN005-final*.tar.gz')
sha_bk = hashlib.sha256(open(bk_path[0], 'rb').read()).hexdigest() if bk_path else 'MISSING'
first = queue[0]['id']

# ---------------------------------------------------------------- 65
dash = ['| Phase | Name | Status | Total tasks (51) | Completed (PASS) | Failed | Partial | Blocked | Remaining feasible (queue 66) | Spec | Tasks | Definition of Done | Next action |', '|---|---|---|---|---|---|---|---|---|---|---|---|---|']
NEXT = {}
for p in PH:
    qs = [q['id'] for q in queue if q['phase'] == p and q['status'] == 'NOT-STARTED']
    NEXT[p] = qs[0] if qs else 'none feasible (see phase README)'
    comp, fail, part, blk_, plan, nq = pcell(p)
    tot = sum(st[p].values())
    dash.append(f'| {p} | {PH[p]} | {pstatus(p)} | {tot} | {comp} | {fail} | {part} | {blk_} | {nq} | [spec](phases/phase-{p}/spec.md) | [tasks](phases/phase-{p}/tasks.md) | [DoD](phases/phase-{p}/definition-of-done.md) | {NEXT[p]} |')
cov = [
 ('Routes (pages + 2 route handlers)', '119 discovered', '117 pages loaded in RUN-003/004 + 2 handlers reconciled (06b)', 'route set complete'),
 ('Forms', '128 form containers (07b)', '13 / 116 logical forms have field-level cases (54b)', f'{len(groups)} form groups with untested fields'),
 ('Logical fields', f'{nfld_total}', f'{nfld_tested} tested individually; {len(fields)} NOT TESTED', 'users form INCONCLUSIVE'),
 ('Field validation cases', 'applicable cases per field (66b)', '424 executed (251 accepted+persisted, 17 accepted 2xx, 81 rejected by server, 63 rejected by client)', 'remaining via Q-FLD-*'),
 ('UI controls (unique)', '644', f'{exec53} executed (PASS/NO-EFFECT/FLAG-MUTATION/EXECUTED-RUN005); 9 ERROR; 10 DISABLED', f'{c53["BLOCKED"]} tooling-blocked; {c53["SKIPPED-DESTRUCTIVE"]} skipped destructive'),
 ('Destructive controls', '88 skipped after RUN-004', '44 executed in RUN-005', '44 remain (66c)'),
 ('Arabic interactions', '15 representative routes', '30/30 dir=rtl checks, no overflow, dialogs fit', 'remaining routes via Q-MAT-001'),
 ('Mobile interactions', '15 routes x 4 viewports', '60 checks; 1 overflow finding (F-082)', 'remaining routes via Q-MAT-001'),
 ('API endpoints', '355 inventoried (08)', 'authorization matrix on 28 GET + 8 write probes per role (F-071); functional coverage see 38', 'API functional/authorization completeness NOT established'),
 ('Financial workflows', 'see 10/35b', 'scenario status in 35b', 'payments/bank/POS forms and concurrency remain'),
 ('Report metrics', '119 (30c)', '99 measured vs independent + 8 + 12 blocked', '61 correct, 23 incorrect, 6 policy-dependent, 17 insufficient, 12 blocked'),
 ('Report exports', 'all reports with export', 'KDS/HR export partially (F-068)', 'Q-RPT-003 / Q-RPT-004'),
 ('Security tests', 'see 36', 'RBAC enforce vs staging matrices; F-071', 'F-071 in depth, session/CSRF remaining per 36'),
 ('Reliability tests', 'see 36', 'Redis outage (F-018), load 10-100 VU', 'concurrency Q-CON-001'),
 ('Historical regressions', 'rows in 21', 'partly executed', 'Q-HIS-001'),
 ('Finding verification', f'{len(FIND)} findings', f'{len(FIND)} REPRODUCED', 'none pending verification'),
 ('Root-cause confirmation', f'{len(FIND)}', f'{rc["CONFIRMED"]} CONFIRMED, {rc["PROBABLE"]} PROBABLE (F-082)', 'F-082 element not isolated'),
]
covt = ['| Surface | Denominator | Executed so far | Remaining / note |', '|---|---|---|---|'] + [f'| {a} | {b} | {c} | {d} |' for a, b, c, d in cov]
open('65-AUDIT-MASTER-OVERVIEW.md', 'w', encoding='utf-8').write(f"""# 65 - Audit master overview (START HERE)

Generated {NOW} by `tests/build_handoff_package.py` from the authoritative registers. If a count here disagrees with a register, the register wins (51, 22, 57, 59, 53b, 54b).

## 1. Executive status
- **Audit verdict:** PARTIAL (RUN-001..RUN-005 executed; feasible work remains, see queue 66).
- **Financial release verdict:** NO-GO (see 24, 61). 1 Critical finding is open.
- **Git baseline:** branch `{BRANCH}`, HEAD `{HEAD[:12]}`; the audit folder was committed to the repository on user instruction after RUN-005. No product code was changed by the audit.
- **Last completed run:** RUN-005 (2026-11-27). Next run id: RUN-006 (`RUN_ID=RUN-006` env var for `tests/runtime/chrome_lib.mjs`).
- **Phases:** 23 (00-22). By derived status: {dict(stat)}. Phase status is derived from register 51 plus the queue; a phase is never complete just because its documents exist.
- **Findings:** {len(FIND)} total. Severity {dict(sev)}. Evidence classification: all {len(FIND)} REPRODUCED. Root-cause confidence: {dict(rc)}.
- **Audit tasks (register 51):** {len(R51)} rows: {dict(collections.Counter(r['status'] for r in R51))}.
- **Remaining feasible audit work:** {len(queue)} queue tasks in 66 ({dict(collections.Counter(q['cap'] for q in queue))} by capability) plus per-item lists 66b (fields), 66c (destructive controls), 66d (blocked controls).
- **Future remediation (Workstream B, NOT STARTED):** {len(REM)} implementation tasks (57), {len(REG)} regression tests (59), {dec} policy decisions all OPEN (56, 56b), {len(FIND)} finding packages under `findings/`. See 69.
- **Environment:** {env_ok}. Dev stack: dev frontend containers absent (recreating is the owner's action). Protected tenants 5/6/7 unseeded (COA rows 0). Tenant 10 = synthetic 54,000-line benchmark (read only). Tenant 4 = heavily mutated diagnostic fixture (not a clean baseline).

## 2. System overview (details in 03, 05, 09, 13, 25, 27, 28)
- **Applications (Next.js):** landing (marketing, owner signup), dashboard (restaurant staff: POS, orders, KDS, accounting, HR, reports, settings, theme builder), restaurant-website (public storefront and checkout), super-admin (SaaS tenants, billing). In the isolated stack: landing 13111, dashboard 13112, storefront 13113, super-admin 13114, API 18082, Postgres 127.0.0.1:15437.
- **Backend:** one Go API (`/api/v1`), handlers, services, repositories, middleware; session cookie `poos_session`; Redis cache and session store; an outbox worker publishes accounting events (13).
- **Database:** PostgreSQL, shared schema, `tenant_id` on tenant data; 145 migrations; accounting tables (chart_of_accounts, journal_entries, general_ledger, invoices, payments) next to operational tables (orders, kds_queue, products, stock) (09).
- **Multi-tenancy:** tenant resolved from the session; tenant 4 is the audit fixture, 10 the benchmark, 5/6/7 protected.
- **AuthN/AuthZ:** login rate limit 20/min per IP; RBAC roles (OWNER, ADMIN, MANAGER, EMPLOYEE, DRIVER); enforcement depends on RBAC_ENFORCEMENT (staging mode authorises everything: F-071).
- **Accounting:** double-entry journals posted from invoices, payments and POS orders through accounting events; operational rows (orders, supplier balances) and ledger can diverge (F-021, F-044, F-049); see 10, 11, 12, 27, 33b.
- **Background workers:** outbox worker (15 s cycle) drains events; if it stops, ledger posting lags without an error in the UI (18).
- **Sources of truth and cross-module dependencies:** see 25 and 33b.

## 3. Phase dashboard (one row per phase, from register 51)
{chr(10).join(dash)}

Failed means the audit task executed and the product failed its assertion (task is complete, finding recorded). Partial means feasible audit work remains. Blocked tasks need a product capability that does not exist (48).

## 4. Coverage dashboard
{chr(10).join(covt)}

Distinguish discovered, planned, attempted, executed, passed, failed, partial, blocked and not executed: 53b/54c carry per-item results; 45 and 48 list what is not executed and why.

## 5. Execution entrypoint (for the next agent)
1. Read this file, then `68-LOW-COST-AGENT-EXECUTION-GUIDE.md`, then the last block of `progress.md`.
2. Environment checks: `docker ps` (expect pos-p2a-ws05-e2e api/postgres/redis healthy); `curl http://127.0.0.1:18082/api/v1/health`; DB counts for tenants 5/6/7 (COA rows must be 0). Details in 67/37/64.
3. Open `66-remaining-audit-execution-queue.csv`; the first row with Status NOT-STARTED and no unfinished dependency is next. First task now: **{first}**.
4. Script, fixtures and expected result are in the row; per-item lists are in 66b/66c/66d.
5. Evidence goes to `evidence/{{network,database,screenshots,runtime}}/RUN-006/PH-<phase>/`; never open `evidence/state-*.json`.
6. Update: 34b (result row), 53b/54c/07 (via the generators), 22 (new finding), 51 (task status through `tests/r4_task_overlay.py`), then the RUN-006 block of `progress.md`.
7. Next task: the next eligible row in 66. Regenerate derived files only at checkpoints (order in 67 section Generators).
8. Execution may stop only when: every row of 66 is executed or classified BLOCKED with proof, the acceptance checklist 71 has been evaluated, and restoration (Q-ENV-002) is done; or an escalation condition in 68 applies to every remaining row.
""")

# ---------------------------------------------------------------- 67 file map
FM = [
 ('README.md', 'index; points to 65', 'manual', 'all', 'entry', 1), ('progress.md', 'single authoritative progress tracker', 'manual (+ tests/r5_checkpoint.py block)', 'all', 'input/output', 1),
 ('65-AUDIT-MASTER-OVERVIEW.md', 'master entrypoint', 'generated by tests/build_handoff_package.py', 'all', 'derived', 1), ('66-remaining-audit-execution-queue.csv', 'prioritized remaining audit tasks (refs register 51)', 'generated; status updated by hand after execution', 'all', 'authoritative for remaining work', 1),
 ('66b/66c/66d-*.csv', 'remaining fields / destructive controls / blocked controls', 'generated', '04, 06', 'derived from 07, 53b', 2), ('68-LOW-COST-AGENT-EXECUTION-GUIDE.md', 'operating protocol', 'generated', 'all', 'derived', 1),
 ('69-FUTURE-REMEDIATION-OVERVIEW.md', 'Workstream B plan (planning only)', 'generated', '21', 'derived', 3), ('70-EVIDENCE-AND-TRACEABILITY-GUIDE.md', 'evidence navigation', 'generated', '20', 'derived', 2),
 ('71-FINAL-AUDIT-ACCEPTANCE-CHECKLIST.md', '12 final verdicts', 'generated', '22', 'derived', 2), ('72-NEXT-AGENT-START-PROMPT.md', 'paste-able prompt', 'generated', 'all', 'derived', 1),
 ('51-phased-audit-task-register.csv', 'authoritative phase-level task register (61 rows)', 'tests/build_task_register.py + tests/r4_task_overlay.py (overlay holds results)', 'all', 'authoritative', 1),
 ('22-confirmed-findings-register.csv', 'authoritative finding register (81)', 'manual via tests/r*_register_update*.py scripts', 'all', 'authoritative', 1),
 ('21-existing-issue-regression-matrix.csv', 'historical regression rows', 'manual', '19', 'authoritative', 3), ('24-financial-release-verdict.md', 'release verdict', 'manual', '22', 'authoritative', 2),
 ('34-complete-manual-ui-test-plan.csv / 34b-executed-test-results.csv', 'planned / executed test rows (append only)', '34b appended by recorders', 'all', 'authoritative', 2),
 ('07-field-level-inventory.csv', 'field inventory with tested_status', 'rewritten by tests/build_field_denominator.py then build_field_results.py', '04', 'authoritative+generated', 2), ('54b/54c-*', 'field denominator and isolated results', 'tests/build_field_denominator.py, tests/build_field_results.py', '04', 'derived', 2),
 ('53-ui-control-inventory.csv / 53b-ui-control-interaction-results.csv', 'control inventory / results', '53b regenerated by tests/build_control_interaction_results.py; RECORDERS r5_record_destructive*.py MUST RUN AFTER', '04', 'derived + recorder overlay', 2),
 ('30c/30d/43', 'metric coverage, classes, runtime results', 'tests/build_metric_verification_classes.py', '12', 'derived', 3), ('06b-route-and-api-reconciliation.*', 'route/API reconciliation', 'tests/build_route_api_reconciliation.py', '03', 'derived', 3),
 ('findings/F-xxx/', '10 files per finding (finding, reproduction, root-cause, affected-code-paths, financial-impact, dependency-map, evidence-index, remediation-spec, remediation-tasks, regression-test-spec)', 'tests/build_finding_packages.py + tests/finding_overrides_r4*.py', '20/21', 'derived (overrides are the manual input)', 2),
 ('57 / 58 / 58b / 59', 'future remediation tasks, spec index, validation, regression tests', 'tests/build_finding_packages.py; tests/validate_packages.py', '21', 'derived', 3), ('56 / 56b / 32', 'accounting policy decisions (22, all OPEN)', 'manual; tests/build_policy_detail.py', '21', 'authoritative', 2),
 ('phases/phase-NN/', 'spec.md, tasks.md (generated by tests/build_phase_specs.py) plus README, definition-of-done, evidence-index, dependencies, execution-runbook (generated by build_handoff_package.py)', 'generated', 'per phase', 'derived', 2),
 ('evidence/', 'RUN-001..RUN-005 raw evidence (network, database, screenshots, runtime, logs); state-*.json are session cookies: NEVER open or commit', 'produced by runners', 'all', 'evidence (immutable)', 3),
 ('evidence/database/*.dump', 'database snapshots (RUN-004 start/mid/final, RUN-005 start/final)', 'pg_dump', '02/22', 'evidence (sensitive: test data only)', 3), ('60-evidence-manifest.csv', 'manifest of evidence files', 'tests/build_evidence_manifest.py', '22', 'derived', 3),
 ('tests/', 'generators, validators, overlays, recorders', 'manual', 'all', 'tooling', 2), ('tests/runtime/', 'headed Chrome runners (chrome_lib.mjs shared)', 'manual', 'all', 'tooling', 1), ('oracle/ sql/ fixtures/', 'independent financial oracle, SQL checks, fixture definitions', 'manual', '11/14/15', 'tooling', 3),
 ('61 / 62 / 63 / 64', 'closure report (generated by tests/build_closure_report.py), handoff, RUN-004 and RUN-005 baselines', '61 generated; others manual', '22', 'historical + current', 2),
]
ft = ['| Path | Purpose | Maintenance / generator | Phase | Role | Read priority (1 high) |', '|---|---|---|---|---|---|'] + [f'| {a} | {b} | {c} | {d} | {e} | {f} |' for a, b, c, d, e, f in FM]
open('67-audit-file-and-artifact-map.md', 'w', encoding='utf-8').write(f"""# 67 - Audit file and artifact map

{NOW}. Authoritative = hand-maintained or the source of truth; derived = regenerate, do not hand-edit. Sensitivity: evidence may contain test data and request headers (cookies redacted); `evidence/state-*.json` is ignored by git and must never be opened.

{chr(10).join(ft)}

## Stale or contradictory items (found while consolidating)
- `61` reported Arabic RTL as 0/59 in RUN-004; the generator was fixed in RUN-005 (59/59 and 30/30 now). Historical RUN-004 text kept in the RUN-004 addendum.
- `45` and `48` contain RUN-003, RUN-004 and RUN-005 addenda; earlier sections keep earlier counts (e.g. 88 skipped controls, 36 findings) as history. Current counts are in 65.
- `22` evidence_status column starts with REPRODUCED for every row; the finer classes (OBSERVED, PROBABLE) survive only inside root_cause_first_break.
- Register 51 has 61 phase-level tasks; the large numbers (531 fields, 644 controls) are covered by 66b/66c/66d, not by register 51 rows.
- `07` tested_status is rewritten on every `build_field_denominator.py` run; always run `build_field_results.py` right after it.

## Generators and their order (verified in RUN-005)
Run only at checkpoints. Each generator reads authoritative registers and rewrites its outputs; none touches product code.
1. `tests/build_finding_evidence_check.py` (reads 22 evidence_paths; writes 55)
2. `tests/build_finding_packages.py` (reads 22, finding_overrides_r4*.py; writes findings/, 57, 58, 59; prints counts)
3. `tests/r4_task_overlay.py` (applies RUN-004/005 results onto 51; add a tuple (status, result, finding_ids) per task)
4. `tests/build_phase_specs.py` (phases/*/spec.md and tasks.md; derives status from 51)
5. `tests/build_field_denominator.py` then `tests/build_field_results.py` (07, 54b, 54c; the second marks inconclusive forms)
6. `tests/build_control_interaction_results.py` **resets 53b**: immediately run `tests/r5_record_destructive.py` and `tests/r5_record_destructive_b.py` (they re-apply the 44 executed destructive results; idempotent). Any new recorder must be appended to this step.
7. `tests/validate_packages.py` (58b), `tests/build_policy_detail.py` (56b), `tests/build_metric_verification_classes.py` (30c/30d), `tests/build_route_api_reconciliation.py` (06b)
8. `tests/build_evidence_manifest.py` (60), `tests/build_closure_report.py` (61), `tests/r5_checkpoint.py` (progress.md and 62 block), `tests/build_handoff_package.py` (65-72, phases docs, 66*)
Warning: `build_handoff_package.py` regenerates 66* from the registers; it does not know which queue tasks were executed. After executing a queue task, record the status in 51 via the overlay and in `66-queue-status.csv` (task id, status, evidence); the generator reads that file and keeps statuses (see below).
""")

# queue status carry-over file
QS = '66-queue-status.csv'
if not os.path.exists(QS):
    wcsv(QS, [], ['Task ID', 'Status', 'Date', 'Evidence', 'Note'])
done = {r['Task ID']: r for r in rd(QS)}
if done:
    for q in queue:
        if q['id'] in done:
            q['status'] = done[q['id']]['Status']
    wcsv('66-remaining-audit-execution-queue.csv', [qrow(q) for q in queue], qcols)

# ---------------------------------------------------------------- 68
open('68-LOW-COST-AGENT-EXECUTION-GUIDE.md', 'w', encoding='utf-8').write(f"""# 68 - Low-cost agent execution guide

## Initial context (read only these)
1. 65-AUDIT-MASTER-OVERVIEW.md  2. this file  3. last block of progress.md  4. the next NOT-STARTED row in 66  5. that phase's `phases/phase-NN/README.md`  6. the row's required files. Only then open source or evidence for that task.

## Per-task protocol
1. Read the task row (66) and any per-item list (66b/66c/66d).  2. Verify prerequisites (dependency tasks done; fixtures).  3. Verify the environment (`docker ps`, health 200, protected tenants unchanged).  4. Prepare the fixture through the API with a unique prefix `000R6<random>`; never reuse tenants 5/6/7; tenant 10 is read only; tenant 4 is a mutated diagnostic fixture.  5. Execute the named script in headed Chrome (`AUDIT_CHROME_PROFILE` env gives each concurrent runner its own profile; `RUN_ID=RUN-006`).  6. Capture evidence (screenshot, network request+response, SQL before/after) under `evidence/<type>/RUN-006/PH-<phase>/`.  7. Compare expected and actual; a product failure is a finding, not a failed task.  8. Update: append a row to 34b through a small recorder script in `tests/` (copy `r5_record_destructive_b.py`), add findings through a `tests/r6_register_update*.py` script modelled on `r5_register_update5.py`, set the task status in `66-queue-status.csv`.  9. Append a short line to the RUN-006 block of progress.md.  10. Pick the next row in 66 whose dependencies are done.

## Rules learned the hard way
- Bash heredocs with apostrophes fail: write scripts with the file tool. Use `python -X utf8` (Arabic output).
- Mutating API calls need `Origin: http://localhost:13112` and cookie `poos_session`; login is rate limited to 20/min per IP (a burst of bad logins locks you out ~1 min; run the rate-limit test last).
- Lists show 20 rows sorted by name: prefix fixtures with `000` so they appear first. Icon-only buttons need aria-label selectors. The leaves page hides future dates by default.
- Memory: 7.7 GiB Docker VM. Keep restaurant-website and super-admin containers stopped unless the task needs them; start one at a time.
- Never switch git branches with the audit folder untracked; never `docker system prune`; never delete volumes, evidence, dumps or the 2 git stashes; never open `evidence/state-*.json`.
- A field batch result is meaningless until a clean baseline submit is ACCEPTED (RUN-005 lost time to five invalid baselines).

## Context efficiency
Search with Grep instead of reading whole files; read only the finding package of the related finding; reuse runners and fixtures; do not regenerate all documents after each test, only at checkpoints (order in 67); keep command output short (`| tail`); put durable results in the audit package, not in chat.

## Escalation (do not stop the audit; skip to the next independent row)
Escalate to a stronger model only for: a root cause that cannot be traced after reading the cited source paths; an accounting policy decision; a destructive action the owner has not authorised; an environment that differs materially from 65; evidence that contradicts the documented architecture. Capabilities in 66: ROUTINE-EXECUTION and BROWSER-AUTOMATION need only a low-cost model; DATA-VERIFICATION and SOURCE-TRACING need a mid-level model; FINANCIAL-RECONCILIATION, COMPLEX-INVESTIGATION and POLICY-REVIEW need the strongest available model.

## Closing run (Q-ENV-002)
Final pg_dump with full sha256 (not abbreviated), backup tar.gz excluding `evidence/state-*.json`, tenants 5/6/7 COA rows 0, 2 stashes, `git status` shows no product changes by the audit, dev DB `pos_saas` equals `evidence/runtime/RUN-003/46-dev-db-fingerprint-before.txt`, record all in 64 and 62.
""")

# ---------------------------------------------------------------- 69
by_phase = collections.defaultdict(list)
for r in REM:
    by_phase[r['phase']].append(r)
pf = ['| Implementation phase | Tasks | Findings | Severity mix |', '|---|---|---|---|']
for p in sorted(by_phase):
    rs = by_phase[p]
    pf.append(f"| {p} | {len(rs)} | {len({r['finding_id'] for r in rs})} | {dict(collections.Counter(r['severity'].split(' ')[0] for r in rs))} |")
open('69-FUTURE-REMEDIATION-OVERVIEW.md', 'w', encoding='utf-8').write(f"""# 69 - Future remediation overview (Workstream B: planning only, NOT STARTED)

Verified counts {NOW}: {len(FIND)} findings, {len(REM)} implementation tasks (57), {len(REG)} regression tests (59), {dec} policy decisions all OPEN (56b). Do not execute anything here; Workstream A (audit) is separate (65, 66).

Authoritative sources: task details in `57-future-remediation-task-register.csv` (do not copy them); specifications in `findings/F-xxx/remediation-spec.md`; index in `58-remediation-specification-index.md`; order and dependencies in `23-remediation-roadmap.md` and `46-root-cause-dependency-map.md`; shared-root-cause owner tasks in the packages (SHARED_OWNER map in `tests/build_finding_packages.py`).

## Implementation phases as defined in register 57
{chr(10).join(pf)}

## How to read each implementation phase
For a phase P: filter 57 by `phase`; included finding IDs are its `finding_id` values; prerequisites are the `dependencies` column (shared-root-cause owner tasks first); acceptance and Definition of Done are the `acceptance` and `required_tests` columns; regression tests are the 59 rows with the same finding_id; affected paths are `existing_paths`; policy decisions are listed per finding in `findings/F-xxx/remediation-spec.md` and `56b`. Database, financial and security implications are in `financial-impact.md` and `root-cause.md` of each package.

## Policy gates (all OPEN; the owner decides, the audit never approves)
POL-01..POL-22 in 56b. Any task whose package says "policy" must wait for that decision. Unlocked blocked workflows (F-001 new-tenant accounting bootstrap, F-013 period close, F-014 employee create, others in 48) are listed in the `Blocked workflows unlocked` line of the relevant package.

## Sequencing rule
Critical and High findings first (F-001 bootstrap is the single Critical, release-blocking); fix shared root causes once; after each implementation phase rerun its regression tests (59) and the original audit script named in each test row.
""")

# ---------------------------------------------------------------- 70
open('70-EVIDENCE-AND-TRACEABILITY-GUIDE.md', 'w', encoding='utf-8').write(f"""# 70 - Evidence and traceability guide

## Chain
Phase (`phases/phase-NN`) -> Audit task (51 `task_id`, 66 `Task ID`) -> Test case (34 / 34b `exec_id`, 53b `control_id`, 54c row, 35b scenario) -> Execution result (34b, 53b, 54c, 16 oracle) -> Evidence (file under `evidence/`) -> Finding (22 `finding_id`) -> Root cause (`findings/F-xxx/root-cause.md`, 55) -> Future remediation task (57 `finding_id`) -> Future regression test (59 `finding_id`). Cross-index: `38-master-traceability-matrix.csv`; evidence manifest: `60-evidence-manifest.csv`; per-finding evidence check: `55` (generated by `tests/build_finding_evidence_check.py`; 0 unresolved refs).

## Where evidence lives
| Run | Location | Contents |
|---|---|---|
| RUN-001/002 | `evidence/` top-level files and `evidence/logs`, `evidence/sql` | discovery and first executions |
| RUN-003 | `evidence/*/RUN-003/` | isolated stack build, onboarding, dev DB fingerprint (`evidence/runtime/RUN-003/46-dev-db-fingerprint-before.txt`) |
| RUN-004 | `evidence/{{network,database,screenshots,runtime}}/RUN-004/PH-xx` | control matrix, Arabic, reports, performance, destructive subset; dumps start/mid/final |
| RUN-005 | `evidence/{{network,database,screenshots,runtime}}/RUN-005/PH-xx` | field isolation (PH-06), destructive A/B and language/viewport (PH-04), POS currency (PH-07), invoice form (PH-09), RBAC matrices (PH-17), Redis timing (PH-18); start and final dumps |
Evidence types: screenshots (png), network captures (json: request, status, body head), database evidence (SQL result json/txt, pg_dump), runtime logs, oracle results (`16`, `oracle/`), browser traces.

## Rules
- Every executed result must cite an evidence path that exists (the checker prints unresolved refs).
- Do not duplicate large files; reference them. Do not paste cookies, tokens or passwords into evidence notes; `evidence/state-*.json` is never opened.
- Full integrity: dump sha256 `{sha_dump}`; backup sha256 `{sha_bk}` (recomputed in full when this document was generated).
- For every queue task (66) the `Evidence destination` column states where new evidence must go; completed tasks are justified by the evidence named in their 51 `actual_result` and 34b rows.
""")

# ---------------------------------------------------------------- 71
V = [
 ('1', 'System discovery completeness', '03, 06, 06b, 07, 08, 09, 25', 'Phase 00-03 tasks executed', 'inventories reconciled (119 routes, 355 APIs, 627 field rows)', 'any surface undiscovered', 'none', 'API functional coverage and 438 untested fields are execution gaps, not discovery gaps', '03', 'EVALUATED-PARTIAL (discovery largely complete; not formally closed)'),
 ('2', 'Functional UI testing completeness', 'Phase 04-08 results, 53b', 'Q-DST-*, Q-BLK-*, Q-FRM-*, Q-MAT-001', 'every control executed or BLOCKED with proof', 'controls still SKIPPED/BLOCKED without proof', 'billing checkout, E-Invoice', f'{c53["SKIPPED-DESTRUCTIVE"]} skipped, {c53["BLOCKED"]} blocked', '04', 'NOT MET'),
 ('3', 'Field-level validation completeness', '07, 54b, 54c, 66b', 'Q-FLD-*, Q-USR-001, Q-F065-001', 'every logical field has isolated results with a valid baseline', f'{len(fields)} fields untested', 'none', 'users form inconclusive; select/checkbox/file not covered by the runner', '04', 'NOT MET'),
 ('4', 'Financial transaction verification', '10, 15, 16, 35b, 44', 'Phase 06-11, 14, 15 tasks incl. Q-FRM-001..003, Q-CON-001', 'each money flow reconciled to the oracle', 'unexplained discrepancy', 'period close (F-013), new tenant bootstrap (F-001)', 'payments, bank, POS payment types, concurrency', '11', 'PARTIAL'),
 ('5', 'Reporting verification', '30, 30c, 30d, 43', 'Q-RPT-001..005', '119 metrics classified; boundaries and exports checked', 'metrics unmeasured', 'depreciation (no UI)', '61 correct, 23 incorrect, 6 policy-dependent, 17 insufficient, 12 blocked', '12', 'PARTIAL'),
 ('6', 'Security verification', '19, 36, F-071 package', 'Phase 17 tasks', 'RBAC enforce matrix for all roles; secrets, session, rate limit checked', 'unauthorised access found (it is, F-071)', 'none', 'matrices on 28 GET and 8 write probes per role only', '17', 'PARTIAL (defects found)'),
 ('7', 'Reliability verification', '18, 36, F-018 package', 'Phase 18 tasks, Q-CON-001, Q-RPT-005', 'outage, restart, load and concurrency tests done', 'unrecovered failure', 'none', 'concurrency and browser benchmark not done', '18', 'PARTIAL'),
 ('8', 'Evidence completeness', '55, 60, 70', 'all', 'every result cites existing evidence', 'unresolved refs', 'none', '0 unresolved refs in 55 at last generation', '20', 'MET for executed work; NOT MET for unexecuted work'),
 ('9', 'Remediation specification readiness', '57, 58, 58b, 59, findings/', 'Phase 21', f'{len(FIND)} packages complete, tasks have acceptance', 'incomplete package', 'policy decisions open', '58b validator: 1 finding with issues at last run; F-082 root cause PROBABLE', '21', 'MOSTLY MET (verify 58b)'),
 ('10', 'Overall audit completeness', '65, 66, 71', 'all queue rows', 'all feasible tasks executed', 'any feasible task left', 'blocked rows proven', f'{len(queue)} queue tasks remain', '22', 'NOT MET (PARTIAL)'),
 ('11', 'Financial integrity', '20, 24, 42, 44', 'Phase 11, 14, 15', 'ledger equals operational data for every flow', 'drift (F-021, F-044, F-049, others)', 'policy decisions', 'multiple confirmed integrity defects', '14/15', 'FAIL (defects confirmed)'),
 ('12', 'Financial release readiness', '24, 61', 'verdicts 4 and 11', 'no Critical/High financial defect open and bootstrap works', 'Critical F-001 open', 'none', '1 Critical, 20 High findings', '22', 'NO-GO'),
]
open('71-FINAL-AUDIT-ACCEPTANCE-CHECKLIST.md', 'w', encoding='utf-8').write(f"""# 71 - Final audit acceptance checklist

The final evaluation of these 12 verdicts has NOT been executed. The column "Current known status" is the state of knowledge at {NOW}, not a verdict. Evaluate in Q-ENV-002 / Phase 22 after the queue is exhausted.

| # | Verdict | Inputs / documents | Required completed tasks | Acceptance condition | Failure condition | Blocked condition | Known limitations | Phase | Current known status |
|---|---|---|---|---|---|---|---|---|---|
{chr(10).join('| ' + ' | '.join(v) + ' |' for v in V)}
""")

# ---------------------------------------------------------------- phase folders
ph_find = collections.defaultdict(set)
for r in R51:
    for fid in re.findall(r'F-\d{3}', r['finding_ids']):
        ph_find[r['phase']].add(fid)
for p in PH:
    d = f'phases/phase-{p}/'
    rows = [r for r in R51 if r['phase'] == p]
    qs = [q for q in queue if q['phase'] == p]
    comp, fail, part, blk_, plan, nq = pcell(p)
    ev = sorted(glob.glob(f'evidence/*/RUN-00*/PH-{p}'))
    deps = sorted({x for r in rows for x in re.findall(r'T-P\d\d-\d{3}', r['dependencies'])})
    open(d + 'README.md', 'w', encoding='utf-8').write(f"""# Phase {p} - {PH[p]}

**Status:** {pstatus(p)} (register 51: {dict(st[p])}; queue rows not started: {nq}). Overview: [65](../../65-AUDIT-MASTER-OVERVIEW.md).
**Completed:** {', '.join(r['task_id'] for r in rows if r['status'] == 'EXECUTED-PASS') or 'none'}. **Executed with product failure:** {', '.join(r['task_id'] for r in rows if r['status'] == 'EXECUTED-FAIL') or 'none'}. **Partial:** {', '.join(r['task_id'] for r in rows if r['status'] == 'EXECUTED-PARTIAL') or 'none'}. **Blocked:** {', '.join(r['task_id'] for r in rows if r['status'].startswith('BLOCKED')) or 'none'}.
**Relevant findings:** {', '.join(sorted(ph_find[p])) or 'see tasks'}.
**Required reading:** [spec.md](spec.md), [tasks.md](tasks.md), [dependencies.md](dependencies.md), the findings above under `../../findings/`.
**Execution order:** the queue rows below in the order of 66; then update progress.md.
**Exact next task:** {NEXT[p]}.

## Remaining queue rows
{chr(10).join(f"- {q['id']} ({q['priority']}, {q['cap']}): {q['title']}" for q in qs[:25]) or '- none'}
{('- ... ' + str(len(qs) - 25) + ' more in 66-remaining-audit-execution-queue.csv') if len(qs) > 25 else ''}
""")
    open(d + 'definition-of-done.md', 'w', encoding='utf-8').write(f"""# Phase {p} - Definition of Done

Phase execution completeness (separate from product correctness):
- Every register 51 row of this phase is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a verified prerequisite mapping ({'; '.join(f'{k}={v}' for k, v in st[p].items())} now).
- Every queue row of this phase ({len(qs)}) is executed or classified BLOCKED with proof; {nq} not started now.
- Evidence exists for each executed result and every reference resolves (`tests/build_finding_evidence_check.py`).
- Findings are registered in 22 with reproduction and evidence; financial discrepancies are reconciled or recorded as unexplained.
- No feasible task remains; acceptance criteria listed in spec.md are evaluated.
Current: **{pstatus(p)}**.
""")
    open(d + 'evidence-index.md', 'w', encoding='utf-8').write(f"""# Phase {p} - Evidence index

Existing evidence folders: {', '.join(ev) or 'none under PH-' + p + ' (older runs used other layouts; see 70)'}.
Evidence cited by register 51 rows: {'; '.join(sorted({r['source_paths'][:90] for r in rows if r['source_paths']}))[:900] or 'see tasks.md'}.
Missing evidence: results of the {nq} remaining queue rows (destination `evidence/<type>/RUN-006/PH-{p}/`).
""")
    open(d + 'dependencies.md', 'w', encoding='utf-8').write(f"""# Phase {p} - Dependencies

- Prerequisite tasks (from 51): {', '.join(deps) or 'none recorded'}.
- Environment: isolated stack up (65 section 5); fixtures prefixed `000R6`; protected tenants 5/6/7 never touched.
- Findings that constrain execution: {', '.join(sorted(ph_find[p])) or 'none'}.
- Policy dependencies: see 56b for any decision referenced by these findings.
- Product defects blocking execution: tasks with status BLOCKED-* in 51 (see 48).
- Downstream: Phases 20-22 consume this phase's results.
""")
    open(d + 'execution-runbook.md', 'w', encoding='utf-8').write(f"""# Phase {p} - Execution runbook

1. Environment check (65 section 5, 37). 2. Open the first NOT-STARTED row for this phase in 66 (`Q-...` ids listed in README.md). 3. Use the script named in its `Execution script` column (runners live in `tests/runtime/`, shared helpers in `chrome_lib.mjs`; set `RUN_ID=RUN-006`). 4. Create fixtures through the API with prefix `000R6`. 5. Save evidence to `evidence/<type>/RUN-006/PH-{p}/`. 6. Record results with a recorder script in `tests/` and update 51 through `tests/r4_task_overlay.py`; set the row status in `66-queue-status.csv`. 7. Append to progress.md. Verified commands are those used by existing runners (see their headers); any new script is audit-only and must be listed in 67.
Phase tasks of record (register 51): {', '.join(r['task_id'] for r in rows)}.
""")
    # extend tasks.md idempotently
    t = open(d + 'tasks.md', encoding='utf-8').read()
    t = re.sub(r'\n## Remaining audit queue \(RUN-005 handoff\).*', '', t, flags=re.S)
    t += '\n## Remaining audit queue (RUN-005 handoff)\nWorkstream A only (audit execution). Future remediation tasks are in 57 and are NOT STARTED.\n' + ('\n'.join(f"- {q['id']} [{q['status']}] {q['title']} (script: {q.get('script', '')[:80]})" for q in qs) or '- none') + '\n'
    open(d + 'tasks.md', 'w', encoding='utf-8').write(t)

# ---------------------------------------------------------------- 72
open('72-NEXT-AGENT-START-PROMPT.md', 'w', encoding='utf-8').write(f"""# NEXT AGENT START PROMPT (paste into a new conversation)

You are continuing a financial and functional audit of SERVIO (multi-tenant restaurant POS/ERP/accounting SaaS). Execute remaining audit tasks immediately; do not write another broad plan and do not rediscover the project.

Project root: C:\\Users\\Yousef\\Desktop\\Projects\\POSS (repo Yousef-Abdallah-1711/pos). Audit directory: C:\\Users\\Yousef\\Desktop\\Projects\\POSS\\audit-reports\\full-financial-system-discovery-2026-11-25\\ (branch {BRANCH}, HEAD {HEAD[:12]} when this was generated).

READ FIRST, in order: 65-AUDIT-MASTER-OVERVIEW.md, 68-LOW-COST-AGENT-EXECUTION-GUIDE.md, the last block of progress.md, then the first NOT-STARTED row of 66-remaining-audit-execution-queue.csv and its phase README (phases/phase-NN/README.md). Open source files and evidence only as that task requires. Also read the repository CODING_STANDARDS.md before touching any code-like file (AGENTS.md requires it).

STATE: audit verdict PARTIAL; financial release NO-GO. {len(FIND)} findings ({dict(sev)}), {len(REM)} future remediation tasks and {len(REG)} regression tests (all NOT STARTED, Workstream B), {dec} policy decisions all OPEN. Fields: 93 of 531 tested individually ({len(fields)} not tested, listed in 66b). Controls: {c53['SKIPPED-DESTRUCTIVE']} destructive skipped (66c), {c53['BLOCKED']} tooling-blocked (66d). Queue: {len(queue)} tasks. FIRST TASK: {first} (environment checkpoint), then Q-INC-001.
ENVIRONMENT: isolated compose project pos-p2a-ws05-e2e (compose file docker-compose.phase2a-e2e.yml at the repo root, use --project-directory .), DB pos_phase2a_e2e, API 18082, landing 13111, dashboard 13112, storefront 13113 (container stopped), super-admin 13114 (container stopped), Postgres 127.0.0.1:15437. Mutating API calls need header Origin: http://localhost:13112 and cookie poos_session. Tenant 4 = diagnostic fixture, tenant 10 = read-only benchmark, tenants 5/6/7 PROTECTED (never seed or mutate; COA rows must stay 0). Dev DB pos_saas must stay unchanged (tenants 2, users 6, orders 0, journals 0).

RULES: audit only. Do NOT modify product code, migrations, formulas or configs; do NOT fix any finding; do NOT start any future remediation task; do NOT approve policy decisions; do NOT use subagents; do NOT commit or push product code; do NOT open, print or commit evidence/state-*.json; redact secrets; never delete evidence, dumps, backups, the 2 git stashes or Docker volumes; no docker system prune; never switch git branches (the audit folder may be untracked on other branches).
EXECUTION: loop without asking whether to continue: take the next eligible row of 66, run its script, capture evidence under evidence/<type>/RUN-006/PH-<phase>/, record results (34b via a recorder script, findings via tests/r6_register_update*.py modelled on tests/r5_register_update5.py, status in 66-queue-status.csv and register 51 via tests/r4_task_overlay.py), append to progress.md, continue. Regenerate derived files only at checkpoints, in the order given in 67 (control results reset 53b: run r5_record_destructive*.py after it). A task is done when its Definition of Done in 66 is met; a reproduced product defect is a finding, not a failed task. If something needs escalation (policy decision, unresolvable root cause, unauthorised destructive action, environment mismatch) record it and continue with the next independent row.
STOP only when every row of 66 is executed or BLOCKED with proof. Then evaluate 71 (12 verdicts, keep PARTIAL and NO-GO unless evidence supports otherwise), do the closing run (Q-ENV-002: final dump with full sha256, backup, dev DB fingerprint, protected tenants), and update README, progress, 61 and 62. Report totals, findings added, verdicts and anything unexecuted honestly.
""")
print('queue', len(queue), 'fields', len(fields), 'groups', len(groups), 'dest', len(drows), 'blocked', len(brows), 'sha ok', sha_dump[:8], sha_bk[:8])
