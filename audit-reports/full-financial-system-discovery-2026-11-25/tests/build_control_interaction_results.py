"""Builds 53b-ui-control-interaction-results.csv/.md from the RUN-004 control matrices (forward + reverse runners, dashboard EN desktop, owner, tenant 4) and the Super Admin run. Audit tooling; idempotent."""
import csv, json, os, collections, hashlib
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
P = 'evidence/network/RUN-004/'
sources = [('dashboard', P + 'PH-04/AUD-R4-CONTROL-MATRIX.json'), ('dashboard', P + 'PH-04/AUD-R4-CONTROL-MATRIX-subset.json'), ('super-admin', P + 'PH-13/AUD-R4-SUPERADMIN-CONTROL-MATRIX.json')]
FIND = {('reports/financial/advanced', 'Refresh'): 'F-009', ('accounting/advanced', 'Create Schedule'): 'F-074', ('settings/appearance', 'Theme Colors'): 'F-077', ('/dashboard/users', 'View'): 'F-076'}
rows = {}
for app, f in sources:
    if not os.path.exists(f):
        continue
    per = collections.Counter()
    for r in json.load(open(f, encoding='utf-8'))['results']:
        route = r['route'] or '(home)'
        per[(route, r['control'])] += 1
        rk = (app, route, r['control'], r.get('type'), per[(route, r['control'])])
        rows.setdefault(rk, dict(r, _src=f, _app=app))  # first source (forward runner) wins for overlapping routes
# keep the first occurrence per (app, route, control, ordinal); forward file wins
out = []
seen = set(); ordinal = collections.Counter()
for (app, route, control, typ, n), r in rows.items():
    cid = 'CTL-' + hashlib.sha1(f'{app}|{route}|{control}|{n}'.encode()).hexdigest()[:8].upper()
    if cid in seen:
        continue
    seen.add(cid)
    label = control.split(':', 1)[-1]
    fid = ''
    for (rt, lb), f_ in FIND.items():
        if rt in route and label.startswith(lb):
            fid = f_
    if app == 'super-admin' and route == 'users' and label.startswith('View') and r['status'] == 'ERROR':
        fid = 'F-076'
    st = r['status']
    expected = 'Control responds with an observable change (navigation, dialog, filter, download) and no console error or 5xx; state-changing controls need a purpose-built fixture' if st != 'SKIPPED-DESTRUCTIVE' else 'Not executed: label implies money movement or state change; needs a fixture and an oracle'
    actual = '; '.join(filter(None, [st, r.get('detail') or '', ('navigated to ' + r['navigated']) if r.get('navigated') else '', ('dialog opened' if r.get('dialog') else ''), ('mutations: ' + ', '.join(r['mutations'])) if r.get('mutations') else '', ('errors: ' + ' | '.join(r['errors'])[:160]) if r.get('errors') else '', ('bad GET: ' + ', '.join(r['badGet'])) if r.get('badGet') else '']))
    blocker = 'tooling: click timeout or tab crash (not attributed to the product)' if st == 'BLOCKED' else ('disabled in this state' if st == 'DISABLED' else ('needs fixture' if st == 'SKIPPED-DESTRUCTIVE' else ''))
    out.append([cid, app, route, r.get('type', ''), label, 'owner', 'tenant 4 fixture, EN desktop 1440x900, headed Chrome', expected, actual, st, os.path.basename(r['_src']), fid, blocker])
cols = ['control_id', 'app', 'route', 'control_type', 'label', 'role', 'preconditions', 'expected', 'actual', 'result', 'evidence', 'finding_id', 'blocker']
with open('53b-ui-control-interaction-results.csv', 'w', newline='', encoding='utf-8') as h:
    w = csv.writer(h); w.writerow(cols); w.writerows(out)
cnt = collections.Counter(x[9] for x in out); byapp = collections.Counter(x[1] for x in out)
routes = {(x[1], x[2]) for x in out}
open('53b-ui-control-interaction-results.md', 'w', encoding='utf-8').write(f"""# 53b - UI control interaction results (RUN-004, headed Chrome)

Data: [53b-ui-control-interaction-results.csv](53b-ui-control-interaction-results.csv). Built by `tests/build_control_interaction_results.py` from the matrices produced by `tests/runtime/ph04_r4_control_matrix.mjs` (dashboard) and `tests/runtime/ph13_r4_superadmin_ui.mjs` (Super Admin).

**What this is:** one safe interaction per control (tab, button, link, select, toggle, search/filter), up to 22 controls per route, on {len({r for a, r in routes if a == 'dashboard'})} of 59 dashboard routes and {len({r for a, r in routes if a == 'super-admin'})} Super Admin routes; {len(out)} controls in total. **What it is not:** an exhaustive functional test. Destructive or money-moving labels were not clicked (SKIPPED-DESTRUCTIVE: need a fixture and an oracle); the 22-per-route cap leaves lower-priority controls on long pages untested; Arabic and mobile interaction was not repeated (only the load matrix of RUN-003).

Result counts: {dict(cnt)}. By app: {dict(byapp)}.

Definitions: PASS = observable change without errors; NO-EFFECT = no observable change in the page signature (may still be correct, for example a toggle that re-renders identically); FLAG-MUTATION = the click caused a non-GET request (reviewed individually: exports use POST, Renew/Acknowledge/Mark-all-read/Create Session/Create Sample Template change fixture data, Appearance colours overwrite the active public theme - F-077); ERROR = console error or 5xx (F-009, F-074); BLOCKED = tooling failure (click timeout, Chrome tab crash), not a product verdict; DISABLED = control disabled in the fixture state.

Findings linked: F-009 (menu-profitability and break-even 500 still reproduced through the UI Refresh button), F-074 (unhandled 400 on Create Schedule), F-076 (Super Admin user View links 404), F-077 (Appearance Theme Colors rewrites the active storefront theme).
""")
print(len(out), dict(cnt))
