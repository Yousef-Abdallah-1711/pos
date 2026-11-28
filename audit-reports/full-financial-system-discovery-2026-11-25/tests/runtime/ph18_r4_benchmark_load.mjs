// RUN-004 / PH-18: builds a dedicated SYNTHETIC benchmark tenant (Audit G) with >= 50,000 POSTED journal lines via the product's own CSV import with auto_post=true.
// Fixture label: SYNTHETIC BENCHMARK (not an end-to-end user transaction set). COA copied by SQL from tenant 4 (F-001 workaround, same technique as tenants 8/9).
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112';
const EMAIL = 'audit.g.owner@example.test';
const call = async (m, p, b, ck) => { const x = await fetch(API + p, { method: m, headers: { Origin: O, 'Content-Type': 'application/json', ...(ck ? { Cookie: ck } : {}) }, body: b ? JSON.stringify(b) : undefined }); const t = await x.text(); let j; try { j = JSON.parse(t); } catch { j = t.slice(0, 200); } return { status: x.status, body: j, headers: x.headers }; };
let tid = sql(`select tenant_id from users where email='${EMAIL}'`);
if (!tid) { const r = await call('POST', '/auth/register', { name: 'Audit Owner G', email: EMAIL, password: PASSWORD, restaurant_name: 'Audit G Benchmark', restaurant_slug: 'audit-g-bench', phone: '+201000000007', billing_period: 'monthly' }); console.log('register', r.status); tid = sql(`select tenant_id from users where email='${EMAIL}'`); }
const rid = sql(`select restaurant_id from users where email='${EMAIL}'`);
console.log('tenant', tid, 'restaurant', rid);
if (sql(`select count(*) from chart_of_accounts where tenant_id=${tid}`) === '0') {
  const COLS = sql("select string_agg(column_name,',' order by ordinal_position) from information_schema.columns where table_name='chart_of_accounts' and column_name not in ('id','tenant_id','restaurant_id','parent_id','created_at','updated_at','created_by','updated_by','current_balance')");
  sql(`insert into chart_of_accounts (tenant_id,restaurant_id,parent_id,${COLS}) select ${tid},${rid},null,${COLS} from chart_of_accounts where tenant_id=4`);
}
const lg = await call('POST', '/auth/login', { email: EMAIL, password: PASSWORD });
const ck = lg.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const up = async (path, csv) => { const fd = new FormData(); fd.append('file', new Blob([csv], { type: 'text/csv' }), 'import.csv'); const x = await fetch(API + path, { method: 'POST', headers: { Origin: O, Cookie: ck }, body: fd, signal: AbortSignal.timeout(400000) }); return { status: x.status, body: (await x.text()).slice(0, 200) }; };
const H = 'entry_key,entry_date,description,entry_type,account_code,line_description,debit_amount,credit_amount\n';
const cents = (i) => 1000 + ((i * 7919) % 99000); // 10.00 .. 999.99 in cents
const date = (i) => { const day = i % 330; const dt = new Date(Date.UTC(2026, 0, 1) + day * 86400000); return dt.toISOString().slice(0, 10); };
const fmt = (c) => (c / 100).toFixed(2);
const TOTAL = Number(process.env.BENCH_JOURNALS || 27000), CHUNK = 1500;
const oracle = { sales_cents: 0, expense_cents: 0, journals: 0 };
const have = Number(sql(`select count(*) from journal_entries where tenant_id=${tid}`));
let start = have; console.log('already', have);
const t0 = performance.now(); const timings = [];
for (let s = start; s < TOTAL; s += CHUNK) {
  const rows = [];
  for (let i = s; i < Math.min(s + CHUNK, TOTAL); i++) {
    const c = cents(i); const sale = i % 5 !== 4;
    rows.push(`G${i},${date(i)},Bench ${i},manual,${sale ? 1110 : 5110},l,${fmt(c)},0`); rows.push(`G${i},${date(i)},Bench ${i},manual,${sale ? 4100 : 1110},l,0,${fmt(c)}`);
  }
  const ts = performance.now(); let r; try { r = await up('/accounting/csv/journal-entries/import?auto_post=true', H + rows.join('\n') + '\n'); } catch (e) { r = { status: 'ERR', body: String(e.name) }; }
  timings.push({ from: s, status: r.status, resp: r.body, ms: Math.round(performance.now() - ts) }); console.log(s, r.status, r.body, Math.round(performance.now() - ts), 'ms');
  if (r.status === 'ERR') { await new Promise((x) => setTimeout(x, 20000)); }
}
for (let i = 0; i < TOTAL; i++) { const c = cents(i); if (i % 5 !== 4) oracle.sales_cents += c; else oracle.expense_cents += c; oracle.journals++; }
const st = { tenant: tid, journals: sql(`select count(*) from journal_entries where tenant_id=${tid}`), posted: sql(`select count(*) from journal_entries where tenant_id=${tid} and status='posted'`), lines: sql(`select count(*) from journal_entry_lines where tenant_id=${tid}`), gl_rows: sql(`select count(*) from general_ledger where tenant_id=${tid}`) };
fs.writeFileSync(`${d.net}/AUD-R4-BENCH-LOAD.json`, JSON.stringify({ run: 'RUN-004', label: 'SYNTHETIC BENCHMARK FIXTURE', oracle: { sales: fmt(oracle.sales_cents), expense: fmt(oracle.expense_cents), journals: oracle.journals }, totalMs: Math.round(performance.now() - t0), timings, st }, null, 1));
console.log(JSON.stringify(st), 'oracle sales', fmt(oracle.sales_cents), 'expense', fmt(oracle.expense_cents));
