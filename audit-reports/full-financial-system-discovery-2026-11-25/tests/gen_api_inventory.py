# Read-only generator: parses backend/cmd/api/*.go (non-test) for mux.Handle/HandleFunc and writes a CSV.
import re, os, csv, sys, glob
ROOT = r"c:\Users\Yousef\Desktop\Projects\POSS"
API = os.path.join(ROOT, "backend", "cmd", "api")
HDIR = os.path.join(ROOT, "backend", "internal", "handler")
OUT = sys.argv[1]

# 1. method index: (Type, Method) -> file:line
meth = {}
bymeth = {}
for dp, _, fs in os.walk(HDIR):
    for f in fs:
        if not f.endswith(".go") or f.endswith("_test.go"): continue
        p = os.path.join(dp, f)
        for i, line in enumerate(open(p, encoding="utf-8", errors="ignore"), 1):
            m = re.match(r"func \(\w+ \*?(\w+)\) (\w+)\(", line)
            if m:
                rel = os.path.relpath(p, ROOT).replace("\\", "/")
                meth[(m.group(1), m.group(2))] = f"{rel}:{i}"
                bymeth.setdefault(m.group(2), []).append((m.group(1), f"{rel}:{i}"))

def find_close(s, start):
    depth = 0; i = start; n = len(s)
    while i < n:
        c = s[i]
        if c == '"':
            i += 1
            while s[i] != '"':
                if ord(s[i]) == 92: i += 1
                i += 1
        elif c == '`':
            i = s.index('`', i + 1)
        elif c == "'":
            i += 1
            while s[i] != "'":
                if ord(s[i]) == 92: i += 1
                i += 1
        elif c == '/' and s[i:i+2] == '//':
            i = s.index('\n', i)
        elif c == '(':
            depth += 1
        elif c == ')':
            depth -= 1
            if depth == 0: return i
        i += 1
    return n

rows = []
files = [f for f in glob.glob(os.path.join(API, "*.go")) if not f.endswith("_test.go")]
for p in files:
    src = open(p, encoding="utf-8", errors="ignore").read()
    # variable -> handler type
    vtype = {}
    for m in re.finditer(r"(\w+)\s*:?=\s*(?:handler|kdsws)\.New(\w+?)\(", src):
        vtype[m.group(1)] = m.group(2)
    rel = os.path.relpath(p, ROOT).replace("\\", "/")
    for m in re.finditer(r"\bmux\.(Handle|HandleFunc)\(", src):
        ls = src.rfind(chr(10), 0, m.start()) + 1
        if src[ls:m.start()].lstrip().startswith("//"): continue
        open_idx = m.end() - 1
        end = find_close(src, open_idx)
        call = src[open_idx+1:end]
        line = src.count("\n", 0, m.start()) + 1
        pm = re.search(r'"((?:GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS) [^"]+)"', call)
        if not pm:
            continue
        method, path = pm.group(1).split(" ", 1)
        rest = call[pm.end():]
        flat = re.sub(r"\s+", " ", rest)
        # handler
        h = None; hfile = ""
        hm = re.search(r"(\w+)\.(\w+)\)", re.sub(r"http\.HandlerFunc\(func.*", "", flat)) if "func(w http.ResponseWriter" not in flat else None
        if "func(w http.ResponseWriter" in flat:
            h = "inline closure"; hfile = f"{rel}:{line}"
        else:
            cands = re.findall(r"\b(\w+)\.(\w+)\b(?=\)|,|$)", flat)
            # prefer variable names that are handlers
            pick = None
            for v, mn in cands:
                if v in vtype or v.endswith("Handler") or v.endswith("Handler"):
                    pick = (v, mn)
            if not pick:
                mm = re.search(r"(?:HandlerFunc|Handle)\(\s*(\w+)\.(\w+)\s*\)", flat) or re.search(r"(\w+)\.(\w+)\s*\)?\s*$", flat)
                if mm: pick = (mm.group(1), mm.group(2))
            if pick:
                v, mn = pick
                h = f"{v}.{mn}"
                t = vtype.get(v)
                loc = meth.get((t, mn)) if t else None
                if not loc:
                    c = bymeth.get(mn, [])
                    loc = c[0][1] if len(c) == 1 else (";".join(x[1] for x in c[:3]) + " (ambiguous)" if c else "")
                hfile = loc
            else:
                h = "unresolved"; hfile = f"{rel}:{line}"
        # wrapper
        wr = []
        mod = lvl = None
        mm = re.search(r"wrap(WithPermission|AuthPermission)\(", flat)
        for name in ["wrapWithPermission", "wrapAuthPermission", "wrapKDSWithPermission", "wrapProtectedNoStatus", "wrapProtected", "wrapAuthOnly", "wrapSuperAdmin"]:
            if name + "(" in flat: wr.append(name)
        if "guards.authenticateTenant" in flat: wr.append("authenticateTenant(wrapProtectedNoStatus)")
        for name in ["middleware.RateLimitByIPWithStore", "middleware.RateLimitByIP", "middleware.RateLimitByUserWindow", "middleware.RateLimitByUser", "middleware.RequireTrustedOrigin", "middleware.TenantStatusMiddleware", "middleware.CommercialCheckoutGate", "middleware.SuperAdminAuthMiddleware", "middleware.RequireInternalToken", "middleware.SessionAuthMiddleware"]:
            if name in flat: wr.append(name.replace("middleware.", ""))
        if "tenantContext(" in flat: wr.append("tenantContext")
        if "authenticateRequest(" in src[m.start():m.start()+120] : wr.append("authenticateRequest")
        pmm = re.search(r"(\d+), \"(READ|WRITE|DELETE|ADMIN|EXPORT|APPROVE|MANAGE)\"\)\)?\s*$", flat)
        if pmm: mod, lvl = pmm.group(1), pmm.group(2)
        pm2 = re.search(r"wrapKDSWithPermission\(.*?, \"(READ|WRITE|ADMIN)\", (true|false)\)", flat)
        if pm2: mod, lvl = "13", pm2.group(1)
        if "guards.requirePermission" in flat:
            q = re.search(r"requirePermission\((\d+), \"(\w+)\"\)", flat)
            if q: mod, lvl = q.group(1), q.group(2)
        # authentication
        if any(x in wr for x in ["wrapWithPermission","wrapAuthPermission","wrapKDSWithPermission","wrapProtectedNoStatus","wrapProtected","wrapAuthOnly","authenticateTenant(wrapProtectedNoStatus)"]) or "authenticateRequest" in flat or "SessionAuthMiddleware" in flat:
            auth = "tenant-session-auth"
        elif "wrapSuperAdmin" in wr or "SuperAdminAuthMiddleware" in flat:
            auth = "super-admin-auth"
        else:
            auth = "none/public"
        status = "statusGuard" if any(x in wr for x in ["wrapWithPermission","wrapKDSWithPermission","wrapProtected"]) or "TenantStatusMiddleware" in flat else ("no-status(blocked tenants allowed)" if "wrapProtectedNoStatus" in wr or "authenticateTenant" in " ".join(wr) and "TenantStatusMiddleware" not in flat else "")
        perm = f"module={mod} level={lvl}" if mod else ""
        mw = "|".join(dict.fromkeys(wr))
        rows.append(dict(method=method, path=path, handler_func=h, handler_file_line=hfile, reg=f"{rel}:{line}", auth=auth, mw=mw, status=status, perm=perm, mod=mod or "", lvl=lvl or ""))
rows.sort(key=lambda r: (r["path"], r["method"]))
import json
json.dump(rows, open(OUT, "w"), indent=0)
print(len(rows))
