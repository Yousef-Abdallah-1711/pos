"""RUN-005: upgrades OBSERVED findings F-023, F-026, F-035, F-036 to REPRODUCED with RUN-005 runtime evidence (idempotent)."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
N17 = 'evidence/network/RUN-005/PH-17/AUD-R5-OBSERVED-REPRO.json'
N7 = 'evidence/network/RUN-005/PH-07/AUD-R5-POS-CURRENCY.json'
U = {
    'F-023': ('REPRODUCED (RUN-005 headed Chrome)', N7 + '; evidence/screenshots/RUN-005/PH-07/AUD-R5-POS-currency.png; tests/runtime/ph07_r5_pos_currency.mjs', ' RUN-005: the POS page shows $12.00, $10.00, $100.00, $10,000.00 while the dashboard home shows "EGP 110.00" for the same tenant (currency EGP).'),
    'F-026': ('REPRODUCED (RUN-005 API-DIAGNOSTIC)', N17 + '; tests/runtime/ph17_r5_observed_repro.mjs', ' RUN-005: card order 5795 (110.00) marked paid -> journal Dr 1113 Card Payment Clearing 110.00 / Cr 4100 100.00 / Cr 2120 10.00; no settlement route exists.'),
    'F-035': ('REPRODUCED (RUN-005 API-DIAGNOSTIC)', N17 + '; tests/runtime/ph17_r5_observed_repro.mjs', ' RUN-005: 25 bad logins -> 19 x 401 then 429 from the 20th attempt; the dashboard Content-Security-Policy (frontend/apps/dashboard/next.config.js:155) contains unsafe-inline and unsafe-eval in the isolated build; the API /health response carries no CSP header and the landing response carried none in this run.'),
    'F-036': ('REPRODUCED (RUN-005 API-DIAGNOSTIC)', N17 + '; tests/runtime/ph17_r5_observed_repro.mjs', ' RUN-005: unauthenticated GET /api/v1/mock/paymob/iframe -> 400 and POST /api/v1/mock/paymob/process -> 400 "payment_token is required": both execute their handler without any session.'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
ix = {k: i for i, k in enumerate(hd)}
for r in rows[1:]:
    u = U.get(r[0])
    if u:
        r[ix['evidence_status']] = u[0]
        if u[1] not in r[ix['evidence_paths']]:
            r[ix['evidence_paths']] += '; ' + u[1]
        if u[2].strip() not in r[ix['actual']]:
            r[ix['actual']] += u[2]
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    if 'EX-RUN005-OBSERVED-REPRO' not in ex:
        csv.writer(h).writerow(['EX-RUN005-OBSERVED-REPRO', 'API-DIAGNOSTIC + UI-HEADED-CHROME', 'Reproduction of OBSERVED findings F-023, F-026, F-035, F-036 (F-024 pending the invoice form run)', 'owner/anonymous', 'FAIL', 'all four reproduced: POS shows $ vs EGP; card payment debits 1113 with no settlement; login limit 20/min and CSP with unsafe-eval; mock Paymob routes execute without a session; evidence ' + N17 + ' and ' + N7])
print('ok')
