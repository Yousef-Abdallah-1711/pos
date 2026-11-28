// RUN-003 / PH-07 / AUD-POS-001: cashier-style cash order (2 x Audit Burger) in HEADED Chrome, tenant 4 owner. Compares UI totals, API response and stored row; checks journal + stock side effects.
import fs from 'fs';
import { openChrome, ensureLoggedIn, evidenceDirs, recordMutations, sql, DASH } from './chrome_lib.mjs';
const d = evidenceDirs('PH-07');
const { ctx, page } = await openChrome();
const shotFailures = [];
const shot = async (path) => { try { await page.bringToFront(); await page.screenshot({ path, timeout: 25000 }); } catch (e) { shotFailures.push(path + ': ' + String(e).slice(0, 80)); } };
const calls = []; recordMutations(page, calls);
await ensureLoggedIn(page, 'audit.c.owner@example.test', `${DASH}/en/dashboard/pos`);
await page.waitForSelector('text=Audit Burger', { timeout: 120000 }); await page.waitForTimeout(2000);
const snap = () => sql("select 'orders',count(*),coalesce(max(id),0) from orders where tenant_id=4 union all select 'journals',count(*),coalesce(max(id),0) from journal_entries where tenant_id=4 union all select 'stock',quantity_in_stock::int,0 from products where id=2 union all select 'outbox',count(*),0 from accounting_inventory_movement_events where tenant_id=4");
const before = snap();
for (let i = 0; i < 2; i++) { await page.click('text=Audit Burger'); await page.waitForTimeout(400); }
const body = await page.locator('body').innerText();
const ui = body.split('\n').filter((l) => /subtotal|tax|total|vat/i.test(l)).slice(0, 10);
await page.click('button:has-text("Cash")').catch(() => {}); await page.waitForTimeout(300);
await shot(`${d.shot}/AUD-POS-001-before-submit.png`);
await page.click('button:has-text("Place Order")'); await page.waitForTimeout(6000);
await shot(`${d.shot}/AUD-POS-001-after-submit.png`);
const after = snap();
const order = sql("select id,status,payment_status,payment_method,subtotal,tax_amount,total_amount from orders where tenant_id=4 order by id desc limit 1");
const je = sql("select id,entry_number,reference_type,status,total_debit from journal_entries where tenant_id=4 order by id desc limit 1");
const lines = sql("select coa.account_code,gl.debit_amount,gl.credit_amount from general_ledger gl join chart_of_accounts coa on coa.id=gl.account_id where gl.tenant_id=4 and gl.journal_entry_id=(select max(id) from journal_entries where tenant_id=4) order by 1");
// Oracle (EGP, integer minor units): 2 x 10,000 = 20,000 subtotal. Server rate observed 10% => tax 2,000 total 22,000; UI hard-coded 14% => tax 2,800 total 22,800.
fs.writeFileSync(`${d.net}/AUD-POS-001.json`, JSON.stringify({ run: 'RUN-003', tenant: 4, role: 'owner', ui, calls }, null, 1));
fs.writeFileSync(`${d.db}/AUD-POS-001-before-after.txt`, `BEFORE\n${before}\nAFTER\n${after}\nORDER: ${order}\nLATEST JOURNAL: ${je}\nLINES:\n${lines}\nORACLE: subtotal 20000; tax10%=2000 total 22000; tax14%=2800 total 22800\n`);
console.log({ shotFailures, ui, before: before.split('\n'), after: after.split('\n'), order, je, lines: lines.split('\n'), calls: calls.map((c) => `${c.method} ${c.url} ${c.status} ${c.response.slice(0, 160)}`) });
await ctx.close();
