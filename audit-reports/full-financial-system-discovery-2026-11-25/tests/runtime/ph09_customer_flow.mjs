// RUN-003 / PH-09: customer invoice/receipt matrix on a FRESH invoice. API-DIAGNOSTIC, isolated tenant 4.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-09');
let cookie = '';
const call = async (method, path, body) => {
  const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined });
  const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; ');
  const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t; } return { status: r.status, body: j };
};
const log = []; const ar = () => sql("select coalesce(sum(l.debit_amount-l.credit_amount),0) from journal_entry_lines l join journal_entries j on j.id=l.journal_entry_id join chart_of_accounts a on a.id=l.account_id where j.tenant_id=4 and j.status='posted' and a.account_code='1120'");
const rec = (step, r, extra = {}) => { const e = { step, status: r?.status, resp: JSON.stringify(r?.body).slice(0, 300), arDebitBalance: ar(), ...extra }; log.push(e); console.log(step, r?.status, 'AR', e.arDebitBalance, JSON.stringify(r?.body).slice(0, 140)); };
await call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD });
let r = await call('POST', '/accounting/customers', { customer_name: 'Audit Customer B', email: 'custb@example.test', phone: '+201000000091' }); rec('create customer', r);
const cid = r.body?.data?.id ?? r.body?.id ?? Number(sql("select max(id) from customers_accounting where tenant_id=4"));
r = await call('POST', '/accounting/customer-invoices', { invoice_date: '2026-11-26', due_date: '2026-12-26', customer_id: cid, currency: 'EGP', status: 'sent', lines: [{ description: 'Probe', quantity: 1, unit_price: 3000, tax_rate: 0, discount_rate: 0, line_order: 1 }] }); rec('create invoice sent 3000 no tax', r);
const inv = r.body?.data?.id ?? r.body?.id;
const pay = (amt) => call('POST', '/accounting/payments/received', { payment_date: '2026-11-26', customer_id: cid, invoice_id: inv, amount: amt, payment_method: 'cash', reference: 'probe' });
rec('single overpayment 5000 on 3000', await pay(5000));
rec('partial receipt 1000', await pay(1000));
const par = await Promise.all([1, 2, 3, 4, 5].map(() => pay(1000))); log.push({ concurrent5x1000_on_balance2000: par.map((x) => x.status) }); rec('after 5 concurrent receipts of 1000 on balance 2000', par[0], { statuses: par.map((x) => x.status) });
log.push({ final: sql(`select status||' paid='||amount_paid||' bal='||balance_due||' total='||total_amount from invoices_customer where id=${inv}`), receiptsSum: sql(`select coalesce(sum(amount),0) from payments_received where invoice_id=${inv}`) });
rec('void paid invoice', await call('POST', `/accounting/customer-invoices/${inv}/void`, { reason: 'probe' }));
log.push({ afterVoid: sql(`select status||' paid='||amount_paid||' bal='||balance_due from invoices_customer where id=${inv}`) });
fs.writeFileSync(`${d.db}/AUD-CINV-FLOW.json`, JSON.stringify({ run: 'RUN-003', customer_id: cid, invoice_id: inv, log }, null, 1));
