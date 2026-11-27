import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [('EX-RUN004-LANDING-SMOKE', 'UI-HEADED-CHROME', 'Landing smoke on the current uncommitted working tree: 12 routes x desktop and 390px, EN/AR, language switch, no form submitted', 'anonymous', 'PASS',
        'all routes 200 (unknown path 404), RTL on /ar, no horizontal overflow, language switch /en -> /ar works; OBSERVATION OBS-LAND-001: a React hydration mismatch warning on /en and /ar auth/register in the dev build (landing is under separate development; not a finding); evidence evidence/network/RUN-004/PH-00/AUD-R4-LANDING-SMOKE.json')]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    w = csv.writer(h)
    for r in new:
        if r[0] not in ex: w.writerow(r)
p = 'tests/r4_task_overlay.py'; s = open(p, encoding='utf-8').read()
s = s.replace("('T-R04-BASE-003', '00', 'landing', 'Verify landing registration/login/language switching on current tree (read-only)', 'PLANNED', 'landing /en, /ar, /en/auth/login return 200; registration/language switching UI check pending', '')",
              "('T-R04-BASE-003', '00', 'landing', 'Verify landing registration/login/language switching on current tree (read-only)', 'EXECUTED-PASS', 'smoke of 12 routes desktop+mobile EN/AR passed; language switch works; registration and login pages render (submission itself was exercised in RUN-003); hydration warning observed (OBS-LAND-001)', '')")
open(p, 'w', encoding='utf-8').write(s); print('ok')
