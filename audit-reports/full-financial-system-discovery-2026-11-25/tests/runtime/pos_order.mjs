import { login, EV, DASH, netLog } from './lib.mjs';
import fs from 'fs';
const [tag, email, product, method, qty, id] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net);
await p.goto(DASH + '/en/dashboard/pos', { waitUntil: 'networkidle' }); await p.waitForTimeout(2000);
for (let i = 0; i < Number(qty || 1); i++) { await p.click(`text=${product}`); await p.waitForTimeout(300); }
const shown = (await p.locator('body').innerText());
const grab = k => (shown.match(new RegExp(k + '\s*\n?\s*([^\n]+)')) || [])[1];
const display = { subtotal: grab('Subtotal'), tax: grab('Tax \(14%\)') || grab('Tax[^\n]*'), total: grab('Total') };
await p.click(`button:has-text("${method}")`); await p.waitForTimeout(300);
await p.screenshot({ path: `${EV}/screenshots/pos-${id}-before-submit.png` });
const t0 = Date.now();
await p.click('button:has-text("Place Order")'); await p.waitForTimeout(5000);
const out = { id, tag, product, method, qty, display, ui_after: (await p.locator('body').innerText()).replace(/\n+/g, ' | ').slice(0, 250), net: net.filter(n => n.m !== 'GET') };
await p.screenshot({ path: `${EV}/screenshots/pos-${id}-after-submit.png` });
fs.writeFileSync(`${EV}/network/pos-${id}.json`, JSON.stringify(out, null, 1));
console.log(JSON.stringify({ display, writes: out.net.map(n => `${n.m} ${n.u} ${n.s} ${n.body.slice(0, 200)}`) }, null, 1));
await b.close();
