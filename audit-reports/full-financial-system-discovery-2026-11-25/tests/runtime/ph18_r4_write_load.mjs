// RUN-004 / PH-18: incremental WRITE load (10, 25, 50, 100 virtual users, 20 s each) on tenant 4 (diagnostic fixture): POS cash orders for an untracked product (3) and a tracked product (4), mixed with reads.
// Safety: aborts a stage if free host memory < 700 MB or error rate > 60%. After the run: waits for the accounting outbox to drain, then checks duplicates, stock = ledger, journals vs paid orders, trial balance. Not a production capacity claim.
import fs from 'fs';
import os from 'os';
import { execSync } from 'child_process';
import { API, PASSWORD, evidenceDirs, sql } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const H = { Origin: O, Cookie: ck, 'Content-Type': 'application/json' };
const sleep = (ms) => new Promise((x) => setTimeout(x, ms));
const freeMB = async () => { try { return +(execSync('powershell -NoProfile -Command "[math]::Round((Get-CimInstance Win32_OperatingSystem).FreePhysicalMemory/1024)"', { encoding: 'utf8' }).trim()); } catch { return 9999; } };
const dock = () => { try { return execSync('docker stats --no-stream --format "{{.Name}}={{.CPUPerc}}/{{.MemUsage}}" pos-p2a-ws05-e2e-api-1 pos-p2a-ws05-e2e-postgres-1 pos-p2a-ws05-e2e-redis-1', { encoding: 'utf8' }).trim().replace(/\n/g, ' | '); } catch { return 'n/a'; } };
const before = { orders: +sql('select count(*) from orders where tenant_id=4'), paid: +sql("select count(*) from orders where tenant_id=4 and payment_status='paid'"), stock4: +sql('select quantity_in_stock from products where id=4'), journals: +sql('select count(*) from journal_entries where tenant_id=4'), outbox: sql("select status||':'||count(*) from accounting_event_outbox where tenant_id=4 group by status order by 1").replace(/\n/g, ',') };
console.log('before', JSON.stringify(before));
const before_pending = +sql("select count(*) from accounting_event_outbox where tenant_id=4 and status in ('pending','processing','retryable','failed')");
const stages = []; let accepted4 = 0; let seq = 0;
const mk = (w) => {
  const i = seq++;
  if (i % 5 === 4) return { m: 'GET', p: ['/products?page=1&page_size=5', '/kds/stats', '/admin/orders?page=1&page_size=5', '/reports/financial/profit-loss?date_from=2026-11-01&date_to=2026-11-30'][i % 4] };
  const tracked = i % 3 === 0;
  return { m: 'POST', p: '/pos/orders', b: { order_source: 'pos', customer_name: 'LOAD', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: tracked ? 4 : 3, quantity: 1 }] }, tracked };
};
for (const vu of [10, 25, 50, 100]) {
  if ((await freeMB()) < 700) { console.log('ABORT before stage', vu, 'low memory'); stages.push({ vu, aborted: 'low memory' }); break; }
  const lat = []; const codes = {}; let n = 0, err = 0; const stop = Date.now() + 20000; let peakConn = 0; let samples = [];
  const mon = setInterval(() => { try { peakConn = Math.max(peakConn, +sql("select count(*) from pg_stat_activity where datname='pos_phase2a_e2e'")); } catch {} }, 3000);
  const dmon = setInterval(() => samples.push(dock()), 7000);
  await Promise.all(Array.from({ length: vu }, async () => {
    while (Date.now() < stop) {
      const q = mk(); const s = performance.now(); n++;
      try { const x = await fetch(API + q.p, { method: q.m, headers: H, body: q.b ? JSON.stringify(q.b) : undefined, signal: AbortSignal.timeout(20000) }); codes[`${q.m}${x.status}`] = (codes[`${q.m}${x.status}`] || 0) + 1; if (x.status >= 500) err++; if (q.tracked && x.status === 201) accepted4++; await x.arrayBuffer(); } catch { err++; codes.NETERR = (codes.NETERR || 0) + 1; }
      lat.push(performance.now() - s);
    }
  }));
  clearInterval(mon); clearInterval(dmon);
  lat.sort((a, b) => a - b); const pc = (p) => Math.round(lat[Math.floor(lat.length * p)] || 0);
  const st = { vu, requests: n, rps: +(n / 20).toFixed(1), server_errors_or_neterr: err, errorRate: +(err / n * 100).toFixed(2), p50: pc(0.5), p95: pc(0.95), p99: pc(0.99), max: Math.round(lat[lat.length - 1] || 0), codes, peak_pg_connections: peakConn, docker_samples: samples, free_mb_after: await freeMB() };
  stages.push(st); console.log(vu, 'VU', JSON.stringify({ ...st, docker_samples: undefined }));
  if (st.errorRate > 60) { console.log('ABORT: error rate'); break; }
  await sleep(8000);
}
console.log('waiting for outbox drain...');
for (let i = 0; i < 40; i++) { const pend = +sql("select count(*) from accounting_event_outbox where tenant_id=4 and status in ('pending','processing','retryable','failed')"); if (pend <= before_pending) break; await sleep(5000); }
const after = { orders: +sql('select count(*) from orders where tenant_id=4'), paid: +sql("select count(*) from orders where tenant_id=4 and payment_status='paid'"), stock4: +sql('select quantity_in_stock from products where id=4'), journals: +sql('select count(*) from journal_entries where tenant_id=4'), outbox: sql("select status||':'||count(*) from accounting_event_outbox where tenant_id=4 group by status order by 1").replace(/\n/g, ',') };
const checks = {
  new_orders: after.orders - before.orders,
  duplicate_order_numbers: sql('select count(*) from (select order_number from orders where tenant_id=4 group by 1 having count(*)>1) x'),
  tracked_orders_201: accepted4, stock_decrease: before.stock4 - after.stock4,
  stock_ne_ledger_note: 'stock decrease should equal accepted tracked orders (each qty 1)',
  new_paid: after.paid - before.paid, new_journals: after.journals - before.journals,
  paid_orders_without_sales_journal: sql("select count(*) from orders o where o.tenant_id=4 and o.payment_status='paid' and o.id>(select max(id) from orders where tenant_id=4 and id<=" + sql('select max(id) from orders where tenant_id=4') + ")-" + (after.orders - before.orders) + " and not exists (select 1 from journal_entries j where j.tenant_id=4 and j.reference_type='order' and j.reference_id=o.id)"),
  unbalanced_journals: sql('select count(*) from (select j.id from journal_entries j join journal_entry_lines l on l.journal_entry_id=j.id where j.tenant_id=4 group by j.id having round(sum(l.debit_amount)-sum(l.credit_amount),2)<>0) x'),
  tb_total_debit_minus_credit: sql('select round(sum(l.debit_amount)-sum(l.credit_amount),2) from journal_entry_lines l join journal_entries j on j.id=l.journal_entry_id where j.tenant_id=4 and j.status=\'posted\''),
};
fs.writeFileSync(`${d.net}/AUD-R4-WRITE-LOAD.json`, JSON.stringify({ run: 'RUN-004', kind: 'API-DIAGNOSTIC load', note: 'isolated stack, tenant 4 fixture, headless fetch clients on the same host as the stack; not a production capacity claim', before, stages, after, checks }, null, 1));
console.log('after', JSON.stringify(after)); console.log('checks', JSON.stringify(checks));
