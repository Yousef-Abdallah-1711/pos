import { login, EV, DASH, netLog } from './lib.mjs';
const [tag, email, id, action] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net); p.on('dialog', d => d.accept());
await p.goto(DASH + '/en/dashboard/hr/payroll', { waitUntil: 'networkidle' }); await p.waitForTimeout(2000);
const d = p.locator('input[type=date]'); await d.nth(0).fill('2026-11-01'); await d.nth(1).fill('2026-11-30'); await p.waitForTimeout(2500);
console.log((await p.locator('body').innerText()).replace(/\n+/g, ' | ').slice(150, 800));
const btns = await p.$$eval('button', els => els.filter(e => e.offsetParent !== null).map(e => (e.innerText || e.title || e.getAttribute('aria-label') || '').trim().slice(0, 25)).filter(Boolean)); console.log('BTN', btns.join(' || '));
if (action) { await p.click(`button:has-text("${action}")`); await p.waitForTimeout(1500); const c = p.locator('button:has-text("Confirm")'); if (await c.count()) { await c.last().click(); await p.waitForTimeout(2500); } }
console.log(JSON.stringify(net.filter(x => x.m !== 'GET').map(x => `${x.m} ${x.u} ${x.s} resp=${x.body.slice(0, 250)}`), null, 1));
await p.screenshot({ path: `${EV}/screenshots/${id}.png` });
await b.close();
