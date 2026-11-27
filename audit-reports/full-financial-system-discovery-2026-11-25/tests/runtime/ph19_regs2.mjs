// RUN-003 / PH-19: HC-10 (payment void/refund transitions on orders), FIN-019 (export content/CSV injection), CFR-006 (report latency on current data). Isolated tenant 4. API-DIAGNOSTIC.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-19'); let cookie = '';
const call = async (method, path, body, raw) => { const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined }); const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; '); const t = await r.text(); let j = t; if (!raw) { try { j = JSON.parse(t); } catch { j = t.slice(0, 300); } } return { status: r.status, body: j, ctype: r.headers.get('content-type'), disp: r.headers.get('content-disposition') }; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 320)); };
await call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD });
// ---- HC-10 transitions
const mk = async (qty = 1) => { const r = await call('POST', '/pos/orders', { order_source: 'pos', customer_name: 'HC10', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: 4, quantity: qty }] }); return r.body?.data?.id; };
const o = await mk(); const stat = (id) => sql(`select payment_status||'/'||status from orders where id=${id}`);
const journals = (id) => sql(`select count(*) from journal_entries where tenant_id=4 and reference_type='order' and reference_id=${id}`);
const tryPay = async (id, st) => { const r = await call('PUT', `/admin/orders/${id}/payment-status`, { payment_status: st }); return { to: st, http: r.status, resp: JSON.stringify(r.body).slice(0, 110), now: stat(id) }; };
const seq = [];
seq.push(await tryPay(o, 'paid')); await new Promise((x) => setTimeout(x, 3500)); seq.push({ journalsAfterPaid: journals(o) });
for (const st of ['refunded', 'voided', 'failed', 'pending']) { seq.push(await tryPay(o, st)); }
await new Promise((x) => setTimeout(x, 3500)); seq.push({ journalsAfterAttempts: journals(o), finalState: stat(o) });
rec('HC10-ORDER', 'paid order -> refunded/voided/failed/pending transitions', { seq });
const o2 = await mk(); const seq2 = [await tryPay(o2, 'refunded'), await tryPay(o2, 'voided'), await tryPay(o2, 'failed')]; rec('HC10-UNPAID', 'pending order -> refunded/voided/failed', { seq2 });
// ---- FIN-019 exports
await call('POST', '/accounting/accounts', { account_code: '5998', account_name: '=1+1 injected', account_type: 'expense', normal_balance: 'debit' });
await call('POST', '/accounting/accounts', { account_code: '5997', account_name: '@SUM(A1:A2)', account_type: 'expense', normal_balance: 'debit', description: '+cmd|calc!A0' });
const coa = await call('GET', '/accounting/csv/chart-of-accounts/export', null, true); const lines = String(coa.body).split('\n');
const inj = lines.filter((l) => /5998|5997/.test(l));
rec('FIN019-CSV', 'COA CSV export formula-injection cells', { status: coa.status, ctype: coa.ctype, rows: lines.length, injectedRows: inj, escaped: inj.every((l) => !/(^|,)"?[=@+\-]/.test(l.split(',').slice(1).join(','))) });
const jex = await call('GET', '/accounting/csv/journal-entries/export', null, true); rec('FIN019-JE', 'journal CSV export', { status: jex.status, ctype: jex.ctype, rows: String(jex.body).split('\n').length, jeDbCount: sql("select count(*) from journal_entries where tenant_id=4") });
const xl = await call('GET', '/accounting/excel/chart-of-accounts/export', null, true); rec('FIN019-XLSX', 'COA excel export', { status: xl.status, ctype: xl.ctype, disp: xl.disp, bytes: String(xl.body).length });
const oex = await call('GET', '/admin/orders/export', null, true); rec('FIN019-ORD', 'orders export', { status: oex.status, ctype: oex.ctype, rows: String(oex.body).split('\n').length, ordersDb: sql('select count(*) from orders where tenant_id=4') });
fs.writeFileSync(`${d.net}/FIN-019-coa-export.csv`, String(coa.body).slice(0, 20000));
// ---- CFR-006 latency (current small dataset)
const eps = ['/reports/financial/profit-loss?date_from=2026-10-28&date_to=2026-11-26', '/reports/financial/cash-flow?date_from=2026-10-28&date_to=2026-11-26', '/reports/financial/balance-sheet?as_of=2026-11-26', '/accounting/trial-balance', '/reports/sales?date_from=2026-10-28&date_to=2026-11-26', '/reports/financial/receivables-aging'];
const lat = {};
for (const e of eps) { const t = []; for (let i = 0; i < 8; i++) { const s = performance.now(); await call('GET', e); t.push(Math.round(performance.now() - s)); } t.sort((a, b) => a - b); lat[e.split('?')[0]] = { p50: t[4], max: t[7], glLines: Number(sql('select count(*) from general_ledger where tenant_id=4')) }; }
rec('CFR006', 'report latency ms on current dataset (8 runs each, cache may apply)', lat);
fs.writeFileSync(`${d.db}/AUD-REG-HC10-FIN019-CFR006.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
