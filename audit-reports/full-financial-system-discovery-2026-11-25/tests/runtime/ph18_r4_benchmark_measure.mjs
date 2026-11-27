// RUN-004 / PH-18: measures reports on the SYNTHETIC BENCHMARK tenant (Audit G, >= 50,000 posted lines): correctness vs an independent integer-cents oracle, API latency (3 runs), payload size, EXPLAIN ANALYZE of the P&L aggregation, index inventory.
// Concurrent load (other audit Chrome sessions) is recorded; not a production capacity claim.
import fs from 'fs';
import { execSync } from 'child_process';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112';
const EMAIL = 'audit.g.owner@example.test'; const TID = sql(`select tenant_id from users where email='${EMAIL}'`);
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: EMAIL, password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const sizes = { journals: sql(`select count(*) from journal_entries where tenant_id=${TID}`), posted: sql(`select count(*) from journal_entries where tenant_id=${TID} and status='posted'`), lines: sql(`select count(*) from journal_entry_lines where tenant_id=${TID}`), gl_rows: sql(`select count(*) from general_ledger where tenant_id=${TID}`) };
const load = { concurrent_audit_node_processes: execSync('powershell -NoProfile -Command "(Get-CimInstance Win32_Process -Filter \\"Name=\'node.exe\'\\" | Where-Object { $_.CommandLine -match \'ph0|ph1\' }).Count"', { encoding: 'utf8' }).trim(), free_mb: execSync('powershell -NoProfile -Command "[math]::Round((Get-CimInstance Win32_OperatingSystem).FreePhysicalMemory/1024)"', { encoding: 'utf8' }).trim() };
// oracle from the generation rule used by ph18_r4_benchmark_load.mjs (integer cents)
const N = Number(sizes.journals); const cents = (i) => 1000 + ((i * 7919) % 99000);
let sales = 0, expense = 0; for (let i = 0; i < Math.min(N, 27000); i++) { if (i % 5 !== 4) sales += cents(i); else expense += cents(i); }
const oracle = { revenue: (sales / 100).toFixed(2), expense: (expense / 100).toFixed(2), net: ((sales - expense) / 100).toFixed(2), cash_1110: ((sales - expense) / 100).toFixed(2), note: 'valid only while journals == generated count (27,000)' };
const eps = {
  trial_balance: '/accounting/trial-balance?as_of=2026-12-31',
  profit_loss: '/reports/financial/profit-loss?date_from=2026-01-01&date_to=2026-12-31',
  balance_sheet: '/reports/financial/balance-sheet?as_of=2026-12-31',
  cash_flow: '/reports/financial/cash-flow?date_from=2026-01-01&date_to=2026-12-31',
  ledger_1110_page: `/accounting/ledger?account_id=${sql(`select id from chart_of_accounts where tenant_id=${TID} and account_code='1110'`)}&start=2026-01-01&end=2026-12-31&limit=100`,
  journals_first_page: '/accounting/journal-entries?page=1&limit=50',
  journals_export_csv: '/accounting/csv/journal-entries/export',
};
const results = {}; const bodies = {};
for (const [k, p] of Object.entries(eps)) {
  const runs = [];
  for (let i = 0; i < 3; i++) { const s = performance.now(); let st, len, body; try { const x = await fetch(API + p, { headers: { Origin: O, Cookie: ck }, signal: AbortSignal.timeout(240000) }); st = x.status; body = await x.text(); len = body.length; } catch (e) { st = 'ERR ' + e.name; len = 0; } runs.push({ st, ms: Math.round(performance.now() - s), bytes: len }); if (i === 0 && body && k !== 'journals_export_csv') bodies[k] = body.slice(0, 200000); }
  results[k] = runs; console.log(k, JSON.stringify(runs));
}
const num = (o, keys) => keys.map((k) => o?.[k]);
let correctness = {};
try { const pl = JSON.parse(bodies.profit_loss).data ?? JSON.parse(bodies.profit_loss); correctness.pl = { total_revenue: pl.total_revenue, total_expenses: pl.total_expenses, net_income: pl.net_income ?? pl.net_profit, oracle_revenue: oracle.revenue, oracle_expense: oracle.expense }; } catch (e) { correctness.pl = 'parse error ' + e.message; }
try { const tb = JSON.parse(bodies.trial_balance); const dd = tb.data ?? tb; correctness.tb = { keys: Object.keys(dd).slice(0, 8), total_debit: dd.total_debit, total_credit: dd.total_credit }; } catch (e) { correctness.tb = 'parse error ' + e.message; }
correctness.sql_cash_1110 = sql(`select round(sum(g.debit_amount-g.credit_amount),2) from general_ledger g join chart_of_accounts a on a.id=g.account_id where g.tenant_id=${TID} and a.account_code='1110'`);
correctness.sql_revenue_4100 = sql(`select round(sum(g.credit_amount-g.debit_amount),2) from general_ledger g join chart_of_accounts a on a.id=g.account_id where g.tenant_id=${TID} and a.account_code='4100'`);
correctness.sql_unbalanced_journals = sql(`select count(*) from (select journal_entry_id from journal_entry_lines where tenant_id=${TID} group by 1 having round(sum(debit_amount)-sum(credit_amount),2)<>0) x`);
correctness.duplicate_entry_numbers = sql(`select count(*) from (select entry_number from journal_entries where tenant_id=${TID} group by 1 having count(*)>1) x`);
const explain = sql(`explain (analyze, buffers, timing off) select a.account_type, sum(l.credit_amount-l.debit_amount) from journal_entry_lines l join journal_entries j on j.id=l.journal_entry_id join chart_of_accounts a on a.id=l.account_id where j.tenant_id=${TID} and j.status='posted' and j.entry_date between '2026-01-01' and '2026-12-31' group by 1`);
const idx = sql("select tablename||'.'||indexname||' : '||indexdef from pg_indexes where tablename in ('journal_entries','journal_entry_lines','general_ledger') order by 1").split('\n').map((l) => l.slice(0, 200));
fs.writeFileSync(`${d.net}/AUD-R4-BENCH-MEASURE.json`, JSON.stringify({ run: 'RUN-004', label: 'SYNTHETIC BENCHMARK FIXTURE', tenant: TID, sizes, load, oracle, results, correctness, explain: explain.split('\n'), indexes: idx }, null, 1));
console.log(JSON.stringify({ sizes, load, correctness }, null, 1)); console.log(explain.split('\n').slice(0, 14).join('\n'));
