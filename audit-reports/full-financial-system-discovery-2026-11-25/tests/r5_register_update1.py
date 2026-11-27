"""RUN-005 root-cause confirmations for F-018, F-021, F-038, F-074 (idempotent). Audit tooling only."""
import csv, os
os.chdir(r'C:\Users\Yousef\Desktop\Projects\POSS\audit-reports\full-financial-system-discovery-2026-11-25')
D = 'evidence/database/RUN-005/PH-06/'
N18 = 'evidence/network/RUN-005/PH-18/'
N6 = 'evidence/network/RUN-005/PH-06/'
U = {
    'F-018': dict(
        rc='CONFIRMED (RUN-005 source + fault injection): AuthSessionStore.Resolve calls the Redis cache with the request context, which has no deadline (backend/internal/middleware/auth.go:60 -> backend/internal/service/auth_session.go:166-171 -> RedisCache.Get cache.go:248); the Redis client is built from ParseURL defaults with no explicit dial, read, write or pool timeouts (cache.go:232-237); the middleware intends a 503 on a resolver error (auth.go:61-67) but the HTTP server WriteTimeout of 60 s (cmd/api/main.go:2004) closes the connection first. Measured with the isolated Redis stopped: GET /products with a valid cookie failed after 68,179 ms and POST /auth/login after 68,650 ms, both with NO HTTP response (client "fetch failed"); /health answered 200 in 5 ms and the public theme endpoint answered in 45 ms. After Redis restarted the old cookie returned 401 (all sessions lost) and a new login worked. Residual not isolated: the exact split of the 68 s between go-redis retries and pool waits',
        refs='backend/internal/middleware/auth.go:60; backend/internal/service/auth_session.go:166; backend/cmd/api/main.go:2004', ev=N18 + 'AUD-R5-F018-REDIS-TIMING.json; tests/runtime/ph18_r5_redis_timing.mjs'),
    'F-021': dict(
        rc='CONFIRMED (RUN-005 source + replay): suppliers.current_balance has exactly one money-movement writer: a supplier payment subtracts the amount (backend/internal/repository/payment_acct_repo.go:302-306), while supplier invoice creation and approval never add to it (no other UPDATE suppliers exists; supplier_repo.go:172-200 writes the value supplied by the handler, which preserves the stored value, supplier_handler.go:73). Replay on a new supplier: invoice 10,000 approved -> current_balance 0.00 while AP 2101 rose by 10,000; payments 4,000 and 6,000 -> current_balance -4,000.00 then -10,000.00 while the invoice balance_due went to 0 and AP returned to its prior 800.00. The stored balance therefore equals minus the sum of payments and ignores invoices',
        refs='backend/internal/repository/payment_acct_repo.go:302; backend/internal/repository/supplier_repo.go:172; backend/internal/handler/http/supplier_handler.go:73', ev=D + 'AUD-R5-F021-SUPPLIER-BALANCE.json; tests/runtime/ph06_r5_supplier_balance.mjs'),
    'F-038': dict(
        rc='CONFIRMED (RUN-005 source + RUN-003 runtime): in the current landing code (frontend/apps/landing/src/app/[locale]/auth/register/page.tsx) onSubmit sets the slug error with setError when availability is not available (lines 59-64) and nothing clears it later: the only useEffect (lines 20-24) handles billing and analytics, availability is read from useDebouncedAvailability (line 39) without an effect that calls clearErrors, and the schema resolver does not know availability, so the manual error persists until the field is edited or the form is submitted again. RUN-003 reproduced the persistence in headed Chrome (2 of 3 early-submit scenarios, error visible >= 8 s while the field showed "available"). The landing source moved under the committed landing workstream; line numbers are for HEAD 35100f7',
        refs='frontend/apps/landing/src/app/[locale]/auth/register/page.tsx:59; frontend/apps/landing/src/app/[locale]/auth/register/page.tsx:39', ev=''),
    'F-074': dict(
        rc='CONFIRMED (RUN-005 source + runtime): the parent customers page shows a toast and then rethrows the error (frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx:95-118, throw error at line 112), and the child form awaits onSubmit without a try/catch (components/accounting/CustomerForm.tsx:69-79; same pattern in SupplierForm.tsx:75-85 and AccountForm.tsx:81-117), so every server rejection becomes an unhandled promise rejection; the contracts page uses try/finally with no catch (app/[locale]/dashboard/hr/contracts/page.tsx:18-32) and ContractForm awaits without catch (components/hr-advanced/ContractForm.tsx:27-29), so no message is shown at all. Runtime (Next dev build): all four forms raised window unhandledrejection "Request failed with status code 400"; customers, suppliers and accounts showed a toast with the raw database text "pq: value too long for type character varying(100)/(20)"; contracts showed no toast',
        refs='frontend/apps/dashboard/src/app/[locale]/dashboard/accounting/customers/page.tsx:95; frontend/apps/dashboard/src/components/accounting/CustomerForm.tsx:69; frontend/apps/dashboard/src/app/[locale]/dashboard/hr/contracts/page.tsx:18; frontend/apps/dashboard/src/components/hr-advanced/ContractForm.tsx:27', ev=N6 + 'AUD-R5-F074-UNHANDLED.json; tests/runtime/ph06_r5_f074.mjs',
        title='Rejected submits on the customer, supplier, account and contract forms raise an unhandled promise rejection; three show the raw database error in a toast and the contract form shows nothing',
        actual='POST returns 400 (for example pq: value too long for type character varying(100)); in every case the page raises an unhandled promise rejection "Request failed with status code 400"; the customer, supplier and account forms show a toast containing the raw database text (see F-025) and the contract form shows no message at all; empty submits show client-side required messages, so only server-side rejections are affected'),
}
rows = list(csv.reader(open('22-confirmed-findings-register.csv', encoding='utf-8')))
hd = rows[0]
ix = {k: i for i, k in enumerate(hd)}
for r in rows[1:]:
    u = U.get(r[0])
    if not u:
        continue
    r[ix['root_cause_first_break']] = u['rc']
    for ref in u['refs'].split('; '):
        if ref not in r[ix['source_refs']]:
            r[ix['source_refs']] = (r[ix['source_refs']] + '; ' if r[ix['source_refs']] else '') + ref
    if u['ev'] and u['ev'] not in r[ix['evidence_paths']]:
        r[ix['evidence_paths']] += '; ' + u['ev']
    if 'title' in u:
        r[ix['title']] = u['title']
    if 'actual' in u:
        r[ix['actual']] = u['actual']
csv.writer(open('22-confirmed-findings-register.csv', 'w', newline='', encoding='utf-8')).writerows(rows)
ex = {r[0] for r in csv.reader(open('34b-executed-test-results.csv', encoding='utf-8'))}
new = [
    ('EX-RUN005-F018-REDIS-TIMING', 'API-DIAGNOSTIC+DOCKER', 'F-018 Redis outage timing: authenticated request, login, health, public endpoint, recovery', 'owner', 'FAIL', 'authenticated request and login failed after about 68 s with no HTTP response (WriteTimeout 60 s), /health 200 in 5 ms, public endpoint fine; after restart old cookie 401, new login works (sessions lost); root cause CONFIRMED; evidence ' + N18 + 'AUD-R5-F018-REDIS-TIMING.json'),
    ('EX-RUN005-F021-SUPPLIER-BALANCE', 'API-DIAGNOSTIC', 'F-021 supplier balance replay: invoice, two payments, balance versus AP', 'owner', 'FAIL', 'current_balance 0 after a 10,000 invoice, -10,000 after payments while invoice balance_due 0: balance equals minus payments; root cause CONFIRMED; evidence ' + D + 'AUD-R5-F021-SUPPLIER-BALANCE.json'),
    ('EX-RUN005-F074-UNHANDLED', 'UI-HEADED-CHROME', 'F-074 unhandled rejection on four forms (customer, supplier, account, contract) with an over-long value', 'owner', 'FAIL', 'all four raise an unhandled rejection; three toast raw pq text, contracts shows nothing; root cause CONFIRMED; evidence ' + N6 + 'AUD-R5-F074-UNHANDLED.json'),
]
with open('34b-executed-test-results.csv', 'a', newline='', encoding='utf-8') as h:
    for r in new:
        if r[0] not in ex:
            csv.writer(h).writerow(r)
print('ok')
