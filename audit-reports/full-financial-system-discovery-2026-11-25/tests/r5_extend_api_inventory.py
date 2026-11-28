"""RUN-005: appends the API routes registered in source but missing from 08-api-inventory.csv (idempotent). Source location is exact; wrappers/module are read from the registration line."""
import csv, os, re, glob
REPO = r'C:\Users\Yousef\Desktop\Projects\POSS'
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
rows = list(csv.reader(open('08-api-inventory.csv', encoding='utf-8')))
hd = rows[0]
have = {re.sub(r'\{[^}]*\}', '{}', f'{r[0]} {r[1]}') for r in rows[1:]}
src = {}
for p in glob.glob(os.path.join(REPO, 'backend', '**', '*.go'), recursive=True):
    if p.endswith('_test.go'):
        continue
    for i, line in enumerate(open(p, encoding='utf-8', errors='replace'), start=1):
        m = re.search(r'"((?:GET|POST|PUT|PATCH|DELETE) (/api/v1/[^"\s]*))"', line)
        if m:
            key = re.sub(r'\{[^}]*\}', '{}', m.group(1))
            src.setdefault(key, (m.group(1), os.path.relpath(p, REPO).replace('\\', '/'), i, line.strip()))
added = 0
for key, (lit, path, ln, line) in sorted(src.items()):
    if key in have:
        continue
    method, route = lit.split(' ', 1)
    wrap = re.search(r'(wrapWithPermission|wrapProtected|authenticateRequest|RequireTrustedOrigin|RateLimit\w*)', line)
    mod = re.search(r',\s*(\d+),\s*"(READ|WRITE|DELETE|EXPORT)"', line)
    rows.append([method, route, '(handler not enriched)', f'{path}:{ln}', ('wrappers=' + wrap.group(1)) if wrap else 'wrappers=unknown (read the registration)', (mod.group(1) + '=module id') if mod else '', 'not classified', f'ADDED RUN-005 from source: registered at {path}:{ln}'])
    added += 1
csv.writer(open('08-api-inventory.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
print('added', added, 'rows now', len(rows) - 1)
