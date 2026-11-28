// RUN-005 / PH-17: reproducible confirmation of OBSERVED findings F-026 (card order paid -> clearing 1113), F-035 (login rate limit and CSP headers), F-036 (unauthenticated mock Paymob routes). API-DIAGNOSTIC, isolated stack.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-17'); const O = 'http://localhost:13112';
const out = {}; await new Promise((r) => setTimeout(r, 70000)); // let the previous per-IP login rate-limit window expire
// F-036
const g = await fetch(API + '/mock/paymob/iframe', { signal: AbortSignal.timeout(20000) }); out.F036_get_iframe_noauth = { status: g.status, bytes: (await g.text()).length };
const p = await fetch(API + '/mock/paymob/process', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: '{}', signal: AbortSignal.timeout(20000) }); out.F036_post_process_noauth = { status: p.status, body: (await p.text()).slice(0, 120) };
// F-035: CSP headers on API and dashboard + rate limit
for (const [k, u] of [['api_health', API + '/health'], ['dashboard', 'http://localhost:13112/en/dashboard'], ['landing', 'http://localhost:13111/en']]) { const r = await fetch(u, { signal: AbortSignal.timeout(120000) }); out['F035_csp_' + k] = { status: r.status, csp: (r.headers.get('content-security-policy') || 'none').slice(0, 260) }; }
// F-026
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; '); const H = { Origin: O, Cookie: ck, 'Content-Type': 'application/json' };
const prod = sql("select id from products where tenant_id=4 and track_inventory=false and status='active' and price>0 order by id limit 1");
const o = await (await fetch(API + '/pos/orders', { method: 'POST', headers: H, body: JSON.stringify({ order_source: 'pos', customer_name: 'R5 card', customer_phone: '0000000000', payment_method: 'card', items: [{ product_id: Number(prod), quantity: 1 }] }) })).json(); const oid = o?.data?.id;
const before = sql(`select count(*) from journal_entries where tenant_id=4 and reference_type='order' and reference_id=${oid}`);
const paid = await fetch(`${API}/admin/orders/${oid}/payment-status`, { method: 'PUT', headers: H, body: JSON.stringify({ payment_status: 'paid' }) });
await new Promise((r) => setTimeout(r, 6000));
out.F026 = { order: oid, product: prod, order_row: sql(`select payment_method||'|'||payment_status||'|'||total_amount from orders where id=${oid}`), paid_http: paid.status, journals_before: before, lines: sql(`select string_agg(a.account_code||' '||a.account_name||' D'||l.debit_amount||' C'||l.credit_amount, ' ; ') from journal_entries j join journal_entry_lines l on l.journal_entry_id=j.id join chart_of_accounts a on a.id=l.account_id where j.tenant_id=4 and j.reference_type='order' and j.reference_id=${oid}`), settlement_routes: 'none (no route or caller references account 1113 except the posting rule; see root cause F-026)' };
const codes = []; for (let i = 0; i < 25; i++) { const r = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'nobody.r5@example.test', password: 'wrong-' + i }) }); codes.push(r.status); await r.arrayBuffer(); }
out.F035_25_bad_logins = { statuses: codes.join(','), first429_at_attempt: codes.indexOf(429) + 1 };
console.log(JSON.stringify(out, null, 1)); fs.writeFileSync(`${d.net}/AUD-R5-OBSERVED-REPRO.json`, JSON.stringify({ run: 'RUN-005', findings: ['F-026', 'F-035', 'F-036'], out }, null, 1));
