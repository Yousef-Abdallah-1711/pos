// usage: node row_action.mjs tag email id path rowText [actionText]   (no actionText => list buttons in row)
import { login, EV, DASH, netLog } from './lib.mjs';
const [tag, email, id, path, rowText, action] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net); p.on('dialog', d => d.accept());
await p.goto(DASH + '/' + path, { waitUntil: 'networkidle' }); await p.waitForTimeout(2500);
const row = p.locator(`tr:has-text("${rowText}")`).first();
const rb = await row.locator('button, a').evaluateAll(els => els.map(e => (e.innerText || e.getAttribute('aria-label') || e.title || '').trim()));
console.log('ROW BUTTONS', JSON.stringify(rb));
if (action) { await row.locator(`text=${action}`).first().click(); await p.waitForTimeout(1500); const c = p.locator('button:has-text("Confirm")'); if (await c.count()) { await c.first().click(); await p.waitForTimeout(2000); } }
console.log(JSON.stringify(net.filter(x => x.m !== 'GET').map(x => `${x.m} ${x.u} ${x.s} req=${x.req.slice(0, 200)} resp=${x.body.slice(0, 300)}`), null, 1));
await p.screenshot({ path: `${EV}/screenshots/${id}.png` });
await b.close();
