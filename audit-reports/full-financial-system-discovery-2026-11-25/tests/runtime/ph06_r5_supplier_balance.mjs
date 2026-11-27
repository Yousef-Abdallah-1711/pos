// RUN-005 / PH-06 (F-021 root-cause confirmation): replays supplier transactions on a NEW supplier and records suppliers.current_balance, the invoice balance_due and the GL AP credit balance (2101) after every step,
// then edits the supplier through the API with an unrelated change to show whether current_balance is overwritten. API-DIAGNOSTIC, tenant 4 fixture.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-06');
let cookie = '';
const call = async (method, path, body) => { const r = await fetch(API + path, { method, headers: { Origin: 'http://localhost:13112', Cookie: cookie, 'Content-Type': 'application/json' }, body: body ? JSON.stringify(body) : undefined }); const sc = r.headers.getSetCookie?.() ?? []; if (sc.length) cookie = sc.map((c) => c.split(';')[0]).join('; '); const t = await r.text(); let j; try { j = JSON.parse(t); } catch { j = t; } return { status: r.status, body: j }; };
const log = []; let sid = 0;
const snap = (step, extra = {}) => { const e = { step, supplier_current_balance: sql(`select current_balance from suppliers where id=${sid}`), invoices_balance_due_sum: sql(`select coalesce(sum(balance_due),0) from invoices_supplier where supplier_id=${sid} and status not in ('void','cancelled','draft')`), ap_2101_credit_balance_all_suppliers: sql("select coalesce(sum(l.credit_amount-l.debit_amount),0) from journal_entry_lines l join journal_entries j on j.id=l.journal_entry_id join chart_of_accounts a on a.id=l.account_id where j.tenant_id=4 and j.status='posted' and a.account_code='2101'"), ...extra }; log.push(e); console.log(JSON.stringify(e)); };
await call('POST', '/auth/login', { email: 'audit.c.owner@example.test', password: PASSWORD });
let r = await call('POST', '/suppliers', { company_name: 'Audit Supplier R5 ' + Date.now().toString().slice(-5), email: 'supr5@example.test', phone: '+201000000091' });
sid = r.body?.id ?? r.body?.data?.id ?? Number(sql("select id from suppliers where tenant_id=4 order by id desc limit 1")); snap('created supplier', { http: r.status, id: sid });
r = await call('POST', '/accounting/supplier-invoices', { supplier_ref: 'R5-' + Date.now().toString().slice(-6), invoice_date: '2026-11-26', due_date: '2026-12-26', supplier_id: sid, currency: 'EGP', status: 'draft', lines: [{ description: 'R5 goods', quantity: 10, unit_price: 1000, tax_rate: 0, discount_rate: 0, line_order: 1 }] });
const inv = r.body?.id ?? r.body?.data?.id; snap('invoice 10,000 created (draft)', { http: r.status, inv });
r = await call('POST', `/accounting/supplier-invoices/${inv}/approve`, {}); snap('invoice approved (AP posted)', { http: r.status });
r = await call('POST', '/accounting/payments/made', { payment_date: '2026-11-26', supplier_id: sid, invoice_id: inv, amount: 4000, payment_method: 'cash', reference: 'r5-1' }); snap('payment 4,000', { http: r.status });
r = await call('POST', '/accounting/payments/made', { payment_date: '2026-11-26', supplier_id: sid, invoice_id: inv, amount: 6000, payment_method: 'cash', reference: 'r5-2' }); snap('payment 6,000 (invoice fully paid)', { http: r.status });
const sup = (await call('GET', `/suppliers/${sid}`)).body; const sd = sup?.data ?? sup;
r = await call('PUT', `/suppliers/${sid}`, { ...sd, company_name: (sd.company_name || 'x') + ' edited', current_balance: 777 }); snap('PUT supplier (name edit, client sends current_balance 777)', { http: r.status, resp: JSON.stringify(r.body).slice(0, 160) });
r = await call('PUT', `/suppliers/${sid}`, { ...sd, company_name: (sd.company_name || 'x') + ' edited2' }); snap('PUT supplier again WITHOUT the field being changed by the user (re-sends the value it read)', { http: r.status });
fs.writeFileSync(`${d.db}/AUD-R5-F021-SUPPLIER-BALANCE.json`, JSON.stringify({ run: 'RUN-005', finding: 'F-021', supplier_id: sid, log }, null, 1));
