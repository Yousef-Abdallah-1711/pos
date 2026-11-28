"""RUN-005 durable checkpoint writer: replaces the '### RUN-005 checkpoint' block in progress.md and writes 62-next-model-handoff.md RUN-005 block (idempotent)."""
import csv, os, re, collections, subprocess, datetime
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
find = list(csv.DictReader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
conf = collections.Counter((re.match(r'\s*(CONFIRMED|PROBABLE)', x['root_cause_first_break']) or [None, 'UNSTATED'])[1] for x in find)
ev = collections.Counter(x['evidence_status'].split(' ')[0] for x in find)
t51 = collections.Counter(r['status'] for r in csv.DictReader(open('51-phased-audit-task-register.csv', encoding='utf-8')))
now = datetime.datetime.utcnow().strftime('%Y-%m-%dT%H:%MZ')
fr = [f for f in os.listdir('evidence/network/RUN-005/PH-06') if f.startswith('AUD-R5-FIELD-') and f.endswith('.json') and 'partial' not in f] if os.path.isdir('evidence/network/RUN-005/PH-06') else []
nf = 0; ncases = 0
import json
for f in fr:
    try:
        j = json.load(open('evidence/network/RUN-005/PH-06/' + f, encoding='utf-8'))['results']
        ncases += len(j); nf += len({(r['form'], r['field_index']) for r in j})
    except Exception:
        pass
B = f"""### RUN-005 checkpoint ({now})
Baseline 64-run005-baseline.md. Done so far: storage safety and verified backup (55 MB, sha256 in 64); audit folder restored from branch audit/financial-discovery-run004 after a branch switch deleted it (do not switch branches blindly); isolated stack rebuilt on the persistent volume; start snapshot `evidence/database/isolated-audit-run005-start-2026-11-27.dump`; route/API reconciliation (06b: 117 pages + 2 handlers = 119; API 355 rows after adding 6); metric classes (30d: 61 correct, 23 incorrect, 6 policy, 17 insufficient, 12 blocked); the four PROBABLE root causes confirmed (F-018 redis timing 68 s, F-021 replay, F-038 source, F-074 runtime) - root cause now {dict(conf)}; five OBSERVED findings: F-023/F-026/F-035/F-036 reproduced, F-024 pending invoice form run; authorization matrix on 4 roles x enforce/staging (F-071 expanded); package validation (58b); finding.md/root-cause.md/task acceptance enriched; field-isolation runner ph06_r5_field_isolation.mjs running ({nf} fields / {ncases} cases recorded so far in evidence/network/RUN-005/PH-06/AUD-R5-FIELD-*.json).
Evidence status counts: {dict(ev)}. Register 51: {dict(t51)}.
Resume instructions: see 62 (RUN-005 block). Not yet done: field isolation for all forms, destructive controls, Arabic and mobile interaction, invoice/payment/POS/settings form matrices, report date-boundary and export verification, browser benchmark, historical regression leftovers, final regeneration (generators order in 62), final dump and restoration check.
"""
p = open('progress.md', encoding='utf-8').read()
p = re.sub(r'\n### RUN-005 checkpoint.*?(?=\n### |\Z)', '', p, flags=re.S)
open('progress.md', 'w', encoding='utf-8').write(p.rstrip('\n') + '\n\n' + B)
h = open('62-next-model-handoff.md', encoding='utf-8').read()
H = f"""
## RUN-005 block (updated {now})
- Read 64-run005-baseline.md first. The audit folder is untracked on the current branch: never switch branches without first copying it; backup in `C:\\Users\\Yousef\\audit-backups\\` (tar.gz + sha256).
- Isolated stack: compose project `pos-p2a-ws05-e2e` (images were rebuilt in RUN-005 because the earlier ones had been removed). Start: from the repo root run docker compose with `--project-directory .`, `-p pos-p2a-ws05-e2e`, `-f audit-reports/full-financial-system-discovery-2026-11-25/tests/runtime/docker-compose.audit-persistent.yml` and `P2A_CORS_ORIGINS` including `http://*.localhost:13113` (`up -d --build postgres redis init-roles api frontend-install landing dashboard restaurant-website super-admin`).
- `tests/runtime/chrome_lib.mjs` now defaults RUN_ID to RUN-005. Rate limit: 20 logins per minute per IP - bursts of failed logins lock out the audit's own Chrome login for a minute.
- Field-isolation results: `evidence/network/RUN-005/PH-06/AUD-R5-FIELD-<form>.json` (one record per field and case). To resume a form: `ONLY=<key> node tests/runtime/ph06_r5_field_isolation.mjs` (keys: customers suppliers accounts drivers products categories roles users tables contracts attendance).
- Roles for authorization tests: owner audit.c.owner@example.test (ADMIN), manager.r5@example.test (MANAGER), cashier.c@example.test (EMPLOYEE), driver.r5@example.test (DRIVER); password in tests/runtime/chrome_lib.mjs (synthetic).
"""
if '## RUN-005 block' in h:
    h = re.sub(r'\n## RUN-005 block.*', '', h, flags=re.S)
open('62-next-model-handoff.md', 'w', encoding='utf-8').write(h.rstrip('\n') + '\n' + H)
print('checkpoint written; fields', nf, 'cases', ncases)
