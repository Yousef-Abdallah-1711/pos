// RUN-003 / PH-16 / AUD-MD-NEGPRICE: can a negative-price product be sold, and what happens to totals/GL? Tenant 4 (products 8 price -5, 9 price 0, 15 price 10.12345).
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-16'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const A = async (m, p, b) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: x.status, body: j }; };
const out = [];
const base = { order_source: 'pos', customer_name: 'NP', customer_phone: '0000000000', payment_method: 'cash' };
for (const [tag, pid, q] of [['negative-price', 8, 1], ['zero-price', 9, 1], ['five-decimal-price', 15, 3]]) {
  const r = await A('POST', '/pos/orders', { ...base, items: [{ product_id: pid, quantity: q }] }); const id = r.body?.data?.id;
  const row = id ? sql(`select subtotal||'|'||tax_amount||'|'||total_amount from orders where id=${id}`) : null;
  let pay = null, je = null;
  if (id) { pay = (await A('PUT', `/admin/orders/${id}/payment-status`, { payment_status: 'paid' })).status; await new Promise((x) => setTimeout(x, 4000)); je = sql(`select string_agg(a.account_code||':'||l.debit_amount||'/'||l.credit_amount,' ') from journal_entries j join journal_entry_lines l on l.journal_entry_id=j.id join chart_of_accounts a on a.id=l.account_id where j.tenant_id=4 and j.reference_type='order' and j.reference_id=${id}`); }
  out.push({ tag, pid, q, http: r.status, resp: id ? undefined : JSON.stringify(r.body).slice(0, 160), order: id, totals: row, payHttp: pay, journal: je });
  console.log(JSON.stringify(out[out.length - 1]));
}
fs.writeFileSync(`${d.db}/AUD-MD-NEGPRICE.json`, JSON.stringify({ run: 'RUN-003', out }, null, 1));
