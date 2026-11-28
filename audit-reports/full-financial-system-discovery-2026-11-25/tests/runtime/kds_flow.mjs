import { login, EV, DASH, netLog } from './lib.mjs';
const { b, ctx } = await login('audit.c.owner@example.test', 'C');
const p = await ctx.newPage(); const net = []; netLog(p, net);
await p.goto(DASH + '/en/dashboard/kds', { waitUntil: 'networkidle' }); await p.waitForTimeout(3000);
const steps = [];
for (let i = 0; i < 4; i++) {
  const btn = p.locator('button').filter({ hasText: /\s*(->|→)\s*|Bump/ }).first();
  if (!(await btn.count())) break;
  const label = (await btn.innerText()).replace(/\s+/g, ' '); await btn.click(); await p.waitForTimeout(2500); steps.push(label);
}
console.log('clicked:', steps.join(' || '));
console.log(JSON.stringify(net.filter(x => x.m !== 'GET').map(x => `${x.m} ${x.u} ${x.s} req=${x.req.slice(0, 80)} ${x.body.slice(0, 120)}`), null, 1));
console.log((await p.locator('body').innerText()).split(String.fromCharCode(10)).join(' | ').slice(0, 300));
await b.close();
