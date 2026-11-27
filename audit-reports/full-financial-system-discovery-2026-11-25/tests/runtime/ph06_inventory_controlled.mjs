// RUN-003 / PH-06 / OBS-INV-001: controlled fresh-product stock-vs-ledger experiment (API-level diagnostic, isolated stack, tenant 4).
// Not a UI test: labelled API-DIAGNOSTIC. Product code untouched.
import fs from 'fs';
import { sql, ROOT, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-06');
const O = 'http://localhost:13112';
let cookie = '';
const call = async (method, path, body, extra = {}) => {
  const r = await fetch(API + path, { method, headers: { Origin: O, Cookie: cookie, ...(body && !(body instanceof FormData) ? { 'Content-Type': 'application/json' } : {}), ...extra }, body: body instanceof FormData ? body : body ? JSON.stringify(body) : undefined });
  const sc = r.headers.getSetCookie?.() ?? [];
  if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; ');
  const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t; }
  return { status: r.status, body: j };
};
const log = [];
const snap = (label, pid) => {
  const stock = sql(`select quantity_in_stock from products where id=${pid}`);
  const mv = sql(`select movement_type||':'||sum(quantity) from inventory_movements where product_id=${pid} group by movement_type order by 1`).split('\n').join(' ');
  const net = sql(`select coalesce(sum(quantity),0) from inventory_movements where product_id=${pid} and movement_type in ('adjustment','reservation','reservation_release','waste','purchase','return')`);
  log.push({ step: label, stock, movementsByType: mv, ledgerNetExcludingSaleRows: net });
  console.log(label, '| stock', stock, '| movements', mv, '| ledger(excl sale rows)', net);
};
let r = await call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD });
const fd = new FormData(); fd.append('data', JSON.stringify({ category_id: 1, name_en: 'Audit Inv Probe B', track_inventory: true, price: 100, cost: 40, quantity_in_stock: 20, is_available: true }));
r = await call('POST', '/products', fd); log.push({ step: 'create', status: r.status, body: JSON.stringify(r.body).slice(0, 300) });
console.log('create', r.status, JSON.stringify(r.body).slice(0, 200));
const pid = sql("select id from products where tenant_id=4 and name_en='Audit Inv Probe B' order by id desc limit 1") || sql("select id from products where tenant_id=4 order by id desc limit 1");
snap('S0 after create(20)', pid);
r = await call('POST', '/pos/orders', { order_source: 'pos', customer_name: 'Probe', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: Number(pid), quantity: 3 }] });
const o1 = r.body?.data?.id; log.push({ step: 'order qty3', status: r.status, id: o1 }); snap('S1 after order qty3 (unpaid)', pid);
r = await call('PUT', `/admin/orders/${o1}/payment-status`, { payment_status: 'paid' }); log.push({ step: 'pay o1', status: r.status }); await new Promise((x) => setTimeout(x, 4000)); snap('S2 after paying o1', pid);
r = await call('POST', '/pos/orders', { order_source: 'pos', customer_name: 'Probe', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: Number(pid), quantity: 2 }] });
const o2 = r.body?.data?.id; snap('S3 after order qty2 (unpaid)', pid);
r = await call('DELETE', `/admin/orders/${o2}`); log.push({ step: 'cancel o2', status: r.status, body: JSON.stringify(r.body).slice(0, 200) }); snap('S4 after cancelling o2', pid);
r = await call('POST', '/inventory/adjust', { product_id: Number(pid), quantity_change: -1, reason: 'waste', notes: 'probe waste' }); log.push({ step: 'waste -1', status: r.status }); snap('S5 after waste -1', pid);
r = await call('POST', '/inventory/adjust', { product_id: Number(pid), quantity_change: 5, reason: 'purchase', notes: 'probe purchase' }); log.push({ step: 'purchase +5', status: r.status }); snap('S6 after purchase +5', pid);
const je = sql(`select count(*) from journal_entries j where j.tenant_id=4 and (j.description ilike '%probe%' or j.reference_type='inventory_adjustment')`);
const jemv = sql(`select j.entry_number||' '||j.reference_type||' '||j.created_at::time(0) from journal_entries j where j.tenant_id=4 order by j.id desc limit 6`).split('\n');
fs.writeFileSync(`${d.db}/OBS-INV-001-controlled.json`, JSON.stringify({ run: 'RUN-003', test: 'OBS-INV-001-controlled', product_id: pid, orders: [o1, o2], log, journalsMentioningInventoryAdjust: je, latestJournals: jemv }, null, 1));
