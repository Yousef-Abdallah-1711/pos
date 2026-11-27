// RUN-003 / PH-19 / CFR-004: durable accounting outbox - paid order with NO chart of accounts -> failure/retry state -> API restart -> seed COA (diagnostic tenant) -> redrive -> posting.
// Uses a NEW diagnostic tenant (Audit E). Tenants 5/6/7 stay untouched. API-DIAGNOSTIC + docker restart of the ISOLATED api container only.
import fs from 'fs';
import { execSync } from 'child_process';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-19'); let cookie = '';
const call = async (method, path, body, isForm) => { const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, ...(isForm ? {} : { 'Content-Type': 'application/json' }) }, body: isForm ? body : body ? JSON.stringify(body) : undefined }); const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; '); const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: r.status, body: j }; };
const sleep = (ms) => new Promise((x) => setTimeout(x, ms));
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 300)); };
const EMAIL = 'audit.e.owner@example.test';
let r = await call('POST', '/auth/register', { name: 'Audit Owner E', email: EMAIL, password: PASSWORD, restaurant_name: 'Audit E Outbox', restaurant_slug: 'audit-e-outbox', phone: '+201000000005', billing_period: 'monthly' }); rec('REG', 'register diagnostic tenant E', { status: r.status });
r = await call('POST', '/auth/login', { email: EMAIL, password: PASSWORD }); rec('LOGIN', 'login', { status: r.status });
const tid = Number(sql(`select tenant_id from users where email='${EMAIL}'`)); rec('TID', 'tenant id', { tid });
r = await call('POST', '/categories', { name: 'Probe', language: 'en' }); const cat = r.body?.id ?? r.body?.data?.id; rec('CAT', 'category', { status: r.status, cat });
const fd = new FormData(); fd.append('data', JSON.stringify({ category_id: cat, name_en: 'Outbox Probe Item', price: 100, cost: 40, quantity_in_stock: 50, track_inventory: true, is_available: true }));
r = await call('POST', '/products', fd, true); const pid = r.body?.id; rec('PROD', 'tracked product', { status: r.status, pid });
r = await call('POST', '/pos/orders', { order_source: 'pos', customer_name: 'Probe', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: pid, quantity: 2 }] }); const oid = r.body?.data?.id; rec('ORDER', 'POS order', { status: r.status, oid });
r = await call('PUT', `/admin/orders/${oid}/payment-status`, { payment_status: 'paid' }); rec('PAID', 'mark paid with NO COA', { status: r.status, resp: JSON.stringify(r.body).slice(0, 120) });
await sleep(6000);
const ob = () => sql(`select id||' status='||status||' attempts='||attempts||'/'||max_attempts||' err='||left(coalesce(last_error,''),80) from accounting_event_outbox where tenant_id=${tid} order by id desc limit 1`);
const jc = () => sql(`select count(*) from journal_entries where tenant_id=${tid}`);
rec('S1', 'outbox after payment (no COA)', { outbox: ob(), journals: jc(), orderPaid: sql(`select payment_status from orders where id=${oid}`) });
execSync('docker restart pos-p2a-ws05-e2e-api-1', { stdio: 'ignore' });
for (let i = 0; i < 40; i++) { try { const h = await fetch(API + '/health'); if (h.status === 200) break; } catch { /* waiting */ } await sleep(2000); }
rec('S2', 'outbox after API restart (durability of pending event)', { outbox: ob(), journals: jc() });
// seed COA + rules for tenant E by copying tenant 1 (diagnostic fixture, same technique as tenant C)
const restId = sql(`select id from restaurants where tenant_id=${tid} limit 1`);
rec('S3', 'restaurant id', { restId });
fs.writeFileSync(`${d.db}/CFR-004-state.json`, JSON.stringify({ run: 'RUN-003', tenant: tid, order: oid, restaurant: restId, log }, null, 1));
