"""RUN-004 register update 2 (idempotent): F-071 (SEC-013 staging gap) + 34b row. Audit tooling only."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
NET = 'evidence/network/RUN-004/PH-17/'
F = {
    'F-071': dict(
        title='ENV=staging is treated as development: RBAC defaults to audit-only (cashier reaches finance, HR and user endpoints), mock payments, open metrics and in-memory sessions are allowed',
        category='Security / configuration', severity='High', evidence_status='REPRODUCED',
        affected_modules='Platform configuration (RBAC, payments, sessions, metrics, cookies)',
        reproduction='Start the isolated API image with ENV=staging and RBAC_ENFORCEMENT unset (same DB/Redis) next to the enforce-mode API; log in as the cashier on both and request the same endpoints; production-mode probe for contrast (API-DIAGNOSTIC)',
        expected='A staging deployment enforces the same guards as production (RBAC enforce, no mock payments, protected metrics, Secure cookies), or the weaker behaviour is explicitly documented',
        actual='Enforce API: cashier gets 403 on /accounting/journal-entries, /accounting/accounts, /hr/salaries, balance sheet, profit-loss, /users. Staging API (RBAC unset): the same requests return 200 and a cashier write reaches the validator (400, not 403). ENV=production refuses to start without PAYMOB_MODE=live; staging accepts mock payments. /health/metrics answers 200 without a token and the session cookie has no Secure flag in both non-production runs. Staging does require independent signing secrets (signing_secrets.go:9) - the only guard staging shares with production',
        monetary_impact_EGP='0 direct; exposes all financial data and writes to any signed-in role on a staging deployment that holds real or production-like data',
        root_cause_first_break='CONFIRMED (source + runtime): config.IsProductionMode() is true only for production/prod (runtime_env.go:20-23); RBAC default mode is audit unless production (middleware/permission.go:36-42); mock Paymob allowed unless production (config/paymob_mode.go:29); metrics token and Redis/session requirements and secure cookies are gated the same way (health_routes.go:52, main.go:148-156, auth_handler.go:223)',
        source_refs='backend/internal/config/runtime_env.go:20; backend/internal/middleware/permission.go:36; backend/internal/config/paymob_mode.go:29; backend/cmd/api/health_routes.go:52; backend/cmd/api/main.go:705; backend/cmd/api/signing_secrets.go:9',
        related_static_ids='SEC-013; F-020; F-033',
        evidence_paths=NET + 'AUD-R4-SEC013.json; ' + NET + 'AUD-R4-SEC013-production-probe.log; tests/runtime/ph17_sec013.mjs',
        remediation_direction='Treat staging as production-like for guards (explicit allow-list of development environments instead of a production allow-list), refuse RBAC_ENFORCEMENT=audit and mock payments outside development, require the metrics token and Secure cookies when not development',
        regression_tests_needed='Start with ENV=staging and no RBAC variable: cashier requests to finance endpoints are 403; mock payment mode is rejected; metrics require a token'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
have = {r[0] for r in rows}
for fid, v in F.items():
    if fid not in have:
        v = dict(v, finding_id=fid)
        rows.append([v.get(k, '') for k in hd])
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
if 'EX-RUN004-SEC013' not in ex:
    with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
        csv.writer(h).writerow(['EX-RUN004-SEC013', 'API-DIAGNOSTIC', 'SEC-013 ENV=staging vs enforce-mode API (cashier matrix) and production startup probe', 'cashier', 'FAIL',
                                'staging default: 6 finance/HR/user endpoints 200 for cashier (enforce: 403); mock Paymob allowed; metrics open; cookie not Secure (F-071)'])
print(len(rows) - 1, 'findings')
