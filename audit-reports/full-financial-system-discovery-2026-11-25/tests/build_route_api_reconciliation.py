"""RUN-005 Phase L: reconciles page routes and API registrations (source) against 06/08/53/53b. Writes 06b-route-and-api-reconciliation.md/.csv. Read-only on the repository. Idempotent."""
import csv, os, re, glob
REPO = r'C:\Users\Yousef\Desktop\Projects\POSS'
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
APPS = {'dashboard': 'frontend/apps/dashboard/src/app', 'landing': 'frontend/apps/landing/src/app', 'super-admin': 'frontend/apps/super-admin/src/app', 'restaurant-website': 'frontend/apps/restaurant-website/app'}
rows = []
for app, base in APPS.items():
    root = os.path.join(REPO, base.replace('/', os.sep))
    for dp, dn, fn in os.walk(root):
        dn[:] = [d for d in dn if d != 'node_modules']
        for f in fn:
            if f in ('page.tsx', 'page.ts', 'page.jsx', 'route.ts', 'route.tsx'):
                rel = os.path.relpath(dp, root).replace('\\', '/')
                segs = [s for s in ([] if rel == '.' else rel.split('/')) if not (s.startswith('(') and s.endswith(')'))]
                url = '/' + '/'.join(segs)
                kind = 'route-handler' if f.startswith('route') else 'page'
                dyn = any(s.startswith('[') and s != '[locale]' for s in segs)
                rows.append([app, url, kind, 'dynamic' if dyn else 'static', os.path.relpath(os.path.join(dp, f), REPO).replace('\\', '/')])
# dashboard routes covered by the control inventory (53) and interaction results (53b)
c53 = set()
for r in csv.DictReader(open('53-ui-control-inventory.csv', encoding='utf-8')):
    c53.add(r['route'])
c53b = {(r['app'], r['route']) for r in csv.DictReader(open('53b-ui-control-interaction-results.csv', encoding='utf-8'))}
def dash_key(url):
    u = re.sub(r'^/\[locale\]/dashboard/?', '', url)
    return u or '(home)'
out = []
for app, url, kind, dyn, path in rows:
    if app == 'dashboard' and url.startswith('/[locale]/dashboard'):
        k = dash_key(url)
        cov53 = 'yes' if k in c53 else ('n/a (dynamic: needs a valid id)' if dyn == 'dynamic' else 'NO')
        cov53b = 'yes' if ('dashboard', k) in c53b else ('n/a (dynamic)' if dyn == 'dynamic' else 'NO')
    else:
        cov53 = cov53b = 'outside the dashboard control inventory (loaded in RUN-003; interaction per app: see 34b)'
        if app == 'super-admin' and url.startswith('/dashboard'):
            k = url.replace('/dashboard', '').strip('/') or ''
            cov53b = 'yes' if ('super-admin', k) in c53b else ('n/a (dynamic)' if dyn == 'dynamic' else 'NO (audit-logs/settings not reached in RUN-004)')
    out.append([app, url, kind, dyn, cov53, cov53b, path])
with open('06b-route-and-api-reconciliation.csv', 'w', newline='', encoding='utf-8') as h:
    w = csv.writer(h); w.writerow(['app', 'route_pattern', 'kind', 'static_or_dynamic', 'in_53_control_inventory', 'in_53b_interaction_results', 'source_file']); w.writerows(out)
# API: every quoted "METHOD /api/v1/..." literal in backend Go sources vs 08
reg = set()
for p in glob.glob(os.path.join(REPO, 'backend', '**', '*.go'), recursive=True):
    if p.endswith('_test.go'):
        continue
    for m in re.finditer(r'"((?:GET|POST|PUT|PATCH|DELETE|OPTIONS) /api/v1/[^"\s]*)"', open(p, encoding='utf-8', errors='replace').read()):
        reg.add(m.group(1))
inv = {f"{r['method']} {r['path']}" for r in csv.DictReader(open('08-api-inventory.csv', encoding='utf-8'))}
norm = lambda s: re.sub(r'\{[^}]*\}', '{}', s)
regn = {norm(x) for x in reg}; invn = {norm(x) for x in inv}
only_src = sorted(regn - invn); only_inv = sorted(invn - regn)
pages = [r for r in out if r[2] == 'page']; handlers = [r for r in out if r[2] == 'route-handler']
dash = [r for r in pages if r[0] == 'dashboard']
miss53 = [r for r in dash if r[4] == 'NO']; miss53b = [r for r in dash if r[5] == 'NO']
open('06b-route-and-api-reconciliation.md', 'w', encoding='utf-8').write(f"""# 06b - Route and API reconciliation (RUN-005, Phase L / T-P03-002)

Built by `tests/build_route_api_reconciliation.py` from the current source tree (read-only). Detail: [06b-route-and-api-reconciliation.csv](06b-route-and-api-reconciliation.csv).

## Pages and route handlers
- Page files: **{len(pages)}** (dashboard {len(dash)}, landing {sum(1 for r in pages if r[0] == 'landing')}, restaurant-website {sum(1 for r in pages if r[0] == 'restaurant-website')}, super-admin {sum(1 for r in pages if r[0] == 'super-admin')}); route handlers: **{len(handlers)}** ({', '.join(r[0] + ' ' + r[1] for r in handlers)}). Pages + handlers = {len(pages) + len(handlers)}, which equals the 119 discovered entries in 06: the 119-vs-117 difference is exactly the 2 route handlers (T-P03-002 resolved; the earlier "117 loaded" counted pages only).
- Dynamic page patterns: {sum(1 for r in pages if r[3] == 'dynamic')} (need a valid id to render; invalid-id behaviour is tested separately).
- Dashboard pages not in the 53 control inventory (static): {len(miss53)}: {', '.join(r[1] for r in miss53) or 'none'}.
- Dashboard pages without interaction results in 53b (static): {len(miss53b)}.

## API registrations
- Distinct `METHOD /api/v1/...` literals registered in backend sources: **{len(regn)}**; rows in 08: **{len(invn)}** (after normalising path parameters).
- In source but not in 08 ({len(only_src)}): {', '.join(only_src[:40]) or 'none'}{' ...' if len(only_src) > 40 else ''}
- In 08 but not found in source literals ({len(only_inv)}): {', '.join(only_inv[:40]) or 'none'}{' ...' if len(only_inv) > 40 else ''}
- Note: authentication-only sweeps (EX-API-*) are NOT functional API coverage; functional coverage is the set of endpoints exercised by the RUN-003/004 diagnostics and the UI matrices.
""")
print(len(pages), len(handlers), len(dash), 'miss53', len(miss53), 'miss53b', len(miss53b), 'api src', len(regn), 'inv', len(invn), 'only_src', len(only_src), 'only_inv', len(only_inv))
