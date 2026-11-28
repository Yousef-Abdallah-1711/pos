// RUN-003 / PH-19: HC-04 (Balance Sheet residual) and HC-05 (journal attachment validation) in HEADED Chrome, tenant 4 owner.
// HC-05 creates DRAFT diagnostic journals on tenant 4 if the API accepts them (recorded, side effect declared).
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, sql, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-19');
const { ctx, page } = await openChrome();
const bs = {};
page.on('response', async (r) => { if (/balance-sheet/.test(r.url()) && r.request().method() === 'GET') { try { bs[r.url().replace('http://localhost:18082/api/v1', '')] = await r.text(); } catch {} } });
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/reports/financial/balance-sheet`);
await page.waitForTimeout(8000);
await page.screenshot({ path: `${d.shot}/AUD-REG-HC04-balance-sheet.png`, fullPage: true });
const bsText = (await page.locator('body').innerText()).split('\n').filter((l) => /asset|liabilit|equity|balanc|differ|total/i.test(l)).slice(0, 25);
const jeBefore = sql("select count(*) from journal_entries where tenant_id=4");
const line = (a, dr, cr) => ({ account_id: a, debit_amount: dr, credit_amount: cr, description: 'AUD-REG-HC05', description_ar: '', line_order: 1, tax_code: '', tax_amount: 0 });
const variants = {
  'v1-wellformed-array': [{ name: 'a.pdf', url: 'https://example.test/a.pdf', size: 10 }],
  'v2-object-not-array': { evil: true },
  'v3-plain-string-xss': '<script>alert(1)</script>',
  'v4-javascript-url': [{ name: 'x', url: 'javascript:alert(1)' }],
  'v5-1MB-blob': [{ name: 'big', data: 'A'.repeat(1024 * 1024) }],
};
const res = await page.evaluate(async (variants) => {
  const out = {};
  for (const [k, att] of Object.entries(variants)) {
    const body = { entry_date: '2026-11-26', description: 'AUD-REG-HC05 ' + k, notes: '', attachments: att, lines: [{ account_id: 107, debit_amount: 1, credit_amount: 0, description: 'AUD-REG-HC05', description_ar: '', line_order: 1, tax_code: '', tax_amount: 0 }, { account_id: 71, debit_amount: 0, credit_amount: 1, description: 'AUD-REG-HC05', description_ar: '', line_order: 2, tax_code: '', tax_amount: 0 }] };
    try { const r = await fetch('http://localhost:18082/api/v1/accounting/journal-entries', { method: 'POST', credentials: 'include', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) }); out[k] = { status: r.status, resp: (await r.text()).slice(0, 220) }; } catch (e) { out[k] = { error: String(e) }; }
  }
  return out;
}, variants);
const jeAfter = sql("select count(*) from journal_entries where tenant_id=4");
const stored = sql("select id,status,left(attachments::text,60),length(attachments::text) from journal_entries where tenant_id=4 and description like 'AUD-REG-HC05%' order by id");
fs.writeFileSync(`${d.net}/AUD-REG-HC04-HC05.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, bsApi: bs, bsVisible: bsText, hc05: res }, null, 1));
fs.writeFileSync(`${d.db}/AUD-REG-HC05-before-after.txt`, `journals before=${jeBefore} after=${jeAfter}\nstored(id|status|attachments_prefix|len):\n${stored}\n`);
console.log(JSON.stringify({ bsVisible: bsText, hc05: Object.fromEntries(Object.entries(res).map(([k, v]) => [k, `${v.status ?? v.error} ${String(v.resp).slice(0, 90)}`])), jeBefore, jeAfter, stored: stored.split('\n') }, null, 1));
await ctx.close();
