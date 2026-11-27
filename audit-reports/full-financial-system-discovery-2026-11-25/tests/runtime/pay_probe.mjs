import { login, EV, DASH, netLog } from './lib.mjs';
const [tag, email, path, btn] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage();
await p.goto(DASH + '/' + path, { waitUntil: 'networkidle' }); await p.waitForTimeout(2500);
await p.click(`button:has-text("${btn}")`); await p.waitForTimeout(1500);
const f = await p.$$eval('input,select,textarea,button', els => els.filter(e => e.offsetParent !== null).map(e => `${e.tagName}[${e.type}] ph=${e.placeholder} txt=${(e.innerText || '').slice(0, 60).replace(/\n/g, ' / ')}`)); console.log(f.slice(12).join('\n'));
await b.close();
