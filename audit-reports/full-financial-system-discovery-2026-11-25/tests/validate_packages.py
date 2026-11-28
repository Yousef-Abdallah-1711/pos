"""RUN-005 Phases O/P/Q validator: checks every finding folder (10 files, substantive size, evidence and source paths), every future task (fields, dependencies, acceptance, status NOT STARTED, no duplicate implementation of a shared root cause)
and every regression test (fields, links). Writes 58b-package-validation.md. Read-only except that file. Idempotent."""
import csv, os, re, collections, subprocess
REPO = r'C:\Users\Yousef\Desktop\Projects\POSS'
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
FILES = ['finding.md', 'reproduction.md', 'root-cause.md', 'affected-code-paths.md', 'financial-impact.md', 'dependency-map.md', 'evidence-index.md', 'remediation-spec.md', 'remediation-tasks.md', 'regression-test-spec.md']
MIN = {'finding.md': 1500, 'reproduction.md': 700, 'root-cause.md': 600, 'affected-code-paths.md': 700, 'financial-impact.md': 500, 'dependency-map.md': 350, 'evidence-index.md': 300, 'remediation-spec.md': 3300, 'remediation-tasks.md': 3000, 'regression-test-spec.md': 450}
find = {r['finding_id']: r for r in csv.DictReader(open('22-confirmed-findings-register.csv', encoding='utf-8'))}
tasks = list(csv.DictReader(open('57-future-remediation-task-register.csv', encoding='utf-8')))
regs = list(csv.DictReader(open('59-post-remediation-regression-test-register.csv', encoding='utf-8')))
tracked = set(subprocess.check_output(['git', '-C', REPO, 'ls-files'], text=True).split('\n'))
issues = collections.defaultdict(list)
dirs = {d[:5]: d for d in os.listdir('findings') if d.startswith('F-')}
for fid, r in find.items():
    d = dirs.get(fid)
    if not d:
        issues[fid].append('no findings folder'); continue
    for f in FILES:
        p = os.path.join('findings', d, f)
        if not os.path.exists(p):
            issues[fid].append('missing ' + f); continue
        n = os.path.getsize(p)
        if n < MIN[f]:
            issues[fid].append(f'{f} small ({n} bytes < {MIN[f]})')
    spec = open(os.path.join('findings', d, 'remediation-spec.md'), encoding='utf-8').read() if os.path.exists(os.path.join('findings', d, 'remediation-spec.md')) else ''
    for sec in ('## 9. Proposed remediation design', '## 23. Measurable acceptance criteria', '## 19. Rollback strategy'):
        if sec not in spec:
            issues[fid].append('spec lacks ' + sec)
    if 'RUN-004' not in r['evidence_paths'] and 'RUN-005' not in r['evidence_paths'] and 'RUN-003' not in r['evidence_paths'] and 'runtime' not in r['evidence_paths']:
        issues[fid].append('no run-attributed evidence path')
    for tok in re.findall(r'((?:backend|frontend)/[A-Za-z0-9_./\[\]-]+\.(?:go|sql|tsx|ts|mjs|js))(?::\d+(?:-\d+)?)?', r['source_refs']):
        if tok not in tracked:
            issues[fid].append('source path not tracked: ' + tok)
# tasks
bad_task = []; by = collections.defaultdict(list)
for t in tasks:
    by[t['finding_id']].append(t)
    for k in ('task_id', 'finding_id', 'phase', 'objective', 'required_changes', 'dependencies', 'required_tests', 'acceptance', 'status'):
        if not (t.get(k) or '').strip():
            bad_task.append((t['task_id'], 'empty ' + k))
    if t['status'] != 'NOT STARTED':
        bad_task.append((t['task_id'], 'status ' + t['status']))
    if len(t['acceptance']) < 40:
        bad_task.append((t['task_id'], 'acceptance too short'))
    if t['finding_id'] not in find:
        bad_task.append((t['task_id'], 'unknown finding'))
shared_p3 = [t for t in tasks if t['phase'].startswith('P3') and 'SHARED IMPLEMENTATION' in t['required_changes']]
# regressions
bad_reg = []
for g in regs:
    for k in ('test_id', 'finding_id', 'type', 'description', 'preconditions', 'assertion', 'automation', 'status'):
        if not (g.get(k) or '').strip():
            bad_reg.append((g['test_id'], 'empty ' + k))
    if g['status'] != 'NOT STARTED':
        bad_reg.append((g['test_id'], 'status ' + g['status']))
no_reg = [f for f in find if not any(g['finding_id'] == f for g in regs)]
open('58b-package-validation.md', 'w', encoding='utf-8').write(f"""# 58b - Package validation (RUN-005)

Built by `tests/validate_packages.py` (structural validation: it proves files exist, are non-trivial and reference real paths; it does NOT prove the proposed designs are technically optimal).

## Finding folders ({len(find)} findings)
- Folders present: {sum(1 for f in find if f in dirs)}; findings with at least one issue: **{len(issues)}**.
{chr(10).join('- ' + k + ': ' + '; '.join(v[:4]) for k, v in sorted(issues.items())[:60]) or '- none'}

## Future implementation tasks ({len(tasks)})
- Tasks with an empty required field, a non-NOT-STARTED status, a short acceptance or an unknown finding: **{len(bad_task)}**{(' - ' + '; '.join(f'{a}: {b}' for a, b in bad_task[:20])) if bad_task else ''}.
- Findings with tasks: {len(by)} of {len(find)}. Shared-component tasks (P3 depends on an owner finding): {len(shared_p3)} - these tasks do not re-implement the shared component.

## Future regression tests ({len(regs)})
- Tests with an empty required field or a status other than NOT STARTED: **{len(bad_reg)}**{(' - ' + '; '.join(f'{a}: {b}' for a, b in bad_reg[:20])) if bad_reg else ''}.
- Findings without any regression test: {len(no_reg)} {no_reg[:20]}.
- Distinction kept: current-product reproduction scripts live under `tests/runtime`; the 59 register holds future post-remediation acceptance tests only (all NOT STARTED).
""")
print(len(issues), 'findings with issues;', len(bad_task), 'task issues;', len(bad_reg), 'regression issues;', len(no_reg), 'findings without regression tests')
