// RUN-003 / PH-06: supplier invoice discount/tax, approval, partial payments, overpayment guard, concurrent payments, void-with-payments. API-DIAGNOSTIC against isolated tenant 4.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-06');
let cookie = '';
const call = async (method, path, body) => {
  const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined });
  const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; ');
  const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t; } return { status: r.status, body: j };
};
const log = []; const ap = () => sql("select coalesce(sum(l.credit_amount-l.debit_amount),0) from journal_entry_lines l join journal_entries j on j.id=l.journal_entry_id join chart_of_accounts a on a.id=l.account_id where j.tenant_id=4 and j.status='posted' and a.account_code='2101'");
const rec = (step, r, extra = {}) => { const e = { step, status: r?.status, resp: JSON.stringify(r?.body).slice(0, 420), apCreditBalance: ap(), ...extra }; log.push(e); console.log(step, r?.status, 'AP', e.apCreditBalance, JSON.stringify(r?.body).slice(0, 160)); };
await call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD });
let r = await call('POST', '/suppliers', { company_name: 'Audit Supplier B', email: 'supb@example.test', phone: '+201000000090' }); rec('create supplier', r);
const sid = r.body?.id ?? r.body?.data?.id ?? Number(sql("select id from suppliers where tenant_id=4 and company_name='Audit Supplier B'"));
r = await call('POST', '/accounting/supplier-invoices', { supplier_ref: 'SUPB-001', invoice_date: '2026-11-26', due_date: '2026-12-26', supplier_id: sid, currency: 'EGP', status: 'draft', lines: [{ description: 'Probe goods', quantity: 10, unit_price: 1000, tax_rate: 14, discount_rate: 10, line_order: 1 }] }); rec('create invoice draft (10x1000, disc 10%, tax 14%)', r);
const inv = r.body?.id ?? r.body?.data?.id; const sub = r.body?.subtotal ?? r.body?.data?.subtotal; const tot = r.body?.total_amount ?? r.body?.data?.total_amount;
log.push({ oracle: 'gross 10000 - disc 1000 = 9000; VAT 14% = 1260; total 10260 (discount-first); pre-discount VAT = 1400 total 10400', apiSubtotal: sub, apiTotal: tot });
r = await call('POST', `/accounting/supplier-invoices/${inv}/approve`, {}); rec('approve', r);
const total = Number(sql(`select total_amount from invoices_supplier where id=${inv}`));
const pay = (amt) => call('POST', '/accounting/payments/made', { payment_date: '2026-11-26', supplier_id: sid, invoice_id: inv, amount: amt, payment_method: 'cash', reference: 'probe' });
rec('pay 4000', await pay(4000)); rec('pay 4000', await pay(4000));
const bal = Number(sql(`select balance_due from invoices_supplier where id=${inv}`)); log.push({ balanceAfter2: bal, total });
rec(`overpay attempt ${bal + 1000} (balance ${bal})`, await pay(bal + 1000));
const par = await Promise.all([1, 2, 3, 4].map(() => pay(Math.floor(bal / 3)))); log.push({ concurrent: par.map((x) => x.status) }); rec('after 4 concurrent payments of floor(balance/3)', par[0], { statuses: par.map((x) => x.status) });
log.push({ final: sql(`select status||' paid='||amount_paid||' bal='||balance_due||' total='||total_amount from invoices_supplier where id=${inv}`) });
rec('void invoice with payments', await call('POST', `/accounting/supplier-invoices/${inv}/void`, { reason: 'probe' }));
log.push({ invoiceAfterVoid: sql(`select status||' paid='||amount_paid||' bal='||balance_due from invoices_supplier where id=${inv}`), jeStates: sql(`select string_agg(entry_number||':'||reference_type||':'||status,', ') from journal_entries where tenant_id=4 and reference_id=${inv} and reference_type like '%supplier%'`) });
fs.writeFileSync(`${d.db}/AUD-SUP-FLOW.json`, JSON.stringify({ run: 'RUN-003', supplier_id: sid, invoice_id: inv, log }, null, 1));
