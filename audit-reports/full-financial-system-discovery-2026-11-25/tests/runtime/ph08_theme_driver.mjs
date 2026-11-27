// RUN-003 / PH-08: theme builder API (create/duplicate/activate/history/injection) and driver/delivery assignment lifecycle on tenant 4. API-DIAGNOSTIC.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-08'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: x.status, body: j }; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 420)); };
let r = await A('GET', '/admin/themes'); rec('THM-01', 'list themes', { http: r.status, body: JSON.stringify(r.body).slice(0, 260) });
r = await A('GET', '/admin/theme-presets'); rec('THM-02', 'presets', { http: r.status, n: JSON.stringify(r.body).length });
r = await A('POST', '/admin/themes', { name: 'RUN3 Theme', primary_color: 'not-a-color' }); rec('THM-03', 'theme with invalid color', { http: r.status, resp: JSON.stringify(r.body).slice(0, 220) });
r = await A('POST', '/admin/themes', { name: 'RUN3 <img src=x onerror=alert(1)>', custom_css: 'body{background:url(javascript:alert(1))}' }); const tid = r.body?.data?.id ?? r.body?.id; rec('THM-04', 'theme name HTML + custom_css javascript: url', { http: r.status, id: tid, resp: JSON.stringify(r.body).slice(0, 220) });
if (tid) { r = await A('POST', `/admin/themes/${tid}/duplicate`, {}); rec('THM-05', 'duplicate', { http: r.status }); r = await A('POST', `/admin/themes/${tid}/activate`, {}); rec('THM-06', 'activate', { http: r.status }); r = await A('GET', `/admin/themes/${tid}/history`); rec('THM-07', 'history', { http: r.status, n: JSON.stringify(r.body).length }); r = await A('DELETE', `/admin/themes/${tid}`); rec('THM-08', 'delete ACTIVE theme', { http: r.status, resp: JSON.stringify(r.body).slice(0, 160) }); }
// drivers / delivery
r = await A('GET', '/admin/drivers'); const drv = r.body?.drivers?.[0]; rec('DRV-01', 'drivers list', { http: r.status, n: r.body?.drivers?.length, first: drv && `${drv.id}:${drv.status}` });
const o = await A('POST', '/pos/orders', { order_source: 'pos', customer_name: 'DEL', customer_phone: '0000000000', payment_method: 'cash', delivery_address: '9 Probe St', items: [{ product_id: 4, quantity: 1 }] }); const oid = o.body?.data?.id;
r = await A('POST', `/admin/orders/${oid}/assign`, { driver_id: drv?.id }); rec('DRV-02', 'assign driver to order', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200), orderStatus: sql(`select status from orders where id=${oid}`), driverStatus: sql(`select status from drivers where id=${drv?.id}`) });
r = await A('POST', `/admin/orders/${oid}/assign`, { driver_id: drv?.id }); rec('DRV-03', 'assign same driver twice / busy driver', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200) });
for (const st of ['preparing', 'ready', 'out_for_delivery', 'delivered']) { r = await A('PUT', `/admin/orders/${oid}/status`, { status: st }); rec('DRV-04-' + st, 'status -> ' + st, { http: r.status, resp: JSON.stringify(r.body).slice(0, 120), now: sql(`select status||'/'||payment_status from orders where id=${oid}`) }); }
rec('DRV-05', 'COD: delivered order payment status and journals (cash on delivery settlement)', { pay: sql(`select payment_status from orders where id=${oid}`), journals: sql(`select count(*) from journal_entries where tenant_id=4 and reference_type='order' and reference_id=${oid}`) });
r = await A('POST', '/admin/drivers', { first_name: 'Bad', last_name: 'Phone', email: 'baddrv@example.test', phone_number: 'abc', license_number: 'L1', vehicle_type: 'bike' }); rec('DRV-06', 'driver with non-numeric phone', { http: r.status, resp: JSON.stringify(r.body).slice(0, 200) });
fs.writeFileSync(`${d.db}/AUD-THEME-DRIVER.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
