// RUN-003 / PH-14: compressed operating-month simulation on a FRESH diagnostic tenant (Audit F). API-DIAGNOSTIC with an INDEPENDENT integer-cents ledger.
// Orders cannot be back-dated by the product, so sales happen "today"; all other events carry real dates across 2026-11-01..11-30.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-14'); const O = 'http://localhost:13112'; const EMAIL = 'audit.f.owner@example.test';
const sleep = (ms) => new Promise((x) => setTimeout(x, ms));
const sess = async (email, tries = 40) => { for (let i = 0; i < tries; i++) { const r = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email, password: PASSWORD }) }); if (r.status === 200) { const ck = r.headers.getSetCookie().map((c) => c.split(';')[0]).join('; '); return async (m, p, b, form) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, ...(form ? {} : { 'Content-Type': 'application/json' }) }, body: form ? b : b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: x.status, body: j }; }; } if (r.status === 401 && i > 2) break; await sleep(3000); } return null; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 300)); };
// --- tenant bootstrap (fresh registration + diagnostic COA copy, same technique as tenants 4/8)
let r0 = await fetch(API + '/auth/register', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ name: 'Audit Owner F', email: EMAIL, password: PASSWORD, restaurant_name: 'Audit F Month', restaurant_slug: 'audit-f-month', phone: '+201000000006', billing_period: 'monthly' }) });
rec('REG', 'register tenant F', { http: r0.status });
const tid = Number(sql(`select tenant_id from users where email='${EMAIL}'`));
if (sql(`select count(*) from chart_of_accounts where tenant_id=${tid}`) === '0') {
  const COLS = sql("select string_agg(column_name,',' order by ordinal_position) from information_schema.columns where table_name='chart_of_accounts' and column_name not in ('id','tenant_id','restaurant_id','parent_id','created_at','updated_at','created_by','updated_by','current_balance')");
  sql(`insert into chart_of_accounts (tenant_id,restaurant_id,parent_id,${COLS}) select ${tid},${tid},null,${COLS} from chart_of_accounts where tenant_id=4`);
}
const A = await sess(EMAIL);
const acct = (code) => Number(sql(`select id from chart_of_accounts where tenant_id=${tid} and account_code='${code}'`));
const L = {}; const post = (code, dr, cr) => { L[code] = (L[code] || 0) + dr - cr; }; // independent ledger: net debit in cents
const jid = (r) => r.body?.data?.id ?? r.body?.id;
const journal = async (date, desc, lines, expect = true) => { const r = await A('POST', '/accounting/journal-entries', { entry_date: date, description: desc, entry_type: 'manual', lines: lines.map(([c, dr, cr], i) => ({ account_id: acct(c), debit_amount: dr / 100, credit_amount: cr / 100, line_order: i + 1 })) }); const id = jid(r); const p = await A('POST', `/accounting/journal-entries/${id}/post`, {}); if (expect && p.status === 200) lines.forEach(([c, dr, cr]) => post(c, dr, cr)); return p.status; };
// Day 1 capital
rec('D01', 'capital contribution 100,000', { http: await journal('2026-11-01', 'Capital contribution', [['1110', 10000000, 0], ['3100', 0, 10000000]]) });
// Supplier + invoices (no discount lines: oracle unambiguous), payments
let r = await A('POST', '/suppliers', { company_name: 'Month Supplier', email: 'ms@example.test', phone: '+201000000100' }); const sid = r.body?.data?.id ?? r.body?.id;
const supInv = async (date, qty, unit, taxPct) => { const rr = await A('POST', '/accounting/supplier-invoices', { supplier_ref: 'M' + date, invoice_date: date, due_date: '2026-12-30', supplier_id: sid, currency: 'EGP', status: 'draft', lines: [{ description: 'Goods', quantity: qty, unit_price: unit, tax_rate: taxPct, discount_rate: 0, line_order: 1 }] }); const id = jid(rr); const ap = await A('POST', `/accounting/supplier-invoices/${id}/approve`, {}); const net = Math.round(qty * unit * 100); const tax = Math.round(net * taxPct / 100); if (ap.status === 200) { post('5110', net, 0); post('1150', tax, 0); post('2101', 0, net + tax); } return { id, net, tax, st: ap.status }; };
const s1 = await supInv('2026-11-03', 50, 200, 14); const s2 = await supInv('2026-11-10', 30, 150, 14); const s3 = await supInv('2026-11-17', 20, 310.5, 14);
rec('SUP', 'three supplier invoices', { s1, s2, s3 });
const supPay = async (invId, amtCents, date) => { const rr = await A('POST', '/accounting/payments/made', { payment_date: date, supplier_id: sid, invoice_id: invId, amount: amtCents / 100, payment_method: 'cash', reference: 'month' }); if (rr.status === 201) { post('2101', amtCents, 0); post('1110', 0, amtCents); } return rr.status; };
rec('SUPPAY', 'payments', { a: await supPay(s1.id, 5700000 / 10, '2026-11-12'), b: await supPay(s2.id, Math.round((s2.net + s2.tax) * 1), '2026-11-20'), c: await supPay(s3.id, 500000, '2026-11-25') });
// Customer invoices + receipts
r = await A('POST', '/accounting/customers', { customer_name: 'Month Customer', email: 'mc@example.test', phone: '+201000000101' }); const cid = r.body?.data?.id ?? r.body?.id;
const custInv = async (date, amt, taxPct) => { const rr = await A('POST', '/accounting/customer-invoices', { invoice_date: date, due_date: '2026-12-30', customer_id: cid, currency: 'EGP', status: 'sent', lines: [{ description: 'Catering', quantity: 1, unit_price: amt, tax_rate: taxPct, discount_rate: 0, line_order: 1 }] }); const id = jid(rr); const net = Math.round(amt * 100); const tax = Math.round(net * taxPct / 100); if (rr.status === 201) { post('1120', net + tax, 0); post('4100', 0, net); post('2120', 0, tax); } return { id, net, tax, st: rr.status }; };
const c1 = await custInv('2026-11-05', 4000, 14); const c2 = await custInv('2026-11-15', 2500, 14);
const rcpt = async (invId, cents, date) => { const rr = await A('POST', '/accounting/payments/received', { payment_date: date, customer_id: cid, invoice_id: invId, amount: cents / 100, payment_method: 'cash', reference: 'month' }); if (rr.status === 201) { post('1110', cents, 0); post('1120', 0, cents); } return rr.status; };
rec('CUST', 'invoices and receipts', { c1, c2, r1: await rcpt(c1.id, 300000, '2026-11-12'), r2: await rcpt(c1.id, 156000, '2026-11-26'), r3: await rcpt(c2.id, 100000, '2026-11-28') });
// Expenses through the month
for (const [date, code, amt] of [['2026-11-02', '5310', 1200000], ['2026-11-08', '5310', 150000], ['2026-11-14', '5310', 90000], ['2026-11-22', '5310', 60000]]) rec('EXP-' + date, 'expense ' + amt, { http: await journal(date, 'Expense ' + date, [[code, amt, 0], ['1110', 0, amt]]) });
// Owner withdrawal
r = await A('POST', '/accounting/accounts', { account_code: '3150', account_name: 'Owner Drawings', account_type: 'equity', normal_balance: 'debit', currency: 'EGP' });
rec('DRAW', 'owner withdrawal 3,000', { http: await journal('2026-11-27', 'Owner withdrawal', [['3150', 300000, 0], ['1110', 0, 300000]]) });
// Sales today: 10 POS cash orders of qty 1-4 of a tracked product price 100 (tax 10%), COGS cost 40
const fd = new FormData(); r = await A('POST', '/categories', { name: 'Month Mains', language: 'en' }); const cat = r.body?.id ?? r.body?.data?.id; fd.append('data', JSON.stringify({ category_id: cat, name_en: 'Month Item', price: 100, cost: 40, quantity_in_stock: 500, track_inventory: true, is_available: true }));
r = await A('POST', '/products', fd, true); const pid = r.body?.id;
let units = 0; const qtys = [1, 2, 3, 4, 1, 2, 3, 4, 2, 3];
for (const q of qtys) { const o = await A('POST', '/pos/orders', { order_source: 'pos', customer_name: 'W', customer_phone: '0000000000', payment_method: 'cash', items: [{ product_id: pid, quantity: q }] }); const oid = o.body?.data?.id; const pr = await A('PUT', `/admin/orders/${oid}/payment-status`, { payment_status: 'paid' }); if (pr.status === 200) { units += q; } }
await sleep(15000);
const net = units * 10000, tax = net / 10, cogs = units * 4000;
post('1110', net + tax, 0); post('4100', 0, net); post('2120', 0, tax); post('5110', cogs, 0); post('1130', 0, cogs);
rec('SALES', 'sales (10 cash orders)', { units, net, tax, cogs });
// payroll
r = await A('POST', '/hr/salaries', { employee_id: 0 }); // employees cannot be created through UI (F-014); payroll for tenant F not executable
rec('PAYROLL', 'payroll needs an employee; employee creation is blocked by F-014 so payroll is NOT part of this month', { http: r.status });
// ---- reconcile independent ledger vs system
const sysBal = {}; for (const row of sql(`select a.account_code||'|'||coalesce(sum(g.debit_amount-g.credit_amount),0) from general_ledger g join chart_of_accounts a on a.id=g.account_id where g.tenant_id=${tid} group by a.account_code`).split('\n').filter(Boolean)) { const [c, v] = row.split('|'); sysBal[c] = Math.round(Number(v) * 100); }
const diffs = {}; for (const c of new Set([...Object.keys(L), ...Object.keys(sysBal)])) { const dlt = (sysBal[c] || 0) - (L[c] || 0); if (dlt !== 0) diffs[c] = dlt / 100; }
const tb = await A('GET', '/accounting/trial-balance'); const pnl = await A('GET', '/reports/financial/profit-loss?date_from=2026-11-01&date_to=2026-11-30'); const bs = await A('GET', '/reports/financial/balance-sheet?as_of=2026-11-30'); const cf = await A('GET', '/reports/financial/cash-flow?date_from=2026-11-01&date_to=2026-11-30');
const rev = -(L['4100'] || 0) / 100, exp = ((L['5110'] || 0) + (L['5310'] || 0)) / 100; const cashOracle = ((L['1110'] || 0)) / 100;
rec('RECON', 'independent ledger vs system GL', { oracleCash: cashOracle, sysBal, oracle: L, diffs });
rec('REPORTS', 'reports vs oracle', { tbBalanced: tb.body?.data?.is_balanced, pnl: JSON.stringify(pnl.body?.data).slice(0, 360), oracleRevenue: rev, oracleExpenses: exp, oracleNetProfit: rev - exp, bsDifference: JSON.stringify(bs.body?.data).slice(-300), cashFlowReported: cf.body?.data?.net_cash_flow, oracleCashChange: cashOracle });
fs.writeFileSync(`${d.db}/AUD-MONTH.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', tenant: tid, log }, null, 1));
