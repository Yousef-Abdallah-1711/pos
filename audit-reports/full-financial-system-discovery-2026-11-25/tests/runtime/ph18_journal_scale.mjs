// RUN-003 / PH-18: list/export latency with ~17,000 draft journals in tenant 4 (created by the import bulk tests). API-DIAGNOSTIC.
import fs from 'fs';
import { sql, API, PASSWORD, evidenceDirs } from './chrome_lib.mjs';
const d = evidenceDirs('PH-18'); const O = 'http://localhost:13112';
const r0 = await fetch(API + '/auth/login', { method: 'POST', headers: { Origin: O, 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'audit.c.owner@example.test', password: PASSWORD }) });
const ck = r0.headers.getSetCookie().map((c) => c.split(';')[0]).join('; ');
const out = { journals: sql("select count(*) from journal_entries where tenant_id=4"), lines: sql("select count(*) from journal_entry_lines where tenant_id=4") };
for (const p of ['/accounting/journal-entries', '/accounting/journal-entries?page=1&limit=50', '/accounting/csv/journal-entries/export', '/accounting/trial-balance', '/reports/financial/profit-loss?date_from=2026-01-01&date_to=2026-12-31']) {
  const s = performance.now(); let st, len; try { const x = await fetch(API + p, { headers: { Origin: O, Cookie: ck }, signal: AbortSignal.timeout(120000) }); st = x.status; len = (await x.text()).length; } catch (e) { st = 'ERR ' + String(e.name); }
  out[p] = { status: st, bytes: len, ms: Math.round(performance.now() - s) }; console.log(p, JSON.stringify(out[p]));
}
fs.writeFileSync(`${d.db}/AUD-JOURNAL-SCALE.json`, JSON.stringify({ run: 'RUN-003', out }, null, 1));
console.log(out.journals, out.lines);
