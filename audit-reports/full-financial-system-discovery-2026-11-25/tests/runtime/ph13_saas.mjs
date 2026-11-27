// RUN-003 / PH-13: SaaS subscription, mock Paymob, Super Admin metrics/refund/disable (API-DIAGNOSTIC, isolated stack, mock provider only, no real credentials).
import fs from 'fs';
import { sql as sqlRaw, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-13');
const sql = (q) => { try { return sqlRaw(q); } catch (e) { return 'SQL-ERROR: ' + String(e.stderr || e).slice(0, 120); } };
const mkSession = async (path, origin, email) => { for (let i = 0; i < 30; i++) { const r = await fetch(API + path, { method: 'POST', headers: { Origin: origin, 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password: PASSWORD }) }); if (r.status === 200) { const ck = r.headers.getSetCookie().map((c) => c.split(';')[0]).join('; '); return async (m, p, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: origin, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: x.status, body: j }; }; } await new Promise((z) => setTimeout(z, 3000)); } throw new Error('login ' + path); };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 460)); };
const SA = await mkSession('/super-admin/auth/login', 'http://localhost:13114', 'superadmin@pos-saas.com');
const E = await mkSession('/auth/login', 'http://localhost:13112', 'audit.e.owner@example.test');
let r = await E('GET', '/subscription/plans'); const plans = r.body?.data ?? r.body; rec('SUB-01', 'plans', { http: r.status, plans: JSON.stringify(plans).slice(0, 380) });
r = await E('GET', '/subscription/current'); rec('SUB-02', 'current subscription of tenant E', { http: r.status, body: JSON.stringify(r.body).slice(0, 300) });
const plan = (Array.isArray(plans) ? plans[0] : plans?.plans?.[0]) ?? {};
r = await E('POST', '/subscription/initiate', { plan_id: plan.id ?? plan.plan_id ?? 1, billing_period: 'monthly' }); rec('SUB-03', 'initiate subscription payment', { http: r.status, body: JSON.stringify(r.body).slice(0, 400) });
const pay = sql("select id||'|'||status||'|'||amount||'|'||coalesce(provider_order_id,'') from subscription_payments order by id desc limit 1");
rec('SUB-04', 'latest subscription payment row', { pay });
// SaaS metrics vs SQL
const m = await SA('GET', '/super-admin/dashboard/stats'); rec('SA-01', 'dashboard stats', { http: m.status, body: JSON.stringify(m.body?.data ?? m.body).slice(0, 700), sql: sql("select (select count(*) from tenants)||' tenants; active='||(select count(*) from tenants where status='active')||'; users='||(select count(*) from users)") });
const pm = await SA('GET', '/super-admin/payment-metrics'); rec('SA-02', 'payment metrics', { http: pm.status, body: JSON.stringify(pm.body?.data ?? pm.body).slice(0, 600), sql: sql("select status||':'||count(*)||':'||coalesce(sum(amount),0) from subscription_payments group by status").split('\n') });
const tl = await SA('GET', '/super-admin/tenants'); rec('SA-03', 'tenant list', { http: tl.status, count: (tl.body?.data?.tenants ?? tl.body?.data ?? []).length, sqlTenants: sql('select count(*) from tenants') });
// disable/enable tenant 7 (unseeded probe tenant) is NOT done: protected. Use tenant 8 (diagnostic E).
r = await SA('POST', '/super-admin/tenants/8/disable', { reason: 'RUN-003 probe' }); rec('SA-04', 'disable tenant 8', { http: r.status, body: JSON.stringify(r.body).slice(0, 200), status: sql('select status from tenants where id=8') });
const dis = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: 'http://localhost:13112', 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.e.owner@example.test', password: PASSWORD }) }); rec('SA-05', 'login of owner of a disabled tenant', { http: dis.status });
const stale = await E('GET', '/products'); rec('SA-06', 'existing session of disabled tenant still works?', { http: stale.status });
r = await SA('POST', '/super-admin/tenants/8/enable', {}); rec('SA-07', 'enable tenant 8', { http: r.status, status: sql('select status from tenants where id=8') });
r = await SA('POST', '/super-admin/tenants/8/extend-trial', { days: 7 }); rec('SA-08', 'extend trial 7 days', { http: r.status, body: JSON.stringify(r.body).slice(0, 200) });
r = await SA('POST', '/super-admin/tenants/8/extend-trial', { days: -30 }); rec('SA-09', 'extend trial by -30 days (validation)', { http: r.status, body: JSON.stringify(r.body).slice(0, 200) });
r = await SA('POST', '/super-admin/tenants/8/extend-trial', { days: 100000 }); rec('SA-10', 'extend trial by 100000 days', { http: r.status, body: JSON.stringify(r.body).slice(0, 200) });
// mock paymob process + webhook
const mp = await fetch(API + '/mock/paymob/process', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ order_id: 'x', success: true }) }); rec('PAY-01', 'mock paymob process unauthenticated', { http: mp.status, body: (await mp.text()).slice(0, 200) });
const wh = await fetch(API + '/webhook/paymob', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ type: 'TRANSACTION', obj: { id: 1, success: true, order: { id: 1 }, amount_cents: 1000 } }) }); rec('PAY-02', 'webhook with unsigned payload', { http: wh.status, body: (await wh.text()).slice(0, 200) });
// refund endpoint on a non-existent / non-success payment
r = await SA('POST', '/super-admin/payments/999999/refund', { reason: 'probe' }); rec('PAY-03', 'refund unknown payment', { http: r.status, body: JSON.stringify(r.body).slice(0, 200) });
// impersonation audit
r = await SA('GET', '/super-admin/audit-logs'); rec('SA-11', 'audit logs list', { http: r.status, count: JSON.stringify(r.body).length });
fs.writeFileSync(`${d.db}/AUD-SAAS-MATRIX.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
