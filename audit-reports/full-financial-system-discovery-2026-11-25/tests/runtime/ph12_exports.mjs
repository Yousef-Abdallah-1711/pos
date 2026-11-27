// RUN-003 / PH-12: report export parity (sales, profitability, inventory; csv and pdf; EN/AR) vs the report API and SQL; invoice PDF. Tenant 4. API-DIAGNOSTIC.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-12'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const post = async (b) => { const x = await fetch(API + '/reports/export', { method: 'POST', headers: { Origin: O, Cookie: ck, 'Content-Type': 'application/json' }, body: JSON.stringify(b) }); const buf = Buffer.from(await x.arrayBuffer()); return { status: x.status, ctype: x.headers.get('content-type'), disp: x.headers.get('content-disposition'), buf }; };
const log = []; const rec = (id, desc, extra = {}) => { log.push({ id, desc, ...extra }); console.log(id, desc, JSON.stringify(extra).slice(0, 420)); };
const Q = { date_from: '2026-10-28', date_to: '2026-11-26', granularity: 'day', top_limit: 10, include_items: true };
const api = async (p) => (await (await fetch(API + p, { headers: { Origin: O, Cookie: ck } })).json()).data;
const sales = await api('/reports/sales?date_from=2026-10-28&date_to=2026-11-26&granularity=day'); const prof = await api('/reports/profitability?date_from=2026-10-28&date_to=2026-11-26&granularity=day'); const inv = await api('/reports/inventory?date_from=2026-10-28&date_to=2026-11-26&granularity=day');
for (const [type, api_total] of [['sales', sales.summary.total_revenue], ['profitability', prof.summary.total_revenue], ['inventory', inv.summary.total_current_inventory_value]]) {
  for (const [fmt, loc] of [['csv', 'en'], ['csv', 'ar'], ['pdf', 'en']]) {
    const r = await post({ ...Q, report_type: type, format: fmt, locale: loc });
    const txt = fmt === 'csv' ? r.buf.toString('utf8') : null;
    const nums = txt ? (txt.match(/\d{1,3}(?:,\d{3})*(?:\.\d+)?|\d+(?:\.\d+)?/g) || []).map((x) => Number(x.replace(/,/g, ''))) : [];
    rec(`EXP-${type}-${fmt}-${loc}`, `${type} export`, { http: r.status, ctype: r.ctype, disp: r.disp, bytes: r.buf.length, apiTotal: api_total, totalInFile: txt ? nums.includes(Number(api_total)) : 'pdf-not-parsed', pdfMagic: fmt === 'pdf' ? r.buf.slice(0, 5).toString() : undefined, bom: txt ? txt.charCodeAt(0) === 0xFEFF : undefined, head: txt ? txt.split('\n').slice(0, 3).join(' // ').slice(0, 200) : undefined });
    if (fmt === 'csv' && loc === 'en') fs.writeFileSync(`${d.net}/AUD-EXPORT-${type}.csv`, txt.slice(0, 20000));
  }
}
let r = await post({ ...Q, report_type: 'balance_sheet', format: 'csv' }); rec('EXP-unsupported', 'unsupported report type (financial statements)', { http: r.status, resp: r.buf.toString().slice(0, 120) });
const invId = Number(sql("select max(id) from invoices_customer where tenant_id=4"));
const pdf = await fetch(API + `/accounting/customer-invoices/${invId}/pdf`, { headers: { Origin: O, Cookie: ck } }); const pb = Buffer.from(await pdf.arrayBuffer()); rec('EXP-INV-PDF', 'customer invoice PDF', { http: pdf.status, ctype: pdf.headers.get('content-type'), bytes: pb.length, magic: pb.slice(0, 5).toString(), totalInDb: sql(`select total_amount from invoices_customer where id=${invId}`) });
fs.writeFileSync(`${d.db}/AUD-EXPORTS.json`, JSON.stringify({ run: 'RUN-003', kind: 'API-DIAGNOSTIC', log }, null, 1));
