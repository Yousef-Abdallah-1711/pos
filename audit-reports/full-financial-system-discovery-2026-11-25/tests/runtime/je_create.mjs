// usage: node je_create.mjs tag email id "description" date(YYYY-MM-DD) "code:dr:cr" "code:dr:cr" ...
import { login, EV, DASH, netLog } from './lib.mjs';
import fs from 'fs';
const [tag, email, id, desc, date, ...lines] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net);
await p.goto(DASH + '/en/dashboard/accounting/journal-entries', { waitUntil: 'networkidle' }); await p.waitForTimeout(2000);
await p.click('button:has-text("New journal entry")'); await p.waitForTimeout(1200);
const dlg = p.locator('form').last();
await p.locator('input[type=date]').last().fill(date);
await p.fill('input[placeholder="Record petty cash replenishment"]', desc);
while ((await p.locator('select:has(option:text-is("Select account"))').count()) < lines.length) { await p.click('button:has-text("Add line")'); await p.waitForTimeout(200); }
for (let i = 0; i < lines.length; i++) {
  const [code, dr, cr] = lines[i].split(':');
  const sel = p.locator('select:has(option:text-is("Select account"))').nth(i);
  const opt = await sel.locator('option').evaluateAll((os, c) => os.find(o => o.textContent.trim().startsWith(c + ' '))?.value, code);
  await sel.selectOption(opt);
  const nums = p.locator('input[type=number]');
  const base = await nums.count() - lines.length * 2; // line inputs are the last 2*n numbers
  if (dr && dr !== '0') await nums.nth(base + i * 2).fill(dr);
  if (cr && cr !== '0') await nums.nth(base + i * 2 + 1).fill(cr);
}
await p.screenshot({ path: `${EV}/screenshots/${id}-form.png` });
await p.click('button:has-text("Save journal entry")'); await p.waitForTimeout(3000);
let out = { id, desc, date, lines, writes: net.filter(n => n.m !== 'GET').map(n => `${n.m} ${n.u} ${n.s} ${n.body.slice(0, 220)}`) };
console.log(JSON.stringify(out, null, 1));
console.log((await p.locator('body').innerText()).replace(/\n+/g, ' | ').slice(0, 400));
await b.close();
