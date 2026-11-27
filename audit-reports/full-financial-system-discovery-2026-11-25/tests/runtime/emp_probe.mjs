import { login, DASH } from './lib.mjs';
const { b, ctx } = await login('audit.c.owner@example.test', 'C');
const p = await ctx.newPage();
await p.goto(DASH + '/en/dashboard/hr/employees', { waitUntil: 'networkidle' }); await p.waitForTimeout(2000);
await p.click('button:has-text("Add Employee")'); await p.waitForTimeout(1000);
for (const t of ['Employment', 'Compensation', 'Contact']) {
  await p.click(`button:has-text("${t}")`); await p.waitForTimeout(600);
  const f = await p.$$eval('input,select,textarea', els => els.filter(e => e.offsetParent !== null).map(e => `${e.tagName}[${e.type}] name=${e.name} ph=${e.placeholder} opts=${e.tagName === 'SELECT' ? [...e.options].map(o => o.text).slice(0, 8).join('/') : ''}`));
  console.log('== ' + t + '\n' + f.join('\n'));
}
await b.close();
