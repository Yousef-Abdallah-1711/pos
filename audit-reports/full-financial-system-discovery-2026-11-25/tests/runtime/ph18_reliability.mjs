// RUN-003 / PH-18: REL-004 (concurrent redrive), REL-009 (concurrent stock adjustment + orders), REL-025 (payment retry with Idempotency-Key), REL-027 (API restart while event retryable). Isolated stack; tenants 8 and 4. API-DIAGNOSTIC + docker restart of isolated api only.
import fs from 'fs';
import { execSync } from 'child_process';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112';
const sess = async (email) => { for (let i = 0; i < 30; i++) { const r = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password: PASSWORD }) }); if (r.status === 200) { const ck = r.headers.getSetCookie().map((c) => c.split(';')[0]).join('; '); return async (m, p, b, h = {}) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json', ...h }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 160); } return { status: x.status, body: j }; }; } await new Promise((z) => setTimeout(z, 3000)); } throw new Error('login failed'); };
const sleep = (ms) => new Promise((x) => setTimeout(x, ms));
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 420)); };
const E = await sess('audit.e.owner@example.test'); const A = await sess('audit.c.owner@example.test');
// ---------- REL-004 concurrent redrive (tenant 8)
let r = await E('POST', '/pos/orders', { order_source: 'pos', customer_name: 'REL4', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: 5, quantity: 1 }] }); const o = r.body?.data?.id;
await E('PUT', `/admin/orders/${o}/payment-status`, { payment_status: 'paid' }); await sleep(5000);
const jeCount = () => Number(sql(`select count(*) from journal_entries where tenant_id=8 and reference_type='order' and reference_id=${o}`));
const evId = sql(`select id from accounting_event_outbox where tenant_id=8 and order_id=${o} order by id desc limit 1`);
const before = { status: sql(`select status from accounting_event_outbox where id=${evId}`), journals: jeCount() };
sql(`update accounting_event_outbox set status='failed', attempts=max_attempts, last_error='REL-004 fixture: forced failed state' where id=${evId}`);
const rd = await Promise.all([1, 2].map(() => E('POST', `/accounting/paid-order-events/${evId}/redrive`, { reason: 'REL-004 concurrent redrive' })));
await sleep(8000);
rec('REL-004', 'two concurrent redrives of one failed event (event forced failed by SQL fixture after it had posted)', { before, http: rd.map((x) => x.status), eventNow: sql(`select status||' attempts='||attempts from accounting_event_outbox where id=${evId}`), journalsAfter: jeCount(), audits: sql(`select count(*) from accounting_event_redrive_audit where outbox_id=${evId}`) });
// ---------- REL-009 concurrent stock adjustment + order consumption (tenant 4 product 4)
const stock = () => Number(sql('select quantity_in_stock from products where id=4'));
const s0 = stock();
const reqs = [];
for (let i = 0; i < 5; i++) reqs.push(A('POST', '/inventory/adjust', { product_id: 4, quantity_change: 2, reason: 'purchase', notes: 'REL-009' }));
for (let i = 0; i < 5; i++) reqs.push(A('POST', '/pos/orders', { order_source: 'pos', customer_name: 'REL9', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: 4, quantity: 1 }] }));
const res = await Promise.all(reqs); await sleep(2000);
const okAdj = res.slice(0, 5).filter((x) => x.status === 200 || x.status === 201).length; const okOrd = res.slice(5).filter((x) => x.status === 201).length;
const ledger = Number(sql("select coalesce(sum(quantity),0) from inventory_movements where product_id=4 and movement_type in ('adjustment','reservation','reservation_release','waste','purchase','return')"));
rec('REL-009', 'concurrent 5 stock +2 adjustments and 5 one-unit orders', { s0, adjOk: okAdj, ordersOk: okOrd, expected: s0 + 2 * okAdj - okOrd, stockNow: stock(), ledgerNet: ledger, statuses: res.map((x) => x.status) });
// ---------- REL-025 payment POST retry with the same Idempotency-Key (tenant 4 customer invoice)
const inv = await A('POST', '/accounting/customer-invoices', { invoice_date: '2026-11-26', due_date: '2026-12-26', customer_id: 2, currency: 'EGP', status: 'sent', lines: [{ description: 'REL25', quantity: 1, unit_price: 1000, tax_rate: 0, discount_rate: 0, line_order: 1 }] }); const invId = inv.body?.data?.id ?? inv.body?.id;
const key = crypto.randomUUID();
const pay = () => A('POST', '/accounting/payments/received', { payment_date: '2026-11-26', customer_id: 2, invoice_id: invId, amount: 400, payment_method: 'cash', reference: 'REL25' }, { 'Idempotency-Key': key });
const p1 = await pay(); const p2 = await pay();
rec('REL-025', 'same payment POST sent twice with identical Idempotency-Key (client retry)', { http: [p1.status, p2.status], receiptsForInvoice: sql(`select count(*)||' sum='||coalesce(sum(amount),0) from payments_received where invoice_id=${invId}`), invoice: sql(`select amount_paid||'/'||balance_due from invoices_customer where id=${invId}`) });
// order creation idempotency (F-017 check)
const k2 = crypto.randomUUID(); const mkO = () => A('POST', '/pos/orders', { order_source: 'pos', customer_name: 'REL1', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: 4, quantity: 1 }] }, { 'Idempotency-Key': k2 });
const before1 = Number(sql('select count(*) from orders where tenant_id=4')); const oa = await mkO(); const ob = await mkO(); rec('REL-001', 'same order POST twice with identical Idempotency-Key', { http: [oa.status, ob.status], newOrders: Number(sql('select count(*) from orders where tenant_id=4')) - before1 });
// ---------- REL-027 restart while event retryable (tenant 8 product 5 order with COA temporarily removed is destructive; instead verify lease columns after restart on a processing event fixture)
const o3 = (await E('POST', '/pos/orders', { order_source: 'pos', customer_name: 'REL27', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: 5, quantity: 1 }] })).body?.data?.id;
await E('PUT', `/admin/orders/${o3}/payment-status`, { payment_status: 'paid' }); await sleep(4000);
const ev3 = sql(`select id from accounting_event_outbox where tenant_id=8 and order_id=${o3}`);
const jBefore = Number(sql(`select count(*) from journal_entries where tenant_id=8 and reference_type='order' and reference_id=${o3}`));
sql(`update accounting_event_outbox set status='processing', lease_expires_at=now()+interval '20 seconds', completed_at=null where id=${ev3}`); // fixture: simulates a worker that died mid-processing (journals already exist)
execSync('docker restart pos-p2a-ws05-e2e-api-1', { stdio: 'ignore' });
for (let i = 0; i < 40; i++) { try { const h = await fetch(API + '/health'); if (h.status === 200) break; } catch { /* wait */ } await sleep(2000); }
await sleep(45000);
rec('REL-027', 'event left in processing with a 20s lease, API restarted (journals already exist)', { eventNow: sql(`select status||' attempts='||attempts||' gen='||lease_generation from accounting_event_outbox where id=${ev3}`), journalsBefore: jBefore, journalsAfter: Number(sql(`select count(*) from journal_entries where tenant_id=8 and reference_type='order' and reference_id=${o3}`)) });
fs.writeFileSync(`${d.db}/AUD-REL-BATTERY.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
