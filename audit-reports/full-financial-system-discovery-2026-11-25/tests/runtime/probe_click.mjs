import { login, EV, DASH, netLog } from './lib.mjs';
const [tag, email, path, clickText] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage(); const net = []; netLog(p, net);
await p.goto(DASH + '/' + path, { waitUntil: 'networkidle' }); await p.waitForTimeout(2500);
if (clickText) { await p.click(`text=${clickText}`).catch(e => console.log('click fail', e.message.slice(0, 60))); await p.waitForTimeout(1500); }
console.log((await p.locator('body').innerText()).replace(/\n+/g, ' | ').slice(0, 1300));
const f = await p.$$eval('input,select,textarea,button', els => els.filter(e => e.offsetParent !== null).map(e => `${e.tagName}[${e.type}] name=${e.name} ph=${e.placeholder} txt=${(e.innerText || '').slice(0, 25).replace(/\n/g, ' ')}`)); console.log(f.slice(8).join('\n'));
await b.close();
