"""Builds 73-SINGLE-FILE-PHASE-TASK-ISSUE-MAP.md: one file with all 23 phases (status, tasks, remaining queue, Definition of Done)
and all findings (closure criteria, tasks, regression tests, package path). Reads registers only. Audit tooling; idempotent."""
import csv, collections, re, os, glob
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
rd = lambda f: list(csv.DictReader(open(f, encoding='utf-8')))
R51 = rd('51-phased-audit-task-register.csv'); FIND = rd('22-confirmed-findings-register.csv'); REM = rd('57-future-remediation-task-register.csv')
REG = rd('59-post-remediation-regression-test-register.csv'); Q = rd('66-remaining-audit-execution-queue.csv')
PH = {}
for i in range(23):
    t = open(f'phases/phase-{i:02d}/spec.md', encoding='utf-8').readline().strip()
    PH[f'{i:02d}'] = re.sub(r'^# Phase \d+ [—-] ', '', t)
def expand(s):
    out = set()
    for a, b in re.findall(r'F-(\d{3})\.\.F-(\d{3})', s):
        out |= {f'F-{n:03d}' for n in range(int(a), int(b) + 1)}
    out |= set(re.findall(r'F-\d{3}', s))
    return out
f2p = collections.defaultdict(set)
for r in R51:
    for f in expand(r['finding_ids']):
        f2p[f].add(r['phase'])
KW = (('super-admin', '13'), ('tenant', '13'), ('payroll', '10'), ('hr', '10'), ('employee', '10'), ('leave', '10'), ('kds', '08'), ('storefront', '08'), ('checkout', '08'), (' tables', '08'), ('profit', '12'),
      ('pos', '07'), ('order', '07'), ('customer', '09'), ('supplier', '06'), ('product', '06'), ('inventory', '06'), ('report', '12'), ('rbac', '17'), ('security', '17'), ('auth', '17'),
      ('redis', '18'), ('performance', '18'), ('accounting', '11'), ('journal', '11'), ('invoice', '11'), ('bank', '11'), ('theme', '05'), ('settings', '05'), ('onboarding', '05'), ('register', '05'))
def phase_of_finding(f):
    if f2p.get(f['finding_id']):
        return sorted(f2p[f['finding_id']])
    s = (f['affected_modules'] + ' ' + f['title']).lower()
    for k, p in KW:
        if k in s:
            return [p]
    return ['11']
st = collections.defaultdict(collections.Counter)
for r in R51:
    st[r['phase']][r['status']] += 1
rem_by_f = collections.defaultdict(list); reg_by_f = collections.defaultdict(list)
for r in REM: rem_by_f[r['finding_id']].append(r)
for r in REG: reg_by_f[r['finding_id']].append(r)
fp = {}
for d in glob.glob('findings/F-*'):
    fp[os.path.basename(d)[:5]] = d.replace('\\', '/')
L = []
L.append('# 73 - SINGLE FILE: phases, tasks, Definition of Done and every issue (start here with a cheaper agent)\n')
L.append('One file to start from. It is generated from the registers by `tests/build_single_file_map.py` (registers win on disagreement). Companion files: queue `66-remaining-audit-execution-queue.csv`, per-item lists 66b/66c/66d, start prompt `72-NEXT-AGENT-START-PROMPT.md`, guide `68`.\n')
L.append('**How to use:** (1) pick the next NOT-STARTED row in section 3 or in 66; (2) execute it with its script; (3) record evidence and results; (4) tick its Definition of Done; (5) when every audit row of a phase is done, evaluate the phase DoD in section 2. Issues (section 4) are **closed for the audit** when reproduced, root-caused and specified (already true for all 81); they are **closed for the product** only after the future remediation tasks pass their regression tests (Workstream B, NOT STARTED, never done by the audit agent).\n')
L.append('## 1. Phase dashboard\n')
L.append('| Phase | Name | Status | Tasks (51) | PASS | FAIL (defect recorded) | PARTIAL | BLOCKED | Queue rows not started | Issues |')
L.append('|---|---|---|---|---|---|---|---|---|---|')
qn = collections.Counter(q['Phase ID'] for q in Q if q['Status'] == 'NOT-STARTED')
pf = collections.defaultdict(list)
for f in FIND:
    for p in phase_of_finding(f):
        pf[p].append(f['finding_id'])
def pstat(p):
    c = st[p]
    if c['PLANNED'] and not (c['EXECUTED-PASS'] + c['EXECUTED-FAIL'] + c['EXECUTED-PARTIAL']): return 'UNSTARTED'
    if qn[p] or c['EXECUTED-PARTIAL'] or c['PLANNED']: return 'PARTIAL'
    return 'EXECUTED'
for p in PH:
    c = st[p]
    L.append(f"| {p} | {PH[p]} | {pstat(p)} | {sum(c.values())} | {c['EXECUTED-PASS']} | {c['EXECUTED-FAIL']} | {c['EXECUTED-PARTIAL']} | {c['BLOCKED-CAPABILITY'] + c['BLOCKED-DEFECT']} | {qn[p]} | {len(pf[p])} |")
L.append('\n## 2. Phase details: tasks and Definition of Done\n')
for p in PH:
    rows = [r for r in R51 if r['phase'] == p]
    qs = [q for q in Q if q['Phase ID'] == p and q['Status'] == 'NOT-STARTED']
    L.append(f"### Phase {p} - {PH[p]}  [{pstat(p)}]")
    L.append(f"Folder: `phases/phase-{p}/` (README, spec, tasks, definition-of-done, evidence-index, dependencies, execution-runbook).")
    L.append('**Audit tasks of record (register 51):**')
    for r in rows:
        L.append(f"- {r['task_id']} [{r['status']}] {r['objective'][:130]}" + (f" - findings {r['finding_ids']}" if r['finding_ids'] else '') + (f" - next: {r['next_action'][:90]}" if r['next_action'] else ''))
    L.append(f"**Remaining queue rows ({len(qs)}):** " + ('; '.join(f"{q['Task ID']} ({q['Required model capability']})" for q in qs[:12]) + (f' ... +{len(qs) - 12} more in 66' if len(qs) > 12 else '') if qs else 'none'))
    L.append(f"**Definition of Done:** every row above is EXECUTED-PASS, EXECUTED-FAIL (task done, defect recorded) or BLOCKED with a proven prerequisite; every queue row of the phase executed or BLOCKED with proof; evidence exists and resolves; findings registered in 22; discrepancies reconciled or marked unexplained. Now: **{pstat(p)}**.\n")
L.append('## 3. Execution order for the cheaper agent (first 25 rows of the queue)\n')
L.append('| Task | Phase | Capability | Title | Script | Evidence destination |\n|---|---|---|---|---|---|')
for q in Q[:25]:
    L.append(f"| {q['Task ID']} | {q['Phase ID']} | {q['Required model capability']} | {q['Title'][:90]} | {q['Execution script'][:70]} | {q['Evidence destination'][:50]} |")
L.append(f"\nThe full queue has {len(Q)} rows (66). Field batches Q-FLD-* (437 fields), destructive controls Q-DST-* (44), blocked controls Q-BLK-* (58) are listed item by item in 66b, 66c, 66d.\n")
L.append('## 4. Every issue (finding): closure map\n')
L.append('Audit status: all findings REPRODUCED. Product status: all NOT FIXED (remediation NOT STARTED). **Closed for the product when:** the root cause is fixed per the remediation spec, every regression test listed passes, and the original reproduction script no longer shows the failure.\n')
L.append('| ID | Sev | Phase(s) | Title | Root cause | Remediation tasks (57) | Regression tests (59) | Package | Closed when |')
L.append('|---|---|---|---|---|---|---|---|---|')
for f in FIND:
    fid = f['finding_id']
    rc = (re.match(r'\s*(CONFIRMED|PROBABLE)', f['root_cause_first_break']) or [0, 'UNSTATED'])[1]
    rem = rem_by_f.get(fid, []); reg = reg_by_f.get(fid, [])
    v = [r for r in rem if r['phase'].startswith('P8')]
    closed = (v[0]['acceptance'] if v else f['regression_tests_needed'])[:140].replace('|', '/').replace('\n', ' ')
    L.append(f"| {fid} | {f['severity'].split(' ')[0]} | {','.join(phase_of_finding(f))} | {f['title'][:85].replace('|', '/')} | {rc} | {len(rem)} | {len(reg)} | {fp.get(fid, 'findings/' + fid)} | {closed} |")
L.append('\nPer-finding detail (reproduction, root cause with file:line, financial impact, dependencies, remediation spec, tasks, regression spec) is in the package folder named in the table. Implementation order and policy gates: `69-FUTURE-REMEDIATION-OVERVIEW.md`, `56b-policy-decision-detail.md` (22 decisions, all OPEN).\n')
L.append('## 5. Rules (short)\nAudit only: no product fixes, no remediation tasks, no policy approvals, no subagents. Never open `evidence/state-*.json`. Protected tenants 5/6/7 stay unseeded; tenant 10 read-only; tenant 4 is a mutated diagnostic fixture. Release stays NO-GO and audit PARTIAL until the checklist in 71 is evaluated with evidence.\n')
open('73-SINGLE-FILE-PHASE-TASK-ISSUE-MAP.md', 'w', encoding='utf-8').write('\n'.join(L))
print(len(L), 'lines;', len(FIND), 'findings')
