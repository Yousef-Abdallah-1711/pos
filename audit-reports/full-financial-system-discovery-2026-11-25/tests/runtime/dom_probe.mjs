import { login, DASH } from './lib.mjs';
const [tag, email, path, text] = process.argv.slice(2);
const { b, ctx } = await login(email, tag);
const p = await ctx.newPage();
await p.goto(DASH + '/' + path, { waitUntil: 'networkidle' }); await p.waitForTimeout(2500);
const h = await p.evaluate(t => { const el = [...document.querySelectorAll('*')].filter(e => e.children.length === 0 && e.textContent.includes(t))[0]; if (!el) return 'NOT FOUND'; let r = el; for (let i = 0; i < 3 && r.parentElement; i++) r = r.parentElement; return r.outerHTML.slice(0, 1500); }, text);
console.log(h);
await b.close();
