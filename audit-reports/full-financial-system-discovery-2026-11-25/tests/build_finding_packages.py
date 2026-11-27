"""Builds findings/F-XXX-*/ packages (10 files each), 57 task register, 58 spec index and 59 regression register
from the authoritative registers (22 findings, 55 evidence check, 34b executed tests). Audit tooling only: reads
the repository (git ls-files) and the audit folder; writes only inside the audit folder. Idempotent."""
import csv, os, re, subprocess, json, collections

AUDIT = r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25'
REPO = r'C:\Users\Yousef\Desktop\Projects\POSS'
os.chdir(AUDIT)
FILES = subprocess.check_output(['git', '-C', REPO, 'ls-files'], text=True).split('\n')
FILESET = set(FILES)

def resolve(token):
    """Return (path or None, how)."""
    t = token.strip().strip('`').replace('\\', '/')
    if t in FILESET:
        return t, 'exact'
    if os.path.exists(os.path.join(REPO, t.replace('/', os.sep))):
        return t, 'untracked-audit-file'
    for pre in ('backend/', 'backend/internal/', 'backend/internal/repository/', 'backend/internal/usecase/', 'backend/internal/handler/http/', 'backend/internal/domain/', 'frontend/apps/dashboard/src/', ''):
        if pre + t in FILESET:
            return pre + t, 'prefix'
    base = t.split('/')[-1]
    m = [f for f in FILES if f.endswith('/' + t) or f.split('/')[-1] == base]
    m = [x for x in m if '/node_modules/' not in x and not x.startswith('audit-reports/')]
    if len(m) == 1:
        return m[0], 'basename'
    if len(m) > 1:
        pref = [x for x in m if x.startswith('backend/internal') or x.startswith('backend/migrations') or x.startswith('frontend/apps')]
        return (pref or m)[0], 'basename-ambiguous(%d)' % len(m)
    return None, 'unresolved'

REF_RE = re.compile(r'([A-Za-z0-9_./\[\]-]+\.(?:go|sql|tsx|ts|mjs|js|yml|yaml)(?![A-Za-z0-9_]))(?::(\d+(?:-\d+)?))?')

def parse_refs(text):
    out, seen = [], set()
    for m in REF_RE.finditer(text or ''):
        tok, ln = m.group(1), m.group(2)
        if tok in seen and not ln:
            continue
        seen.add(tok)
        p, how = resolve(tok)
        out.append((tok, ln, p, how))
    return out

FRONTEND_MAP = [
    (r'registration|onboarding|landing', ['frontend/apps/landing/src/app/[locale]/auth/register/page.tsx']),
    (r'\bpos\b|cashier', ['frontend/apps/dashboard/src/app/[locale]/dashboard/pos/page.tsx']),
    (r'kds|kitchen', ['frontend/apps/dashboard/src/app/[locale]/dashboard/kds']),
    (r'customer invoice|receipt', ['frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customer-invoices']),
    (r'supplier', ['frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/supplier-invoices', 'frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/suppliers']),
    (r'journal|ledger|accounting|bank|chart of accounts|period', ['frontend/apps/dashboard/src/app/[locale]/dashboard/accounting']),
    (r'report|kpi|cash flow|vat|inventory valuation|turnover|break-even|menu profitability', ['frontend/apps/dashboard/src/app/[locale]/dashboard/reports']),
    (r'\bhr\b|payroll|salary|salaries|employee|attendance|leave', ['frontend/apps/dashboard/src/app/[locale]/dashboard/hr']),
    (r'storefront|public menu', ['frontend/apps/restaurant-website/app/[locale]']),
    (r'super admin|saas|subscription|trial|billing', ['frontend/apps/super-admin']),
    (r'settings', ['frontend/apps/dashboard/src/app/[locale]/dashboard/settings']),
    (r'product|categor|inventory', ['frontend/apps/dashboard/src/app/[locale]/dashboard/products']),
    (r'order', ['frontend/apps/dashboard/src/app/[locale]/dashboard/orders']),
]

def frontend_paths(f):
    blob = (f['affected_modules'] + ' ' + f['title']).lower()
    out = []
    for pat, paths in FRONTEND_MAP:
        if re.search(pat, blob):
            for p in paths:
                full = os.path.join(REPO, p.replace('/', os.sep))
                out.append((p, os.path.exists(full)))
    seen, res = set(), []
    for p, ok in out:
        if p not in seen:
            seen.add(p); res.append((p, ok))
    return res

ACCOUNT_RE = re.compile(r'\b(1110|1111|1112|1113|1120|1130|1150|2101|2110|2120|3100|3150|4100|5110|5210|5220|5310)\b')
TABLES = ['orders', 'order_items', 'journal_entries', 'journal_entry_lines', 'general_ledger', 'chart_of_accounts', 'invoices_customer', 'invoices_supplier', 'invoice_lines_customer', 'invoice_lines_supplier', 'payments_received', 'payments_made', 'products', 'inventory_movements', 'kds_queue', 'salaries', 'attendance', 'leaves', 'employees', 'tenants', 'users', 'roles', 'role_permissions', 'bank_accounts_ledger', 'accounting_event_outbox', 'fiscal_periods', 'tax_periods', 'auto_posting_rules', 'recurring_journal_templates', 'exchange_rates', 'budget_entries', 'restaurants', 'customers_accounting', 'suppliers']

def sev_rank(s):
    s = s.lower()
    return 0 if s.startswith('critical') else 1 if s.startswith('high') else 2 if s.startswith('medium') else 3

def sev_short(s):
    return s.split(' ')[0]

def conf(rc):
    rc = rc or ''
    m = re.match(r'\s*(CONFIRMED|PROBABLE|CANDIDATE|UNKNOWN)', rc or '')
    if m:
        return m.group(1)
    if re.search(r'CONFIRMED', rc or ''):
        return 'CONFIRMED (in appended RUN-003 text)'
    return 'UNSTATED in register (treated as CANDIDATE until verified)'

def slug(title):
    s = re.sub(r'[^a-z0-9]+', '-', title.lower()).strip('-')
    return s[:48].rstrip('-')

def kinds(f):
    blob = ' '.join([f['category'], f['title'], f['affected_modules']]).lower()
    k = set()
    if 'security' in blob or 'tenant' in blob or 'rbac' in blob or 'privilege' in blob or 'cors' in blob:
        k.add('security')
    if re.search(r'journal|ledger|posting|cogs|vat|tax|invoice|receipt|payment|cash|revenue|payroll|salary|salaries|accounting|inventory valuation|refund|bank', blob):
        k.add('accounting')
    if re.search(r'report|kpi|metric|stats|turnover|forecast', blob):
        k.add('report')
    if 'race' in blob or 'concurren' in blob or 'idempoten' in blob or 'duplicate' in blob:
        k.add('concurrency')
    if 'missing capability' in blob or 'no refund' in blob or 'no api' in blob:
        k.add('capability')
    if 'validation' in blob or 'input' in blob:
        k.add('validation')
    if re.search(r'ui|ux|label|currency symbol|toast|page|form', blob):
        k.add('ui')
    if 'reliab' in blob or 'redis' in blob or 'environment' in blob:
        k.add('reliability')
    return k

def load():
    rows = list(csv.DictReader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
    ev = {r['finding_id']: r for r in csv.DictReader(open('55-finding-evidence-and-root-cause-verification.csv', encoding='utf-8'))} if os.path.exists('55-finding-evidence-and-root-cause-verification.csv') else {}
    ex = list(csv.DictReader(open('34b-executed-test-results.csv', encoding='utf-8')))
    return rows, ev, ex


SHARED_OWNER = {  # member -> (owner finding, shared component) : implement the component once under the owner
    'F-031': ('F-025', 'error-mapping layer'), 'F-050': ('F-025', 'error-mapping layer'), 'F-064': ('F-025', 'error-mapping layer'),
    'F-033': ('F-020', 'role/permission matrix'), 'F-027': ('F-020', 'role/permission matrix'), 'F-071': ('F-020', 'role/permission matrix'),
    'F-017': ('F-052', 'idempotency store'), 'F-047': ('F-052', 'idempotency store'),
    'F-042': ('F-012', 'shared revenue definitions'), 'F-043': ('F-012', 'shared revenue definitions'),
    'F-041': ('F-007', 'discount-first tax calculation'), 'F-024': ('F-007', 'discount-first tax calculation'), 'F-004': ('F-007', 'discount-first tax calculation'),
    'F-058': ('F-062', 'server-side input validation helpers'), 'F-037': ('F-062', 'server-side input validation helpers'), 'F-056': ('F-062', 'server-side input validation helpers'), 'F-053': ('F-062', 'server-side input validation helpers'), 'F-048': ('F-062', 'server-side input validation helpers'), 'F-059': ('F-062', 'server-side input validation helpers'),
    'F-045': ('F-003', 'cash-flow method'),
}
NON_MONEY_CATEGORY = r'UI|validation|error|security|configuration|time zone|export|report logic|aggregation|API contract'

OVERRIDES = {}
if os.path.exists(os.path.join(AUDIT, 'tests', 'finding_overrides.json')):
    OVERRIDES = json.load(open(os.path.join(AUDIT, 'tests', 'finding_overrides.json'), encoding='utf-8'))

import glob as _glob
for _r4 in sorted(_glob.glob(os.path.join(AUDIT, 'tests', 'finding_overrides_r4*.py'))):  # RUN-004 hand-written designs override the generated defaults (and the older json for the same id)
    _ns = {}
    exec(compile(open(_r4, encoding='utf-8').read(), _r4, 'exec'), _ns)
    for _k, _v in _ns.get('R4', {}).items():
        OVERRIDES[_k] = {**OVERRIDES.get(_k, {}), **_v}

def tests_for(fid, ex, f):
    ids = []
    for r in ex:
        blob = ' '.join(r.values())
        if re.search(r'\b%s\b' % re.escape(fid), blob):
            ids.append(r['exec_id'])
    return ids

def verify_refs(f):
    refs = [r for r in parse_refs(f['source_refs'] + ' ' + f['root_cause_first_break']) if '<next>' not in r[0]]
    good, bad = [], []
    for tok, ln, p, how in refs:
        if p:
            full = os.path.join(REPO, p.replace('/', os.sep))
            exists = os.path.exists(full)
            n_lines = None
            if exists:
                try:
                    n_lines = sum(1 for _ in open(full, encoding='utf-8', errors='replace'))
                except Exception:
                    n_lines = None
            rng_ok = True
            if ln and n_lines:
                hi = int(ln.split('-')[-1])
                rng_ok = hi <= n_lines
            (good if exists else bad).append((tok, ln, p, how, n_lines, rng_ok))
        else:
            bad.append((tok, ln, None, how, None, False))
    _seen = set(); _g = []
    for _x in good:
        if (_x[2], _x[1]) not in _seen:
            _seen.add((_x[2], _x[1])); _g.append(_x)
    return _g, bad

def db_tables(f):
    blob = ' '.join(f.values())
    return [t for t in TABLES if re.search(r'\b%s\b' % t, blob)]

def md_table(rows, head):
    out = ['| ' + ' | '.join(head) + ' |', '|' + '|'.join(['---'] * len(head)) + '|']
    for r in rows:
        out.append('| ' + ' | '.join(str(c).replace('\n', ' ').replace('|', '/') for c in r) + ' |')
    return '\n'.join(out)

def related_findings(f, all_rows):
    ids = set(re.findall(r'F-\d{3}', f['related_static_ids'] + ' ' + f['root_cause_first_break'] + ' ' + f['actual'] + ' ' + f['evidence_paths']))
    ids.discard(f['finding_id'])
    # shared source file
    mine = {p for _, _, p, _ in parse_refs(f['source_refs']) if p}
    shared = collections.defaultdict(list)
    for o in all_rows:
        if o['finding_id'] == f['finding_id']:
            continue
        theirs = {p for _, _, p, _ in parse_refs(o['source_refs']) if p}
        inter = mine & theirs
        if inter:
            shared[o['finding_id']] = sorted(inter)
    return sorted(ids), shared

def build():
    rows, ev, ex = load()
    os.makedirs('findings', exist_ok=True)
    task_rows, reg_rows, index_rows = [], [], []
    summary = []
    for f in rows:
        fid = f['finding_id']
        ov = OVERRIDES.get(fid, {})
        d = os.path.join('findings', '%s-%s' % (fid, slug(f['title'])))
        os.makedirs(d, exist_ok=True)
        K = kinds(f)
        good, bad = verify_refs(f)
        fe = frontend_paths(f)
        tabs = db_tables(f)
        accts = sorted(set(ACCOUNT_RE.findall(' '.join(f.values()))))
        rel_ids, shared = related_findings(f, rows)
        tids = tests_for(fid, ex, f)
        c = conf(f['root_cause_first_break'])
        evs = [e.strip() for e in re.split(r';\s*', f['evidence_paths']) if e.strip()]
        sev = sev_short(f['severity'])
        scripts = [e for e in evs if e.startswith('tests/runtime/')]
        # ---------- finding.md
        open(os.path.join(d, 'finding.md'), 'w', encoding='utf-8').write(f"""# {fid} — {f['title']}

- **Severity:** {f['severity']}  ·  **Category:** {f['category']}  ·  **Evidence status:** {f['evidence_status']}
- **Affected modules:** {f['affected_modules']}
- **Root-cause confidence:** {c}
- **Executed audit tests referencing this finding:** {', '.join(tids) if tids else 'none recorded in 34b (evidence is in the register/evidence files)'}
- **Remediation status:** NOT STARTED (audit-only deliverable; no product code was changed)

## Problem statement
{f['expected'] and 'Expected: ' + f['expected']}

Observed: {f['actual']}

## Severity rationale
Rated **{f['severity']}** (category: {f['category']}; evidence: {f['evidence_status']}). Basis recorded in the register: {f['monetary_impact_EGP'] or 'no monetary effect recorded'}. Root-cause confidence: {c}. Severity is the auditor's rating from the recorded impact and reproducibility; it is not an approved business priority.

## Business impact
Affected modules: {f['affected_modules']}. {f['monetary_impact_EGP'] or 'The register records no direct monetary effect; the impact is on correctness, usability or control quality of the listed modules.'}

## Technical impact
{f['root_cause_first_break'][:420]}

Data integrity: {'stored monetary or ledger data can become inconsistent (see financial-impact.md)' if (('accounting' in K or 'concurrency' in K) and not re.search(NON_MONEY_CATEGORY, f['category'], re.I)) else 'no stored monetary or ledger corruption is recorded for this finding; the effect is on presentation, validation, availability or access'}. Security relevance: {'yes, see category' if 'security' in K else 'none recorded'}. Affected tables, APIs, components, workers and reports are listed in [affected-code-paths.md](affected-code-paths.md) and [financial-impact.md](financial-impact.md).

## Monetary / operational impact
{f['monetary_impact_EGP'] or 'Not quantified in the register.'}

## Where it was reproduced
{f['reproduction']}

See [reproduction.md](reproduction.md), [root-cause.md](root-cause.md), [affected-code-paths.md](affected-code-paths.md), [financial-impact.md](financial-impact.md), [dependency-map.md](dependency-map.md), [evidence-index.md](evidence-index.md), [remediation-spec.md](remediation-spec.md), [remediation-tasks.md](remediation-tasks.md), [regression-test-spec.md](regression-test-spec.md).
""")
        # ---------- reproduction.md
        env = 'Isolated Docker stack `pos-p2a-ws05-e2e` (DB `pos_phase2a_e2e`, API :18082, dashboard :13112). Synthetic accounts only.'
        open(os.path.join(d, 'reproduction.md'), 'w', encoding='utf-8').write(f"""# {fid} — Reproduction

## Environment
{env}

## Steps (as recorded in the register)
{f['reproduction']}

## Expected result
{f['expected']}

## Actual result (measured)
{f['actual']}

## Re-runnable scripts / evidence
{chr(10).join('- `%s`' % s for s in scripts) if scripts else '- No dedicated script; reproduction was performed through the UI/API/SQL steps above and is captured in the evidence files listed in evidence-index.md.'}

## Preconditions and cautions
- Use a diagnostic tenant (tenant 4 fixture or a new tenant); tenants 5, 6, 7 are protected evidence for F-001/F-038 and must not be seeded or modified.
- Mutating calls through the API require the header `Origin: http://localhost:13112` (CSRF origin check) and a session cookie from `POST /api/v1/auth/login`.
- Record the dataset state before/after (row counts, GL balances) so that the result can be compared after a future fix.
""")
        # ---------- root-cause.md
        rc = f['root_cause_first_break']
        extra_rc = ov.get('root_cause_notes', '')
        missing = '' if c.startswith('CONFIRMED') else '\n## Missing evidence and required next investigation\n- The register does not mark the root cause CONFIRMED to a file:line. Next step: read the listed source files at the cited lines, add a failing automated test that reproduces the defect at the code level, and promote the confidence level.\n'
        open(os.path.join(d, 'root-cause.md'), 'w', encoding='utf-8').write(f"""# {fid} — Root cause

**Confidence:** {c}

## First break (from the register)
{rc}

{extra_rc}
## Origin and propagation
- Origin (first verified source location): {('`%s`%s' % (good[0][2], (':' + good[0][1]) if good[0][1] else '')) if good else 'see the first break text above'}.
- Propagation: affected modules {f['affected_modules']}; related findings {', '.join(rel_ids) or 'none recorded'}; observed consequence: {f['actual'][:240]}
- Cause versus symptom: the symptom is what users or reports show ({f['actual'][:140]}); the cause is the first break above.

## Investigation performed
- Executed audit tests: {', '.join(tids) if tids else 'none recorded in 34b'}; scripts: {', '.join(scripts[:3]) if scripts else 'reproduced through UI/API/SQL steps in reproduction.md'}.
- Evidence verification: see evidence-index.md (file 55 resolves every evidence path).

## Verified source locations
{md_table([(tok, ln or '-', p, how, 'yes' if rng else 'RANGE BEYOND FILE') for tok, ln, p, how, n, rng in good], ['reference', 'lines', 'resolved path', 'resolution', 'range valid']) if good else 'No source reference could be resolved to an existing file.'}
{('Unresolved references: ' + ', '.join(t for t, _, _, _, _, _ in bad)) if bad else ''}
{missing}""")
        # ---------- affected-code-paths.md
        fe_rows = [(p, 'exists' if ok else 'NOT FOUND') for p, ok in fe]
        open(os.path.join(d, 'affected-code-paths.md'), 'w', encoding='utf-8').write(f"""# {fid} — Affected code paths

Existing paths were checked against the repository at the audited SHA (`git ls-files`). Paths in the *Proposed* section do not exist yet.

## Backend / database (existing, verified)
{md_table([(p, ln or '-', how) for tok, ln, p, how, n, rng in good], ['path', 'lines', 'resolution']) if good else 'None resolved from the register source references.'}

## Frontend (existing, module-derived — verify the exact component before editing)
{md_table(fe_rows, ['path', 'status']) if fe_rows else 'No frontend surface is involved or none could be derived.'}

## Database objects involved
{', '.join('`%s`' % t for t in tabs) if tabs else 'None identified from the finding text.'}

## API contracts involved
{chr(10).join('- ' + a for a in ov.get('api', [])) or 'See 08-api-inventory.csv for the endpoints in the affected modules; no contract change is proposed beyond the remediation spec.'}

## Proposed (future) files
{chr(10).join('- ' + a for a in ov.get('proposed_files', [])) or '- A regression test file next to the affected package (see regression-test-spec.md). Any new migration must use the next free number after the current highest in backend/migrations.'}
""")
        # ---------- financial-impact.md
        acct_txt = ', '.join(accts) if accts else 'none named'
        open(os.path.join(d, 'financial-impact.md'), 'w', encoding='utf-8').write(f"""# {fid} — Financial impact

- **Quantified impact (register):** {f['monetary_impact_EGP'] or 'not quantified'}
- **Accounts / statements named in the evidence:** {acct_txt}
- **Impact classes:** {', '.join(sorted(K)) or 'general'}

## What is wrong in the books or reports
{f['actual']}

## Which reports and decisions are affected
{('Financial reports that read the affected data: ' + f['affected_modules']) if 'report' in K or 'accounting' in K else 'No financial statement is directly changed; impact is operational/security (see remediation-spec.md).'}

## Reconciliation note
Do not offset this discrepancy against unrelated residuals. Independent reconciliation for the dataset used in the audit is in 44-transaction-level-reconciliation.csv and 43-report-metric-runtime-results.csv.
{ov.get('financial_notes', '')}
""")
        # ---------- dependency-map.md
        blockers = {'F-001': 'Unlocks every financial workflow on a newly registered tenant (nothing posts without a chart of accounts).', 'F-014': 'Unlocks UI employee creation (HR/payroll from scratch).', 'F-022': 'Unlocks the public storefront and online ordering tests.', 'F-013': 'Unlocks period-close testing.', 'F-005': 'Unlocks card/wallet payment accounting tests.', 'F-026': 'Unlocks card/wallet payment accounting tests.'}
        open(os.path.join(d, 'dependency-map.md'), 'w', encoding='utf-8').write(f"""# {fid} — Dependency map

## Findings explicitly related (register text / related ids)
{chr(10).join('- ' + r for r in rel_ids) or '- none'}

## Findings sharing a source file (candidate shared fix / merge conflicts)
{chr(10).join('- %s via %s' % (k, ', '.join(v)) for k, v in sorted(shared.items())) or '- none'}

## Unblocks (if this finding is fixed)
{blockers.get(fid, '- No blocked audit workflow depends on this finding.')}

## Depends on
{ov.get('depends', '- No prerequisite fix identified. Policy decisions listed in remediation-spec.md (section Required policy decisions) must be made first where applicable.')}
""")
        # ---------- evidence-index.md
        e55 = ev.get(fid, {})
        rows_ev = []
        for e in evs:
            full = os.path.join(AUDIT, e.replace('/', os.sep))
            rows_ev.append((e, 'exists' if (os.path.exists(full) or e.startswith('tests/') and os.path.exists(os.path.join(AUDIT, e))) else ('row-reference / pattern' if '*' in e or e.endswith('/') else 'NOT FOUND')))
        open(os.path.join(d, 'evidence-index.md'), 'w', encoding='utf-8').write(f"""# {fid} — Evidence index

- Register evidence status: **{f['evidence_status']}**
- Evidence verification (file 55): {e55.get('evidence_check', e55.get('status', 'see 55-finding-evidence-and-root-cause-verification.csv'))}
- Executed tests: {', '.join(tids) if tids else 'none in 34b'}

## Evidence files (relative to the audit folder)
{md_table(rows_ev, ['path', 'status']) if rows_ev else 'No evidence paths listed.'}

State files `evidence/state-*.json` are session cookies and are intentionally not referenced.
""")
        # ---------- remediation-spec.md
        design = ov.get('design') or f['remediation_direction']
        alternatives = ov.get('alternatives') or ('Alternative A: enforce in the database (constraint/trigger/lock) so every writer is covered; Alternative B: enforce in the application layer only (smaller change, but direct writers and future code paths stay unprotected). Prefer A for monetary invariants, B for presentation/validation defects.' if ('accounting' in K or 'concurrency' in K) else 'Alternative A: fix at the source layer named in root-cause.md; Alternative B: add a compensating check in the calling layer (weaker, not recommended unless the source layer is vendor/owner-controlled).')
        policy = ov.get('policy') or ('Decide and record in 56-accounting-policy-decision-package.md before implementation: ' + ('revenue-recognition basis and tax base (discount-first vs gross)' if re.search(r'tax|vat|revenue|discount', f['title'].lower()) else 'the accounting treatment of the affected transaction') if 'accounting' in K or 'report' in K else 'No accounting policy decision is required; product decisions are limited to the behaviour described in Required behavior.')
        migr = ov.get('migration') or ('Yes — a numbered migration after the current highest in backend/migrations (data-correction + constraint/trigger change). Existing rows must be reconciled first.' if re.search(r'trigger|migration|constraint|\.sql', f['source_refs'] + f['root_cause_first_break']) or 'accounting' in K and 'concurrency' in K else 'No schema change expected; confirm during implementation.')
        hist = ov.get('historical') or ('Existing records created under the defect keep wrong values. Produce a read-only exception report (query by the evidence criteria) and have the accountant approve any correcting journal; never mutate posted journals in place.' if 'accounting' in K or 'report' in K else 'No financial history is affected; existing data needs no correction unless noted in the evidence.')
        idem = ov.get('idempotency') or ('Every mutating endpoint touched by this fix must be safe to retry: accept an Idempotency-Key (see F-052) or be naturally idempotent, and must be covered by a concurrent-call test.' if 'concurrency' in K or 'accounting' in K else 'Not applicable beyond keeping existing behaviour.')
        authz = ov.get('authz') or ('Keep RBAC module/level checks unchanged; add negative tests for roles without permission.' if 'security' in K else 'No authorization change; keep existing wrapWithPermission levels.')
        open(os.path.join(d, 'remediation-spec.md'), 'w', encoding='utf-8').write(f"""# {fid} — Remediation specification (implementation-ready, NOT IMPLEMENTED)

> Confirmed facts come from the audit evidence. Proposed design decisions are marked **Proposed**. Source paths listed under "Affected paths" exist; anything under "Proposed files" does not.

## 1. Exact problem statement
{f['title']}.

## 2. Current behavior (verified)
{f['actual']}

## 3. Required behavior
{f['expected']}

## 4. Reproduction conditions
{f['reproduction']}

## 5. Verified evidence
{chr(10).join('- ' + e for e in evs) or '- see evidence-index.md'}

## 6. Root cause and confidence
{c} — {re.sub(r'^\s*(CONFIRMED|PROBABLE|CANDIDATE|UNKNOWN)[^:]*:?\s*', '', rc)[:700]}

## 7. Affected paths
- Backend/database: {', '.join('`%s`%s' % (p, (':' + ln) if ln else '') for _, ln, p, _, _, _ in good) or 'none resolved'}
- Frontend: {', '.join('`%s`' % p for p, ok in fe if ok) or 'none'}
- API contracts: {'; '.join(ov.get('api', [])) or 'unchanged unless the design below says otherwise (see 08-api-inventory.csv)'}
- Accounting logic / posting rules: {'accounts ' + acct_txt if accts else 'none named'}
- Reports affected: {f['affected_modules'] if 'report' in K or 'accounting' in K else 'none'}

## 8. Impact
- Financial: {f['monetary_impact_EGP'] or 'not quantified'}
- Operational: {f['affected_modules']} workflows behave as in section 2.
- Security: {'Yes — see category: ' + f['category'] if 'security' in K else 'No direct security impact identified.'}
- Data integrity: {'Yes — monetary/ledger data can become inconsistent (see financial-impact.md).' if (('accounting' in K or 'concurrency' in K) and not re.search(r'UI|validation|error|security|configuration|time zone|export|report logic|aggregation|API contract', f['category'], re.I)) else 'No stored monetary/ledger data is corrupted by this defect; impact is on presentation, validation, availability or access (see financial-impact.md).'}

## 9. Proposed remediation design (Proposed)
{design}

## 10. Alternative designs
{alternatives}

## 11. Required policy decisions
{policy}

## 12. Required data migrations
{migr}

## 13. Historical data implications
{hist}

## 14. Transaction-integrity requirements
{'Perform the check and the write in one database transaction with row-level locking (SELECT … FOR UPDATE) or a constraint; no read-then-write across transactions.' if 'concurrency' in K or 'accounting' in K else 'Single-request change; keep existing transaction boundaries.'}

## 15. Idempotency requirements
{idem}

## 16. Authorization requirements
{authz}

## 17. Compatibility requirements
Existing API clients (dashboard, storefront, super-admin, mobile if any) must keep working: do not rename response fields; add fields additively; keep error shapes `{{"error":..., "success":false}}`.

## 18. Recovery requirements
{'If posting fails, the existing durable outbox (accounting_event_outbox) must retry; the fix must not create a path that marks an operation successful without its accounting event.' if 'accounting' in K else 'Failed requests must leave no partial state.'}

## 19. Rollback strategy
Ship behind a reversible change: for migrations provide a down-migration or a forward-fix plan; for application changes keep the previous binary deployable. Snapshot the database (checksummed pg_dump) before any data correction.

## 20. Dependencies
{', '.join(rel_ids) or 'none'}{' and shared-file findings ' + ', '.join(sorted(shared)) if shared else ''}

## 21. Risks
- Regression in modules sharing the same source files ({', '.join(sorted({p for _, _, p, _, _, _ in good})[:4]) or 'n/a'}).
- Historical data that no longer satisfies a new constraint (migration must handle it first).

## 22. Required tests
{f['regression_tests_needed']} (see regression-test-spec.md and 59-post-remediation-regression-test-register.csv)

## 23. Measurable acceptance criteria
{chr(10).join('- ' + a for a in ov.get('acceptance', [])) or '- Re-running the reproduction steps yields the Expected result: ' + f['expected'] + chr(10) + '- The independent oracle/SQL comparison for the evidence dataset shows zero residual for this finding.' + chr(10) + '- All tests in regression-test-spec.md pass on the isolated stack with real headed Chrome for UI-facing checks.'}
""")
        # ---------- remediation-tasks.md
        phases = []
        n = [0]
        def T(phase, obj, paths, changes, deps='none', tests='see regression-test-spec.md', deliver='code + tests', acc=None, risk='low'):
            n[0] += 1
            tid = 'REM-%s-%02d' % (fid, n[0])
            # RUN-005 tailored acceptance: replace generic phase criteria with finding-specific, measurable ones
            _exp = (f['expected'] or '')[:150]
            if phase.startswith('P4'):
                acc = (ov.get('acceptance') or [None])[0] or ('UI, API and stored values agree for the reproduction dataset: ' + _exp)
                if len(acc) < 60:
                    acc = acc + ' (expected behaviour: ' + _exp + ')'
            elif phase.startswith('P5'):
                acc = 'Journals, General Ledger and the affected reports reproduce the expected result for the reproduction dataset (' + _exp + ') and the residual for this finding is 0.00 EGP; unrelated residuals are unchanged'
            elif phase.startswith('P6'):
                acc = 'Negative authorization tests pass: roles without the permission and other tenants receive 401/403/404 on the affected routes for ' + fid + '; no cross-tenant data and no raw database text in responses'
            elif phase.startswith('P7'):
                acc = 'Read-only exception report lists the affected historical records (' + (hist or '')[:110].replace('\n', ' ') + '); totals reconcile before and after; accountant sign-off recorded; no in-place mutation of posted journals'
            elif phase.startswith('P8'):
                acc = 'Every REG-' + fid + '-* test in regression-test-spec.md passes on the isolated stack, including the expected result: ' + _exp + '; evidence stored under the release evidence folder'
            elif phase.startswith('P9'):
                acc = 'On staging the reproduction steps of ' + fid + ' no longer reproduce (expected: ' + _exp + '); the rollback in remediation-spec section 19 is rehearsed and its duration recorded'
            phases.append((phase, tid, obj, paths, changes, deps, tests, deliver, acc or ('Re-run reproduction: ' + f['expected'][:160]), risk))
            task_rows.append({'task_id': tid, 'finding_id': fid, 'severity': sev, 'phase': phase, 'objective': obj, 'existing_paths': '; '.join(paths) if isinstance(paths, list) else paths, 'required_changes': changes, 'dependencies': deps, 'required_tests': tests, 'acceptance': acc or f['expected'], 'risk': risk, 'status': 'NOT STARTED'})
            return tid
        bpaths = [p for _, _, p, _, _, _ in good][:4] or ['(locate per root-cause.md)']
        t0 = None
        if policy.startswith('Decide') or 'POL-' in policy:
            t0 = T('P0 Requirements and policy', 'Record the policy decision for ' + fid, ['audit-reports/full-financial-system-discovery-2026-11-25/56-accounting-policy-decision-package.md'], 'Owner/accountant selects one option; recorded with date and approver. No code until approved.', deps='none', tests='n/a', deliver='signed decision entry', acc='Decision recorded; remediation spec section 11 updated.')
        t1 = T('P1 Failing test first', 'Add automated test that reproduces ' + fid + ' at code level', bpaths, 'Write a failing unit/integration test using the evidence dataset (' + (scripts[0] if scripts else 'reproduction steps') + ').', deps=t0 or 'none', tests='the new test must fail before the fix', deliver='failing test committed on a branch', acc='Test fails with the observed Actual value.')
        if migr.startswith('Yes'):
            T('P2 Database and data integrity', 'Add migration/constraint for ' + fid, ['backend/migrations (next free number)'], 'Constraint/trigger/lock change per remediation-spec sections 12-14; include data pre-check query.', deps=t1, tests='migration applies on a copy of the isolated dump; RLS tables unchanged', deliver='migration + rollback note', acc='Migration applies cleanly on the isolated dump and on an empty database.', risk='high')
        money = ('accounting' in K) and not re.search(NON_MONEY_CATEGORY, f['category'], re.I)
        hist_needed = not re.match(r'\s*(None|No stored|No)', hist)
        own = SHARED_OWNER.get(fid)
        p3_acc = ' | '.join(ov.get('acceptance', [])) or ('Failing test from ' + t1 + ' passes.')
        if own:
            T('P3 Backend implementation', 'Apply the shared ' + own[1] + ' to cover ' + fid + ' (core component delivered once by the P3 task of ' + own[0] + ')', bpaths, 'SHARED IMPLEMENTATION - do not re-implement the component; add only the ' + fid + '-specific cases. Design: ' + design[:600], deps=own[0] + ' P3 task; ' + t1, tests='unit + integration', deliver='backend change (finding-specific cases)', acc=p3_acc, risk='medium')
        else:
            T('P3 Backend implementation', 'Fix root cause of ' + fid, bpaths, design[:900], deps=t1, tests='unit + integration', deliver='backend change', acc=p3_acc, risk='medium')
        api_txt = '; '.join(ov.get('api', []))
        contract_change = bool(api_txt) and not api_txt.lower().startswith('unchanged')
        if (fe and any(ok for _, ok in fe) and ('ui' in K) and contract_change) or (contract_change and re.search(r'UI', f['category'])):
            T('P4 Frontend and API contract', 'Align UI and API contract for ' + fid, [p for p, ok in fe if ok][:3] or bpaths, 'API contract: ' + api_txt[:300] + '. Update the affected screens to the corrected contract (EN and AR); keep response fields additive.', deps='P3 task', tests='headed Chrome check of the affected screens (EN and AR)', deliver='frontend change + contract test', acc='UI value equals API value equals stored value for the reproduction dataset.')
        if money:
            T('P5 Posting and reporting integration', 'Verify posting rules and reports for ' + fid, ['backend/internal/usecase/auto_posting_engine.go', 'backend/internal/usecase/financial_report_usecase.go'], 'Confirm journal lines, GL, trial balance and affected reports reflect the corrected behavior; update expectations in 43/44 registers.', deps='P3 task', tests='oracle reconciliation (Decimal / integer minor units)', deliver='reconciliation evidence', acc='Residual for this finding equals 0.00 EGP.')
        if 'security' in K:
            T('P6 Security and authorization', 'Add negative authorization tests for ' + fid, bpaths, 'Cross-role and cross-tenant tests; confirm no new information exposure.', deps='P3 task', tests='API security battery (ph17_security.mjs pattern)', deliver='security tests', acc='Denied roles receive 401/403; no cross-tenant data.')
        if hist_needed and (money or 'report' in K):
            T('P7 Historical data remediation', 'Produce exception report and approved correction plan for ' + fid, ['audit-reports/full-financial-system-discovery-2026-11-25/evidence (queries)'], hist[:300], deps='P3 task', tests='report reconciles before/after', deliver='read-only exception report', acc='Accountant sign-off; no in-place mutation of posted journals.')
        T('P8 Verification', 'Run regression suite for ' + fid + ' (API, real Chrome, reconciliation' + (', concurrency' if 'concurrency' in K else '') + ')', ['tests/runtime (audit scripts) + new automated tests'], f['regression_tests_needed'], deps='all previous', tests='regression-test-spec.md', deliver='test report', acc='All REG tests pass; evidence stored under evidence/ for the release.')
        T('P9 Release and rollback', 'Release verification for ' + fid, ['deployment pipeline'], 'Deploy to staging, repeat reproduction, verify rollback plan (spec section 19).', deps='P8', tests='staging smoke + reproduction', deliver='release note', acc='Reproduction no longer fails on staging; rollback rehearsed.')
        lines = [f"# {fid} — Implementation phases and tasks (all NOT STARTED)\n", f"Finding: {f['title']} · Severity {f['severity']} · Root-cause confidence {c}\n", 'Tasks are derived from the verified evidence and source locations. Existing paths exist at the audited SHA; refine task granularity during implementation planning. No task has been started.\n']
        cur = None
        for ph, tid, obj, paths, changes, deps, tests, deliver, acc, risk in phases:
            if ph != cur:
                lines.append('\n## %s\n' % ph); cur = ph
            lines.append(f"### {tid}\n- **Finding:** {fid}\n- **Phase:** {ph}\n- **Objective:** {obj}\n- **Existing source paths:** {', '.join(paths) if isinstance(paths, list) else paths}\n- **Proposed paths:** {'; '.join(ov.get('proposed_files', [])) or 'tests next to the affected package; new migration if required'}\n- **Required changes:** {changes}\n- **Dependencies:** {deps}\n- **Prerequisites:** isolated stack reproduction passes (see reproduction.md)\n- **Expected deliverables:** {deliver}\n- **Required tests:** {tests}\n- **Acceptance criteria:** {acc}\n- **Risks:** {risk}\n- **Status:** NOT STARTED\n")
        open(os.path.join(d, 'remediation-tasks.md'), 'w', encoding='utf-8').write('\n'.join(lines))
        # ---------- regression-test-spec.md
        ui_test = 'ui' in K or 'report' in K or 'capability' in K
        tests = []
        k = 0
        def R(ttype, desc, pre, assertion, auto, orig):
            nonlocal k
            k += 1
            tid = 'REG-%s-%02d' % (fid, k)
            tests.append((tid, ttype, desc, pre, assertion, auto, orig))
            reg_rows.append({'test_id': tid, 'finding_id': fid, 'severity': sev, 'type': ttype, 'description': desc, 'preconditions': pre, 'assertion': assertion, 'automation': auto, 'original_audit_script': orig, 'status': 'NOT STARTED'})
        orig = scripts[0] if scripts else 'n/a'
        R('regression (API/DB)', f['regression_tests_needed'], 'isolated stack + dataset from reproduction.md', 'Expected: ' + f['expected'][:220], 'automated (Go integration / API)', orig)
        if ui_test:
            R('real Chrome UI', 'Repeat the user-visible steps in headed Google Chrome (EN and AR where text/labels are involved)', 'diagnostic tenant, logged in via real login UI', 'UI value equals API value equals stored value; screenshot captured', 'Playwright headed Chrome', orig)
        if money or 'report' in K:
            R('financial reconciliation', 'Independent oracle comparison (Decimal / integer minor units) for the dataset', 'GL, TB and affected report captured before and after', 'Residual 0.00 EGP for this finding; unrelated residuals unchanged', 'SQL + oracle script', orig)
        if 'concurrency' in K:
            R('concurrency', 'Fire N parallel requests (N>=5) for the affected operation', 'fresh invoice/order with known balance', 'Exactly the arithmetically allowed number succeed; totals equal oracle; no duplicate rows', 'node/k6 parallel script', orig)
        if 'security' in K:
            R('security negative', 'Cross-role and cross-tenant attempts for the affected route', 'tenants 4 and 8, roles owner/employee', '401/403/404; no foreign data; no raw DB error text', 'API battery', 'tests/runtime/ph17_security.mjs')
        if 'validation' in K or 'ui' in K:
            R('boundary/negative', 'Invalid, empty, boundary and oversized inputs for the affected field/endpoint', 'authenticated session', '4xx with a clear message (never 500/raw DB text); no partial record', 'API + Chrome', orig)
        open(os.path.join(d, 'regression-test-spec.md'), 'w', encoding='utf-8').write(f"# {fid} — Regression test specification\n\n" + md_table([(a, b, c_, e, g) for a, b, c_, p_, e, g, o in tests], ['test id', 'type', 'description', 'assertion', 'automation']) + f"\n\n## Preconditions per test\n" + '\n'.join('- %s: %s (original audit script: %s)' % (a, p_, o) for a, b, c_, p_, e, g, o in tests) + "\n\n## Pass criteria\nAll tests above pass and the independent reconciliation for this finding is zero. Status: NOT STARTED.\n")
        index_rows.append((fid, f['title'], sev, c, os.path.basename(d), len(phases), len(tests), len(bad)))
        summary.append((fid, len(good), len(bad)))
    # 57 / 59 / 58
    with open('57-future-remediation-task-register.csv', 'w', newline='', encoding='utf-8') as h:
        w = csv.DictWriter(h, fieldnames=['task_id', 'finding_id', 'severity', 'phase', 'objective', 'existing_paths', 'required_changes', 'dependencies', 'required_tests', 'acceptance', 'risk', 'status'])
        w.writeheader(); w.writerows(task_rows)
    with open('59-post-remediation-regression-test-register.csv', 'w', newline='', encoding='utf-8') as h:
        w = csv.DictWriter(h, fieldnames=['test_id', 'finding_id', 'severity', 'type', 'description', 'preconditions', 'assertion', 'automation', 'original_audit_script', 'status'])
        w.writeheader(); w.writerows(reg_rows)
    open('58-remediation-specification-index.md', 'w', encoding='utf-8').write('# 58 — Remediation specification index (generated by tests/build_finding_packages.py)\n\nEvery finding has a folder under `findings/` with 10 files. All remediation tasks are NOT STARTED.\n\n' + md_table([(fid, '[%s](findings/%s/remediation-spec.md)' % (t[:70], d), s, cf[:12], ph, rt, ub) for fid, t, s, cf, d, ph, rt, ub in index_rows], ['finding', 'spec', 'severity', 'root-cause confidence', 'tasks', 'regression tests', 'unresolved source refs']) + '\n')
    print('findings', len(rows), 'tasks', len(task_rows), 'regression tests', len(reg_rows), 'unresolved refs total', sum(b for _, _, b in summary))
    return summary

if __name__ == '__main__':
    build()
