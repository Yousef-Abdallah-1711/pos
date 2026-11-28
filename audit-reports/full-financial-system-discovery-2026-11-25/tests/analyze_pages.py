"""Read-only page analyzer: for each Next.js page, walk local imports (depth<=3) and extract API endpoint literals + permission checks."""
import os, re, json, sys, csv
ROOT = r"c:\Users\Yousef\Desktop\Projects\POSS\frontend\apps"
APPS = {
 'dashboard': ('src/app', 'src'),
 'landing': ('src/app', 'src'),
 'super-admin': ('src/app', 'src'),
 'restaurant-website': ('app', ''),
}
routes = json.load(open('tests/_routes_raw.json'))
def to_re(p):
    p = p.replace('/api/v1', '')
    p = re.sub(r'\{[^}]+\}', r'[^/]+', p)
    return re.compile('^' + p + '$')
route_res = [(x['method'], x['path'].replace('/api/v1', ''), to_re(x['path'])) for x in routes]
SEG = set('accounting hr orders admin reports products categories inventory kds tables pos public rbac users tenant subscription auth notifications super-admin suppliers me audit impersonation translations webhook health'.split())
LIT = re.compile(r"""(?:(\.(get|post|put|patch|delete|request)\s*(?:<[^>]*>)?\s*\(\s*)|(fetch\w*\(\s*)|[`'"])\s*(?:\$\{[^}]*\}|[A-Za-z_\.]*(?:API|BASE|URL|Url|url)[A-Za-z_]*\s*\+\s*)?[`'"]?(/(?:api/v1/)?[a-z][a-z\-]*(?:/[^`'"\s?)]*)?)""")
STR = re.compile(r"""[`'"](/(?:api/v1/)?([a-z][a-z\-]*)(/[^`'"\s?]*)?)(?:\?[^`'"]*)?[`'"]""")
IMP = re.compile(r"""(?:import|export)\s+(type\s+)?([^'";]*?)\s*from\s+['"]([^'"]+)['"]|import\(\s*['"]([^'"]+)['"]\s*\)""")
EXPB = re.compile(r"^export\s+(?:default\s+)?(?:async\s+)?(?:function\*?|const|class|let)\s+(\w+)", re.M)
def blocks(txt):
    idx = [(m.start(), m.group(1)) for m in EXPB.finditer(txt)]
    res = []
    prev = 0; prevname = None
    for st, nm in idx:
        res.append((prevname, txt[prev:st])); prev = st; prevname = nm
    res.append((prevname, txt[prev:]))
    return res
def select(txt, names):
    if names is None: return txt
    keep = []
    for nm, body in blocks(txt):
        if nm is None or nm in names: keep.append(body)
    return ''.join(keep)
def parse_names(clause):
    if clause is None: return None
    clause = clause.strip()
    if '*' in clause: return None
    m = re.search(r'\{([^}]*)\}', clause)
    if not m: return None
    out = set()
    for part in m.group(1).split(','):
        part = part.strip()
        if not part: continue
        part = re.sub(r'^type\s+', '', part)
        out.add(part.split(' as ')[0].strip())
    return out
PERM = re.compile(r"checkPermission\(\s*MODULE_IDS\.(\w+)\s*,\s*PERMISSION_LEVELS\.(\w+)")
PERM2 = re.compile(r"(?:hasPermission|canAccess|usePermission|requirePermission|PermissionGuard)\b[^\n]{0,100}")
def resolve(frm, spec, appdir, srcdir):
    cands = []
    if spec.startswith('@/'):
        base = os.path.join(appdir, srcdir, spec[2:]) if srcdir else os.path.join(appdir, spec[2:])
        cands.append(base)
    elif spec.startswith('.'):
        cands.append(os.path.normpath(os.path.join(os.path.dirname(frm), spec)))
    else:
        return None
    for b in cands:
        for ext in ('', '.ts', '.tsx', '/index.ts', '/index.tsx'):
            p = b + ext
            if os.path.isfile(p): return p
    return None
def scan(path, appdir, srcdir, depth, seen, names=None):
    key = (path, tuple(sorted(names)) if names else None)
    if key in seen or depth > 4: return
    seen.add(key)
    try: txt = open(path, encoding='utf-8', errors='ignore').read()
    except: return
    if '.test.' in path or '.spec.' in path: return
    use = select(txt, names)
    yield path, use
    for m in IMP.finditer(use):
        spec = m.group(3) or m.group(4)
        if not spec: continue
        p = resolve(path, spec, appdir, srcdir)
        if p:
            nm = parse_names(m.group(2)) if m.group(3) else None
            yield from scan(p, appdir, srcdir, depth + 1, seen, nm)
def extract(txt):
    eps = {}
    for m in STR.finditer(txt):
        full, seg, rest = m.group(1), m.group(2), m.group(3) or ''
        if seg not in SEG: continue
        p = full.replace('/api/v1', '')
        p = re.sub(r'\$\{[^}]*\}', '{}', p)
        pre = txt[max(0, m.start() - 60):m.start()]
        mm = re.findall(r"\.(get|post|put|patch|delete)\s*(?:<[^>]*>)?\s*\(\s*$", pre)
        method = mm[-1].upper() if mm else ''
        if not method:
            mm = re.search(r"method:\s*['\"](GET|POST|PUT|PATCH|DELETE)['\"]", txt[m.end():m.end() + 160])
            if mm: method = mm.group(1)
        eps[(method, p)] = 1
    return eps
def match(method, p):
    p2 = p.rstrip('/')
    hits = [(rm, rp) for rm, rp, rx in route_res if rx.match(p2) and (not method or rm == method)]
    if not hits:
        # template with {} segments - try loose
        q = re.sub(r'\{\}', 'X', p2)
        hits = [(rm, rp) for rm, rp, rx in route_res if rx.match(q) and (not method or rm == method)]
    return hits
out = []
for app, (adir, sdir) in APPS.items():
    appdir = os.path.join(ROOT, app)
    base = os.path.join(appdir, adir)
    for dp, _, fs in os.walk(base):
        for f in fs:
            if f not in ('page.tsx', 'page.ts', 'route.ts'): continue
            pf = os.path.join(dp, f)
            rel = os.path.relpath(dp, base).replace('\\', '/')
            url = '/' + ('' if rel == '.' else rel)
            url = re.sub(r'/\([^)]+\)', '', url) or '/'
            seen = set(); eps = {}; perms = set(); files = []
            for p, txt in scan(pf, appdir, sdir, 0, seen):
                files.append(os.path.relpath(p, appdir).replace('\\', '/'))
                for k in extract(txt): eps[k] = 1
                for m in PERM.finditer(txt): perms.add(f"{m.group(1)}:{m.group(2)}")
            # layouts: collect permission/status guards from enclosing layouts
            lay = []
            d = dp
            while d.startswith(base):
                lf = os.path.join(d, 'layout.tsx')
                if os.path.isfile(lf): lay.append(os.path.relpath(lf, appdir).replace('\\', '/'))
                if d == base: break
                d = os.path.dirname(d)
            verified = []
            unmatched = []
            for (m, p) in sorted(eps):
                h = match(m, p)
                if h: verified.append(f"{m or '?'} {p}")
                else: unmatched.append(f"{m or '?'} {p}")
            out.append(dict(app=app, url=url, file=os.path.relpath(pf, appdir).replace('\\', '/'), endpoints=verified, unmatched=unmatched, perms=sorted(perms), layouts=lay, nfiles=len(files)))
out.sort(key=lambda x: (x['app'], x['url']))
json.dump(out, open('tests/_pages_raw.json', 'w'), indent=1)
print(len(out))
