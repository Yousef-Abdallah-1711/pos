"""RUN-004 finalisation (run LAST, after the generators): 34b rows, 61 section 8, 62 handoff rewrite, README/progress entries. Audit tooling only."""
import csv, os, re, collections, subprocess
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
# 34b: landing smoke repeat
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    if 'EX-RUN004-LANDING-SMOKE-HEAD' not in ex:
        csv.writer(h).writerow(['EX-RUN004-LANDING-SMOKE-HEAD', 'UI-HEADED-CHROME', 'Landing smoke repeated against HEAD 0de839b (final committed tree), same 12 routes x desktop/390px', 'anonymous', 'PASS', 'identical results to the pre-commit run (same status codes, text lengths, directions, no overflow); hydration warning on register pages unchanged (OBS-LAND-001); evidence evidence/network/RUN-004/PH-00/AUD-R4-LANDING-SMOKE.json (pre-commit copy: ...pre-commit-tree.json)'])
find = list(csv.DictReader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
sev = collections.Counter(x['severity'].split(' ')[0] for x in find)
conf = collections.Counter((re.match(r'\s*(CONFIRMED|PROBABLE|CANDIDATE|UNKNOWN)', x['root_cause_first_break']) or [None, 'UNSTATED'])[1] for x in find)
tasks = list(csv.DictReader(open('57-future-remediation-task-register.csv', encoding='utf-8')))
regs = list(csv.DictReader(open('59-post-remediation-regression-test-register.csv', encoding='utf-8')))
t51 = collections.Counter(r['status'] for r in csv.DictReader(open('51-phased-audit-task-register.csv', encoding='utf-8')))
head = subprocess.check_output(['git', '-C', r'C:\Users\Yousef\Desktop\Projects\POSS', 'rev-parse', 'HEAD'], text=True).strip()
S8 = f"""
### RUN-004 environment closure ({'2026-11-27'})
- Final isolated snapshot: `evidence/database/isolated-audit-run004-final-2026-11-27.dump` (sha256 in `evidence/database/RUN-004/T-R04-checkpoint-final.txt`); start and mid snapshots and all earlier dumps preserved. Protected tenants 5, 6, 7 remain unseeded (COA 0). Isolated stack stopped, all volumes kept.
- Development database `pos_saas`: **unchanged** (tenants 2, users 6, orders 0, journals 0, COA 62, migrations 145 = baseline in `evidence/runtime/RUN-003/46-dev-db-fingerprint-before.txt`; record `evidence/runtime/RUN-004/02-restoration-and-dev-fingerprint-after.txt`).
- Git: HEAD `{head}` (the landing workstream committed three landing commits during RUN-004; backend untouched), two stashes intact, only 6 tooling files modified (tsconfig/next-env). **Product files modified by the audit: 0.**
- **Development containers NOT restored.** `pos-backend`, `pos-landing`, `pos-dashboard`, `pos-restaurant-website`, `pos-super-admin` are exited; `docker start pos-backend` fails with "Bind for 127.0.0.1:3001 failed: port is already allocated" (pos-backend owns the shared network namespace and publishes 3001-3004). Holders: containers `pos-landing-ci-check` (127.0.0.1:3001) and `pos-landing-ci-recheck` (127.0.0.1:3004), `pnpm --filter ./apps/landing start`, started 03:01Z and 03:14Z, no log output in the last 15 minutes, apparently landing build checks of the finished frontend work. My attempt to stop them was declined by the permission classifier and was not retried; they are untouched. **To finish restoration:** `docker stop pos-landing-ci-check pos-landing-ci-recheck` then `docker start pos-backend pos-landing pos-dashboard pos-restaurant-website pos-super-admin`, then re-verify the dev URLs.
- Earlier in RUN-004 the port-3001 holder was a host `next start` (pid 9064, landing); it was stopped under the user's explicit authorization after identity verification (see 63).
"""
s = open('61-final-exhaustive-audit-closure-report.md', encoding='utf-8').read()
s = re.sub(r'## 8\. Environment restoration\n.*?(?=\n## 9\.)', '## 8. Environment restoration' + S8, s, flags=re.S)
s = s.replace('| Original environment restored | SEE §8 | |', '| Original environment restored | NOT MET | dev database verified unchanged; dev containers still stopped (port 3001/3004 held by two landing CI containers, see §8) |') if '| Original environment restored | SEE §8 | |' in s else s
open('61-final-exhaustive-audit-closure-report.md', 'w', encoding='utf-8').write(s)
H62 = f"""# 62 - Next-model handoff (updated at the end of RUN-004, 2026-11-27)

Resume order: README.md -> progress.md -> 63-resume-baseline-run004.md -> this file -> 61 -> 51 -> 45/48 addenda -> current phase tasks.

## State
- Audit directory: `...\\audit-reports\\full-financial-system-discovery-2026-11-25\\` (the folder is untracked by git: it exists only on this machine; back it up).
- Git HEAD `{head}` (audit baseline for RUN-001..003 was `133260db...`). No product code modified by the audit. Two stashes preserved.
- Findings {len(find)} (Critical {sev['Critical']}, High {sev['High']}, Medium {sev['Medium']}, Low {sev['Low']}; F-051 unused); root cause {conf.get('CONFIRMED', 0)} CONFIRMED / {conf.get('PROBABLE', 0)} PROBABLE (F-018, F-021, F-038, F-074: each states the exact gap). Retracted claim: duplicate-table-number part of F-064.
- Future tasks {len(tasks)} (NOT STARTED; only applicable phases per finding; shared root causes depend on one owner task via `SHARED_OWNER` in `tests/build_finding_packages.py`); future regression tests {len(regs)}; policy decisions 19 (POL-16..19 are proposals; none approved).
- Register 51: {dict(t51)}. Overall classification: PARTIAL. Financial release: NO-GO.

## Environment
- Compose project `pos-p2a-ws05-e2e`, persistent volume `pos-p2a-ws05-e2e_p2a_audit_pgdata`, DB `pos_phase2a_e2e`. Start from the repo root with `--project-directory .` and `P2A_CORS_ORIGINS` including `http://*.localhost:13113` (compose file `tests/runtime/docker-compose.audit-persistent.yml`); recreating the API requires recreating the four frontends. URLs: API 18082, landing 13111, dashboard 13112, storefront 13113, super-admin 13114.
- The isolated stack is currently STOPPED. Dev containers are stopped because ports 3001/3004 are held by `pos-landing-ci-check` and `pos-landing-ci-recheck` (see 61 section 8 for the exact restoration commands).
- Tenants: 4 diagnostic fixture (heavily mutated: 5,794 orders, themes 4/11/13, products 12-18, customers/suppliers/accounts/roles/users created by the form matrix); 8 and 9 earlier fixtures; **10 = synthetic benchmark tenant (54,000 posted lines, label SYNTHETIC BENCHMARK FIXTURE, owner audit.g.owner@example.test)**; 5, 6, 7 protected unseeded (F-001 evidence); 3 unseeded.
- Snapshots (sha256 files in `evidence/database/RUN-004/`): start, mid, final 2026-11-27; earlier dumps untouched.
- Never open `evidence/state-*.json` (session cookies).

## Regenerate derived files (order matters; r4_* scripts are RUN-004 overlays)
`python -I tests/build_finding_evidence_check.py` -> `python -I tests/build_finding_packages.py` (reads `tests/finding_overrides.json` and `tests/finding_overrides_r4*.py`) -> `python -I tests/task_status_overlay.py` then `python -I tests/r4_task_overlay.py` (51) -> `python -I -X utf8 tests/build_phase_specs.py` -> `python -I -X utf8 tests/build_field_denominator.py` -> `python -I -X utf8 tests/build_control_interaction_results.py` -> `python -I tests/build_evidence_manifest.py` -> `python -I tests/build_closure_report.py` -> `python -I -X utf8 tests/r4_finalize.py`. Register edits for RUN-004 are in `tests/r4_register_update*.py` and `tests/r4_rootcause_update.py` (idempotent).
- Tooling notes: heredocs containing apostrophes break the Bash tool, write files with the Write tool; `python -I` ignores PYTHONIOENCODING, use `-X utf8`; Chrome profiles for the audit live in the session scratchpad (separate profile per concurrent runner); Chrome helper `tests/runtime/chrome_lib.mjs` honours `RUN_ID` (default RUN-004).

## Remaining work (feasible, in priority order)
1. Finish environment restoration (61 section 8) and re-run the landing smoke if the landing changes again.
2. Single-field isolation tests for the 531 logical fields (54b), starting with money forms; extend the form matrix to invoices, payments, bank accounts, POS, settings sub-forms, storefront checkout, super-admin; add the forms missing from the field inventory (driver, user, role, employee fields).
3. Arabic and mobile interaction passes for the control matrix; rerun the 58 tooling-blocked controls (increase timeouts, one runner at a time); purpose-built fixtures for the 88 destructive controls.
4. Confirm root causes F-018, F-021, F-038, F-074; audit other date-bounded reports for the UTC cut-off (F-070 candidate); measure browser rendering on the 54k tenant.
5. Reconcile route counts 119 vs 117 (T-P03-002) and the 10 dynamic routes.
6. After remediation exists (not started): execute the 233 regression tests and re-run blocked workflows (payroll, period close, refunds, subscription billing, depreciation).

## Do not
Seed tenants 5, 6 or 7; commit; modify product code; trust SQL agreement alone as financial correctness; mark future regression tests as passed.
"""
open('62-next-model-handoff.md', 'w', encoding='utf-8').write(H62)
PR = f"""
### RUN-004 closure (2026-11-27)
Resumed with the finish-the-audit authorization (no remediation, no commits). Baseline in 63 (HEAD advanced to {head[:7]} during the run: landing commits only). Executed: all six remaining metrics (F-066..F-070), SEC-013 (F-071), write load / payment burst / 54k-line benchmark (all PASS, reconciled), bank reconciliation (F-072, F-073), real-UI form matrix over 17 forms x 5 modes (F-074, F-075, F-079), per-control interaction matrix (644 controls; F-076, F-077), Theme Builder UI (F-077, F-078), landing smoke (twice), Super Admin UI partial. Root causes: {conf.get('CONFIRMED', 0)} CONFIRMED / {conf.get('PROBABLE', 0)} PROBABLE. Remediation packages rewritten with hand-written designs; tasks {len(tasks)}, regression tests {len(regs)}. Findings now {len(find)}. Dev DB unchanged; dev containers NOT restored (ports 3001/3004 held by two landing CI containers; stop attempt declined, see 61 section 8). Overall PARTIAL; financial release NO-GO.
"""
p = open('progress.md', encoding='utf-8').read()
if '### RUN-004 closure' not in p:
    open('progress.md', 'a', encoding='utf-8').write(PR)
r = open('README.md', encoding='utf-8').read()
if 'RUN-004' not in r:
    open('README.md', 'a', encoding='utf-8').write("\n## RUN-004 (2026-11-27)\nSee [63-resume-baseline-run004.md](63-resume-baseline-run004.md), [53b-ui-control-interaction-results.md](53b-ui-control-interaction-results.md), [54b-field-denominator-and-form-coverage.md](54b-field-denominator-and-form-coverage.md), the addenda in 45/48, and [62-next-model-handoff.md](62-next-model-handoff.md). Evidence for the run is under `evidence/*/RUN-004/`.\n")
print('finalised')
