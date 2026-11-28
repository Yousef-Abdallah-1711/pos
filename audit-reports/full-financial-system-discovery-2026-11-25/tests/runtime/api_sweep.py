import csv, re, json, subprocess, urllib.request, urllib.error, os, sys, collections
ROOT = r'c:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25'
API = 'http://127.0.0.1:18082'
routes = list(csv.DictReader(open(os.path.join(ROOT, '08-api-inventory.csv'), encoding='utf-8')))
def login(email):
    req = urllib.request.Request(API + '/api/v1/auth/login', data=json.dumps({'email': email, 'password': 'AuditPass#2026A'}).encode(), headers={'Content-Type': 'application/json', 'Origin': 'http://localhost:13111'})
    r = urllib.request.urlopen(req, timeout=20)
    for h in r.headers.get_all('Set-Cookie') or []:
        m = re.match(r'poos_session=([^;]+)', h)
        if m: return m.group(1)
cashier = login('cashier.c@example.test'); owner = login('audit.c.owner@example.test')
SKIP = re.compile(r'/auth/(login|register|logout|refresh)|/super-admin/auth/login|/subscription/cancel')
def sub(path):
    p = re.sub(r'\{(email)\}', 'x@example.test', path)
    p = re.sub(r'\{(slug)\}', 'audit-c-diag', p)
    p = re.sub(r'\{(number|orderNumber)\}', 'ORD-X', p)
    return re.sub(r'\{[^}]+\}', '999999', p)
def call(method, path, cookie):
    h = {'Origin': 'http://localhost:13112', 'Content-Type': 'application/json'}
    if cookie: h['Cookie'] = 'poos_session=' + cookie
    data = b'{}' if method in ('POST', 'PUT', 'PATCH') else None
    req = urllib.request.Request(API + path, data=data, headers=h, method=method)
    try:
        r = urllib.request.urlopen(req, timeout=15); return r.status
    except urllib.error.HTTPError as e: return e.code
    except Exception: return -1
rows = []
for r in routes:
    path = sub(r['path']); m = r['method']
    if SKIP.search(r['path']): rows.append([m, r['path'], 'SKIPPED', 'SKIPPED', 'SKIPPED', r['module']]); continue
    rows.append([m, r['path'], call(m, path, None), call(m, path, cashier), call(m, path, owner), r['module']])
with open(os.path.join(ROOT, 'evidence', 'api-sweep-results.csv'), 'w', newline='', encoding='utf-8') as f:
    w = csv.writer(f); w.writerow(['method', 'path', 'no_token', 'cashier', 'owner_invalid_id_empty_body', 'module']); w.writerows(rows)
c = collections.Counter();
for r in rows: c[('none', r[2])] += 1
print('routes', len(rows)); print('no_token', collections.Counter(str(r[2]) for r in rows)); print('cashier', collections.Counter(str(r[3]) for r in rows)); print('owner', collections.Counter(str(r[4]) for r in rows))
