"""Builds 60-evidence-manifest.csv: path, bytes, sha256, kind, run, phase for every file under evidence/ (and tests/, oracle/, sql/).
Session-state files (evidence/state-*.json) are listed but NEVER opened or hashed. Audit tooling; idempotent."""
import csv, hashlib, os, re
AUDIT = r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25'
os.chdir(AUDIT)
rows = []
def kind(p):
    ext = os.path.splitext(p)[1].lower()
    return {'.png': 'screenshot', '.json': 'network/json', '.txt': 'text/sql-output', '.dump': 'pg_dump', '.mjs': 'script', '.py': 'script', '.sql': 'sql', '.csv': 'csv', '.md': 'markdown', '.log': 'log'}.get(ext, ext.lstrip('.') or 'other')
for root in ('evidence', 'tests', 'oracle', 'sql', 'findings', 'phases', 'diagrams'):
    for dp, dn, fn in os.walk(root):
        for f in sorted(fn):
            p = os.path.join(dp, f).replace('\\', '/')
            m = re.search(r'RUN-(\d+)/PH-(\d+)', p)
            run, ph = (m.group(1), m.group(2)) if m else ('', '')
            if re.match(r'evidence/state-.*\.json$', p):
                rows.append([p, '', 'EXCLUDED-NOT-READ', 'session-state (cookies)', '', ''])
                continue
            try:
                b = open(p, 'rb').read() if os.path.getsize(p) < 600_000_000 else None
            except Exception:
                b = None
            if b is None:
                rows.append([p, os.path.getsize(p), 'NOT-HASHED', kind(p), run, ph]); continue
            rows.append([p, len(b), hashlib.sha256(b).hexdigest(), kind(p), run, ph])
with open('60-evidence-manifest.csv', 'w', newline='', encoding='utf-8') as h:
    w = csv.writer(h); w.writerow(['path', 'bytes', 'sha256', 'kind', 'run', 'phase']); w.writerows(rows)
print(len(rows), 'files; excluded', sum(1 for r in rows if r[2] == 'EXCLUDED-NOT-READ'))
