// RUN-003 / PH-09 / AUD-CINV-001: new customer invoice (qty 1 x 8,000, tax 14%, discount 10%) in HEADED Chrome, tenant 4 owner.
// Checks F-008 (journal posted while status draft?) and F-007 (UI-displayed tax/total vs stored). Independent oracle in integer piasters.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, recordMutations, sql, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-09');
const { ctx, page } = await openChrome();
const calls = []; recordMutations(page, calls);
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/accounting/customer-invoices`);
await page.waitForSelector('button:has-text("New customer invoice")', { timeout: 120000 });
const snap = () => sql("select 'invoices',count(*),coalesce(max(id),0) from invoices_customer where tenant_id=4 union all select 'journals',count(*),coalesce(max(id),0) from journal_entries where tenant_id=4 union all select 'gl_lines',count(*),0 from general_ledger where tenant_id=4");
const before = snap();
await page.click('button:has-text("New customer invoice")'); await page.waitForTimeout(1000);
const dates = page.locator('input[type=date]'); const n = await dates.count();
await dates.nth(n - 2).fill('2026-11-26'); await dates.nth(n - 1).fill('2026-11-26');
await page.locator('select').filter({ hasText: 'Select customer' }).first().selectOption({ index: 1 });
await page.fill('input[placeholder="Consulting services / Raw materials"]', 'AUD-CINV-001 draft check');
const nums = page.locator('input[type=number]'); const m = await nums.count();
await nums.nth(m - 4).fill('1'); await nums.nth(m - 3).fill('8000'); await nums.nth(m - 2).fill('14'); await nums.nth(m - 1).fill('10');
await page.waitForTimeout(800);
const formText = await page.locator('body').innerText();
const statusControl = await page.evaluate(() => [...document.querySelectorAll('select')].map((s) => ({ opts: [...s.options].map((o) => o.text).join('|'), val: s.value })));
await page.screenshot({ path: `${d.shot}/AUD-CINV-001-form-before-submit.png` });
await page.click('button:has-text("Create customer invoice")'); await page.waitForTimeout(4000);
await page.screenshot({ path: `${d.shot}/AUD-CINV-001-after-submit.png` });
const after = snap();
const inv = sql("select id,invoice_number,status,subtotal,tax_amount,total_amount from invoices_customer where tenant_id=4 order by id desc limit 1");
const je = sql("select id,entry_number,reference_type,reference_id,status,total_debit from journal_entries where tenant_id=4 order by id desc limit 1");
const lines = sql("select coa.account_code,gl.debit_amount,gl.credit_amount from general_ledger gl join chart_of_accounts coa on coa.id=gl.account_id where gl.tenant_id=4 and gl.journal_entry_id=(select max(id) from journal_entries where tenant_id=4) order by 1");
// Oracle (integer minor units): subtotal 800000; discount 10% => 80000; discount-first base 720000; VAT 14% on base = 100800; total 820800. VAT-before-discount: tax 112000, total 832000.
const oracle = { subtotal: 8000, discount: 800, base: 7200, tax_discount_first: 1008, total_discount_first: 8208, tax_before_discount: 1120 };
fs.writeFileSync(`${d.net}/AUD-CINV-001.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, role: 'owner', calls, formTextExcerpt: formText.split('\n').filter((l) => /tax|total|discount|subtotal|vat/i.test(l)).slice(0, 20), statusControl }, null, 1));
fs.writeFileSync(`${d.db}/AUD-CINV-001-before-after.txt`, `BEFORE\n${before}\nAFTER\n${after}\nINVOICE: ${inv}\nJOURNAL: ${je}\nLINES(account|debit|credit):\n${lines}\nORACLE(EGP): ${JSON.stringify(oracle)}\n`);
console.log({ before: before.split('\n'), after: after.split('\n'), inv, je, lines: lines.split('\n'), ui: formText.split('\n').filter((l) => /tax|total|discount|subtotal|vat/i.test(l)).slice(0, 12), statusControl, calls: calls.map((c) => `${c.method} ${c.url} ${c.status}`) });
await ctx.close();
