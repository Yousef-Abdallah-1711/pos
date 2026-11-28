// usage: node supinv_create.mjs tag email id "supplier" "desc" accountCode|default qty unit tax% date
import { login, EV, DASH, netLog } from './lib.mjs';
const [tag, email, id, supplier, desc, acct, qty, unit, tax, date] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net);
await p.goto(DASH + '/en/dashboard/accounting/supplier-invoices', { waitUntil: 'networkidle' }); await p.waitForTimeout(2000);
await p.click('button:has-text("New supplier invoice")'); await p.waitForTimeout(1000);
const dates = p.locator('input[type=date]'); const n = await dates.count();
await dates.nth(n - 2).fill(date); await dates.nth(n - 1).fill(date);
await p.locator('select:has(option:text-is("Select supplier"))').selectOption({ label: supplier });
await p.fill('input[placeholder="Consulting services / Raw materials"]', desc);
if (acct !== 'default') { const s = p.locator('select:has(option:text-is("Use system default"))'); const v = await s.locator('option').evaluateAll((os, c) => os.find(o => o.textContent.trim().startsWith(c + ' '))?.value, acct); await s.selectOption(v); }
const nums = p.locator('input[type=number]'); const m = await nums.count();
await nums.nth(m - 4).fill(qty); await nums.nth(m - 3).fill(unit); await nums.nth(m - 2).fill(tax);
await p.screenshot({ path: `${EV}/screenshots/${id}-form.png` });
await p.click('button:has-text("Create supplier invoice")'); await p.waitForTimeout(3500);
console.log(JSON.stringify(net.filter(x => x.m !== 'GET').map(x => `${x.m} ${x.u} ${x.s} req=${x.req.slice(0, 350)} resp=${x.body.slice(0, 380)}`), null, 1));
await b.close();
