// usage: node je_post.mjs tag email id entry_number [action=Post]
import { login, EV, DASH, netLog } from './lib.mjs';
const [tag, email, id, num, action = 'Post'] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net);
p.on('dialog', d => d.accept());
await p.goto(DASH + '/en/dashboard/accounting/journal-entries', { waitUntil: 'networkidle' }); await p.waitForTimeout(2000);
await p.locator(`tr:has-text("${num}") >> text=Open`).click(); await p.waitForTimeout(1500);
const btns = await p.$$eval('button', els => els.map(e => (e.innerText || '').trim().slice(0, 25)).filter(Boolean)); console.log('BTN', btns.join(' || '));
await p.click(`button:has-text("${action}")`); await p.waitForTimeout(1500);
const conf = p.locator('button:has-text("Confirm")'); if (await conf.count()) { await conf.first().click(); await p.waitForTimeout(2500); }
console.log(JSON.stringify(net.filter(n => n.m !== 'GET').map(n => `${n.m} ${n.u} ${n.s} ${n.body.slice(0, 160)}`), null, 1));
await p.screenshot({ path: `${EV}/screenshots/${id}-${action}.png` });
await b.close();
