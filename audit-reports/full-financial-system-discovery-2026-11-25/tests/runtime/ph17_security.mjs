// RUN-003 / PH-17: security battery on the ISOLATED stack (API-DIAGNOSTIC). Synthetic accounts only.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-17');
const O = 'http://localhost:13112';
const mk = () => { let cookie = ''; const call = async (method, path, body, opts = {}) => { const r = await fetch(API + path, { method, headers: { ...(opts.noOrigin ? {} : { Origin: O }), Cookie: opts.cookie ?? cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined }); const sc = r.headers.getSetCookie?.() ?? []; if (sc.length && !opts.keep) cookie = sc.map((c) => c.split(';')[0]).join('; '); const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 160); } return { status: r.status, body: j, setCookie: sc, headers: r.headers }; }; return { call, get cookie() { return cookie; }, set cookie(v) { cookie = v; } }; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 360)); };
const A = mk(); const E = mk(); const C = mk();
await A.call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD });
const resE = await E.call('POST', '/auth/login', { email: 'audit.e.owner@example.test', password: PASSWORD });
const resC = await C.call('POST', '/auth/login', { email: 'cashier.c@example.test', password: PASSWORD });
rec('LOGIN', 'sessions', { owner4: !!A.cookie, ownerE: resE.status, cashier: resC.status });
// SEC-014 stub login
const s1 = await mk().call('POST', '/auth/login', { email: 'nobody@example.test', password: 'password123' }); const s2 = await mk().call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: 'password123' });
rec('SEC-014', 'stub-mode credentials password123', { unknownUser: s1.status, knownOwnerWrongPw: s2.status });
// cookie flags
const lg = await mk().call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD }); rec('SEC-COOKIE', 'Set-Cookie attributes', { setCookie: lg.setCookie.map((c) => c.replace(/=[^;]+/, '=<redacted>')) });
// SEC-011 tampered / empty cookie
const ck = A.cookie; const tampered = ck.slice(0, -3) + (ck.slice(-3) === 'AAA' ? 'BBB' : 'AAA');
const t1 = await mk().call('GET', '/accounting/trial-balance', null, { cookie: tampered }); const t2 = await mk().call('GET', '/accounting/trial-balance', null, { cookie: '' }); const t3 = await mk().call('GET', '/accounting/trial-balance', null, { cookie: 'poos_session=abc.def.ghi' });
rec('SEC-011', 'tampered/empty/garbage session cookie on financial route', { tampered: t1.status, empty: t2.status, garbage: t3.status });
// IDOR: tenant E (owner) reading tenant 4 objects by id
const ids = { invoice4: Number(sql("select max(id) from invoices_customer where tenant_id=4")), order: Number(sql("select max(id) from orders where tenant_id=4 and id<43")), journal: Number(sql("select max(id) from journal_entries where tenant_id=4")), supplier: 2, product: 2, acct: 70 };
const probes = [['GET', `/accounting/customer-invoices/${ids.invoice4}`], ['GET', `/admin/orders/${ids.order}`], ['GET', `/accounting/journal-entries/${ids.journal}`], ['GET', `/suppliers/${ids.supplier}`], ['GET', `/products/${ids.product}`], ['GET', `/accounting/ledger?account_id=${ids.acct}`], ['GET', '/accounting/accounts/70'], ['GET', `/admin/orders/${ids.order}/history`]];
const idor = {}; for (const [m, p] of probes) { const r = await E.call(m, p); idor[p] = { status: r.status, leaked: r.status === 200 && JSON.stringify(r.body).includes('"tenant_id":4') }; }
rec('SEC-IDOR-READ', 'tenant E reads tenant 4 objects by id', idor);
// IDOR writes
const w = {};
let r = await E.call('PUT', `/products/${ids.product}`, { name_en: 'HACKED' }); w['PUT product'] = r.status;
r = await E.call('POST', `/accounting/journal-entries/${ids.journal}/void`, { reason: 'x' }); w['void foreign journal'] = r.status;
r = await E.call('PUT', `/admin/orders/${ids.order}/payment-status`, { payment_status: 'paid' }); w['pay foreign order'] = r.status;
r = await E.call('POST', '/accounting/payments/received', { payment_date: '2026-11-26', customer_id: 2, invoice_id: ids.invoice4, amount: 1, payment_method: 'cash' }); w['receipt on foreign invoice'] = r.status;
r = await E.call('DELETE', `/suppliers/${ids.supplier}`); w['delete foreign supplier'] = r.status;
rec('SEC-IDOR-WRITE', 'tenant E writes to tenant 4 objects', { ...w, productNameNow: sql('select name_en from products where id=2'), foreignReceiptRows: sql(`select count(*) from payments_received where invoice_id=${ids.invoice4} and amount=1`) });
// mass assignment
const fd = new FormData(); fd.append('data', JSON.stringify({ category_id: 2, name_en: 'MassAssign', price: 1, tenant_id: 4, restaurant_id: 4 }));
const mr = await fetch(API + '/products', { method: 'POST', headers: { Origin: O, Cookie: E.cookie }, body: fd }); const mj = await mr.json().catch(() => ({}));
rec('SEC-MASS', 'tenant E creates product with tenant_id=4 in body', { status: mr.status, tenantOfCreated: sql("select tenant_id from products where name_en='MassAssign' order by id desc limit 1") });
// role: cashier calling accounting routes
const cw = {};
for (const [m, p, b] of [['GET', '/accounting/trial-balance'], ['POST', '/accounting/journal-entries', { entry_date: '2026-11-26', description: 'x', entry_type: 'manual', lines: [] }], ['GET', '/reports/financial/profit-loss?date_from=2026-11-01&date_to=2026-11-26'], ['GET', '/hr/employees'], ['POST', '/accounting/accounts', { account_code: '9001', account_name: 'x', account_type: 'asset', normal_balance: 'debit' }], ['GET', '/accounting/paid-order-events/failed'], ['GET', '/super-admin/dashboard/stats']]) { const r2 = await C.call(m, p, b); cw[`${m} ${p.split('?')[0]}`] = r2.status; }
rec('SEC-009', 'cashier role on accounting/HR/admin routes', cw);
// origin check + rate limit
const no = await mk().call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD }, { noOrigin: true }); rec('SEC-012a', 'login without Origin', { status: no.status });
const evil = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: 'http://evil.example', 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'x@example.test', password: 'x' }) }); rec('SEC-012b', 'login with foreign Origin', { status: evil.status });
const codes = []; for (let i = 0; i < 30; i++) { const rr = await mk().call('POST', '/auth/login', { email: 'ratelimit.probe@example.test', password: 'wrong' + i }); codes.push(rr.status); }
rec('SEC-012c', '30 rapid failed logins for a nonexistent email', { statusCounts: codes.reduce((a, c) => ({ ...a, [c]: (a[c] || 0) + 1 }), {}) });
// injection probes on list/search params
const inj = {}; for (const q of ["/products?search=' OR 1=1--", "/admin/orders?search=%27%3B+DROP+TABLE+orders%3B--", "/accounting/customers?search=' UNION SELECT 1--", "/suppliers?search=\\", "/accounting/journal-entries?status=posted' OR '1'='1"]) { const r3 = await A.call('GET', q); inj[q.slice(0, 48)] = { status: r3.status, body: JSON.stringify(r3.body).slice(0, 80) }; }
rec('SEC-SQLI', 'quote/injection payloads in list filters (tenant 4)', inj); rec('SEC-SQLI-DB', 'orders table intact', { orders: sql('select count(*) from orders') });
// security headers
const h = await fetch(API.replace('/api/v1', '') + '/api/v1/health'); rec('SEC-HDR', 'API security headers', { hsts: h.headers.get('strict-transport-security'), csp: !!h.headers.get('content-security-policy'), xcto: h.headers.get('x-content-type-options'), xfo: h.headers.get('x-frame-options'), server: h.headers.get('server') });
fs.writeFileSync(`${d.net}/AUD-SEC-BATTERY.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
