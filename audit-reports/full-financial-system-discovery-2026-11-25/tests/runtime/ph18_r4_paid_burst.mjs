// RUN-004 / PH-18: concurrent payment burst. Marks 600 of the LOAD orders (tenant 4, untracked product 3, 100 + 10% tax) paid at 50 concurrent clients, then checks outbox drain, one sales journal per order, balanced ledger and revenue/AR totals. API-DIAGNOSTIC.
import fs from 'fs';
import { API, PASSWORD, evidenceDirs, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const H = { Origin: O, Cookie: ck, 'Content-Type': 'application/json' };
const sleep = (ms) => new Promise((x) => setTimeout(x, ms));
const N = 600;
const ids = sql(`select id from orders where tenant_id=4 and customer_name='LOAD' and payment_status<>'paid' and status<>'cancelled' order by id limit ${N}`).split('\n').filter(Boolean).map(Number);
const sums = sql(`select count(*), round(sum(total_amount),2), round(sum(subtotal),2), round(sum(tax_amount),2) from orders where id in (${ids.join(',')})`);
const before = { journals: +sql('select count(*) from journal_entries where tenant_id=4'), outbox: sql("select status||':'||count(*) from accounting_event_outbox where tenant_id=4 group by status order by 1").replace(/\n/g, ','), gl_rev: sql("select round(coalesce(sum(l.credit_amount-l.debit_amount),0),2) from journal_entry_lines l join journal_entries j on j.id=l.journal_entry_id join chart_of_accounts a on a.id=l.account_id where j.tenant_id=4 and j.status='posted' and a.account_type='revenue'") };
const codes = {}; const lat = []; let next = 0; const t0 = Date.now();
await Promise.all(Array.from({ length: 50 }, async () => { while (next < ids.length) { const id = ids[next++]; const s = performance.now(); try { const x = await fetch(`${API}/admin/orders/${id}/payment-status`, { method: 'PUT', headers: H, body: JSON.stringify({ payment_status: 'paid' }), signal: AbortSignal.timeout(30000) }); codes[x.status] = (codes[x.status] || 0) + 1; await x.arrayBuffer(); } catch { codes.NETERR = (codes.NETERR || 0) + 1; } lat.push(performance.now() - s); } }));
const burstMs = Date.now() - t0; lat.sort((a, b) => a - b);
const pend = () => +sql("select count(*) from accounting_event_outbox where tenant_id=4 and status in ('pending','processing') ");
let drained = false; const t1 = Date.now();
for (let i = 0; i < 60; i++) { if (pend() === 0 && +sql(`select count(*) from journal_entries where tenant_id=4 and reference_type='order' and reference_id in (${ids.join(',')})`) >= ids.length) { drained = true; break; } await sleep(4000); }
const drainMs = Date.now() - t1;
const idl = ids.join(',');
const after = {
  paid_now: sql(`select count(*) from orders where id in (${idl}) and payment_status='paid'`),
  order_journals: sql(`select count(*) from journal_entries where tenant_id=4 and reference_type='order' and reference_id in (${idl})`),
  orders_with_2plus_journals: sql(`select count(*) from (select reference_id from journal_entries where tenant_id=4 and reference_type='order' and reference_id in (${idl}) group by 1 having count(*)>1) x`),
  orders_without_journal: sql(`select count(*) from orders o where o.id in (${idl}) and not exists (select 1 from journal_entries j where j.tenant_id=4 and j.reference_type='order' and j.reference_id=o.id)`),
  outbox: sql("select status||':'||count(*) from accounting_event_outbox where tenant_id=4 group by status order by 1").replace(/\n/g, ','),
  unbalanced: sql('select count(*) from (select j.id from journal_entries j join journal_entry_lines l on l.journal_entry_id=j.id where j.tenant_id=4 group by j.id having round(sum(l.debit_amount)-sum(l.credit_amount),2)<>0) x'),
  gl_rev: sql("select round(coalesce(sum(l.credit_amount-l.debit_amount),0),2) from journal_entry_lines l join journal_entries j on j.id=l.journal_entry_id join chart_of_accounts a on a.id=l.account_id where j.tenant_id=4 and j.status='posted' and a.account_type='revenue'"),
  cash_debit_for_these: sql(`select round(coalesce(sum(l.debit_amount),0),2) from journal_entry_lines l join journal_entries j on j.id=l.journal_entry_id where j.tenant_id=4 and j.reference_type='order' and j.reference_id in (${idl}) and l.debit_amount>0`),
};
const out = { run: 'RUN-004', kind: 'API-DIAGNOSTIC', orders: ids.length, order_sums_count_total_subtotal_tax: sums, expected_revenue_delta: sql(`select round(sum(subtotal),2) from orders where id in (${idl})`), before, burst: { clients: 50, ms: burstMs, rps: +(ids.length / (burstMs / 1000)).toFixed(1), codes, p50: Math.round(lat[Math.floor(lat.length * .5)]), p95: Math.round(lat[Math.floor(lat.length * .95)]), max: Math.round(lat[lat.length - 1]) }, drained, drainMs, after };
fs.writeFileSync(`${d.net}/AUD-R4-PAID-BURST.json`, JSON.stringify(out, null, 1)); console.log(JSON.stringify(out, null, 1));
